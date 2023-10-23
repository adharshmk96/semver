package verman_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.M) {
	tempDir, _ := os.MkdirTemp("", "test")
	os.Chdir(tempDir)

	err := os.WriteFile(".version", []byte("v1.0.0-rc.1"), 0644)
	if err != nil {
		panic(err)
	}

	code := t.Run()
	os.Exit(code)
}

func TestBuildContext(t *testing.T) {
	t.Run("Build context with current version.", func(t *testing.T) {
		args := []string{}
		ctx := verman.BuildContext(args, false)

		ctxVersion := ctx.CurrentVersion

		assert.Equal(t, 1, ctxVersion.Major)
		assert.Equal(t, 0, ctxVersion.Minor)
		assert.Equal(t, 0, ctxVersion.Patch)
		assert.Equal(t, 0, ctxVersion.Alpha)
		assert.Equal(t, 0, ctxVersion.Beta)
		assert.Equal(t, 1, ctxVersion.RC)

		assert.False(t, ctx.IsGitRepo)
		assert.True(t, ctx.IsVerman)

	})

	t.Run("build context in a git repo.", func(t *testing.T) {
		exec.Command("git", "init").Run() //nolint:gosec // This is a test and we need to run git commands.

		args := []string{}
		ctx := verman.BuildContext(args, false)

		assert.True(t, ctx.IsGitRepo)
		assert.True(t, ctx.IsVerman)
	})
}
