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

var untagRemote bool

var untagCmd = &cobra.Command{
	Use:   "untag",
	Short: "Delete a specific tag from git (default: current tag)",
	Run: func(cmd *cobra.Command, args []string) {
		dry := cmd.Flag("dry").Value.String() == "true"
		ctx := verman.BuildContext(dry)

		if ctx.SemverSource != core.SourceGit {
			fmt.Println("not a git repository.")
			return
		}

		var tagToRemove []string
		if len(args) > 0 {
			tagToRemove = args
		} else {
			tagToRemove = []string{ctx.CurrentVersion.String()}
		}

		fmt.Println("untagging...", tagToRemove)

		if !ctx.DryRun {
			err := verman.UntagVersions(tagToRemove, untagRemote)
			if err != nil {
				fmt.Println("error: untagging versions.", err)
				return
			}
		}

		fmt.Println("done.")

	},
}

func init() {
	untagCmd.Flags().BoolVarP(&untagRemote, "remote", "r", false, "remove remote tag as well")
	untagCmd.Flags().BoolP("dry", "d", false, "dry run")
	rootCmd.AddCommand(untagCmd)
}
