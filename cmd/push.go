/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Git push the current version of the project",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := verman.GetVersionFromConfig()
		if err != nil {
			fmt.Println("error reading configuration file.")
			return
		}

		fmt.Println("pushing git tag:", version.String())
		err = verman.GitPushTag(version)
		if err != nil {
			fmt.Println("error pushing git tag.")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
