package main

import (
	"os"
	"fmt"
	"path/filepath"
)

type Crate struct {
	vrsn    []byte //Version
	osrt    []byte //Sorting
	tvcn    []byte //Column Name
	brev    []byte // ???
	columns []Column
	tracks  []Track
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
	for {
		_, err := ReadBytes(f, 1)
		if err != nil {
			break
		} else {
			f.Seek(-1, 1)
			if string(ReadBytesWithOffset(f, 0, 4)) == "ovct" {
				f.Seek(-4, 1)
				columns = append(columns, readColumn(f))
			} else {
				f.Seek(-4, 1)
				break
			}
		}
	}
	return columns
}

func readTracks(f *os.File) []Track {
	var tracks []Track
	for {
		_, err := ReadBytes(f, 1)
		if err != nil {
			break
		} else {
			f.Seek(-1, 1)
			tracks = append(tracks, ReadTrack(f))
		}
	}
	return tracks

}

func (c *Crate) AddTrack(f *os.File) {
	path, _ := filepath.Abs(f.Name())
	t := NewTrack(path)
	if !c.ContainsTrack(t) {
		c.tracks = append(c.tracks, t)
	} else {
		fmt.Printf("Track already in crate !!!")
	}
}

func (c *Crate) TrackList() []string {
	var output []string
	for _, track := range c.tracks {
		output = append(output, string(track.CleanTrackName()))
	}
	return output
}

func (c *Crate) ContainsTrack(t Track) bool {
	for _, track := range c.tracks {
		if track.Equals(t) {
			return true
		}
	}
	return false
}
