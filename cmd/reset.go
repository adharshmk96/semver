/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "(CAUTION) Reset all tags and remove the semver configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("resetting semver configuration...")
		err := verman.RemoveConfig()
		if err != nil {
			fmt.Println("error removing configuration file.")
		}

		if !verman.IsGitRepository() {
			fmt.Println("not a git repository.")
			return
		}

		fmt.Println("removing all git tags...")
		err = verman.GitRemoveAllLocalTags()
		if err != nil {
			fmt.Println("error removing git tags.", err)
			return
		}

		fmt.Println("done. run `semver init` to initialize again...")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
