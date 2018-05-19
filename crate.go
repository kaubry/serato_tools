package main

import (
	"os"
	"fmt"
	"strings"
)

const version = "81.0/Serato ScratchLive Crate"

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
		vrsn:    ReadBytesWithOffset(f, 4, 60),
		osrt:    ReadBytesWithOffset(f, 4, 4),
		tvcn:    ReadBytesWithDynamicLength(f, 4, 4),
		brev:    ReadBytesWithOffset(f, 4, 5),
		columns: readColumns(f),
		tracks:  readTracks(f),
	}
	return &crate
}

func NewEmptyCrate(columnNames []ColumnName) *Crate {
	crate := Crate{
		vrsn:    PadForLength(PadByteArray([]byte(version)), 60),
		osrt:    Int32ToByteArray(4, 19),
		tvcn:    GetBytesWithDynamicLength(PadByteArray([]byte("#")), 4),
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

func (c *Crate) AddTrack(path string) {
	if strings.Contains(path, "Regresar") {
		fmt.Printf("hello")
	}
	//path, _ := filepath.Abs(f.Name())
	t := NewTrack(path)
	if !c.ContainsTrack(t) {
		c.tracks = append(c.tracks, t)
	} else {
		fmt.Printf("Track already in crate !!!")
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
	output = append(output, c.vrsn...)
	//Sorting
	output = append(output, []byte("osrt")...)
	output = append(output, c.osrt...)
	//Column Sort
	output = append(output, []byte("tvcn")...)
	output = append(output, c.tvcn...)
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

func (c Crate) String() string {
	return fmt.Sprintf("Vrsn: %s\n Osrt: %d\n Tvcn: %s", string(UnPadByteArray(c.vrsn)), ReadInt32(c.osrt), string(UnPadByteArray(c.tvcn)))
}
