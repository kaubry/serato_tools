package serato

import (
	"os"
	"bytes"
	"fmt"
	"github.com/watershine/serato_tools/files"
	"github.com/watershine/serato_tools/encoding"
)

type Track struct {
	otrk []byte //Track length
	ptrk []byte //Track name (location)
}

func ReadTrack(f *os.File) Track {
	return Track{
		otrk: files.ReadBytesWithOffset(f, 4, 4),
		ptrk: files.ReadBytesWithDynamicLength(f, 4, 4),
	}
}

func NewTrack(path string) Track {
	ptrk := encoding.EncodeUTF16(path, false)
	otrk := encoding.Int32ToByteArray(4, uint32(len(ptrk) + 8))
	//h := Int32ToByteArray(4, uint32(len(ptrk)))
	//ptrk = append(h, ptrk...)
	t := Track {
		otrk: otrk,
		ptrk: ptrk,
	}
	return t
}

func (t *Track) Equals(t2 Track) bool {
	return bytes.Equal(t.otrk, t2.otrk) && bytes.Equal(t.ptrk, t2.ptrk)
}

func (track *Track) CleanTrackName() string {
	name, _ := encoding.DecodeUTF16(track.ptrk)
	return name
}

func (track *Track) GetTrackBytes() []byte {
	var output []byte
	output = append(output, []byte("otrk")...)
	//otrk := Int32ToByteArray(4, uint32(len(track.ptrk) + 8))
	output = append(output, track.otrk...)

	output = append(output, []byte("ptrk")...)
	output = append(output, files.GetBytesWithDynamicLength(track.ptrk, 4)...)

	return output
}

func (track Track) String() string {
	return fmt.Sprintf("otrk: %d  //  cleaned ptrk: %s", files.ReadInt32(track.otrk), string(track.CleanTrackName()))
}