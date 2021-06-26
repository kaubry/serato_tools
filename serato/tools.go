package serato

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gibsn/serato_tools/files"

	"gopkg.in/fatih/set.v0"
)

const DARWIN_VOLUME_REGEX = `(\/Volumes\/[\d\w\s]+).*`

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
	if runtime.GOOS == "windows" {
		volume := filepath.VolumeName(path)
		return strings.Replace(path, volume+string(os.PathSeparator), "", 1), nil
	} else if runtime.GOOS == "darwin" {
		r, _ := regexp.Compile(DARWIN_VOLUME_REGEX)
		if !r.MatchString(path) {
			return strings.Replace(path, string(os.PathSeparator), "", 1), nil
		} else {
			matches := r.FindStringSubmatch(path)
			return strings.Replace(path, matches[1]+string(os.PathSeparator), "", 1), nil
		}
	}
	return "", errors.New("OS not supported")
}

func GetSeratoDir(c *Config) (string, error) {
	if runtime.GOOS == "windows" {
		volume := filepath.VolumeName(c.MusicPath)
		if volume == "C:" {
			return filepath.Join(getHomeDir(), "_Serato_"), nil
		} else {
			return filepath.Join(volume, "/_Serato_"), nil
		}
	} else if runtime.GOOS == "darwin" {

		r, _ := regexp.Compile(DARWIN_VOLUME_REGEX)
		if !r.MatchString(c.MusicPath) {
			return filepath.Join(getHomeDir(), "_Serato_"), nil
		} else {
			matches := r.FindStringSubmatch(c.MusicPath)
			volume := matches[1]
			return filepath.Join(volume, "_Serato_"), nil
		}
	}
	return "", errors.New("OS not supported")
}

func getHomeDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, "Music")
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
	if runtime.GOOS == "windows" {
		volume := filepath.VolumeName(seratoDir)
		return filepath.Join(volume, path), nil
	} else if runtime.GOOS == "darwin" {
		r, _ := regexp.Compile(DARWIN_VOLUME_REGEX)
		if r.MatchString(seratoDir) {
			matches := r.FindStringSubmatch(seratoDir)
			volume := matches[1]
			return filepath.Join(volume, path), nil
		} else {
			return path, nil
		}
	}
	return "", errors.New("OS not supported")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
