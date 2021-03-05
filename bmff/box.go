package bmff

import (
	"errors"
	"fmt"
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
	TypeAvcC            // 'avcC'
	TypeCCDT            // 'CCDT'
	TypeCCTP            // 'CCTP'
	TypeCdsc            // 'cdsc'
	TypeClap            // 'clap'
	TypeCMT1            // 'CMT1'
	TypeCMT2            // 'CMT2'
	TypeCMT3            // 'CMT3'
	TypeCMT4            // 'CMT4'
	TypeCNCV            // 'CNCV'
	TypeColr            // 'colr'
	TypeCrtt            // 'crtt'
	TypeCTBO            // 'CTBO'
	TypeDimg            // 'dimg'
	TypeDinf            // 'dinf'
	TypeEtyp            // 'etyp'
	TypeFree            // 'free'
	TypeFtyp            // 'ftyp'
	TypeGrpl            // 'grpl'
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
	TypeLhvC            // 'lhvC'
	TypeMdat            // 'mdat'
	TypeMdft            // 'mdft'
	TypeMeta            // 'meta'
	TypeMoov            // 'moov'
	TypeMvhd            // 'mvhd'
	TypeOinf            // 'oinf'
	TypePasp            // 'pasp'
	TypePitm            // 'pitm'
	TypePixi            // 'pixi'
	TypePRVW            // 'PRVW'
	TypeThmb            // 'thmb'
	TypeTHMB            // 'THMB'
	TypeTols            // 'tols'
	TypeTrak            // 'trak'
	TypeUUID            // 'uuid'
)

var mapStringBoxType = map[string]BoxType{
	"auxC": TypeAuxC,
	"auxl": TypeAuxl,
	"av01": TypeAv01,
	"av1C": TypeAv1C,
	"avcC": TypeAvcC,
	"CCDT": TypeCCDT,
	"CCTP": TypeCCTP,
	"cdsc": TypeCdsc,
	"clap": TypeClap,
	"CMT1": TypeCMT1,
	"CMT2": TypeCMT2,
	"CMT3": TypeCMT3,
	"CMT4": TypeCMT4,
	"CNCV": TypeCNCV,
	"colr": TypeColr,
	"crtt": TypeCrtt,
	"CTBO": TypeCTBO,
	"dimg": TypeDimg,
	"dinf": TypeDinf,
	"etyp": TypeEtyp,
	"free": TypeFree,
	"ftyp": TypeFtyp,
	"grpl": TypeGrpl,
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
	"lhvC": TypeLhvC,
	"mdat": TypeMdat,
	"mdft": TypeMdft,
	"meta": TypeMeta,
	"moov": TypeMoov,
	"mvhd": TypeMvhd,
	"oinf": TypeOinf,
	"pasp": TypePasp,
	"pitm": TypePitm,
	"pixi": TypePixi,
	"PRVW": TypePRVW,
	"thmb": TypeThmb,
	"THMB": TypeTHMB,
	"tols": TypeTols,
	"trak": TypeTrak,
	"uuid": TypeUUID,
}

var mapBoxTypeString = map[BoxType]string{
	TypeAuxC: "auxC",
	TypeAuxl: "auxl",
	TypeAv01: "av01",
	TypeAv1C: "av1C",
	TypeAvcC: "avcC",
	TypeCCDT: "CCDT",
	TypeCCTP: "CCTP",
	TypeCdsc: "cdsc",
	TypeClap: "clap",
	TypeCMT1: "CMT1",
	TypeCMT2: "CMT2",
	TypeCMT3: "CMT3",
	TypeCMT4: "CMT4",
	TypeCNCV: "CNCV",
	TypeColr: "colr",
	TypeCrtt: "crtt",
	TypeCTBO: "CTBO",
	TypeDimg: "dimg",
	TypeDinf: "dinf",
	TypeEtyp: "etyp",
	TypeFree: "free",
	TypeFtyp: "ftyp",
	TypeGrpl: "grpl",
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
	TypeLhvC: "lhvC",
	TypeMdat: "mdat",
	TypeMdft: "mdft",
	TypeMeta: "meta",
	TypeMoov: "moov",
	TypeMvhd: "mvhd",
	TypeOinf: "oinf",
	TypePasp: "pasp",
	TypePitm: "pitm",
	TypePixi: "pixi",
	TypePRVW: "PRVW",
	TypeThmb: "thmb",
	TypeTHMB: "THMB",
	TypeTols: "tols",
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
	if debugFlag {
		log.Debug("BoxType '%s' not found", buf)
	}
	return TypeUnknown
}

// Box is an interface for different BMFF boxes.
type Box interface {
	//Size() int64 // 0 means unknown (will read to end of file)
	Type() BoxType
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

func (f Flags) String() string {
	return fmt.Sprintf("Flags:%d, Version:%d", f.Flags(), f.Version())
}

// Box is a BMFF box
type box struct {
	bufReader
	size    int64 // 0 means unknown, will read to end of file (box container)
	boxType BoxType
}

func (b box) String() string {
	return fmt.Sprintf("(Box) type:'%s', offset:%d, size:%d", b.boxType, b.offset, b.size)
}

func (b box) Size() int64   { return b.size }
func (b box) Type() BoxType { return b.boxType }

func (b *box) Parse() (Box, error) {
	if !b.anyRemain() {
		return nil, ErrBufLength
	}
	parser, ok := parsers[b.Type()]
	if !ok {
		if debugFlag {
			log.Debug("Unknown Parser. Box:'%s', BoxSize:'%d'", b.Type(), b.Size())
		}
		return UnknownBox{t: b.Type(), s: b.size}, nil
	}
	return parser(b)
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
		TypeIrot: parseIrot,
		TypeIspe: parseIspe,
		TypeMeta: parseMeta,
		TypeMoov: parseMoov,
		TypePitm: parsePitm,
		TypeHvcC: parseUnknownBox,
		TypeColr: parseUnknownBox,
		TypePixi: parseUnknownBox,
	}
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
	//br.parseAppendBoxes(&dib.Children)
	return dib, outer.discard(outer.remain)
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
	return fmt.Sprintf("Type: %s, Size: %d", ub.t, ub.s)
}

// Process Boxes that are not implemented
func parseUnknownBox(outer *box) (Box, error) {
	return UnknownBox{outer.boxType, outer.size}, outer.discard(outer.remain)
}
