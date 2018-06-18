package serato

import (
	"os"
	"fmt"
	"github.com/watershine/serato_crates/encoding"
	"github.com/watershine/serato_crates/files"
)

const version = "1.0/Serato ScratchLive Crate"

type Crate struct {
	vrsn    []byte //Version
	osrt    []byte //Sorting (always int 19 or 25)
	tvcn    []byte //Default sorting column name (song, artist, etc...). Default #
	brev    []byte // ??? Always 5 bytes 0 0 0 1 0, maybe a delimiter
	columns []Column
	tracks  []Track
}

func NewCrate(f *os.File) *Crate {
	crate := Crate{
		vrsn:    files.ReadBytesWithDynamicLength(f, 4, 4),
		osrt:    files.ReadBytesWithOffset(f, 4, 4),
		tvcn:    files.ReadBytesWithDynamicLength(f, 4, 4),
		brev:    files.ReadBytesWithOffset(f, 4, 5),
		columns: readColumns(f),
		tracks:  readTracks(f),
	}
	return &crate
}

func NewEmptyCrate(columnNames []ColumnName) *Crate {
	crate := Crate{
		vrsn:    encoding.EncodeUTF16(version, false),
		osrt:    encoding.Int32ToByteArray(4, 19),
		tvcn:    encoding.EncodeUTF16("#", false),
		brev:    []byte{0, 0, 0, 1, 0},
		columns: make([]Column, 0),
		tracks:  make([]Track, 0),
	}

	for _, c := range columnNames {
		crate.AddColumn(c)
	}
	return &crate
}

func readColumns(f *os.File) []Column {
	var columns []Column
	for {
		_, err := files.ReadBytes(f, 1)
		if err != nil {
			break
		} else {
			f.Seek(-1, 1)
			if string(files.ReadBytesWithOffset(f, 0, 4)) == "ovct" {
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
		_, err := files.ReadBytes(f, 1)
		if err != nil {
			break
		} else {
			f.Seek(-1, 1)
			tracks = append(tracks, ReadTrack(f))
		}
	}
	return tracks

}

func (c *Crate) AddTrack(path string) {
	//path, _ := filepath.Abs(f.Name())
	t := NewTrack(path)
	if !c.ContainsTrack(t) {
		c.tracks = append(c.tracks, t)
	} else {
		fmt.Printf("Track already in crate !!!")
	}
}

func (c *Crate) RemoveTrack(path string) {
	t := NewTrack(path)
	if i := c.IndexOfTrack(t); i >= 0 {
		c.tracks = append(c.tracks[:i], c.tracks[i+1:]...)
	} else {
		fmt.Printf("Track not in crate !!!")
	}
}

func (c *Crate) AddColumn(name ColumnName) {
	column := NewColumn(name, 0)
	if !c.ContainsColumn(column) {
		c.columns = append(c.columns, column)
	} else {
		fmt.Printf("Column already in crate !!!")
	}
}

func (c *Crate) NumberOfTracks() int {
	return len(c.tracks)
}

func (c *Crate) TrackList() []string {
	var output []string
	for _, track := range c.tracks {
		output = append(output, track.CleanTrackName())
	}
	return output
}

func (c *Crate) ContainsTrack(t Track) bool {
	return c.IndexOfTrack(t) >= 0
}

func (c *Crate) IndexOfTrack(t Track) int {
	for index, track := range c.tracks {
		if track.Equals(t) {
			return index
		}
	}
	return -1
}

func (c *Crate) ContainsColumn(col Column) bool {
	for _, column := range c.columns {
		if column.Equals(col) {
			return true
		}
	}
	return false
}

func (c *Crate) GetCrateBytes() []byte {
	var output []byte
	//Version
	output = append(output, []byte("vrsn")...)
	output = append(output, files.GetBytesWithDynamicLength(c.vrsn, 4)...)
	//Sorting
	output = append(output, []byte("osrt")...)
	output = append(output, c.osrt...)
	//Column Sort
	output = append(output, []byte("tvcn")...)
	output = append(output, files.GetBytesWithDynamicLength(c.tvcn, 4)...)
	//Brev
	output = append(output, []byte("brev")...)
	output = append(output, c.brev...)

	//Columns
	for _, col := range c.columns {
		output = append(output, col.GetColumnBytes()...)
	}

	//Tracks
	for _, track := range c.tracks {
		output = append(output, track.GetTrackBytes()...)
	}

	return output
}

func (c *Crate) String() string {
	return fmt.Sprintf("Vrsn: %s\n Osrt: %d\n Tvcn: %s", c.getVrsn(), files.ReadInt32(c.osrt), c.getTvcn())
}

func (c *Crate) getVrsn() string {
	s, _ := encoding.DecodeUTF16(c.vrsn)
	return s
}

func (c *Crate) getTvcn() string {
	s, _ := encoding.DecodeUTF16(c.tvcn)
	return s
}
