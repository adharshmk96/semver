package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// VersionFile holds the version for projects that are not git repositories.
const VersionFile = ".version"

// Source tells where the current version was read from.
type Source int

const (
	SourceNone Source = iota
	SourceGit
	SourceFile
)

func (s Source) String() string {
	switch s {
	case SourceGit:
		return "Source: Git"
	case SourceFile:
		return "Source: File"
	default:
		return "semver config not found. run `semver init` to initialize the semver configuration."
	}
}

// Context is the state of the project in the current working directory.
type Context struct {
	Version   *Semver
	Source    Source
	IsGitRepo bool
	DryRun    bool
}

// BuildContext moves to the repository root (when in a git repo) and reads the
// current version from git tags, or from the .version file otherwise.
func BuildContext(dry bool) *Context {
	ctx := &Context{Version: &Semver{}, DryRun: dry}

	if top, err := git("rev-parse", "--show-toplevel"); err == nil {
		ctx.IsGitRepo = true
		os.Chdir(top)

		tag, err := git("describe", "--tags", "--abbrev=0")
		if err != nil {
			return ctx
		}
		version, err := ParseSemver(tag)
		if err != nil {
			fmt.Println("error reading git tag:", err)
			return ctx
		}
		ctx.Version, ctx.Source = version, SourceGit
		return ctx
	}

	content, err := os.ReadFile(VersionFile)
	if err != nil {
		return ctx
	}
	version, err := ParseSemver(string(content))
	if err != nil {
		fmt.Println("error reading version file:", err)
		return ctx
	}
	ctx.Version, ctx.Source = version, SourceFile

	return ctx
}

// Commit records the given version as a git tag, or in the .version file.
func (c *Context) Commit(version string) error {
	if c.IsGitRepo {
		return gitRun("tag", version)
	}
	return os.WriteFile(VersionFile, []byte(version), 0644)
}

// Push pushes the current version's tag to origin.
func (c *Context) Push() error {
	if !c.IsGitRepo {
		return errors.New("not a git repository")
	}
	fmt.Println("pushing git tag:", c.Version)
	return gitRun("push", "origin", c.Version.String())
}

// Reset removes every tag locally, remotely, or remote-only, or the .version file.
func (c *Context) Reset(remote, remoteOnly bool) error {
	if !c.IsGitRepo {
		return os.Remove(VersionFile)
	}

	tags, err := git("tag", "-l")
	if err != nil || tags == "" {
		return err
	}

	names := strings.Fields(tags)
	if remoteOnly {
		return UntagRemoteOnly(names)
	}
	return Untag(names, remote)
}

// UntagRemoteOnly deletes tags on origin without removing local tags.
func UntagRemoteOnly(versions []string) error {
	if OriginURL() == "" {
		return errors.New("no origin remote is configured")
	}
	for _, version := range versions {
		if err := gitRun("push", "--delete", "origin", version); err != nil {
			return err
		}
	}
	return nil
}

// Untag deletes the given tags locally, and remotely when remote is set.
func Untag(versions []string, remote bool) error {
	for _, version := range versions {
		if remote {
			if err := gitRun("push", "--delete", "origin", version); err != nil {
				return err
			}
		}
		if err := gitRun("tag", "-d", version); err != nil {
			return err
		}
	}
	return nil
}

// FetchTags fetches tags from the remote.
func FetchTags() error { return gitRun("fetch", "--tags") }

// RemoteSync fetches branches and tags from origin and prunes stale remote refs.
func RemoteSync() error {
	if OriginURL() == "" {
		return errors.New("no origin remote is configured")
	}
	if err := gitRun("fetch", "origin", "--prune"); err != nil {
		return err
	}
	return gitRun("fetch", "--tags")
}

// Branch is a local or remote-tracking branch ref.
type Branch struct {
	Name    string `json:"name"`
	Commit  string `json:"commit"`
	Subject string `json:"subject"`
	Current bool   `json:"current"`
}

// ListBranches returns local branches and origin remote-tracking branches.
func ListBranches() (local []Branch, remote []Branch, err error) {
	current, _ := git("rev-parse", "--abbrev-ref", "HEAD")

	local, err = parseBranches("refs/heads", current, false)
	if err != nil {
		return nil, nil, err
	}

	if OriginURL() != "" {
		remote, err = parseBranches("refs/remotes/origin", current, true)
		if err != nil {
			return local, nil, err
		}
	}

	return local, remote, nil
}

