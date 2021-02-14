package bmff

import (
	"fmt"
)

// Brand of ISOBMFF ftyp
type Brand uint8

// Major and Minor Brands
const (
	brandUnknown Brand = iota // unknown ISOBMFF brand
	brandAvif                 // 'avif': AVIF
	brandHeic                 // 'heic': the usual HEIF images
	brandHeim                 // 'heim': multiview
	brandHeis                 // 'heis': scalable
	brandHeix                 // 'heix': 10bit images, or anything that uses h265 with range extension
	brandHevc                 // 'hevc': brand for image sequences
	brandHevm                 // 'hevm': multiview sequence
	brandHevs                 // 'hevs': scalable sequence
	brandHevx                 // 'hevx': image sequence
	brandMeta                 // 'meta': meta
	//brandMA1B                 // 'MA1B': ?
	brandCrx  // 'crx ' : Canon CR3
	brandIsom // 'isom' : ?
	brandMiaf // 'miaf' :
	brandMif1 // 'mif1': image
	brandMiHB // 'MiHB' :
	brandMiHE // 'MiHE' :
	brandMsf1 // 'msf1': sequence
)

var mapStringBrand = map[string]Brand{
	"avif": brandAvif,
	"heic": brandHeic,
	"heim": brandHeim,
	"heis": brandHeis,
	"heix": brandHeix,
	"hevc": brandHevc,
	"hevm": brandHevm,
	"hevs": brandHevs,
	"hevx": brandHevx,
	"meta": brandMeta,
	"miaf": brandMiaf,
	"mif1": brandMif1,
	"MiHB": brandMiHB,
	"MiHE": brandMiHE,
	"msf1": brandMsf1,
	"crx ": brandCrx,
	"isom": brandIsom,
}

var mapBrandString = map[Brand]string{
	brandAvif: "avif",
	brandHeic: "heic",
	brandHeim: "heim",
	brandHeis: "heis",
	brandHeix: "heix",
	brandHevc: "hevc",
	brandHevm: "hevm",
	brandHevs: "hevs",
	brandHevx: "hevx",
	brandMeta: "meta",
	brandMiaf: "miaf",
	brandMif1: "mif1",
	brandMiHB: "MiHB",
	brandMiHE: "MiHE",
	brandMsf1: "msf1",
	brandCrx:  "crx ",
	brandIsom: "isom",
}

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
	if Debug {
		fmt.Println("Unknown Brand: ", string(buf), buf)
	}
	return brandUnknown
}

// FileTypeBox is a BMFF FileTypeBox
type FileTypeBox struct {
	MinorVersion string   // 4 bytes
	MajorBrand   Brand    // 4 bytes
	Compatible   [6]Brand // all 4 bytes
}

// Type returns TypeFtyp
func (ftyp FileTypeBox) Type() BoxType {
	return TypeFtyp
}

func parseFtyp(outer *box) (Box, error) {
	return parseFileTypeBox(outer)
}

func parseFileTypeBox(outer *box) (ftyp FileTypeBox, err error) {
	var buf []byte
	if buf, err = outer.Peek(8); err != nil {
		return
	}
	ftyp.MajorBrand = brand(buf[:4])
	ftyp.MinorVersion = processString(buf[4:8])
	if err = outer.discard(8); err != nil {
		return
	}
	// Read maximum 6 Compatible brands
	for i := 0; i < 6; i++ {
		if outer.remain < 4 {
			break
		}

		ftyp.Compatible[i], _ = outer.readBrand()
	}
	err = outer.discard(outer.remain)
	return ftyp, err
}

func (ftyp FileTypeBox) String() string {
	str := fmt.Sprintf("(Box) ftyp | Major Brand: %s, Minor Version: %s, Compatible: ", ftyp.MajorBrand, "")
	for _, b := range ftyp.Compatible {
		str += b.String() + " "
	}
	return str
}

func processString(buf []byte) string {
	if buf[0] == 0 && buf[1] == 0 && buf[2] == 0 && buf[3] == 0 {
		return ""
	}
	return string(buf)
}
