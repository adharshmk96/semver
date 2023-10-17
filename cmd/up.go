/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

var dry bool
var writeVersionConst bool

func incrementVersion(versionType string) (*verman.Semver, error) {
	printIfDryRun("dry run...")
	projectVersion, err := verman.GetVersionFromConfig()
	if err != nil {
		fmt.Println("error reading configuration file.")
		return projectVersion, err
	}
	fmt.Println("current version:", projectVersion.String())
	projectVersion.UpdateSemver(versionType)
	return projectVersion, nil
}

func printIfDryRun(msg string) {
	if dry {
		fmt.Println(msg)
	}
}

func commitAndPrintVersion(version *verman.Semver) {
	err := commitUpdatedVersion(version)
	if err != nil {
		return
	}
	fmt.Println("updated version:", version.String())
}

func commitUpdatedVersion(version *verman.Semver) error {
	printIfDryRun(fmt.Sprintf("updated version: %s", version.String()))
	if dry {
		return nil
	}

	err := verman.WriteVersionToConfig(version)
	if err != nil {
		fmt.Println("error writing to configuration file.")
		return err
	}

	if !verman.IsGitRepository() {
		return nil
	}

	err = verman.GitCommitVersionConfig(version)
	if err != nil {
		fmt.Println("error committing configuration file.")
		return err
	}

	err = verman.GitTagVersion(version)
	if err != nil {
		fmt.Println("error creating git tag: check if the tag already exists.")
		return err
	}

	return nil
}

func getSubVersionFromArgs(args []string) string {
	if len(args) > 0 {
		switch args[0] {
		case "alpha", "beta", "rc":
			return args[0]
		}
	}
	return ""
}

func updateAndCommitVersion(versionType string, args []string) {
	projectVersion, err := incrementVersion(versionType)
	if err != nil {
		fmt.Println("error incrementing version.")
		return
	}

	subVersion := getSubVersionFromArgs(args)
	if subVersion != "" {
		projectVersion.UpdateSemver(subVersion)
	}

	commitAndPrintVersion(projectVersion)
}

var versionCmds = []*cobra.Command{
	{
		Use:   "major",
		Short: "Increment the major version by one",
		Run: func(cmd *cobra.Command, args []string) {
			updateAndCommitVersion("major", args)
		},
	},
	{
		Use:   "minor",
		Short: "Increment the minor version by one",
		Run: func(cmd *cobra.Command, args []string) {
			updateAndCommitVersion("minor", args)
		},
	},
	{
		Use:   "patch",
		Short: "Increment the patch version by one",
		Run: func(cmd *cobra.Command, args []string) {
			updateAndCommitVersion("patch", args)
		},
	},
	{
		Use:   "alpha",
		Short: "Increment the alpha version by one",
		Run: func(cmd *cobra.Command, args []string) {
			updateAndCommitVersion("alpha", args)
		},
	},
	{
		Use:   "beta",
		Short: "Increment the beta version by one",
		Run: func(cmd *cobra.Command, args []string) {
			updateAndCommitVersion("beta", args)
		},
	},
	{
		Use:   "rc",
		Short: "Increment the rc version by one",
		Run: func(cmd *cobra.Command, args []string) {
			updateAndCommitVersion("rc", args)
		},
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Increment the version by one",
	Args:  cobra.MaximumNArgs(2),
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.PersistentFlags().BoolVarP(&dry, "dry", "d", false, "dry run mode")
	upCmd.PersistentFlags().BoolVarP(&writeVersionConst, "write-version", "w", false, "write the version to cmd/version_constant.go file")

	for _, cmd := range versionCmds {
		upCmd.AddCommand(cmd)
	}
}
