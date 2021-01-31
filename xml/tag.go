package xml

import (
	"bufio"
	"fmt"

	"github.com/evanoberholster/image-meta/xml/xmlname"
)

// Tag -
type Tag struct {
	parent xmlname.Namespace
	ns     xmlname.Namespace
	raw    []byte
	val    []byte
	t      TagType
}

func (t Tag) String() string {
	return fmt.Sprintf("%s: \t (%s) %s \t Val:%s", t.t.String(), t.parent.String(), t.ns.String(), string(t.val))
}

// Namespace returns the tag's namespace
func (t Tag) Namespace() xmlname.Namespace {
	return t.ns
}

// Parent returns the tag's parent namespace
func (t Tag) Parent() xmlname.Namespace {
	return t.parent
}

// Space returns the Tag's XMP Namespace
func (t Tag) Space() xmlname.Space {
	return t.ns.Space()
}

// Name returns the Tag's Name
func (t Tag) Name() xmlname.Name {
	return t.ns.Name()
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
	t.ns = xmlname.IdentifyNamespace(a, b)
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
	attr.ns = xmlname.IdentifyNamespace(a, b)

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
	parent xmlname.Namespace
	ns     xmlname.Namespace
	value  []byte
	raw    []byte
}

func (attr Attribute) String() string {
	return fmt.Sprintf("Attribute: %s=\"%s\"", attr.ns.String(), string(attr.value))
}

func (attr Attribute) Parent() xmlname.Namespace {
	return attr.parent
}

func (attr Attribute) Space() xmlname.Space {
	return attr.ns.Space()
}

func (attr Attribute) Name() xmlname.Name {
	return attr.ns.Name()
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
	return t.ns.Equals(xmlname.RDFSeq)
}

func (t Tag) isRDFAlt() bool {
	return t.ns.Equals(xmlname.RDFAlt)
}
func (t Tag) isRDFLi() bool {
	return t.ns.Equals(xmlname.RDFLi)
}
func (t Tag) isRootStopTag() bool {
	return t.ns.Equals(xmlname.XMPRootNamespace) && t.t == StopTag
}

func (t Tag) isEndTag(namespace xmlname.Namespace) bool {
	return t.t == StopTag && t.ns.Equals(namespace)
}

// ---------------------------------------------------

//
//func readNS(buf []byte) (ns xmlname.NS, b []byte) {
//	buf, b = readUntil(buf, colon)
//	ns = xmlname.IdentifyNS(buf)
//	return
//}
//
