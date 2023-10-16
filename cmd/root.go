/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configExists bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Manage your project's semver configuration",
	Long: `A CLI tool to manage your project's semantic version.

semver uses a version.yaml file to store the current version of the project.
It also uses git tags to manage the version of the project. Updating the version will also tag the git repository with the new version.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func gitTopLevelDirectory() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "."
	}
	return strings.TrimSpace(string(out))
}

func init() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(gitTopLevelDirectory())
	viper.SetConfigName("version")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			escape := []string{"init", "reset", "version"}
			if len(os.Args) > 1 && !slices.Contains(escape, os.Args[1]) {
				fmt.Println("no version.yaml configuration found. run `semver init` to initialize the configuration.")
				os.Exit(0)
			} else {
				configExists = false
			}

		} else {
			os.Exit(1)
		}
	} else {
		configExists = true
	}
}
