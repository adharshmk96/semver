package verman_test

import (
	"errors"
	"testing"

	"github.com/adharshmk96/semver/pkg/verman"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	fs := afero.NewMemMapFs()

	repo := verman.NewFileRepo(fs)
	fs.Mkdir("testdir", 0755)
	afero.WriteFile(fs, "testfile.txt", []byte("content"), 0644)

	tests := []struct {
		path   string
		exists bool
	}{
		{"testfile.txt", true},
		{"nonexistent.txt", false},
		{"testdir", false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.exists, repo.FileExists(tt.path))
	}
}

func TestReadFileContent(t *testing.T) {
	fs := afero.NewMemMapFs()
	repo := verman.NewFileRepo(fs)
	afero.WriteFile(fs, "testfile.txt", []byte("content"), 0644)

	content, err := repo.ReadFileContent("testfile.txt")
	assert.NoError(t, err)
	assert.Equal(t, "content", content)

	_, err = repo.ReadFileContent("nonexistent.txt")
	assert.True(t, errors.Is(err, afero.ErrFileNotFound))
}

func TestWriteToFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	repo := verman.NewFileRepo(fs)

	err := repo.WriteToFile("testfile.txt", "content")
	assert.NoError(t, err)

	content, err := afero.ReadFile(fs, "testfile.txt")
	assert.NoError(t, err)
	assert.Equal(t, "content", string(content))
}

func TestDeleteFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	repo := verman.NewFileRepo(fs)
	afero.WriteFile(fs, "testfile.txt", []byte("content"), 0644)

	err := repo.DeleteFile("testfile.txt")
	assert.NoError(t, err)

	_, err = fs.Stat("testfile.txt")
	assert.True(t, errors.Is(err, afero.ErrFileNotFound))
}
