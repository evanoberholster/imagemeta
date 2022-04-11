package bmff

import (
	"fmt"
)

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
	TypeCo64            // 'co64'
	TypeColr            // 'colr'
	TypeCRAW            // 'CRAW'
	TypeCrtt            // 'crtt'
	TypeCTBO            // 'CTBO'
	TypeCTMD            // 'CTMD'
	TypeDimg            // 'dimg'
	TypeDinf            // 'dinf'
	TypeDref            // 'dref'
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
	TypeMdhd            // 'mdhd'
	TypeMdia            // 'mdia'
	TypeMeta            // 'meta'
	TypeMinf            // 'minf'
	TypeMoov            // 'moov'
	TypeMvhd            // 'mvhd'
	TypeNmhd            // 'nmhd'
	TypeOinf            // 'oinf'
	TypePasp            // 'pasp'
	TypePitm            // 'pitm'
	TypePixi            // 'pixi'
	TypePRVW            // 'PRVW'
	TypeStbl            // 'stbl'
	TypeStsc            // 'stsc'
	TypeStsd            // 'stsd'
	TypeStsz            // 'stsz'
	TypeStts            // 'stts'
	TypeThmb            // 'thmb'
	TypeTHMB            // 'THMB'
	TypeTkhd            // 'tkhd'
	TypeTols            // 'tols'
	TypeTrak            // 'trak'
	TypeUUID            // 'uuid'
	TypeVmhd            // 'vmhd'
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
	"co64": TypeCo64,
	"colr": TypeColr,
	"CRAW": TypeCRAW,
	"crtt": TypeCrtt,
	"CTBO": TypeCTBO,
	"CTMD": TypeCTMD,
	"dimg": TypeDimg,
	"dinf": TypeDinf,
	"dref": TypeDref,
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
	"mdhd": TypeMdhd,
	"mdia": TypeMdia,
	"meta": TypeMeta,
	"minf": TypeMinf,
	"moov": TypeMoov,
	"mvhd": TypeMvhd,
	"nmhd": TypeNmhd,
	"oinf": TypeOinf,
	"pasp": TypePasp,
	"pitm": TypePitm,
	"pixi": TypePixi,
	"PRVW": TypePRVW,
	"stbl": TypeStbl,
	"stsc": TypeStsc,
	"stsd": TypeStsd,
	"stsz": TypeStsz,
	"stts": TypeStts,
	"thmb": TypeThmb,
	"THMB": TypeTHMB,
	"tkhd": TypeTkhd,
	"tols": TypeTols,
	"trak": TypeTrak,
	"uuid": TypeUUID,
	"vmhd": TypeVmhd,
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
	TypeCo64: "co64",
	TypeColr: "colr",
	TypeCRAW: "CRAW",
	TypeCrtt: "crtt",
	TypeCTBO: "CTBO",
	TypeCTMD: "CTMD",
	TypeDimg: "dimg",
	TypeDinf: "dinf",
	TypeDref: "dref",
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
	TypeMdhd: "mdhd",
	TypeMdia: "mdia",
	TypeMeta: "meta",
	TypeMinf: "minf",
	TypeMoov: "moov",
	TypeMvhd: "mvhd",
	TypeNmhd: "nmhd",
	TypeOinf: "oinf",
	TypePasp: "pasp",
	TypePitm: "pitm",
	TypePixi: "pixi",
	TypePRVW: "PRVW",
	TypeStbl: "stbl",
	TypeStsc: "stsc",
	TypeStsd: "stsd",
	TypeStsz: "stsz",
	TypeStts: "stts",
	TypeThmb: "thmb",
	TypeTHMB: "THMB",
	TypeTkhd: "tkhd",
	TypeTols: "tols",
	TypeTrak: "trak",
	TypeUUID: "uuid",
	TypeVmhd: "vmhd",
}

func (t BoxType) String() string {
	str, ok := mapBoxTypeString[t]
	if ok {
		return str
	}
	return "nnnn"
}

func boxType(buf []byte) BoxType {
	if string(buf[:4]) == "infe" { // inital check for performance reasons
		return TypeInfe
	}
	if b, ok := mapStringBoxType[string(buf)]; ok {
		return b
	}
	if debugFlag {
		log.Debug("BoxType '%s' not found", buf)
	}
	return TypeUnknown
}

// Box is an interface for different BMFF boxes.
type Box interface {
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

func (b *box) read() (buf []byte, err error) {
	if buf, err = b.peek(b.remain); err != nil {
		return
	}
	err = b.discard(len(buf))
	return
}

func (b box) Size() int64   { return b.size }
func (b box) Type() BoxType { return b.boxType }
func (b box) String() string {
	return fmt.Sprintf("(Box) type:'%s', offset:%d, size:%d", b.boxType, b.offset, b.size)
}

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
		//boxType("dref"): parseDataReferenceBox,
		TypeDinf: parseDinf,
		//TypeFtyp: parseFtyp,
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

func parseDinf(outer *box) (Box, error) {
	return outer.parseDataInformationBox()
}

func (b *box) parseDataInformationBox() (Box, error) {
	dib := DataInformationBox{}
	return dib, b.discard(b.remain)
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
