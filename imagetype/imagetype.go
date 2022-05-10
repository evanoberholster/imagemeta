// Package imagetype provides types and functions for identifying Image document types
package imagetype

import (
	"errors"
	"strings"
)

var (
	// ErrDataLength is an error for data length
	ErrDataLength = errors.New("error the data is not long enough")

	// ImageType stringer Index
	_ImageTypeIndex = [...]uint{0, 24, 34, 43, 52, 61, 71, 81, 90, 100, 117, 134, 155, 171, 188, 205, 222, 239, 264, 283, 293, 316, 325, 338, 350}

	// ImageType extension Index
	_ImageTypeExtIndex = [...]uint{0, 0, 3, 6, 9, 12, 16, 20, 23, 27, 30, 33, 36, 39, 42, 45, 48, 51, 54, 57, 61, 64, 67, 70, 76}
)

const (
	// ImageType stringer Names
	_ImageTypeString = "application/octet-streamimage/jpegimage/pngimage/gifimage/bmpimage/webpimage/heifimage/rawimage/tiffimage/x-adobe-dngimage/x-nikon-nefimage/x-panasonic-rawimage/x-sony-arwimage/x-canon-crwimage/x-gopro-gprimage/x-canon-cr3image/x-canon-cr2image/vnd.adobe.photoshopapplication/rdf+xmlimage/avifimage/x-portable-pixmapimage/jp2image/svg+xmlimage/magick"

	// ImageType extension Names
	_ImageTypeExtString = "jpgpnggifbmpwebpheifRAWTIFFDNGNEFRW2ARWCRWGPRCR3CR2PSDXMPavifppmjp2svgmagick"
)

//go:generate msgp

// ImageType is type of Image or Metadata file
//		ImageUnknown: "application/octet-stream"
//		ImageJPEG:    "image/jpeg"
//		ImagePNG:     "image/png"
//		ImageGIF:     "image/gif"
//		ImageBMP:     "image/bmp"
//		ImageWebP:    "image/webp"
//		ImageHEIF:    "image/heif"
//		ImageRAW:     "image/raw"
//		ImageTiff:    "image/tiff"
//		ImageDNG:     "image/x-adobe-dng"
//		ImageNEF:     "image/x-nikon-nef"
//		ImagePanaRAW: "image/x-panasonic-raw"
//		ImageARW:     "image/x-sony-arw"
//		ImageCRW:     "image/x-canon-crw"
//		ImageGPR:     "image/x-gopro-gpr"
//		ImageCR3:     "image/x-canon-cr3"
//		ImageCR2:     "image/x-canon-cr2"
//		ImagePSD:     "image/vnd.adobe.photoshop"
//		ImageXMP:     "application/rdf+xml"
//		ImageAVIF:    "image/avif"
//		ImagePPM:     "image/x-portable-pixmap"
type ImageType uint8

