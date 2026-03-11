package tag

import (
	"strconv"
	"strings"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
)

// ValueNameFor returns the symbolic name for a numeric enum value.
func ValueNameFor(directoryType ifd.Type, id ID, value uint32) string {
	if names := valueNames(directoryType, id); names != nil {
		if name, ok := names[value]; ok {
			return name
		}
	}
	return strconv.FormatUint(uint64(value), 10)
}

// ParseValueID parses either a numeric string or known symbolic enum value.
func ParseValueID(directoryType ifd.Type, id ID, raw string) (uint32, bool) {
	if v, ok := parseUint(raw); ok {
		return v, true
	}
	byName := valueIDs(directoryType, id)
	if len(byName) == 0 {
		return 0, false
	}
	v, ok := byName[normalizeValueString(raw)]
	return v, ok
}

func valueNames(directoryType ifd.Type, id ID) map[uint32]string {
	switch directoryType {
	case ifd.IFD0, ifd.IFD1, ifd.IFD2:
		return ifd0ValueNames[id]
	case ifd.ExifIFD, ifd.SubIFD0, ifd.SubIFD1, ifd.SubIFD2, ifd.SubIFD3, ifd.SubIFD4, ifd.SubIFD5, ifd.SubIFD6, ifd.SubIFD7:
		return exifValueNames[id]
	}
	return nil
}

func valueIDs(directoryType ifd.Type, id ID) map[string]uint32 {
	switch directoryType {
	case ifd.IFD0, ifd.IFD1, ifd.IFD2:
		return ifd0ValueIDs[id]
	case ifd.ExifIFD, ifd.SubIFD0, ifd.SubIFD1, ifd.SubIFD2, ifd.SubIFD3, ifd.SubIFD4, ifd.SubIFD5, ifd.SubIFD6, ifd.SubIFD7:
		return exifValueIDs[id]
	}
	return nil
}

func normalizeValueString(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.ReplaceAll(v, "_", " ")
	return strings.Join(strings.Fields(v), " ")
}

func parseUint(v string) (uint32, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, false
	}
	if i := strings.IndexByte(v, '('); i > 0 {
		v = strings.TrimSpace(v[:i])
	}
	fields := strings.Fields(v)
	if len(fields) > 0 {
		v = fields[0]
	}
	v = strings.Trim(v, ",;")
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
		n, err := strconv.ParseUint(v[2:], 16, 32)
		return uint32(n), err == nil
	}
	n, err := strconv.ParseUint(v, 10, 32)
	return uint32(n), err == nil
}

func invertValueMap(src map[uint32]string, aliases map[string]uint32) map[string]uint32 {
	dst := make(map[string]uint32, len(src)+len(aliases))
	for id, name := range src {
		dst[normalizeValueString(name)] = id
	}
	for alias, id := range aliases {
		dst[normalizeValueString(alias)] = id
	}
	return dst
}

var ifd0ValueNames = map[ID]map[uint32]string{
	TagPhotometricInterpretation: {
		0:     "WhiteIsZero",
		1:     "BlackIsZero",
		2:     "RGB",
		3:     "RGB Palette",
		4:     "Transparency Mask",
		5:     "CMYK",
		6:     "YCbCr",
		8:     "CIELab",
		32803: "Color Filter Array",
		34892: "Linear Raw",
	},
	TagPlanarConfiguration: {
		1: "Chunky",
		2: "Planar",
	},
	TagOrientation: {
		1: "Horizontal (normal)",
		2: "Mirror horizontal",
		3: "Rotate 180",
		4: "Mirror vertical",
		5: "Mirror horizontal and rotate 270 CW",
		6: "Rotate 90 CW",
		7: "Mirror horizontal and rotate 90 CW",
		8: "Rotate 270 CW",
	},
	TagResolutionUnit: {
		1: "None",
		2: "inches",
		3: "cm",
	},
	TagCompression: {
		1:     "Uncompressed",
		2:     "CCITT 1D",
		3:     "T4/Group 3 Fax",
		4:     "T6/Group 4 Fax",
		5:     "LZW",
		6:     "JPEG (old-style)",
		7:     "JPEG",
		8:     "Adobe Deflate",
		9:     "JBIG B&W",
		10:    "JBIG Color",
		99:    "JPEG",
		262:   "Kodak 262",
		32766: "NeXt or Sony ARW Compressed 2",
		32767: "Sony ARW Compressed",
		32769: "Packed RAW",
		32770: "Samsung SRW Compressed",
		32771: "CCIRLEW",
		32772: "Samsung SRW Compressed 2",
		32773: "PackBits",
		32809: "Thunderscan",
		32867: "Kodak KDC Compressed",
		32895: "IT8CTPAD",
		32896: "IT8LW",
		32897: "IT8MP",
		32898: "IT8BL",
		32908: "PixarFilm",
		32909: "PixarLog",
		32946: "Deflate",
		32947: "DCS",
		33003: "Aperio JPEG 2000 YCbCr",
		33005: "Aperio JPEG 2000 RGB",
		34661: "JBIG",
		34676: "SGILog",
		34677: "SGILog24",
		34712: "JPEG 2000",
		34713: "Nikon NEF Compressed",
		34715: "JBIG2 TIFF FX",
		34718: "Microsoft Document Imaging (MDI) Binary Level Codec",
		34719: "Microsoft Document Imaging (MDI) Progressive Transform Codec",
		34720: "Microsoft Document Imaging (MDI) Vector",
		34887: "ESRI Lerc",
		34892: "Lossy JPEG",
		34925: "LZMA2",
		34926: "Zstd (old)",
		34927: "WebP (old)",
		34933: "PNG",
		34934: "JPEG XR",
		50000: "Zstd",
		50001: "WebP",
		50002: "JPEG XL (old)",
		52546: "JPEG XL",
		65000: "Kodak DCR Compressed",
		65535: "Pentax PEF Compressed",
	},
}

