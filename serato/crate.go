package serato

import (
	"fmt"
	"io"
	"os"

	"github.com/kaubry/serato_tools/encoding"
	"github.com/kaubry/serato_tools/files"
	"github.com/kaubry/serato_tools/logger"
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

func NewCrate(f *os.File) (*Crate, error) {
	vrsn, err := files.ReadBytesWithDynamicLength(f, 4, 4)
	if err != nil {
		return nil, err
	}

	osrt, err := files.ReadBytesWithOffset(f, 4, 4)
	if err != nil {
		return nil, err
	}

	tvcn, err := files.ReadBytesWithDynamicLength(f, 4, 4)
	if err != nil {
		return nil, err
	}

	brev, err := files.ReadBytesWithOffset(f, 4, 5)
	if err != nil {
		return nil, err
	}

	columns, err := readColumns(f)
	if err != nil {
		return nil, err
	}

	tracks, err := readTracks(f)
	if err != nil {
		return nil, err
	}

	crate := Crate{
		vrsn:    vrsn,
		osrt:    osrt,
		tvcn:    tvcn,
		brev:    brev,
		columns: columns,
		tracks:  tracks,
	}

	return &crate, nil
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

func readColumns(f *os.File) ([]Column, error) {
	var columns []Column
	for {
		_, err := files.ReadBytes(f, 1)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		_, err = f.Seek(-1, 1)
		if err != nil {
			return nil, err
		}

		possibleOvct, err := files.ReadBytesWithOffset(f, 0, 4)
		if err != nil {
			return nil, err
		}

		if string(possibleOvct) != "ovct" {
			_, err = f.Seek(-4, 1)
			if err != nil {
				return nil, err
			}

			break
		}

		_, err = f.Seek(-4, 1)
		if err != nil {
			return nil, err
		}

		column, err := readColumn(f)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil
}

func readTracks(f *os.File) ([]Track, error) {
	var tracks []Track
	for {
		_, err := files.ReadBytes(f, 1)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		_, err = f.Seek(-1, 1)
		if err != nil {
			return nil, err
		}

		track, err := ReadTrack(f)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (c *Crate) AddTrack(path string) {
	//path, _ := filepath.Abs(f.Name())
	t := NewTrack(path)
	if !c.ContainsTrack(t) {
		c.tracks = append(c.tracks, t)
	} else {
		logger.Logger.Error("Track already in crate !!!")
	}
}

func (c *Crate) RemoveTrack(path string) {
	t := NewTrack(path)
	if i := c.IndexOfTrack(t); i >= 0 {
		c.tracks = append(c.tracks[:i], c.tracks[i+1:]...)
	} else {
		logger.Logger.Error("Track not in crate !!!")
	}
}

func (c *Crate) AddColumn(name ColumnName) {
	column := NewColumn(name, 0)
	if !c.ContainsColumn(column) {
		c.columns = append(c.columns, column)
	} else {
		logger.Logger.Error("Column already in crate !!!")
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
