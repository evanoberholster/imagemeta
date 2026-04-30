package xmp

import (
	"strings"
)

// property is an XML property
type property struct {
	val         []byte
	parent      Property
	self        Property
	pt          pType
	regionIndex int16
}

// Property returns the property's XMP Property
func (p property) Property() Property {
	return p.self
}

// SetParent sets the Property's parent
func (p *property) SetParent(parent Property) {
	p.parent = parent
}

// Parent returns the property's parent's XMP Property
func (p property) Parent() Property {
	return p.parent
}

// RegionIndex returns the zero-based region item index for mwg-rs region-list
// values. It returns -1 when the property is not inside a region-list item.
func (p property) RegionIndex() int {
	return int(p.regionIndex)
}

// Name returns the property's XMP local-name.
func (p property) Name() Name {
	return p.self.Name()
}

// Namespace returns the property's XMP namespace prefix.
func (p property) Namespace() Namespace {
	return p.self.Namespace()
}

// TagName is retained as a compatibility wrapper around Name.
func (p property) TagName() Name {
	return p.Name()
}

// NameSpace is retained as a compatibility wrapper around Namespace.
func (p property) NameSpace() Namespace {
	return p.Namespace()
}

// Value returns the property's Value
func (p property) Value() []byte {
	return p.val
}

// Is
func (p property) Is(p1 Property) bool {
	return p.self.Equals(p1)
}

func (p property) String() string {
	var b strings.Builder
	b.Grow(len(p.val) + 32)
	b.WriteByte('(')
	b.WriteString(p.parent.String())
	b.WriteString(") ")
	b.WriteString(p.self.String())
	b.WriteString(" \t Val:")
	b.Write(p.val)
	return b.String()
}

func (p property) valuePreview(limit int) string {
	if len(p.val) == 0 || limit <= 0 {
		return ""
	}
	if len(p.val) <= limit {
		return string(p.val)
	}
	return string(p.val[:limit]) + "..."
}

// Tag is an xmp Tag
type Tag struct {
	property
	t tagType
}

func (t Tag) String() string {
	var b strings.Builder
	b.Grow(len(t.val) + 40)
	b.WriteString(t.t.String())
	b.WriteString(": \t (")
	b.WriteString(t.parent.String())
	b.WriteString(") ")
	b.WriteString(t.self.String())
	b.WriteString(" \t Val:")
	b.Write(t.val)
	return b.String()
}

// Attribute is an xmp Attribute
type Attribute struct {
	property
}

func (attr Attribute) String() string {
	var b strings.Builder
	b.Grow(len(attr.val) + 40)
	b.WriteString("Attribute: (")
	b.WriteString(attr.parent.String())
	b.WriteString(") ")
	b.WriteString(attr.self.String())
	b.WriteString("=\"")
	b.Write(attr.val)
	b.WriteByte('"')
	return b.String()
}

// pType represents a property's type.
type pType uint8

const (
	noPType pType = iota
	attrPType
	tagPType
)

// ---------------------------------------------------

// tagType represents the Tag's type.
type tagType uint8

// Tag Types
const (
	noTag tagType = iota
	startTag
	soloTag
	stopTag
)

func (tt tagType) String() string {
	if int(tt) < len(tagTypeStrings) {
		return tagTypeStrings[tt]
	}
	return tagTypeStrings[noTag]
}

var tagTypeStrings = [...]string{
	noTag:    "No Tag",
	startTag: "Start Tag",
	soloTag:  "Solo Tag",
	stopTag:  "Stop Tag",
}
