package files

import (
	"bytes"
	"encoding/binary"
	"gopkg.in/fatih/set.v0"
	"io"
	"os"
	"path/filepath"
	"watershine/serato_tools/encoding"
)

func WriteToFile(path string, data []byte) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	f.Write(data)
	f.Sync()
}

func ReadBytes(r io.Reader, nbrOfBytes int) ([]byte, error) {
	b := make([]byte, nbrOfBytes)
	_, err := io.ReadAtLeast(r, b, nbrOfBytes)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReadInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func ReadBytesWithDynamicLength(f *os.File, offset int64, headerLength int64) []byte {
	f.Seek(offset, 1)
	l, err := ReadBytes(f, int(headerLength))
	check(err)
	length := ReadInt32(l)
	returnValue, err2 := ReadBytes(f, int(length))
	check(err2)
	return returnValue
}

func GetBytesWithDynamicLength(value []byte, headerLength int) []byte {
	header := encoding.Int32ToByteArray(headerLength, uint32(len(value)))
	return append(header, value...)
}

func ReadBytesWithOffset(f *os.File, offset int64, length int64) []byte {
	f.Seek(offset, 1)
	returnValue, err := ReadBytes(f, int(length))
	check(err)
	return returnValue
}

func ListFiles(dir string, supporterExtension *set.Set) map[string][]string {
	output := make(map[string][]string)
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if output[path] == nil {
				output[path] = []string{}
			}
		} else if supporterExtension.Has(filepath.Ext(path)) {
			key := filepath.Dir(path)
			if output[key] == nil {
				output[key] = []string{}
			}
			output[key] = append(output[key], path)
		}
		return nil
	})

	return output
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
