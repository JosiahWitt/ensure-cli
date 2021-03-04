package fswrite

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//nolint:golint // Ignore stuttering concerns
type FSWriteIface interface {
	WriteFile(filename string, data string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	GlobRemoveAll(pattern string) error
}

type FSWrite struct{}

var _ FSWriteIface = &FSWrite{}

// WriteFile wraps ioutil.WriteFile.
func (*FSWrite) WriteFile(filename string, data string, perm os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), perm)
}

// MkdirAll wraps os.MkdirAll.
func (*FSWrite) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// GlobRemoveAll deletes everything matching the provided pattern.
func (*FSWrite) GlobRemoveAll(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			return err
		}
	}

	return nil
}
