package serato

import (
	"os"
	"github.com/watershine/serato_crates/files"
	"fmt"
	"github.com/watershine/serato_crates/encoding"
	"github.com/watershine/serato_crates/logger"
)

type Database struct {
	vrsn              []byte //Version
	DatabaseMusicFile []DatabaseMusicFile
}

func NewDatabase(f *os.File) *Database {
	database := Database{
		vrsn:              files.ReadBytesWithDynamicLength(f, 4, 4),
		DatabaseMusicFile: readDatabaseMusicFiles(f),
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
			logger.Logger.Debug(mf.String())
		}
	}
	//df = append(df, ReadMusicFile(f))

	return df
}

func (d *Database) String() string {
	return fmt.Sprintf("Vrsn: %s\n", d.getVrsn())
}

func (d *Database) getVrsn() string {
	s, _ := encoding.DecodeUTF16(d.vrsn)
	return s
}
