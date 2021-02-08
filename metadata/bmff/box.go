package bmff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ErrUnknownBox is returned by Box.Parse for unrecognized box types.
var ErrUnknownBox = errors.New("heif: unknown box")

// BoxType is an ISOBMFF box
type BoxType uint8

// Common box types.
const (
	TypeUnknown BoxType = iota
	TypeFtyp            // 'ftyp'
	TypeMeta            // 'meta'
	TypeInfe            // 'infe'
	TypeDinf            // 'dinf'
	TypeHdlr            // 'hdlr'
	TypeIinf            // 'iinf'
	TypePitm            // 'pitm'
	TypeIref            // 'iref'
	TypeIprp            // 'iprp'
	TypeIdat            // 'idat'
	TypeIloc            // 'iloc'
	TypeUUID            // 'uuid'
	TypeImag            // 'Imag'
)

var mapStringBoxType = map[string]BoxType{
	"ftyp": TypeFtyp,
	"meta": TypeMeta,
	"infe": TypeInfe,
	"dinf": TypeDinf,
	"hdlr": TypeHdlr,
	"iinf": TypeIinf,
	"pitm": TypePitm,
	"iref": TypeIref,
	"iprp": TypeIprp,
	"idat": TypeIdat,
	"iloc": TypeIloc,
	"uuid": TypeUUID,
	"Imag": TypeImag,
}

var mapBoxTypeString = map[BoxType]string{
	TypeFtyp: "ftyp",
	TypeMeta: "meta",
	TypeInfe: "infe",
	TypeDinf: "dinf",
	TypeHdlr: "hdlr",
	TypeIinf: "iinf",
	TypePitm: "pitm",
	TypeIref: "iref",
	TypeIprp: "iprp",
	TypeIdat: "idat",
	TypeIloc: "iloc",
	TypeUUID: "uuid",
	TypeImag: "Imag",
}

func (t BoxType) String() string {
	str, ok := mapBoxTypeString[t]
	if ok {
		return str
	}
	return "nnnn"
}

func boxType(buf []byte) BoxType {
	if buf[0] == 'i' {
		if buf[1] == 'n' && buf[2] == 'f' && buf[3] == 'e' {
			return TypeInfe
		}
	}
	if len(buf) == 4 {
		b, ok := mapStringBoxType[string(buf)]
		if ok {
			return b
		}
	}
	if Debug {
		fmt.Println(string(buf))
	}
	return TypeUnknown
}

type box struct {
	r bufReader
	//r       *bufio.Reader
	size    int64 // 0 means unknown, will read to end of file (box container)
	err     error
	parsed  Box // if non-nil, the Parsed result
	boxType BoxType
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
		TypeDinf: parseDataInformationBox,
		//boxType("dref"): parseDataReferenceBox,
		TypeFtyp: parseFileTypeBox,
		TypeHdlr: parseHandlerBox,
		TypeIinf: parseItemInfoBox,
		TypeInfe: parseItemInfoEntry,
		//boxType("iloc"): parseItemLocationBox,
		//boxType("ipco"): parseItemPropertyContainerBox,
		//boxType("ipma"): parseItemPropertyAssociation,
		//boxType("iprp"): parseItemPropertiesBox,
		//boxType("irot"): parseImageRotation,
		//boxType("ispe"): parseImageSpatialExtentsProperty,
		TypeMeta: parseMetaBox,
		TypePitm: parsePrimaryItemBox,
	}
}

type FullBox struct {
	box
	Version uint8
	Flags   uint32 // 24 bits
}

type MetaBox struct {
	FullBox
	Handler     HandlerBox
	PrimaryItem PrimaryItemBox
	ItemInfo    ItemInfoBox
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
	case ItemInfoBox:
		mb.ItemInfo = v
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
	fb.box = *outer
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
			boxr.br.discard(int(inner.r.remain))
			//fmt.Println(err, inner.r.remain)

		} else {
			mb.setBox(p)
		}
		mb.r.remain -= inner.size

		if Debug {
			fmt.Println(inner, mb.r.remain, inner.r.remain, inner.size)
		}
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
