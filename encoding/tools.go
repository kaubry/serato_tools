package encoding

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type ByteOrder int8

const (
	InvalidBOM   ByteOrder = -1
	BigEndian    ByteOrder = 1
	LittleEndian ByteOrder = 2
	NoBOM        ByteOrder = 0
)

func Int32ToByteArray(arrayLength int, number uint32) []byte {
	bs := make([]byte, arrayLength)
	binary.BigEndian.PutUint32(bs, number)
	return bs
}

// EncodeUTF16 get a utf8 string and translate it into a slice of bytes of ucs2
func EncodeUTF16(s string, addBom bool) []byte {
	r := []rune(s)
	iresult := utf16.Encode(r)
	var bytes []byte
	if addBom {
		bytes = []byte{0xFE, 0xFF} //BIG-Endian BOM
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
		return "", fmt.Errorf("must have even length byte slice")
	}

	bom := UTF16Bom(b)
	if bom == InvalidBOM {
		return "", fmt.Errorf("invalid BOM, buffer is too small")
	}

	runeSlice := decodeUTF16Runes(b, bom)

	return string(runeSlice), nil
}

func decodeUTF16Runes(b []byte, bom ByteOrder) []rune {
	runeSlice := make([]rune, 0)
	for i := 0; i < len(b); i += 2 {
		var u16 uint16
		if bom == NoBOM || bom == BigEndian {
			u16 = uint16(b[i+1]) + (uint16(b[i]) << 8)
		}
		if bom == LittleEndian {
			u16 = uint16(b[i]) + (uint16(b[i+1]) << 8)
		}
		runeValue := utf16.Decode([]uint16{u16})
		runeSlice = append(runeSlice, runeValue...)
	}
	return runeSlice
}

// UTF16Bom returns 0 for no BOM, 1 for Big Endian and 2 for little endian
// it will return -1 if b is too small for having BOM
func UTF16Bom(b []byte) ByteOrder {
	if len(b) < 2 {
		return InvalidBOM
	}

	if b[0] == 0xFE && b[1] == 0xFF {
		return BigEndian
	}

	if b[0] == 0xFF && b[1] == 0xFE {
		return LittleEndian
	}

	return NoBOM
}
