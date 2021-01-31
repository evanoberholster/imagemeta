package xmpns

import "fmt"

var (
	// XMPRootProperty is the root Property for an XMP file
	XMPRootProperty = NewProperty(X, XmpMeta)
	// RDFSeq is the rdf:Seq Property
	RDFSeq = NewProperty(Rdf, Seq)
	// RDFAlt is the rdf:Alt Property
	RDFAlt = NewProperty(Rdf, Alt)
	// RDFLi is the rdf:Li Property
	RDFLi = NewProperty(Rdf, Li)
)

// Property is an XMP Namespace with the associated Property Name
type Property [2]uint8

// NewProperty returns the correspoding Property for the given Namespace and Name.
func NewProperty(ns Namespace, name Name) Property {
	return Property{uint8(ns), uint8(name)}
}

// Equals returns true if one property is equal to the other.
func (p1 Property) Equals(p2 Property) bool {
	return p1[0] == p2[0] && p1[1] == p2[1]
}

// IdentifyProperty returns a Property correspondent to the "space" and "name" byte values.
func IdentifyProperty(space []byte, name []byte) Property {
	return Property{uint8(IdentifyNamespace(space)), uint8(IdentifyName(name))}
}

// Namespace returns the property's XMP Namespace
func (p Property) Namespace() Namespace {
	return Namespace(p[0])
}

// Name returns the property's name
func (p Property) Name() Name {
	return Name(p[1])
}

func (p Property) String() string {
	return fmt.Sprintf("%s:%s", p.Namespace().String(), p.Name().String())
}
