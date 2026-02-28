package xmp

import (
	"bufio"
	"bytes"
	"io"
	"sync"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/isobmff"
	"github.com/evanoberholster/imagemeta/meta/jpeg"
	"github.com/orisano/gosax"
)

const (
	parserBufferSize = 64 * 1024
	scanBufferSize   = 4 * 1024
)

var rdfRootProperty = NewProperty(RdfNS, RDF)
var xmpRootToken = []byte("<x:xmpmeta")
var rdfRootToken = []byte("<rdf:RDF")

var saxReaderPool = sync.Pool{
	New: func() interface{} {
		return gosax.NewReaderSize(nil, parserBufferSize)
	},
}

// Parse reads XMP metadata from sidecar files and embedded image payloads.
//
// Supported embedded sources include JPEG APP1 packets and ISOBMFF files such as
// CR3, HEIC/HEIF, AVIF and JXL containers.
func Parse(r io.Reader) (XMP, error) {
	br := asBufferedReader(r)

	fileType, err := imagetype.ScanBuf(br)
	if err != nil {
		return ParseXmp(br)
	}

	switch fileType {
	case imagetype.ImageJPEG:
		return parseJPEG(br)
	case imagetype.ImageCR3, imagetype.ImageHEIC, imagetype.ImageHEIF, imagetype.ImageAVIF:
		return parseISOBMFF(br)
	case imagetype.ImageJXL:
		// JPEG XL may be codestream-based and not always wrapped in BMFF boxes.
		// Direct XMP stream parsing is robust for both sidecar-like and embedded packets.
		return ParseXmp(br)
	default:
		return ParseXmp(br)
	}
}

func parseJPEG(r io.Reader) (XMP, error) {
	var out XMP
	found := false

	err := jpeg.ScanJPEG(
		r,
		nil,
		func(packet io.Reader) error {
			x, err := ParseXmp(packet)
			if err != nil {
				return err
			}
			out = x
			found = true
			return nil
		},
	)
	if err != nil {
		return XMP{}, err
	}
	if !found {
		return XMP{}, ErrNoXMP
	}
	return out, nil
}

func parseISOBMFF(r io.Reader) (XMP, error) {
	var out XMP
	found := false

	bmffReader := isobmff.NewReader(
		r,
		nil,
		func(packet io.Reader, _ isobmff.XPacketHeader) error {
			x, err := ParseXmp(packet)
			if err != nil {
				return err
			}
			out = x
			found = true
			return nil
		},
		nil,
	)
	defer bmffReader.Close()

	if err := bmffReader.ReadFTYP(); err != nil {
		return XMP{}, err
	}

	for {
		err := bmffReader.ReadMetadata()
		if err == nil {
			continue
		}
		if err == io.EOF {
			break
		}
		return XMP{}, err
	}

	if !found {
		return XMP{}, ErrNoXMP
	}
	return out, nil
}

func asBufferedReader(r io.Reader) *bufio.Reader {
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < scanBufferSize {
		return bufio.NewReaderSize(r, scanBufferSize)
	}
	return br
}

func parseXMPStream(r io.Reader) (x XMP, err error) {
	br := asBufferedReader(r)
	start, found, err := seekXMPStart(br)
	if err != nil {
		return XMP{}, err
	}
	if !found {
		return XMP{}, ErrNoXMP
	}

	sax := saxReaderPool.Get().(*gosax.Reader)
	sax.Reset(start)
	sax.EmitSelfClosingTag = true
	defer saxReaderPool.Put(sax)

	parser := xmpStreamParser{r: sax}
	return parser.parse()
}

func seekXMPStart(br *bufio.Reader) (io.Reader, bool, error) {
	iXMP := 0
	iRDF := 0

	for {
		b, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil, false, nil
			}
			return nil, false, err
		}

		iXMP = advancePatternMatch(xmpRootToken, iXMP, b)
		if iXMP == len(xmpRootToken) {
			return io.MultiReader(bytes.NewReader(xmpRootToken), br), true, nil
		}

		iRDF = advancePatternMatch(rdfRootToken, iRDF, b)
		if iRDF == len(rdfRootToken) {
			return io.MultiReader(bytes.NewReader(rdfRootToken), br), true, nil
		}
	}
}

func advancePatternMatch(pattern []byte, idx int, b byte) int {
	if b == pattern[idx] {
		idx++
		if idx == len(pattern) {
			return idx
		}
		return idx
	}
	if b == pattern[0] {
		return 1
	}
	return 0
}

