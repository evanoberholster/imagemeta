package xml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/evanoberholster/image-meta/xml/xmlname"
)

// xmpRootTag starts with "<x:xmpmeta"
var xmpRootTag = [10]byte{60, 120, 58, 120, 109, 112, 109, 101, 116, 97}

const (
	newLine      byte = 10   // "\n"
	colon        byte = 58   // ":"
	startTag     byte = 60   // "<"
	equals       byte = 61   // "="
	endTag       byte = 62   // ">"
	space        byte = 0x20 // " "
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
		if t, err = xmp.readTag(xmlname.XMPRootNamespace); err != nil {
			fmt.Println(err)
			break
		}
		if t.isRootStopTag() {
			break
		}
	}

	return xmp, err
}
func (xmp *XMP) readRootTag() (discarded uint, err error) {
	var buf []byte
	for {
		if buf, err = xmp.br.Peek(10); err != nil {
			if err == io.EOF {
				err = ErrNoXMP
			}
			return
		}
		//if buf[0] == xmpRootTag[0] {
		//
		//}
		if bytes.EqualFold(xmpRootTag[:], buf) {
			// Read until end of the StartTag (RootTag)
			_, err = readUntilByte(xmp.br, endTag)
			return
		}
		discarded++
		xmp.br.Discard(1)
	}
}

func (xmp *XMP) decodeSeq(namespace xmlname.Namespace, val []byte) (err error) {

	switch namespace.Name() {
	case xmlname.ISOSpeedRatings:
		xmp.Exif.ISOSpeedRatings = uint32(parseUint(val))
	case xmlname.Creator:
		xmp.DC.Creator = append(xmp.DC.Creator, string(val))
	case xmlname.Rights:
		xmp.DC.Rights = append(xmp.DC.Rights, string(val))
	case xmlname.Title:
		xmp.DC.Title = append(xmp.DC.Rights, string(val))
	}
	//fmt.Println("Seq:", string(val))
	return
}

func (xmp *XMP) readRdfSeq(tag Tag) (Tag, error) {
	var err error
	if (tag.ns.Equals(xmlname.RDFSeq) || tag.ns.Equals(xmlname.RDFAlt)) && tag.TagType() == StartTag {
		//fmt.Println(tag.parent)
		// Read till end of sequence
		var child Tag
		for {
			// Start Tag
			if child, err = xmp.decodeTag(tag.parent); err != nil {
				return Tag{}, err
			}
			if child.t == SoloTag {
				// read Attributes
			}
			if child.t == StartTag {
				if err = child.readVal(xmp.br); err != nil {
					return tag, err
				}
				// ISOSpeed
				if err = xmp.decodeSeq(tag.parent, child.val); err != nil {
					fmt.Println(err)
				}
				continue
			}
			if child.isEndTag(xmlname.RDFSeq) || child.isEndTag(xmlname.RDFAlt) {
				//fmt.Println(tag)
				return tag, nil
			}
		}
	}
	return tag, nil
}

func (xmp *XMP) readTag(parent xmlname.Namespace) (Tag, error) {
	t, err := xmp.decodeTag(parent)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error Here", err)
		}
		return Tag{}, err
	}

	var attr Attribute
	for t.nextAttr() {
		attr, _ = t.attr()
		_ = attr
		//fmt.Println(attr)
	}

	if err = t.readVal(xmp.br); err != nil {
		return t, err
	}

	//fmt.Println(t)

	// read Next Tag
	// if tag is start Tag, read next tag
	if t.TagType() == StartTag {
		switch {
		case t.isRDFSeq(), t.isRDFAlt():
			return xmp.readRdfSeq(t)
		default:
			var child Tag
			for {
				if child, err = xmp.readTag(t.ns); err != nil {
					fmt.Println("Tags here", err)
					break
				}
				if child.isEndTag(t.ns) {
					break
				}
			}
		}
	}
	//if t.TagType() != StopTag {
	//xmp.readTag(br, t)
	//}

	return t, err
}

func (xmp *XMP) setValue(t Tag, attr Attribute) {
	//if t.ns == xmlname.Rdf {
	//	switch t.name {
	//	case xmlname.Description:
	//		xmp.setDescription(attr)
	//	default:
	//		return
	//	}
	//}
}

func (xmp *XMP) setDescription(attr Attribute) {
	if attr.ns.Space() == xmlname.Tiff {
		switch attr.ns.Name() {
		case xmlname.Make:
			xmp.Tiff.Make = string(attr.value)
		case xmlname.Model:
			xmp.Tiff.Model = string(attr.value)
		case xmlname.Orientation:
			//xmp.Tiff.Orientation = attr.parseUint8()
		case xmlname.ImageWidth:
			//xmp.Tiff.ImageWidth = attr.parseUint16()
		case xmlname.ImageLength:
			//xmp.Tiff.ImageLength = attr.parseUint16()
		default:
			fmt.Println("Not supported:", attr)
		}
	}
}

func (xmp *XMP) decodeTag(parentNS xmlname.Namespace) (t Tag, err error) {
	t.parent = parentNS
	if _, err = xmp.br.ReadSlice(startTag); err != nil {
		return
	}

	if t.raw, err = xmp.br.ReadSlice(endTag); err != nil {
		return
	}

	// StopTag
	if t.raw[0] == forwardSlash {
		t.t = StopTag     // set type StopTag
		t.raw = t.raw[1:] // remove forward slash

		t.readNamespace()
		return
	}

	// StartTag or Solo Tag
	if t.raw[len(t.raw)-2] == forwardSlash {
		t.t = SoloTag // set type SoloTag
	} else {
		t.t = StartTag
	}

	// Read Namespace
	t.readNamespace()
	return
}

func (xmp *XMP) handleTag(t Tag) {
	//switch t.name {
	//case xmlname.Flash:
	//	//xmp.Exif.Flash.read(t)
	//	//fmt.Println(t)
	//	// Read Flash
	//default:
	//	// if seg, bag, or alt
	//	// read Attributes
	//	var attr Attribute
	//	for t.nextAttr() {
	//		attr, _ = t.attr()
	//		_ = attr
	//		//fmt.Println(attr)
	//	}
	//}
}