func parseBranches(refPrefix string, current string, isRemote bool) ([]Branch, error) {
	out, err := run("git", "for-each-ref", "--format=%(refname:short)%09%(objectname:short)%09%(contents:subject)", refPrefix)
	if err != nil {
		return nil, err
	}

	var branches []Branch
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 3)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		if isRemote && (name == "origin" || name == "origin/HEAD") {
			continue
		}
		b := Branch{Name: name, Commit: fields[1]}
		if len(fields) == 3 {
			b.Subject = fields[2]
		}
		if !isRemote && name == current {
			b.Current = true
		}
		branches = append(branches, b)
	}

	sort.Slice(branches, func(i, j int) bool {
		if branches[i].Current {
			return true
		}
		if branches[j].Current {
			return false
		}
		return branches[i].Name < branches[j].Name
	})

	return branches, nil
}

// CheckoutBranch switches to a local branch or creates one from a remote-tracking ref.
func CheckoutBranch(name string, isRemote bool) error {
	if WorkingTree().Any() {
		return errors.New("uncommitted changes — commit or stash before switching branches")
	}

	if isRemote {
		ref := name
		if !strings.HasPrefix(ref, "origin/") {
			ref = "origin/" + ref
		}
		local := strings.TrimPrefix(ref, "origin/")
		if local == "HEAD" {
			return errors.New("cannot checkout origin/HEAD")
		}
		if _, err := git("rev-parse", "--verify", "refs/heads/"+local); err == nil {
			return gitRun("checkout", local)
		}
		return gitRun("checkout", "-b", local, ref)
	}

	return gitRun("checkout", name)
}

// DeleteBranch removes a branch locally or on origin.
func DeleteBranch(name string, remote bool) error {
	if remote {
		if OriginURL() == "" {
			return errors.New("no origin remote is configured")
		}
		branch := strings.TrimPrefix(name, "origin/")
		if branch == "HEAD" || branch == "" {
			return errors.New("cannot delete origin/HEAD")
		}
		return gitRun("push", "origin", "--delete", branch)
	}

	current, _ := git("rev-parse", "--abbrev-ref", "HEAD")
	if name == current {
		return errors.New("cannot delete the current branch")
	}
	if err := gitRun("branch", "-d", name); err != nil {
		return gitRun("branch", "-D", name)
	}
	return nil
}

// CreateBranch creates a new local branch, optionally from a starting ref.
func CreateBranch(name, from string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("branch name is required")
	}
	if from == "" {
		return gitRun("branch", name)
	}
	return gitRun("branch", name, from)
}

// CommitEntry is one commit on a branch log.
type CommitEntry struct {
	Hash    string   `json:"hash"`
	Subject string   `json:"subject"`
	Author  string   `json:"author"`
	Date    string   `json:"date"`
	Tags    []string `json:"tags"`
}

// ListCommits returns recent commits on the given branch.
func ListCommits(branch string, limit, skip int) ([]CommitEntry, error) {
	if branch == "" {
		head := HeadCommit()
		branch = head.Branch
	}
	if branch == "" {
		return nil, errors.New("no branch to log")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if skip < 0 {
		skip = 0
	}

	out, err := run("git", "log", branch,
		"--format=%H%x09%h%x09%s%x09%an%x09%ad",
		"--date=short",
		"-n", fmt.Sprintf("%d", limit),
		"--skip", fmt.Sprintf("%d", skip))
	if err != nil {
		return nil, err
	}

	var commits []CommitEntry
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 5)
		if len(fields) < 5 {
			continue
		}
		entry := CommitEntry{
			Hash:    fields[1],
			Subject: fields[2],
			Author:  fields[3],
			Date:    fields[4],
			Tags:    []string{},
		}
		if tagOut, err := git("tag", "--points-at", fields[0]); err == nil {
			for _, tag := range strings.Fields(tagOut) {
				if _, err := ParseSemver(tag); err == nil {
					entry.Tags = append(entry.Tags, tag)
				}
			}
		}
		commits = append(commits, entry)
	}

	return commits, nil
}

