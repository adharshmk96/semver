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

var resetRemote bool

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "(CAUTION) Reset all tags and remove the semver configuration",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := verman.BuildContext(false)

		if ctx.SemverSource == core.SourceNone {
			fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
			return
		}

		fmt.Println("resetting semver configuration...")
		err := verman.ResetSemver(ctx, resetRemote)
		if err != nil {
			fmt.Println("error: resetting semver.", err)
			return
		}

		var initVersion string
		if len(args) > 0 {
			initVersion = args[0]
		} else {
			fmt.Println("done. run `semver init` to initialize again...")
			return
		}

		fmt.Println("re-initializing semver configuration...")
		err = verman.InitializeSemver(ctx, initVersion)
		if err != nil {
			fmt.Println("error: initalizing semver.", err)
			return
		}

		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")

	},
}

func init() {
	resetCmd.Flags().BoolVarP(&resetRemote, "remote", "r", false, "remove remote tags as well")

	rootCmd.AddCommand(resetCmd)
}
