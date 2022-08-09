package serato

import "path/filepath"

func GetVolume(path string) string {
	return filepath.VolumeName(path)
}
