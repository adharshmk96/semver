/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// SemVer is set at build time via ldflags.
var SemVer = "development"

const notInitialized = "semver config not found. run `semver init` to initialize the semver configuration."

var (
	dry    bool
	push   bool
	sync   bool
	remote bool
)

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Manage your project's semver configuration",
	Long: `A CLI tool to manage your project's semantic version.

semver uses git tags or .version file (for non-git projects) to manage the version of the project.

Running 'semver' with no arguments starts the local web UI (same as 'semver ui').
Use 'semver --help' to list the available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			if SemVer == "development" {
				if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
					SemVer = info.Main.Version
				}
			}
			fmt.Println(SemVer)
			return
		}
		if help, _ := cmd.Flags().GetBool("help"); help {
			cmd.Help()
			return
		}
		runUI(cmd)
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display the current version of the project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := BuildContext(false)
		if ctx.Source == SourceNone {
			fmt.Println(notInitialized)
			return
		}
		if source, _ := cmd.Flags().GetBool("source"); source {
			fmt.Println(ctx.Source)
		}
		fmt.Println(ctx.Version)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the semver configuration",
	Long: `Initialize the semver configuration. This will create a .version file in the current directory,
or tag the git repository, and set it as the current version of the project.

If no version is given, it defaults to v0.0.1.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := BuildContext(false)
		if ctx.Source != SourceNone {
			fmt.Println("semver config found, run `semver get` to view the version.")
			return
		}
		if err := ctx.Commit(initialVersion(args)); err != nil {
			fmt.Println("error: initializing semver.", err)
			return
		}
		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")
	},
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Removes the pre-release and tags the release version",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := buildContextWithSync(dry)
		if ctx.Source == SourceNone {
			fmt.Println(notInitialized)
			return
		}
		if !ctx.Version.IsPreRelease() {
			fmt.Println("current version is not a pre-release.")
			return
		}
		ctx.Version.Bump("release")
		commitAndPush(ctx)
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Git push the current version of the project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := BuildContext(false)
		if ctx.Source == SourceNone {
			fmt.Println(notInitialized)
			return
		}
		if err := ctx.Push(); err != nil {
			fmt.Println("error pushing git tag:", err)
		}
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetch tags from the remote repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := FetchTags(); err != nil {
			fmt.Println("error fetching tags:", err)
		}
	},
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start a local web UI to manage the project's tags",
	Run: func(cmd *cobra.Command, args []string) {
		runUI(cmd)
	},
}

// runUI starts the local web UI using the port/no-open flags of cmd.
func runUI(cmd *cobra.Command) {
	port, _ := cmd.Flags().GetString("port")
	noOpen, _ := cmd.Flags().GetBool("no-open")

	BuildContext(false) // move to the repository root
	if err := Serve("127.0.0.1:"+port, !noOpen); err != nil {
		fmt.Println("error starting ui:", err)
		os.Exit(1)
	}
}

var refsCmd = &cobra.Command{
	Use:   "refs",
	Short: "Display references of the current version in the project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := BuildContext(false)
		if ctx.Source == SourceNone {
			fmt.Println(notInitialized)
			return
		}
		fmt.Println(ctx.Version)
		refs, err := ctx.References()
		if err != nil {
			fmt.Println("no references found.")
			return
		}
		fmt.Println("References:")
		fmt.Println(refs)
	},
}

var untagCmd = &cobra.Command{
	Use:   "untag",
	Short: "Delete a specific tag from git (default: current tag)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := buildContextWithSync(dry)
		if ctx.Source != SourceGit {
			fmt.Println("not a git repository.")
			return
		}

		tags := args
		if len(tags) == 0 {
			tags = []string{ctx.Version.String()}
		}
		fmt.Println("untagging...", tags)

		if !ctx.DryRun {
			if err := Untag(tags, remote); err != nil {
				fmt.Println("error: untagging versions.", err)
				return
			}
		}
		fmt.Println("done.")
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "(CAUTION) Reset all tags and remove the semver configuration",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := buildContextWithSync(false)
		if ctx.Source == SourceNone {
			fmt.Println(notInitialized)
			return
		}

		fmt.Println("resetting semver configuration...")
		if err := ctx.Reset(remote); err != nil {
			fmt.Println("error: resetting semver.", err)
			return
		}

		if len(args) == 0 {
			fmt.Println("done. run `semver init` to initialize again...")
			return
		}

		fmt.Println("re-initializing semver configuration...")
		if err := ctx.Commit(args[0]); err != nil {
			fmt.Println("error: initializing semver.", err)
			return
		}
		fmt.Println("semver configuration initialized successfully. run `semver get` to display the current version.")
	},
}

