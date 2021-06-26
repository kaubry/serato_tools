package serato

import (
	"fmt"
	"io"
	"os"

	"github.com/gibsn/serato_tools/encoding"
	"github.com/gibsn/serato_tools/files"
	"github.com/gibsn/serato_tools/logger"
)

type Database struct {
	vrsn []byte //Version
	Dmfs []DatabaseMusicFile
}

func NewDatabase(f *os.File) (*Database, error) {
	vrsn, err := files.ReadBytesWithDynamicLength(f, 4, 4)
	if err != nil {
		return nil, err
	}

	dmfs, err := readDatabaseMusicFiles(f)
	if err != nil {
		return nil, err
	}

	database := Database{
		vrsn: vrsn,
		Dmfs: dmfs,
	}

	return &database, nil
}

func readDatabaseMusicFiles(f *os.File) ([]DatabaseMusicFile, error) {
	var df []DatabaseMusicFile

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

		mf, err := ReadMusicFile(f)
		if err != nil {
			return nil, err
		}

		df = append(df, mf)
	}

	return df, nil
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

func (d *Database) GetMusicFiles() []string {
	var output []string
	for _, dmf := range d.Dmfs {
		s, err := dmf.getFilePath()
		if err != nil {
			logger.Logger.Error(err.Error())
		} else {
			output = append(output, s)
		}
	}
	return output
}

func (d *Database) RemoveMusicFile(path string) {
	if i := d.IndexOfMusicFile(path); i >= 0 {
		d.Dmfs = append(d.Dmfs[:i], d.Dmfs[i+1:]...)
	} else {
		logger.Logger.Error("Music File not in Database !!!")
	}
}

func (d *Database) AddMusicFile(file DatabaseMusicFile) {
	d.Dmfs = append(d.Dmfs, file)
}

func (d *Database) IndexOfMusicFile(path string) int {
	for index, dmf := range d.Dmfs {
		p, _ := dmf.getFilePath()
		if p == path {
			return index
		}
	}
	return -1
}
