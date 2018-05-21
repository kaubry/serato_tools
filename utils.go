package main

import (
	"io"
	"bytes"
	"encoding/binary"
	"path/filepath"
	"os"
	"unicode/utf16"
	"unicode/utf8"
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

// EncodeUTF16 get a utf8 string and translate it into a slice of bytes of ucs2
func EncodeUTF16(s string, add_bom bool) []byte {
	r := []rune(s)
	iresult := utf16.Encode(r)
	var bytes []byte
	if add_bom {
		bytes = make([]byte, 2)
		bytes = []byte{254, 255}
	}
	for _, i := range iresult {
		temp := make([]byte, 2)
		binary.BigEndian.PutUint16(temp, i)
		bytes = append(bytes, temp...)
	}
	return bytes
}

// DecodeUTF16 get a slice of bytes and decode it to UTF-8
func DecodeUTF16(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	bom := UTF16Bom(b)
	if bom < 0 {
		return "", fmt.Errorf("Buffer is too small")
	}

	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	lb := len(b)

	for i := 0; i < lb; i += 2 {
		//assuming bom is big endian if 0 returned
		if bom == 0 || bom == 1 {
			u16s[0] = uint16(b[i+1]) + (uint16(b[i]) << 8)
		}
		if bom == 2 {
			u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		}
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write([]byte(string(b8buf[:n])))
	}

	return ret.String(), nil
}

// UTF16Bom returns 0 for no BOM, 1 for Big Endian and 2 for little endian
// it will return -1 if b is too small for having BOM
func UTF16Bom(b []byte) int8 {
	if len(b) < 2 {
		return -1
	}

	if b[0] == 0xFE && b[1] == 0xFF {
		return 1
	}

	if b[0] == 0xFF && b[1] == 0xFE {
		return 2
	}

	return 0
}
