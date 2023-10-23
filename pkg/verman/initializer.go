package verman

import (
	"fmt"
	"log"
	"os"

	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/spf13/afero"
)

func BuildContext(args []string, dry bool) *Context {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	fs := afero.NewOsFs()
	fileRepo := NewFileRepo(fs)

	workDir := determineWorkDir(gitCmd)
	os.Chdir(workDir)

	source, currentVersion := getVersion(gitCmd, fileRepo)

	return &Context{
		WorkDir:        workDir,
		CurrentVersion: currentVersion,
		DryRun:         dry,
		SemverSource:   source,
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

func getVersion(gitCmd commands.GitCmd, filerepo FileRepo) (Source, *Semver) {
	if gitCmd.IsRepo() {
		semver, err := getVersionFromGit(gitCmd)
		if err != nil {
			fmt.Println("error reading git tags:", err)
			return SourceNone, &Semver{}
		}
		return SourceGit, semver
	}

	if filerepo.FileExists(VERSION_FILE) {
		semver, err := getVersionFromFile(filerepo)
		if err != nil {
			fmt.Println("error reading version file:", err)
			return SourceNone, &Semver{}
		}
		return SourceFile, semver
	}

	return SourceNone, &Semver{}

}

func getVersionFromGit(gitCmd commands.GitCmd) (*Semver, error) {
	latestTag, err := gitCmd.LatestTag()
	if err != nil {
		return nil, err
	}

	return ParseSemver(latestTag)
}

func getVersionFromFile(filerepo FileRepo) (*Semver, error) {
	content, err := filerepo.ReadFileContent(VERSION_FILE)
	if err != nil {
		return nil, err
	}

	return ParseSemver(string(content))
}
