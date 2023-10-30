package verman

import (
	"github.com/spf13/afero"
)

type FileRepo interface {
	FileExists(path string) bool
	ReadFileContent(path string) (string, error)
	WriteToFile(path string, content string) error
	DeleteFile(path string) error
}

type fileRepo struct {
	fs afero.Fs
}

func NewFileRepo(fs afero.Fs) FileRepo {
	return &fileRepo{fs: fs}
}

func (r *fileRepo) FileExists(path string) bool {
	exists, err := r.fs.Stat(path)
	if err != nil {
		return false
	}
	return !exists.IsDir()
}

func (r *fileRepo) ReadFileContent(path string) (string, error) {
	content, err := afero.ReadFile(r.fs, path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (r *fileRepo) WriteToFile(path string, content string) error {
	return afero.WriteFile(r.fs, path, []byte(content), 0644)
}

func (r *fileRepo) DeleteFile(path string) error {
	return r.fs.Remove(path)
}
