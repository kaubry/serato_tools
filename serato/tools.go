package serato

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kaubry/serato_tools/files"

	"gopkg.in/fatih/set.v0"
)

const (
	seratoDirName          = "_Serato_"
	darwinVolumesPrefix    = "/Volumes/"
	darwinRootVolumePrefix = "/Users/"
	darwinRootVolume       = "/"
)

var ErrInvalidPath = errors.New("invalid path")

type homeDirGetter interface {
	getHomeDir() string
}

type localHomeDirGetter struct {
}

func (_ localHomeDirGetter) getHomeDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, "Music")
}

var (
	defaultHomeDirGetter homeDirGetter = localHomeDirGetter{}
)

type Config struct {
	MusicPath string
	RootCrate string
}

func CreateCrates(files map[string][]string, c *Config) {
	ensureDirectories(c)
	for key, tracks := range files {
		cratePath := getCratePath(key, c)
		createCrate(cratePath, GetDefaultColumn(), c, tracks...)
		//info, err := os.Stat(crateFilePath)
		//crateFile, err := os.Create()
		//if err != nil {
		//	log.Printf("Can't create file %s", file)
		//}

	}
}

func ensureDirectories(c *Config) {
	path, err := GetSubcrateFolder(c)
	check(err)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e := os.MkdirAll(path, os.ModePerm)
		check(e)
	}
}

func createCrate(path string, columns []ColumnName, c *Config, tracks ...string) {
	crate := NewEmptyCrate(columns)
	for _, t := range tracks {
		//f, err := os.Open(t)
		//if err != nil {
		//	log.Printf("File %s is not a track", t)
		//}
		trackPath, err := RemoveVolumeFromPath(t)
		check(err)
		trackPath = uniformPathSeparator(trackPath)
		crate.AddTrack(trackPath)
	}
	files.WriteToFile(path, crate.GetCrateBytes())
}

func getCratePath(file string, c *Config) string {
	subcrateFolder, err := GetSubcrateFolder(c)
	check(err)
	newPath := removeMusicPathFromPath(file, c)
	newPath = strings.Replace(newPath, string(os.PathSeparator), "%%", -1)
	newPath = strings.Replace(newPath, "-", "_", -1)
	if len(c.RootCrate) > 0 {
		newPath = c.RootCrate + newPath
	}
	return filepath.Join(subcrateFolder, newPath+".crate")
	//return ""
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

func GetSeratoDir(c *Config) (string, error) {
	volume := GetVolume(c.MusicPath)
	if volume == "" {
		return "", fmt.Errorf("%w '%s'", ErrInvalidPath, c.MusicPath)
	}
	if runtime.GOOS == "windows" {
		if volume == "C:" {
			return filepath.Join(getHomeDir(), seratoDirName), nil
		}
		return filepath.Join(volume, string(os.PathSeparator)+seratoDirName), nil
	}

	if runtime.GOOS == "darwin" {
		if volume == darwinRootVolume {
			return filepath.Join(getHomeDir(), seratoDirName), nil
		}

		return filepath.Join(volume, seratoDirName), nil
	}

	return "", errors.New("OS not supported")
}

func getHomeDir() string {
	return defaultHomeDirGetter.getHomeDir()
}

func GetSubcrateFolder(c *Config) (string, error) {
	s, err := GetSeratoDir(c)
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
