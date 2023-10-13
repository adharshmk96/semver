package verman

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func GitTagVersion(semver *Semver) error {
	cmd := exec.Command("git", "tag", semver.String())

	out, err := cmd.Output()
	if err != nil || len(out) > 0 {
		return ErrCreatingGitTag
	}

	return nil
}

func GitRemoveTag(semver *Semver) error {
	cmd := exec.Command("git", "tag", "-d", semver.String())

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

func GetVersionFromGitTag() (*Semver, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrGettingGitTag
	}
	versionString := strings.TrimSpace(string(output))
	version, err := Parse(versionString)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func GitPushTag(semver *Semver) error {
	cmd := exec.Command("git", "push", "origin", semver.String())

	err := cmd.Run()
	if err != nil {
		return ErrCreatingGitTag
	}

	return nil
}

// remove tag
func GitRemoveRemoteTag(semver *Semver) error {
	cmd := exec.Command("git", "push", "--delete", "origin", semver.String())

	err := cmd.Run()
	if err != nil {
		return ErrCreatingGitTag
	}

	return nil
}

func GitRemoveLocalTag(semver *Semver) error {
	cmd := exec.Command("git", "tag", "-d", semver.String())

	err := cmd.Run()
	if err != nil {
		return ErrCreatingGitTag
	}

	return nil
}

// reset all tags
func GitRemoveAllLocalTags() error {
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
func GitRemoveAllRemoteTags() error {
	cmd := exec.Command("git", "push", "--delete", "origin", "$(git tag -l)")

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func gitAdd(file string) error {
	cmd := exec.Command("git", "add", file)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func GitCommitVersionConfig(version *Semver) error {
	err := gitAdd("version.yaml")
	if err != nil {
		return err
	}
	// Commits all staged.
	return gitCommit("bump version to " + version.String())
}
