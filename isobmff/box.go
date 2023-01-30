package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// box is an ISOBMFF box
type box struct {
	size    int64
	remain  int
	offset  int
	flags   flags
	boxType BoxType
	outer   *box
	reader  *Reader
}

func (b box) Size() int64            { return b.size }
func (b box) Type() BoxType          { return b.boxType }
func (b box) isType(bt BoxType) bool { return b.boxType == bt }
func (b box) String() string {
	return fmt.Sprintf("(Box) type:'%s', offset:%d, size:%d", b.boxType, b.offset, b.size)
}

func (b *box) Peek(n int) ([]byte, error) {
	if b.remain >= n {
		if b.outer != nil {
			return b.outer.Peek(n)
		}
		return b.reader.peek(n)
	}
	return nil, ErrRemainLengthInsufficient
}

func (b *box) Discard(n int) error {
	if b.remain >= n {
		b.remain -= n
		if b.outer != nil {
			return b.outer.Discard(n)
		}
		return b.reader.discard(n)
	}
	return ErrRemainLengthInsufficient
}

func (b *box) adjust(n int) {
	b.remain -= n
	if b.outer != nil {
		b.outer.adjust(n)
	}
}

func (b *box) Read(p []byte) (n int, err error) {
	if b.remain >= len(p) {
		//fmt.Println(b.remain)
		n, err = b.reader.br.Read(p)
		b.adjust(n)
		return n, err
	}
	return 0, ErrRemainLengthInsufficient
}

func (b *box) close() error {
	if b.remain == 0 {
		return nil
	}
	return b.Discard(b.remain)
}

func (b *box) readInnerBox() (inner box, next bool, err error) {
	if b.remain < 16 {
		return inner, false, nil
	}
	// Read box size and box type
	var buf []byte
	if buf, err = b.Peek(8); err != nil {
		return inner, false, errors.Wrap(ErrBufLength, "readBox")
	}
	inner.reader = b.reader
	inner.outer = b
	inner.size = int64(bmffEndian.Uint32(buf[:4]))
	inner.remain = int(inner.size)
	inner.boxType = boxTypeFromBuf(buf[4:8])
	inner.offset = int(b.size) - b.remain + b.offset

	switch inner.size {
	case 1:
		// 1 means it's actually a 64-bit size, after the type.
		if buf, err = b.Peek(16); err != nil {
			return inner, false, errors.Wrap(ErrBufLength, "readBox")
		}
		inner.size = int64(bmffEndian.Uint32(buf[8:16]))
		if inner.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return inner, false, errors.Wrapf(errLargeBox, "readBox '%s'", inner.boxType)
		}
		inner.remain = int(inner.size)
		return inner, true, inner.Discard(16)
		//case 0:
		// 0 means unknown & to read to end of file. No more boxes.
		// r.noMoreBoxes = true
		// TODO: error
	}
	return inner, true, inner.Discard(8)
}

// readUint16 from box
func (b *box) readUint16() (uint16, error) {
	buf, err := b.Peek(2)
	if err != nil {
		return 0, errors.Wrap(ErrBufLength, "readUint16")
	}
	return bmffEndian.Uint16(buf[:2]), b.Discard(2)
}

// readUUID reads a 16 byte UUID from the box.
func (b *box) readUUID() (u meta.UUID, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		return u, errors.Wrap(ErrBufLength, "readUUID")
	}
	if err = u.UnmarshalBinary(buf); err != nil {
		return u, err
	}
	return u, b.Discard(16)
}

// BoxType is an ISOBMFF box
type BoxType uint8

// String is Stringer interface for boxType
func (t BoxType) String() string {
	str, ok := mapBoxTypeString[t]
	if ok {
		return str
	}
	return "nnnn"
}

func boxTypeFromBuf(buf []byte) BoxType {
	if string(buf[:4]) == "infe" { // inital check for performance reasons
		return typeInfe
	}
	if b, ok := mapStringBoxType[string(buf)]; ok {
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
		return errors.Wrap(ErrBufLength, "readFlags")
	}
	b.readFlagsFromBuf(buf)
	return b.Discard(4)
}

func (b *box) readFlagsFromBuf(buf []byte) {
	b.flags = flags(bmffEndian.Uint32(buf[:4]))
}

// Flags returns underlying Flags after removing version.
// Flags are 24 bits.
func (f flags) Flags() uint32 {
	// Left Shift
	f = f << 8
	// Right Shift
	return uint32(f >> 8)
}

// Version returns a uint8 version.
func (f flags) Version() uint8 {
	return uint8(f >> 24)
}

// String is Stringer interface for Flags
func (f flags) String() string {
	return fmt.Sprintf("Flags:%d, Version:%d", f.Flags(), f.Version())
}

// func (f flags)
// MarshalZerologObject is a zerolog interface for logging
func (f flags) MarshalZerologObject(e *zerolog.Event) {
	e.Uint8("version", f.Version()).Uint32("flags", f.Flags())
}

// Common box types.
const (
	typeUnknown BoxType = iota
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
)

var mapStringBoxType = map[string]BoxType{
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
}

var mapBoxTypeString = map[BoxType]string{
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
}
