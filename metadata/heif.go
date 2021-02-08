package metadata

import (
	"bufio"
	"fmt"

	"github.com/evanoberholster/imagemeta/metadata/bmff"
)

type HeifMetadata struct {
	Decoder

	FileType    bmff.FileTypeBox
	Handler     *bmff.HandlerBox
	PrimaryItem *bmff.PrimaryItemBox

	// Reader
	br        *bufio.Reader
	discarded uint32
	pos       uint8
}

func NewHeifMetadata(br *bufio.Reader) *HeifMetadata {
	return &HeifMetadata{br: br}
}

func (hm *HeifMetadata) GetMeta() {
	bmr := bmff.NewReader(hm.br)
	p, err := bmr.ReadAndParseBox(bmff.TypeFtyp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p.(bmff.FileTypeBox))
	hm.setBox(p)

	p, err = bmr.ReadAndParseBox(bmff.TypeMeta)
	if err != nil {
		return
	}
	fmt.Println(p.(bmff.MetaBox))
	hm.setBox(p)
}

func (hm *HeifMetadata) setBox(b bmff.Box) {
	switch box := b.(type) {
	case bmff.FileTypeBox:
		hm.FileType = box
	case *bmff.HandlerBox:
		hm.Handler = box
	case *bmff.PrimaryItemBox:
		hm.PrimaryItem = box
	}
}
