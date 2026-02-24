package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/rs/zerolog"
)

// box is an ISOBMFF box
type box struct {
	size    int64
	remain  int
	offset  int
	flags   flags
	boxType boxType
	outer   *box
	reader  *Reader
}

// isType returns the boxType
func (b box) isType(bt boxType) bool { return b.boxType == bt }

// Peek returns []byte without advancing the reader. Is limited by the
// constrains of the box.
func (b *box) Peek(n int) ([]byte, error) {
	if b.remain >= n {
		if b.outer != nil {
			return b.outer.Peek(n)
		}
		return b.reader.peek(n)
	}
	return nil, ErrRemainLengthInsufficient
}

// Discard advances the reader. Is limited by the
// constrains of the box.
func (b *box) Discard(n int) (int, error) {
	if b.remain >= n {
		b.remain -= n
		if b.outer != nil {
			return b.outer.Discard(n)
		}
		return b.reader.discard(n)
	}
	return 0, ErrRemainLengthInsufficient
}

// Read the bytes from underlying reader. Is limited by the
// constrains of the box
func (b *box) Read(p []byte) (n int, err error) {
	if b.remain >= len(p) {
		//fmt.Println(b.remain)
		n, err = b.reader.br.Read(p)
		b.adjust(n)
		return n, err
	}
	return 0, ErrRemainLengthInsufficient
}

func (b *box) adjust(n int) {
	if n > b.remain {
		n = b.remain
	}
	b.remain -= n
	if b.outer != nil {
		b.outer.adjust(n)
	}
}

func (b *box) close() error {
	if b.remain == 0 {
		return nil
	}
	_, err := b.Discard(b.remain)
	return err
}

func (b *box) readInnerBox() (inner box, next bool, err error) {
	if b.remain < 8 {
		return inner, false, nil
	}
	fmt.Println("Read Box", b.boxType)
	buf, err := b.Peek(8)
	if err != nil {
		return inner, false, fmt.Errorf("readBox: %w", ErrBufLength)
	}
	inner.reader = b.reader
	inner.outer = b
	inner.offset = int(b.size) - b.remain + b.offset
	// Read box size and box type
	inner.size = int64(bmffEndian.Uint32(buf[:4]))
	inner.boxType = boxTypeFromBuf(buf[4:8])

	headerSize := 8
	if inner.size == 1 {
		buf, err = b.Peek(16)
		if err != nil {
			return inner, false, fmt.Errorf("readBox: %w", ErrBufLength)
		}

		// 1 means it's actually a 64-bit size, after the type.
		inner.size = int64(bmffEndian.Uint64(buf[8:16]))
		if inner.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return inner, false, fmt.Errorf("readBox '%s': %w", inner.boxType, errLargeBox)
		}
		headerSize = 16
	}

	if inner.size < int64(headerSize) {
		return inner, false, fmt.Errorf("readBox invalid size %d for '%s': %w", inner.size, inner.boxType, ErrBufLength)
	}

	maxInt := int64(^uint(0) >> 1)
	if inner.size > maxInt {
		return inner, false, fmt.Errorf("readBox '%s': %w", inner.boxType, errLargeBox)
	}

	inner.remain = int(inner.size)
	_, err = inner.Discard(headerSize)
	return inner, true, err
}

// readUint16 from box
func (b *box) readUint16() (uint16, error) {
	buf, err := b.Peek(2)
	if err != nil {
		return 0, fmt.Errorf("readUint16: %w", ErrBufLength)
	}
	_, err = b.Discard(2)
	return bmffEndian.Uint16(buf[:2]), err
}

// readUUID reads a 16 byte UUID from the box.
func (b *box) readUUID() (u meta.UUID, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		return u, fmt.Errorf("readUUID: %w", ErrBufLength)
	}
	if err = u.UnmarshalBinary(buf); err != nil {
		return u, err
	}
	_, err = b.Discard(16)
	return u, err
}

// MarshalZerologObject is a zerolog interface for logging
func (b box) MarshalZerologObject(e *zerolog.Event) {
	e.Str("boxType", b.boxType.String()).Int("offset", b.offset).Int64("size", b.size)
	if b.flags != 0 {
		e.Object("flags", b.flags)
	}
}

// BoxType is an ISOBMFF box
type boxType uint8

// String is Stringer interface for boxType
func (t boxType) String() string {
	str, ok := mapBoxTypeString[t]
	if ok {
		return str
	}
	return "nnnn"
}

func fourCCFromString(str string) uint32 {
	if len(str) < 4 {
		return 0
	}
	return uint32(str[0])<<24 | uint32(str[1])<<16 | uint32(str[2])<<8 | uint32(str[3])
}

var boxTypeInfeFourCC = fourCCFromString("infe")

