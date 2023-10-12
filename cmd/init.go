/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

func initializeVersion() (version *verman.Semver, useGitTag bool, err error) {
	useGitTag = false
	lastGitTag, err := verman.GetVersionFromGitTag()
	if err != nil {
		fmt.Println("no git tags found. setting version to v0.0.0")
		lastGitTag = "v0.0.0"
	} else {
		fmt.Println("found git tag:", lastGitTag)
		useGitTag = true
	}

	version, err = verman.Parse(lastGitTag)
	if err != nil {
		fmt.Println("invalid version from git tag. setting version to v0.0.0")
	}
	return version, useGitTag, err
}

func setVersion(version *verman.Semver, useGitTag bool) error {
	fmt.Println("setting version to ", version.String(), "...")
	if err := verman.WriteVersionToConfig(version); err != nil {
		return fmt.Errorf("error writing to configuration file: %w", err)
	}

	if !useGitTag {
		fmt.Println("creating git tag...")
		if err := verman.GitTagVersion(version); err != nil {
			return fmt.Errorf("error creating git tag: check if the tag already exists: %w", err)
		}
	}

	return nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the semver configuration",
	Long: `Initialize the semver configuration. This will create a version.yaml file in the current directory.
This file will contain the current version of the project.

It will get latest tag from git and set it as the current version, if the git tag is a semver tag.
If no git tags are found, it will set the version to 0.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		if configExists {
			fmt.Println("configuration already exists. run `semver get` to display the current version.")
			return
		}

		fmt.Println("initializing configuration...")

		version, useGitTag, err := initializeVersion()
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