// bumpCmd increments part (major, minor, patch, alpha, beta or rc) of the
// current version. Release parts additionally accept --alpha/--beta/--rc to
// start a pre-release of the new version.
func bumpCmd(part string) *cobra.Command {
	isRelease := part == "major" || part == "minor" || part == "patch"

	cmd := &cobra.Command{
		Use:   part,
		Short: fmt.Sprintf("Increment the %s version by one", part),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := buildContextWithSync(dry)
			if ctx.Source == SourceNone {
				fmt.Println(notInitialized)
				return
			}
			if !isRelease && !ctx.Version.IsPreRelease() {
				fmt.Println("current version is not a pre-release. run `semver ( major | minor | patch ) --( alpha | beta | rc )` to create a pre-release.")
				fmt.Println("hint: you can't go back to pre-release from an existing release version.")
				return
			}

			ctx.Version.Bump(part)
			if isRelease {
				for _, label := range preLabels {
					if pre, _ := cmd.Flags().GetBool(label); pre {
						ctx.Version.Bump(label)
						break
					}
				}
			}

			commitAndPush(ctx)
		},
	}

	setCommonFlags(cmd)
	if isRelease {
		for _, label := range preLabels {
			cmd.Flags().Bool(label, false, fmt.Sprintf("start a %s pre-release", label))
		}
	}

	return cmd
}

// commitAndPush records the updated version, unless this is a dry run.
func commitAndPush(ctx *Context) {
	if ctx.DryRun {
		fmt.Println(ctx.Version)
		return
	}
	if err := ctx.Commit(ctx.Version.String()); err != nil {
		fmt.Println("error: updating version.", err)
		return
	}
	if push {
		if err := ctx.Push(); err != nil {
			fmt.Println("error pushing git tag:", err)
			return
		}
	}
	fmt.Println(ctx.Version)
}

func buildContextWithSync(dry bool) *Context {
	if sync {
		fmt.Println("fetching remote...")
		FetchTags()
	}
	return BuildContext(dry)
}

func initialVersion(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "v0.0.1"
}

func setCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&dry, "dry", "d", false, "dry run mode")
	cmd.Flags().BoolVarP(&push, "push", "p", false, "push the git tag")
	cmd.Flags().BoolVar(&sync, "sync", false, "fetch remote tags first")
}

func main() {
	rootCmd.Flags().BoolP("version", "v", false, "display current version")
	getCmd.Flags().BoolP("source", "s", false, "display source info")

	setCommonFlags(releaseCmd)

	untagCmd.Flags().BoolVarP(&remote, "remote", "r", false, "remove the remote tag as well")
	untagCmd.Flags().BoolVarP(&dry, "dry", "d", false, "dry run mode")
	untagCmd.Flags().BoolVar(&sync, "sync", false, "fetch remote tags first")

	resetCmd.Flags().BoolVarP(&remote, "remote", "r", false, "remove remote tags as well")
	resetCmd.Flags().BoolVar(&sync, "sync", false, "fetch remote tags first")

	for _, cmd := range []*cobra.Command{rootCmd, uiCmd} {
		cmd.Flags().StringP("port", "P", "7420", "port to serve the ui on")
		cmd.Flags().Bool("no-open", false, "do not open the browser")
	}

	rootCmd.AddCommand(getCmd, initCmd, releaseCmd, pushCmd, syncCmd, refsCmd, untagCmd, resetCmd, uiCmd)
	for _, part := range []string{"major", "minor", "patch", "alpha", "beta", "rc"} {
		rootCmd.AddCommand(bumpCmd(part))
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
