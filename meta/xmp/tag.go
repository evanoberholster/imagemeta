package xmp

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/meta/xmp/xmpns"
)

// property is an XML property
type property struct {
	val    []byte
	parent xmpns.Property
	self   xmpns.Property
	pt     pType
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

// Name returns the property's XMP Name
func (p property) Name() xmpns.Name {
	return p.self.Name()
}

// Namespace returns the property's XMP Namespace
func (p property) Namespace() xmpns.Namespace {
	return p.self.Namespace()
}

// Value returns the property's Value
func (p property) Value() []byte {
	return p.val
}

// Is
func (p property) Is(p1 xmpns.Property) bool {
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

// Tags
// ---------------------------------------------------

func (t Tag) isStartTag() bool {
	return t.t == startTag
}

func (t Tag) isRootStopTag() bool {
	return t.self.Equals(xmpns.XMPRootProperty) && t.t == stopTag
}

func (t Tag) isEndTag(p xmpns.Property) bool {
	return t.t == stopTag && t.self.Equals(p)
}

// ---------------------------------------------------

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
	return mapTagTypeString[tt]
}

var mapTagTypeString = map[tagType]string{
	noTag:    "No Tag",
	startTag: "Start Tag",
	soloTag:  "Solo Tag",
	stopTag:  "Stop Tag",
}
