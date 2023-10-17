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
	if dry {
		fmt.Println("dry run...")
	}
	projectVersion, err := verman.GetVersionFromConfig()
	fmt.Println("current version:", projectVersion.String())
	if err != nil {
		fmt.Println("error reading configuration file.")
		return projectVersion, err
	}

	projectVersion.UpdateSemver(versionType)

	return projectVersion, nil

}

func commitUpdatedVersion(projectVersion *verman.Semver) error {

	if dry {
		fmt.Println("updated version:", projectVersion.String())
		return nil
	}

	err := verman.WriteVersionToConfig(projectVersion)
	if err != nil {
		fmt.Println("error writing to configuration file.")
		return err
	}

	if verman.IsGitRepository() {
		err = verman.GitCommitVersionConfig(projectVersion)
		if err != nil {
			fmt.Println("error committing configuration file.")
			fmt.Println(err)
			return err
		}

		err = verman.GitTagVersion(projectVersion)
		if err != nil {
			fmt.Println("error creating git tag: check if the tag already exists.")
			return err
		}
	}

	fmt.Println("updated version:", projectVersion.String())
	return nil
}

func getSubVersionFromArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	if args[0] == "alpha" || args[0] == "beta" || args[0] == "rc" {
		return args[0]
	}
	return ""
}

var majorCmd = &cobra.Command{
	Use:   "major",
	Short: "Increment the major version by one",
	Run: func(cmd *cobra.Command, args []string) {
		projectVersion, err := incrementVersion("major")
		if err != nil {
			fmt.Println("error incrementing version.")
		}

		subVersion := getSubVersionFromArgs(args)
		if subVersion != "" {
			projectVersion.UpdateSemver(subVersion)
		}

		commitUpdatedVersion(projectVersion)
	},
}

var minorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Increment the minor version by one",
	Run: func(cmd *cobra.Command, args []string) {
		projectVersion, err := incrementVersion("minor")
		if err != nil {
			fmt.Println("error incrementing version.")
		}

		subVersion := getSubVersionFromArgs(args)
		if subVersion != "" {
			projectVersion.UpdateSemver(subVersion)
		}

		commitUpdatedVersion(projectVersion)
	},
}

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Increment the patch version by one",
	Run: func(cmd *cobra.Command, args []string) {
		projectVersion, err := incrementVersion("patch")
		if err != nil {
			fmt.Println("error incrementing version.")
		}

		subVersion := getSubVersionFromArgs(args)
		if subVersion != "" {
			projectVersion.UpdateSemver(subVersion)
		}

		commitUpdatedVersion(projectVersion)
	},
}

var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "Increment the alpha version by one",
	Run: func(cmd *cobra.Command, args []string) {
		projectVersion, err := incrementVersion("alpha")
		if err != nil {
			fmt.Println("error incrementing version.")
		}
		commitUpdatedVersion(projectVersion)
	},
}

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Increment the beta version by one",
	Run: func(cmd *cobra.Command, args []string) {
		projectVersion, err := incrementVersion("beta")
		if err != nil {
			fmt.Println("error incrementing version.")
		}
		commitUpdatedVersion(projectVersion)
	},
}

var rcCmd = &cobra.Command{
	Use:   "rc",
	Short: "Increment the rc version by one",
	Run: func(cmd *cobra.Command, args []string) {
		projectVersion, err := incrementVersion("rc")
		if err != nil {
			fmt.Println("error incrementing version.")
		}
		commitUpdatedVersion(projectVersion)
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
