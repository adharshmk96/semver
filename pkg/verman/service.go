package verman

import (
	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/adharshmk96/semver/pkg/verman/core"
	"github.com/spf13/afero"
)

func InitializeSemver(ctx *core.Context, tag string) error {
	if tag == "" {
		tag = ctx.CurrentVersion.String()
	}
	if ctx.IsGitRepo {
		gitCmd := commands.NewGitCmd(commands.NewGitExec())
		return gitCmd.TagVersion(tag)
	}

	return writeToFile(core.VERSION_FILE, tag)
}

func PushGitTag(ctx *core.Context) error {
	if !ctx.IsGitRepo {
		return core.ErrNotGitRepo
	}
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	return gitCmd.PushTag(ctx.CurrentVersion.String())
}

func CommitVersionLocally(ctx *core.Context) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	if ctx.IsGitRepo {
		return gitCmd.TagVersion(ctx.CurrentVersion.String())
	}

	fileContent := ctx.CurrentVersion.String()
	return writeToFile(core.VERSION_FILE, fileContent)
}

func ResetSemver(ctx *core.Context, remote bool) error {

	if !ctx.IsGitRepo {
		fsRepo := NewFileRepo(afero.NewOsFs())
		return fsRepo.DeleteFile(core.VERSION_FILE)
	}

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

func UntagVersions(versions []string, remote bool) error {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())

	for _, version := range versions {
		if remote {
			err := gitCmd.RemoveRemoteTag(version)
			if err != nil {
				return err
			}
		}
		err := gitCmd.RemoveTag(version)
		if err != nil {
			return err
		}
	}

	return nil

}
