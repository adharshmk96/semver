package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

// Reset removes every tag (and optionally every remote tag), or the .version file.
func (c *Context) Reset(remote bool) error {
	if !c.IsGitRepo {
		return os.Remove(VersionFile)
	}

	tags, err := git("tag", "-l")
	if err != nil || tags == "" {
		return err
	}

	return Untag(strings.Fields(tags), remote)
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

// References greps the project for occurrences of the current version.
func (c *Context) References() (string, error) {
	version := strings.TrimPrefix(c.Version.String(), "v")
	return run("grep", "-rnH", "--exclude-dir=.git", "--exclude=go.mod", "--exclude=go.sum", version, ".")
}

func git(args ...string) (string, error) { return run("git", args...) }

func gitRun(args ...string) error {
	_, err := git(args...)
	return err
}

// run executes a command and returns its trimmed first line of output.
func run(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", errors.New(msg)
		}
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
