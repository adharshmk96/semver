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

func repoWithTag(tag string) {
	exec.Command("git", "init").Run() //nolint:gosec // This is a test and we need to run git commands.
	os.WriteFile("test.txt", []byte("test"), 0644)
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "initial commit").Run()
	exec.Command("git", "tag", tag).Run()
}

func TestBuildContext(t *testing.T) {
	t.Run("empty dir.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		args := []string{}
		ctx := verman.BuildContext(args, false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("untagged git repo.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		exec.Command("git", "init").Run() //nolint:gosec // This is a test and we need to run git commands.

		args := []string{}
		ctx := verman.BuildContext(args, false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("valid tagged git repo.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag("v1.0.0")

		args := []string{}
		ctx := verman.BuildContext(args, false)

		assert.Equal(t, core.SourceGit, ctx.SemverSource)
		assert.Equal(t, 1, ctx.CurrentVersion.Major)
		assert.Equal(t, 0, ctx.CurrentVersion.Minor)
		assert.Equal(t, 0, ctx.CurrentVersion.Patch)
		assert.Equal(t, 0, ctx.CurrentVersion.Alpha)
		assert.Equal(t, 0, ctx.CurrentVersion.Beta)
		assert.Equal(t, 0, ctx.CurrentVersion.RC)
	})

	t.Run("invalid tagged git repo.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()
		exec.Command("git", "init").Run() //nolint:gosec // This is a test and we need to run git commands.
		exec.Command("git", "tag", "v1.0.0").Run()

		args := []string{}
		ctx := verman.BuildContext(args, false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Equal(t, 0, ctx.CurrentVersion.Major)
		assert.Equal(t, 0, ctx.CurrentVersion.Minor)
		assert.Equal(t, 0, ctx.CurrentVersion.Patch)
		assert.Equal(t, 0, ctx.CurrentVersion.Alpha)
		assert.Equal(t, 0, ctx.CurrentVersion.Beta)
		assert.Equal(t, 0, ctx.CurrentVersion.RC)
	})

	t.Run("valid version file.", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
		assert.NoError(t, err)

		args := []string{}
		ctx := verman.BuildContext(args, false)

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

		args := []string{}
		ctx := verman.BuildContext(args, false)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
		assert.Empty(t, ctx.CurrentVersion)
	})

	t.Run("invalid git tag and valid .version file", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag("invalid")

		err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
		assert.NoError(t, err)

		args := []string{}
		ctx := verman.BuildContext(args, false)

		ctxVersion := ctx.CurrentVersion

		assert.Empty(t, ctxVersion)

		assert.Equal(t, core.SourceNone, ctx.SemverSource)
	})

	t.Run("valid git tag and invalid .version file", func(t *testing.T) {
		cleanUp := testutils.SetupTempDir(t)
		defer cleanUp()

		repoWithTag("v1.0.0")

		err := os.WriteFile(".version", []byte("invalid"), 0644)
		assert.NoError(t, err)

		args := []string{}
		ctx := verman.BuildContext(args, false)

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

		repoWithTag("v1.0.0")

		err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
		assert.NoError(t, err)

		args := []string{}
		ctx := verman.BuildContext(args, false)

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
