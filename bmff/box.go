package bmff

import (
	"errors"
	"fmt"
	"io"
)

// ErrUnknownParser is returned by Box.Parse for unrecognized box parser.
var ErrUnknownParser = errors.New("error no parser for box")

// BoxType is an ISOBMFF box
type BoxType uint8

// Common box types.
const (
	TypeUnknown BoxType = iota
	TypeAuxC            // 'auxC'
	TypeAuxl            // 'auxl'
	TypeAv01            // 'av01'
	TypeAv1C            // 'av1C'
	TypeCdsc            // 'cdsc'
	TypeClap            // 'clap'
	TypeColr            // 'colr'
	TypeDimg            // 'dimg'
	TypeDinf            // 'dinf'
	TypeFtyp            // 'ftyp'
	TypeHdlr            // 'hdlr'
	TypeHvcC            // 'hvcC'
	TypeIdat            // 'idat'
	TypeIinf            // 'iinf'
	TypeIloc            // 'iloc'
	TypeImir            // 'imir'
	TypeInfe            // 'infe'
	TypeIovl            // 'iovl'
	TypeIpco            // 'ipco
	TypeIpma            // 'ipma'
	TypeIprp            // 'iprp'
	TypeIref            // 'iref'
	TypeIrot            // 'irot'
	TypeIspe            // 'ispe'
	TypeMeta            // 'meta'
	TypeMoov            // 'moov'
	TypePasp            // 'pasp'
	TypePitm            // 'pitm'
	TypePixi            // 'pixi'
	TypeThmb            // 'thmb'
	TypeTrak            // 'trak'
	TypeUUID            // 'uuid'
)

var mapStringBoxType = map[string]BoxType{
	"auxC": TypeAuxC,
	"auxl": TypeAuxl,
	"av01": TypeAv01,
	"av1C": TypeAv1C,
	"cdsc": TypeCdsc,
	"clap": TypeClap,
	"colr": TypeColr,
	"dimg": TypeDimg,
	"dinf": TypeDinf,
	"ftyp": TypeFtyp,
	"hdlr": TypeHdlr,
	"hvcC": TypeHvcC,
	"idat": TypeIdat,
	"iinf": TypeIinf,
	"iloc": TypeIloc,
	"imir": TypeImir,
	"infe": TypeInfe,
	"iovl": TypeIovl,
	"ipco": TypeIpco,
	"ipma": TypeIpma,
	"iprp": TypeIprp,
	"iref": TypeIref,
	"irot": TypeIrot,
	"ispe": TypeIspe,
	"meta": TypeMeta,
	"moov": TypeMoov,
	"pasp": TypePasp,
	"pitm": TypePitm,
	"pixi": TypePixi,
	"thmb": TypeThmb,
	"trak": TypeTrak,
	"uuid": TypeUUID,
}

var mapBoxTypeString = map[BoxType]string{
	TypeAuxC: "auxC",
	TypeAuxl: "auxl",
	TypeAv01: "av01",
	TypeAv1C: "av1C",
	TypeCdsc: "cdsc",
	TypeClap: "clap",
	TypeColr: "colr",
	TypeDimg: "dimg",
	TypeDinf: "dinf",
	TypeFtyp: "ftyp",
	TypeHdlr: "hdlr",
	TypeHvcC: "hvcC",
	TypeIdat: "idat",
	TypeIinf: "iinf",
	TypeIloc: "iloc",
	TypeImir: "imir",
	TypeInfe: "infe",
	TypeIovl: "iovl",
	TypeIpco: "ipco",
	TypeIpma: "ipma",
	TypeIprp: "iprp",
	TypeIref: "iref",
	TypeIrot: "irot",
	TypeIspe: "ispe",
	TypeMeta: "meta",
	TypeMoov: "moov",
	TypePasp: "pasp",
	TypePitm: "pitm",
	TypePixi: "pixi",
	TypeThmb: "thmb",
	TypeTrak: "trak",
	TypeUUID: "uuid",
}

func (t BoxType) String() string {
	str, ok := mapBoxTypeString[t]
	if ok {
		return str
	}
	return "nnnn"
}

func isBoxinfe(buf []byte) bool {
	return buf[0] == 'i' && buf[1] == 'n' && buf[2] == 'f' && buf[3] == 'e'
}

