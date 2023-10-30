/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var SemVer = "development"

func displaySemverInfo() {
	if SemVer != "development" {
		fmt.Println(SemVer)
		return
	}
	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != "(devel)" && version.Main.Version != "" {
		SemVer = version.Main.Version
	}
	fmt.Println(SemVer)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Manage your project's semver configuration",
	Long: `A CLI tool to manage your project's semantic version.

semver uses git tags or .version file (for non-git projects) to manage the version of the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			displaySemverInfo()
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "display current version")
}
