package xmpns

import "fmt"

// Namespace is the Namespace portion of an XMP Property
type Namespace uint8

func (ns Namespace) String() string {
	return fmt.Sprintf(mapNSString[ns])
}

// IdentifyNamespace returns the (Namespace) XML Namespace correspondent to buf.
// If NS was not identified returns UnknownNS.
func IdentifyNamespace(buf []byte) (n Namespace) {
	return mapStringNS[string(buf)]
}

// XML Namespaces supported
const (
	UnknownNS   Namespace = iota
	AuxNS                 // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	CrsNS                 // xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
	DarktableNS           // xmlns:darktable="http://darktable.sf.net/
	DcNS                  // xmlns:dc="http://purl.org/dc/elements/1.1/"
	ExifNS                // xmlns:exif="http://ns.adobe.com/exif/1.0/"
	ExifEXNS              // xmlns:exifEX="http://cipa.jp/exif/1.0/"
	LrNS                  // xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
	PhotoshopNS           // xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
	PmiNS                 // xmlns:pmi='http://prismstandard.org/namespaces/pmi/2.2/'
	RdfNS                 // xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	StEvtNS               // xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
	StRefNS               // xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
	TiffNS                // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	XNS                   // xmlns:x="adobe:ns:meta/"
	XMLNS
	XMLnsNS
	XmpNS   // xmlns:xmp="http://ns.adobe.com/xap/1.0/"
	XmpMMNS // xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
)

var mapStringNS = map[string]Namespace{
	"Unknown":   UnknownNS,
	"aux":       AuxNS,
	"crs":       CrsNS,
	"darktable": DarktableNS,
	"dc":        DcNS,
	"exif":      ExifNS,
	"exifEX":    ExifEXNS,
	"lr":        LrNS,
	"photoshop": PhotoshopNS,
	"pmi":       PmiNS,
	"rdf":       RdfNS,
	"stEvt":     StEvtNS,
	"stRef":     StRefNS,
	"tiff":      TiffNS,
	"x":         XNS,
	"xml":       XMLNS,
	"xmlns":     XMLnsNS,
	"xmp":       XmpNS,
	"xmpMM":     XmpMMNS,
}

var mapNSString = map[Namespace]string{
	UnknownNS:   "Unknown",
	AuxNS:       "aux",
	CrsNS:       "crs",
	DarktableNS: "darktable",
	DcNS:        "dc",
	ExifNS:      "exif",
	ExifEXNS:    "exifEX",
	LrNS:        "lr",
	PhotoshopNS: "photoshop",
	PmiNS:       "pmi",
	RdfNS:       "rdf",
	StEvtNS:     "stEvt",
	StRefNS:     "stRef",
	TiffNS:      "tiff",
	XNS:         "x",
	XMLNS:       "xml",
	XMLnsNS:     "xmlns",
	XmpNS:       "xmp",
	XmpMMNS:     "xmpMM",
}
