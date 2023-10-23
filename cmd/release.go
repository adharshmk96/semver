/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

// releaseCmd represents the get command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "removes the pre-release and tags the release version",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := verman.BuildContext(args, false)

		if ctx.SemverSource == verman.SourceNone {
			fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
			return
		}

		if ctx.CurrentVersion.IsPreRelease() {
			fmt.Println("current version is not a pre-release.")
			return
		}

		fmt.Println("current version:", ctx.CurrentVersion.String())

		ctx.CurrentVersion.Release()

		verman.CommitVersionLocally(ctx)

		fmt.Println("updated version:", ctx.CurrentVersion.String())
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