func boxTypeFromBuf(buf []byte) boxType {
	if len(buf) < 4 {
		return typeUnknown
	}

	fourCC := bmffEndian.Uint32(buf[:4])
	if fourCC == boxTypeInfeFourCC { // initial check for performance reasons
		return typeInfe
	}

	if b, ok := mapFourCCBoxType[fourCC]; ok {
		return b
	}
	if logLevelError() {
		logErrorMsg("BoxType", "error BoxType '%s' unknown", buf)
	}
	return typeUnknown
}

// flags for a FullBox
// 8 bits -> Version
// 24 bits -> Flags
type flags uint32

// readFlags reads the Flags from a FullBox header.
func (b *box) readFlags() error {
	buf, err := b.Peek(4)
	if err != nil {
		return fmt.Errorf("readFlags: %w", ErrBufLength)
	}
	b.readFlagsFromBuf(buf)
	_, err = b.Discard(4)
	return err
}

func (b *box) readFlagsFromBuf(buf []byte) {
	b.flags = flags(bmffEndian.Uint32(buf[:4]))
}

// Flags returns underlying Flags after removing version.
// Flags are 24 bits.
func (f flags) flags() uint32 {
	// Left Shift
	f = f << 8
	// Right Shift
	return uint32(f >> 8)
}

// Version returns a uint8 version.
func (f flags) version() uint8 {
	return uint8(f >> 24)
}

// MarshalZerologObject is a zerolog interface for logging
func (f flags) MarshalZerologObject(e *zerolog.Event) {
	e.Uint8("version", f.version()).Uint32("flags", f.flags())
}

// Common box types.
const (
	typeUnknown boxType = iota
	typeAuxC            // 'auxC'
	typeAuxl            // 'auxl'
	typeAv01            // 'av01'
	typeAv1C            // 'av1C'
	typeAvcC            // 'avcC'
	typeCCDT            // 'CCDT'
	typeCCTP            // 'CCTP'
	typeCdsc            // 'cdsc'
	typeClap            // 'clap'
	typeCMT1            // 'CMT1'
	typeCMT2            // 'CMT2'
	typeCMT3            // 'CMT3'
	typeCMT4            // 'CMT4'
	typeCNCV            // 'CNCV'
	typeCo64            // 'co64'
	typeColr            // 'colr'
	typeCRAW            // 'CRAW'
	typeCrtt            // 'crtt'
	typeCTBO            // 'CTBO'
	typeCTMD            // 'CTMD'
	typeDimg            // 'dimg'
	typeDinf            // 'dinf'
	typeDref            // 'dref'
	typeEtyp            // 'etyp'
	typeFree            // 'free'
	typeFtyp            // 'ftyp'
	typeGrpl            // 'grpl'
	typeHdlr            // 'hdlr'
	typeHvcC            // 'hvcC'
	typeIdat            // 'idat'
	typeIinf            // 'iinf'
	typeIloc            // 'iloc'
	typeImir            // 'imir'
	typeInfe            // 'infe'
	typeIovl            // 'iovl'
	typeIpco            // 'ipco
	typeIpma            // 'ipma'
	typeIprp            // 'iprp'
	typeIref            // 'iref'
	typeIrot            // 'irot'
	typeIspe            // 'ispe'
	typeJXL             // 'JXL '
	typeJumb            // 'jumb'
	typeJxlc            // 'jxlc'
	typeJxll            // 'jxll'
	typeJxlp            // 'jxlp'
	typeLhvC            // 'lhvC'
	typeMdat            // 'mdat'
	typeMdft            // 'mdft'
	typeMdhd            // 'mdhd'
	typeMdia            // 'mdia'
	typeMeta            // 'meta'
	typeMinf            // 'minf'
	typeMoov            // 'moov'
	typeMvhd            // 'mvhd'
	typeNmhd            // 'nmhd'
	typeOinf            // 'oinf'
	typePasp            // 'pasp'
	typePitm            // 'pitm'
	typePixi            // 'pixi'
	typePRVW            // 'PRVW'
	typeStbl            // 'stbl'
	typeStsc            // 'stsc'
	typeStsd            // 'stsd'
	typeStsz            // 'stsz'
	typeStts            // 'stts'
	typeThmb            // 'thmb'
	typeTHMB            // 'THMB'
	typeTkhd            // 'tkhd'
	typeTols            // 'tols'
	typeTrak            // 'trak'
	typeUUID            // 'uuid'
	typeVmhd            // 'vmhd'
	typeExif            // 'Exif
)

