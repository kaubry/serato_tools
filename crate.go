package main

import (
	"os"
	"fmt"
)

type Crate struct {
	vrsn    []byte //Version
	osrt    []byte //Sorting
	tvcn    []byte //Column Name
	brev    []byte // ???
	columns []Column
	tracks  []Track
}

type Column struct {
	ovct []byte // ???
	tvcn []byte //Column Name
	tvcw []byte
}

type Track struct {
	otrk []byte //Track length
	ptrk []byte //Track name (location)
}

func NewCrate(f *os.File) *Crate {
	crate := Crate{
		vrsn:    ReadBytesWithOffset(f, 4, 60),
		osrt:    ReadBytesWithOffset(f, 4, 4),
		tvcn:    ReadBytesWithDynamicLength(f, 4, 4),
		brev:    ReadBytesWithOffset(f, 4, 5),
		columns: readColumns(f),
		tracks:  readTracks(f),
	}
	return &crate
}

func readColumns(f *os.File) []Column {
	var columns []Column
	for string(ReadBytesWithOffset(f, 0, 4)) == "ovct" {
		f.Seek(-4, 1)
		columns = append(columns, readColumn(f))
	}
	f.Seek(-4, 1)
	return columns
}

func readColumn(f *os.File) Column {
	return Column{
		ovct: ReadBytesWithOffset(f, 4, 4),
		tvcn: ReadBytesWithDynamicLength(f, 4, 4),
		tvcw: ReadBytesWithDynamicLength(f, 4, 4),
	}
}

func readTracks(f *os.File) []Track {
	var tracks []Track
	for {
		_, err := ReadBytes(f, 1)
		if err != nil {
			break
		} else {
			f.Seek(-1, 1)
			tracks = append(tracks, readTrack(f))
		}
	}
	return tracks

}

func readTrack(f *os.File) Track {
	return Track{
		otrk: ReadBytesWithOffset(f, 4, 4),
		ptrk: ReadBytesWithDynamicLength(f, 4, 4),
	}
}

func (c *Crate) ListTracks() {
	for _, track := range c.tracks {
		fmt.Printf("length: %d / %d, %s\n", ReadInt32(track.otrk), len(track.ptrk)+8, string(cleanTrackName(track.ptrk)))
	}
}

func cleanTrackName(track []byte) []byte {
	var t []byte
	for _, b := range track {
		if b != byte(0) {
			t = append(t, b)
		}
	}
	return t
}
