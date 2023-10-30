package core

import "fmt"

var (
	ErrGitTagNotFound = fmt.Errorf("no git tags found")
	ErrGettingGitTag  = fmt.Errorf("error getting last git tag")
	ErrCreatingGitTag = fmt.Errorf("error creating git tag")
	ErrNotGitRepo     = fmt.Errorf("not a git repository")

	ErrInvalidVersionFormat = fmt.Errorf("invalid version format")
)
