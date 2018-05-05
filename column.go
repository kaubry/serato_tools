package main

import (
	"os"
	"fmt"
)

type ColumnName int

const (
	added       ColumnName = iota + 1
	album
	artist
	bitrate
	bpm
	comment
	composer
	filename
	genre
	grouping
	key
	label
	length
	location
	remixer
	samplerate
	size
	track
	video_track
	year
)

func (c ColumnName) String() string {
	names := [...]string{
		"added",
		"album",
		"artist",
		"bitrate",
		"bpm",
		"comment",
		"composer",
		"filename",
		"genre",
		"grouping",
		"key",
		"label",
		"length",
		"location",
		"remixer",
		"samplerate",
		"size",
		"track",
		"video track",
		"year",
	}
	return names[c]
}

type Column struct {
	ovct []byte // total size of column bytes data (without ovct)
	tvcn []byte //Column Name (4 + 4 + dynamic length)
	tvcw []byte // Column Width (in pixel ???) (4 + 4 + dynamic length) // 0 padded string
}

func readColumn(f *os.File) Column {
	return Column{
		ovct: ReadBytesWithOffset(f, 4, 4),
		tvcn: ReadBytesWithDynamicLength(f, 4, 4),
		tvcw: ReadBytesWithDynamicLength(f, 4, 4),
	}
}

func (c Column) String() string {
	return fmt.Sprintf("ovct: %d  //  cleaned tvcn: %s  //  tvcw: %s", ReadInt32(c.ovct), string(UnPadByteArray(c.tvcn)), string(UnPadByteArray(c.tvcw)))
}