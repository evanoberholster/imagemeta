package xml

import (
	"bufio"
	"fmt"

	"github.com/evanoberholster/image-meta/xml/xmpns"
)

// Tag -
type Tag struct {
	parent xmpns.Property
	self   xmpns.Property
	raw    []byte
	val    []byte
	t      TagType
}

func (t Tag) String() string {
	return fmt.Sprintf("%s: \t (%s) %s \t Val:%s", t.t.String(), t.parent.String(), t.self.String(), string(t.val))
}

// Property returns the tag's XMP Property
func (t Tag) Property() xmpns.Property {
	return t.self
}

// Parent returns the tag's parent's XMP Property
func (t Tag) Parent() xmpns.Property {
	return t.parent
}

// Namespace returns the Tag's XMP Property's Namespace
func (t Tag) Namespace() xmpns.Namespace {
	return t.self.Namespace()
}

// Name returns the Tag's XMP Property's Name
func (t Tag) Name() xmpns.Name {
	return t.self.Name()
}

func (t *Tag) readNamespace() {
	var a []byte
	var b []byte
	a, t.raw = readUntil(t.raw, colon)
	b, t.raw = readUntil(t.raw, space)

	// Clean up possible new line.
	if len(b) > 0 {
		if b[len(b)-1] == newLine {
			b = b[:len(b)-1]
		}
	}
	t.self = xmpns.IdentifyProperty(a, b)
}

func (t *Tag) nextAttr() bool {
	if !t.isStopTag() {
		for {
			if len(t.raw) > 0 {
				if t.raw[0] == space {
					t.raw = t.raw[1:]
					continue
				}
				if t.raw[0] == newLine {
					t.raw = t.raw[1:]
					continue
				}
				if t.raw[0] == endTag {
					return false
				}
				if t.raw[0] == forwardSlash && t.raw[1] == endTag {
					return false
				}
				return true
			}
			return false
		}
	}
	return false
}

func (t *Tag) attr() (attr Attribute, err error) {
	var a []byte
	var b []byte
	// Read Space
	a, t.raw = readUntil(t.raw, colon)
	// Read Name
	b, t.raw = readUntil(t.raw, equals)

	// Clean up possible new line.
	if len(b) > 0 {
		if b[len(b)-1] == newLine {
			b = b[:len(b)-1]
		}
	}
	attr.self = xmpns.IdentifyProperty(a, b)

	//var a []byte
	//// Read Namespace
	//a, t.raw = readUntil(t.raw, colon)
	//attr.ns = xmlname.IdentifyNS(a)
	//
	//// Read Name
	//a, t.raw = readUntil(t.raw, equals)
	//if len(a) > 0 {
	//	if attr.ns == xmlname.XMLns {
	//		//attr.name = xmlname.IdentifyNS(a)
	//	} else {
	//		attr.name = xmlname.IdentifyTagName(a)
	//	}
	//}

	// Read Value
	if len(t.raw) > 1 {
		if t.raw[0] == quotes {
			attr.value, t.raw = readUntil(t.raw[1:], quotes)
		} else if t.raw[0] == quotesAlt {
			attr.value, t.raw = readUntil(t.raw[1:], quotesAlt)
		}
	}
	return
}

func (t *Tag) readVal(br *bufio.Reader) (err error) {
	if t.t == StartTag {
		// read TagValue
		var a []byte
		if a, err = br.ReadSlice(startTag); err != nil {
			return
		}
		err = br.UnreadByte()
		if len(a) > 1 {
			a = a[:len(a)-1]
			for i := 0; i < len(a); i++ {
				if a[i] == newLine {
					a = a[:i]
				}
			}
			t.val = a
		}
	}
	return
}

// Attribute -
type Attribute struct {
	parent xmpns.Property
	self   xmpns.Property
	value  []byte
	raw    []byte
}

func (attr Attribute) String() string {
	return fmt.Sprintf("Attribute: %s=\"%s\"", attr.self.String(), string(attr.value))
}

// Parent returns the Parent's XMP Property
func (attr Attribute) Parent() xmpns.Property {
	return attr.parent
}

// Namespace returns the Attribute's Namespace
func (attr Attribute) Namespace() xmpns.Namespace {
	return attr.self.Namespace()
}

// Name returns the Attribute's Name
func (attr Attribute) Name() xmpns.Name {
	return attr.self.Name()
}

// TagType represents the Tag's type.
type TagType uint8

// Tag Types
const (
	StartTag TagType = iota
	SoloTag
	StopTag
)

func (tt TagType) String() string {
	return mapTagTypeString[tt]
}

var mapTagTypeString = map[TagType]string{
	StartTag: "Start Tag",
	SoloTag:  "Solo Tag",
	StopTag:  "Stop Tag",
}

func (t Tag) isStopTag() bool {
	return t.t == StopTag
}

// TagType returns the Tag's type.
// Either StartTag, SoloTag or StopTag.
func (t Tag) TagType() TagType {
	return t.t
}

// Sequence Tags
// ---------------------------------------------------

func (t Tag) isRDFSeq() bool {
	return t.self.Equals(xmpns.RDFSeq)
}

func (t Tag) isRDFAlt() bool {
	return t.self.Equals(xmpns.RDFAlt)
}
func (t Tag) isRDFLi() bool {
	return t.self.Equals(xmpns.RDFLi)
}
func (t Tag) isRootStopTag() bool {
	return t.self.Equals(xmpns.XMPRootProperty) && t.t == StopTag
}

func (t Tag) isEndTag(p xmpns.Property) bool {
	return t.t == StopTag && t.self.Equals(p)
}

// ---------------------------------------------------

//
//func readNS(buf []byte) (ns xmlname.NS, b []byte) {
//	buf, b = readUntil(buf, colon)
//	ns = xmlname.IdentifyNS(buf)
//	return
//}
//
