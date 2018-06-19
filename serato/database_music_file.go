package serato

import (
	"github.com/watershine/serato_crates/files"
	"os"
	"github.com/watershine/serato_crates/encoding"
	"fmt"
)

type DatabaseMusicFile struct {
	otrk []byte //length
	ttyp []byte
	pfil []byte
	tsng []byte
	tart []byte
	talb []byte
	tgen []byte
	tlen []byte
	tbit []byte
	tsmp []byte
	tbpm []byte
	tcom []byte
	tgrp []byte
	tadd []byte
	tiid []byte
	uadd []byte //??
	ulbl []byte //??
	utme []byte //??
	sbav []byte //??
	bhrt []byte //??
	bmis []byte //??
	bply []byte //??
	blop []byte //??
	bitu []byte //??
	bovc []byte //??
	bcrt []byte //??
	biro []byte //??
	bwlb []byte //??
	bwll []byte //??
	buns []byte //??
	bbgl []byte //??
	bkrk []byte //??
}

func ReadMusicFile(f *os.File) DatabaseMusicFile {
	df := DatabaseMusicFile{
		otrk: files.ReadBytesWithOffset(f, 4, 4),
		ttyp: files.ReadBytesWithDynamicLength(f, 4, 4),
		pfil: files.ReadBytesWithDynamicLength(f, 4, 4),
		tsng: files.ReadBytesWithDynamicLength(f, 4, 4),
		tart: files.ReadBytesWithDynamicLength(f, 4, 4),
		talb: files.ReadBytesWithDynamicLength(f, 4, 4),
		tgen: files.ReadBytesWithDynamicLength(f, 4, 4),
		tlen: files.ReadBytesWithDynamicLength(f, 4, 4),
		tbit: files.ReadBytesWithDynamicLength(f, 4, 4),
		tsmp: files.ReadBytesWithDynamicLength(f, 4, 4),
		tbpm: files.ReadBytesWithDynamicLength(f, 4, 4),
		tcom: files.ReadBytesWithDynamicLength(f, 4, 4),
		tgrp: files.ReadBytesWithDynamicLength(f, 4, 4),
		tadd: files.ReadBytesWithDynamicLength(f, 4, 4),
		tiid: files.ReadBytesWithDynamicLength(f, 4, 4),
		uadd: files.ReadBytesWithDynamicLength(f, 4, 4),
		ulbl: files.ReadBytesWithDynamicLength(f, 4, 4),
		utme: files.ReadBytesWithDynamicLength(f, 4, 4),
		sbav: files.ReadBytesWithDynamicLength(f, 4, 4),
		bhrt: files.ReadBytesWithDynamicLength(f, 4, 4),
		bmis: files.ReadBytesWithDynamicLength(f, 4, 4),
		bply: files.ReadBytesWithDynamicLength(f, 4, 4),
		blop: files.ReadBytesWithDynamicLength(f, 4, 4),
		bitu: files.ReadBytesWithDynamicLength(f, 4, 4),
		bovc: files.ReadBytesWithDynamicLength(f, 4, 4),
		bcrt: files.ReadBytesWithDynamicLength(f, 4, 4),
		biro: files.ReadBytesWithDynamicLength(f, 4, 4),
		bwlb: files.ReadBytesWithDynamicLength(f, 4, 4),
		bwll: files.ReadBytesWithDynamicLength(f, 4, 4),
		buns: files.ReadBytesWithDynamicLength(f, 4, 4),
		bbgl: files.ReadBytesWithDynamicLength(f, 4, 4),
		bkrk: files.ReadBytesWithDynamicLength(f, 4, 4),
	}
	return df
}

func (d *DatabaseMusicFile) String() string {
	return fmt.Sprintf("Otrk: %d\nTtyp: %s\nPfil: %s\nTsng: %s\nTart: %s\n"+
		"Talb: %s\nTgen: %s\nTlen: %s\nTbit: %s\nTsmp: %s\nTbpm: %s\n"+
		"Tcom: %s\nTgrp: %s\nTadd: %s\nTiid: %s\nUadd: %d\nUlbl: %d\nUtme: %d\n" +
		"Sbav: %d\nBhrt: %d\nBmis: %d\nBply: %d\nBlop: %d\nBitu: %d\n" +
		"Bovc: %d\nBcrt: %d\nBiro: %d\nBwlb: %d\nBwll: %d\nBuns: %d\n" +
		"Bbgl: %d\nBkrk: %d\n",
		getIntField(d.otrk), getStringField(d.ttyp), getStringField(d.pfil), getStringField(d.tsng), getStringField(d.tart),
		getStringField(d.talb), getStringField(d.tgen), getStringField(d.tlen), getStringField(d.tbit), getStringField(d.tsmp), getStringField(d.tbpm),
		getStringField(d.tcom), getStringField(d.tgrp), getStringField(d.tadd), getStringField(d.tiid), getIntField(d.uadd), getIntField(d.ulbl), getIntField(d.utme),
		getIntField(d.sbav), getIntField(d.bhrt), getIntField(d.bmis), getIntField(d.bply), getIntField(d.blop), getIntField(d.bitu),
		getIntField(d.bovc), getIntField(d.bcrt), getIntField(d.biro), getIntField(d.bwlb), getIntField(d.bwll), getIntField(d.buns),
		getIntField(d.bbgl), getIntField(d.bkrk) )
}

func getIntField(field []byte) int32 {
	return files.ReadInt32(field)
}

func getStringField(field []byte) string {
	s, _ := encoding.DecodeUTF16(field)
	return s
}
