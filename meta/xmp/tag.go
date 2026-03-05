package xmp

import (
	"fmt"
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
	return fmt.Sprintf("(%s) %s \t Val:%s", p.parent.String(), p.self.String(), string(p.val))
}

// Tag is an xmp Tag
type Tag struct {
	property
	t tagType
}

func (t Tag) String() string {
	return fmt.Sprintf("%s: \t (%s) %s \t Val:%s", t.t.String(), t.parent.String(), t.self.String(), string(t.val))
}

// Attribute is an xmp Attribute
type Attribute struct {
	property
}

func (attr Attribute) String() string {
	return fmt.Sprintf("Attribute: (%s) %s=\"%s\"", attr.parent.String(), attr.self.String(), string(attr.val))
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