// IsUnknown returns true if the Image Type is unknown
func (it ImageType) IsUnknown() bool {
	return it == ImageUnknown
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (it ImageType) MarshalText() (text []byte, err error) {
	return []byte(it.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface that is
// used by encoding/json
func (it *ImageType) UnmarshalText(text []byte) (err error) {
	*it = FromString(string(text))
	return nil
}

func (it ImageType) String() string {
	if int(it) < len(_ImageTypeIndex)-1 {
		return _ImageTypeString[_ImageTypeIndex[it]:_ImageTypeIndex[it+1]]
	}
	return _ImageTypeString[:_ImageTypeIndex[1]]
}

// Extension returns the default extension for the Imagetype
func (it ImageType) Extension() string {
	if int(it) < len(_ImageTypeExtIndex)-1 {
		return _ImageTypeExtString[_ImageTypeExtIndex[it]:_ImageTypeExtIndex[it+1]]
	}
	return _ImageTypeExtString[:_ImageTypeExtIndex[1]]
}

// FromString returns an ImageType for the given content-type string or common filename extension
func FromString(str string) ImageType {
	// from content-type
	if it, ok := imageTypeValues[str]; ok {
		return it
	}
	// from extension
	if it, ok := imageTypeExtensions[strings.ToLower(str)]; ok {
		return it
	}
	return ImageUnknown
}

// Image file types Raw/Compressed/JPEG
const (
	ImageUnknown ImageType = iota
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
)

// ImageTypeValues maps a content-type string with an imagetype.
var imageTypeValues = map[string]ImageType{
	"application/octet-stream":  ImageUnknown,
	"image/jpeg":                ImageJPEG,
	"image/png":                 ImagePNG,
	"image/gif":                 ImageGIF,
	"image/bmp":                 ImageBMP,
	"image/webp":                ImageWebP,
	"image/heif":                ImageHEIF,
	"image/raw":                 ImageRAW,
	"image/tiff":                ImageTiff,
	"image/x-adobe-dng":         ImageDNG,
	"image/x-nikon-nef":         ImageNEF,
	"image/x-panasonic-raw":     ImagePanaRAW,
	"image/x-sony-arw":          ImageARW,
	"image/x-canon-crw":         ImageCRW,
	"image/x-gopro-gpr":         ImageGPR,
	"image/x-canon-cr3":         ImageCR3,
	"image/x-canon-cr2":         ImageCR2,
	"image/vnd.adobe.photoshop": ImagePSD,
	"application/rdf+xml":       ImageXMP,
	"image/avif":                ImageAVIF,
	"image/x-portable-pixmap":   ImagePPM,
	"image/jp2":                 ImageJP2K,
	"image/svg+xml":             ImageSVG,
	"image/magick":              ImageMAGICK,
}

// ImageTypeExtensions maps filename extensions with an imagetype.
var imageTypeExtensions = map[string]ImageType{
	"":        ImageUnknown,
	".jpg":    ImageJPEG,
	".png":    ImagePNG,
	".gif":    ImageGIF,
	".bmp":    ImageBMP,
	".webp":   ImageWebP,
	".heif":   ImageHEIF,
	".raw":    ImageRAW,
	".tiff":   ImageTiff,
	".dng":    ImageDNG,
	".nef":    ImageNEF,
	".rw2":    ImagePanaRAW,
	".arw":    ImageARW,
	".crw":    ImageCRW,
	".gpr":    ImageGPR,
	".cr3":    ImageCR3,
	".cr2":    ImageCR2,
	".psd":    ImagePSD,
	".xmp":    ImageXMP,
	".avif":   ImageAVIF,
	".ppm":    ImagePPM,
	".jp2":    ImageJP2K,
	".svg":    ImageSVG,
	".magick": ImageMAGICK,
}

// isTiff() Checks to see if an Image has the tiff format header.
//
func isTiff(buf []byte) bool {
	return len(buf) > 4 &&
		// BigEndian Tiff Image Header
		IsTiffBigEndian(buf[:4]) ||
		// LittleEndian Tiff Image Header
		IsTiffLittleEndian(buf[:4])
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func IsTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func IsTiffBigEndian(buf []byte) bool {
	return buf[0] == 0x4d &&
		buf[1] == 0x4d &&
		buf[2] == 0x00 &&
		buf[3] == 0x2a
}

// isCRW returns true if it matches an image/x-canon-crw with 14 bytes of the header.
//
// CanonCRWHeader is the file Header for a Canon CRW file. Currently only Little Endian support
// Reference: https://exiftool.org/canon_raw.html
func isCRW(buf []byte) bool {
	// ByteOrder: LittleEndian
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		// Signature: HEAPCCDR
		buf[6] == 0x48 &&
		buf[7] == 0x45 &&
		buf[8] == 0x41 &&
		buf[9] == 0x50 &&
		buf[10] == 0x43 &&
		buf[11] == 0x43 &&
		buf[12] == 0x44 &&
		buf[13] == 0x52
}

// isCR2 returns true if it matches an image/x-canon-cr2.
//
// CanonCR2Header is the Header for a Canon CR2 file
// 4 bytes after TiffSignature and before the beginning of IFDO
func isCR2(buf []byte) bool {
	return isTiff(buf) &&
		buf[8] == 0x43 &&
		buf[9] == 0x52 &&
		buf[10] == 0x02 &&
		buf[11] == 0x00
}

// isCR3 returns true if it matches an image/x-canon-cr3.
//
// ftyp box with major_brand: 'crx ' and compatible_brands: 'crx ' 'isom'
func isCR3(buf []byte) bool {
	return isFTYPBox(buf) &&
		isFTYPBrand(buf[8:12], "crx ")
}

// isHeif returns true if the header matches the start of a HEIF file.
//
// Major brands: heic, mif1, heix
// Minor brand contains:
func isHeif(buf []byte) bool {
	return isFTYPBox(buf) &&
		(isFTYPBrand(buf[8:12], "heic") ||
			isFTYPBrand(buf[8:12], "heix") ||
			(isFTYPBrand(buf[8:12], "mif1") && isFTYPBrand(buf[16:20], "heic")) ||
			(isFTYPBrand(buf[8:12], "mif1") && isFTYPBrand(buf[20:24], "heic")) ||
			(isFTYPBrand(buf[8:12], "msf1") && isFTYPBrand(buf[20:24], "hevc")))
}

// isFTYPBrand returns true if the Brand in []byte matches the brand in str.
// the Limit is 4 bytes
func isFTYPBrand(buf []byte, str string) bool {
	return buf[0] == str[0] && buf[1] == str[1] && buf[2] == str[2] && buf[3] == str[3]
}

// isFTYPBox returns true if the header matches an ftyp box.
// This indicates an ISO Base Media File Format.
func isFTYPBox(buf []byte) bool {
	return buf[0] == 0x0 &&
		buf[1] == 0x0 &&
		// buf[0:4] is 'ftyp' box size
		buf[4] == 0x66 &&
		buf[5] == 0x74 &&
		buf[6] == 0x79 &&
		buf[7] == 0x70
}

// isAVIF returns true if the header matches an ftyp box and
// an avif box.
//
func isAVIF(buf []byte) bool {
	return isFTYPBox(buf) &&
		(isFTYPBrand(buf[8:12], "avif") ||
			(isFTYPBrand(buf[8:12], "mif1") && isFTYPBrand(buf[20:24], "avif")))
}

// isBMP returns true if the header matches the start of a BMP file
// Bitmap Image
func isBMP(buf []byte) bool {
	return buf[0] == 0x42 &&
		buf[1] == 0x4D
}

// isRW2 returns true if the first 4 bytes match the Panasonic Tiff alternate
// header and bytes 8 through 12 match the RW2 header
//
func isRW2(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x55 &&
		buf[3] == 0x00 &&
		buf[8] == 0x88 &&
		buf[9] == 0xe7 &&
		buf[10] == 0x74 &&
		buf[11] == 0xd8
}

// isJPEG returns true if the first 2 bytes match a JPEG file header
//
// JPEG SOI Marker (FF D8)
func isJPEG(buf []byte) bool {
	return buf[0] == 0xff &&
		buf[1] == 0xd8
}

// isPNG returns true if the first 4 bytes match a PNG file header.
//
func isPNG(buf []byte) bool {
	return buf[0] == 0x89 &&
		buf[1] == 0x50 &&
		buf[2] == 0x4E &&
		buf[3] == 0x47
}

// isWebP returns true is the first 12 bytes match a WebP file header.
// RIFF and WebP
func isWebP(buf []byte) bool {
	return buf[0] == 0x52 &&
		buf[1] == 0x49 &&
		buf[2] == 0x46 &&
		buf[3] == 0x46 &&
		buf[8] == 0x57 &&
		buf[9] == 0x45 &&
		buf[10] == 0x42 &&
		buf[11] == 0x50
}

// isJPEG2000 returns true if the first 12 bytes match a JPEG2000 file header
//
func isJPEG2000(buf []byte) bool {
	return buf[0] == 0x0 &&
		buf[1] == 0x0 &&
		buf[2] == 0x0 &&
		buf[3] == 0xC &&
		buf[4] == 0x6A &&
		buf[5] == 0x50 &&
		buf[6] == 0x20 &&
		buf[7] == 0x20 &&
		buf[8] == 0xD &&
		buf[9] == 0xA &&
		buf[10] == 0x87 &&
		buf[11] == 0xA
}

// isPSD returns true if the header matches a PSDImage.
//
// PSD Photoshop document
func isPSD(buf []byte) bool {
	return buf[0] == 0x38 && buf[1] == 0x42 &&
		buf[2] == 0x50 && buf[3] == 0x53
}

// isXMP returns true if the header matches "<x:xmpmeta" start of a file.
//
// XMP sidecar files. The XMPHeader are the first 10bytes of an XMP sidecar.
func isXMP(buf []byte) bool {
	return buf[0] == 0x3c &&
		buf[1] == 0x78 &&
		buf[2] == 0x3a &&
		buf[3] == 0x78 &&
		buf[4] == 0x6d &&
		buf[5] == 0x70 &&
		buf[6] == 0x6d &&
		buf[7] == 0x65 &&
		buf[8] == 0x74 &&
		buf[9] == 0x61
}

// isGIF returns true if the header matches the header of a GIF version 87a
// or 89a.
func isGIF(buf []byte) bool {
	return buf[0] == 'G' &&
		buf[1] == 'I' &&
		buf[2] == 'F' &&
		buf[3] == '8' &&
		(buf[4] == '7' || buf[4] == '9') &&
		buf[5] == 'a'
}

func isPPM(buf []byte) bool {
	return buf[0] == 'P' &&
		(buf[1] == '3' || buf[1] == '6') &&
		(buf[2] == '\n' || buf[2] == '\r' || buf[2] == '\t' || buf[2] == ' ')
}
