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

func initializeVersion(args []string) (version *verman.Semver, useGitTag bool, err error) {
	useGitTag = false

	version, err = getVersionFromArg(args)
	if err == nil {
		return version, useGitTag, nil
	}

	version, err = verman.GetVersionFromGitTag()
	if err != nil {
		if errors.Is(err, verman.ErrGettingGitTag) {
			fmt.Println("no git tags found. setting version to v0.0.0")
		}
		if errors.Is(err, verman.ErrInvalidVersionFormat) {
			fmt.Println("latest git tag is not a valid semver tag. setting version to v0.0.0")
		}
	} else {
		fmt.Println("latest git tag found:", version.String())
		useGitTag = true
	}

	return version, useGitTag, err
}

func setVersion(version *verman.Semver, useGitTag bool) error {
	fmt.Println("setting current version:", version.String(), "...")
	if err := verman.WriteVersionToConfig(version); err != nil {
		return fmt.Errorf("error writing to configuration file: %w", err)
	}

	if !useGitTag {
		fmt.Printf("creating git tag %s...\n", version.String())
		if err := verman.GitTagVersion(version); err != nil {
			return fmt.Errorf("error creating git tag: check if the tag already exists: %w", err)
		}
	}

	return nil
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
	Long: `Initialize the semver configuration. This will create a version.yaml file in the current directory.
This file will contain the current version of the project.

It will get latest tag from git and set it as the current version, if the git tag is a semver tag.
If no git tags are found, it will set the version to 0.0.0`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if configExists {
			fmt.Println("configuration already exists. run `semver get` to display the current version.")
			return
		}

		fmt.Println("initializing configuration...")

		version, useGitTag, err := initializeVersion(args)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := setVersion(version, useGitTag); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
