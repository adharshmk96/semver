/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
			fmt.Println("configuration already exists. Please run `semver get` to display the current version.")
			return
		}

		fmt.Println("initializing semver configuration...")
		lastGitTag, err := verman.GetLastTagFromGit()
		if err != nil {
			fmt.Println("no git tags found. setting version to v0.0.0")
			viper.Set("version", "v0.0.0")
			lastGitTag = "v0.0.0"
		}

		fmt.Println("found git tag:", lastGitTag)

		version, err := verman.Parse(lastGitTag)
		if err != nil {
			fmt.Println("invalid version from git tag. setting version to v0.0.0")
			viper.Set("version", "v0.0.0")
		}

		fmt.Print("setting version to ", version.String(), "...")

		viper.Set("major", version.Major)
		viper.Set("minor", version.Minor)
		viper.Set("patch", version.Patch)
		viper.Set("alpha", version.Alpha)
		viper.Set("beta", version.Beta)
		viper.Set("rc", version.RC)

		viper.WriteConfigAs("version.yaml")

		fmt.Println("creating git tag...")
		err = verman.GitTagVersion(version)
		if err != nil {
			fmt.Println("error creating git tag: check if the tag already exists.")
			return
		}

		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
