package main

import (
	"os"
	"fmt"
	"bytes"
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
	ovct []byte // Total size of column bytes data (without ovct)
	tvcn []byte // Column Name (4 + 4 + dynamic length)
	tvcw []byte // Column Width (in pixel ???) (4 + 4 + dynamic length) // 0 padded string
}

func NewColumn(name ColumnName, width int32) Column {
	c := Column{
		tvcn: PadByteArray([]byte(string(name))),
		tvcw: PadByteArray([]byte(string(width))),
	}
	return c
}

func readColumn(f *os.File) Column {
	return Column{
		ovct: ReadBytesWithOffset(f, 4, 4),
		tvcn: ReadBytesWithDynamicLength(f, 4, 4),
		tvcw: ReadBytesWithDynamicLength(f, 4, 4),
	}
}

func (c *Column) Equals(c2 Column) bool {
	return bytes.Equal(c.tvcn, c2.tvcn)
}

func (c *Column) GetColumnBytes() []byte {
	var output []byte
	output = append(output, []byte("ovct")...)
	length := len(c.tvcn) + len(c.tvcw) + 16
	output = append(output, Int32ToByteArray(4, uint32(length))...)

	output = append(output, []byte("tvcn")...)
	output = append(output, GetBytesWithDynamicLength(c.tvcn, 4)...)

	output = append(output, []byte("tvcw")...)
	output = append(output, GetBytesWithDynamicLength(c.tvcw, 4)...)

	return output
}

func (c Column) String() string {
	return fmt.Sprintf("ovct: %d  //  cleaned tvcn: %s  //  tvcw: %s", ReadInt32(c.ovct), string(UnPadByteArray(c.tvcn)), string(UnPadByteArray(c.tvcw)))
}
