package verman_test

import (
	"os"
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

		ctx := verman.BuildContext([]string{}, false)

		err := verman.InitializeSemver(ctx)
		assert.NoError(t, err)

		content, err := os.ReadFile(".version")
		assert.NoError(t, err)
		assert.Equal(t, "v0.0.1", string(content))
	})

	t.Run("initialize semver in git repo without tags", func(t *testing.T) {
		testDir := testutils.SetupTempDir(t)
		defer testDir()

		gitCmd := commands.NewGitCmd(commands.NewGitExec())
		gitCmd.Run("init")

		ctx := verman.BuildContext([]string{}, false)

		err := verman.InitializeSemver(ctx)
		assert.NoError(t, err)

		assert.NoFileExists(t, ".version")

		latestTag, err := gitCmd.LatestTag()
		assert.NoError(t, err)
		assert.Equal(t, "v0.0.1", latestTag)

	})

}
