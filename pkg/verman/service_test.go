package verman_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/adharshmk96/semver/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInitializeSemver(t *testing.T) {
	t.Run("initialize semver in empty folder without git", func(t *testing.T) {
		testDir := testutils.SetupTempDir(t)
		defer testDir()

		ctx := verman.BuildContext(false)

		err := verman.InitializeSemver(ctx, "v0.0.1")
		assert.NoError(t, err)

		content, err := os.ReadFile(".version")
		assert.NoError(t, err)
		assert.Equal(t, "v0.0.1", string(content))
	})

	t.Run("initialize semver in git repo without tags", func(t *testing.T) {
		testDir := testutils.SetupTempDir(t)
		defer testDir()

		gitCmd := commands.NewGitCmd(commands.NewGitExec())
		assert.NoError(t, exec.Command("git", "init").Run())
		assert.NoError(t, exec.Command("git", "config", "user.email", "user@email.com").Run())
		assert.NoError(t, exec.Command("git", "config", "user.name", "user").Run())
		assert.NoError(t, os.WriteFile("test.txt", []byte("test"), 0644))
		assert.NoError(t, exec.Command("git", "add", ".").Run())
		assert.NoError(t, exec.Command("git", "commit", "-m", "initial commit").Run())

		ctx := verman.BuildContext(false)

		err := verman.InitializeSemver(ctx, "v0.0.1")
		assert.NoError(t, err)

		assert.NoFileExists(t, ".version")

		latestTag, err := gitCmd.LatestTag()
		assert.NoError(t, err)
		assert.Equal(t, "v0.0.1", latestTag)

	})

}
