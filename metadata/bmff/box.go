package bmff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ErrUnknownBox is returned by Box.Parse for unrecognized box types.
var ErrUnknownBox = errors.New("heif: unknown box")

// Common box types.
var (
	TypeFtyp = BoxType{'f', 't', 'y', 'p'}
	TypeMeta = BoxType{'m', 'e', 't', 'a'}
)

// BoxType is an ISOBMFF box
type BoxType [4]byte

func (t BoxType) String() string {
	return string(t[:])
}

var boxTypeUnknown = BoxType{0, 0, 0, 0}

func boxType(s string) BoxType {
	if len(s) == 4 {
		return BoxType{s[0], s[1], s[2], s[3]}
	}
	return boxTypeUnknown
}

type box struct {
	r bufReader
	//r       *bufio.Reader
	size    int64 // 0 means unknown, will read to end of file (box container)
	err     error
	boxType BoxType
	parsed  Box // if non-nil, the Parsed result
}

func (b box) String() string {
	return fmt.Sprintf("(Box) type: \"%s\", size: %d", b.boxType, b.size)
}

func (b box) Size() int64   { return b.size }
func (b box) Type() BoxType { return b.boxType }

func (b *box) Parse() (Box, error) {
	if b.parsed != nil {
		return b.parsed, nil
	}
	parser, ok := parsers[b.Type()]
	if !ok {
		err := fmt.Errorf("Error: Unknown Box %s", b.Type()) //ErrUnknownBox
		return nil, err
	}

	v, err := parser(b, b.r)
	if err != nil {
		return nil, err
	}
	b.parsed = v
	return v, nil
}

type parserFunc func(b *box, br bufReader) (Box, error)

var parsers map[BoxType]parserFunc

func init() {
	parsers = map[BoxType]parserFunc{
		boxType("dinf"): parseDataInformationBox,
		//boxType("dref"): parseDataReferenceBox,
		boxType("ftyp"): parseFileTypeBox,
		boxType("hdlr"): parseHandlerBox,
		//boxType("iinf"): parseItemInfoBox,
		//boxType("infe"): parseItemInfoEntry,
		//boxType("iloc"): parseItemLocationBox,
		//boxType("ipco"): parseItemPropertyContainerBox,
		//boxType("ipma"): parseItemPropertyAssociation,
		//boxType("iprp"): parseItemPropertiesBox,
		//boxType("irot"): parseImageRotation,
		//boxType("ispe"): parseImageSpatialExtentsProperty,
		boxType("meta"): parseMetaBox,
		boxType("pitm"): parsePrimaryItemBox,
	}
}

type FullBox struct {
	*box
	Version uint8
	Flags   uint32 // 24 bits
}

type MetaBox struct {
	FullBox
	Handler     HandlerBox
	PrimaryItem PrimaryItemBox
	//ItemInfo       ItemInfoBox
	//ItemProperties ItemPropertiesBox
	//ItemLocation   ItemLocationBox
	Children []Box
}

func (mb MetaBox) String() string {
	str := fmt.Sprintf("(Box) bmff.Metabox, %d Children", len(mb.Children))
	return str
}

func (mb *MetaBox) setBox(b Box) error {
	switch v := b.(type) {
	case HandlerBox:
		mb.Handler = v
	case PrimaryItemBox:
		mb.PrimaryItem = v
	//case *bmff.ItemInfoBox:
	//	meta.ItemInfo = v
	//case *bmff.ItemPropertiesBox:
	//	meta.Properties = v
	//case *bmff.ItemLocationBox:
	//	meta.ItemLocation = v
	default:
		mb.Children = append(mb.Children, b)
	}
	return nil
}

func readFullBox(outer *box) (fb FullBox, err error) {
	fb.box = outer
	// Parse FullBox header.
	buf, err := fb.box.r.Peek(4)
	if err != nil {
		return FullBox{}, fmt.Errorf("failed to read 4 bytes of FullBox: %v", err)
	}
	fb.Version = buf[0]
	buf[0] = 0
	fb.Flags = binary.BigEndian.Uint32(buf[:4])
	err = fb.box.r.discard(4)
	return fb, err
}

func parseMetaBox(outer *box, br bufReader) (Box, error) {
	fb, err := readFullBox(outer)
	if err != nil {
		return nil, err
	}
	mb := MetaBox{FullBox: fb}
	boxr := mb.newReader(mb.r.remain)
	var inner box
	for mb.r.remain > 0 {
		inner, err = boxr.ReadBox()
		if err != nil {
			if err == io.EOF {
				return mb, nil
			}
			boxr.br.err = err
			return mb, err
		}
		p, err := inner.Parse()
		if err != nil {
			fmt.Println(err)
			break
		}
		mb.setBox(p)
		mb.r.remain -= inner.size
		//boxr.br.discard(int(inner.r.remain))

		fmt.Println(inner, mb.r.remain, inner.r.remain, inner.size)
	}
	//mb.Children, err = fb.parseAppendBoxes()
	//fmt.Println(mb, mb.r.remain)
	return mb, err
}

func (fb *FullBox) parseAppendBoxes() (dst []Box, err error) {
	if fb.err != nil {
		err = fb.err
		return
	}

	boxr := fb.newReader(fb.r.remain)
	var inner box
	i := 5
	for i > 0 {
		if inner, err = boxr.ReadBox(); err != nil {
			if err == io.EOF {
				return dst, nil
			}
			boxr.br.err = err
			return dst, err
		}
		i--
		//fmt.Println("Inner", inner, inner.r.remain)
		boxr.br.discard(int(inner.r.remain))
		//slurp, err := ioutil.ReadAll(inner.Body())
		//if err != nil {
		//	br.err = err
		//	return err
		//}
		//inner.(*box).slurp = slurp
		//dst = append(dst, inner)
	}
	return
}

// PrimaryItemBox is a "pitm" box
type PrimaryItemBox struct {
	FullBox
	ItemID uint16
}

func parsePrimaryItemBox(gen *box, br bufReader) (Box, error) {
	fb, err := readFullBox(gen)
	if err != nil {
		return nil, err
	}
	pib := PrimaryItemBox{FullBox: fb}
	pib.ItemID, _ = gen.r.readUint16()
	if !br.ok() {
		return nil, br.err
	}
	return pib, nil
}

// DataInformationBox is a "dinf" box
type DataInformationBox struct {
	*box
	Children []Box
}

func parseDataInformationBox(gen *box, br bufReader) (Box, error) {
	dib := &DataInformationBox{box: gen}
	gen.r.discard(int(gen.r.remain))
	return dib, nil //br.parseAppendBoxes(&dib.Children)
}
