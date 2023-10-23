package commands

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

// TODO: move this to another repo / package

type CmdExecutor interface {
	RunCmd(args ...string) (string, error)
}

type GitCmd interface {
	Run(args ...string) (string, error)
	IsRepo() bool
	GetTopLevel() (string, error)
	Revparse(ref string) (string, error)
	TagVersion(version string) error
	RemoveTag(version string) error
	RemoveRemoteTag(version string) error
	RemoveAllLocalTags() error
	RemoveAllRemoteTags() error
	LatestTag() (string, error)
	PushTag(version string) error
	PullTags(version string) error
	Add(files []string) error
	Commit(message string) error
}

type gitCommands struct {
	exec CmdExecutor
}

func NewGitCmd(e CmdExecutor) GitCmd {
	return &gitCommands{exec: e}
}

type gitExec struct{}

func NewGitExec() CmdExecutor {
	return &gitExec{}
}

func (g *gitExec) RunCmd(args ...string) (string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	c := exec.Command("git", args...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return Clean(stdout.String(), nil)
}

func (g *gitCommands) Run(args ...string) (string, error) {
	return g.exec.RunCmd(args...)
}

func (g *gitCommands) Revparse(ref string) (string, error) {
	out, err := g.exec.RunCmd("rev-parse", ref)
	if err != nil {
		return "", err
	}

	return out, nil
}

// Local commands
func (g *gitCommands) IsRepo() bool {
	out, err := g.Revparse("--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

func (g *gitCommands) GetTopLevel() (string, error) {
	out, err := g.Revparse("--show-toplevel")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func (g *gitCommands) TagVersion(version string) error {
	_, err := g.exec.RunCmd("tag", version)
	return err
}

func (g *gitCommands) RemoveTag(version string) error {
	_, err := g.exec.RunCmd("tag", "-d", version)
	return err
}

func (g *gitCommands) RemoveAllLocalTags() error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "for /f \"delims=\" %i in ('git tag -l') do git tag -d %i")
	} else {
		cmd = exec.Command("bash", "-c", "git tag -d $(git tag -l)")
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (g *gitCommands) LatestTag() (string, error) {
	out, err := g.exec.RunCmd("describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func (g *gitCommands) Add(files []string) error {
	_, err := g.exec.RunCmd("add", strings.Join(files, " "))
	return err
}

func (g *gitCommands) Commit(message string) error {
	_, err := g.exec.RunCmd("commit", "-m", message)
	return err
}

// Remote commands
func (g *gitCommands) PushTag(version string) error {
	_, err := g.exec.RunCmd("push", "origin", version)
	return err
}

func (g *gitCommands) PullTags(version string) error {
	_, err := g.exec.RunCmd("fetch", "--tags")
	return err
}

func (g *gitCommands) RemoveRemoteTag(version string) error {
	_, err := g.exec.RunCmd("push", "--delete", "origin", version)
	return err
}

func (g *gitCommands) RemoveAllRemoteTags() error {
	var cmd *exec.Cmd

	g.exec.RunCmd("fetch", "--tags")

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "for /f \"delims=\" %i in ('git tag -l') do git push --delete origin %i")
	} else {
		cmd = exec.Command("bash", "-c", "git push --delete origin $(git tag -l)")
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
