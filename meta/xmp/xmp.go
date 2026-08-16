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

	xmpRootCloseTag = [...]byte{'<', '/', 'x', ':', 'x', 'm', 'p', 'm', 'e', 't', 'a', '>'}
)

// ParseOptions configures parser behavior.
type ParseOptions struct {
	// Debug enables warning output for non-fatal parse issues.
	Debug bool
}

// nsFlag is a bitset of parsed namespace presence flags.
//
// Each canonical Namespace value maps to one bit position (1 << ns).
// Aliased namespaces are normalized via canonicalNamespace before checking/setting.
type nsFlag uint64

// XMP contains parsed XML namespace groups.
// Namespace values remain zero when not present in the source metadata.
type XMP struct {
	Aux          Aux          // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	Exif         Exif         // xmlns:exifEX="http://cipa.jp/exif/1.0/" and xmlns:exif="http://ns.adobe.com/exif/1.0/"
	Tiff         Tiff         // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	Basic        Basic        // xmlns:xmp="http://ns.adobe.com/xap/1.0/" and xmlns:x="adobe:ns:meta/"
	DC           DublinCore   // xmlns:dc="http://purl.org/dc/elements/1.1/"
	CRS          CRS          // xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
	MM           XMPMM        // xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
	Photoshop    Photoshop    // xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
	DynamicMedia DynamicMedia // xmlns:xmpDM="http://ns.adobe.com/xmp/1.0/DynamicMedia/"
	Lightroom    Lightroom    // xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
	Regions      RegionInfo   // xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/"
	nsParsed     nsFlag
}

// IsParsed reports whether a namespace group was present in the parsed XMP document.
func (x XMP) IsParsed(ns Namespace) bool {
	canonical, ok := canonicalNamespace(ns)
	if !ok {
		return false
	}
	return x.hasNamespaceFlag(canonical)
}

func (x *XMP) markParsed(ns Namespace) {
	canonical, ok := canonicalNamespace(ns)
	if !ok {
		return
	}
	x.setNamespaceFlag(canonical)
}

func (x *XMP) setNamespaceFlag(ns Namespace) {
	if ns >= 64 {
		return
	}
	bit := nsFlag(1) << ns
	x.nsParsed |= bit
}

func (x XMP) hasNamespaceFlag(ns Namespace) bool {
	if ns >= 64 {
		return false
	}
	bit := nsFlag(1) << ns
	return x.nsParsed&bit != 0
}

func canonicalNamespace(ns Namespace) (Namespace, bool) {
	switch ns {
	case AuxNS:
		return AuxNS, true
	case ExifNS, ExifEXNS:
		return ExifNS, true
	case TiffNS:
		return TiffNS, true
	case XNS, XapNS, XmpNS:
		return XmpNS, true
	case DcNS:
		return DcNS, true
	case CrsNS:
		return CrsNS, true
	case XmpMMNS, XapMMNS, StEvtNS, StRefNS:
		return XmpMMNS, true
	case PhotoshopNS:
		return PhotoshopNS, true
	case XmpDMNS:
		return XmpDMNS, true
	case LrNS:
		return LrNS, true
	case MwgRSNS, StDimNS, StAreaNS, AppleFiNS:
		return MwgRSNS, true
	default:
		return 0, false
	}
}

// ParseXmp reads XMP Metadata from the given reader and returns XMP.
func ParseXmp(r io.Reader) (xmp XMP, err error) {
	return ParseXmpWithOptions(r, ParseOptions{})
}

// ParseXmpWithOptions reads XMP metadata from the given reader and returns XMP.
func ParseXmpWithOptions(r io.Reader, opts ParseOptions) (xmp XMP, err error) {
	return parseXMPStream(r, opts)
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
