package verman

import (
	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/adharshmk96/semver/pkg/verman/core"
)

func InitializeSemver(ctx *core.Context, tag string) error {
	if tag == "" {
		tag = ctx.CurrentVersion.String()
	}
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	if gitCmd.IsRepo() {
		return gitCmd.TagVersion(tag)
	}

	return writeToFile(core.VERSION_FILE, tag)
}

func PushGitTag(ctx *core.Context) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	return gitCmd.PushTag(ctx.CurrentVersion.String())
}

func CommitVersionLocally(ctx *core.Context) error {
	if ctx.DryRun {
		return nil
	}
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	if gitCmd.IsRepo() {
		return gitCmd.TagVersion(ctx.CurrentVersion.String())
	}

	fileContent := ctx.CurrentVersion.String()
	return writeToFile(core.VERSION_FILE, fileContent)
}

func ResetSemver(ctx *core.Context, remote bool) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())

	if remote {
		err := gitCmd.RemoveAllRemoteTags()
		if err != nil {
			return err
		}
	}

	err := gitCmd.RemoveAllLocalTags()
	if err != nil {
		return err
	}

	return nil
}

func UntagVersion(version string, remote bool) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())

	if remote {
		err := gitCmd.RemoveRemoteTag(version)
		if err != nil {
			return err
		}
	}

	return gitCmd.RemoveTag(version)
}

func UpdateAndCommitVersion(ctx *core.Context, incType string) error {
	ctx.CurrentVersion.UpdateSemver(incType)
	if ctx.DryRun {
		return nil
	}
	return CommitVersionLocally(ctx)
}
