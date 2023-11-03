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

var (
	dry  bool
	sync bool
	push bool
)

func createReleaseCommand(versionType string) *cobra.Command {
	return &cobra.Command{
		Use:   versionType,
		Short: fmt.Sprintf("Increment the %s version by one", versionType),
		Run: func(cmd *cobra.Command, args []string) {

			alpha, _ := cmd.Flags().GetBool("alpha")
			beta, _ := cmd.Flags().GetBool("beta")
			rc, _ := cmd.Flags().GetBool("rc")

			if sync {
				fmt.Println("Fetching remote...")
				verman.FetchTags()
			}

			ctx := verman.BuildContext(dry)
			if ctx.SemverSource == core.SourceNone {
				fmt.Println("semver config not found. run `semver init` to initialize the semver configuration.")
				return
			}

			ctx.CurrentVersion.UpdateSemver(versionType)

			if rc {
				ctx.CurrentVersion.IncrementRC()
			} else if beta {
				ctx.CurrentVersion.IncrementBeta()
			} else if alpha {
				ctx.CurrentVersion.IncrementAlpha()
			}

			if ctx.DryRun {
				fmt.Println(ctx.CurrentVersion.String())
				return
			}

			verman.CommitVersionLocally(ctx)
			if push {
				fmt.Println("pushing git tag:", ctx.CurrentVersion.String())
				err := verman.PushGitTag(ctx)
				if err != nil {
					fmt.Println("error pushing git tag:", err)
					return
				}
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
			if sync {
				fmt.Println("Fetching remote...")
				verman.FetchTags()
			}

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

			ctx.CurrentVersion.UpdateSemver(versionType)

			if ctx.DryRun {
				fmt.Println(ctx.CurrentVersion.String())
				return
			}

			verman.CommitVersionLocally(ctx)

			if push {
				fmt.Println("pushing git tag:", ctx.CurrentVersion.String())
				err := verman.PushGitTag(ctx)
				if err != nil {
					fmt.Println("error pushing git tag:", err)
					return
				}
			}

			fmt.Println(ctx.CurrentVersion.String())
		},
	}
}

func setReleaseCmdFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&dry, "dry", "d", false, "dry run mode")
	cmd.Flags().BoolVarP(&push, "push", "p", false, "push the git tag")
	cmd.Flags().BoolVar(&sync, "sync", false, "sync with remote")

	cmd.Flags().Bool("alpha", false, "increment alpha version")
	cmd.Flags().Bool("beta", false, "increment beta version")
	cmd.Flags().Bool("rc", false, "increment rc version")

}

func setPreReleaseCmdFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&dry, "dry", "d", false, "dry run mode")
	cmd.Flags().BoolVarP(&push, "push", "p", false, "push the git tag")
	cmd.Flags().BoolVar(&sync, "sync", false, "sync with remote")
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
	betaCmd := createPreReleaseCommand("beta")
	rcCmd := createPreReleaseCommand("rc")

	setPreReleaseCmdFlags(alphaCmd)
	setPreReleaseCmdFlags(betaCmd)
	setPreReleaseCmdFlags(rcCmd)

	rootCmd.AddCommand(alphaCmd, betaCmd, rcCmd)
}
