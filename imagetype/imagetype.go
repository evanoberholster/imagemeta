// Package imagetype provides types and functions for identifying Image document types
package imagetype

import (
	"errors"
)

// Errors
var (
	// ErrDataLength is an error for data length
	ErrDataLength = errors.New("error the data is not long enough")
)

// TiffBigEndianSignature is the Tiff Signature for BigEndian encoded images
//TiffBigEndianSignature = []byte{0x4d, 0x4d, 0x00, 0x2a}

// TiffLittleEndianSignature is the Tiff Signature for LittleEndian encoded images
//TiffLittleEndianSignature = []byte{0x49, 0x49, 0x2a, 0x00}

// JPEGStartOfImageMarker is the JPEG Start of Image Marker.
// JPEG SOI Marker
//JPEGStartOfImageMarker = []byte{0xff, 0xd8}

// PNGImageSignature is the marker for the start of a PNG Image.
// 4 Bytes
//PNGImageSignature = []byte{0x89, 0x50, 0x4E, 0x47}

// ImageType -
type ImageType uint8

// IsUnknown returns true if the Image Type is unknown
func (it ImageType) IsUnknown() bool {
	return it == ImageUnknown
}

// MarshalText implements the TextMarshaler interface that is
// used by encoding/json
func (it *ImageType) MarshalText() (text []byte, err error) {
	return []byte(imageTypeStrings[*it]), nil
}

func (it ImageType) String() string {
	return imageTypeStrings[it]
}

func FromString(str string) ImageType {
	return imageTypeValues[str]
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
)

// ImageTypeStrings - Map accepting ImageType and returning string
var imageTypeStrings = map[ImageType]string{
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
}

// ImageTypeValues - Map accepting string and returning Image Type
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
}

// isTiff() Checks to see if an Image has the tiff format header.
//
func isTiff(buf []byte) bool {
	return len(buf) > 4 &&
		// BigEndian Tiff Image Header
		isTiffBigEndian(buf[:4]) ||
		// LittleEndian Tiff Image Header
		isTiffLittleEndian(buf[:4])
}

// IsTiffLittleEndian checks the buf for the Tiff LittleEndian Signature
func isTiffLittleEndian(buf []byte) bool {
	return buf[0] == 0x49 &&
		buf[1] == 0x49 &&
		buf[2] == 0x2a &&
		buf[3] == 0x00
}

// IsTiffBigEndian checks the buf for the TiffBigEndianSignature
func isTiffBigEndian(buf []byte) bool {
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
// TODO: missing major brand and minor brand
// major_brand: crx // minor_version   : 1 // compatible_brands: crx isom
// ftyp
func isCR3(buf []byte) bool {
	return buf[0] == 0x0 &&
		buf[1] == 0x0 &&
		buf[2] == 0x0 &&
		buf[3] == 0x18 &&
		buf[4] == 0x66 &&
		buf[5] == 0x74 &&
		buf[6] == 0x79 &&
		buf[7] == 0x70 &&
		buf[8] == 0x63 &&
		buf[9] == 0x72 &&
		buf[10] == 0x78 &&
		buf[11] == 0x20
}

// isHeif returns true if the header matches the start of a HEIF file.
// TODO: missing major brand and minor brand
// ftyp
func isHeif(buf []byte) bool {
	return buf[0] == 0x0 &&
		buf[1] == 0x0 &&
		buf[2] == 0x0 &&
		//(buf[3] == 0x18 || buf[3] == 0x20) &&
		buf[4] == 0x66 &&
		buf[5] == 0x74 &&
		buf[6] == 0x79 &&
		buf[7] == 0x70 &&
		((buf[8] == 0x68 &&
			buf[9] == 0x65 &&
			buf[10] == 0x69 &&
			buf[11] == 0x63) ||
			(buf[8] == 0x6D &&
				buf[9] == 0x69 &&
				buf[10] == 0x66 &&
				buf[11] == 0x31))

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
