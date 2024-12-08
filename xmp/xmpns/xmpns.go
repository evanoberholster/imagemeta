// Package xmpns provides XMP Namespace information
// Based on https://github.com/exiftool/exiftool
package xmpns

// Namespace represents the namespace of an XMP property.
type Namespace uint8

// Namespace constants.
const (
	UnknownNamespace Namespace = iota
	AuxNamespace
	AlbumNamespace
	CcNamespace
	CrdNamespace
	CrsNamespace
	CrssNamespace
	DcNamespace
	ExifNamespace
	ExifEXNamespace
	IXNamespace
	PdfNamespace
	PdfxNamespace
	PhotoshopNamespace
	RdfNamespace
	RdfsNamespace
	StDimNamespace
	StEvtNamespace
	StFntNamespace
	StJobNamespace
	StRefNamespace
	StVerNamespace
	StMfsNamespace
	StCameraNamespace
	CrlcpNamespace
	TiffNamespace
	XNamespace
	XmpGNamespace
	XmpGImgNamespace
	XmpNamespace
	XmpBJNamespace
	XmpDMNamespace
	XmpMMNamespace
	XmpRightsNamespace
	XmpNoteNamespace
	XmpTPgNamespace
	XmpidqNamespace
	XmpPLUSNamespace
	PanoramaNamespace
	DexNamespace
	MediaproNamespace
	ExpressionmediaNamespace
	Iptc4xmpCoreNamespace
	Iptc4xmpExtNamespace
	MicrosoftPhotoNamespace
	MP1Namespace
	MPNamespace
	MPRINamespace
	MPRegNamespace
	LrNamespace
	DICOMNamespace
	DroneDjiNamespace
	SvgNamespace
	EtNamespace
	PlusNamespace
	PrismNamespace
	PrlNamespace
	PurNamespace
	PmiNamespace
	PrmNamespace
	AcdseeNamespace
	AcdseeRsNamespace
	DigiKamNamespace
	SwfNamespace
	CellNamespace
	AasNamespace
	MwgRsNamespace
	MwgKwNamespace
	MwgCollNamespace
	ExtensisNamespace
	IcsNamespace
	FpvNamespace
	CreatorAtomNamespace
	AppleFiNamespace
	GAudioNamespace
	GImageNamespace
	GPanoNamespace
	GSphericalNamespace
	GDepthNamespace
	GFocusNamespace
	GCameraNamespace
	GCreationsNamespace
	DwcNamespace
	GettyImagesGIFTNamespace
	LImageNamespace
	ProfileNamespace
	SdcNamespace
	AstNamespace
	NineNamespace
	HdrMetadataNamespace
	HdrgmNamespace
	XmpDSANamespace
	SealNamespace
	GContainerNamespace
	HDRGainMapNamespace
	ApdiNamespace
	XMLNamespace
)

