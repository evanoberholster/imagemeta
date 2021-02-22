// Package xmp provides functions for decoding .xmp sidecar files and XMP embedded within image files
package xmp

import (
	"bufio"
	"errors"
	"io"

	"github.com/evanoberholster/imagemeta/xmp/xmpns"
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

// Common Errors
var (

	// ErrNoXMP is returned when no XMP Root Tag is found.
	ErrNoXMP          = errors.New("xmp: error XMP not found")
	ErrPropertyNotSet = errors.New("xmp: error property not set")
)

// DebugMode when true would print items not parsed in XMP
var DebugMode = false

const (
	xmpBufferLength = 1024 * (3 / 2) // (1.5kb)
)

// ParseXmp reads XMP Metadata from the given reader and returns XMP.
//
func ParseXmp(r io.Reader) (xmp XMP, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok || br.Size() < xmpBufferLength {
		br = bufio.NewReaderSize(r, xmpBufferLength)
	}
	bufR := bufReader{r: br}
	rootTag, err := bufR.readRootTag()
	if err != nil {
		return XMP{}, err
	}

	var tag Tag
	for {
		tag, err = bufR.readTag(&xmp, rootTag)
		if err != nil {
			return xmp, err
		}
		if tag.isRootStopTag() {
			return
		}
	}
}

func (br *bufReader) readTag(xmp *XMP, parent Tag) (tag Tag, err error) {
	for {
		if tag, err = br.readTagHeader(parent); err != nil {
			break
		}
		if tag.isEndTag(parent.self) {
			break
		}
		var attr Attribute
		for br.hasAttribute() {
			if attr, err = br.readAttribute(&tag); err != nil {
				return
			}
			// Parse Attribute Value
			if err = xmp.parser(attr.property); err != nil {
				return
			}
		}
		if tag.isStartTag() {
			if tag.Is(xmpns.RDFSeq) || tag.Is(xmpns.RDFAlt) || tag.Is(xmpns.RDFBag) {
				if err = br.readSeqTags(xmp, tag); err != nil {
					return
				}
			} else {
				tag.val, err = br.readTagValue()
				if err != nil {
					return
				}
				// Parse Tag Value
				if err = xmp.parser(tag.property); err != nil {
					return
				}

				if tag, err = br.readTag(xmp, tag); err != nil {
					return
				}
			}
		}
		if tag.isRootStopTag() {
			return
		}
	}
	return
}

// Special Tags
// xmpMM:History -> stEvt
// rdf:Bag -> rdf:li
// rdf:Seq -> rdf:li
// rdf:Alt -> rdf:li
func (br *bufReader) readSeqTags(xmp *XMP, parent Tag) (err error) {
	var tag Tag
	for {
		if tag, err = br.readTagHeader(parent); err != nil {
			return
		}

		if tag.isEndTag(parent.self) {
			break
		}
		if tag.isStartTag() {
			var attr Attribute
			for br.hasAttribute() {
				attr, err = br.readAttribute(&tag)
				if err != nil {
					return
				}

				attr.parent = attr.self
				attr.self = parent.parent
				// Parse Attribute Value
				if err = xmp.parser(attr.property); err != nil {
					return
				}
			}

			if tag.val, err = br.readTagValue(); err != nil {
				return
			}
			tag.self = parent.parent
			tag.parent = parent.self
			// Parse Tag Value
			if err = xmp.parser(tag.property); err != nil {
				return
			}
		}

	}
	return
}
