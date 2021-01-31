package xmlname

import "fmt"

// Space is the Namespace part of an XMP Namespace
type Space uint8

func (s Space) String() string {
	return fmt.Sprintf(mapSpaceString[s])
}

// IdentifySpace returns the (Space) XML Namespace correspondent to buf.
// If NS was not identified returns UnknownNS.
func IdentifySpace(buf []byte) (n Space) {
	return mapStringSpace[string(buf)]
}

// XML Namespaces supported
const (
	UnknownNS Space = iota
	Aux             // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	Crs             // xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
	Darktable       // xmlns:darktable="http://darktable.sf.net/
	Dc              // xmlns:dc="http://purl.org/dc/elements/1.1/"
	Exif            // xmlns:exif="http://ns.adobe.com/exif/1.0/"
	ExifEX          // xmlns:exifEX="http://cipa.jp/exif/1.0/"
	Lr              // xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
	Photoshop       // xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
	Pmi             // xmlns:pmi='http://prismstandard.org/namespaces/pmi/2.2/'
	Rdf             // xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	StEvt           // xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
	StRef           // xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
	Tiff            // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	X               // xmlns:x="adobe:ns:meta/"
	XML
	XMLns
	Xmp   // xmlns:xmp="http://ns.adobe.com/xap/1.0/"
	XmpMM // xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
)

var mapStringSpace = map[string]Space{
	"Unknown":   UnknownNS,
	"aux":       Aux,
	"crs":       Crs,
	"darktable": Darktable,
	"dc":        Dc,
	"exif":      Exif,
	"exifEX":    ExifEX,
	"lr":        Lr,
	"photoshop": Photoshop,
	"pmi":       Pmi,
	"rdf":       Rdf,
	"stEvt":     StEvt,
	"stRef":     StRef,
	"tiff":      Tiff,
	"x":         X,
	"xml":       XML,
	"xmlns":     XMLns,
	"xmp":       Xmp,
	"xmpMM":     XmpMM,
}

var mapSpaceString = map[Space]string{
	UnknownNS: "Unknown",
	Aux:       "aux",
	Crs:       "crs",
	Darktable: "darktable",
	Dc:        "dc",
	Exif:      "exif",
	ExifEX:    "exifEX",
	Lr:        "lr",
	Photoshop: "photoshop",
	Pmi:       "pmi",
	Rdf:       "rdf",
	StEvt:     "stEvt",
	StRef:     "stRef",
	Tiff:      "tiff",
	X:         "x",
	XML:       "xml",
	XMLns:     "xmlns",
	Xmp:       "xmp",
	XmpMM:     "xmpMM",
}
