package verman

import (
	"log"
	"os"

	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/spf13/afero"
)

func BuildContext(args []string, dry bool) *Context {
	gitCmd := commands.NewGitCmd(commands.NewGitExec())
	fs := afero.NewOsFs()
	fileRepo := NewFileRepo(fs)

	// Determine the working directory
	workDir := determineWorkDir(gitCmd)

	// Check if current directory has a .version file and parse it
	semver, isVerman := fetchSemverIfPresent(fileRepo)

	return &Context{
		WorkDir:        workDir,
		CurrentVersion: semver,
		DryRun:         dry,
		IsGitRepo:      gitCmd.IsRepo(),
		IsVerman:       isVerman,
	}
}

func determineWorkDir(gitCmd commands.GitCmd) string {
	if gitCmd.IsRepo() {
		workDir, err := gitCmd.GetTopLevel()
		if err != nil {
			return "."
		}
		return workDir
	}
	return "."
}

func fetchSemverIfPresent(fileRepo FileRepo) (*Semver, bool) {
	if fileRepo.FileExists(".version") {
		fileContent, err := fileRepo.ReadFileContent(".version")
		if err != nil {
			return nil, false
		}

		semver, err := ParseSemver(fileContent)
		if err != nil {
			log.Fatalln("error parsing version file:", err)
			os.Exit(1)
		}
		return semver, true
	}
	return nil, false
}
