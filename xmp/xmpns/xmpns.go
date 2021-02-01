// Package xmpns provides XMP Namespace information
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
	UnknownNS Namespace = iota
	// xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	AuxNS
	// xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
	CrsNS
	// xmlns:darktable="http://darktable.sf.net/"
	DarktableNS
	// xmlns:dc="http://purl.org/dc/elements/1.1/"
	DcNS
	// xmlns:exif="http://ns.adobe.com/exif/1.0/"
	ExifNS
	// xmlns:exifEX="http://cipa.jp/exif/1.0/"
	ExifEXNS
	// xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
	LrNS
	// xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
	PhotoshopNS
	// xmlns:pmi="http://prismstandard.org/namespaces/pmi/2.2/"
	PmiNS
	// xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	RdfNS
	// xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
	StEvtNS
	// xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
	StRefNS
	// xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	TiffNS
	// xmlns:x="adobe:ns:meta/"
	XNS
	// xmlns:xap="http://ns.adobe.com/xap/1.0/"
	XapNS
	// xmlns:xapMM="http://ns.adobe.com/xap/1.0/mm/"
	XapMMNS
	XMLNS
	XMLnsNS
	// xmlns:xmp="http://ns.adobe.com/xap/1.0/"
	XmpNS
	// xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
	XmpMMNS
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
	"xap":       XapNS,
	"xapMM":     XapMMNS,
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
	XapNS:       "xap",
	XapMMNS:     "xapMM",
	XMLNS:       "xml",
	XMLnsNS:     "xmlns",
	XmpNS:       "xmp",
	XmpMMNS:     "xmpMM",
}
