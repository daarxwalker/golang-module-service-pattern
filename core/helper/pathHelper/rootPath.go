package pathHelper

import (
	"os"
	"strings"
)

const (
	rootDir = "example"
)

func GetRoot() (string, error) {
	var rootPath []string

	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	split := strings.Split(path, "/")

	for _, item := range split {
		rootPath = append(rootPath, item)
		if item == rootDir {
			break
		}
	}

	return strings.Join(rootPath, "/"), nil
}
