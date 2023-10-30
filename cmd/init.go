/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/adharshmk96/semver/pkg/verman/core"
	"github.com/spf13/cobra"
)

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
		ctx := verman.BuildContext(false)

		if ctx.SemverSource != core.SourceNone {
			fmt.Println("semver config found, run `semver get` to view the version.")
			return
		}

		var initVersion string
		if len(args) > 0 {
			initVersion = args[0]
		} else {
			initVersion = "v0.0.1"
		}

		fmt.Println("initializing...")
		err := verman.InitializeSemver(ctx, initVersion)
		if err != nil {
			fmt.Println("error: initalizing semver.", err)
			return
		}

		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
