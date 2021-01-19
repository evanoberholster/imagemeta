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
	colon    byte = 58   // ":"
	startTag byte = 60   // "<"
	equals   byte = 61   // "="
	endTag   byte = 62   // ">"
	space         = 0x20 // " "
	quotes        = 0x22 // \"
)

// Read -
func Read(r io.Reader) (xmp XMP, err error) {
	br := bufio.NewReader(r)
	// find start of XML
	_, err = findRootTag(br)
	if err != nil {
		return
	}

	// find end of tag
	n, err := readUntilByte(br, endTag)
	if err != nil {
		return
	}
	_ = n
	//fmt.Println(n)

	xmp.readTag(br)
	xmp.readTag(br)
	//fmt.Println(br.ReadByte())
	return
}

func findRootTag(br *bufio.Reader) (discarded int, err error) {
	var tag []byte
	for {
		tag, err = br.Peek(10)
		if err != nil {
			return
		}
		if bytes.EqualFold(xmpRootTag[:], tag) {
			return
		}
		discarded++
		br.Discard(1)
	}
}

func readUntilByte(br *bufio.Reader, end byte) (n int, err error) {
	var b byte
	for {
		b, err = br.ReadByte()
		if err != nil {
			return
		}
		n++
		if b == end {
			return
		}
	}
}

// Attribute -
type Attribute struct {
	ns    ns
	name  xmlname.Name
	value []byte
}

func (attr Attribute) String() string {
	return fmt.Sprintf("Attribute: %s:%s=\"%s\"", mapNSString[attr.ns], attr.name.String(), string(attr.value))
}

// Tag -
type Tag struct {
	ns   ns
	name xmlname.Name
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag: %s:%s", mapNSString[t.ns], t.name)
}

func (t *Tag) readNS(br *bufio.Reader) error {
	buf, err := br.ReadSlice(colon)
	if err != nil {
		return err
	}
	buf = buf[:len(buf)-1]
	t.ns = identifyNS(buf)
	return nil
}

func (t *Tag) readName(br *bufio.Reader) error {
	buf, err := br.ReadSlice(space)
	if err != nil {
		return err
	}
	buf = buf[:len(buf)-1]
	t.name = identifyName(buf)
	return nil
}

func readUntil(buf []byte, delimiter byte) (a []byte, b []byte) {
	for i := 0; i < len(buf); i++ {
		if buf[i] == delimiter {
			return buf[:i], buf[i+1:]
		}
	}
	return nil, nil
}

func (attr *Attribute) readName(buf []byte) []byte {
	var a []byte
	a, buf = readUntil(buf, equals)
	attr.name = identifyName(a)
	return buf
}

func (attr *Attribute) readNS(buf []byte) []byte {
	var a []byte
	a, buf = readUntil(buf, colon)
	attr.ns = identifyNS(a)
	return buf
}

func (attr *Attribute) readValue(buf []byte) (b []byte) {
	if buf[0] == quotes {
		val, buf := readUntil(buf[1:], quotes)
		attr.value = val
		return buf
	}
	// TODO: write error
	panic("Attribute Error")
}

func (attr *Attribute) parseUint8() uint8 {
	val := uint8(parseUint(attr.value))
	return val
}

func (attr *Attribute) parseUint16() uint16 {
	val := uint16(parseUint(attr.value))
	return val
}

func (attr *Attribute) parseUint32() uint32 {
	val := uint32(parseUint(attr.value))
	return val
}

func parseUint(buf []byte) (u uint64) {
	for i := 0; i < len(buf); i++ {
		u *= 10
		u += uint64(buf[i] - '0')
	}
	return
}

func readAttr(buf []byte) (b []byte, attr Attribute) {

	buf = attr.readNS(buf)    // Read Namespace
	buf = attr.readName(buf)  // Read Name
	buf = attr.readValue(buf) // Read Value

	return buf, attr
}

func (xmp *XMP) setValue(t Tag, attr Attribute) {
	if t.ns == rdfNS {
		switch t.name {
		case xmlname.Description:
			xmp.setDescription(attr)
		default:
			return
		}
	}
}

func (xmp *XMP) setDescription(attr Attribute) {
	if attr.ns == tiffNS {
		switch attr.name {
		case xmlname.Make:
			xmp.Tiff.Make = string(attr.value)
		case xmlname.Model:
			xmp.Tiff.Model = string(attr.value)
		case xmlname.Orientation:
			xmp.Tiff.Orientation = attr.parseUint8()
		case xmlname.ImageWidth:
			xmp.Tiff.ImageWidth = attr.parseUint16()
		case xmlname.ImageLength:
			xmp.Tiff.ImageLength = attr.parseUint16()
		default:
			fmt.Println("Not supported:", attr)
		}
	}
}

func (xmp *XMP) readTag(br *bufio.Reader) (tag Tag, err error) {
	_, err = readUntilByte(br, startTag)
	if err != nil {
		return
	}

	// Read Tag
	if err = tag.readNS(br); err != nil {
		return
	}
	if err = tag.readName(br); err != nil {
		return
	}

	var buf []byte
	buf, err = br.ReadSlice(endTag)
	if err != nil {
		return
	}
	var attr Attribute
	for len(buf) > 0 {
		if buf[0] == space {
			buf = buf[1:]
			continue
		}
		if buf[0] == []byte("\n")[0] {
			buf = buf[1:]
			continue
		}
		buf, attr = readAttr(buf)
		xmp.setValue(tag, attr)
		//fmt.Println(attr)
		if buf[0] == endTag {
			break
		}
	}
	return
}

// ns is an XML Namespace
type ns uint8

const (
	unknownNS ns = iota
	xNS
	xmlnsNS
	xmpNS
	xmpMMNS
	tiffNS
	exifNS
	exifEXNS
	dcNS
	auxNS
	photoshopNS
	crsNS
	lrNS
	rdfNS
)

func identifyNS(buf []byte) (n ns) {
	return mapStringNS[string(buf)]
}

func identifyName(buf []byte) (n xmlname.Name) {
	return xmlname.MapStringName[string(buf)]
}

var mapStringNS = map[string]ns{
	"Unknown":   unknownNS,
	"x":         xNS,
	"xmlns":     xmlnsNS,
	"xmp":       xmpNS,
	"xmpMM":     xmpMMNS,
	"tiff":      tiffNS,
	"exif":      exifNS,
	"exifEX":    exifEXNS,
	"dc":        dcNS,
	"aux":       auxNS,
	"photoshop": photoshopNS,
	"crs":       crsNS,
	"lr":        lrNS,
	"rdf":       rdfNS,
}

var mapNSString = map[ns]string{
	unknownNS:   "Unknown",
	xNS:         "x",
	xmlnsNS:     "xmlns",
	xmpNS:       "xmp",
	xmpMMNS:     "xmpMM",
	tiffNS:      "tiff",
	exifNS:      "exif",
	exifEXNS:    "exifEX",
	dcNS:        "dc",
	auxNS:       "aux",
	photoshopNS: "photoshop",
	crsNS:       "crs",
	lrNS:        "lr",
	rdfNS:       "rdf",
}
