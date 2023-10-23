package verman

import (
	"fmt"
	"os"

	"github.com/adharshmk96/semver/pkg/verman/core"
)

func writeToFile(filePath string, fileContent string) error {
	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// New stuff

func DisplaySource(ct *core.Context) {
	switch ct.SemverSource {
	case core.SourceNone:
		fmt.Println("no version source found.")
	case core.SourceGit:
		fmt.Println("version source: git tag.")
	case core.SourceFile:
		fmt.Println("version source: .version file.")
	}
}
