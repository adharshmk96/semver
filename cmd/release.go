/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

// releaseCmd represents the get command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "removes the pre-release and tags the release version",
	Run: func(cmd *cobra.Command, args []string) {

		semver, err := verman.GetVersionFromConfig()
		if err != nil {
			fmt.Println("error reading configuration file.")
			return
		}

		isPreRelease := verman.IsPreRelease(semver)
		if !isPreRelease {
			fmt.Println("not a pre-release.")
			return
		}

		semver.Release()

		if verman.GitTagExists(semver.String()) {
			fmt.Println("tag already exists. please run `semver up` to increment the version.")
			return
		}

		err = verman.WriteVersionToConfig(semver)
		if err != nil {
			fmt.Println("error writing to configuration file.")
			return
		}

		commitUpdatedVersion(semver)

		fmt.Println(semver.String())
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
