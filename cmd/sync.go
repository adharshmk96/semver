/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Display the current version of the project",
	Run: func(cmd *cobra.Command, args []string) {
		err := verman.FetchTags()
		if err != nil {
			fmt.Println("error fetching tags:", err)
			return
		}
	},
}

func init() {
	syncCmd.Flags().BoolP("source", "s", false, "display source info")
	rootCmd.AddCommand(syncCmd)
}
