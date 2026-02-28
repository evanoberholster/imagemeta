// Package xmp provides functions for decoding .xmp sidecar files and XMP embedded within image files
package xmp

import (
	"bytes"
	"errors"
	"io"
)

// Common Errors
var (
	// ErrNoXMP is returned when no XMP Root Tag is found.
	ErrNoXMP          = errors.New("xmp: error no XMP Tag found")
	ErrPropertyNotSet = errors.New("xmp: error property not set")

	// DebugMode when true would print items not parsed in XMP
	DebugMode = false

	xmpRootCloseTag = [...]byte{'<', '/', 'x', ':', 'x', 'm', 'p', 'm', 'e', 't', 'a', '>'}
)

// XMP contains parsed XML namespace groups.
// All namespaces except Basic are optional pointers and remain nil when absent.
type XMP struct {
	Aux          *Aux          // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	Exif         *Exif         // xmlns:exifEX="http://cipa.jp/exif/1.0/" and xmlns:exif="http://ns.adobe.com/exif/1.0/"
	Tiff         *Tiff         // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	Basic        Basic         // xmlns:xmp="http://ns.adobe.com/xap/1.0/" and xmlns:x="adobe:ns:meta/"
	DC           *DublinCore   // xmlns:dc="http://purl.org/dc/elements/1.1/"
	CRS          *CRS          // xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
	MM           *XMPMM        // xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
	Photoshop    *Photoshop    // xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
	DynamicMedia *DynamicMedia // xmlns:xmpDM="http://ns.adobe.com/xmp/1.0/DynamicMedia/"
	Lightroom    *Lightroom    // xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
}

// ParseXmp reads XMP Metadata from the given reader and returns XMP.
func ParseXmp(r io.Reader) (xmp XMP, err error) {
	return parseXMPStream(r)
}

// CleanXMPSuffixWhiteSpace returns the same slice with the whitespace after "</x:xmpmeta>" removed.
func CleanXMPSuffixWhiteSpace(buf []byte) []byte {
	for i := len(buf) - 1; i > 12; i-- {
		if buf[i] == '>' && buf[i-1] == 'a' {
			// </x:xmpmeta>
			if bytes.Equal(xmpRootCloseTag[:], buf[i-11:i+1]) {
				buf = buf[:i+1]
				return buf
			}
		}
	}
	return buf
}
