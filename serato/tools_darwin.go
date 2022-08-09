package serato

import (
	"os"
	"path/filepath"
	"strings"
)

func GetVolume(path string) string {
	if strings.HasPrefix(path, darwinVolumesPrefix) {
		pathSplit := strings.Split(path, string(os.PathSeparator))
		return string(os.PathSeparator) + filepath.Join(pathSplit[1], pathSplit[2])
	}

	if len(path) > 0 && string(path[0]) == darwinRootVolume {
		return darwinRootVolume
	}

	return ""
}
