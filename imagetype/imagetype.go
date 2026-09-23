// Package imagetype provides types and functions for identifying Image document types.
//
// Detection model: Scan/ScanBuf/ReadAt classify from magic bytes in a
// 64-byte header window, with a bounded deeper probe for TIFF Make/Model
// sniffing. Types without a byte signature (TGA, FPX, MAGICK, MPO and
// friends) resolve through extensions and MIME types only via FromString.
// SVG detection is best-effort inside the window. Extension() returns the
// canonical extension without a leading dot.
package imagetype

import (
	"bytes"
	"errors"
	"mime"
	"path/filepath"
	"strings"
)

var (
	// ErrDataLength is an error for data length
	ErrDataLength = errors.New("error the data is not long enough")
)

//go:generate msgp

// FileType is type of Image or Metadata file
type FileType uint8

// MIMEType is the canonical MIME type for a file type.
type MIMEType string

// FileTypeExtension is the canonical filename extension for a file type.
type FileTypeExtension string

// ImageType is kept as an alias for backward compatibility.
type ImageType = FileType

func (mt MIMEType) String() string {
	return string(mt)
}

func (ext FileTypeExtension) String() string {
	return string(ext)
}

// IsUnknown returns true if the file type is unknown.
func (ft FileType) IsUnknown() bool {
	return ft == ImageUnknown
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (ft FileType) MarshalText() (text []byte, err error) {
	return []byte(ft.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (ft *FileType) UnmarshalText(text []byte) (err error) {
	*ft = FromString(string(text))
	return nil
}

// String returns the canonical MIME type for the file type.
func (ft FileType) String() string {
	return ft.MIMEType().String()
}

// MIMEType returns the canonical MIME type for the file type.
func (ft FileType) MIMEType() MIMEType {
	if mimeType, ok := fileTypeCanonicalMIME[ft]; ok {
		return mimeType
	}
	return fileTypeCanonicalMIME[ImageUnknown]
}

// Extension returns the canonical file extension for the file type,
// without a leading dot ("jpg", not ".jpg").
func (ft FileType) Extension() string {
	return ft.FileTypeExtension().String()
}

// FileTypeExtension returns the canonical filename extension for the file type.
func (ft FileType) FileTypeExtension() FileTypeExtension {
	if ext, ok := fileTypeCanonicalExtension[ft]; ok {
		return ext
	}
	return fileTypeCanonicalExtension[ImageUnknown]
}

// Family returns the broad semantic family for the image type.
func (ft FileType) Family() MediaType {
	return ft.MediaType()
}

// MediaType returns the broad semantic media class for the file type.
func (ft FileType) MediaType() MediaType {
	switch {
	case ft.IsUnknown():
		return MediaTypeUnknown
	case ft == ImageXMP:
		return MediaTypeMetadata
	case ft == ImageSVG:
		return MediaTypeVector
	case ft.IsRAW():
		return MediaTypeRaw
	default:
		return MediaTypeRaster
	}
}

// BaseType returns the container/signature family used to classify this type.
func (ft FileType) BaseType() BaseType {
	if ft.IsUnknown() {
		return BaseTypeUnknown
	}

	switch ft {
	case ImageJPEG:
		return BaseTypeJPEG
	case ImagePNG, ImageAPNG:
		return BaseTypePNG
	case ImageGIF:
		return BaseTypeGIF
	case ImageBMP:
		return BaseTypeBMP
	case ImageWebP:
		return BaseTypeRIFF
	case ImageHEIF, ImageHEIC, ImageAVIF, ImageCR3:
		return BaseTypeISOBMFF
	case ImageTiff, ImageDNG, ImageNEF, ImagePanaRAW, ImageARW, ImageCR2, ImageGPR,
		ImageRAF, ImageORF, ImageSRW, ImagePEF, ImageRWL, ImageIIQ, Image3FR, ImageX3F,
		ImageMRW, ImageKDC, ImageDCR, ImageERF, ImageNRW, ImageSR2, ImageSRF, ImageFFF,
		ImageMOS, ImageK25, ImageMEF:
		return BaseTypeTIFF
	case ImageCRW:
		return BaseTypeCIFF
	case ImagePSD:
		return BaseTypePSD
	case ImageXMP:
		return BaseTypeXML
	case ImagePPM, ImagePBM, ImagePGM, ImagePNM, ImagePAM:
		return BaseTypeNetpbm
	case ImageJP2K:
		return BaseTypeJP2
	case ImageSVG:
		return BaseTypeSVG
	case ImageMAGICK:
		return BaseTypeMagick
	case ImageICO, ImageCUR:
		return BaseTypeICO
	case ImageTGA:
		return BaseTypeTGA
	case ImageDDS:
		return BaseTypeDDS
	case ImageEXR:
		return BaseTypeEXR
	case ImageHDR:
		return BaseTypeHDR
	case ImageJXL:
		return BaseTypeJXL
	case ImageJXR:
		return BaseTypeJXR
	case ImageMNG:
		return BaseTypeMNG
	case ImageJNG:
		return BaseTypeJNG
	case ImageMPO:
		return BaseTypeMPO
	case ImageDPX:
		return BaseTypeDPX
	case ImageFITS:
		return BaseTypeFITS
	case ImageDCM:
		return BaseTypeDICOM
	case ImageFPX:
		return BaseTypeFPX
	case ImageDJVU:
		return BaseTypeDJVU
	case ImagePCX:
		return BaseTypePCX
	case ImageWPG:
		return BaseTypeWPG
	case ImagePICT:
		return BaseTypePICT
	case ImagePCD:
		return BaseTypePCD
	case ImageBPG:
		return BaseTypeBPG
	case ImageFLIF:
		return BaseTypeFLIF
	case ImagePGF:
		return BaseTypePGF
	case ImageXCF:
		return BaseTypeXCF
	case ImageQTIF:
		return BaseTypeQTIF
	case ImageJ2C:
		return BaseTypeJP2
	case ImageXISF:
		return BaseTypeXISF
	case ImageRAW:
		return BaseTypeUnknown
	default:
		return BaseTypeUnknown
	}
}

// Container is kept for backward compatibility.
func (ft FileType) Container() BaseType {
	return ft.BaseType()
}

// IsRAW returns true for camera raw and raw-family image types.
func (ft FileType) IsRAW() bool {
	switch ft {
	case ImageRAW, ImageDNG, ImageNEF, ImagePanaRAW, ImageARW, ImageCRW, ImageGPR,
		ImageCR3, ImageCR2, ImageRAF, ImageORF, ImageSRW, ImagePEF, ImageRWL, ImageIIQ,
		Image3FR, ImageX3F, ImageMRW, ImageKDC, ImageDCR, ImageERF, ImageNRW, ImageSR2,
		ImageSRF, ImageFFF, ImageMOS, ImageK25, ImageMEF:
		return true
	default:
		return false
	}
}

// IsISOBMFF returns true when the file type is based on the ISO Base Media File Format.
func (ft FileType) IsISOBMFF() bool {
	return ft.BaseType() == BaseTypeISOBMFF
}

// FromString returns a FileType for the given content-type string, extension,
// or filename.
//
// Resolution order: exact MIME match, MIME with parameters, dotted
// extension suffix (covers bare tokens, dotted tokens, paths and URLs
// with query/fragment suffixes). Anything else is ImageUnknown.
func FromString(str string) FileType {
	str = strings.TrimSpace(str)
	if str == "" {
		return ImageUnknown
	}

	normalized := strings.ToLower(str)

	// from content-type, exact or with parameters
	if it, ok := mimeTypeValues[MIMEType(normalized)]; ok {
		return it
	}
	if idx := strings.IndexByte(normalized, ';'); idx > 0 {
		if mediaType, _, err := mime.ParseMediaType(normalized); err == nil {
			if it, ok := mimeTypeValues[MIMEType(mediaType)]; ok {
				return it
			}
		}
		if it, ok := mimeTypeValues[MIMEType(strings.TrimSpace(normalized[:idx]))]; ok {
			return it
		}
	}

	// from extension: strip query/fragment suffixes, then take the dotted
	// suffix. A bare token ("jpeg") probes its dotted form (".jpeg").
	candidate := normalized
	if idx := strings.IndexAny(candidate, "?#"); idx > 0 {
		candidate = candidate[:idx]
	}
	if ext := filepath.Ext(candidate); ext != "" {
		if it, ok := fileTypeExtensions[FileTypeExtension(ext)]; ok {
			return it
		}
	}
	if !strings.HasPrefix(candidate, ".") {
		if it, ok := fileTypeExtensions[FileTypeExtension("."+candidate)]; ok {
			return it
		}
	}

	return ImageUnknown
}

// maxTokenLen bounds the stack buffer used for case-insensitive matching.
// MIME types and extensions are well under this size; longer inputs (file
// paths, URLs) fall back to FromString.
const maxTokenLen = 64

// FromBytes returns a FileType for content-type bytes, extensions, or filenames.
//
// It never modifies buf. Token lookups allocate nothing; only inputs that
// need the full FromString fallback (paths, parameters beyond the fast
// paths) may allocate.
func FromBytes(buf []byte) FileType {
	buf = bytes.TrimSpace(buf)
	if len(buf) == 0 {
		return ImageUnknown
	}

	if it, ok := lookupToken(buf); ok {
		return it
	}

	// Case-insensitive retry on a bounded stack copy; the caller's buffer
	// is never modified.
	if len(buf) > maxTokenLen {
		return FromString(string(buf))
	}
	var tmp [maxTokenLen]byte
	lower := tmp[:len(buf)]
	for i, c := range buf {
		if c >= 'A' && c <= 'Z' {
			c |= 0x20
		}
		lower[i] = c
	}
	if it, ok := lookupToken(lower); ok {
		return it
	}

	if idx := bytes.IndexByte(lower, ';'); idx > 0 {
		if it, ok := lookupToken(bytes.TrimSpace(lower[:idx])); ok {
			return it
		}
	}

	return FromString(string(buf))
}

// lookupToken matches a MIME type, extension, or dotted extension against
// the canonical tables. The string conversions in map-index position do not
// allocate.
func lookupToken(buf []byte) (FileType, bool) {
	if it, ok := mimeTypeValues[MIMEType(buf)]; ok {
		return it, true
	}
	if it, ok := fileTypeExtensions[FileTypeExtension(buf)]; ok {
		return it, true
	}
	if len(buf) > 0 && buf[0] != '.' {
		// Dotted-extension retry without allocating: bound the probe to
		// tokens that fit alongside the dot.
		var tmp [maxTokenLen]byte
		if len(buf)+1 <= len(tmp) {
			tmp[0] = '.'
			copy(tmp[1:], buf)
			if it, ok := fileTypeExtensions[FileTypeExtension(tmp[:len(buf)+1])]; ok {
				return it, true
			}
		}
	}
	return ImageUnknown, false
}

// Image file types Raw/Compressed/JPEG
const (
	ImageUnknown FileType = iota
	ImageJPEG
	ImagePNG
	ImageGIF
	ImageBMP
	ImageWebP
	ImageHEIF
	ImageRAW
	ImageTiff
	ImageDNG
	ImageNEF
	ImagePanaRAW
	ImageARW
	ImageCRW
	ImageGPR
	ImageCR3
	ImageCR2
	ImagePSD
	ImageXMP
	ImageAVIF
	ImagePPM
	ImageJP2K   // JP2K represents the JPEG 2000 image type.
	ImageSVG    // SVG represents the SVG image type.
	ImageMAGICK // MAGICK represents the libmagick compatible genetic image type.

	// --- additions ---
	ImageICO  // Windows icon
	ImageCUR  // Windows cursor
	ImageTGA  // Truevision TGA
	ImageDDS  // DirectDraw Surface
	ImageEXR  // OpenEXR
	ImageHDR  // Radiance HDR (RGBE)
	ImageJXL  // JPEG XL
	ImageHEIC // HEIC (HEIF family, but common as its own label)
	ImageAPNG // Animated PNG (still PNG container, but useful to distinguish)
	ImagePBM  // Netpbm PBM
	ImagePGM  // Netpbm PGM
	ImagePNM  // Netpbm PNM (generic)
	ImagePAM  // Netpbm PAM

	// More camera RAWs
	ImageRAF // Fujifilm RAF
	ImageORF // Olympus ORF
	ImageSRW // Samsung SRW
	ImagePEF // Pentax PEF
	ImageRWL // Leica RWL
	ImageIIQ // Phase One IIQ
	Image3FR // Hasselblad 3FR
	ImageX3F // Sigma X3F
	ImageMRW // Minolta MRW
	ImageKDC // Kodak KDC
	ImageDCR // Kodak DCR
	ImageERF // Epson ERF

	// --- added (ExifTool-supported image formats you were missing) ---
	ImageJXR  // JPEG XR / HD Photo (JXR, WDP, HDP)
	ImageMNG  // Multiple-image Network Graphics
	ImageJNG  // JPEG Network Graphics
	ImageMPO  // Multi Picture Object (multi-frame JPEG)
	ImageDPX  // Digital Picture Exchange
	ImageFITS // FITS (astronomy)
	ImageDCM  // DICOM (medical imaging container)
	ImageFPX  // FlashPix
	ImageDJVU // DjVu
	ImagePCX  // PC Paintbrush
	ImageWPG  // WordPerfect Graphics
	ImagePICT // Apple PICT
	ImagePCD  // Photo CD
	ImageBPG  // Better Portable Graphics
	ImageFLIF // Free Lossless Image Format
	ImagePGF  // Progressive Graphics File
	ImageXCF  // GIMP native
	ImageQTIF // QuickTime Image File (QTIF)

	// --- added RAWs ExifTool supports ---
	ImageNRW // Nikon NRW
	ImageSR2 // Sony SR2
	ImageSRF // Sony SRF
	ImageFFF // Hasselblad FFF
	ImageMOS // Leaf MOS
	ImageK25 // Kodak K25

	// NOTE: FileType values are serialized (msgp) and must stay stable:
	// always append new types here, never insert or reorder.
	ImageJ2C  // JPEG 2000 codestream (SOC+SIZ), distinct from the JP2 container
	ImageXISF // XISF (astronomy, "XISF0100" signature)
	ImageMEF  // Mamiya MEF
)

// MediaType groups file types into high-level semantic classes.
type MediaType uint8

const (
	MediaTypeUnknown MediaType = iota
	MediaTypeRaster
	MediaTypeVector
	MediaTypeMetadata
	MediaTypeRaw
)

func (m MediaType) String() string {
	switch m {
	case MediaTypeRaster:
		return "raster"
	case MediaTypeVector:
		return "vector"
	case MediaTypeMetadata:
		return "metadata"
	case MediaTypeRaw:
		return "raw"
	default:
		return "unknown"
	}
}

// BaseType represents the container/signature class used for classification.
type BaseType uint8

const (
	BaseTypeUnknown BaseType = iota
	BaseTypeJPEG
	BaseTypePNG
	BaseTypeGIF
	BaseTypeBMP
	BaseTypeRIFF
	BaseTypeISOBMFF
	BaseTypeTIFF
	BaseTypeCIFF
	BaseTypePSD
	BaseTypeXML
	BaseTypeNetpbm
	BaseTypeJP2
	BaseTypeSVG
	BaseTypeMagick
	BaseTypeICO
	BaseTypeTGA
	BaseTypeDDS
	BaseTypeEXR
	BaseTypeHDR
	BaseTypeJXL
	BaseTypeJXR
	BaseTypeMNG
	BaseTypeJNG
	BaseTypeMPO
	BaseTypeDPX
	BaseTypeFITS
	BaseTypeDICOM
	BaseTypeFPX
	BaseTypeDJVU
	BaseTypePCX
	BaseTypeWPG
	BaseTypePICT
	BaseTypePCD
	BaseTypeBPG
	BaseTypeFLIF
	BaseTypePGF
	BaseTypeXCF
	BaseTypeQTIF
	BaseTypeXISF
)

func (b BaseType) String() string {
	switch b {
	case BaseTypeJPEG:
		return "jpeg"
	case BaseTypePNG:
		return "png"
	case BaseTypeGIF:
		return "gif"
	case BaseTypeBMP:
		return "bmp"
	case BaseTypeRIFF:
		return "riff"
	case BaseTypeISOBMFF:
		return "isobmff"
	case BaseTypeTIFF:
		return "tiff"
	case BaseTypeCIFF:
		return "ciff"
	case BaseTypePSD:
		return "psd"
	case BaseTypeXML:
		return "xml"
	case BaseTypeNetpbm:
		return "netpbm"
	case BaseTypeJP2:
		return "jp2"
	case BaseTypeSVG:
		return "svg"
	case BaseTypeMagick:
		return "magick"
	case BaseTypeICO:
		return "ico"
	case BaseTypeTGA:
		return "tga"
	case BaseTypeDDS:
		return "dds"
	case BaseTypeEXR:
		return "exr"
	case BaseTypeHDR:
		return "hdr"
	case BaseTypeJXL:
		return "jxl"
	case BaseTypeJXR:
		return "jxr"
	case BaseTypeMNG:
		return "mng"
	case BaseTypeJNG:
		return "jng"
	case BaseTypeMPO:
		return "mpo"
	case BaseTypeDPX:
		return "dpx"
	case BaseTypeFITS:
		return "fits"
	case BaseTypeDICOM:
		return "dicom"
	case BaseTypeFPX:
		return "fpx"
	case BaseTypeDJVU:
		return "djvu"
	case BaseTypePCX:
		return "pcx"
	case BaseTypeWPG:
		return "wpg"
	case BaseTypePICT:
		return "pict"
	case BaseTypePCD:
		return "pcd"
	case BaseTypeBPG:
		return "bpg"
	case BaseTypeFLIF:
		return "flif"
	case BaseTypePGF:
		return "pgf"
	case BaseTypeXCF:
		return "xcf"
	case BaseTypeQTIF:
		return "qtif"
	case BaseTypeXISF:
		return "xisf"
	default:
		return "unknown"
	}
}

var fileTypeCanonicalMIME = map[FileType]MIMEType{
	ImageUnknown: "application/octet-stream",
	ImageJPEG:    "image/jpeg",
	ImagePNG:     "image/png",
	ImageGIF:     "image/gif",
	ImageBMP:     "image/bmp",
	ImageWebP:    "image/webp",
	ImageHEIF:    "image/heif",
	ImageRAW:     "image/raw",
	ImageTiff:    "image/tiff",
	ImageDNG:     "image/x-adobe-dng",
	ImageNEF:     "image/x-nikon-nef",
	ImagePanaRAW: "image/x-panasonic-raw",
	ImageARW:     "image/x-sony-arw",
	ImageCRW:     "image/x-canon-crw",
	ImageGPR:     "image/x-gopro-gpr",
	ImageCR3:     "image/x-canon-cr3",
	ImageCR2:     "image/x-canon-cr2",
	ImagePSD:     "image/vnd.adobe.photoshop",
	ImageXMP:     "application/rdf+xml",
	ImageAVIF:    "image/avif",
	ImagePPM:     "image/x-portable-pixmap",
	ImageJP2K:    "image/jp2",
	ImageSVG:     "image/svg+xml",
	ImageMAGICK:  "image/magick",
	ImageICO:     "image/vnd.microsoft.icon",
	ImageCUR:     "image/x-cursor",
	ImageTGA:     "image/x-tga",
	ImageDDS:     "image/vnd-ms.dds",
	ImageEXR:     "image/x-exr",
	ImageHDR:     "image/vnd.radiance",
	ImageJXL:     "image/jxl",
	ImageHEIC:    "image/heic",
	ImageAPNG:    "image/apng",
	ImagePBM:     "image/x-portable-bitmap",
	ImagePGM:     "image/x-portable-graymap",
	ImagePNM:     "image/x-portable-anymap",
	ImagePAM:     "image/x-portable-arbitrarymap",
	ImageRAF:     "image/x-fuji-raf",
	ImageORF:     "image/x-olympus-orf",
	ImageSRW:     "image/x-samsung-srw",
	ImagePEF:     "image/x-pentax-pef",
	ImageRWL:     "image/x-leica-rwl",
	ImageIIQ:     "image/x-phaseone-iiq",
	Image3FR:     "image/x-hasselblad-3fr",
	ImageX3F:     "image/x-sigma-x3f",
	ImageMRW:     "image/x-minolta-mrw",
	ImageKDC:     "image/x-kodak-kdc",
	ImageDCR:     "image/x-kodak-dcr",
	ImageERF:     "image/x-epson-erf",
	ImageJXR:     "image/vnd.ms-photo",
	ImageMNG:     "image/x-mng",
	ImageJNG:     "image/x-jng",
	ImageMPO:     "image/mpo",
	ImageDPX:     "image/x-dpx",
	ImageFITS:    "image/fits",
	ImageDCM:     "application/dicom",
	ImageFPX:     "image/vnd.fpx",
	ImageDJVU:    "image/vnd.djvu",
	ImagePCX:     "image/x-pcx",
	ImageWPG:     "application/x-wpg",
	ImagePICT:    "image/x-pict",
	ImagePCD:     "image/x-photo-cd",
	ImageBPG:     "image/bpg",
	ImageFLIF:    "image/flif",
	ImagePGF:     "image/pgf",
	ImageXCF:     "image/x-xcf",
	ImageQTIF:    "image/qtif",
	ImageNRW:     "image/x-nikon-nrw",
	ImageSR2:     "image/x-sony-sr2",
	ImageSRF:     "image/x-sony-srf",
	ImageFFF:     "image/x-hasselblad-fff",
	ImageMOS:     "image/x-leaf-mos",
	ImageK25:     "image/x-kodak-k25",
	ImageJ2C:     "image/j2c",
	ImageXISF:    "image/xisf",
	ImageMEF:     "image/x-mamiya-mef",
}

var fileTypeCanonicalExtension = map[FileType]FileTypeExtension{
	ImageUnknown: "",
	ImageJPEG:    "jpg",
	ImagePNG:     "png",
	ImageGIF:     "gif",
	ImageBMP:     "bmp",
	ImageWebP:    "webp",
	ImageHEIF:    "heif",
	ImageRAW:     "raw",
	ImageTiff:    "tiff",
	ImageDNG:     "dng",
	ImageNEF:     "nef",
	ImagePanaRAW: "rw2",
	ImageARW:     "arw",
	ImageCRW:     "crw",
	ImageGPR:     "gpr",
	ImageCR3:     "cr3",
	ImageCR2:     "cr2",
	ImagePSD:     "psd",
	ImageXMP:     "xmp",
	ImageAVIF:    "avif",
	ImagePPM:     "ppm",
	ImageJP2K:    "jp2",
	ImageSVG:     "svg",
	ImageMAGICK:  "magick",
	ImageICO:     "ico",
	ImageCUR:     "cur",
	ImageTGA:     "tga",
	ImageDDS:     "dds",
	ImageEXR:     "exr",
	ImageHDR:     "hdr",
	ImageJXL:     "jxl",
	ImageHEIC:    "heic",
	ImageAPNG:    "apng",
	ImagePBM:     "pbm",
	ImagePGM:     "pgm",
	ImagePNM:     "pnm",
	ImagePAM:     "pam",
	ImageRAF:     "raf",
	ImageORF:     "orf",
	ImageSRW:     "srw",
	ImagePEF:     "pef",
	ImageRWL:     "rwl",
	ImageIIQ:     "iiq",
	Image3FR:     "3fr",
	ImageX3F:     "x3f",
	ImageMRW:     "mrw",
	ImageKDC:     "kdc",
	ImageDCR:     "dcr",
	ImageERF:     "erf",
	ImageJXR:     "jxr",
	ImageMNG:     "mng",
	ImageJNG:     "jng",
	ImageMPO:     "mpo",
	ImageDPX:     "dpx",
	ImageFITS:    "fits",
	ImageDCM:     "dcm",
	ImageFPX:     "fpx",
	ImageDJVU:    "djvu",
	ImagePCX:     "pcx",
	ImageWPG:     "wpg",
	ImagePICT:    "pict",
	ImagePCD:     "pcd",
	ImageBPG:     "bpg",
	ImageFLIF:    "flif",
	ImagePGF:     "pgf",
	ImageXCF:     "xcf",
	ImageQTIF:    "qtif",
	ImageNRW:     "nrw",
	ImageSR2:     "sr2",
	ImageSRF:     "srf",
	ImageFFF:     "fff",
	ImageMOS:     "mos",
	ImageK25:     "k25",
	ImageJ2C:     "j2c",
	ImageXISF:    "xisf",
	ImageMEF:     "mef",
}

// mimeTypeValues maps a content-type string with a file type.
var mimeTypeValues = map[MIMEType]FileType{
	"application/dicom":             ImageDCM,
	"application/dicom+json":        ImageDCM, // sometimes used; still DICOM family
	"application/dicom+xml":         ImageDCM,
	"application/fits":              ImageFITS,
	"application/octet-stream":      ImageUnknown,
	"application/rdf+xml":           ImageXMP,
	"application/x-pcx":             ImagePCX,
	"application/x-wpg":             ImageWPG,
	"application/x-xcf":             ImageXCF,
	"image/aces":                    ImageEXR,  // uncommon; many servers just octet-stream
	"image/apng":                    ImageAPNG, // official-ish; often served as image/png
	"image/avif":                    ImageAVIF,
	"image/bmp":                     ImageBMP,
	"image/bpg":                     ImageBPG,
	"image/fits":                    ImageFITS,
	"image/flif":                    ImageFLIF,
	"image/gif":                     ImageGIF,
	"image/heic":                    ImageHEIC, // commonly used; HEIC is part of HEIF
	"image/heic-sequence":           ImageHEIC,
	"image/heif":                    ImageHEIF,
	"image/heif-sequence":           ImageHEIF,
	"image/jpeg":                    ImageJPEG,
	"image/jp2":                     ImageJP2K,
	"image/jxl":                     ImageJXL,
	"image/jxr":                     ImageJXR, // uncommon, but seen
	"image/j2c":                     ImageJ2C,
	"image/xisf":                    ImageXISF,
	"image/x-mamiya-mef":            ImageMEF,
	"image/magick":                  ImageMAGICK,
	"image/mpo":                     ImageMPO,
	"image/pgf":                     ImagePGF,
	"image/png":                     ImagePNG,
	"image/qtif":                    ImageQTIF,
	"image/raw":                     ImageRAW,
	"image/svg+xml":                 ImageSVG,
	"image/tiff":                    ImageTiff,
	"image/vnd-ms.dds":              ImageDDS,
	"image/vnd.adobe.photoshop":     ImagePSD,
	"image/vnd.djvu":                ImageDJVU,
	"image/vnd.fpx":                 ImageFPX,
	"image/vnd.microsoft.icon":      ImageICO,
	"image/vnd.ms-photo":            ImageJXR, // common for .wdp/.jxr
	"image/vnd.radiance":            ImageHDR,
	"image/webp":                    ImageWebP,
	"image/x-adobe-dng":             ImageDNG,
	"image/x-canon-cr2":             ImageCR2,
	"image/x-canon-cr3":             ImageCR3,
	"image/x-canon-crw":             ImageCRW,
	"image/x-cursor":                ImageCUR,
	"image/x-dds":                   ImageDDS,
	"image/x-djvu":                  ImageDJVU,
	"image/x-dpx":                   ImageDPX,
	"image/x-epson-erf":             ImageERF,
	"image/x-exr":                   ImageEXR,
	"image/x-flif":                  ImageFLIF,
	"image/x-fuji-raf":              ImageRAF,
	"image/x-gopro-gpr":             ImageGPR,
	"image/x-hasselblad-3fr":        Image3FR,
	"image/x-hasselblad-fff":        ImageFFF,
	"image/x-hdp":                   ImageJXR,
	"image/x-hdr":                   ImageHDR,
	"image/x-icon":                  ImageICO, // also used for .ico
	"image/x-jng":                   ImageJNG,
	"image/x-jxr":                   ImageJXR,
	"image/x-kodak-dcr":             ImageDCR,
	"image/x-kodak-k25":             ImageK25,
	"image/x-kodak-kdc":             ImageKDC,
	"image/x-leaf-mos":              ImageMOS,
	"image/x-leica-rwl":             ImageRWL,
	"image/x-minolta-mrw":           ImageMRW,
	"image/x-mng":                   ImageMNG,
	"image/x-mpo":                   ImageMPO,
	"image/x-nikon-nef":             ImageNEF,
	"image/x-nikon-nrw":             ImageNRW,
	"image/x-olympus-orf":           ImageORF,
	"image/x-panasonic-raw":         ImagePanaRAW,
	"image/x-pcx":                   ImagePCX,
	"image/x-pentax-pef":            ImagePEF,
	"image/x-phaseone-iiq":          ImageIIQ,
	"image/x-photo-cd":              ImagePCD,
	"image/x-pic":                   ImagePICT,
	"image/x-pict":                  ImagePICT,
	"image/x-pgf":                   ImagePGF,
	"image/x-portable-anymap":       ImagePNM,
	"image/x-portable-arbitrarymap": ImagePAM,
	"image/x-portable-bitmap":       ImagePBM,
	"image/x-portable-graymap":      ImagePGM,
	"image/x-portable-pixmap":       ImagePPM,
	"image/x-qtif":                  ImageQTIF,
	"image/x-samsung-srw":           ImageSRW,
	"image/x-sigma-x3f":             ImageX3F,
	"image/x-sony-arw":              ImageARW,
	"image/x-sony-sr2":              ImageSR2,
	"image/x-sony-srf":              ImageSRF,
	"image/x-targa":                 ImageTGA,
	"image/x-tga":                   ImageTGA,
	"image/x-wdp":                   ImageJXR,
	"image/x-win-bitmap":            ImageICO, // seen occasionally
	"image/x-xcf":                   ImageXCF,
	"video/x-mng":                   ImageMNG, // often mislabeled as video/*
}

// fileTypeExtensions maps filename extensions with a file type.
var fileTypeExtensions = map[FileTypeExtension]FileType{
	"":        ImageUnknown,
	".3fr":    Image3FR,
	".apng":   ImageAPNG, // if you want to distinguish; otherwise map to ImagePNG
	".arw":    ImageARW,
	".avif":   ImageAVIF,
	".bmp":    ImageBMP,
	".bpg":    ImageBPG,
	".cr2":    ImageCR2,
	".cr3":    ImageCR3,
	".crw":    ImageCRW,
	".cur":    ImageCUR,
	".dcm":    ImageDCM,
	".dcr":    ImageDCR,
	".dds":    ImageDDS,
	".djv":    ImageDJVU,
	".djvu":   ImageDJVU,
	".dng":    ImageDNG,
	".dpx":    ImageDPX,
	".erf":    ImageERF,
	".exr":    ImageEXR,
	".fff":    ImageFFF,
	".fit":    ImageFITS,
	".fits":   ImageFITS,
	".flif":   ImageFLIF,
	".fpx":    ImageFPX,
	".fts":    ImageFITS,
	".gif":    ImageGIF,
	".gpr":    ImageGPR,
	".hdp":    ImageJXR,
	".hdr":    ImageHDR,
	".heic":   ImageHEIC,
	".heics":  ImageHEIC,
	".heif":   ImageHEIF,
	".heifs":  ImageHEIF,
	".ico":    ImageICO,
	".iiq":    ImageIIQ,
	".j2k":    ImageJP2K,
	".jfif":   ImageJPEG,
	".jng":    ImageJNG,
	".jpe":    ImageJPEG,
	".jpeg":   ImageJPEG,
	".jp2":    ImageJP2K,
	".jpg":    ImageJPEG,
	".jpm":    ImageJP2K,
	".jpx":    ImageJP2K,
	".jxl":    ImageJXL,
	".jxr":    ImageJXR,
	".j2c":    ImageJ2C,
	".xisf":   ImageXISF,
	".mef":    ImageMEF,
	".k25":    ImageK25,
	".kdc":    ImageKDC,
	".magick": ImageMAGICK,
	".mng":    ImageMNG,
	".mos":    ImageMOS,
	".mpo":    ImageMPO,
	".mrw":    ImageMRW,
	".nef":    ImageNEF,
	".nrw":    ImageNRW,
	".orf":    ImageORF,
	".pam":    ImagePAM,
	".pbm":    ImagePBM,
	".pcd":    ImagePCD,
	".pct":    ImagePICT,
	".pcx":    ImagePCX,
	".pef":    ImagePEF,
	".pgf":    ImagePGF,
	".pgm":    ImagePGM,
	".pic":    ImagePICT,
	".pict":   ImagePICT,
	".pnm":    ImagePNM,
	".png":    ImagePNG,
	".ppm":    ImagePPM,
	".psd":    ImagePSD,
	".qti":    ImageQTIF,
	".qtif":   ImageQTIF,
	".raf":    ImageRAF,
	".raw":    ImageRAW,
	".rw2":    ImagePanaRAW,
	".rwl":    ImageRWL,
	".sr2":    ImageSR2,
	".srf":    ImageSRF,
	".srw":    ImageSRW,
	".svg":    ImageSVG,
	".svgz":   ImageSVG, // gzipped svg
	".tga":    ImageTGA,
	".tif":    ImageTiff,
	".tiff":   ImageTiff,
	".wdp":    ImageJXR,
	".webp":   ImageWebP,
	".wpg":    ImageWPG,
	".x3f":    ImageX3F,
	".xcf":    ImageXCF,
	".xmp":    ImageXMP,
}

var (
	tiffLittleEndianSignature = []byte{0x49, 0x49, 0x2A, 0x00}
	tiffBigEndianSignature    = []byte{0x4D, 0x4D, 0x00, 0x2A}

	// BigTIFF (magic 43) signatures, allowed for DNG since spec 1.7.
	bigTiffLittleEndianSignature = []byte{0x49, 0x49, 0x2B, 0x00}
	bigTiffBigEndianSignature    = []byte{0x4D, 0x4D, 0x00, 0x2B}

	crwByteOrderSignature = []byte{0x49, 0x49}
	crwHeapSignature      = []byte("HEAPCCDR")
	cr2Signature          = []byte{0x43, 0x52, 0x02, 0x00}

	ftypBoxType = []byte("ftyp")
	brandCRX    = []byte("crx ")
	brandHEIC   = []byte("heic")
	brandHEIF   = []byte("heif")
	brandHEIM   = []byte("heim")
	brandHEIS   = []byte("heis")
	brandHEIX   = []byte("heix")
	brandHEVM   = []byte("hevm")
	brandHEVS   = []byte("hevs")
	brandHEVX   = []byte("hevx")
	brandMIAF   = []byte("miaf")
	brandMIF1   = []byte("mif1")
	brandMSF1   = []byte("msf1")
	brandHEVC   = []byte("hevc")
	brandAVIF   = []byte("avif")
	brandAVIS   = []byte("avis")
	brandJXL    = []byte("jxl ")

	bmpSignature         = []byte("BM")
	icoSignature         = []byte{0x00, 0x00, 0x01, 0x00}
	curSignature         = []byte{0x00, 0x00, 0x02, 0x00}
	ddsSignature         = []byte("DDS ")
	exrSignature         = []byte{0x76, 0x2F, 0x31, 0x01}
	dpxBigSignature      = []byte("SDPX")
	dpxLittleSignature   = []byte("XPDS")
	mngSignature         = []byte{0x8A, 0x4D, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	jngSignature         = []byte{0x8B, 0x4A, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	fitsSignature        = []byte("SIMPLE  =")
	rafSignature         = []byte("FUJIFILMCCD-RAW ")
	xcfSignature         = []byte("gimp xcf ")
	flifSignature        = []byte("FLIF")
	bpgSignature         = []byte{0x42, 0x50, 0x47, 0xFB}
	hdrRadianceSignature = []byte("#?RADIANCE")
	hdrRGBESignature     = []byte("#?RGBE")
	djvuFormSignature    = []byte("AT&TFORM")
	djvuTypeDJVU         = []byte("DJVU")
	djvuTypeDJVM         = []byte("DJVM")
	djvuTypeDJVI         = []byte("DJVI")
	rw2TiffSignature     = []byte{0x49, 0x49, 0x55, 0x00}
	rw2RawSignature      = []byte{0x88, 0xE7, 0x74, 0xD8}
	jpegSignature        = []byte{0xFF, 0xD8}
	pngSignature         = []byte{0x89, 0x50, 0x4E, 0x47}
	riffSignature        = []byte("RIFF")
	webpSignature        = []byte("WEBP")
	jpeg2000Signature    = []byte{0x00, 0x00, 0x00, 0x0C, 0x6A, 0x50, 0x20, 0x20, 0x0D, 0x0A, 0x87, 0x0A}

	jxlCodestreamSignature = []byte{0xFF, 0x0A}
	jxlContainerSignature  = []byte{0x00, 0x00, 0x00, 0x0C, 0x4A, 0x58, 0x4C, 0x20, 0x0D, 0x0A, 0x87, 0x0A}

	psdSignature     = []byte("8BPS")
	utf8BOMSignature = []byte{0xEF, 0xBB, 0xBF}
	svgTagSignature  = []byte("<svg")
	svgXMLSignature  = []byte("<?xml")
	svgDocTypeSig    = []byte("<!doctype")
	svgCommentSig    = []byte("<!--")
	xmpSignature     = []byte("<x:xmpmeta")
	gif87aSignature  = []byte("GIF87a")
	gif89aSignature  = []byte("GIF89a")

	pgfSignature  = []byte("PGF")
	wpgSignature  = []byte{0xFF, 0x57, 0x50, 0x43}
	j2cSignature  = []byte{0xFF, 0x4F, 0xFF, 0x51}
	xisfSignature = []byte("XISF0100")
)

func hasPrefix(buf, sig []byte) bool {
	return len(buf) >= len(sig) && bytes.Equal(buf[:len(sig)], sig)
}

func hasPrefixFold(buf, sig []byte) bool {
	return len(buf) >= len(sig) && bytes.EqualFold(buf[:len(sig)], sig)
}

func indexFold(buf, token []byte) int {
	if len(token) == 0 || len(buf) < len(token) {
		return -1
	}

	limit := len(buf) - len(token)
	for i := 0; i <= limit; i++ {
		if bytes.EqualFold(buf[i:i+len(token)], token) {
			return i
		}
	}

	return -1
}

func hasAt(buf []byte, offset int, sig []byte) bool {
	return offset >= 0 && len(buf) >= offset+len(sig) && bytes.Equal(buf[offset:offset+len(sig)], sig)
}

func hasCompatibleBrand(buf []byte, brand []byte) bool {
	// Walk every compatible-brand slot in the ftyp box instead of only
	// the first two: files with long brand lists otherwise miss.
	if len(buf) < 8 {
		return false
	}
	size := uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	limit := len(buf)
	// A zero box size means "to end of file" per the spec.
	if size != 0 && int(size) < limit {
		limit = int(size)
	}
	for off := 16; off+4 <= limit; off += 4 {
		if hasAt(buf, off, brand) {
			return true
		}
	}
	return false
}

func hasAnyCompatibleBrand(buf []byte, brands ...[]byte) bool {
	for _, brand := range brands {
		if hasCompatibleBrand(buf, brand) {
			return true
		}
	}
	return false
}

// isTiff() Checks to see if an Image has the tiff format header (classic or
// BigTIFF).
func isTiff(buf []byte) bool {
	return IsTiffBigEndian(buf) || IsTiffLittleEndian(buf) || isBigTiff(buf)
}

// isBigTiff reports either BigTIFF (magic 43) byte-order signature.
func isBigTiff(buf []byte) bool {
	return hasPrefix(buf, bigTiffLittleEndianSignature) || hasPrefix(buf, bigTiffBigEndianSignature)
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func IsTiffLittleEndian(buf []byte) bool {
	return hasPrefix(buf, tiffLittleEndianSignature)
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func IsTiffBigEndian(buf []byte) bool {
	return hasPrefix(buf, tiffBigEndianSignature)
}

// isCRW returns true if it matches an image/x-canon-crw with 14 bytes of the header.
//
// CanonCRWHeader is the file Header for a Canon CRW file. Currently only Little Endian support
// Reference: https://exiftool.org/canon_raw.html
func isCRW(buf []byte) bool {
	// ByteOrder: LittleEndian + Signature: HEAPCCDR
	return hasPrefix(buf, crwByteOrderSignature) &&
		hasAt(buf, 6, crwHeapSignature)
}

// isCR2 returns true if it matches an image/x-canon-cr2.
//
// CanonCR2Header is the Header for a Canon CR2 file
// 4 bytes after TiffSignature and before the beginning of IFDO
func isCR2(buf []byte) bool {
	return isTiff(buf) && hasAt(buf, 8, cr2Signature)
}

// isCR3 returns true if it matches an image/x-canon-cr3.
//
// ftyp box with major_brand: 'crx ' and compatible_brands: 'crx ' 'isom'
func isCR3(buf []byte) bool {
	return isFTYPBox(buf) && isFTYPBrand(buf[8:], brandCRX)
}

// isobmffSubtype returns a specific ISOBMFF-based file type where possible.
func isobmffSubtype(buf []byte) FileType {
	if !isFTYPBox(buf) {
		return ImageUnknown
	}

	if isCR3(buf) {
		return ImageCR3
	}

	if isAVIF(buf) {
		return ImageAVIF
	}

	// JPEG XL brand fallback for containers without the leading 'JXL ' signature box.
	if isFTYPBrand(buf[8:], brandJXL) ||
		hasAnyCompatibleBrand(buf, brandJXL) {
		return ImageJXL
	}

	// HEIC-branded variants.
	if isFTYPBrand(buf[8:], brandHEIC) ||
		isFTYPBrand(buf[8:], brandHEIM) ||
		isFTYPBrand(buf[8:], brandHEIS) ||
		isFTYPBrand(buf[8:], brandHEIX) ||
		isFTYPBrand(buf[8:], brandHEVC) ||
		isFTYPBrand(buf[8:], brandHEVM) ||
		isFTYPBrand(buf[8:], brandHEVS) ||
		isFTYPBrand(buf[8:], brandHEVX) ||
		hasAnyCompatibleBrand(buf, brandHEIC, brandHEIM, brandHEIS, brandHEIX, brandHEVC, brandHEVM, brandHEVS, brandHEVX) {
		return ImageHEIC
	}

	// Generic HEIF (mif1/msf1/heif without explicit HEIC branding).
	if isFTYPBrand(buf[8:], brandHEIF) ||
		isFTYPBrand(buf[8:], brandMIAF) ||
		isFTYPBrand(buf[8:], brandMIF1) ||
		isFTYPBrand(buf[8:], brandMSF1) ||
		hasAnyCompatibleBrand(buf, brandHEIF, brandMIAF, brandMIF1, brandMSF1) {
		return ImageHEIF
	}

	if isHeif(buf) {
		return ImageHEIF
	}

	return ImageUnknown
}

// isHeif returns true if the header matches the start of a HEIF file.
//
// Major brands: heic/heix (+ variants), mif1/msf1, heif/miaf
// Minor brand contains:
func isHeif(buf []byte) bool {
	if !isFTYPBox(buf) {
		return false
	}

	if isFTYPBrand(buf[8:], brandHEIC) ||
		isFTYPBrand(buf[8:], brandHEIM) ||
		isFTYPBrand(buf[8:], brandHEIS) ||
		isFTYPBrand(buf[8:], brandHEIX) ||
		isFTYPBrand(buf[8:], brandHEVC) ||
		isFTYPBrand(buf[8:], brandHEVM) ||
		isFTYPBrand(buf[8:], brandHEVS) ||
		isFTYPBrand(buf[8:], brandHEVX) ||
		isFTYPBrand(buf[8:], brandHEIF) ||
		isFTYPBrand(buf[8:], brandMIAF) {
		return true
	}

	return (isFTYPBrand(buf[8:], brandMIF1) && hasAnyCompatibleBrand(buf, brandHEIC, brandHEIM, brandHEIS, brandHEIX, brandHEVC, brandHEVM, brandHEVS, brandHEVX, brandHEIF, brandMIAF)) ||
		(isFTYPBrand(buf[8:], brandMSF1) && hasAnyCompatibleBrand(buf, brandHEIC, brandHEIM, brandHEIS, brandHEIX, brandHEVC, brandHEVM, brandHEVS, brandHEVX, brandHEIF, brandMIAF))
}

// isFTYPBrand returns true if the Brand in []byte matches the provided brand.
// the Limit is 4 bytes
func isFTYPBrand(buf []byte, brand []byte) bool {
	return hasPrefix(buf, brand)
}

// isFTYPBox returns true if the header matches an ftyp box.
// This indicates an ISO Base Media File Format.
func isFTYPBox(buf []byte) bool {
	// buf[0:4] is 'ftyp' box size
	return len(buf) >= 8 &&
		buf[0] == 0x00 &&
		buf[1] == 0x00 &&
		hasAt(buf, 4, ftypBoxType)
}

// isAVIF returns true if the header matches an ftyp box and
// an avif box.
func isAVIF(buf []byte) bool {
	return isFTYPBox(buf) &&
		(isFTYPBrand(buf[8:], brandAVIF) ||
			isFTYPBrand(buf[8:], brandAVIS) ||
			((isFTYPBrand(buf[8:], brandMIF1) || isFTYPBrand(buf[8:], brandMSF1)) &&
				hasAnyCompatibleBrand(buf, brandAVIF, brandAVIS)))
}

// isBMP returns true if the header matches the start of a BMP file
// Bitmap Image
func isBMP(buf []byte) bool {
	return hasPrefix(buf, bmpSignature)
}

func isICO(buf []byte) bool {
	return hasPrefix(buf, icoSignature)
}

func isCUR(buf []byte) bool {
	return hasPrefix(buf, curSignature)
}

func isDDS(buf []byte) bool {
	return hasPrefix(buf, ddsSignature)
}

func isEXR(buf []byte) bool {
	return hasPrefix(buf, exrSignature)
}

func isDPX(buf []byte) bool {
	return hasPrefix(buf, dpxBigSignature) || hasPrefix(buf, dpxLittleSignature)
}

func isMNG(buf []byte) bool {
	return hasPrefix(buf, mngSignature)
}

func isJNG(buf []byte) bool {
	return hasPrefix(buf, jngSignature)
}

func isFITS(buf []byte) bool {
	return hasPrefix(buf, fitsSignature)
}

func isRAF(buf []byte) bool {
	return hasPrefix(buf, rafSignature)
}

func isXCF(buf []byte) bool {
	return hasPrefix(buf, xcfSignature)
}

func isFLIF(buf []byte) bool {
	return hasPrefix(buf, flifSignature)
}

func isBPG(buf []byte) bool {
	return hasPrefix(buf, bpgSignature)
}

func isHDR(buf []byte) bool {
	return hasPrefix(buf, hdrRadianceSignature) || hasPrefix(buf, hdrRGBESignature)
}

func isDJVU(buf []byte) bool {
	return hasPrefix(buf, djvuFormSignature) &&
		(hasAt(buf, 12, djvuTypeDJVU) || hasAt(buf, 12, djvuTypeDJVM) || hasAt(buf, 12, djvuTypeDJVI))
}

// isRW2 returns true if the first 4 bytes match the Panasonic Tiff alternate
// header and bytes 8 through 12 match the RW2 header
func isRW2(buf []byte) bool {
	return hasPrefix(buf, rw2TiffSignature) && hasAt(buf, 8, rw2RawSignature)
}

func tiffReadU16(buf []byte, offset int, littleEndian bool) (uint16, bool) {
	if offset < 0 || len(buf) < offset+2 {
		return 0, false
	}
	if littleEndian {
		return uint16(buf[offset]) | (uint16(buf[offset+1]) << 8), true
	}
	return (uint16(buf[offset]) << 8) | uint16(buf[offset+1]), true
}

func tiffReadU32(buf []byte, offset int, littleEndian bool) (uint32, bool) {
	if offset < 0 || len(buf) < offset+4 {
		return 0, false
	}
	if littleEndian {
		return uint32(buf[offset]) |
			(uint32(buf[offset+1]) << 8) |
			(uint32(buf[offset+2]) << 16) |
			(uint32(buf[offset+3]) << 24), true
	}
	return (uint32(buf[offset]) << 24) |
		(uint32(buf[offset+1]) << 16) |
		(uint32(buf[offset+2]) << 8) |
		uint32(buf[offset+3]), true
}

// tiffMakeModel scans the first IFD of a TIFF or BigTIFF image in buf for
// Make (0x010F), Model (0x0110) and DNGVersion (0xC612) tags. It returns
// ok=false when the structure is truncated or malformed. All reads are
// bounds-checked and iterative: no recursion, no unbounded loops.
func tiffMakeModel(buf []byte) (manufacturer, model string, hasDNGVersion, ok bool) {
	if len(buf) < 8 {
		return "", "", false, false
	}
	big := isBigTiff(buf)
	var littleEndian bool
	switch {
	case IsTiffLittleEndian(buf) || (big && buf[0] == 0x49):
		littleEndian = true
	case IsTiffBigEndian(buf) || (big && buf[0] == 0x4D):
		littleEndian = false
	default:
		return "", "", false, false
	}

	var ifdOffset uint64
	var entrySize, countSize int
	if big {
		// BigTIFF: byte order, magic 43, offset byte size, 0, 8-byte IFD offset.
		if len(buf) < 16 {
			return "", "", false, false
		}
		off, readOK := tiffReadU64(buf, 8, littleEndian)
		if !readOK {
			return "", "", false, false
		}
		ifdOffset, entrySize, countSize = off, 20, 8
	} else {
		off, readOK := tiffReadU32(buf, 4, littleEndian)
		if !readOK {
			return "", "", false, false
		}
		ifdOffset, entrySize, countSize = uint64(off), 12, 2
	}

	count, ok := tiffReadCount(buf, ifdOffset, countSize, littleEndian)
	if !ok || count == 0 || count > maxTiffEntries {
		return "", "", false, false
	}
	entriesAt := ifdOffset + uint64(countSize)
	for i := uint64(0); i < count; i++ {
		base := entriesAt + i*uint64(entrySize)
		tag, ok := tiffReadU16At(buf, base, littleEndian)
		if !ok {
			return "", "", false, false
		}
		switch tag {
		case tiffTagMake, tiffTagModel, tiffTagDNGVersion:
		default:
			continue
		}
		typ, ok := tiffReadU16At(buf, base+2, littleEndian)
		if !ok {
			return "", "", false, false
		}
		if tag == tiffTagDNGVersion {
			hasDNGVersion = true
			continue
		}
		if typ != tiffTypeASCII {
			continue
		}
		var elen uint64
		if big {
			elen, ok = tiffReadU64(buf, base+4, littleEndian)
		} else {
			var elen32 uint32
			//nolint:gosec // G115: base+4 <= len(buf) is proven by the tag/type/count reads above.
			elen32, ok = tiffReadU32(buf, int(base)+4, littleEndian)
			elen = uint64(elen32)
		}
		if !ok || elen == 0 || elen > maxTiffString {
			continue
		}
		// Inline values live in the entry tail (4 bytes classic, 8
		// bytes BigTIFF); longer ones sit at an offset.
		var valAt uint64
		var inlineCap uint64
		if big {
			valAt, inlineCap = base+12, 8
		} else {
			valAt, inlineCap = base+8, 4
		}
		var s string
		if elen <= inlineCap {
			s, ok = tiffASCII(buf, valAt, elen)
		} else {
			var off uint64
			if big {
				off, ok = tiffReadU64(buf, valAt, littleEndian)
			} else {
				var off32 uint32
				//nolint:gosec // G115: valAt = base+8 with base proven readable above; tiffReadU32 re-checks bounds.
				off32, ok = tiffReadU32(buf, int(valAt), littleEndian)
				off = uint64(off32)
			}
			if !ok {
				continue
			}
			s, ok = tiffASCII(buf, off, elen)
		}
		if !ok || s == "" {
			continue
		}
		if tag == tiffTagMake {
			manufacturer = s
		} else {
			model = s
		}
		if manufacturer != "" && model != "" && hasDNGVersion {
			break
		}
	}
	return manufacturer, model, hasDNGVersion, true
}

// TIFF IFD constants for subtype detection.
const (
	tiffTagMake       = 0x010F
	tiffTagModel      = 0x0110
	tiffTagDNGVersion = 0xC612
	tiffTypeASCII     = 2

	// maxTiffEntries caps the first-IFD walk; real IFDs stay far below it
	// and garbage counts must not cause long loops.
	maxTiffEntries = 512
	// maxTiffString caps Make/Model value reads.
	maxTiffString = 64
)

// tiffReadCount reads an entry count of countSize bytes at offset.
func tiffReadCount(buf []byte, offset uint64, countSize int, littleEndian bool) (uint64, bool) {
	if countSize == 8 {
		return tiffReadU64(buf, offset, littleEndian)
	}
	v, ok := tiffReadU16At(buf, offset, littleEndian)
	return uint64(v), ok
}

// tiffReadU16At reads a uint16 at an absolute uint64 offset.
func tiffReadU16At(buf []byte, offset uint64, littleEndian bool) (uint16, bool) {
	if offset > uint64(len(buf))-2 {
		return 0, false
	}
	//nolint:gosec // G115: offset+2 <= len(buf) <= maxInt per the check above.
	return tiffReadU16(buf, int(offset), littleEndian)
}

// tiffReadU64 reads a uint64 at offset with bounds checking.
func tiffReadU64(buf []byte, offset uint64, littleEndian bool) (uint64, bool) {
	if offset > uint64(len(buf))-8 {
		return 0, false
	}
	//nolint:gosec // G115: offset+8 <= len(buf) <= maxInt per the check above.
	off := int(offset)
	if littleEndian {
		return uint64(buf[off]) |
			(uint64(buf[off+1]) << 8) |
			(uint64(buf[off+2]) << 16) |
			(uint64(buf[off+3]) << 24) |
			(uint64(buf[off+4]) << 32) |
			(uint64(buf[off+5]) << 40) |
			(uint64(buf[off+6]) << 48) |
			(uint64(buf[off+7]) << 56), true
	}
	return (uint64(buf[off]) << 56) |
		(uint64(buf[off+1]) << 48) |
		(uint64(buf[off+2]) << 40) |
		(uint64(buf[off+3]) << 32) |
		(uint64(buf[off+4]) << 24) |
		(uint64(buf[off+5]) << 16) |
		(uint64(buf[off+6]) << 8) |
		uint64(buf[off+7]), true
}

// tiffASCII reads an ASCII value of elen bytes at offset, stopping at the
// first NUL. The value must lie entirely inside buf.
func tiffASCII(buf []byte, offset, elen uint64) (string, bool) {
	if elen == 0 || offset > uint64(len(buf)) || elen > uint64(len(buf))-offset {
		return "", false
	}
	//nolint:gosec // G115: offset+elen <= len(buf) <= maxInt per the check above.
	raw := buf[int(offset):int(offset+elen)]
	if i := bytes.IndexByte(raw, 0); i >= 0 {
		raw = raw[:i]
	}
	return string(raw), true
}

// tiffMakeModelType maps Make/Model strings to a RAW type. hasDNGVersion
// wins for any maker. Makes without a defensible mapping return
// ImageUnknown so callers fall back to the entry-count heuristics.
func tiffMakeModelType(manufacturer, model string, hasDNGVersion bool) FileType {
	if hasDNGVersion {
		return ImageDNG
	}
	mk := strings.ToUpper(strings.Trim(manufacturer, " \x00"))
	md := strings.ToUpper(strings.Trim(model, " \x00"))
	switch {
	case strings.HasPrefix(mk, "OLYMPUS") || strings.HasPrefix(mk, "OM DIGITAL"):
		return ImageORF
	case mk == "SIGMA":
		return ImageX3F
	case mk == "PENTAX" || strings.HasPrefix(mk, "RICOH"):
		return ImagePEF
	case mk == "SAMSUNG":
		return ImageSRW
	case strings.Contains(mk, "MINOLTA"):
		return ImageMRW
	case strings.Contains(mk, "EPSON"):
		return ImageERF
	case strings.Contains(mk, "LEAF"):
		return ImageMOS
	case strings.Contains(mk, "PHASE ONE") || mk == "PHASEONE":
		return ImageIIQ
	case strings.Contains(mk, "GOPRO"):
		return ImageGPR
	case strings.Contains(mk, "MAMIYA"):
		return ImageMEF
	case strings.HasPrefix(mk, "NIKON"):
		if strings.HasPrefix(md, "COOLPIX P") || md == "COOLPIX A" {
			return ImageNRW
		}
		return ImageNEF
	case mk == "SONY":
		if md == "DSLR-A100" {
			return ImageSR2
		}
		if strings.HasPrefix(md, "DSC-") {
			return ImageSRF
		}
		return ImageARW
	case strings.Contains(mk, "KODAK"):
		if strings.Contains(md, "DC25") {
			return ImageK25
		}
		if strings.Contains(md, "DCS PRO") {
			return ImageDCR
		}
		return ImageKDC
	}
	return ImageUnknown
}

// tiffSecondarySubtype performs best-effort subtype detection for TIFF-based
// files from the 64-byte header using IFD entry-count heuristics. It stays as
// the fallback wherever Make/Model sniffing (tiffMakeModelType) finds nothing:
// the heuristics are firmware-sensitive, but they need no string pool.
func tiffSecondarySubtype(buf []byte) FileType {
	// RW2 has its own TIFF-like signature.
	if isRW2(buf) {
		return ImagePanaRAW
	}

	if !isTiff(buf) {
		return ImageUnknown
	}

	// CR2 has a fixed marker in bytes 8..11.
	if isCR2(buf) {
		return ImageCR2
	}

	littleEndian := IsTiffLittleEndian(buf)
	if !littleEndian && !IsTiffBigEndian(buf) {
		return ImageUnknown
	}

	entryCount, ok := tiffReadU16(buf, 8, littleEndian)
	if !ok {
		return ImageUnknown
	}
	firstTag, ok := tiffReadU16(buf, 10, littleEndian)
	if !ok {
		return ImageUnknown
	}
	firstType, ok := tiffReadU16(buf, 12, littleEndian)
	if !ok {
		return ImageUnknown
	}
	firstCount, ok := tiffReadU32(buf, 14, littleEndian)
	if !ok {
		return ImageUnknown
	}
	firstValue, ok := tiffReadU32(buf, 18, littleEndian)
	if !ok {
		return ImageUnknown
	}
	secondTag, ok := tiffReadU16(buf, 22, littleEndian)
	if !ok {
		return ImageUnknown
	}

	// Common first IFD structure used by many camera RAW formats.
	if firstTag != 0x00FE || firstType != 4 || firstCount != 1 {
		return ImageUnknown
	}

	// GoPro GPR samples have SubfileType=0 and high IFD entry count.
	if (entryCount == 0x0039 || entryCount == 0x003A) && firstValue == 0 && secondTag == 0x0100 {
		return ImageGPR
	}

	// DNG samples have high IFD entry count with SubfileType=1.
	if entryCount == 0x003F && firstValue == 1 && secondTag == 0x0100 {
		return ImageDNG
	}

	// Nikon NEF samples.
	if (entryCount == 0x001B || entryCount == 0x001C) && firstValue == 1 && secondTag == 0x0100 {
		return ImageNEF
	}

	// Sony ARW samples use Compression (0x0103) as second tag.
	if (entryCount == 0x0012 || entryCount == 0x0013) && firstValue == 1 && secondTag == 0x0103 {
		return ImageARW
	}

	return ImageUnknown
}

// isJPEG returns true if the first 2 bytes match a JPEG file header
//
// JPEG SOI Marker (FF D8)
func isJPEG(buf []byte) bool {
	return hasPrefix(buf, jpegSignature)
}

// isPNG returns true if the first 4 bytes match a PNG file header.
func isPNG(buf []byte) bool {
	return hasPrefix(buf, pngSignature)
}

// isWebP returns true is the first 12 bytes match a WebP file header.
// RIFF and WebP
func isWebP(buf []byte) bool {
	return hasPrefix(buf, riffSignature) && hasAt(buf, 8, webpSignature)
}

// isJPEG2000 returns true if the first 12 bytes match a JPEG2000 file header
func isJPEG2000(buf []byte) bool {
	return hasPrefix(buf, jpeg2000Signature)
}

// isJXL returns true for JPEG XL codestream or container signatures.
func isJXL(buf []byte) bool {
	return hasPrefix(buf, jxlCodestreamSignature) || hasPrefix(buf, jxlContainerSignature)
}

// isPSD returns true if the header matches a PSDImage.
//
// PSD Photoshop document
func isPSD(buf []byte) bool {
	return hasPrefix(buf, psdSignature)
}

// isXMP returns true if the header matches "<x:xmpmeta" start of a file.
//
// XMP sidecar files. The XMPHeader are the first 10bytes of an XMP sidecar.
func isXMP(buf []byte) bool {
	return hasPrefix(buf, xmpSignature)
}

// isSVG returns true if the header appears to be an SVG XML document.
func isSVG(buf []byte) bool {
	if hasPrefix(buf, utf8BOMSignature) {
		buf = buf[len(utf8BOMSignature):]
	}

	buf = bytes.TrimLeft(buf, " \t\r\n\f")
	if len(buf) == 0 {
		return false
	}

	if hasPrefixFold(buf, svgTagSignature) {
		return true
	}

	if !hasPrefixFold(buf, svgXMLSignature) &&
		!hasPrefixFold(buf, svgDocTypeSig) &&
		!hasPrefixFold(buf, svgCommentSig) {
		return false
	}

	return indexFold(buf, svgTagSignature) >= 0
}

// isGIF returns true if the header matches the header of a GIF version 87a
// or 89a.
func isGIF(buf []byte) bool {
	return hasPrefix(buf, gif87aSignature) || hasPrefix(buf, gif89aSignature)
}

// isPCX returns true for the PCX manufacturer byte and a defined version
// (0: 2.5, 2/3: 2.8, 5: 3.0).
func isPCX(buf []byte) bool {
	return len(buf) >= 2 && buf[0] == 0x0A &&
		(buf[1] == 0 || buf[1] == 2 || buf[1] == 3 || buf[1] == 5)
}

// isPGF returns true for the Progressive Graphics File magic.
func isPGF(buf []byte) bool {
	return hasPrefix(buf, pgfSignature)
}

// isWPG returns true for the WordPerfect Graphics magic ("ÿWPC").
func isWPG(buf []byte) bool {
	return hasPrefix(buf, wpgSignature)
}

// isJ2C returns true for a JPEG 2000 codestream (SOC followed by SIZ),
// distinct from the JP2 container.
func isJ2C(buf []byte) bool {
	return hasPrefix(buf, j2cSignature)
}

// isXISF returns true for the XISF 1.0 file signature.
func isXISF(buf []byte) bool {
	return hasPrefix(buf, xisfSignature)
}

func netpbmType(buf []byte) (FileType, bool) {
	if len(buf) < 3 || buf[0] != 'P' {
		return ImageUnknown, false
	}

	if buf[2] != '\n' && buf[2] != '\r' && buf[2] != '\t' && buf[2] != ' ' {
		return ImageUnknown, false
	}

	switch buf[1] {
	case '1', '4':
		return ImagePBM, true
	case '2', '5':
		return ImagePGM, true
	case '3', '6':
		return ImagePPM, true
	case '7':
		return ImagePAM, true
	default:
		return ImageUnknown, false
	}
}