type xmpStreamParser struct {
	r       *gosax.Reader
	xmp     XMP
	started bool
	root    Property

	stack [64]Property
	depth int
}

func (p *xmpStreamParser) parse() (XMP, error) {
	for {
		ev, err := p.r.Event()
		if err != nil {
			return XMP{}, err
		}

		switch ev.Type() {
		case gosax.EventEOF:
			if !p.started {
				return XMP{}, ErrNoXMP
			}
			return p.xmp, nil
		case gosax.EventStart:
			p.parseStart(ev.Bytes)
		case gosax.EventEnd:
			if p.parseEnd(ev.Bytes) {
				return p.xmp, nil
			}
		case gosax.EventText:
			p.parseText(ev.Bytes)
		case gosax.EventCData:
			p.parseCData(ev.Bytes)
		}
	}
}

func (p *xmpStreamParser) parseStart(tag []byte) {
	tagName, attrs := gosax.Name(tag)
	tagProp := identifyProperty(tagName)

	if !p.started {
		if !tagProp.Equals(XMPRootProperty) && !tagProp.Equals(rdfRootProperty) {
			return
		}
		p.started = true
		p.root = tagProp
	}

	for len(attrs) > 0 {
		attr, next, err := gosax.NextAttribute(attrs)
		if err != nil {
			break
		}
		attrs = next
		if len(attr.Key) == 0 {
			break
		}
		p.parseAttribute(tagProp, attr.Key, attr.Value)
	}

	p.push(tagProp)
}

func (p *xmpStreamParser) parseEnd(tag []byte) bool {
	if !p.started {
		return false
	}
	tagName, _ := gosax.Name(tag)
	tagProp := identifyProperty(tagName)
	p.pop(tagProp)
	return tagProp.Equals(p.root)
}

func (p *xmpStreamParser) parseAttribute(tagProp Property, attrName, attrValue []byte) {
	if len(attrValue) < 2 {
		return
	}

	attrProp := identifyProperty(attrName)
	if attrProp.Equals((Property{})) {
		return
	}

	attrVal := unescapeXML(attrValue[1 : len(attrValue)-1])

	prop := property{
		pt:     attrPType,
		parent: tagProp,
		self:   attrProp,
		val:    attrVal,
	}

	if listTarget, ok := p.listTargetFor(tagProp); ok {
		prop.parent = attrProp
		prop.self = listTarget
	}

	_ = p.xmp.parser(prop)
}

func (p *xmpStreamParser) parseText(text []byte) {
	if !p.started || p.depth == 0 {
		return
	}

	value := trimSpace(text)
	if len(value) == 0 {
		return
	}
	value = unescapeXML(value)

	self := p.stack[p.depth-1]
	parent := Property{}
	if p.depth >= 2 {
		parent = p.stack[p.depth-2]
	}

	if listTarget, ok := p.listTargetFor(self); ok {
		self = listTarget
	}

	prop := property{
		pt:     tagPType,
		parent: parent,
		self:   self,
		val:    value,
	}
	_ = p.xmp.parser(prop)
}

func (p *xmpStreamParser) parseCData(cdata []byte) {
	// <![CDATA[value]]>
	if len(cdata) >= 12 {
		p.parseText(cdata[9 : len(cdata)-3])
		return
	}
	p.parseText(cdata)
}

func (p *xmpStreamParser) push(prop Property) {
	if p.depth >= len(p.stack) {
		return
	}
	p.stack[p.depth] = prop
	p.depth++
}

func (p *xmpStreamParser) pop(prop Property) {
	for p.depth > 0 {
		top := p.stack[p.depth-1]
		p.depth--
		if top.Equals(prop) {
			return
		}
	}
}

func (p *xmpStreamParser) listTargetFor(tagProp Property) (Property, bool) {
	if !tagProp.Equals(RDFLi) {
		return Property{}, false
	}

	// Attribute parsing: current tag is not yet pushed on stack.
	if p.depth >= 2 && isListContainer(p.stack[p.depth-1]) {
		return p.stack[p.depth-2], true
	}

	// Text parsing: current tag is already on stack.
	if p.depth >= 3 && p.stack[p.depth-1].Equals(RDFLi) && isListContainer(p.stack[p.depth-2]) {
		return p.stack[p.depth-3], true
	}

	return Property{}, false
}

