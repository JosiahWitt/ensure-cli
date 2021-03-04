package fswrite_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure-cli/internal/fswrite"
)

func TestWriteFile(t *testing.T) {
	ensure := ensure.New(t)

	const contents = "testing"
	fileName := filepath.Join(t.TempDir(), "file.txt")

	fsWrite := fswrite.FSWrite{}
	err := fsWrite.WriteFile(fileName, contents, 0655)
	ensure(err).IsNotError()

	actualContents, err := ioutil.ReadFile(fileName)
	ensure(err).IsNotError()
	ensure(string(actualContents)).Equals(contents)
}

func TestMkdirAll(t *testing.T) {
	ensure := ensure.New(t)

	dirName := filepath.Join(t.TempDir(), "nested", "dir")
	fsWrite := fswrite.FSWrite{}
	err := fsWrite.MkdirAll(dirName, 0755)
	ensure(err).IsNotError()

	// Ensure we can write to the directory
	err = ioutil.WriteFile(filepath.Join(dirName, "file.txt"), []byte("testing"), 0600)
	ensure(err).IsNotError()
}

func TestGlobRemoveAll(t *testing.T) {
	t.Run("deletes all matching files", func(t *testing.T) {
		ensure := ensure.New(t)
		dirName := t.TempDir()

		file1Path := filepath.Join(dirName, "it_matches_123.txt")
		err := ioutil.WriteFile(file1Path, []byte("testing123"), 0600)
		ensure(err).IsNotError()

		file2Path := filepath.Join(dirName, "it_matches_456.txt")
		err = ioutil.WriteFile(file2Path, []byte("testing456"), 0600)
		ensure(err).IsNotError()

		file3Path := filepath.Join(dirName, "it_does_not_match_789.txt")
		err = ioutil.WriteFile(file3Path, []byte("testing789"), 0600)
		ensure(err).IsNotError()

		fsWrite := fswrite.FSWrite{}
		err = fsWrite.GlobRemoveAll(filepath.Join(dirName, "it_matches_*.txt"))
		ensure(err).IsNotError()

		_, err = ioutil.ReadFile(file1Path)
		ensure(err != nil).IsTrue() // Expect error reading file, since it was deleted

		_, err = ioutil.ReadFile(file2Path)
		ensure(err != nil).IsTrue() // Expect error reading file, since it was deleted

		actualContents, err := ioutil.ReadFile(file3Path)
		ensure(err).IsNotError()
		ensure(string(actualContents)).Equals("testing789")
	})

	t.Run("when glob is invalid", func(t *testing.T) {
		ensure := ensure.New(t)

		fsWrite := fswrite.FSWrite{}
		err := fsWrite.GlobRemoveAll(`\`)
		ensure(err).IsError(filepath.ErrBadPattern)
	})
}
