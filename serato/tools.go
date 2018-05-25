package serato

import (
	"github.com/watershine/serato_crates/files"
	"strings"
	"os"
	"path/filepath"
)

type Config struct {
	MusicPath  string
	VolumePath string
	RootCrate  string
}

func CreateCrates(files map[string][]string, c *Config) {
	ensureDirectories(c)
	for key, tracks := range files {
		//info, err := os.Stat(crateFilePath)
		createCrate(getCratePath(key, c), GetDefaultColumn(), c,  tracks...)
		//crateFile, err := os.Create()
		//if err != nil {
		//	log.Printf("Can't create file %s", file)
		//}

	}
}

func ensureDirectories(c *Config) {
	path := getSubcrateFolder(c)
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
		crate.AddTrack(removeVolumeFromPath(t, c))
	}
	files.WriteToFile(path, crate.GetCrateBytes())
}

func getCratePath(file string, c *Config) string {
	newPath := removeVolumeFromPath(file, c)
	newPath = removeMusicPathFromPath(newPath, c)
	newPath = strings.Replace(newPath, string(os.PathSeparator), "%%", -1)
	newPath = strings.Replace(newPath, "-", "_",  -1)
	if len(c.RootCrate) > 0 {
		newPath = c.RootCrate + newPath
	}
	return filepath.Join(getSubcrateFolder(c),  newPath+".crate")
}

func removeMusicPathFromPath(file string, c *Config) string {
	return strings.Replace(file, c.MusicPath, "", 1)
}

func removeVolumeFromPath(file string, c *Config) string {
	return strings.Replace(file, c.VolumePath, "", 1)
}

func getSeratoDir(c *Config) string{
	return filepath.Join(c.VolumePath, "_Serato_")
}

func getSubcrateFolder(c *Config) string {
	return filepath.Join(getSeratoDir(c), "Subcrates")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}