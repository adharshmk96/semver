/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var SemVer = "v0.0.0"

func GetSemverInfo() string {
	if SemVer != "v0.0.0" {
		return SemVer
	}
	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != "(devel)" && version.Main.Version != "" {
		return version.Main.Version
	}
	return SemVer
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of semver",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetSemverInfo())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
