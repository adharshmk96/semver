/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/adharshmk96/semver/pkg/parser"
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
			fmt.Println("Configuration already exists. Please run `semver get` to display the current version.")
			return
		}

		// get latest tag from git
		// if no tags found, set version to 0.0.0
		// if tag found, set version to tag
		// if tag is not semver, set version to 0.0.0
		// if tag is semver, set version to tag
		lastGitTag, err := GetLastTagFromGit()
		if err != nil {
			fmt.Println("No git tags found. Setting version to v0.0.0")
			viper.Set("version", "v0.0.0")
		}

		version, err := parser.Parse(lastGitTag)
		if err != nil {
			fmt.Println("Invalid version from git tag. Setting version to v0.0.0")
			viper.Set("version", "v0.0.0")
		}

		viper.Set("major", version.Major)
		viper.Set("minor", version.Minor)
		viper.Set("patch", version.Patch)
		viper.Set("alpha", version.Alpha)
		viper.Set("beta", version.Beta)
		viper.Set("rc", version.RC)

		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func GetLastTagFromGit() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting last git tag: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
