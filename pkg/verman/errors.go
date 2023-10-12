package verman

import "fmt"

var (
	ErrGitTagNotFound = fmt.Errorf("no git tags found")

	ErrInvalidVersionFormat = fmt.Errorf("invalid version format")

	ErrCreatingGitTag = fmt.Errorf("error creating git tag")
)
