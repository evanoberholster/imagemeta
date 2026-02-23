// Package imagetype provides types and functions for identifying Image document types
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

// Extension returns the canonical file extension for the file type.
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
		ImageMOS, ImageK25:
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
		ImageSRF, ImageFFF, ImageMOS, ImageK25:
		return true
	default:
		return false
	}
}

// FromString returns a FileType for the given content-type string, extension,
// or filename.
func FromString(str string) FileType {
	str = strings.TrimSpace(str)
	if str == "" {
		return ImageUnknown
	}

	normalized := strings.ToLower(str)

	// from content-type
	if it, ok := mimeTypeValues[MIMEType(normalized)]; ok {
		return it
	}

	// from content-type with optional parameters
	if mediaType, _, err := mime.ParseMediaType(normalized); err == nil {
		if it, ok := mimeTypeValues[MIMEType(mediaType)]; ok {
			return it
		}
	}

	// from extension
	if it, ok := fileTypeExtensions[FileTypeExtension(normalized)]; ok {
		return it
	}
	if !strings.HasPrefix(normalized, ".") {
		if it, ok := fileTypeExtensions[FileTypeExtension("."+normalized)]; ok {
			return it
		}
	}

	// from file path / file name extension
	if ext := strings.ToLower(filepath.Ext(normalized)); ext != "" {
		if it, ok := fileTypeExtensions[FileTypeExtension(ext)]; ok {
			return it
		}
	}

	// from common MIME shorthand "image/jpg; charset=..."
	if idx := strings.IndexByte(normalized, ';'); idx > 0 {
		if it, ok := mimeTypeValues[MIMEType(strings.TrimSpace(normalized[:idx]))]; ok {
			return it
		}
	}

	// from common query-like suffixes: "file.jpg?foo=bar"
	if idx := strings.IndexAny(normalized, "?#"); idx > 0 {
		if ext := strings.ToLower(filepath.Ext(normalized[:idx])); ext != "" {
			if it, ok := fileTypeExtensions[FileTypeExtension(ext)]; ok {
				return it
			}
		}
	}

	// from extension token in path-like inputs without a leading dot
	if idx := strings.LastIndexByte(normalized, '.'); idx > 0 && idx < len(normalized)-1 {
		if it, ok := fileTypeExtensions[FileTypeExtension(normalized[idx:])]; ok {
			return it
		}
	}

	// from file names where ext is the full token (e.g. "jpeg")
	if it, ok := fileTypeExtensions[FileTypeExtension("."+normalized)]; ok {
		return it
	}

	return ImageUnknown
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
	return hasAt(buf, 16, brand) || hasAt(buf, 20, brand)
}

func hasAnyCompatibleBrand(buf []byte, brands ...[]byte) bool {
	for _, brand := range brands {
		if hasCompatibleBrand(buf, brand) {
			return true
		}
	}
	return false
}

// isTiff() Checks to see if an Image has the tiff format header.
func isTiff(buf []byte) bool {
	return IsTiffBigEndian(buf) || IsTiffLittleEndian(buf)
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

// tiffSecondarySubtype performs best-effort subtype detection for TIFF-based files.
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
