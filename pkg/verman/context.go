package verman

import "errors"

type Source int

const (
	SourceNone Source = iota
	SourceGit
	SourceFile
)

type Context struct {
	WorkDir        string
	CurrentVersion *Semver

	DryRun bool

	IsGitRepo    bool
	SemverSource Source
}

const (
	VERSION_FILE    = ".version"
	INITIAL_VERSION = "v0.0.1"
)

var (
	ErrParsingSemver = errors.New("error parsing semver")
)
