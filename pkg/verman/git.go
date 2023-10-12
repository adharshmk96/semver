package verman

import (
	"fmt"
	"os/exec"
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

func GetVersionFromGitTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting last git tag: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
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
func GitRemoveAllTags() error {
	cmd := exec.Command("git", "tag", "-d", "$(git tag -l)")

	err := cmd.Run()
	if err != nil {
		return ErrCreatingGitTag
	}

	return nil
}

func GitRemoveAllRemoteTags() error {
	cmd := exec.Command("git", "push", "--delete", "origin", "$(git tag -l)")

	err := cmd.Run()
	if err != nil {
		return ErrCreatingGitTag
	}

	return nil
}
