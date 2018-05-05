package main

import (
	"os"
)

func main() {
	//listFiles()
	readCrate("./ALFA18 Fri.crate")
	//readCrate("./DJ Yuma.crate")
	//readCrate("./Drop.crate")
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
	crate.ListTracks()
	//display := crate.tracks[len(crate.tracks)-1].ptrk
	//display := crate.brev
	//fmt.Printf("%d bytes: %s\n", len(display), string(display))
}
