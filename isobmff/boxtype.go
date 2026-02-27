package isobmff

import (
	"fmt"
)

// boxType identifies an ISOBMFF box by its FourCC.
type boxType uint8

// String returns the canonical FourCC text for a box type.
func (t boxType) String() string {
	switch t {
	case typeAuxC:
		return "auxC"
	case typeAuxl:
		return "auxl"
	case typeAv01:
		return "av01"
	case typeAv1C:
		return "av1C"
	case typeAvcC:
		return "avcC"
	case typeCCDT:
		return "CCDT"
	case typeCCTP:
		return "CCTP"
	case typeCdsc:
		return "cdsc"
	case typeClap:
		return "clap"
	case typeCMT1:
		return "CMT1"
	case typeCMT2:
		return "CMT2"
	case typeCMT3:
		return "CMT3"
	case typeCMT4:
		return "CMT4"
	case typeCNCV:
		return "CNCV"
	case typeCo64:
		return "co64"
	case typeColr:
		return "colr"
	case typeCRAW:
		return "CRAW"
	case typeCrtt:
		return "crtt"
	case typeCTBO:
		return "CTBO"
	case typeCTMD:
		return "CTMD"
	case typeDimg:
		return "dimg"
	case typeDinf:
		return "dinf"
	case typeDref:
		return "dref"
	case typeEtyp:
		return "etyp"
	case typeFree:
		return "free"
	case typeFtyp:
		return "ftyp"
	case typeGrpl:
		return "grpl"
	case typeHdlr:
		return "hdlr"
	case typeHvcC:
		return "hvcC"
	case typeIdat:
		return "idat"
	case typeIinf:
		return "iinf"
	case typeIloc:
		return "iloc"
	case typeImir:
		return "imir"
	case typeInfe:
		return "infe"
	case typeIovl:
		return "iovl"
	case typeIpco:
		return "ipco"
	case typeIpma:
		return "ipma"
	case typeIprp:
		return "iprp"
	case typeIref:
		return "iref"
	case typeIrot:
		return "irot"
	case typeIspe:
		return "ispe"
	case typeJXL:
		return "JXL "
	case typeJumb:
		return "jumb"
	case typeJxlc:
		return "jxlc"
	case typeJxll:
		return "jxll"
	case typeJxlp:
		return "jxlp"
	case typeLhvC:
		return "lhvC"
	case typeMdat:
		return "mdat"
	case typeMdft:
		return "mdft"
	case typeMdhd:
		return "mdhd"
	case typeMdia:
		return "mdia"
	case typeMeta:
		return "meta"
	case typeMinf:
		return "minf"
	case typeMoov:
		return "moov"
	case typeMvhd:
		return "mvhd"
	case typeNmhd:
		return "nmhd"
	case typeOinf:
		return "oinf"
	case typePasp:
		return "pasp"
	case typePitm:
		return "pitm"
	case typePixi:
		return "pixi"
	case typePRVW:
		return "PRVW"
	case typeStbl:
		return "stbl"
	case typeStsc:
		return "stsc"
	case typeStsd:
		return "stsd"
	case typeStsz:
		return "stsz"
	case typeStts:
		return "stts"
	case typeThmb:
		return "thmb"
	case typeTHMB:
		return "THMB"
	case typeTkhd:
		return "tkhd"
	case typeTols:
		return "tols"
	case typeTrak:
		return "trak"
	case typeUUID:
		return "uuid"
	case typeVmhd:
		return "vmhd"
	case typeExif:
		return "Exif"
	default:
		return "nnnn"
	}
}

// fourCCFromString packs the first 4 bytes of a string into a big-endian
// uint32 FourCC value used by BMFF box and handler types.
func fourCCFromString(str string) uint32 {
	if len(str) < 4 {
		return 0
	}
	return uint32(str[0])<<24 | uint32(str[1])<<16 | uint32(str[2])<<8 | uint32(str[3])
}

