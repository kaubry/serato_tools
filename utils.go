package main

import (
	"io"
	"bytes"
	"encoding/binary"
	"path/filepath"
	"os"
	"fmt"
)

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

func Int32ToByteArray(arrayLength int, number uint32) []byte {
	bs := make([]byte, arrayLength)
	binary.BigEndian.PutUint32(bs, number)
	return bs
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
	header := Int32ToByteArray(headerLength, uint32(len(value)))
	return append(header, value...)
}

func ReadBytesWithOffset(f *os.File, offset int64, length int64) []byte {
	f.Seek(offset, 1)
	returnValue, err := ReadBytes(f, int(length))
	check(err)
	return returnValue
}

func PadByteArray(input []byte) []byte {
	var output []byte
	for _, b := range input {
		output = append(output, byte(0), b)
	}
	return output
}

func PadForLength(input []byte, length int) []byte {
	for {
		if len(input) < length {
			input = append([]byte{0}, input...)
		} else {
			break
		}
	}
	return input
}

func UnPadByteArray(input []byte) []byte {
	var t []byte
	for _, b := range input {
		if b != byte(0) {
			t = append(t, b)
		}
	}
	return t
}

func StringToPaddedByteArray(s string) []byte {
	return PadByteArray([]byte(s))
}

func ListFiles(dir string) map[string][]string {
	output := make(map[string][]string)
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if output[path] == nil {
				output[path] = []string{}
			}
		} else if filepath.Ext(path) == ".mp3" {
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
