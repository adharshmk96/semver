package commands_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/stretchr/testify/assert"
)

type MockCmdExecutor struct {
	mockRunCmd func(args ...string) (string, error)
}

func (m *MockCmdExecutor) RunCmd(args ...string) (string, error) {
	return m.mockRunCmd(args...)
}

func tempWorkDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, _ := os.MkdirTemp("", "test")

	os.Chdir(tempDir)
	cleanUp := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanUp
}

func TestGitCommands(t *testing.T) {
	t.Run("test isRepo on real directory", func(t *testing.T) {
		tempDir, cleanUp := tempWorkDir(t)

		defer cleanUp()

		gitCmd := commands.NewGitCmd(commands.NewGitExec())

		_, err := gitCmd.Run("init")
		assert.NoError(t, err)

		fmt.Println(tempDir)
		assert.True(t, gitCmd.IsRepo())

	})

	t.Run("test isRepo returns error", func(t *testing.T) {
		mockExec := MockCmdExecutor{
			mockRunCmd: func(args ...string) (string, error) {
				return "", fmt.Errorf("error")
			},
		}

		gitCmd := commands.NewGitCmd(&mockExec)

		assert.False(t, gitCmd.IsRepo())
	})
}
