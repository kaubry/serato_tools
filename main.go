package main

import (
	"os"
	"strings"
	"path/filepath"
	log "github.com/watershine/serato_crates/logger"
)

const volumePath = "E:\\"
const musicPath = "GDrive\\DJ\\"

func main() {
	log.Logger.Debug("test")
	//version := EncodeUTF16(version, true)
	//fmt.Printf("%v", version)
	//files := ListFiles("E:/GDrive/DJ")
	//files := ListFiles("E:/Test")
	//for k,v := range files {
	//	fmt.Printf("dir: %s  //   val:%s\n", k, v)
	//}
	//createCrates(files)
	//createCrate()
	//readCrate("./cuban.crate")
	//fmt.Printf("\n")
	//readCrate("./Kevin.crate")
	//listFiles()
	//f, _ := os.Open("./o.mp3")
	//addTrack(f)
	//readCrate("./ALFA18 Fri.crate")
	//readCrate("./Test.crate")
	//readCrate("./SubCrate4.crate")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getSeratoDir() string{
	return volumePath+"/_Serato_/"

}

func createCrates(files map[string][]string) {
	for key, tracks := range files {

		//info, err := os.Stat(crateFilePath)
		createCrate(getCratePath(key), GetDefaultColumn(), tracks...)
		//crateFile, err := os.Create()
		//if err != nil {
		//	log.Printf("Can't create file %s", file)
		//}

	}
}

func getCratePath(file string) string {
	newPath := removeVolumeFromPath(file)
	newPath = removeMusicPathFromPath(newPath)
	newPath = strings.Replace(newPath, string(os.PathSeparator), "%%", -1)
	newPath = strings.Replace(newPath, "-", "_",  -1)
	return filepath.Join(getSeratoDir(), "Subcrates",  "DJ Yuma%%"+newPath+".crate")
}

func removeMusicPathFromPath(file string) string {
	return strings.Replace(file, musicPath, "", 1)
}

func removeVolumeFromPath(file string) string {
	return strings.Replace(file, volumePath, "", 1)
}

func readCrate(path string) {
	f, err := os.Open(path)
	check(err)
	NewCrate(f)
	//track := crate.TrackList()[0]

	//fmt.Printf("%s\n", crate)
	//for _, c := range crate.columns {
	//	fmt.Printf("%s\n", c)
	//}
	//fmt.Printf("---- Tracks ----\n")
	//for _, t := range crate.tracks {
	//	fmt.Printf("%s\n", t)
	//
	//}
	//fmt.Printf("%s", crate.tracks[0])

	//t := NewTrack("Users/kevin/Downloads/Haila - De Donde Vengo.mp3")
	//fmt.Printf("%v", t.Equals(crate.tracks[0]))
	//display := crate.brev
	//fmt.Printf("%d bytes: %s\n", len(display), string(display))
}

//func createCrateWithTracks(path string, columns []ColumnName, files []string) {
//	crate
//}

func createCrate(path string, columns []ColumnName, tracks ...string) {
	crate := NewEmptyCrate(columns)
	for _, t := range tracks {
		//f, err := os.Open(t)
		//if err != nil {
		//	log.Printf("File %s is not a track", t)
		//}
		crate.AddTrack(removeVolumeFromPath(t))
	}
	//fmt.Printf("%s\n", crate)
	//for _, c := range crate.columns {
	//	fmt.Printf("%s\n", c)
	//}
	//fmt.Printf("---- Tracks ----\n")
	//for _, t := range crate.tracks {
	//	fmt.Printf("%s\n", t)
	//
	//}
	writeToFile(path, crate.GetCrateBytes())
}

func writeToFile(path string, data []byte) {
	f, _ := os.Create(path)
	//check(err)
	defer f.Close()
	f.Write(data)
	f.Sync()
}

func getCrates(map[string][]string) map[string][]Crate {
 return nil
}
