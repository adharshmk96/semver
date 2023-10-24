/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/adharshmk96/semver/pkg/verman/core"
	"github.com/spf13/cobra"
)

func createReleaseCommand(versionType string) *cobra.Command {
	return &cobra.Command{
		Use:   versionType,
		Short: fmt.Sprintf("Increment the %s version by one", versionType),
		Run: func(cmd *cobra.Command, args []string) {
			dry, _ := cmd.Flags().GetBool("dry")

			alpha, _ := cmd.Flags().GetBool("alpha")
			beta, _ := cmd.Flags().GetBool("beta")
			rc, _ := cmd.Flags().GetBool("rc")

			ctx := verman.BuildContext(dry)
			if ctx.SemverSource == core.SourceNone {
				fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
				return
			}

			verman.UpdateAndCommitVersion(ctx, versionType)

			if rc {
				ctx.CurrentVersion.IncrementRC()
			} else if beta {
				ctx.CurrentVersion.IncrementBeta()
			} else if alpha {
				ctx.CurrentVersion.IncrementAlpha()
			}

			fmt.Println(ctx.CurrentVersion.String())
		},
	}
}

func createPreReleaseCommand(versionType string) *cobra.Command {
	return &cobra.Command{
		Use:   versionType,
		Short: fmt.Sprintf("Increment the %s version by one", versionType),
		Run: func(cmd *cobra.Command, args []string) {
			dry, _ := cmd.Flags().GetBool("dry")

			ctx := verman.BuildContext(dry)
			if ctx.SemverSource == core.SourceNone {
				fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
				return
			}

			if ctx.CurrentVersion.IsRelease() {
				fmt.Println("current veresion is not a pre-release. run `semver ( major | minor | patch ) --( alpha | beta | rc )` to create a pre-release.")
				fmt.Println("hint: you can't go back to pre-release from existing release version. not allowed to perform a pre-release action on a release version.")
				return
			}

			verman.UpdateAndCommitVersion(ctx, versionType)

			fmt.Println(ctx.CurrentVersion.String())
		},
	}
}

func setReleaseCmdFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("dry", "d", false, "dry run mode")
	cmd.Flags().Bool("alpha", false, "increment alpha version")
	cmd.Flags().Bool("beta", false, "increment beta version")
	cmd.Flags().Bool("rc", false, "increment rc version")
}

func init() {
	majorCmd := createReleaseCommand("major")
	minorCmd := createReleaseCommand("minor")
	patchCmd := createReleaseCommand("patch")

	setReleaseCmdFlags(majorCmd)
	setReleaseCmdFlags(minorCmd)
	setReleaseCmdFlags(patchCmd)

	rootCmd.AddCommand(majorCmd, minorCmd, patchCmd)

	alphaCmd := createPreReleaseCommand("alpha")
	alphaCmd.Flags().BoolP("dry", "d", false, "dry run mode")

	betaCmd := createPreReleaseCommand("beta")
	betaCmd.Flags().BoolP("dry", "d", false, "dry run mode")

	rcCmd := createPreReleaseCommand("rc")
	rcCmd.Flags().BoolP("dry", "d", false, "dry run mode")

	rootCmd.AddCommand(alphaCmd, betaCmd, rcCmd)
}
