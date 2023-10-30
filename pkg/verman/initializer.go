package verman

import (
	"fmt"
	"log"
	"os"

	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/adharshmk96/semver/pkg/verman/core"
	"github.com/spf13/afero"
)

func BuildContext(dry bool) *core.Context {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	fs := afero.NewOsFs()
	fileRepo := NewFileRepo(fs)

	workDir := determineWorkDir(gitCmd)
	os.Chdir(workDir)

	IsGitRepo := gitCmd.IsRepo()
	source, currentVersion := getVersion(gitCmd, fileRepo)

	if IsGitRepo {
		gitCmd.Run("fetch", "--tags")
	}

	return &core.Context{
		WorkDir:        workDir,
		CurrentVersion: currentVersion,
		DryRun:         dry,
		SemverSource:   source,
		IsGitRepo:      IsGitRepo,
	}
}

func determineWorkDir(gitCmd commands.GitCmd) string {
	if !gitCmd.IsRepo() {
		return "."
	}

	workDir, err := gitCmd.GetTopLevel()
	if err != nil {
		log.Println("Failed to get git top level directory:", err)
		return "."
	}
	return workDir
}

func getVersion(gitCmd commands.GitCmd, filerepo FileRepo) (core.Source, *core.Semver) {
	if gitCmd.IsRepo() {
		semver, err := getVersionFromGit(gitCmd)
		if err != nil {
			fmt.Println("error reading git tags:", err)
			return core.SourceNone, &core.Semver{}
		}
		return core.SourceGit, semver
	}

	if filerepo.FileExists(core.VERSION_FILE) {
		semver, err := getVersionFromFile(filerepo)
		if err != nil {
			fmt.Println("error reading version file:", err)
			return core.SourceNone, &core.Semver{}
		}
		return core.SourceFile, semver
	}

	return core.SourceNone, &core.Semver{}

}

func getVersionFromGit(gitCmd commands.GitCmd) (*core.Semver, error) {
	latestTag, err := gitCmd.LatestTag()
	if err != nil {
		return nil, err
	}

	return core.ParseSemver(latestTag)
}

func getVersionFromFile(filerepo FileRepo) (*core.Semver, error) {
	content, err := filerepo.ReadFileContent(core.VERSION_FILE)
	if err != nil {
		return nil, err
	}

	return core.ParseSemver(string(content))
}
