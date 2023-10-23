package verman

type Context struct {
	WorkDir        string
	CurrentVersion *Semver

	DryRun bool

	IsGitRepo bool
	IsVerman  bool
}
