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

func incrementVersion(versionType string) error {
	if dry {
		fmt.Println("dry run...")
	}
	projectVersion, err := verman.GetVersionFromConfig()
	fmt.Println("current version:", projectVersion.String())
	if err != nil {
		fmt.Println("error reading configuration file.")
		return err
	}

	switch versionType {
	case "major":
		projectVersion.IncrementMajor()
	case "minor":
		projectVersion.IncrementMinor()
	case "patch":
		projectVersion.IncrementPatch()
	case "alpha":
		projectVersion.IncrementAlpha()
	case "beta":
		projectVersion.IncrementBeta()
	case "rc":
		projectVersion.IncrementRC()
	default:
		fmt.Println("invalid version type")
		return fmt.Errorf("invalid version type")
	}

	if !dry {

		err = verman.WriteVersionToConfig(projectVersion)
		if err != nil {
			fmt.Println("error writing to configuration file.")
			return err
		}

		if verman.IsGitRepository() {
			err = verman.GitCommitVersionConfig(projectVersion)
			if err != nil {
				fmt.Println("error committing configuration file.")
				return err
			}

			err = verman.GitTagVersion(projectVersion)
			if err != nil {
				fmt.Println("error creating git tag: check if the tag already exists.")
				return err
			}
		}
	}

	fmt.Println("updated version:", projectVersion.String())
	return nil
}

var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Increment the major version by one",
	Run: func(cmd *cobra.Command, args []string) {
		incrementVersion("major")
	},
}

var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Increment the minor version by one",
	Run: func(cmd *cobra.Command, args []string) {
		incrementVersion("minor")
	},
}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Increment the patch version by one",
	Run: func(cmd *cobra.Command, args []string) {
		incrementVersion("patch")
	},
}

var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "Increment the alpha version by one",
	Run: func(cmd *cobra.Command, args []string) {
		incrementVersion("alpha")
	},
}

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Increment the beta version by one",
	Run: func(cmd *cobra.Command, args []string) {
		incrementVersion("beta")
	},
}

var rcCmd = &cobra.Command{
	Use:   "rc",
	Short: "Increment the rc version by one",
	Run: func(cmd *cobra.Command, args []string) {
		incrementVersion("rc")
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Increment the semver version by one",
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.PersistentFlags().BoolVarP(&dry, "dry", "d", false, "dry run mode")
	upCmd.PersistentFlags().BoolVarP(&writeVersionConst, "write-version", "w", false, "write the version to cmd/version_constant.go file")

	upCmd.AddCommand(majorCmd)
	upCmd.AddCommand(minorCmd)
	upCmd.AddCommand(patchCmd)
	upCmd.AddCommand(alphaCmd)
	upCmd.AddCommand(betaCmd)
	upCmd.AddCommand(rcCmd)
}
