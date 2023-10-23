package testutils

import (
	"os"
	"testing"
)

func SetupTempDir(t *testing.T) func() {
	tempDir, _ := os.MkdirTemp("", "test")
	os.Chdir(tempDir)

	return func() {
		os.RemoveAll(tempDir)
	}

}
