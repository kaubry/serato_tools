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
		vrsn:    readBytesWithOffset(f, 4, 60),
		osrt:    readBytesWithOffset(f, 4, 4),
		tvcn:    readBytesWithDynamicLength(f, 4, 4),
		brev:    readBytesWithOffset(f, 4, 5),
		columns: readColumns(f),
		tracks:  readTracks(f),
	}
	return &crate
}

func readColumns(f *os.File) []Column {
	var columns []Column
	for string(readBytesWithOffset(f, 0, 4)) == "ovct" {
		f.Seek(-4, 1)
		columns = append(columns, readColumn(f))
	}
	f.Seek(-4, 1)
	return columns
}

func readColumn(f *os.File) Column {
	return Column{
		ovct: readBytesWithOffset(f, 4, 4),
		tvcn: readBytesWithDynamicLength(f, 4, 4),
		tvcw: readBytesWithDynamicLength(f, 4, 4),
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
		otrk: readBytesWithOffset(f, 4, 4),
		ptrk: readBytesWithDynamicLength(f, 4, 4),
	}
}

func readBytesWithDynamicLength(f *os.File, offset int64, headerLength int64) []byte {
	f.Seek(offset, 1)
	l, err := ReadBytes(f, int(headerLength))
	check(err)
	length := ReadInt32(l)
	returnValue, err2 := ReadBytes(f, int(length))
	check(err2)
	return returnValue
}

func readBytesWithOffset(f *os.File, offset int64, length int64) []byte {
	f.Seek(offset, 1)
	returnValue, err := ReadBytes(f, int(length))
	check(err)
	return returnValue
}
