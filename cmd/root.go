/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Manage your project's semver configuration",
	Long: `A CLI tool to manage your project's semantic version.

semver uses a .version.yaml file to store the current version of the project.
It also uses git tags to manage the version of the project. Updating the version will also tag the git repository with the new version.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
