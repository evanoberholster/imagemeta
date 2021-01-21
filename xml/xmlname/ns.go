package xmlname

import "fmt"

// NS is an XML Namespace
type NS uint8

func (ns NS) String() string {
	return fmt.Sprintf(mapNSString[ns])
}

// IdentifyNS returns the (NS) XML Namespace correspondent to buf.
// If NS was not identified returns UnknownNS.
func IdentifyNS(buf []byte) (n NS) {
	return mapStringNS[string(buf)]
}

// XML Namespaces supported
const (
	UnknownNS NS = iota
	Rdf
	Xmlns
	X         // xmlns:x="adobe:ns:meta/"
	Xmp       // xmlns:xmp="http://ns.adobe.com/xap/1.0/"
	XmpMM     // xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
	Tiff      // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	Exif      // xmlns:exif="http://ns.adobe.com/exif/1.0/"
	ExifEX    // xmlns:exifEX="http://cipa.jp/exif/1.0/"
	Dc        // xmlns:dc="http://purl.org/dc/elements/1.1/"
	Aux       // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	Crs       // xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
	Lr        // xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
	StEvt     // xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
	StRef     // xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
	Photoshop // xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
	Darktable // xmlns:darktable="http://darktable.sf.net/
)

var mapStringNS = map[string]NS{
	"Unknown":   UnknownNS,
	"x":         X,
	"xmlns":     Xmlns,
	"xmp":       Xmp,
	"xmpMM":     XmpMM,
	"tiff":      Tiff,
	"exif":      Exif,
	"exifEX":    ExifEX,
	"dc":        Dc,
	"aux":       Aux,
	"photoshop": Photoshop,
	"crs":       Crs,
	"lr":        Lr,
	"rdf":       Rdf,
	"stEvt":     StEvt,
	"stRef":     StRef,
	"darktable": Darktable,
}

var mapNSString = map[NS]string{
	UnknownNS: "Unknown",
	X:         "x",
	Xmlns:     "xmlns",
	Xmp:       "xmp",
	XmpMM:     "xmpMM",
	Tiff:      "tiff",
	Exif:      "exif",
	ExifEX:    "exifEX",
	Dc:        "dc",
	Aux:       "aux",
	Photoshop: "photoshop",
	Crs:       "crs",
	Lr:        "lr",
	Rdf:       "rdf",
	StEvt:     "stEvt",
	StRef:     "stRef",
	Darktable: "darktable",
}
