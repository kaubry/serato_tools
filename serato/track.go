package serato

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kaubry/serato_tools/encoding"
	"github.com/kaubry/serato_tools/files"
)

type Track struct {
	otrk []byte //Track length
	ptrk []byte //Track name (location)
}

func ReadTrack(f *os.File) (Track, error) {
	otrk, err := files.ReadBytesWithOffset(f, 4, 4)
	if err != nil {
		return Track{}, err
	}

	ptrk, err := files.ReadBytesWithDynamicLength(f, 4, 4)
	if err != nil {
		return Track{}, err
	}

	return Track{
		otrk: otrk,
		ptrk: ptrk,
	}, nil
}

func NewTrack(path string) Track {
	ptrk := encoding.EncodeUTF16(path, false)
	otrk := encoding.Int32ToByteArray(4, uint32(len(ptrk)+8))
	//h := Int32ToByteArray(4, uint32(len(ptrk)))
	//ptrk = append(h, ptrk...)
	t := Track{
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