// RemoteInfo describes a configured git remote.
type RemoteInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ListRemotes returns every configured remote and its fetch url.
func ListRemotes() ([]RemoteInfo, error) {
	out, err := run("git", "remote", "-v")
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var remotes []RemoteInfo
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name, url := fields[0], fields[1]
		if fields[len(fields)-1] != "(fetch)" || seen[name] {
			continue
		}
		seen[name] = true
		remotes = append(remotes, RemoteInfo{Name: name, URL: url})
	}

	sort.Slice(remotes, func(i, j int) bool {
		return remotes[i].Name < remotes[j].Name
	})

	return remotes, nil
}

// FetchRemote fetches all refs from a single remote and prunes deleted branches.
func FetchRemote(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("remote name is required")
	}
	if err := gitRun("fetch", name, "--prune"); err != nil {
		return err
	}
	if name == "origin" {
		return gitRun("fetch", "--tags")
	}
	return nil
}

// RepoName is the base name of the repository root directory.
func RepoName() string {
	top, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		wd, _ := os.Getwd()
		return filepath.Base(wd)
	}
	return filepath.Base(top)
}

// OriginURL is the url of the origin remote, empty when there is none.
func OriginURL() string {
	url, _ := git("remote", "get-url", "origin")
	return url
}

// Tag is a semver tag of the repository, and whether it exists on origin.
type Tag struct {
	Name    string `json:"name"`
	Commit  string `json:"commit"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
	Pushed  bool   `json:"pushed"`
	Current bool   `json:"current"`
}

// taggedVersion keeps a tag next to its parsed version, so that sorting by
// semver precedence moves both together.
type taggedVersion struct {
	tag     Tag
	version *Semver
}

// ListTags returns every semver tag of the repository, newest version first.
// Tags that are not valid semver are ignored.
func ListTags() ([]Tag, error) {
	out, err := run("git", "for-each-ref", "--format=%(refname:short)%09%(objectname:short)%09%(creatordate:short)%09%(contents:subject)", "refs/tags")
	if err != nil {
		return nil, err
	}

	remote := remoteTags()

	var entries []taggedVersion

	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 4)
		if len(fields) < 3 {
			continue
		}
		version, err := ParseSemver(fields[0])
		if err != nil {
			continue
		}
		tag := Tag{Name: fields[0], Commit: fields[1], Date: fields[2], Pushed: remote[fields[0]]}
		if len(fields) == 4 {
			tag.Subject = fields[3]
		}
		entries = append(entries, taggedVersion{tag, version})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return Compare(entries[i].version, entries[j].version) > 0
	})

	tags := make([]Tag, len(entries))
	for i, entry := range entries {
		tags[i] = entry.tag
	}
	if len(tags) > 0 {
		tags[0].Current = true
	}

	return tags, nil
}

// Head is the commit a new tag would be created on. Tags lists the semver tags
// already pointing at it.
type Head struct {
	Branch  string   `json:"branch"`
	Commit  string   `json:"commit"`
	Subject string   `json:"subject"`
	Tags    []string `json:"tags"`
}

// HeadCommit describes HEAD. Fields are empty when the repository has no
// commits yet.
func HeadCommit() Head {
	head := Head{Tags: []string{}}
	head.Branch, _ = git("rev-parse", "--abbrev-ref", "HEAD")
	head.Commit, _ = git("rev-parse", "--short", "HEAD")
	head.Subject, _ = git("log", "-1", "--pretty=%s")

	// --points-at resolves annotated tags to their commit, which comparing
	// object names would not.
	if out, err := git("tag", "--points-at", "HEAD"); err == nil {
		for _, name := range strings.Fields(out) {
			if _, err := ParseSemver(name); err == nil {
				head.Tags = append(head.Tags, name)
			}
		}
	}

	return head
}

// Changes summarises the working tree: how many paths are staged, modified but
// unstaged, or untracked.
type Changes struct {
	Staged    int      `json:"staged"`
	Unstaged  int      `json:"unstaged"`
	Untracked int      `json:"untracked"`
	Files     []string `json:"files"`
}

func (c Changes) Any() bool { return c.Staged+c.Unstaged+c.Untracked > 0 }

// maxChangedFiles caps how many paths the ui lists, so a huge working tree
// does not produce an unusable page.
const maxChangedFiles = 50

// WorkingTree reads the status of the working tree. Paths are reported in
// `git status --short` form, e.g. " M repo.go".
func WorkingTree() Changes {
	out, err := runRaw("git", "status", "--porcelain")
	if err != nil || strings.TrimSpace(out) == "" {
		return Changes{Files: []string{}}
	}

	changes := Changes{Files: []string{}}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 3 {
			continue
		}
		index, worktree := line[0], line[1]
		switch {
		case index == '?' && worktree == '?':
			changes.Untracked++
		default:
			if index != ' ' {
				changes.Staged++
			}
			if worktree != ' ' {
				changes.Unstaged++
			}
		}
		if len(changes.Files) < maxChangedFiles {
			changes.Files = append(changes.Files, line)
		}
	}

	return changes
}

// CommitChanges commits the working tree with the given message. When all is
// set every change (including untracked files) is staged first, otherwise only
// what is already staged is committed.
func CommitChanges(message string, all bool) error {
	if strings.TrimSpace(message) == "" {
		return errors.New("a commit message is required")
	}

	changes := WorkingTree()
	if all {
		if !changes.Any() {
			return errors.New("nothing to commit")
		}
		if err := gitRun("add", "-A"); err != nil {
			return err
		}
	} else if changes.Staged == 0 {
		return errors.New("no staged changes to commit")
	}

	return gitRun("commit", "-m", message)
}

// HighestVersion is the greatest semver tag of the repository, or nil when
// there is none. Unlike `git describe`, it considers tags on every branch and
// orders them by semver precedence.
func HighestVersion() *Semver {
	out, err := git("tag", "-l")
	if err != nil {
		return nil
	}

	var highest *Semver
	for _, name := range strings.Fields(out) {
		version, err := ParseSemver(name)
		if err != nil {
			continue
		}
		if highest == nil || Compare(version, highest) > 0 {
			highest = version
		}
	}

	return highest
}

// remoteTags is the set of tag names that exist on origin. It falls back to
// the cached remote refs when origin is unreachable.
func remoteTags() map[string]bool {
	tags := map[string]bool{}

	out, err := runTimeout(remoteTimeout, "git", "ls-remote", "--tags", "origin")
	if err != nil {
		return tags
	}

	for _, line := range strings.Split(out, "\n") {
		_, ref, ok := strings.Cut(line, "refs/tags/")
		if !ok {
			continue
		}
		tags[strings.TrimSuffix(strings.TrimSpace(ref), "^{}")] = true
	}

	return tags
}

// PushTag pushes a single tag to origin.
func PushTag(version string) error { return gitRun("push", "origin", version) }

// PushAllTags pushes every local tag to origin.
func PushAllTags() error { return gitRun("push", "origin", "--tags") }

// RenameTag re-tags the commit of from as to, and deletes from. When remote is
// set the rename is mirrored on origin.
func RenameTag(from, to string, remote bool) error {
	commit, err := git("rev-list", "-n", "1", from)
	if err != nil {
		return err
	}
	if err := gitRun("tag", to, commit); err != nil {
		return err
	}
	if err := Untag([]string{from}, remote); err != nil {
		return err
	}
	if remote {
		return PushTag(to)
	}
	return nil
}

// References greps the project for occurrences of the current version.
func (c *Context) References() (string, error) {
	version := strings.TrimPrefix(c.Version.String(), "v")
	return run("grep", "-rnH", "--exclude-dir=.git", "--exclude=go.mod", "--exclude=go.sum", version, ".")
}

// remoteTimeout bounds the reachability check for origin, so an unreachable
// or credential-prompting remote does not stall the ui.
const remoteTimeout = 5 * time.Second

func git(args ...string) (string, error) { return run("git", args...) }

func gitRun(args ...string) error {
	_, err := git(args...)
	return err
}

// run executes a command and returns its trimmed output.
func run(name string, args ...string) (string, error) {
	out, err := output(exec.Command(name, args...))
	return strings.TrimSpace(out), err
}

// runRaw is run without trimming, for output whose leading whitespace is
// significant — `git status --porcelain` encodes state in the first columns.
func runRaw(name string, args ...string) (string, error) {
	return output(exec.Command(name, args...))
}

// runTimeout is run, aborting the command after timeout. Credential prompts
// are disabled so a command can never block on user input.
func runTimeout(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_SSH_COMMAND=ssh -oBatchMode=yes")

	out, err := output(cmd)
	return strings.TrimSpace(out), err
}

func output(cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", errors.New(msg)
		}
		return "", err
	}

	return stdout.String(), nil
}
