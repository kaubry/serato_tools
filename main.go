package main

import (
	"os"
	"fmt"
)

func main() {
	createCrate()
	//readCrate("./Kevin.crate")
	//fmt.Printf("\n")
	readCrate("./Kevin.crate")
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

func readCrate(path string) {
	f, err := os.Open(path)
	check(err)
	crate := NewCrate(f)
	//track := crate.TrackList()[0]

	fmt.Printf("%s\n", crate)
	for _, c := range crate.columns {
		fmt.Printf("%s\n", c)
	}
	fmt.Printf("---- Tracks ----\n")
	for _, t := range crate.tracks {
		fmt.Printf("%s\n", t)

	}
	//fmt.Printf("%s", crate.tracks[0])

	//t := NewTrack("Users/kevin/Downloads/Haila - De Donde Vengo.mp3")
	//fmt.Printf("%v", t.Equals(crate.tracks[0]))
	//display := crate.brev
	//fmt.Printf("%d bytes: %s\n", len(display), string(display))
}

func createCrate() {
	crate := NewEmptyCrate([]ColumnName{song, artist, length})
	f, _ := os.Open("/Users/kevin/Downloads/Haila - De Donde Vengo.mp3")
	crate.AddTrack(f)
	f, _ = os.Open("/Users/kevin/Downloads/Issac Delgado - Toro Mata.mp3")
	crate.AddTrack(f)
	f, _ = os.Open("/Users/kevin/Downloads/Johnny Polanco Y Su Conjunto Amistad - Happy Birthday.mp3")
	crate.AddTrack(f)
	writeToFile("./Kevin.crate", crate.GetCrateBytes())
}

func writeToFile(path string, data []byte) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	f.Write(data)
	f.Sync()
}
