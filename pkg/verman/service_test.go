package verman_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/semver/pkg/commands"
	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/adharshmk96/semver/testutils"
	"github.com/stretchr/testify/assert"
)

func TestWriteConfigToFile(t *testing.T) {
	t.Run("write config to file", func(t *testing.T) {
		semver := verman.Semver{
			Major: 1,
			Minor: 1,
		}

		err := verman.WriteVersionToFile(&semver)

		assert.NoError(t, err)

		content, err := os.ReadFile(".version")
		assert.NoError(t, err)
		assert.Equal(t, "v1.1.0", string(content))
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		semver := verman.Semver{
			Major: 1,
		}

		err := verman.WriteVersionToFile(&semver)

		assert.NoError(t, err)

		content, err := os.ReadFile(".version")
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", string(content))

		semver = verman.Semver{
			Major: 2,
		}
		err = verman.WriteVersionToFile(&semver)

		assert.NoError(t, err)

		content, err = os.ReadFile(".version")
		assert.NoError(t, err)
		assert.Equal(t, "v2.0.0", string(content))

	})
}

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

	t.Run("initialize semver in git repo with tags", func(t *testing.T) {
		testDir := testutils.SetupTempDir(t)
		defer testDir()

		gitCmd := commands.NewGitCmd(commands.NewGitExec())
		gitCmd.Run("init")
		gitCmd.Run("tag", "v1.0.0")

		ctx := verman.BuildContext([]string{}, false)

		err := verman.InitializeSemver(ctx)
		assert.NoError(t, err)

		assert.NoFileExists(t, ".version")

		latestTag, err := gitCmd.LatestTag()
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", latestTag)
	})
}
