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

func GetLastTagFromGit() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting last git tag: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
