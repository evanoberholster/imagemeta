package xmp

// Namespace is the namespace-prefix portion of an XMP property (for example "exif").
type Namespace uint8

// NameSpace is retained as a compatibility alias for Namespace.
type NameSpace = Namespace

func (ns Namespace) String() string {
	if int(ns) < len(mapNSString) {
		return mapNSString[ns]
	}
	return mapNSString[UnknownNS]
}

// IdentifyNamespace resolves an XML prefix token to its internal namespace identifier.
// Unknown prefixes resolve to UnknownNS.
func IdentifyNamespace(buf []byte) (n Namespace) {
	if n = identifyNamespaceFast(buf); n != UnknownNS {
		return
	}
	for i := 1; i < len(mapNSString); i++ {
		s := mapNSString[i]
		if s != "" && eqString(buf, s) {
			return Namespace(i)
		}
	}
	return UnknownNS
}

// IdentifyNameSpace is retained as a compatibility wrapper around IdentifyNamespace.
func IdentifyNameSpace(buf []byte) Namespace {
	return IdentifyNamespace(buf)
}

// XML Namespaces supported
const (
	UnknownNS Namespace = iota
	AppleFiNS
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
	// xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/"
	MwgRSNS
	// xmlns:pmi="http://prismstandard.org/namespaces/pmi/2.2/"
	PmiNS
	// xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	RdfNS
	StAreaNS
	StDimNS
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
	XmpDMNS
	// xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
	XmpMMNS
)

var mapNSString = [...]string{
	UnknownNS:   "Unknown",
	AppleFiNS:   "apple-fi",
	AuxNS:       "aux",
	CrsNS:       "crs",
	DarktableNS: "darktable",
	DcNS:        "dc",
	ExifNS:      "exif",
	ExifEXNS:    "exifEX",
	LrNS:        "lr",
	PhotoshopNS: "photoshop",
	MwgRSNS:     "mwg-rs",
	PmiNS:       "pmi",
	RdfNS:       "rdf",
	StAreaNS:    "stArea",
	StDimNS:     "stDim",
	StEvtNS:     "stEvt",
	StRefNS:     "stRef",
	TiffNS:      "tiff",
	XNS:         "x",
	XapNS:       "xap",
	XapMMNS:     "xapMM",
	XMLNS:       "xml",
	XMLnsNS:     "xmlns",
	XmpNS:       "xmp",
	XmpDMNS:     "xmpDM",
	XmpMMNS:     "xmpMM",
}

func identifyNamespaceFast(buf []byte) Namespace {
	switch string(buf) {
	case "x":
		return XNS
	case "dc":
		return DcNS
	case "lr":
		return LrNS
	case "aux":
		return AuxNS
	case "crs":
		return CrsNS
	case "rdf":
		return RdfNS
	case "xml":
		return XMLNS
	case "xap":
		return XapNS
	case "xmp":
		return XmpNS
	case "tiff":
		return TiffNS
	case "exif":
		return ExifNS
	case "xmlns":
		return XMLnsNS
	case "stArea":
		return StAreaNS
	case "stDim":
		return StDimNS
	case "stEvt":
		return StEvtNS
	case "stRef":
		return StRefNS
	case "xapMM":
		return XapMMNS
	case "xmpDM":
		return XmpDMNS
	case "xmpMM":
		return XmpMMNS
	case "exifEX":
		return ExifEXNS
	case "photoshop":
		return PhotoshopNS
	case "mwg-rs":
		return MwgRSNS
	case "apple-fi":
		return AppleFiNS
	}
	return UnknownNS
}

func eqString(buf []byte, s string) bool {
	if len(buf) != len(s) {
		return false
	}
	for i := 0; i < len(buf); i++ {
		if buf[i] != s[i] {
			return false
		}
	}
	return true
}
