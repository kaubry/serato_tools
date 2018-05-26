package serato

import (
	"strings"
	"os"
	"path/filepath"
	"gopkg.in/fatih/set.v0"
	"runtime"
	"os/user"
	"errors"
	"regexp"
	"github.com/watershine/serato_crates/files"
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
		trackPath, err := removeVolumeFromPath(t)
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

func removeVolumeFromPath(path string) (string, error) {
	if runtime.GOOS == "windows" {
		volume := filepath.VolumeName(path)
		return strings.Replace(path, volume+string(os.PathSeparator), "", 1), nil
	} else if runtime.GOOS == "darwin" {
		r, _ := regexp.Compile(`(\/Volume\/(.+)\/).+`)
		if !r.MatchString(path) {
			return strings.Replace(path, string(os.PathSeparator), "", 1), nil
		} else {
			volume := r.FindStringSubmatch(path)[1]
			return strings.Replace(path, volume+string(os.PathSeparator), "", 1), nil
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

		r, _ := regexp.Compile(`(\/Volume\/(.+)\/).+`)
		if !r.MatchString(c.MusicPath) {
			return filepath.Join(getHomeDir(), "_Serato_"), nil
		} else {
			volume := r.FindStringSubmatch(c.MusicPath)[1]
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

func GetSupportedExtension() *set.Set {
	return set.New(".mp3",
		".ogg",
		".alac", //Only on MAC
		".flac",
		".aif",
		".wav",
		".mp4",
		".m4a")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