func isListContainer(prop Property) bool {
	return prop.Equals(RDFSeq) || prop.Equals(RDFBag) || prop.Equals(RDFAlt)
}

func trimSpace(buf []byte) []byte {
	i := 0
	for i < len(buf) {
		switch buf[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			goto right
		}
	}
	return buf[:0]

right:
	j := len(buf) - 1
	for j >= i {
		switch buf[j] {
		case ' ', '\n', '\r', '\t':
			j--
		default:
			return buf[i : j+1]
		}
	}
	return buf[:0]
}

func unescapeXML(buf []byte) []byte {
	if len(buf) == 0 {
		return buf
	}
	if !containsXMLEscape(buf) {
		return buf
	}
	unescaped, err := gosax.Unescape(buf)
	if err != nil {
		return decodeXMLEntities(buf)
	}
	return unescaped
}

func containsXMLEscape(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		switch buf[i] {
		case '&', '\r':
			return true
		}
	}
	return false
}

func identifyProperty(name []byte) Property {
	colon := -1
	for i := 0; i < len(name); i++ {
		if name[i] == ':' {
			colon = i
			break
		}
	}
	if colon <= 0 || colon+1 >= len(name) {
		return Property{}
	}

	prefix := name[:colon]
	local := name[colon+1:]

	ns := IdentifyNamespace(prefix)
	n := IdentifyName(local)

	if ns == UnknownNS || n == UnknownPropertyName {
		return Property{}
	}

	return NewProperty(ns, n)
}

func identifyNamespace(buf []byte) Namespace {
	switch len(buf) {
	case 1:
		switch {
		case eqString(buf, "x"):
			return XNS
		}
	case 2:
		switch {
		case eqString(buf, "dc"):
			return DcNS
		case eqString(buf, "lr"):
			return LrNS
		}
	case 3:
		switch {
		case eqString(buf, "aux"):
			return AuxNS
		case eqString(buf, "crs"):
			return CrsNS
		case eqString(buf, "rdf"):
			return RdfNS
		case eqString(buf, "xml"):
			return XMLNS
		case eqString(buf, "xap"):
			return XapNS
		case eqString(buf, "xmp"):
			return XmpNS
		}
	case 4:
		switch {
		case eqString(buf, "tiff"):
			return TiffNS
		case eqString(buf, "exif"):
			return ExifNS
		}
	case 5:
		switch {
		case eqString(buf, "xmlns"):
			return XMLnsNS
		case eqString(buf, "stEvt"):
			return StEvtNS
		case eqString(buf, "stRef"):
			return StRefNS
		case eqString(buf, "xapMM"):
			return XapMMNS
		case eqString(buf, "xmpDM"):
			return XmpDMNS
		case eqString(buf, "xmpMM"):
			return XmpMMNS
		}
	case 6:
		switch {
		case eqString(buf, "exifEX"):
			return ExifEXNS
		}
	case 9:
		switch {
		case eqString(buf, "photoshop"):
			return PhotoshopNS
		}
	}
	return UnknownNS
}

