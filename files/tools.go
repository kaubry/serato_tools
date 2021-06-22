package files

import (
	"bytes"
	"encoding/binary"
	"github.com/gibsn/serato_tools/encoding"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/fatih/set.v0"
)

func WriteToFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
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

func ReadBytesWithDynamicLength(f *os.File, offset int64, headerLength int64) ([]byte, error) {
	_, err := f.Seek(offset, 1)
	if err != nil {
		return nil, err
	}

	l, err := ReadBytes(f, int(headerLength))
	if err != nil {
		return nil, err
	}

	length := ReadInt32(l)
	returnValue, err := ReadBytes(f, int(length))
	if err != nil {
		return nil, err
	}

	return returnValue, nil
}

func GetBytesWithDynamicLength(value []byte, headerLength int) []byte {
	header := encoding.Int32ToByteArray(headerLength, uint32(len(value)))
	return append(header, value...)
}

func ReadBytesWithOffset(f *os.File, offset int64, length int64) ([]byte, error) {
	_, err := f.Seek(offset, 1)
	if err != nil {
		return nil, err
	}

	returnValue, err := ReadBytes(f, int(length))
	if err != nil {
		return nil, err
	}

	return returnValue, nil
}

func ListFiles(dir string, supporterExtension set.Interface) map[string][]string {
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