func boxType(buf []byte) BoxType {
	if isBoxinfe(buf) { // inital check for performance reasons, TODO: confirm with benchmarks
		return TypeInfe
	}
	b, ok := mapStringBoxType[string(buf)]
	if ok {
		return b
	}
	if Debug {
		fmt.Println(string(buf))
	}
	return TypeUnknown
}

type box struct {
	bufReader
	size    int64 // 0 means unknown, will read to end of file (box container)
	boxType BoxType
}

func (b box) String() string {
	return fmt.Sprintf("(Box) type: \"%s\", size: %d", b.boxType, b.size)
}

func (b box) Size() int64   { return b.size }
func (b box) Type() BoxType { return b.boxType }

func (b *box) Parse() (Box, error) {
	if !b.anyRemain() {
		return nil, ErrBufReaderLength
	}
	parser, ok := parsers[b.Type()]
	if !ok {
		return UnknownBox{t: b.Type(), s: b.size}, ErrUnknownParser
	}

	v, err := parser(b)
	if err != nil {
		// Write error with parser
		return nil, err
	}
	return v, nil
}

// Parsers
type parserFunc func(b *box) (Box, error)

var parsers map[BoxType]parserFunc

func init() {
	parsers = map[BoxType]parserFunc{
		TypeDinf: parseDataInformationBox,
		//boxType("dref"): parseDataReferenceBox,
		TypeFtyp: parseFtyp,
		TypeHdlr: parseHdlr,
		TypeIinf: parseIinf,
		TypeInfe: parseInfe,
		TypeIloc: parseIloc,
		TypeIpco: parseIpco,
		TypeIpma: parseIpma,
		TypeIprp: parseIprp,
		TypeIref: parseIref,
		TypeIrot: parseImageRotation,
		TypeIspe: parseImageSpatialExtentsProperty,
		TypeMeta: parseMeta,
		TypePitm: parsePitm,
	}
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

func parseMeta(outer *box) (Box, error) {
	return parseMetaBox(outer)
}

func parseMetaBox(outer *box) (mb MetaBox, err error) {
	mb = MetaBox{size: uint32(outer.size)}
	mb.Flags, err = outer.readFlags()
	if err != nil {
		return mb, err
	}

	var inner box
	for outer.anyRemain() {
		inner, err = outer.readInnerBox()
		if err != nil {
			if err == io.EOF {
				return mb, nil
			}
			return mb, err
		}
		switch inner.boxType {
		case TypeIdat, TypeDinf, TypeUUID, TypeIref:
			// Do not parse

		//case TypeIref:
		//	_, err = inner.Parse()
		case TypePitm:
			mb.Primary, err = parsePrimaryItemBox(&inner)
		case TypeIinf:
			mb.ItemInfo, err = parseItemInfoBox(&inner)
		case TypeHdlr:
			mb.Handler, err = parseHandlerBox(&inner)
		case TypeIprp:
			mb.Properties, err = parseItemPropertiesBox(&inner)
		case TypeIloc:
			mb.Location, err = parseItemLocationBox(&inner)
		default:
			p, err := inner.Parse()
			if err == nil {
				mb.Children = append(mb.Children, p)
			}
		}
		if err != nil {
			if Debug {
				fmt.Println(err)
			}
		}
		outer.remain -= int(inner.size)
		if err = inner.discard(inner.remain); err != nil {
			if Debug {
				fmt.Println(err)
			}
			// TODO: improve error handling
			break
		}

		if Debug {
			fmt.Println(inner, outer.remain, inner.remain, inner.size)
		}
	}
	err = outer.discard(outer.remain)
	return mb, err
}

// DataInformationBox is a "dinf" box
type DataInformationBox struct {
	Children []Box
}

// Type returns TypeDinf
func (dinf DataInformationBox) Type() BoxType {
	return TypeDinf
}

func parseDataInformationBox(outer *box) (Box, error) {
	dib := DataInformationBox{}
	err := outer.discard(outer.remain)
	return dib, err //br.parseAppendBoxes(&dib.Children)
}

// UnknownBox is a box that was unable to be parsed.
type UnknownBox struct {
	t BoxType
	s int64
}

// Type returns the BoxType of the UnknownBox
func (ub UnknownBox) Type() BoxType {
	return ub.t
}

// Size returns the Size of the UnknownBox
func (ub UnknownBox) Size() int64 {
	return ub.s
}

func (ub UnknownBox) String() string {
	return fmt.Sprintf(" Type: %s, Size: %d", ub.t, ub.s)
}