// boxTypeFromBuf maps a 4-byte box type code to an internal boxType enum.
// Unknown types are allowed and returned as typeUnknown.
func boxTypeFromBuf(buf []byte) boxType {
	if len(buf) < 4 {
		return typeUnknown
	}

	switch bmffEndian.Uint32(buf[:4]) {
	case boxTypeAuxCFourCC:
		return typeAuxC
	case boxTypeAuxlFourCC:
		return typeAuxl
	case boxTypeAv01FourCC:
		return typeAv01
	case boxTypeAv1CFourCC:
		return typeAv1C
	case boxTypeAvcCFourCC:
		return typeAvcC
	case boxTypeCCDTFourCC:
		return typeCCDT
	case boxTypeCCTPFourCC:
		return typeCCTP
	case boxTypeCdscFourCC:
		return typeCdsc
	case boxTypeClapFourCC:
		return typeClap
	case boxTypeCMT1FourCC:
		return typeCMT1
	case boxTypeCMT2FourCC:
		return typeCMT2
	case boxTypeCMT3FourCC:
		return typeCMT3
	case boxTypeCMT4FourCC:
		return typeCMT4
	case boxTypeCNCVFourCC:
		return typeCNCV
	case boxTypeCo64FourCC:
		return typeCo64
	case boxTypeColrFourCC:
		return typeColr
	case boxTypeCRAWFourCC:
		return typeCRAW
	case boxTypeCrttFourCC:
		return typeCrtt
	case boxTypeCTBOFourCC:
		return typeCTBO
	case boxTypeCTMDFourCC:
		return typeCTMD
	case boxTypeDimgFourCC:
		return typeDimg
	case boxTypeDinfFourCC:
		return typeDinf
	case boxTypeDrefFourCC:
		return typeDref
	case boxTypeEtypFourCC:
		return typeEtyp
	case boxTypeFreeFourCC:
		return typeFree
	case boxTypeFtypFourCC:
		return typeFtyp
	case boxTypeGrplFourCC:
		return typeGrpl
	case boxTypeHdlrFourCC:
		return typeHdlr
	case boxTypeHvcCFourCC:
		return typeHvcC
	case boxTypeIdatFourCC:
		return typeIdat
	case boxTypeIinfFourCC:
		return typeIinf
	case boxTypeIlocFourCC:
		return typeIloc
	case boxTypeImirFourCC:
		return typeImir
	case boxTypeInfeFourCC:
		return typeInfe
	case boxTypeIovlFourCC:
		return typeIovl
	case boxTypeIpcoFourCC:
		return typeIpco
	case boxTypeIpmaFourCC:
		return typeIpma
	case boxTypeIprpFourCC:
		return typeIprp
	case boxTypeIrefFourCC:
		return typeIref
	case boxTypeIrotFourCC:
		return typeIrot
	case boxTypeIspeFourCC:
		return typeIspe
	case boxTypeJXLFourCC:
		return typeJXL
	case boxTypeJumbFourCC:
		return typeJumb
	case boxTypeJxlcFourCC:
		return typeJxlc
	case boxTypeJxllFourCC:
		return typeJxll
	case boxTypeJxlpFourCC:
		return typeJxlp
	case boxTypeLhvCFourCC:
		return typeLhvC
	case boxTypeMdatFourCC:
		return typeMdat
	case boxTypeMdftFourCC:
		return typeMdft
	case boxTypeMdhdFourCC:
		return typeMdhd
	case boxTypeMdiaFourCC:
		return typeMdia
	case boxTypeMetaFourCC:
		return typeMeta
	case boxTypeMinfFourCC:
		return typeMinf
	case boxTypeMoovFourCC:
		return typeMoov
	case boxTypeMvhdFourCC:
		return typeMvhd
	case boxTypeNmhdFourCC:
		return typeNmhd
	case boxTypeOinfFourCC:
		return typeOinf
	case boxTypePaspFourCC:
		return typePasp
	case boxTypePitmFourCC:
		return typePitm
	case boxTypePixiFourCC:
		return typePixi
	case boxTypePRVWFourCC:
		return typePRVW
	case boxTypeStblFourCC:
		return typeStbl
	case boxTypeStscFourCC:
		return typeStsc
	case boxTypeStsdFourCC:
		return typeStsd
	case boxTypeStszFourCC:
		return typeStsz
	case boxTypeSttsFourCC:
		return typeStts
	case boxTypeThmbFourCC:
		return typeThmb
	case boxTypeTHMBFourCC:
		return typeTHMB
	case boxTypeTkhdFourCC:
		return typeTkhd
	case boxTypeTolsFourCC:
		return typeTols
	case boxTypeTrakFourCC:
		return typeTrak
	case boxTypeUUIDFourCC:
		return typeUUID
	case boxTypeVmhdFourCC:
		return typeVmhd
	case boxTypeExifFourCC:
		return typeExif
	}
	if logLevelDebug() {
		logDebug().Str("boxType", string(buf[:4])).Msg("unknown box type")
	}
	return typeUnknown
}

// flags stores the 32-bit FullBox field:
// upper 8 bits are version, lower 24 bits are flags.
type flags uint32

const (
	boxStringReadChunk = 4096
	maxBoxStringLength = 64 * 1024
)

// readFlags parses and consumes a FullBox version/flags field.
func (b *box) readFlags() error {
	buf, err := b.Peek(4)
	if err != nil {
		return fmt.Errorf("readFlags: %w", ErrBufLength)
	}
	b.readFlagsFromBuf(buf)
	_, err = b.Discard(4)
	return err
}

// readFlagsFromBuf decodes a FullBox 32-bit version/flags field
// (version in top 8 bits, flags in lower 24 bits).
func (b *box) readFlagsFromBuf(buf []byte) {
	b.flags = flags(bmffEndian.Uint32(buf[:4]))
}

// flags returns the lower 24-bit FullBox flags value.
func (f flags) flags() uint32 {
	return uint32(f) & 0x00FFFFFF
}

// version returns the FullBox version byte.
func (f flags) version() uint8 {
	return uint8(f >> 24)
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
	typeIpco            // 'ipco'
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
	typeExif            // 'Exif'
)

