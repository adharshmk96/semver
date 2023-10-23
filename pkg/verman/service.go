package verman

import "github.com/adharshmk96/semver/pkg/commands"

func InitializeSemver(ctx *Context) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	if gitCmd.IsRepo() {
		latestTag := ctx.CurrentVersion.String()
		return gitCmd.TagVersion(latestTag)
	}

	return WriteVersionToFile(ctx.CurrentVersion)
}

func WriteVersionToFile(semver *Semver) error {
	fileContent := semver.String()
	return writeToFile(VERSION_FILE, fileContent)
}
