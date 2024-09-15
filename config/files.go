package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// findAbsoluteFilePath searches for a file with the given name starting from the current working
// directory and going up to the root directory. It returns the path of the file if found, or an
// os.ErrNotExist if it does not exist.
func findAbsoluteFilePath(fileName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "getting current working directory")
	}

	for dir := cwd; dir != "/"; dir = filepath.Dir(dir) {
		possibleFilePath := filepath.Join(dir, fileName)
		if _, err := os.Stat(possibleFilePath); err == nil {
			return possibleFilePath, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", errors.Wrapf(os.ErrNotExist, "%s file not found", fileName)
}

// mustFindAbsoluteFilePath is like findFilePath, just that it panics if the file is not found.
func mustFindAbsoluteFilePath(fileName string) string {
	filePath, err := findAbsoluteFilePath(fileName)
	if err != nil {
		panic(errors.Wrapf(err, "finding file path %s", fileName))
	}
	return filePath
}
