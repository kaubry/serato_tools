package serato

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/kaubry/serato_tools/encoding"
	"github.com/kaubry/serato_tools/files"
)

type ColumnName int

const (
	song ColumnName = iota
	added
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
		"song",
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

func NewColumn(name ColumnName, width int) Column {
	c := Column{
		tvcn: encoding.EncodeUTF16(name.String(), false),
		tvcw: encoding.EncodeUTF16(strconv.Itoa(width), false),
	}
	return c
}

func readColumn(f *os.File) (Column, error) {
	ovct, err := files.ReadBytesWithOffset(f, 4, 4)
	if err != nil {
		return Column{}, err
	}

	tvcn, err := files.ReadBytesWithDynamicLength(f, 4, 4)
	if err != nil {
		return Column{}, err
	}

	tvcw, err := files.ReadBytesWithDynamicLength(f, 4, 4)
	if err != nil {
		return Column{}, err
	}

	return Column{
		ovct: ovct,
		tvcn: tvcn,
		tvcw: tvcw,
	}, nil
}

func (c *Column) Equals(c2 Column) bool {
	return bytes.Equal(c.tvcn, c2.tvcn)
}

func (c *Column) GetColumnBytes() []byte {
	var output []byte
	output = append(output, []byte("ovct")...)
	length := len(c.tvcn) + len(c.tvcw) + 16
	output = append(output, encoding.Int32ToByteArray(4, uint32(length))...)

	output = append(output, []byte("tvcn")...)
	output = append(output, files.GetBytesWithDynamicLength(c.tvcn, 4)...)

	output = append(output, []byte("tvcw")...)
	output = append(output, files.GetBytesWithDynamicLength(c.tvcw, 4)...)

	return output
}

func (c Column) String() string {
	return fmt.Sprintf("ovct: %d  //  cleaned tvcn: %s  //  tvcw: %s", files.ReadInt32(c.ovct), c.getTvcn(), c.getTvcw())
}

func GetDefaultColumn() []ColumnName {
	return []ColumnName{song, artist, length, bpm, key, comment, grouping}
}

func (c Column) getTvcn() string {
	s, _ := encoding.DecodeUTF16(c.tvcn)
	return s
}

func (c Column) getTvcw() string {
	s, _ := encoding.DecodeUTF16(c.tvcw)
	return s
}
