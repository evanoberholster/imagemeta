package xmp

import (
	"bufio"
	"fmt"

	"github.com/evanoberholster/image-meta/xmp/xmpns"
)

// property -
type property struct {
	parent xmpns.Property
	self   xmpns.Property
	val    []byte
}

// Property returns the property's XMP Property
func (p property) Property() xmpns.Property {
	return p.self
}

// SetParent sets the Property's parent
func (p *property) SetParent(parent xmpns.Property) {
	p.parent = parent
}

// Parent returns the property's parent's XMP Property
func (p property) Parent() xmpns.Property {
	return p.parent
}

// Name returns the property's XMP Property's Name
func (p property) Name() xmpns.Name {
	return p.self.Name()
}

// Namespace returns the property's XMP Property's Namespace
func (p property) Namespace() xmpns.Namespace {
	return p.self.Namespace()
}

// Is
func (p property) Is(p1 xmpns.Property) bool {
	return p.self.Equals(p1)
}

// Tag -
type Tag struct {
	property
	raw []byte
	t   tagType
}

func (t Tag) String() string {
	return fmt.Sprintf("%s: \t (%s) %s \t Val:%s", t.t.String(), t.parent.String(), t.self.String(), string(t.val))
}

func (t *Tag) readNamespace() {
	var a []byte
	var b []byte
	a, t.raw = readUntil(t.raw, markerCo)
	b, t.raw = readUntil(t.raw, markerSp)

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
				if t.raw[0] == markerSp {
					t.raw = t.raw[1:]
					continue
				}
				if t.raw[0] == newLine {
					t.raw = t.raw[1:]
					continue
				}
				if t.raw[0] == markerGt {
					return false
				}
				if t.raw[0] == forwardSlash && t.raw[1] == markerGt {
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
	a, t.raw = readUntil(t.raw, markerCo)
	// Read Name
	b, t.raw = readUntil(t.raw, markerEq)

	// Clean up possible new line.
	if len(b) > 0 {
		if b[len(b)-1] == newLine {
			b = b[:len(b)-1]
		}
	}
	attr.parent = t.self
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
			attr.val, t.raw = readUntil(t.raw[1:], quotes)
		} else if t.raw[0] == quotesAlt {
			attr.val, t.raw = readUntil(t.raw[1:], quotesAlt)
		}
	}
	return
}

func (t *Tag) readVal(br *bufio.Reader) (err error) {
	if t.t == startTag {
		// read TagValue
		var a []byte
		if a, err = br.ReadSlice(markerLt); err != nil {
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
	property
	raw []byte
}

func (attr Attribute) String() string {
	return fmt.Sprintf("Attribute: (%s) %s=\"%s\"", attr.parent.String(), attr.self.String(), string(attr.val))
}

// Tags
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

func (t Tag) isSoloTag() bool {
	return t.t == soloTag
}

func (t Tag) isStartTag() bool {
	return t.t == startTag
}

func (t Tag) isStopTag() bool {
	return t.t == stopTag
}

func (t Tag) isRootStopTag() bool {
	return t.self.Equals(xmpns.XMPRootProperty) && t.t == stopTag
}

func (t Tag) isEndTag(p xmpns.Property) bool {
	return t.t == stopTag && t.self.Equals(p)
}

// ---------------------------------------------------

// tagType represents the Tag's type.
type tagType uint8

// Tag Types
const (
	startTag tagType = iota
	soloTag
	stopTag
)

func (tt tagType) String() string {
	return mapTagTypeString[tt]
}

var mapTagTypeString = map[tagType]string{
	startTag: "Start Tag",
	soloTag:  "Solo Tag",
	stopTag:  "Stop Tag",
}

// tagType returns the Tag's type.
// Either startTag, soloTag or stopTag.
func (t Tag) tagType() tagType {
	return t.t
}
