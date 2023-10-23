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

// releaseCmd represents the get command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "removes the pre-release and tags the release version",
	Run: func(cmd *cobra.Command, args []string) {
		dry, _ := cmd.Flags().GetBool("dry")
		ctx := verman.BuildContext(dry)

		if ctx.SemverSource == core.SourceNone {
			fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
			return
		}

		if !ctx.CurrentVersion.IsPreRelease() {
			fmt.Println("current version is not a pre-release.")
			return
		}

		verman.UpdateAndCommitVersion(ctx, "release")

		fmt.Println(ctx.CurrentVersion.String())
	},
}

func init() {
	releaseCmd.Flags().BoolP("dry", "d", false, "dry run mode")
	rootCmd.AddCommand(releaseCmd)
}
