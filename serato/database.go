package serato

import (
	"os"
	"github.com/watershine/serato_crates/files"
	"fmt"
	"github.com/watershine/serato_crates/encoding"
)

type Database struct {
	vrsn []byte //Version
	Dmfs []DatabaseMusicFile
}

func NewDatabase(f *os.File) *Database {
	database := Database{
		vrsn: files.ReadBytesWithDynamicLength(f, 4, 4),
		Dmfs: readDatabaseMusicFiles(f),
	}
	return &database
}
func readDatabaseMusicFiles(f *os.File) []DatabaseMusicFile {
	var df []DatabaseMusicFile
	for {
		_, err := files.ReadBytes(f, 1)
		if err != nil {
			break
		} else {
			f.Seek(-1, 1)
			mf := ReadMusicFile(f)
			df = append(df, mf)
		}
	}
	return df
}

func (d *Database) GetBytes() []byte {
	var output []byte
	output = append(output, []byte("vrsn")...)
	output = append(output, files.GetBytesWithDynamicLength(d.vrsn, 4)...)

	//MusicFiles
	for _, dmf := range d.Dmfs {
		output = append(output, dmf.GetBytes()...)
	}

	return output
}

func (d *Database) String() string {
	return fmt.Sprintf("Vrsn: %s\n", d.getVrsn())
}

func (d *Database) getVrsn() string {
	s, _ := encoding.DecodeUTF16(d.vrsn)
	return s
}
