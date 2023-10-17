/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

var resetRemote bool

// logAndExecute logs the given message, then executes the given function
func logAndExecute(message string, action func() error) {
	fmt.Println(message)
	if err := action(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "(CAUTION) Reset all tags and remove the semver configuration",
	Run: func(cmd *cobra.Command, args []string) {
		logAndExecute("resetting semver configuration...", verman.RemoveConfig)

		if !verman.IsGitRepository() {
			fmt.Println("not a git repository.")
			return
		}

		if resetRemote {
			logAndExecute("removing all remote git tags...", verman.GitRemoveAllRemoteTags)
		}

		logAndExecute("removing all local git tags...", verman.GitRemoveAllLocalTags)
		fmt.Println("done. run `semver init` to initialize again...")
	},
}

func init() {
	resetCmd.Flags().BoolVarP(&resetRemote, "remote", "r", false, "remove remote tags as well")

	rootCmd.AddCommand(resetCmd)
}