func identifyName(buf []byte) Name {
	switch {
	case eqFoldString(buf, "xmpmeta"):
		return XmpMeta
	case eqFoldString(buf, "xmpDM"):
		return XmpDM
	case eqFoldString(buf, "XMPToolkit"):
		return XMPToolkit
	case eqFoldString(buf, "xmptk"):
		return XMPToolkit
	case eqFoldString(buf, "xap"):
		return Xap
	case eqFoldString(buf, "RDF"):
		return RDF
	case eqFoldString(buf, "rdf"):
		return RDF
	case eqFoldString(buf, "Description"):
		return Description
	case eqFoldString(buf, "dc"):
		return Dc
	case eqFoldString(buf, "about"):
		return About
	case eqFoldString(buf, "Seq"):
		return Seq
	case eqFoldString(buf, "Bag"):
		return Bag
	case eqFoldString(buf, "Alt"):
		return Alt
	case eqFoldString(buf, "altTapeName"):
		return AltTapeName
	case eqFoldString(buf, "altTimecode"):
		return AltTimecode
	case eqFoldString(buf, "li"):
		return Li
	case eqFoldString(buf, "lang"):
		return Lang
	case eqFoldString(buf, "parseType"):
		return ParseType
	case eqFoldString(buf, "h"):
		return H
	case eqFoldString(buf, "w"):
		return W
	case eqFoldString(buf, "CreateDate"):
		return CreateDate
	case eqFoldString(buf, "CreatorTool"):
		return CreatorTool
	case eqFoldString(buf, "Label"):
		return Label
	case eqFoldString(buf, "NativeDigest"):
		return NativeDigest
	case eqFoldString(buf, "MetadataDate"):
		return MetadataDate
	case eqFoldString(buf, "ModifyDate"):
		return ModifyDate
	case eqFoldString(buf, "Rating"):
		return Rating
	case eqFoldString(buf, "DocumentID"):
		return DocumentID
	case eqFoldString(buf, "DerivedFromDocumentID"):
		return DerivedFromDocumentID
	case eqFoldString(buf, "DerivedFromOriginalDocumentID"):
		return DerivedFromOriginalDocumentID
	case eqFoldString(buf, "DerivedFrom"):
		return DerivedFrom
	case eqFoldString(buf, "OriginalDocumentID"):
		return OriginalDocumentID
	case eqFoldString(buf, "PreservedFileName"):
		return PreservedFileName
	case eqFoldString(buf, "InstanceID"):
		return InstanceID
	case eqFoldString(buf, "PixelXDimension"):
		return PixelXDimension
	case eqFoldString(buf, "PixelYDimension"):
		return PixelYDimension
	case eqFoldString(buf, "ExifImageWidth"):
		return PixelXDimension
	case eqFoldString(buf, "ExifImageHeight"):
		return PixelYDimension
	case eqFoldString(buf, "DateTimeOriginal"):
		return DateTimeOriginal
	case eqFoldString(buf, "DateTimeDigitized"):
		return DateTimeDigitized
	case eqFoldString(buf, "ApertureValue"):
		return ApertureValue
	case eqFoldString(buf, "BrightnessValue"):
		return BrightnessValue
	case eqFoldString(buf, "CameraOwnerName"):
		return CameraOwnerName
	case eqFoldString(buf, "BodySerialNumber"):
		return BodySerialNumber
	case eqFoldString(buf, "ColorSpace"):
		return ColorSpace
	case eqFoldString(buf, "ComponentsConfiguration"):
		return ComponentsConfiguration
	case eqFoldString(buf, "CompressedBitsPerPixel"):
		return CompressedBitsPerPixel
	case eqFoldString(buf, "BitsPerSample"):
		return BitsPerSample
	case eqFoldString(buf, "CustomRendered"):
		return CustomRendered
	case eqFoldString(buf, "DigitalZoomRatio"):
		return DigitalZoomRatio
	case eqFoldString(buf, "ExifVersion"):
		return ExifVersion
	case eqFoldString(buf, "ExposureTime"):
		return ExposureTime
	case eqFoldString(buf, "ExposureProgram"):
		return ExposureProgram
	case eqFoldString(buf, "ExposureMode"):
		return ExposureMode
	case eqFoldString(buf, "FileSource"):
		return FileSource
	case eqFoldString(buf, "FlashpixVersion"):
		return FlashpixVersion
	case eqFoldString(buf, "ExposureBiasValue"):
		return ExposureBiasValue
	case eqFoldString(buf, "ExposureCompensation"):
		return ExposureBiasValue
	case eqFoldString(buf, "FocalLength"):
		return FocalLength
	case eqFoldString(buf, "FocalLengthIn35mmFilm"):
		return FocalLengthIn35mmFilm
	case eqFoldString(buf, "FocalPlaneResolutionUnit"):
		return FocalPlaneResolutionUnit
	case eqFoldString(buf, "FocalPlaneXResolution"):
		return FocalPlaneXResolution
	case eqFoldString(buf, "FocalPlaneYResolution"):
		return FocalPlaneYResolution
	case eqFoldString(buf, "GainControl"):
		return GainControl
	case eqFoldString(buf, "SubjectDistance"):
		return SubjectDistance
	case eqFoldString(buf, "MeteringMode"):
		return MeteringMode
	case eqFoldString(buf, "FNumber"):
		return FNumber
	case eqFoldString(buf, "ISOSpeedRatings"):
		return ISOSpeedRatings
	case eqFoldString(buf, "ISO"):
		return ISOSpeedRatings
	case eqFoldString(buf, "PhotographicSensitivity"):
		return PhotographicSensitivity
	case eqFoldString(buf, "GPSLatitude"):
		return GPSLatitude
	case eqFoldString(buf, "GPSLongitude"):
		return GPSLongitude
	case eqFoldString(buf, "GPSAltitude"):
		return GPSAltitude
	case eqFoldString(buf, "GPSAltitudeRef"):
		return GPSAltitudeRef
	case eqFoldString(buf, "GPSDifferential"):
		return GPSDifferential
	case eqFoldString(buf, "GPSMapDatum"):
		return GPSMapDatum
	case eqFoldString(buf, "GPSStatus"):
		return GPSStatus
	case eqFoldString(buf, "GPSDOP"):
		return GPSDOP
	case eqFoldString(buf, "GPSMeasureMode"):
		return GPSMeasureMode
	case eqFoldString(buf, "GPSSatellites"):
		return GPSSatellites
	case eqFoldString(buf, "GPSTimeStamp"):
		return GPSTimeStamp
	case eqFoldString(buf, "GPSVersionID"):
		return GPSVersionID
	case eqFoldString(buf, "InteroperabilityIndex"):
		return InteroperabilityIndex
	case eqFoldString(buf, "LightSource"):
		return LightSource
	case eqFoldString(buf, "MaxApertureValue"):
		return MaxApertureValue
	case eqFoldString(buf, "PhotometricInterpretation"):
		return PhotometricInterpretation
	case eqFoldString(buf, "RecommendedExposureIndex"):
		return RecommendedExposureIndex
	case eqFoldString(buf, "SamplesPerPixel"):
		return SamplesPerPixel
	case eqFoldString(buf, "PlanarConfiguration"):
		return PlanarConfiguration
	case eqFoldString(buf, "ResolutionUnit"):
		return ResolutionUnit
	case eqFoldString(buf, "Saturation"):
		return Saturation
	case eqFoldString(buf, "Contrast"):
		return Contrast
	case eqFoldString(buf, "SceneCaptureType"):
		return SceneCaptureType
	case eqFoldString(buf, "SceneType"):
		return SceneType
	case eqFoldString(buf, "SensitivityType"):
		return SensitivityType
	case eqFoldString(buf, "Sharpness"):
		return Sharpness
	case eqFoldString(buf, "ShutterSpeedValue"):
		return ShutterSpeedValue
	case eqFoldString(buf, "WhiteBalance"):
		return WhiteBalance
	case eqFoldString(buf, "Temperature"):
		return Temperature
	case eqFoldString(buf, "ColorTemperature"):
		return Temperature
	case eqFoldString(buf, "UserComment"):
		return UserComment
	case eqFoldString(buf, "SubsecTime"):
		return SubsecTime
	case eqFoldString(buf, "SubsecTimeDigitized"):
		return SubsecTimeDigitized
	case eqFoldString(buf, "SubsecTimeOriginal"):
		return SubsecTimeOriginal
	case eqFoldString(buf, "Fired"):
		return Fired
	case eqFoldString(buf, "FlashFired"):
		return Fired
	case eqFoldString(buf, "Return"):
		return Return
	case eqFoldString(buf, "FlashReturn"):
		return Return
	case eqFoldString(buf, "Mode"):
		return Mode
	case eqFoldString(buf, "FlashMode"):
		return Mode
	case eqFoldString(buf, "Function"):
		return Function
	case eqFoldString(buf, "FlashFunction"):
		return Function
	case eqFoldString(buf, "RedEyeMode"):
		return RedEyeMode
	case eqFoldString(buf, "FlashRedEyeMode"):
		return RedEyeMode
	case eqFoldString(buf, "FlashCompensation"):
		return FlashCompensation
	case eqFoldString(buf, "Flash"):
		return FlashTag
	case eqFoldString(buf, "ApproximateFocusDistance"):
		return ApproximateFocusDistance
	case eqFoldString(buf, "ImageNumber"):
		return ImageNumber
	case eqFoldString(buf, "SerialNumber"):
		return SerialNumber
	case eqFoldString(buf, "Lens"):
		return Lens
	case eqFoldString(buf, "LensInfo"):
		return LensInfo
	case eqFoldString(buf, "LensID"):
		return LensID
	case eqFoldString(buf, "LensSerialNumber"):
		return LensSerialNumber
	case eqFoldString(buf, "Firmware"):
		return Firmware
	case eqFoldString(buf, "DistortionCorrectionAlreadyApplied"):
		return DistortionCorrectionAlreadyApplied
	case eqFoldString(buf, "LateralChromaticAberrationCorrectionAlreadyApplied"):
		return LateralChromaticAberrationCorrectionAlreadyApplied
	case eqFoldString(buf, "VignetteCorrectionAlreadyApplied"):
		return VignetteCorrectionAlreadyApplied
	case eqFoldString(buf, "LensModel"):
		return LensModel
	case eqFoldString(buf, "format"):
		return Format
	case eqFoldString(buf, "creator"):
		return Creator
	case eqFoldString(buf, "subject"):
		return Subject
	case eqFoldString(buf, "rights"):
		return Rights
	case eqFoldString(buf, "title"):
		return Title
	case eqFoldString(buf, "description"):
		return Description
	case eqFoldString(buf, "Make"):
		return Make
	case eqFoldString(buf, "Model"):
		return Model
	case eqFoldString(buf, "ImageWidth"):
		return ImageWidth
	case eqFoldString(buf, "ImageLength"):
		return ImageLength
	case eqFoldString(buf, "ImageHeight"):
		return ImageLength
	case eqFoldString(buf, "ImageDescription"):
		return ImageDescription
	case eqFoldString(buf, "Orientation"):
		return Orientation
	case eqFoldString(buf, "Compression"):
		return Compression
	case eqFoldString(buf, "RawFileName"):
		return RawFileName
	case eqFoldString(buf, "Software"):
		return Software
	case eqFoldString(buf, "DateCreated"):
		return DateCreated
	case eqFoldString(buf, "EmbeddedXMPDigest"):
		return EmbeddedXMPDigest
	case eqFoldString(buf, "SidecarForExtension"):
		return SidecarForExtension
	case eqFoldString(buf, "ColorMode"):
		return ColorMode
	case eqFoldString(buf, "ICCProfile"):
		return ICCProfile
	case eqFoldString(buf, "LegacyIPTCDigest"):
		return LegacyIPTCDigest
	case eqFoldString(buf, "History"):
		return HistoryTag
	case eqFoldString(buf, "HistoryAction"):
		return Action
	case eqFoldString(buf, "HistoryChanged"):
		return Changed
	case eqFoldString(buf, "HistoryInstanceID"):
		return InstanceID
	case eqFoldString(buf, "HistoryParameters"):
		return Parameters
	case eqFoldString(buf, "HistorySoftwareAgent"):
		return SoftwareAgent
	case eqFoldString(buf, "HistoryWhen"):
		return When
	case eqFoldString(buf, "parameters"):
		return Parameters
	case eqFoldString(buf, "action"):
		return Action
	case eqFoldString(buf, "changed"):
		return Changed
	case eqFoldString(buf, "softwareAgent"):
		return SoftwareAgent
	case eqFoldString(buf, "startTimecode"):
		return StartTimecode
	case eqFoldString(buf, "stDim"):
		return StDim
	case eqFoldString(buf, "tapeName"):
		return TapeName
	case eqFoldString(buf, "timeValue"):
		return TimeValue
	case eqFoldString(buf, "when"):
		return When
	case eqFoldString(buf, "pick"):
		return Pick
	case eqFoldString(buf, "good"):
		return Good
	case eqFoldString(buf, "hierarchicalSubject"):
		return HierarchicalSubject
	case eqFoldString(buf, "HueAdjustmentRed"):
		return HueAdjustmentRed
	case eqFoldString(buf, "HueAdjustmentOrange"):
		return HueAdjustmentOrange
	case eqFoldString(buf, "HueAdjustmentYellow"):
		return HueAdjustmentYellow
	case eqFoldString(buf, "HueAdjustmentGreen"):
		return HueAdjustmentGreen
	case eqFoldString(buf, "HueAdjustmentAqua"):
		return HueAdjustmentAqua
	case eqFoldString(buf, "HueAdjustmentBlue"):
		return HueAdjustmentBlue
	case eqFoldString(buf, "HueAdjustmentPurple"):
		return HueAdjustmentPurple
	case eqFoldString(buf, "HueAdjustmentMagenta"):
		return HueAdjustmentMagenta
	case eqFoldString(buf, "RedHue"):
		return HueAdjustmentRed
	case eqFoldString(buf, "GreenHue"):
		return HueAdjustmentGreen
	case eqFoldString(buf, "BlueHue"):
		return HueAdjustmentBlue
	case eqFoldString(buf, "SaturationAdjustmentRed"):
		return SaturationAdjustmentRed
	case eqFoldString(buf, "SaturationAdjustmentOrange"):
		return SaturationAdjustmentOrange
	case eqFoldString(buf, "SaturationAdjustmentYellow"):
		return SaturationAdjustmentYellow
	case eqFoldString(buf, "SaturationAdjustmentGreen"):
		return SaturationAdjustmentGreen
	case eqFoldString(buf, "SaturationAdjustmentAqua"):
		return SaturationAdjustmentAqua
	case eqFoldString(buf, "SaturationAdjustmentBlue"):
		return SaturationAdjustmentBlue
	case eqFoldString(buf, "SaturationAdjustmentPurple"):
		return SaturationAdjustmentPurple
	case eqFoldString(buf, "SaturationAdjustmentMagenta"):
		return SaturationAdjustmentMagenta
	case eqFoldString(buf, "RedSaturation"):
		return SaturationAdjustmentRed
	case eqFoldString(buf, "GreenSaturation"):
		return SaturationAdjustmentGreen
	case eqFoldString(buf, "BlueSaturation"):
		return SaturationAdjustmentBlue
	case eqFoldString(buf, "LuminanceAdjustmentRed"):
		return LuminanceAdjustmentRed
	case eqFoldString(buf, "LuminanceAdjustmentOrange"):
		return LuminanceAdjustmentOrange
	case eqFoldString(buf, "LuminanceAdjustmentYellow"):
		return LuminanceAdjustmentYellow
	case eqFoldString(buf, "LuminanceAdjustmentGreen"):
		return LuminanceAdjustmentGreen
	case eqFoldString(buf, "LuminanceAdjustmentAqua"):
		return LuminanceAdjustmentAqua
	case eqFoldString(buf, "LuminanceAdjustmentBlue"):
		return LuminanceAdjustmentBlue
	case eqFoldString(buf, "LuminanceAdjustmentPurple"):
		return LuminanceAdjustmentPurple
	case eqFoldString(buf, "LuminanceAdjustmentMagenta"):
		return LuminanceAdjustmentMagenta
	case eqFoldString(buf, "weightedFlatSubject"):
		return WeightedFlatSubject
	case eqFoldString(buf, "videoFieldOrder"):
		return VideoFieldOrder
	case eqFoldString(buf, "videoFrameRate"):
		return VideoFrameRate
	case eqFoldString(buf, "videoFrameSize"):
		return VideoFrameSize
	case eqFoldString(buf, "videoPixelAspectRatio"):
		return VideoPixelAspectRatio
	case eqFoldString(buf, "videoPixelDepth"):
		return VideoPixelDepth
	case eqFoldString(buf, "Title"):
		return Title
	case eqFoldString(buf, "Unknown"):
		return UnknownPropertyName
	case eqFoldString(buf, "instanceID"):
		return InstanceID
	case eqFoldString(buf, "XResolution"):
		return XResolution
	case eqFoldString(buf, "YCbCrPositioning"):
		return YCbCrPositioning
	case eqFoldString(buf, "YResolution"):
		return YResolution
	case eqFoldString(buf, "AlreadyApplied"):
		return AlreadyApplied
	case eqFoldString(buf, "ToneCurve"):
		return ToneCurve
	case eqFoldString(buf, "ToneCurveRed"):
		return ToneCurveRed
	case eqFoldString(buf, "ToneCurveGreen"):
		return ToneCurveGreen
	case eqFoldString(buf, "ToneCurveBlue"):
		return ToneCurveBlue
	case eqFoldString(buf, "ToneCurvePV2012"):
		return ToneCurvePV2012
	case eqFoldString(buf, "ToneCurvePV2012Red"):
		return ToneCurvePV2012Red
	case eqFoldString(buf, "ToneCurvePV2012Green"):
		return ToneCurvePV2012Green
	case eqFoldString(buf, "ToneCurvePV2012Blue"):
		return ToneCurvePV2012Blue
	default:
		return UnknownPropertyName
	}
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

func eqFoldString(buf []byte, s string) bool {
	if len(buf) != len(s) {
		return false
	}
	for i := 0; i < len(buf); i++ {
		a := buf[i]
		b := s[i]
		if a == b {
			continue
		}
		if 'A' <= a && a <= 'Z' {
			a += 'a' - 'A'
		}
		if 'A' <= b && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}
