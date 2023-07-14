package serato

import (
	"errors"
	"fmt"
	"gopkg.in/fatih/set.v0"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kaubry/serato_tools/files"
)

var ErrInvalidPath = errors.New("invalid path")

type Config struct {
	MusicPath string
	RootCrate string
}

func CreateCrates(files map[string][]string, c *Config) {
	ensureDirectories(c)
	for key, tracks := range files {
		cratePath := getCratePath(key, c)
		createCrate(cratePath, GetDefaultColumn(), c, tracks...)
	}
}

func ensureDirectories(c *Config) {
	path, err := GetSubcrateFolder(c, FolderTool{})
	check(err)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e := os.MkdirAll(path, os.ModePerm)
		check(e)
	}
}

func createCrate(path string, columns []ColumnName, c *Config, tracks ...string) {
	crate := NewEmptyCrate(columns)
	for _, t := range tracks {
		trackPath, err := RemoveVolumeFromPath(t)
		check(err)
		trackPath = uniformPathSeparator(trackPath)
		crate.AddTrack(trackPath)
	}
	files.WriteToFile(path, crate.GetCrateBytes())
}

func getCratePath(file string, c *Config) string {
	subcrateFolder, err := GetSubcrateFolder(c, FolderTool{})
	check(err)
	newPath := removeMusicPathFromPath(file, c)
	newPath = strings.Replace(newPath, string(os.PathSeparator), "%%", -1)
	newPath = strings.Replace(newPath, "-", "_", -1)
	if len(c.RootCrate) > 0 {
		newPath = c.RootCrate + newPath
	}
	return filepath.Join(subcrateFolder, newPath+".crate")
}

func removeMusicPathFromPath(file string, c *Config) string {
	return strings.Replace(file, c.MusicPath, "", 1)
}

func RemoveVolumeFromPath(path string) (string, error) {
	volume := GetVolume(path)
	if volume == "" {
		return "", fmt.Errorf("%w '%s'", ErrInvalidPath, path)
	}
	if runtime.GOOS == "windows" {
		return strings.Replace(path, volume+string(os.PathSeparator), "", 1), nil
	}

	if runtime.GOOS == "darwin" {
		if volume == darwinRootVolume {
			return path[1:], nil
		}

		return path[len(volume)+1:], nil
	}

	return "", errors.New("OS not supported")
}

//Serato uses "/" as path separator and not the default OS. Needs to change to work between OS X and Windows
func uniformPathSeparator(path string) string {
	return strings.Replace(path, string(os.PathSeparator), "/", -1)
}

func GetSubcrateFolder(c *Config, sdg SeratoDirGetter) (string, error) {
	s, err := sdg.GetSeratoDir(c, FolderTool{})
	if err != nil {
		return "", err
	}
	return filepath.Join(s, "/Subcrates"), nil
}

func GetSupportedExtension() set.Interface {
	s := set.New(set.ThreadSafe)
	s.Add(".mp3",
		".ogg",
		".alac", //Only on MAC
		".flac",
		".aif",
		".wav",
		".mp4",
		".m4a")
	return s
}

func GetFilePath(path string, seratoDir string) (string, error) {
	volume := GetVolume(seratoDir)
	if runtime.GOOS == "windows" {
		return filepath.Join(volume, path), nil
	}

	if runtime.GOOS == "darwin" {
		if volume == "" {
			return "", fmt.Errorf("%w '%s'", ErrInvalidPath, path)
		}

		return filepath.Join(volume, path), nil
	}

	return "", errors.New("OS not supported")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
