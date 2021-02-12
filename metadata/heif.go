package metadata

import (
	"bufio"
	"fmt"

	"github.com/evanoberholster/imagemeta/metadata/bmff"
)

type HeifMetadata struct {
	Decoder

	FileType bmff.FileTypeBox
	Meta     bmff.MetaBox

	// Reader
	br *bufio.Reader
}

func NewHeifMetadata(br *bufio.Reader) *HeifMetadata {
	return &HeifMetadata{br: br}
}

func (hm *HeifMetadata) GetMeta() {
	bmr := bmff.NewReader(hm.br)

	ftyp, err := bmr.ReadFtypBox()
	if err != nil {
		fmt.Println(err)
		return
	}
	hm.FileType = ftyp
	fmt.Println(ftyp)

	mb, err := bmr.ReadMetaBox()
	if err != nil {
		return
	}
	hm.Meta = mb
	fmt.Println(mb)
}
