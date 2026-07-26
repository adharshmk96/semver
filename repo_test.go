package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// tempDir switches to an empty directory for the duration of the test.
func tempDir(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	assert.NoError(t, err)

	dir := t.TempDir()
	assert.NoError(t, os.Chdir(dir))

	t.Cleanup(func() { os.Chdir(wd) })
}

// gitRepo initializes a git repo with one commit, tagged when tag is not empty.
func gitRepo(t *testing.T, tag string) {
	t.Helper()

	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "you@example.com"},
		{"config", "user.name", "Your Name"},
		{"commit", "--allow-empty", "-m", "initial commit"},
	} {
		assert.NoError(t, exec.Command("git", args...).Run(), args)
	}

	if tag != "" {
		assert.NoError(t, exec.Command("git", "tag", tag).Run())
	}
}

func TestBuildContext(t *testing.T) {
	t.Run("empty dir", func(t *testing.T) {
		tempDir(t)

		ctx := BuildContext(false)

		assert.Equal(t, SourceNone, ctx.Source)
		assert.False(t, ctx.IsGitRepo)
		assert.Equal(t, &Semver{}, ctx.Version)
	})

	t.Run("untagged git repo", func(t *testing.T) {
		tempDir(t)
		gitRepo(t, "")

		ctx := BuildContext(false)

		assert.Equal(t, SourceNone, ctx.Source)
		assert.True(t, ctx.IsGitRepo)
	})

	t.Run("invalid tag", func(t *testing.T) {
		tempDir(t)
		gitRepo(t, "invalid")

		ctx := BuildContext(false)

		assert.Equal(t, SourceNone, ctx.Source)
	})

	t.Run("tagged git repo", func(t *testing.T) {
		tempDir(t)
		gitRepo(t, "v1.0.0")

		ctx := BuildContext(false)

		assert.Equal(t, SourceGit, ctx.Source)
		assert.Equal(t, "v1.0.0", ctx.Version.String())
	})

	t.Run("git tag wins over version file", func(t *testing.T) {
		tempDir(t)
		gitRepo(t, "v1.0.0")
		assert.NoError(t, os.WriteFile(VersionFile, []byte("v2.0.0"), 0644))

		ctx := BuildContext(false)

		assert.Equal(t, SourceGit, ctx.Source)
		assert.Equal(t, "v1.0.0", ctx.Version.String())
	})

	t.Run("valid version file", func(t *testing.T) {
		tempDir(t)
		assert.NoError(t, os.WriteFile(VersionFile, []byte("v1.0.0-rc.1"), 0644))

		ctx := BuildContext(false)

		assert.Equal(t, SourceFile, ctx.Source)
		assert.Equal(t, "v1.0.0-rc.1", ctx.Version.String())
	})

	t.Run("invalid version file", func(t *testing.T) {
		tempDir(t)
		assert.NoError(t, os.WriteFile(VersionFile, []byte("invalid"), 0644))

		ctx := BuildContext(false)

		assert.Equal(t, SourceNone, ctx.Source)
	})
}

func TestCommit(t *testing.T) {
	t.Run("writes the version file outside a git repo", func(t *testing.T) {
		tempDir(t)

		ctx := BuildContext(false)
		assert.NoError(t, ctx.Commit("v0.0.1"))

		content, err := os.ReadFile(VersionFile)
		assert.NoError(t, err)
		assert.Equal(t, "v0.0.1", string(content))
	})

	t.Run("tags the git repo", func(t *testing.T) {
		tempDir(t)
		gitRepo(t, "")

		ctx := BuildContext(false)
		assert.NoError(t, ctx.Commit("v0.0.1"))

		assert.NoFileExists(t, VersionFile)
		assert.Equal(t, "v0.0.1", BuildContext(false).Version.String())
	})
}

func TestReset(t *testing.T) {
	t.Run("removes the version file", func(t *testing.T) {
		tempDir(t)
		assert.NoError(t, os.WriteFile(VersionFile, []byte("v1.0.0"), 0644))

		assert.NoError(t, BuildContext(false).Reset(false))
		assert.NoFileExists(t, VersionFile)
	})

	t.Run("removes all local tags", func(t *testing.T) {
		tempDir(t)
		gitRepo(t, "v1.0.0")
		assert.NoError(t, exec.Command("git", "tag", "v1.1.0").Run())

		assert.NoError(t, BuildContext(false).Reset(false))
		assert.Equal(t, SourceNone, BuildContext(false).Source)
	})
}
