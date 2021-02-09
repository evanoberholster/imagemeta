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
	r       bufReader
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

	v, err := parser(b)
	if err != nil {
		return nil, err
	}
	//b.parsed = v
	return v, nil
}

type parserFunc func(b *box) (Box, error)

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

// FullBox is a type of box that contains Flags
type FullBox struct {
	box
	F Flags
}

// Flags for a FullBox
// 8 bits -> Version
// 24 bits -> Flags
type Flags uint32

// Flags returns underlying Flags after removing version.
// Flags are 24 bits.
func (f Flags) Flags() uint32 {
	// Left Shift
	f = f << 8
	// Right Shift
	return uint32(f >> 8)
}

// Version returns a uint8 version.
func (f Flags) Version() uint8 {
	return uint8(f >> 24)
}

func (f *Flags) Read(buf []byte) {
	*f = Flags(binary.BigEndian.Uint32(buf[:4]))
}

// MetaBox is a 'meta' box
type MetaBox struct {
	size  uint32
	Flags Flags
	//FullBox
	Handler     HandlerBox
	PrimaryItem PrimaryItemBox
	ItemInfo    ItemInfoBox
	//ItemProperties ItemPropertiesBox
	//ItemLocation   ItemLocationBox
	Children []Box
}

// Size returns the size of the MetaBox
func (mb MetaBox) Size() int64 {
	return int64(mb.size)
}

// Type returns TypeMeta
func (mb MetaBox) Type() BoxType {
	return TypeMeta
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

//func readFullBox(outer *box) (fb FullBox, err error) {
//	fb.box = *outer
//	// Parse FullBox header.
//	buf, err := fb.box.r.Peek(4)
//	if err != nil {
//		return FullBox{}, fmt.Errorf("failed to read 4 bytes of FullBox: %v", err)
//	}
//	fb.F.Read(buf)
//	err = fb.box.r.discard(4)
//	return fb, err
//}

func parseMetaBox(outer *box) (Box, error) {
	flags, err := outer.r.readFlags()
	if err != nil {
		return nil, err
	}
	mb := MetaBox{
		size:  uint32(outer.size),
		Flags: flags}

	boxr := outer.newReader(outer.r.remain)
	var inner box
	for outer.r.remain > 0 {
		inner, err = boxr.readBox()
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
		} else {
			mb.setBox(p)
		}
		outer.r.remain -= inner.size

		if Debug {
			fmt.Println(inner, outer.r.remain, inner.r.remain, inner.size)
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
		if inner, err = boxr.readBox(); err != nil {
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
	//FullBox
	Flags  Flags
	ItemID uint16
}

// Size returns the size of the PrimaryItemBox
func (pitm PrimaryItemBox) Size() int64 {
	return 0
}

// Type returns TypePitm
func (pitm PrimaryItemBox) Type() BoxType {
	return TypePitm
}

func parsePrimaryItemBox(outer *box) (Box, error) {
	flags, err := outer.r.readFlags()
	//fb, err := readFullBox(gen)
	if err != nil {
		return nil, err
	}
	pib := PrimaryItemBox{Flags: flags}
	pib.ItemID, err = outer.r.readUint16()
	if !outer.r.ok() {
		return nil, outer.r.err
	}
	return pib, nil
}

// DataInformationBox is a "dinf" box
type DataInformationBox struct {
	//*box
	Children []Box
}

func (dinf DataInformationBox) Size() int64 {
	return 0
}

func (dinf DataInformationBox) Type() BoxType {
	return TypeDinf
}

func parseDataInformationBox(outer *box) (Box, error) {
	dib := DataInformationBox{}
	outer.r.discard(int(outer.r.remain))
	return dib, nil //br.parseAppendBoxes(&dib.Children)
}
