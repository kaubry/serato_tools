package main

import (
	"os"
	"fmt"
)

func main() {
	//listFiles()
	//f, _ := os.Open("./o.mp3")
	//addTrack(f)
	//readCrate("./ALFA18 Fri.crate")
	readCrate("./Test.crate")
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

	fmt.Printf("'%s\n", crate)
	for _, c := range crate.columns {
		fmt.Printf("%s\n", c)
	}
	//fmt.Printf("%s", crate.tracks[0])

	//t := NewTrack("Users/kevin/Downloads/Haila - De Donde Vengo.mp3")
	//fmt.Printf("%v", t.Equals(crate.tracks[0]))
	//display := crate.brev
	//fmt.Printf("%d bytes: %s\n", len(display), string(display))
}
