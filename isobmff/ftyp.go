package isobmff

import (
	"github.com/pkg/errors"
)

// brandCount is the number of compatible brands supported.
const maxBrandCount = 8

// ReadFTYP reads an 'ftyp' box from a BMFF file.
//
// This should be the first read function called.
func (r *Reader) ReadFTYP() (err error) {
	b, err := r.readBox()
	if err != nil {
		return errors.Wrapf(err, "ReadFTYPBox")
	}
	r.ftyp, err = parseFileTypeBox(&b)
	return err
}

func parseFileTypeBox(b *box) (ftyp FileTypeBox, err error) {
	if !b.isType(typeFtyp) {
		return ftyp, ErrWrongBoxType
	}
	buf, err := b.Peek(b.remain)
	if err != nil {
		return ftyp, err
	}
	ftyp.MajorBrand = brandFromBuf(buf[:4])
	copy(ftyp.MinorVersion[:4], buf[4:8])

	// Read maximum 7 Compatible brands
	for i, compatibleBrand := 8, 0; i < b.remain && compatibleBrand < maxBrandCount; compatibleBrand++ {
		ftyp.Compatible[compatibleBrand] = brandFromBuf(buf[i : i+4])
		i += 4
	}
	if logLevelInfo() {
		logInfoBox(b).Str("MajorBrand", ftyp.MajorBrand.String()).Str("MinorVersion", string(ftyp.MinorVersion[:])).Strs("MinorBrands", minorBrandsToString(ftyp)).Send()
	}
	return ftyp, b.Discard(b.remain)
}

// FileTypeBox is a BMFF FileTypeBox
type FileTypeBox struct {
	Compatible   [maxBrandCount]Brand // all 4 bytes
	MinorVersion [4]byte              // 4 bytes
	MajorBrand   Brand                // 4 bytes
}

// Brand of ISOBMFF ftyp
type Brand uint8

// String is the Stringer interface for Brand
func (b Brand) String() string {
	if str, found := mapBrandString[b]; found {
		return str
	}
	return "nnnn"
}
func brandFromBuf(buf []byte) Brand {
	if len(buf) == 4 {
		if b, ok := mapStringBrand[string(buf)]; ok {
			return b
		}
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
	brandAvci                 // 'avci'
	brandAvif                 // 'avif': AVIF
	brandCrx                  // 'crx ' : Canon CR3
	brandHeic                 // 'heic': the usual HEIF images
	brandHeim                 // 'heim': multiview
	brandHeis                 // 'heis': scalable
	brandHeix                 // 'heix': 10bit images, or anything that uses h265 with range extension
	brandHevc                 // 'hevc': brand for image sequences
	brandHevm                 // 'hevm': multiview sequence
	brandHevs                 // 'hevs': scalable sequence
	brandHevx                 // 'hevx': image sequence
	brandIso8                 // 'iso8': sequence
	brandIsom                 // 'isom' : ?
	brandM4A                  // 'M4A '
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
	brandMsf1                 // 'msf1': sequence
)

var (
	mapStringBrand = map[string]Brand{
		"avci": brandAvci,
		"avif": brandAvif,
		"crx ": brandCrx,
		"heic": brandHeic,
		"heim": brandHeim,
		"heis": brandHeis,
		"heix": brandHeix,
		"hevc": brandHevc,
		"hevm": brandHevm,
		"hevs": brandHevs,
		"hevx": brandHevx,
		"iso8": brandIso8,
		"isom": brandIsom,
		"M4A ": brandM4A,
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
		"msf1": brandMsf1,
	}

	mapBrandString = map[Brand]string{
		brandAvci: "avci",
		brandAvif: "avif",
		brandCrx:  "crx ",
		brandHeic: "heic",
		brandHeim: "heim",
		brandHeis: "heis",
		brandHeix: "heix",
		brandHevc: "hevc",
		brandHevm: "hevm",
		brandHevs: "hevs",
		brandHevx: "hevx",
		brandIso8: "iso8",
		brandIsom: "isom",
		brandM4A:  "M4A ",
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
		brandMsf1: "msf1",
	}
)
