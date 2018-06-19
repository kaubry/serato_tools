package serato

import (
	"github.com/watershine/serato_crates/files"
	"os"
	"github.com/watershine/serato_crates/encoding"
	"bytes"
	"strings"
	"strconv"
)

type DatabaseMusicFile struct {
	otrk []byte //length
	fields map[string][]byte
}

func ReadMusicFile(f *os.File) DatabaseMusicFile {
	df := DatabaseMusicFile{
		otrk: files.ReadBytesWithOffset(f, 4, 4),
		fields: make(map[string][]byte),
	}
	readLength := 0
	for readLength < getIntField(df.otrk) {
		readLength += readNextField(f, &df)
	}
	return df
}

func readNextField(f *os.File, dmf *DatabaseMusicFile) int {
	key, _ := files.ReadBytes(f, 4)
	dmf.fields[string(key)] = files.ReadBytesWithDynamicLength(f, 0, 4)
	return len(dmf.fields[string(key)]) + 4 + len(key)
}

func (d *DatabaseMusicFile) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("otrk: ")
	buffer.WriteString(strconv.Itoa(getIntField(d.otrk)))
	buffer.WriteString("\n")
	for k, v := range d.fields {
		buffer.WriteString(k)
		buffer.WriteString(": ")
		if strings.HasPrefix(k, "t") || k == "pfil" {
			buffer.WriteString(getStringField(v))
		} else {
			buffer.WriteString(strconv.Itoa(getIntField(v)))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func getIntField(field []byte) int {
	return int(files.ReadInt32(field))
}

func getStringField(field []byte) string {
	s, _ := encoding.DecodeUTF16(field)
	return s
}
