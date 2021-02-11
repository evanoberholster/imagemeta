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
	//"MA1B": brandMA1B,
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
	//brandMA1B: "MA1B",
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
	//*box
	MinorVersion string   // 4 bytes
	MajorBrand   Brand    // 4 bytes
	Compatible   [6]Brand // all 4 bytes
}

func (ftyp FileTypeBox) Type() BoxType {
	return TypeFtyp
}

func parseFileTypeBox(outer *box) (b Box, err error) {
	var buf []byte
	if buf, err = outer.r.Peek(8); err != nil {
		return nil, err
	}
	ft := FileTypeBox{
		//box:          outer,
		MajorBrand:   brand(buf[:4]),
		MinorVersion: string(buf[4:8]),
	}
	if err = outer.r.discard(8); err != nil {
		return
	}
	for i := 0; i < 6; i++ {
		if outer.r.remain < 4 {
			break
		}

		ft.Compatible[i], err = outer.r.readBrand()
		if err != nil {
			break
		}
	}
	return ft, outer.r.discard(int(outer.r.remain))
}

func (ftb FileTypeBox) String() string {
	str := fmt.Sprintf("(Box) ftyp | Major Brand: %s, Minor Version: %s, Compatible: ", ftb.MajorBrand, "")
	for _, b := range ftb.Compatible {
		str += b.String() + " "
	}
	return str
}
