package xmp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

// xmpRootTag starts with "<x:xmpmeta"
var xmpRootTag = [10]byte{60, 120, 58, 120, 109, 112, 109, 101, 116, 97}

const (
	newLine      byte = 10   // "\n"
	markerCo     byte = 58   // ":"
	markerLt     byte = 60   // "<"
	markerEq     byte = 61   // "="
	markerGt     byte = 62   // ">"
	markerSp     byte = 0x20 // " "
	quotesAlt    byte = 0x27 // "'"
	quotes       byte = 0x22 // """
	forwardSlash byte = 0x2f // "/"
)

// Read -
func Read(r io.Reader) (XMP, error) {
	var err error
	xmp := XMP{
		br: bufio.NewReaderSize(r, 1024*6),
	}
	// find start of XML
	_, err = xmp.readRootTag()
	if err != nil {
		return xmp, err
	}

	// read Tags
	var t Tag
	for {
		if t, err = xmp.readTag(xmpns.XMPRootProperty); err != nil {
			fmt.Println(err)
			break
		}
		if t.isRootStopTag() {
			break
		}
	}

	return xmp, err
}

// Needs optimization
func (xmp *XMP) readRootTag() (discarded uint, err error) {
	var buf []byte
	for {
		if buf, err = xmp.br.Peek(10); err != nil {
			if err == io.EOF {
				err = ErrNoXMP
			}
			return
		}
		if buf[0] == xmpRootTag[0] {
			if bytes.EqualFold(xmpRootTag[:], buf) {
				// Read until end of the StartTag (RootTag)
				_, err = readUntilByte(xmp.br, markerGt)
				return
			}
		}
		discarded++
		xmp.br.Discard(1)
	}
}

func (xmp *XMP) decodeTag(tag Tag) (err error) {
	switch tag.Namespace() {
	case xmpns.AuxNS:
		return xmp.Aux.decode(tag.property)
	case xmpns.ExifNS:
		return xmp.Exif.decode(tag)
	case xmpns.TiffNS:
		return xmp.Tiff.decode(tag.self, tag.val)
	case xmpns.XmpNS:
		return xmp.Basic.decode(tag.property)
	case xmpns.DcNS:
		return xmp.DC.decode(tag.property)
	case xmpns.RdfNS:
		switch tag.Name() {
		case xmpns.Description:
			// decode Attributes
			var attr Attribute
			for tag.nextAttr() {
				attr, _ = tag.attr()
				if err := xmp.decodeAttr(attr); err != nil {
					if err == ErrPropertyNotSet && DebugMode {
						fmt.Println("Attr NotSet:", attr)
					}
				}
			}
			return nil
		}
	default:
		return ErrPropertyNotSet
	}
	return
}

func (xmp *XMP) decodeAttr(attr Attribute) (err error) {
	switch attr.Namespace() {
	case xmpns.XMLnsNS, xmpns.RdfNS:
		// Null operation
		return
	case xmpns.DcNS:
		return xmp.DC.decode(attr.property)
	case xmpns.AuxNS:
		return xmp.Aux.decode(attr.property)
	case xmpns.XmpNS:
		return xmp.Basic.decode(attr.property)
	case xmpns.TiffNS:
		return xmp.Tiff.decode(attr.self, attr.val)
	}
	return ErrPropertyNotSet
}

func (xmp *XMP) decodeSeq(p property) (err error) {
	switch p.Namespace() {
	case xmpns.CrsNS:
		return xmp.CRS.decode(p)
	case xmpns.DcNS:
		return xmp.DC.decode(p)
	case xmpns.ExifNS:
		switch p.Name() {
		case xmpns.ISOSpeedRatings: // Needs fixing
			xmp.Exif.ISOSpeedRatings = uint32(parseUint(p.val))
			return
		}
	}
	return ErrPropertyNotSet
}

func (xmp *XMP) readRdfSeq(tag Tag) (Tag, error) {
	var err error
	if tag.isStartTag() && (tag.Is(xmpns.RDFSeq) || tag.Is(xmpns.RDFAlt)) {
		//fmt.Println(tag.parent)
		// Read till end of sequence
		var child Tag
		for {
			// Start Tag
			if child, err = xmp.readTagHeader(tag.parent); err != nil {
				return Tag{}, err
			}
			if !child.isStopTag() {
				// read Attributes
				var attr Attribute
				for child.nextAttr() {
					attr, _ = child.attr()
					attr.SetParent(tag.parent)
					xmp.decodeAttr(attr)
				}
			}
			if child.isStartTag() {
				// ISOSpeed
				if err = xmp.decodeSeq(property{self: tag.parent, val: child.val}); err != nil {
					if err == ErrPropertyNotSet && DebugMode {
						fmt.Println("Seq NotSet:", tag.Parent(), string(child.val))
					}
				}
				continue
			}
			if child.isStopTag() && child.Is(tag.self) {
				//fmt.Println(tag)
				return tag, nil
			}
		}
	}
	return tag, nil
}

func (xmp *XMP) readTag(parent xmpns.Property) (Tag, error) {
	// Read Tag Header
	tag, err := xmp.readTagHeader(parent)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error Here", err)
		}
		return Tag{}, err
	}

	if tag.isStopTag() {
		return tag, nil
	}

	// DebugMode
	if DebugMode {
		fmt.Println(tag)
	}

	// Process Tag
	if err := xmp.decodeTag(tag); err != nil {
		if err != ErrPropertyNotSet {
			return tag, err
		}
		if DebugMode {
			fmt.Println("Tag NotSet:", tag)
		}
	}

	// Process Sequential Tags
	if tag.Namespace() == xmpns.RdfNS {
		if tag.Is(xmpns.RDFSeq) || tag.Is(xmpns.RDFAlt) {
			return xmp.readRdfSeq(tag)
		}
	}

	// Process Child tags
	if tag.isStartTag() {
		var child Tag
		for {
			if child, err = xmp.readTag(tag.self); err != nil {
				return tag, err
			}
			if child.isEndTag(tag.self) {
				break
			}
		}
	}
	return tag, err
}

func (xmp *XMP) readTagHeader(parent xmpns.Property) (t Tag, err error) {
	t.SetParent(parent)
	if _, err = xmp.br.ReadSlice(markerLt); err != nil {
		return
	}

	if t.raw, err = xmp.br.ReadSlice(markerGt); err != nil {
		return
	}

	// StopTag
	if t.raw[0] == forwardSlash {
		t.t = stopTag     // set type StopTag
		t.raw = t.raw[1:] // remove forward slash

		// Read Namespace
		t.readNamespace()
		return
	}

	// StartTag or Solo Tag
	if t.raw[len(t.raw)-2] == forwardSlash {
		t.t = soloTag // set type SoloTag
	} else {
		t.t = startTag
	}

	// Read Namespace
	t.readNamespace()

	// Read Tag Value
	err = t.readVal(xmp.br)
	return
}
