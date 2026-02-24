package isobmff

import (
	"fmt"
)

// brandCount is the number of compatible brands supported.
const maxBrandCount = 8

// ReadFTYP reads an 'ftyp' box from a BMFF file.
//
// This should be the first read function called.
func (r *Reader) ReadFTYP() (err error) {
	b, err := r.readBox()
	if err != nil {
		return fmt.Errorf("ReadFTYPBox: %w", err)
	}

	// JPEG XL containers may start with a signature box ('JXL ') before 'ftyp'.
	for b.isType(typeJXL) {
		if err = b.close(); err != nil {
			return fmt.Errorf("ReadFTYPBox: failed to skip JXL signature box: %w", err)
		}
		b, err = r.readBox()
		if err != nil {
			return fmt.Errorf("ReadFTYPBox: %w", err)
		}
	}

	r.ftyp, err = parseFileTypeBox(&b)
	return err
}

func parseFileTypeBox(b *box) (ftyp FileTypeBox, err error) {
	if !b.isType(typeFtyp) {
		return ftyp, ErrWrongBoxType
	}

	if b.remain < 8 {
		return ftyp, fmt.Errorf("parseFileTypeBox: %w", ErrBufLength)
	}

	peekLen := b.remain
	maxFTYPPeek := 8 + (4 * maxBrandCount)
	if peekLen > maxFTYPPeek {
		peekLen = maxFTYPPeek
	}

	buf, err := b.Peek(peekLen)
	if err != nil {
		return ftyp, err
	}

	ftyp.MajorBrand = brandFromBuf(buf[:4])
	copy(ftyp.MinorVersion[:4], buf[4:8])

	for i, compatibleBrand := 8, 0; i+4 <= len(buf) && compatibleBrand < maxBrandCount; compatibleBrand++ {
		ftyp.Compatible[compatibleBrand] = brandFromBuf(buf[i : i+4])
		i += 4
	}
	if logLevelInfo() {
		logInfoBox(b).Str("MajorBrand", ftyp.MajorBrand.String()).Str("MinorVersion", string(ftyp.MinorVersion[:])).Strs("MinorBrands", minorBrandsToString(ftyp)).Send()
	}
	return ftyp, b.close()
}

// FileTypeBox is a BMFF FileTypeBox
type FileTypeBox struct {
	Compatible   [maxBrandCount]Brand // all 4 bytes
	MinorVersion [4]byte              // 4 bytes
	MajorBrand   Brand                // 4 bytes
}

// Brand of ISOBMFF ftyp
type Brand uint32

// String is the Stringer interface for Brand
func (b Brand) String() string {
	if str, found := mapBrandString[b]; found {
		return str
	}
	return "nnnn"
}

func brandFromBuf(buf []byte) Brand {
	if len(buf) < 4 {
		return brandUnknown
	}

	if b, ok := mapFourCCBrand[bmffEndian.Uint32(buf[:4])]; ok {
		return b
	}

	if logLevelError() {
		logErrorMsg("Brand", "error Brand '%s' unknown", buf)
	}
	return brandUnknown
}

func minorBrandsToString(ftyp FileTypeBox) []string {
	brands := make([]string, maxBrandCount)
	j := 0
	for _, b := range ftyp.Compatible {
		if b != brandUnknown {
			brands[j] = b.String()
			j++
		}
	}
	return brands[:j]
	//return nil
}

// Major and Minor Brands
const (
	brandUnknown Brand = iota // unknown ISOBMFF brand
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
	mapStringBrand = map[string]Brand{
		"3g2a": brand3G2A,
		"3g2b": brand3G2B,
		"3g2c": brand3G2C,
		"3gp4": brand3GP4,
		"3gp5": brand3GP5,
		"3gp6": brand3GP6,
		"3gp7": brand3GP7,
		"avci": brandAvci,
		"avif": brandAvif,
		"avis": brandAvis,
		"crx ": brandCrx,
		"dash": brandDash,
		"heic": brandHeic,
		"heif": brandHeif,
		"heim": brandHeim,
		"heis": brandHeis,
		"heix": brandHeix,
		"hevc": brandHevc,
		"hevm": brandHevm,
		"hevs": brandHevs,
		"hevx": brandHevx,
		"iso2": brandIso2,
		"iso3": brandIso3,
		"iso4": brandIso4,
		"iso5": brandIso5,
		"iso6": brandIso6,
		"iso8": brandIso8,
		"isom": brandIsom,
		"jxl ": brandJxl,
		"M4A ": brandM4A,
		"M4V ": brandM4V,
		"M4VH": brandM4VH,
		"M4VP": brandM4VP,
		"MA1B": brandMA1B,
		"meta": brandMeta,
		"miaf": brandMiaf,
		"MiAn": brandMiAn,
		"MiBr": brandMiBr,
		"mif1": brandMif1,
		"mif2": brandMif2,
		"MiHA": brandMiHA,
		"MiHB": brandMiHB,
		"MiHE": brandMiHE,
		"MiPr": brandMiPr,
		"mp41": brandMp41,
		"mp42": brandMp42,
		"mp71": brandMp71,
		"msf1": brandMsf1,
		"qt  ": brandQt,
	}

	mapFourCCBrand = func() map[uint32]Brand {
		m := make(map[uint32]Brand, len(mapStringBrand))
		for k, v := range mapStringBrand {
			if len(k) == 4 {
				m[fourCCFromString(k)] = v
			}
		}
		return m
	}()

	mapBrandString = map[Brand]string{
		brand3G2A: "3g2a",
		brand3G2B: "3g2b",
		brand3G2C: "3g2c",
		brand3GP4: "3gp4",
		brand3GP5: "3gp5",
		brand3GP6: "3gp6",
		brand3GP7: "3gp7",
		brandAvci: "avci",
		brandAvif: "avif",
		brandAvis: "avis",
		brandCrx:  "crx ",
		brandDash: "dash",
		brandHeic: "heic",
		brandHeif: "heif",
		brandHeim: "heim",
		brandHeis: "heis",
		brandHeix: "heix",
		brandHevc: "hevc",
		brandHevm: "hevm",
		brandHevs: "hevs",
		brandHevx: "hevx",
		brandIso2: "iso2",
		brandIso3: "iso3",
		brandIso4: "iso4",
		brandIso5: "iso5",
		brandIso6: "iso6",
		brandIso8: "iso8",
		brandIsom: "isom",
		brandJxl:  "jxl ",
		brandM4A:  "M4A ",
		brandM4V:  "M4V ",
		brandM4VH: "M4VH",
		brandM4VP: "M4VP",
		brandMA1B: "MA1B",
		brandMeta: "meta",
		brandMiaf: "miaf",
		brandMiAn: "MiAn",
		brandMiBr: "MiBr",
		brandMif1: "mif1",
		brandMif2: "mif2",
		brandMiHA: "MiHA",
		brandMiHB: "MiHB",
		brandMiHE: "MiHE",
		brandMiPr: "MiPr",
		brandMp41: "mp41",
		brandMp42: "mp42",
		brandMp71: "mp71",
		brandMsf1: "msf1",
		brandQt:   "qt  ",
	}
)