var NamespaceArray = []string{
	"unknown",
	"http://ns.adobe.com/exif/1.0/aux/",
	"http://ns.adobe.com/album/1.0/",
	"http://creativecommons.org/ns#",
	"http://ns.adobe.com/camera-raw-defaults/1.0/",
	"http://ns.adobe.com/camera-raw-settings/1.0/",
	"http://ns.adobe.com/camera-raw-saved-settings/1.0/",
	"http://purl.org/dc/elements/1.1/",
	"http://ns.adobe.com/exif/1.0/",
	"http://cipa.jp/exif/1.0/",
	"http://ns.adobe.com/iX/1.0/",
	"http://ns.adobe.com/pdf/1.3/",
	"http://ns.adobe.com/pdfx/1.3/",
	"http://ns.adobe.com/photoshop/1.0/",
	"http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"http://www.w3.org/2000/01/rdf-schema#",
	"http://ns.adobe.com/xap/1.0/sType/Dimensions#",
	"http://ns.adobe.com/xap/1.0/sType/ResourceEvent#",
	"http://ns.adobe.com/xap/1.0/sType/Font#",
	"http://ns.adobe.com/xap/1.0/sType/Job#",
	"http://ns.adobe.com/xap/1.0/sType/ResourceRef#",
	"http://ns.adobe.com/xap/1.0/sType/Version#",
	"http://ns.adobe.com/xap/1.0/sType/ManifestItem#",
	"http://ns.adobe.com/photoshop/1.0/camera-profile",
	"http://ns.adobe.com/camera-raw-embedded-lens-profile/1.0/",
	"http://ns.adobe.com/tiff/1.0/",
	"adobe:ns:meta/",
	"http://ns.adobe.com/xap/1.0/g/",
	"http://ns.adobe.com/xap/1.0/g/img/",
	"http://ns.adobe.com/xap/1.0/",
	"http://ns.adobe.com/xap/1.0/bj/",
	"http://ns.adobe.com/xmp/1.0/DynamicMedia/",
	"http://ns.adobe.com/xap/1.0/mm/",
	"http://ns.adobe.com/xap/1.0/rights/",
	"http://ns.adobe.com/xmp/note/",
	"http://ns.adobe.com/xap/1.0/t/pg/",
	"http://ns.adobe.com/xmp/Identifier/qual/1.0/",
	"http://ns.adobe.com/xap/1.0/PLUS/",
	"http://ns.adobe.com/photoshop/1.0/panorama-profile",
	"http://ns.optimasc.com/dex/1.0/",
	"http://ns.iview-multimedia.com/mediapro/1.0/",
	"http://ns.microsoft.com/expressionmedia/1.0/",
	"http://iptc.org/std/Iptc4xmpCore/1.0/xmlns/",
	"http://iptc.org/std/Iptc4xmpExt/2008-02-29/",
	"http://ns.microsoft.com/photo/1.0",
	"http://ns.microsoft.com/photo/1.1",
	"http://ns.microsoft.com/photo/1.2/",
	"http://ns.microsoft.com/photo/1.2/t/RegionInfo#",
	"http://ns.microsoft.com/photo/1.2/t/Region#",
	"http://ns.adobe.com/lightroom/1.0/",
	"http://ns.adobe.com/DICOM/",
	"http://www.dji.com/drone-dji/1.0/",
	"http://www.w3.org/2000/svg",
	"http://ns.exiftool.org/1.0/",
	"http://ns.useplus.org/ldf/xmp/1.0/",
	"http://prismstandard.org/namespaces/basic/2.0/",
	"http://prismstandard.org/namespaces/prl/2.1/",
	"http://prismstandard.org/namespaces/prismusagerights/2.1/",
	"http://prismstandard.org/namespaces/pmi/2.2/",
	"http://prismstandard.org/namespaces/prm/3.0/",
	"http://ns.acdsee.com/iptc/1.0/",
	"http://ns.acdsee.com/regions/",
	"http://www.digikam.org/ns/1.0/",
	"http://ns.adobe.com/swf/1.0/",
	"http://developer.sonyericsson.com/cell/1.0/",
	"http://ns.apple.com/adjustment-settings/1.0/",
	"http://www.metadataworkinggroup.com/schemas/regions/",
	"http://www.metadataworkinggroup.com/schemas/keywords/",
	"http://www.metadataworkinggroup.com/schemas/collections/",
	"http://ns.extensis.com/extensis/1.0/",
	"http://ns.idimager.com/ics/1.0/",
	"http://ns.fastpictureviewer.com/fpv/1.0/",
	"http://ns.adobe.com/creatorAtom/1.0/",
	"http://ns.apple.com/faceinfo/1.0/",
	"http://ns.google.com/photos/1.0/audio/",
	"http://ns.google.com/photos/1.0/image/",
	"http://ns.google.com/photos/1.0/panorama/",
	"http://ns.google.com/videos/1.0/spherical/",
	"http://ns.google.com/photos/1.0/depthmap/",
	"http://ns.google.com/photos/1.0/focus/",
	"http://ns.google.com/photos/1.0/camera/",
	"http://ns.google.com/photos/1.0/creations/",
	"http://rs.tdwg.org/dwc/index.htm",
	"http://xmp.gettyimages.com/gift/1.0/",
	"http://ns.leiainc.com/photos/1.0/image/",
	"http://ns.google.com/photos/dd/1.0/profile/",
	"http://ns.nikon.com/sdc/1.0/",
	"http://ns.nikon.com/asteroid/1.0/",
	"http://ns.nikon.com/nine/1.0/",
	"http://ns.adobe.com/hdr-metadata/1.0/",
	"http://ns.adobe.com/hdr-gain-map/1.0/",
	"http://leica-camera.com/digital-shift-assistant/1.0/",
	"http://ns.seal/2024/1.0/",
	"http://ns.google.com/photos/1.0/container/",
	"http://ns.apple.com/HDRGainMap/1.0/",
	"http://ns.apple.com/pixeldatainfo/1.0/",
	"xmlnamespace",
}

