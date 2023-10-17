/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

func defaultVersion() *verman.Semver {
	return &verman.Semver{Patch: 1}
}

func initializeVersion(args []string) (version *verman.Semver, useGitTag bool, err error) {
	if version, err = getVersionFromArg(args); err == nil {
		return
	}

	if !verman.IsGitRepository() {
		version = defaultVersion()
		return
	}

	version, err = verman.GetVersionFromGitTag()
	switch {
	case errors.Is(err, verman.ErrGettingGitTag), errors.Is(err, verman.ErrInvalidVersionFormat):
		version = defaultVersion()
		err = nil
	case err == nil:
		useGitTag = true
	}
	return
}

func printVersionInitialization(version *verman.Semver, useGitTag bool, err error) {
	if err != nil {
		fmt.Println("Error initializing version:", err)
		return
	}

	if useGitTag {
		fmt.Println("latest git tag found:", version.String())
	} else {
		fmt.Println("Setting version to", version.String())
	}
}

var (
	ErrWritingConfig = errors.New("error writing to configuration file")
	ErrCreatingTag   = errors.New("error creating git tag: check if the tag already exists")
)

func setVersion(version *verman.Semver, useGitTag bool) error {
	if err := verman.WriteVersionToConfig(version); err != nil {
		return ErrWritingConfig
	}

	if !useGitTag && verman.IsGitRepository() {
		if err := verman.GitCommitVersionConfig(version); err != nil {
			return err
		}

		if err := verman.GitTagVersion(version); err != nil {
			return ErrCreatingTag
		}
	}

	return nil
}

func printSetVersionActions(version *verman.Semver, err error) {
	fmt.Println("setting current version:", version.String(), "...")

	if err != nil {
		switch {
		case errors.Is(err, ErrWritingConfig):
			fmt.Println("error writing to configuration file.")
		case errors.Is(err, ErrCreatingTag):
			fmt.Println("error creating git tag: check if the tag already exists.")
		}
		return
	}

	if verman.IsGitRepository() {
		fmt.Printf("creating git tag %s...\n", version.String())
	}
}

func getVersionFromArg(args []string) (*verman.Semver, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("version not provided")
	}
	version, err := verman.Parse(args[0])
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}

	return version, nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the semver configuration",
	Long: `Initialize the semver configuration. This will create a .version.yaml file in the current directory.
This file will contain the current version of the project.

It will get latest tag from git and set it as the current version, if the git tag is a semver tag.
If no git tags are found, it will set the version to 0.0.1`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if configExists {
			fmt.Println("configuration already exists. run `semver get` to display the current version or `semver reset` to reset all tags and config.")
			return
		}

		fmt.Println("initializing configuration...")

		projectVersion, useGitTag, err := initializeVersion(args)
		printVersionInitialization(projectVersion, useGitTag, err)

		err = setVersion(projectVersion, useGitTag)
		printSetVersionActions(projectVersion, err)

		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
