package verman

import "github.com/adharshmk96/semver/pkg/commands"

func InitializeSemver(ctx *Context) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	if gitCmd.IsRepo() {
		return gitCmd.TagVersion(ctx.CurrentVersion.String())
	}

	fileContent := ctx.CurrentVersion.String()
	return writeToFile(VERSION_FILE, fileContent)
}

func PushGitTag(ctx *Context) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	return gitCmd.PushTag(ctx.CurrentVersion.String())
}

func CommitVersionLocally(ctx *Context) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	if gitCmd.IsRepo() {
		return gitCmd.TagVersion(ctx.CurrentVersion.String())
	}

	fileContent := ctx.CurrentVersion.String()
	return writeToFile(VERSION_FILE, fileContent)
}
