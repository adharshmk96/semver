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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display the current version of the project",
	Long:  `Reads the .version.yaml file and displays the current version of the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := verman.BuildContext(false)

		showSource, _ := cmd.Flags().GetBool("source")

		if ctx.SemverSource == core.SourceNone {
			fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
			return
		}

		if showSource {
			fmt.Println(ctx.GetSourceInfo())
		}
		fmt.Println(ctx.CurrentVersion.String())
	},
}

func init() {
	getCmd.Flags().BoolP("source", "s", false, "display source info")
	rootCmd.AddCommand(getCmd)
}
