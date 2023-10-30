package core

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

func (c *Context) GetSourceInfo() string {
	switch c.SemverSource {
	case SourceNone:
		return "semver config not found. run `semver init` to initialize the semver configuration."
	case SourceGit:
		return "Source: Git"
	case SourceFile:
		return "Source: File"
	default:
		return ""
	}
}

const (
	VERSION_FILE    = ".version"
	INITIAL_VERSION = "v0.0.1"
)

var (
	ErrParsingSemver = errors.New("error parsing semver")
)