var exifValueNames = map[ID]map[uint32]string{
	TagExposureProgram: {
		0: "Not Defined",
		1: "Manual",
		2: "Program AE",
		3: "Aperture-priority AE",
		4: "Shutter speed priority AE",
		5: "Creative (Slow speed)",
		6: "Action (High speed)",
		7: "Portrait",
		8: "Landscape",
		9: "Bulb",
	},
	TagMeteringMode: {
		0:   "Unknown",
		1:   "Average",
		2:   "Center-weighted average",
		3:   "Spot",
		4:   "Multi-spot",
		5:   "Multi-segment",
		6:   "Partial",
		255: "Other",
	},
	TagLightSource: {
		0:   "Unknown",
		1:   "Daylight",
		2:   "Fluorescent",
		3:   "Tungsten (Incandescent)",
		4:   "Flash",
		9:   "Fine Weather",
		10:  "Cloudy",
		11:  "Shade",
		12:  "Daylight Fluorescent",
		13:  "Day White Fluorescent",
		14:  "Cool White Fluorescent",
		15:  "White Fluorescent",
		16:  "Warm White Fluorescent",
		17:  "Standard Light A",
		18:  "Standard Light B",
		19:  "Standard Light C",
		20:  "D55",
		21:  "D65",
		22:  "D75",
		23:  "D50",
		24:  "ISO Studio Tungsten",
		255: "Other",
	},
	TagExposureMode: {
		0: "Auto",
		1: "Manual",
		2: "Auto bracket",
	},
	TagColorSpace: {
		1:      "sRGB",
		0xffff: "Uncalibrated",
	},
	TagCustomRendered: {
		0: "Normal",
		1: "Custom",
	},
	TagWhiteBalance: {
		0: "Auto",
		1: "Manual",
	},
	TagSceneCaptureType: {
		0: "Standard",
		1: "Landscape",
		2: "Portrait",
		3: "Night scene",
	},
	TagGainControl: {
		0: "None",
		1: "Low gain up",
		2: "High gain up",
		3: "Low gain down",
		4: "High gain down",
	},
	TagContrast: {
		0: "Normal",
		1: "Low",
		2: "High",
	},
	TagSaturation: {
		0: "Normal",
		1: "Low",
		2: "High",
	},
	TagSharpness: {
		0: "Normal",
		1: "Soft",
		2: "Hard",
	},
	TagSubjectDistanceRange: {
		0: "Unknown",
		1: "Macro",
		2: "Close View",
		3: "Distant View",
	},
	TagSensingMethod: {
		1: "Not defined",
		2: "One-chip color area",
		3: "Two-chip color area",
		4: "Three-chip color area",
		5: "Color sequential area",
		7: "Trilinear",
		8: "Color sequential linear",
	},
	TagSceneType: {
		1: "Directly photographed",
	},
	TagFileSource: {
		3: "Digital Camera",
	},
	TagCompositeImage: {
		0: "Unknown",
		1: "NonComposite",
		2: "GeneralComposite",
		3: "CompositeCapturedWhenShooting",
	},
	TagSensitivityType: {
		0: "Unknown",
		1: "Standard Output Sensitivity",
		2: "Recommended Exposure Index",
		3: "ISO Speed",
		4: "Standard Output Sensitivity and Recommended Exposure Index",
		5: "Standard Output Sensitivity and ISO Speed",
		6: "Recommended Exposure Index and ISO Speed",
		7: "Standard Output Sensitivity, Recommended Exposure Index and ISO Speed",
	},
	TagFlash: {
		0x0:  "No Flash",
		0x1:  "Fired",
		0x5:  "Fired, Return not detected",
		0x7:  "Fired, Return detected",
		0x8:  "On, Did not fire",
		0x9:  "On, Fired",
		0xd:  "On, Return not detected",
		0xf:  "On, Return detected",
		0x10: "Off, Did not fire",
		0x14: "Off, Did not fire, Return not detected",
		0x18: "Auto, Did not fire",
		0x19: "Auto, Fired",
		0x1d: "Auto, Fired, Return not detected",
		0x1f: "Auto, Fired, Return detected",
		0x20: "No flash function",
		0x30: "Off, No flash function",
		0x41: "Fired, Red-eye reduction",
		0x45: "Fired, Red-eye reduction, Return not detected",
		0x47: "Fired, Red-eye reduction, Return detected",
		0x49: "On, Red-eye reduction",
		0x4d: "On, Red-eye reduction, Return not detected",
		0x4f: "On, Red-eye reduction, Return detected",
		0x50: "Off, Red-eye reduction",
		0x58: "Auto, Did not fire, Red-eye reduction",
		0x59: "Auto, Fired, Red-eye reduction",
		0x5d: "Auto, Fired, Red-eye reduction, Return not detected",
		0x5f: "Auto, Fired, Red-eye reduction, Return detected",
	},
}

