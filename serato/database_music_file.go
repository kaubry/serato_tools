package serato

import (
	"bytes"
	"os"
	"strconv"
	"strings"

	"github.com/gibsn/serato_tools/encoding"
	"github.com/gibsn/serato_tools/files"
)

type DatabaseMusicFile struct {
	otrk   []byte //length
	fields map[string][]byte
	keys   []string //Keys to keep the order in the fields map
}

//ttyp / pfil / tsng / tart / talb / tgen / tlen / tsiz / tbit / tsmp / tbpm / tcom / tgrp / trmx / tlbl / tcmp / ttyr / tadd / tkey / uadd / utkn / ulbl / utme / ufsb / sbav / bhrt / bmis / bply / blop / bitu / bovc / bcrt / biro / bwlb / bwll / buns / bbgl / bkrk /

func ReadMusicFile(f *os.File) (DatabaseMusicFile, error) {
	otrk, err := files.ReadBytesWithOffset(f, 4, 4)
	if err != nil {
		return DatabaseMusicFile{}, err
	}

	df := DatabaseMusicFile{
		otrk:   otrk,
		fields: make(map[string][]byte),
	}

	readLength := 0
	for readLength < getIntField(df.otrk) {
		nextFieldLength, err := readNextField(f, &df)
		if err != nil {
			return DatabaseMusicFile{}, err
		}

		readLength += nextFieldLength
	}

	return df, nil
}

func readNextField(f *os.File, dmf *DatabaseMusicFile) (int, error) {
	k, _ := files.ReadBytes(f, 4)
	key := string(k)

	nextField, err := files.ReadBytesWithDynamicLength(f, 0, 4)
	if err != nil {
		return 0, err
	}

	dmf.fields[key] = nextField
	dmf.keys = append(dmf.keys, key)

	return len(dmf.fields[string(key)]) + 4 + len(key), nil
}

func (d *DatabaseMusicFile) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("otrk: ")
	buffer.WriteString(strconv.Itoa(getIntField(d.otrk)))
	buffer.WriteString("\n")
	for k, v := range d.fields {
		buffer.WriteString(k)
		buffer.WriteString(": ")
		if isString(k) {
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

func (dmf *DatabaseMusicFile) GetBytes() []byte {
	var output, fields []byte

	//Fields
	for _, k := range dmf.keys {
		fields = append(fields, []byte(k)...)
		fields = append(fields, files.GetBytesWithDynamicLength(dmf.fields[k], 4)...)
	}

	output = append(output, []byte("otrk")...)
	otrk := encoding.Int32ToByteArray(4, uint32(len(fields)))
	output = append(output, otrk...)
	output = append(output, fields...)
	return output
}

func isString(s string) bool {
	return strings.HasPrefix(s, "t") || s == "pfil"
}

func (dmf *DatabaseMusicFile) getFilePath() (string, error) {
	return encoding.DecodeUTF16(dmf.fields["pfil"])
}

func (dmf *DatabaseMusicFile) SetFilePath(path string) {
	dmf.fields["pfil"] = encoding.EncodeUTF16(path, true)
}
