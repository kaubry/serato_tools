package main

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"encoding/binary"
)

func main() {
	readCrate("./ALFA18 Fri.crate")
	//readCrate("./Drop.crate")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readCrate(path string) {
	f, err := os.Open(path)
	check(err)
	crate := NewCrate(f)
	display := crate.tracks[len(crate.tracks)-1].ptrk
	//display := crate.brev
	fmt.Printf("%d bytes: %s\n", len(display), string(display))
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