var ifd0ValueIDs = map[ID]map[string]uint32{
	TagPhotometricInterpretation: invertValueMap(ifd0ValueNames[TagPhotometricInterpretation], nil),
	TagPlanarConfiguration: invertValueMap(ifd0ValueNames[TagPlanarConfiguration], map[string]uint32{
		"contiguous": 1,
		"separate":   2,
	}),
	TagOrientation: invertValueMap(ifd0ValueNames[TagOrientation], map[string]uint32{
		"Horizontal": 1,
	}),
	TagResolutionUnit: invertValueMap(ifd0ValueNames[TagResolutionUnit], nil),
	TagCompression:    invertValueMap(ifd0ValueNames[TagCompression], nil),
}

var exifValueIDs = map[ID]map[string]uint32{
	TagExposureProgram:      invertValueMap(exifValueNames[TagExposureProgram], nil),
	TagMeteringMode:         invertValueMap(exifValueNames[TagMeteringMode], nil),
	TagLightSource:          invertValueMap(exifValueNames[TagLightSource], nil),
	TagExposureMode:         invertValueMap(exifValueNames[TagExposureMode], nil),
	TagColorSpace:           invertValueMap(exifValueNames[TagColorSpace], nil),
	TagCustomRendered:       invertValueMap(exifValueNames[TagCustomRendered], nil),
	TagWhiteBalance:         invertValueMap(exifValueNames[TagWhiteBalance], nil),
	TagSceneCaptureType:     invertValueMap(exifValueNames[TagSceneCaptureType], nil),
	TagGainControl:          invertValueMap(exifValueNames[TagGainControl], nil),
	TagContrast:             invertValueMap(exifValueNames[TagContrast], nil),
	TagSaturation:           invertValueMap(exifValueNames[TagSaturation], nil),
	TagSharpness:            invertValueMap(exifValueNames[TagSharpness], nil),
	TagSubjectDistanceRange: invertValueMap(exifValueNames[TagSubjectDistanceRange], nil),
	TagSensingMethod:        invertValueMap(exifValueNames[TagSensingMethod], nil),
	TagFileSource:           invertValueMap(exifValueNames[TagFileSource], nil),
	TagSceneType:            invertValueMap(exifValueNames[TagSceneType], nil),
	TagCompositeImage:       invertValueMap(exifValueNames[TagCompositeImage], nil),
	TagSensitivityType:      invertValueMap(exifValueNames[TagSensitivityType], nil),
	TagFlash:                invertValueMap(exifValueNames[TagFlash], nil),
}
