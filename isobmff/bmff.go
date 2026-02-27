package isobmff

import (
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
)

// ExifReader is invoked with Exif payload bytes and its parsed header.
type ExifReader func(r io.Reader, h meta.ExifHeader) error

// XMPReader is invoked with XMP payload bytes and XPacket metadata.
type XMPReader func(r io.Reader, h XPacketHeader) error

// PreviewImageReader is invoked with preview image bytes and metadata header.
type PreviewImageReader func(r io.Reader, h meta.PreviewHeader) error

// XPacketHeader describes discovered XPacket payload metadata.
type XPacketHeader struct {
	Offset       uint64
	Length       uint32
	HasXPacketPI bool
	HasXMPMeta   bool
}

// Errors
var (
	ErrBufLength                = errors.New("insufficient buffer length")
	ErrRemainLengthInsufficient = errors.New("remain length insufficient")
	ErrUnsupportedFieldSize     = errors.New("unsupported field size")
	ErrBoxStringTooLong         = errors.New("box string too long")
	errLargeBox                 = errors.New("unexpectedly large box")
	ErrWrongBoxType             = errors.New("wrong box type")
)

const (
	fourCCSize     = 4
	ftypHeaderSize = 8

	// maxBrandCount is the maximum number of compatible brands retained.
	maxBrandCount = 8
)

// ReadFTYP reads an 'ftyp' box from a BMFF file.
//
// This should be the first read function called.
func (r *Reader) ReadFTYP() (err error) {
	b, err := r.readBox()
	if err != nil {
		return fmt.Errorf("ReadFTYP: %w", err)
	}

	// JPEG XL containers may start with a signature box ('JXL ') before 'ftyp'.
	for b.isType(typeJXL) {
		if err = b.close(); err != nil {
			return fmt.Errorf("ReadFTYP: failed to skip JXL signature box: %w", err)
		}
		b, err = r.readBox()
		if err != nil {
			return fmt.Errorf("ReadFTYP: %w", err)
		}
	}

	r.ftyp, err = parseFileTypeBox(&b)
	if err != nil {
		return err
	}
	r.initMetadataGoals()
	return nil
}

func parseFileTypeBox(b *box) (ftyp fileTypeBox, err error) {
	if !b.isType(typeFtyp) {
		return ftyp, ErrWrongBoxType
	}

	if b.remain < ftypHeaderSize {
		return ftyp, fmt.Errorf("parseFileTypeBox: %w", ErrBufLength)
	}

	peekLen := int64(ftypHeaderSize + (fourCCSize * maxBrandCount))
	if b.remain < peekLen {
		peekLen = b.remain
	}

	buf, err := b.Peek(int(peekLen))
	if err != nil {
		return ftyp, err
	}

	ftyp.MajorBrand = brandFromBuf(buf[:4])
	copy(ftyp.MinorVersion[:], buf[4:8])

	for i, compatibleBrand := ftypHeaderSize, 0; i+fourCCSize <= len(buf) && compatibleBrand < maxBrandCount; compatibleBrand++ {
		ftyp.Compatible[compatibleBrand] = brandFromBuf(buf[i : i+fourCCSize])
		i += fourCCSize
	}
	if logLevelInfo() {
		logInfoBox(b).Str("MajorBrand", ftyp.MajorBrand.String()).Str("MinorVersion", string(ftyp.MinorVersion[:])).Strs("MinorBrands", minorBrandsToString(ftyp)).Send()
	}
	return ftyp, b.close()
}

// fileTypeBox is a BMFF fileTypeBox
type fileTypeBox struct {
	Compatible   [maxBrandCount]brand // all 4 bytes
	MinorVersion [fourCCSize]byte     // 4 bytes
	MajorBrand   brand                // 4 bytes
}

// brand of ISOBMFF ftyp
type brand uint32

// String is the Stringer interface for brand
func (b brand) String() string {
	i := int(b)
	if i >= 0 && i < len(brandCodes) {
		if code := brandCodes[i]; len(code) == fourCCSize {
			return code
		}
	}
	return "nnnn"
}

func brandFromBuf(buf []byte) brand {
	if len(buf) < 4 {
		return brandUnknown
	}

	if b, ok := mapFourCCBrand[bmffEndian.Uint32(buf[:4])]; ok {
		return b
	}

	if logLevelDebug() {
		logDebug().Str("brand", string(buf[:4])).Msg("unknown brand")
	}
	return brandUnknown
}

func minorBrandsToString(ftyp fileTypeBox) []string {
	brands := make([]string, 0, maxBrandCount)
	for _, b := range ftyp.Compatible {
		if b != brandUnknown {
			brands = append(brands, b.String())
		}
	}
	return brands
}

