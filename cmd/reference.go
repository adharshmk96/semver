/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/cobra"
)

// referenceCmd represents the get command
var referenceCmd = &cobra.Command{
	Use:   "refs",
	Short: "Display references of the current version in the project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := verman.BuildContext(false)

		fmt.Println(ctx.CurrentVersion.String())
		fmt.Println("References:")
		res, err := verman.VerifyTagReferences(ctx)
		if err != nil {
			fmt.Println("error verifying tag references:", err)
			return
		}

		fmt.Println(res)

	},
}

func init() {
	referenceCmd.Flags().BoolP("source", "s", false, "display source info")
	rootCmd.AddCommand(referenceCmd)
}