var PrefixToNamespace = map[string]Namespace{
	"unknown":         UnknownNamespace,
	"aux":             AuxNamespace,
	"album":           AlbumNamespace,
	"cc":              CcNamespace,
	"crd":             CrdNamespace,
	"crs":             CrsNamespace,
	"crss":            CrssNamespace,
	"dc":              DcNamespace,
	"exif":            ExifNamespace,
	"exifEX":          ExifEXNamespace,
	"iX":              IXNamespace,
	"pdf":             PdfNamespace,
	"pdfx":            PdfxNamespace,
	"photoshop":       PhotoshopNamespace,
	"rdf":             RdfNamespace,
	"rdfs":            RdfsNamespace,
	"stDim":           StDimNamespace,
	"stEvt":           StEvtNamespace,
	"stFnt":           StFntNamespace,
	"stJob":           StJobNamespace,
	"stRef":           StRefNamespace,
	"stVer":           StVerNamespace,
	"stMfs":           StMfsNamespace,
	"stCamera":        StCameraNamespace,
	"crlcp":           CrlcpNamespace,
	"tiff":            TiffNamespace,
	"x":               XNamespace,
	"xmpG":            XmpGNamespace,
	"xmpGImg":         XmpGImgNamespace,
	"xmp":             XmpNamespace,
	"xap":             XmpNamespace, // Translate older tag name "xap -> xmp"
	"xmpBJ":           XmpBJNamespace,
	"xapBJ":           XmpBJNamespace, // Translate older tag name "xapBJ -> xmpBJ"
	"xmpDM":           XmpDMNamespace,
	"xmpMM":           XmpMMNamespace,
	"xapMM":           XmpMMNamespace, // Translate older tag name "xapMM -> xmpMM"
	"xmpRights":       XmpRightsNamespace,
	"xmpNote":         XmpNoteNamespace,
	"xmpTPg":          XmpTPgNamespace,
	"xmpidq":          XmpidqNamespace,
	"xmpPLUS":         XmpPLUSNamespace,
	"panorama":        PanoramaNamespace,
	"dex":             DexNamespace,
	"mediapro":        MediaproNamespace,
	"expressionmedia": ExpressionmediaNamespace,
	"Iptc4xmpCore":    Iptc4xmpCoreNamespace,
	"Iptc4xmpExt":     Iptc4xmpExtNamespace,
	"MicrosoftPhoto":  MicrosoftPhotoNamespace,
	"MP1":             MP1Namespace,
	"MP":              MPNamespace,
	"MPRI":            MPRINamespace,
	"MPReg":           MPRegNamespace,
	"lr":              LrNamespace,
	"DICOM":           DICOMNamespace,
	"drone-dji":       DroneDjiNamespace,
	"svg":             SvgNamespace,
	"et":              EtNamespace,
	"plus":            PlusNamespace,
	"prism":           PrismNamespace,
	"prl":             PrlNamespace,
	"pur":             PurNamespace,
	"pmi":             PmiNamespace,
	"prm":             PrmNamespace,
	"acdsee":          AcdseeNamespace,
	"acdsee-rs":       AcdseeRsNamespace,
	"digiKam":         DigiKamNamespace,
	"swf":             SwfNamespace,
	"cell":            CellNamespace,
	"aas":             AasNamespace,
	"mwg-rs":          MwgRsNamespace,
	"mwg-kw":          MwgKwNamespace,
	"mwg-coll":        MwgCollNamespace,
	"extensis":        ExtensisNamespace,
	"ics":             IcsNamespace,
	"fpv":             FpvNamespace,
	"creatorAtom":     CreatorAtomNamespace,
	"apple-fi":        AppleFiNamespace,
	"GAudio":          GAudioNamespace,
	"GImage":          GImageNamespace,
	"GPano":           GPanoNamespace,
	"GSpherical":      GSphericalNamespace,
	"GDepth":          GDepthNamespace,
	"GFocus":          GFocusNamespace,
	"GCamera":         GCameraNamespace,
	"GCreations":      GCreationsNamespace,
	"dwc":             DwcNamespace,
	"GettyImagesGIFT": GettyImagesGIFTNamespace,
	"LImage":          LImageNamespace,
	"Profile":         ProfileNamespace,
	"sdc":             SdcNamespace,
	"ast":             AstNamespace,
	"nine":            NineNamespace,
	"hdr_metadata":    HdrMetadataNamespace,
	"hdrgm":           HdrgmNamespace,
	"xmpDSA":          XmpDSANamespace,
	"seal":            SealNamespace,
	"GContainer":      GContainerNamespace,
	"HDRGainMap":      HDRGainMapNamespace,
	"apdi":            ApdiNamespace,
	"xml":             XMLNamespace,
}