// Major and Minor Brands
const (
	brandUnknown brand = iota // unknown ISOBMFF brand
	brand3G2A                 // '3g2a': 3GPP2
	brand3G2B                 // '3g2b': 3GPP2
	brand3G2C                 // '3g2c': 3GPP2
	brand3GP4                 // '3gp4': 3GPP Release 4
	brand3GP5                 // '3gp5': 3GPP Release 5
	brand3GP6                 // '3gp6': 3GPP Release 6
	brand3GP7                 // '3gp7': 3GPP Release 7
	brandAvci                 // 'avci'
	brandAvif                 // 'avif': AVIF
	brandAvis                 // 'avis': AVIF image sequence
	brandCrx                  // 'crx ' : Canon CR3
	brandDash                 // 'dash': MPEG-DASH
	brandHeic                 // 'heic': the usual HEIF images
	brandHeif                 // 'heif': generic HEIF
	brandHeim                 // 'heim': multiview
	brandHeis                 // 'heis': scalable
	brandHeix                 // 'heix': 10bit images, or anything that uses h265 with range extension
	brandHevc                 // 'hevc': brand for image sequences
	brandHevm                 // 'hevm': multiview sequence
	brandHevs                 // 'hevs': scalable sequence
	brandHevx                 // 'hevx': image sequence
	brandIso2                 // 'iso2': MP4 version 2
	brandIso3                 // 'iso3': MP4 version 3
	brandIso4                 // 'iso4': MP4 version 4
	brandIso5                 // 'iso5': MP4 version 5
	brandIso6                 // 'iso6': MP4 version 6
	brandIso8                 // 'iso8': sequence
	brandIsom                 // 'isom' : ?
	brandJxl                  // 'jxl ': JPEG XL
	brandM4A                  // 'M4A '
	brandM4V                  // 'M4V ': MPEG-4 video
	brandM4VH                 // 'M4VH': MPEG-4 video
	brandM4VP                 // 'M4VP': MPEG-4 video
	brandMA1B                 // 'MA1B'
	brandMeta                 // 'meta': meta
	brandMiaf                 // 'miaf' :
	brandMiAn                 // 'MiAn'
	brandMiBr                 // 'MiBr'
	brandMif1                 // 'mif1': image
	brandMif2                 // 'mif2'
	brandMiHA                 // 'MiHA'
	brandMiHB                 // 'MiHB' :
	brandMiHE                 // 'MiHE' :
	brandMiPr                 // 'MiPr'
	brandMp41                 // 'mp41'
	brandMp42                 // 'mp42'
	brandMp71                 // 'mp71'
	brandMsf1                 // 'msf1': sequence
	brandQt                   // 'qt  ': QuickTime
)

var (
	brandCodes = [...]string{
		brandUnknown: "nnnn",
		brand3G2A:    "3g2a",
		brand3G2B:    "3g2b",
		brand3G2C:    "3g2c",
		brand3GP4:    "3gp4",
		brand3GP5:    "3gp5",
		brand3GP6:    "3gp6",
		brand3GP7:    "3gp7",
		brandAvci:    "avci",
		brandAvif:    "avif",
		brandAvis:    "avis",
		brandCrx:     "crx ",
		brandDash:    "dash",
		brandHeic:    "heic",
		brandHeif:    "heif",
		brandHeim:    "heim",
		brandHeis:    "heis",
		brandHeix:    "heix",
		brandHevc:    "hevc",
		brandHevm:    "hevm",
		brandHevs:    "hevs",
		brandHevx:    "hevx",
		brandIso2:    "iso2",
		brandIso3:    "iso3",
		brandIso4:    "iso4",
		brandIso5:    "iso5",
		brandIso6:    "iso6",
		brandIso8:    "iso8",
		brandIsom:    "isom",
		brandJxl:     "jxl ",
		brandM4A:     "M4A ",
		brandM4V:     "M4V ",
		brandM4VH:    "M4VH",
		brandM4VP:    "M4VP",
		brandMA1B:    "MA1B",
		brandMeta:    "meta",
		brandMiaf:    "miaf",
		brandMiAn:    "MiAn",
		brandMiBr:    "MiBr",
		brandMif1:    "mif1",
		brandMif2:    "mif2",
		brandMiHA:    "MiHA",
		brandMiHB:    "MiHB",
		brandMiHE:    "MiHE",
		brandMiPr:    "MiPr",
		brandMp41:    "mp41",
		brandMp42:    "mp42",
		brandMp71:    "mp71",
		brandMsf1:    "msf1",
		brandQt:      "qt  ",
	}

	mapFourCCBrand = func() map[uint32]brand {
		m := make(map[uint32]brand, len(brandCodes))
		for i, code := range brandCodes {
			if len(code) != fourCCSize || code == "nnnn" {
				continue
			}
			m[fourCCFromString(code)] = brand(i)
		}
		return m
	}()
)
