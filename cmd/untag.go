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
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := verman.BuildContext(false)

		if ctx.SemverSource != core.SourceGit {
			fmt.Println("not a git repository.")
			return
		}

		var tagToRemove string
		if len(args) > 0 {
			tagToRemove = args[0]
		} else {
			tagToRemove = ctx.CurrentVersion.String()
		}

		fmt.Println("untagging...", tagToRemove)

		verman.UntagVersion(tagToRemove, untagRemote)

		fmt.Println("done.")

	},
}

func init() {
	untagCmd.Flags().BoolVarP(&untagRemote, "remote", "r", false, "remove remote tag as well")
	rootCmd.AddCommand(untagCmd)
}
