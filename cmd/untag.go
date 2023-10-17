/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

var untagRemote bool

func getVersionToRemove(args []string) string {
	if len(args) == 0 {
		version, err := verman.GetVersionFromConfig()
		if err != nil {
			fmt.Println("error reading configuration file.")
			return ""
		}
		return version.String()
	}
	return args[0]
}

var untagCmd = &cobra.Command{
	Use:   "untag",
	Short: "Delete a specific tag from git (default: current tag)",
	Run: func(cmd *cobra.Command, args []string) {
		versionToRemove := getVersionToRemove(args)
		fmt.Println("untagging...", versionToRemove)

		if !verman.IsGitRepository() {
			fmt.Println("not a git repository.")
			return
		}

		if untagRemote {
			err := verman.GitRemoveRemoteTag(versionToRemove)
			if err != nil {
				fmt.Println("error removing remote git tag.", err)
				return
			}
		}

		err := verman.GitRemoveLocalTag(versionToRemove)
		if err != nil {
			fmt.Println("error removing git tag.", err)
			return
		}

		tag, err := verman.GetVersionFromGitTag()
		if err != nil {
			fmt.Println("error getting latest git tag.", err)
			tag = &verman.Semver{Patch: 1}
		}

		verman.WriteVersionToConfig(tag)
		fmt.Println("current version:", tag.String())

	},
}

func init() {
	untagCmd.Flags().BoolVarP(&untagRemote, "remote", "r", false, "remove remote tag as well")

	rootCmd.AddCommand(untagCmd)
}