var (
	boxTypeAuxCFourCC = fourCCFromString("auxC")
	boxTypeAuxlFourCC = fourCCFromString("auxl")
	boxTypeAv01FourCC = fourCCFromString("av01")
	boxTypeAv1CFourCC = fourCCFromString("av1C")
	boxTypeAvcCFourCC = fourCCFromString("avcC")
	boxTypeCCDTFourCC = fourCCFromString("CCDT")
	boxTypeCCTPFourCC = fourCCFromString("CCTP")
	boxTypeCdscFourCC = fourCCFromString("cdsc")
	boxTypeClapFourCC = fourCCFromString("clap")
	boxTypeCMT1FourCC = fourCCFromString("CMT1")
	boxTypeCMT2FourCC = fourCCFromString("CMT2")
	boxTypeCMT3FourCC = fourCCFromString("CMT3")
	boxTypeCMT4FourCC = fourCCFromString("CMT4")
	boxTypeCNCVFourCC = fourCCFromString("CNCV")
	boxTypeCo64FourCC = fourCCFromString("co64")
	boxTypeColrFourCC = fourCCFromString("colr")
	boxTypeCRAWFourCC = fourCCFromString("CRAW")
	boxTypeCrttFourCC = fourCCFromString("crtt")
	boxTypeCTBOFourCC = fourCCFromString("CTBO")
	boxTypeCTMDFourCC = fourCCFromString("CTMD")
	boxTypeDimgFourCC = fourCCFromString("dimg")
	boxTypeDinfFourCC = fourCCFromString("dinf")
	boxTypeDrefFourCC = fourCCFromString("dref")
	boxTypeEtypFourCC = fourCCFromString("etyp")
	boxTypeFreeFourCC = fourCCFromString("free")
	boxTypeFtypFourCC = fourCCFromString("ftyp")
	boxTypeGrplFourCC = fourCCFromString("grpl")
	boxTypeHdlrFourCC = fourCCFromString("hdlr")
	boxTypeHvcCFourCC = fourCCFromString("hvcC")
	boxTypeIdatFourCC = fourCCFromString("idat")
	boxTypeIinfFourCC = fourCCFromString("iinf")
	boxTypeIlocFourCC = fourCCFromString("iloc")
	boxTypeImirFourCC = fourCCFromString("imir")
	boxTypeInfeFourCC = fourCCFromString("infe")
	boxTypeIovlFourCC = fourCCFromString("iovl")
	boxTypeIpcoFourCC = fourCCFromString("ipco")
	boxTypeIpmaFourCC = fourCCFromString("ipma")
	boxTypeIprpFourCC = fourCCFromString("iprp")
	boxTypeIrefFourCC = fourCCFromString("iref")
	boxTypeIrotFourCC = fourCCFromString("irot")
	boxTypeIspeFourCC = fourCCFromString("ispe")
	boxTypeJXLFourCC  = fourCCFromString("JXL ")
	boxTypeJumbFourCC = fourCCFromString("jumb")
	boxTypeJxlcFourCC = fourCCFromString("jxlc")
	boxTypeJxllFourCC = fourCCFromString("jxll")
	boxTypeJxlpFourCC = fourCCFromString("jxlp")
	boxTypeLhvCFourCC = fourCCFromString("lhvC")
	boxTypeMdatFourCC = fourCCFromString("mdat")
	boxTypeMdftFourCC = fourCCFromString("mdft")
	boxTypeMdhdFourCC = fourCCFromString("mdhd")
	boxTypeMdiaFourCC = fourCCFromString("mdia")
	boxTypeMetaFourCC = fourCCFromString("meta")
	boxTypeMinfFourCC = fourCCFromString("minf")
	boxTypeMoovFourCC = fourCCFromString("moov")
	boxTypeMvhdFourCC = fourCCFromString("mvhd")
	boxTypeNmhdFourCC = fourCCFromString("nmhd")
	boxTypeOinfFourCC = fourCCFromString("oinf")
	boxTypePaspFourCC = fourCCFromString("pasp")
	boxTypePitmFourCC = fourCCFromString("pitm")
	boxTypePixiFourCC = fourCCFromString("pixi")
	boxTypePRVWFourCC = fourCCFromString("PRVW")
	boxTypeStblFourCC = fourCCFromString("stbl")
	boxTypeStscFourCC = fourCCFromString("stsc")
	boxTypeStsdFourCC = fourCCFromString("stsd")
	boxTypeStszFourCC = fourCCFromString("stsz")
	boxTypeSttsFourCC = fourCCFromString("stts")
	boxTypeThmbFourCC = fourCCFromString("thmb")
	boxTypeTHMBFourCC = fourCCFromString("THMB")
	boxTypeTkhdFourCC = fourCCFromString("tkhd")
	boxTypeTolsFourCC = fourCCFromString("tols")
	boxTypeTrakFourCC = fourCCFromString("trak")
	boxTypeUUIDFourCC = fourCCFromString("uuid")
	boxTypeVmhdFourCC = fourCCFromString("vmhd")
	boxTypeExifFourCC = fourCCFromString("Exif")
)
