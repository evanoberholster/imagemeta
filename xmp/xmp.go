// Package xmp provides functions for decoding .xmp sidecar files and XMP embedded within image files
package xmp

// DebugMode when true would print items not parsed in XMP
var DebugMode = false

// XMP contains the XML namespaces represented
type XMP struct {
	br    Reader
	Aux   Aux        // xmlns:aux="http://ns.adobe.com/exif/1.0/aux/"
	Exif  Exif       // xmlns:exifEX="http://cipa.jp/exif/1.0/" and xmlns:exif="http://ns.adobe.com/exif/1.0/"
	Tiff  Tiff       // xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
	Basic Basic      // xmlns:xmp="http://ns.adobe.com/xap/1.0/"
	DC    DublinCore // xmlns:dc="http://purl.org/dc/elements/1.1/"
	CRS   CRS
	MM    XMPMM
}