var mapStringBoxType = map[string]boxType{
	"auxC": typeAuxC,
	"auxl": typeAuxl,
	"av01": typeAv01,
	"av1C": typeAv1C,
	"avcC": typeAvcC,
	"CCDT": typeCCDT,
	"CCTP": typeCCTP,
	"cdsc": typeCdsc,
	"clap": typeClap,
	"CMT1": typeCMT1,
	"CMT2": typeCMT2,
	"CMT3": typeCMT3,
	"CMT4": typeCMT4,
	"CNCV": typeCNCV,
	"co64": typeCo64,
	"colr": typeColr,
	"CRAW": typeCRAW,
	"crtt": typeCrtt,
	"CTBO": typeCTBO,
	"CTMD": typeCTMD,
	"dimg": typeDimg,
	"dinf": typeDinf,
	"dref": typeDref,
	"etyp": typeEtyp,
	"free": typeFree,
	"ftyp": typeFtyp,
	"grpl": typeGrpl,
	"hdlr": typeHdlr,
	"hvcC": typeHvcC,
	"idat": typeIdat,
	"iinf": typeIinf,
	"iloc": typeIloc,
	"imir": typeImir,
	"infe": typeInfe,
	"iovl": typeIovl,
	"ipco": typeIpco,
	"ipma": typeIpma,
	"iprp": typeIprp,
	"iref": typeIref,
	"irot": typeIrot,
	"ispe": typeIspe,
	"JXL ": typeJXL,
	"jumb": typeJumb,
	"jxlc": typeJxlc,
	"jxll": typeJxll,
	"jxlp": typeJxlp,
	"lhvC": typeLhvC,
	"mdat": typeMdat,
	"mdft": typeMdft,
	"mdhd": typeMdhd,
	"mdia": typeMdia,
	"meta": typeMeta,
	"minf": typeMinf,
	"moov": typeMoov,
	"mvhd": typeMvhd,
	"nmhd": typeNmhd,
	"oinf": typeOinf,
	"pasp": typePasp,
	"pitm": typePitm,
	"pixi": typePixi,
	"PRVW": typePRVW,
	"stbl": typeStbl,
	"stsc": typeStsc,
	"stsd": typeStsd,
	"stsz": typeStsz,
	"stts": typeStts,
	"thmb": typeThmb,
	"THMB": typeTHMB,
	"tkhd": typeTkhd,
	"tols": typeTols,
	"trak": typeTrak,
	"uuid": typeUUID,
	"vmhd": typeVmhd,
	"Exif": typeExif,
}

var mapFourCCBoxType = func() map[uint32]boxType {
	m := make(map[uint32]boxType, len(mapStringBoxType))
	for k, v := range mapStringBoxType {
		if len(k) == 4 {
			m[fourCCFromString(k)] = v
		}
	}
	return m
}()

var mapBoxTypeString = map[boxType]string{
	typeAuxC: "auxC",
	typeAuxl: "auxl",
	typeAv01: "av01",
	typeAv1C: "av1C",
	typeAvcC: "avcC",
	typeCCDT: "CCDT",
	typeCCTP: "CCTP",
	typeCdsc: "cdsc",
	typeClap: "clap",
	typeCMT1: "CMT1",
	typeCMT2: "CMT2",
	typeCMT3: "CMT3",
	typeCMT4: "CMT4",
	typeCNCV: "CNCV",
	typeCo64: "co64",
	typeColr: "colr",
	typeCRAW: "CRAW",
	typeCrtt: "crtt",
	typeCTBO: "CTBO",
	typeCTMD: "CTMD",
	typeDimg: "dimg",
	typeDinf: "dinf",
	typeDref: "dref",
	typeEtyp: "etyp",
	typeFree: "free",
	typeFtyp: "ftyp",
	typeGrpl: "grpl",
	typeHdlr: "hdlr",
	typeHvcC: "hvcC",
	typeIdat: "idat",
	typeIinf: "iinf",
	typeIloc: "iloc",
	typeImir: "imir",
	typeInfe: "infe",
	typeIovl: "iovl",
	typeIpco: "ipco",
	typeIpma: "ipma",
	typeIprp: "iprp",
	typeIref: "iref",
	typeIrot: "irot",
	typeIspe: "ispe",
	typeJXL:  "JXL ",
	typeJumb: "jumb",
	typeJxlc: "jxlc",
	typeJxll: "jxll",
	typeJxlp: "jxlp",
	typeLhvC: "lhvC",
	typeMdat: "mdat",
	typeMdft: "mdft",
	typeMdhd: "mdhd",
	typeMdia: "mdia",
	typeMeta: "meta",
	typeMinf: "minf",
	typeMoov: "moov",
	typeMvhd: "mvhd",
	typeNmhd: "nmhd",
	typeOinf: "oinf",
	typePasp: "pasp",
	typePitm: "pitm",
	typePixi: "pixi",
	typePRVW: "PRVW",
	typeStbl: "stbl",
	typeStsc: "stsc",
	typeStsd: "stsd",
	typeStsz: "stsz",
	typeStts: "stts",
	typeThmb: "thmb",
	typeTHMB: "THMB",
	typeTkhd: "tkhd",
	typeTols: "tols",
	typeTrak: "trak",
	typeUUID: "uuid",
	typeVmhd: "vmhd",
	typeExif: "Exif",
}
