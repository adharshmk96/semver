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

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Git push the current version of the project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := verman.BuildContext(false)

		if ctx.SemverSource == core.SourceNone {
			fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
			return
		}

		if !ctx.IsGitRepo {
			fmt.Println("not a git repository.")
			return
		}

		fmt.Println("pushing git tag:", ctx.CurrentVersion.String())
		err := verman.PushGitTag(ctx)
		if err != nil {
			fmt.Println("error pushing git tag:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
