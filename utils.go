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

func ReadBytesWithDynamicLength(f *os.File, offset int64, headerLength int64) []byte {
	f.Seek(offset, 1)
	l, err := ReadBytes(f, int(headerLength))
	check(err)
	length := ReadInt32(l)
	returnValue, err2 := ReadBytes(f, int(length))
	check(err2)
	return returnValue
}

func ReadBytesWithOffset(f *os.File, offset int64, length int64) []byte {
	f.Seek(offset, 1)
	returnValue, err := ReadBytes(f, int(length))
	check(err)
	return returnValue
}

func listFiles() {
	searchDir := "/Users/kevin/Desktop/test-mp3"

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	for _, file := range fileList {
		fmt.Println(file)
	}
}