var NamespaceToPrefix = map[Namespace]string{
	UnknownNamespace:         "unknown",
	AuxNamespace:             "aux",
	AlbumNamespace:           "album",
	CcNamespace:              "cc",
	CrdNamespace:             "crd",
	CrsNamespace:             "crs",
	CrssNamespace:            "crss",
	DcNamespace:              "dc",
	ExifNamespace:            "exif",
	ExifEXNamespace:          "exifEX",
	IXNamespace:              "iX",
	PdfNamespace:             "pdf",
	PdfxNamespace:            "pdfx",
	PhotoshopNamespace:       "photoshop",
	RdfNamespace:             "rdf",
	RdfsNamespace:            "rdfs",
	StDimNamespace:           "stDim",
	StEvtNamespace:           "stEvt",
	StFntNamespace:           "stFnt",
	StJobNamespace:           "stJob",
	StRefNamespace:           "stRef",
	StVerNamespace:           "stVer",
	StMfsNamespace:           "stMfs",
	StCameraNamespace:        "stCamera",
	CrlcpNamespace:           "crlcp",
	TiffNamespace:            "tiff",
	XNamespace:               "x",
	XmpGNamespace:            "xmpG",
	XmpGImgNamespace:         "xmpGImg",
	XmpNamespace:             "xmp",
	XmpBJNamespace:           "xmpBJ",
	XmpDMNamespace:           "xmpDM",
	XmpMMNamespace:           "xmpMM",
	XmpRightsNamespace:       "xmpRights",
	XmpNoteNamespace:         "xmpNote",
	XmpTPgNamespace:          "xmpTPg",
	XmpidqNamespace:          "xmpidq",
	XmpPLUSNamespace:         "xmpPLUS",
	PanoramaNamespace:        "panorama",
	DexNamespace:             "dex",
	MediaproNamespace:        "mediapro",
	ExpressionmediaNamespace: "expressionmedia",
	Iptc4xmpCoreNamespace:    "Iptc4xmpCore",
	Iptc4xmpExtNamespace:     "Iptc4xmpExt",
	MicrosoftPhotoNamespace:  "MicrosoftPhoto",
	MP1Namespace:             "MP1",
	MPNamespace:              "MP",
	MPRINamespace:            "MPRI",
	MPRegNamespace:           "MPReg",
	LrNamespace:              "lr",
	DICOMNamespace:           "DICOM",
	DroneDjiNamespace:        "drone-dji",
	SvgNamespace:             "svg",
	EtNamespace:              "et",
	PlusNamespace:            "plus",
	PrismNamespace:           "prism",
	PrlNamespace:             "prl",
	PurNamespace:             "pur",
	PmiNamespace:             "pmi",
	PrmNamespace:             "prm",
	AcdseeNamespace:          "acdsee",
	AcdseeRsNamespace:        "acdsee-rs",
	DigiKamNamespace:         "digiKam",
	SwfNamespace:             "swf",
	CellNamespace:            "cell",
	AasNamespace:             "aas",
	MwgRsNamespace:           "mwg-rs",
	MwgKwNamespace:           "mwg-kw",
	MwgCollNamespace:         "mwg-coll",
	ExtensisNamespace:        "extensis",
	IcsNamespace:             "ics",
	FpvNamespace:             "fpv",
	CreatorAtomNamespace:     "creatorAtom",
	AppleFiNamespace:         "apple-fi",
	GAudioNamespace:          "GAudio",
	GImageNamespace:          "GImage",
	GPanoNamespace:           "GPano",
	GSphericalNamespace:      "GSpherical",
	GDepthNamespace:          "GDepth",
	GFocusNamespace:          "GFocus",
	GCameraNamespace:         "GCamera",
	GCreationsNamespace:      "GCreations",
	DwcNamespace:             "dwc",
	GettyImagesGIFTNamespace: "GettyImagesGIFT",
	LImageNamespace:          "LImage",
	ProfileNamespace:         "Profile",
	SdcNamespace:             "sdc",
	AstNamespace:             "ast",
	NineNamespace:            "nine",
	HdrMetadataNamespace:     "hdr_metadata",
	HdrgmNamespace:           "hdrgm",
	XmpDSANamespace:          "xmpDSA",
	SealNamespace:            "seal",
	GContainerNamespace:      "GContainer",
	HDRGainMapNamespace:      "HDRGainMap",
	ApdiNamespace:            "apdi",
	XMLNamespace:             "xml",
}

// String returns the string representation of the Namespace.
func (ns Namespace) String() string {
	if int(ns) < len(NamespaceArray) {
		return NamespaceArray[ns]
	}
	return "unknown"
}

// Prefix returns the prefix representation of the Namespace.
func (ns Namespace) Prefix() string {
	if prefix, found := NamespaceToPrefix[ns]; found {
		return prefix
	}
	return "unknown"
}

// IdentifyNamespace returns the Namespace for a given prefix.
func IdentifyNamespace(prefix string) (Namespace, bool) {
	ns, found := PrefixToNamespace[prefix]
	return ns, found
}
