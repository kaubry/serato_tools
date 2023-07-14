package serato

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	seratoDirName = "_Serato_"
)

type FolderTool struct{}

type SeratoDirGetter interface {
	GetSeratoDir(c *Config, hdg HomeDirGetter) (string, error)
}

type HomeDirGetter interface {
	GetHomeDir() string
}

func (f FolderTool) GetHomeDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, "Music")
}

func (f FolderTool) GetSeratoDir(c *Config, hdg HomeDirGetter) (string, error) {
	volume := GetVolume(c.MusicPath)
	if volume == "" {
		return "", fmt.Errorf("%w '%s'", ErrInvalidPath, c.MusicPath)
	}
	switch runtime.GOOS {
	case "windows":
		if volume == "C:" {
			return filepath.Join(hdg.GetHomeDir(), seratoDirName), nil
		}
		return filepath.Join(volume, string(os.PathSeparator)+seratoDirName), nil
	case "darwin":
		if volume == darwinRootVolume {
			return filepath.Join(hdg.GetHomeDir(), seratoDirName), nil
		}
		return filepath.Join(volume, seratoDirName), nil
	default:
		return "", errors.New("OS not supported")
	}

}
