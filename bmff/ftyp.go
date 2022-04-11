package bmff

import (
	"fmt"
	"strings"
)

// Brand of ISOBMFF ftyp
type Brand uint8

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

func (b Brand) String() string {
	str, found := mapBrandString[b]
	if found {
		return str
	}
	return "nnnn"
}
func brand(buf []byte) Brand {
	if len(buf) == 4 {
		if b, found := mapStringBrand[string(buf)]; found {
			return b
		}
	}
	if debugFlag {
		log.Debug("Brand '%s' not found", buf)
	}
	return brandUnknown
}

var MinorVersionNul = [4]byte{}

// FileTypeBox is a BMFF FileTypeBox
type FileTypeBox struct {
	Compatible   [brandCount]Brand // all 4 bytes
	MinorVersion [4]byte           // 4 bytes
	MajorBrand   Brand             // 4 bytes
}

// IsCR3 returns true if major brand is crx (Canon CR3)
func (ftyp FileTypeBox) IsCR3() bool {
	return ftyp.MajorBrand == brandCrx
}

func (b *box) parseFileTypeBox() (ftyp FileTypeBox, err error) {
	if b.boxType != TypeFtyp {
		return ftyp, ErrWrongBoxType
	}
	buf, err := b.peek(b.remain)
	if err != nil {
		return ftyp, err
	}
	ftyp.MajorBrand = brand(buf[:4])
	copy(ftyp.MinorVersion[:4], buf[4:8])

	// Read maximum 7 Compatible brands
	for i, compatibleBrand := 8, 0; i < b.remain && compatibleBrand < brandCount; compatibleBrand++ {
		ftyp.Compatible[compatibleBrand] = brand(buf[i : i+4])
		i += 4
	}

	if debugFlag {
		traceBoxWithMsg(*b, ftyp.String())
	}
	return ftyp, b.discard(b.remain)
}

func (ftyp FileTypeBox) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ftyp | Major Brand:'%s', Minor Version:'%s', Compatible: ", ftyp.MajorBrand, ""))
	for _, b := range ftyp.Compatible {
		sb.WriteString("'")
		sb.WriteString(b.String())
		sb.WriteString("' ")
	}
	return sb.String()
}
