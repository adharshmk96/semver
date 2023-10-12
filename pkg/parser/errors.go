package parser

import "fmt"

var (
	ErrGitTagNotFound = fmt.Errorf("no git tags found")

	ErrInvalidVersionFormat = fmt.Errorf("invalid version format")
)
