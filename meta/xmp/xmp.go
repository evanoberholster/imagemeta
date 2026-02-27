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
)

// XMP contains the XML namespaces represented
type XMP struct {
	Aux   Aux        // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	Exif  Exif       // xmlns:exifEX="http://cipa.jp/exif/1.0/" and xmlns:exif="http://ns.adobe.com/exif/1.0/"
	Tiff  Tiff       // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	Basic Basic      // xmlns:xmp="http://ns.adobe.com/xap/1.0/"
	DC    DublinCore // xmlns:dc="http://purl.org/dc/elements/1.1/"
	CRS   CRS
	MM    XMPMM
}

// ParseXmp reads XMP Metadata from the given reader and returns XMP.
//
func ParseXmp(r io.Reader) (xmp XMP, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	xr := newXMPReader(r)
	rootTag, err := xr.readRootTag()
	if err != nil {
		return XMP{}, err
	}

	var tag Tag
	for {
		if tag, err = xr.readTag(&xmp, rootTag); err != nil {
			return xmp, err
		}
		if tag.isRootStopTag() {
			return
		}
	}
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
