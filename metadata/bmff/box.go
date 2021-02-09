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
	TypeAv01            // 'av01'
	TypeAv1C            // 'av1C'
	TypeAuxC            // 'auxC'
	TypeFtyp            // 'ftyp'
	TypeMeta            // 'meta'
	TypeInfe            // 'infe'
	TypeDinf            // 'dinf'
	TypeHdlr            // 'hdlr'
	TypeIinf            // 'iinf'
	TypePitm            // 'pitm'
	TypeIref            // 'iref'
	TypeIpco            // 'ipco
	TypeIprp            // 'iprp'
	TypeIdat            // 'idat'
	TypeIloc            // 'iloc'
	TypeUUID            // 'uuid'
	TypeColr            // 'colr'
	TypeHvcC            // 'hvcC'
	TypeIspe            // 'ispe'
	TypeIrot            // 'irot'
	TypePixi            // 'pixi'
	TypeIpma            // 'ipma'
	TypePasp            // 'pasp'
)

var mapStringBoxType = map[string]BoxType{
	"av01": TypeAv01,
	"av1C": TypeAv1C,
	"auxC": TypeAuxC,
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
	"ipco": TypeIpco,
	"colr": TypeColr,
	"hvcC": TypeHvcC,
	"ispe": TypeIspe,
	"irot": TypeIrot,
	"pixi": TypePixi,
	"ipma": TypeIpma,
	"pasp": TypePasp,
}

var mapBoxTypeString = map[BoxType]string{
	TypeAv01: "av01",
	TypeAv1C: "av1C",
	TypeAuxC: "auxC",
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
	TypeIpco: "ipco",
	TypeColr: "colr",
	TypeHvcC: "hvcC",
	TypeIspe: "ispe",
	TypeIrot: "irot",
	TypePixi: "pixi",
	TypeIpma: "ipma",
	TypePasp: "pasp",
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
	r    bufReader
	size int64 // 0 means unknown, will read to end of file (box container)
	err  error
	//parsed  Box // if non-nil, the Parsed result
	boxType BoxType
}

func (b box) String() string {
	return fmt.Sprintf("(Box) type: \"%s\", size: %d", b.boxType, b.size)
}

func (b box) Size() int64   { return b.size }
func (b box) Type() BoxType { return b.boxType }

func (b *box) Parse() (Box, error) {
	if !b.r.anyRemain() {
		return nil, ErrBufReaderLength
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
		TypeIloc: parseItemLocationBox,
		TypeIpco: parseItemPropertyContainerBox,
		TypeIpma: parseItemPropertyAssociation,
		TypeIprp: parseItemPropertiesBox,
		TypeIrot: parseImageRotation,
		TypeIspe: parseImageSpatialExtentsProperty,
		TypeMeta: parseMetaBox,
		TypePitm: parsePrimaryItemBox,
	}
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

	Handler    HandlerBox
	Primary    PrimaryItemBox
	ItemInfo   ItemInfoBox
	Properties ItemPropertiesBox
	Location   ItemLocationBox
	Children   []Box
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
	str := fmt.Sprintf("(Box) bmff.Metabox, %d Children\n", len(mb.Children))
	str += "\t" + mb.Primary.String() + "\n"
	str += "\t" + mb.ItemInfo.String() + "\n"
	str += "\t" + mb.Properties.String() + "\n"
	str += "\t" + mb.Location.String() + "\n"
	return str
}

func (mb *MetaBox) setBox(b Box) error {
	switch v := b.(type) {
	case HandlerBox:
		mb.Handler = v
	case PrimaryItemBox:
		mb.Primary = v
	case ItemInfoBox:
		mb.ItemInfo = v
	case ItemPropertiesBox:
		mb.Properties = v
	case ItemLocationBox:
		mb.Location = v
	default:
		//mb.Children = append(mb.Children, b)
	}
	return nil
}

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

		//if Debug {
		//	fmt.Println(inner, outer.r.remain, inner.r.remain, inner.size)
		//}
	}
	//mb.Children, err = fb.parseAppendBoxes()
	return mb, err
}

// PrimaryItemBox is a "pitm" box
type PrimaryItemBox struct {
	Flags  Flags
	ItemID uint16
}

// Size returns the size of the PrimaryItemBox
func (pitm PrimaryItemBox) String() string {
	return fmt.Sprintf("pitm | ItemID: %d, Flags: %d, Version: %d ", pitm.ItemID, pitm.Flags.Flags(), pitm.Flags.Version())
}

// Type returns TypePitm
func (pitm PrimaryItemBox) Type() BoxType {
	return TypePitm
}

func parsePrimaryItemBox(outer *box) (Box, error) {
	flags, err := outer.r.readFlags()
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
