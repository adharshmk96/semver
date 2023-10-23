package verman

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

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
	version, err := ParseSemver(versionString)
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
func GitRemoveRemoteTag(tag string) error {
	cmd := exec.Command("git", "push", "--delete", "origin", tag)

	err := cmd.Run()
	if err != nil {
		return ErrCreatingGitTag
	}

	return nil
}

// remove tag
func GitRemoveLocalTag(tag string) error {
	cmd := exec.Command("git", "tag", "-d", tag)

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

// reset all remote tags
func GitRemoveAllRemoteTags() error {
	var cmd *exec.Cmd

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
	err := gitAdd(".version.yaml")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Commits all staged.
	return gitCommit("bump version to " + version.String())
}

func GitTagExists(tag string) bool {
	cmd := exec.Command("git", "rev-parse", tag)
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
