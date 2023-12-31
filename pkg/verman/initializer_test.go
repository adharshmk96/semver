package verman_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/adharshmk96/semver/pkg/verman/core"
	"github.com/adharshmk96/semver/testutils"
	"github.com/stretchr/testify/assert"
)

func repoWithTag(t *testing.T, tag string) {
	t.Helper()

	assert.NoError(t, exec.Command("git", "init").Run(), "git init failed")
	assert.DirExists(t, ".git", "git init failed")

	assert.NoError(t, exec.Command("git", "config", "user.email", "you@example.com").Run())
	assert.NoError(t, exec.Command("git", "config", "user.name", "Your Name").Run())
	assert.NoError(t, os.WriteFile("test.txt", []byte("test"), 0644), "failed to create test file")
	assert.NoError(t, exec.Command("git", "add", ".").Run(), "failed to add test file")
	assert.NoError(t, exec.Command("git", "commit", "-m", "initial commit").Run(), "failed to commit test file")
	assert.NoError(t, exec.Command("git", "tag", tag).Run(), "failed to tag git repo")
}

func TestBuildContext(t *testing.T) {
	t.Run("empty dir.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		ctx := verman.BuildContext(false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("untagged git repo.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		exec.Command("git", "init").Run() //nolint:gosec // This is a test and we need to run git commands.

		ctx := verman.BuildContext(false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("invalid tagged git repo.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag(t, "invalid")

		ctx := verman.BuildContext(false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("valid version file.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
		assert.NoError(t, err)

		ctx := verman.BuildContext(false)

		ctxVersion := ctx.CurrentVersion

		assert.Equal(t, 1, ctxVersion.Major)
		assert.Equal(t, 0, ctxVersion.Minor)
		assert.Equal(t, 0, ctxVersion.Patch)
		assert.Equal(t, 0, ctxVersion.Alpha)
		assert.Equal(t, 0, ctxVersion.Beta)
		assert.Equal(t, 1, ctxVersion.RC)

		assert.Equal(t, core.SourceFile, ctx.SemverSource)

	})

	t.Run("invalid version file", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		err := os.WriteFile(".version", []byte("invalid"), 0644)
		assert.NoError(t, err)

		ctx := verman.BuildContext(false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("invalid git tag and valid .version file", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag(t, "invalid")

		err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
		assert.NoError(t, err)

		ctx := verman.BuildContext(false)

		ctxVersion := ctx.CurrentVersion

		assert.Empty(t, ctxVersion)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
	})

	t.Run("valid tagged git repo.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag(t, "v1.0.0")

		ctx := verman.BuildContext(false)

		assert.Equal(t, core.SourceGit, ctx.SemverSource)
		assert.Equal(t, 1, ctx.CurrentVersion.Major)
		assert.Equal(t, 0, ctx.CurrentVersion.Minor)
		assert.Equal(t, 0, ctx.CurrentVersion.Patch)
		assert.Equal(t, 0, ctx.CurrentVersion.Alpha)
		assert.Equal(t, 0, ctx.CurrentVersion.Beta)
		assert.Equal(t, 0, ctx.CurrentVersion.RC)
	})

	t.Run("valid git tag and invalid .version file", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag(t, "v1.0.0")

		err := os.WriteFile(".version", []byte("invalid"), 0644)
		assert.NoError(t, err)

		ctx := verman.BuildContext(false)

		ctxVersion := ctx.CurrentVersion

		assert.Equal(t, 1, ctxVersion.Major)
		assert.Equal(t, 0, ctxVersion.Minor)
		assert.Equal(t, 0, ctxVersion.Patch)
		assert.Equal(t, 0, ctxVersion.Alpha)
		assert.Equal(t, 0, ctxVersion.Beta)
		assert.Equal(t, 0, ctxVersion.RC)

		assert.Equal(t, core.SourceGit, ctx.SemverSource)
	})

	t.Run("valid git tag and valid .version file", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag(t, "v1.0.0")

		err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
		assert.NoError(t, err)

		ctx := verman.BuildContext(false)

		ctxVersion := ctx.CurrentVersion

		assert.Equal(t, 1, ctxVersion.Major)
		assert.Equal(t, 0, ctxVersion.Minor)
		assert.Equal(t, 0, ctxVersion.Patch)
		assert.Equal(t, 0, ctxVersion.Alpha)
		assert.Equal(t, 0, ctxVersion.Beta)
		assert.Equal(t, 0, ctxVersion.RC)

		assert.Equal(t, core.SourceGit, ctx.SemverSource)
	})

}
