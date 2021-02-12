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
	err     error
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
		return UnknownBox{t: b.Type(), s: b.size}, ErrUnknownParser //ErrUnknownBox
	}

	v, err := parser(b)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (b *box) discardRemaining(inner box) {

}

// Parsers
type parserFunc func(b *box) (Box, error)

var parsers map[BoxType]parserFunc

func init() {
	parsers = map[BoxType]parserFunc{
		TypeDinf: parseDataInformationBox,
		//boxType("dref"): parseDataReferenceBox,
		TypeFtyp: parseFileTypeBox,
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
		TypeMeta: parseMetaBox,
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
		switch inner.boxType {
		case TypeIdat, TypeDinf:
			// Do not read
			boxr.br.discard(inner.r.remain)
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
			boxr.br.discard(inner.r.remain)
		}
		if err != nil {
			boxr.br.discard(inner.r.remain)
		}
		outer.r.remain -= int(inner.size)

		//if inner.boxType == TypeIdat || inner.boxType == TypeIref || inner.boxType == TypeDinf {
		//	boxr.br.discard(int(inner.r.remain))
		//	outer.r.remain -= int(inner.size)
		//	continue
		//}
		//if inner.boxType == TypePitm {
		//	mb.Primary, err = parsePrimaryItemBox(&inner)
		//	if err != nil {
		//		break
		//	}
		//} else if inner.boxType == TypeIinf {
		//	mb.ItemInfo, err = parseItemInfoBox(&inner)
		//	if err != nil {
		//		break
		//	}
		//} else if inner.boxType == TypeHdlr {
		//	mb.Handler, err = parseHandlerBox(&inner)
		//	if err != nil {
		//		break
		//	}
		//} else {
		//	p, err := inner.Parse()
		//	if err != nil {
		//		boxr.br.discard(int(inner.r.remain))
		//	} else {
		//
		//	}
		//}
		//
		//outer.r.remain -= int(inner.size)

		if Debug {
			fmt.Println(inner, outer.r.remain, inner.r.remain, inner.size)
		}
	}
	boxr.br.discard(outer.r.remain)
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
	outer.r.discard(int(outer.r.remain))
	return dib, nil //br.parseAppendBoxes(&dib.Children)
}
