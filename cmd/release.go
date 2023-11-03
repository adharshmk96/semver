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
		push, _ := cmd.Flags().GetBool("push")
		ctx := verman.BuildContext(dry)

		if ctx.SemverSource == core.SourceNone {
			fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
			return
		}

		if !ctx.CurrentVersion.IsPreRelease() {
			fmt.Println("current version is not a pre-release.")
			return
		}

		if ctx.DryRun {
			fmt.Println(ctx.CurrentVersion.String())
			return
		}

		verman.CommitVersionLocally(ctx)
		if push {
			fmt.Println("pushing git tag:", ctx.CurrentVersion.String())
			err := verman.PushGitTag(ctx)
			if err != nil {
				fmt.Println("error pushing git tag:", err)
				return
			}
		}

		fmt.Println(ctx.CurrentVersion.String())
	},
}

func init() {
	releaseCmd.Flags().BoolP("dry", "d", false, "dry run mode")
	releaseCmd.Flags().BoolP("push", "p", false, "push the tag to remote")
	releaseCmd.Flags().BoolVar(&sync, "sync", false, "fetch remote tags")
	rootCmd.AddCommand(releaseCmd)
}
