/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display the current version of the project",
	Long:  `Reads the version.yaml file and displays the current version of the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !configExists {
			fmt.Println("No configuration found. Please run `semver init` to initialize the configuration.")
			return
		}

		version, err := verman.GetVersionFromConfig()
		if err != nil {
			fmt.Println("Error reading configuration file.")
			return
		}

		fmt.Println("Project version: ", version.String())
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
