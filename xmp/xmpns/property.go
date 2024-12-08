package xmpns

import "fmt"

var (
	// XMPRootProperty is the root Property for an XMP file
	XMPRootProperty = NewProperty(XNamespace, XmpMeta)
	// RDFDescription is the rdf:Description Property
	RDFDescription = NewProperty(RdfNamespace, Description)
	// RDFSeq is the rdf:Seq Property
	RDFSeq = NewProperty(RdfNamespace, Seq)
	// RDFBag is the rdf:Seq Property
	RDFBag = NewProperty(RdfNamespace, Bag)
	// RDFAlt is the rdf:Alt Property
	RDFAlt = NewProperty(RdfNamespace, Alt)
	// RDFLi is the rdf:Li Property
	RDFLi = NewProperty(RdfNamespace, Li)
)

// Property is an XMP Namespace with the associated Property Name
type Property struct {
	ns   Namespace
	name Name
}

// NewProperty returns the correspoding Property for the given Namespace and Name.
func NewProperty(ns Namespace, name Name) Property {
	return Property{ns, name}
}

// Equals returns true if one property is equal to the other.
func (p Property) Equals(p1 Property) bool {
	return p.ns == p1.ns && p.name == p1.name
}

// IdentifyProperty returns a Property correspondent to the "space" and "name" byte values.
func IdentifyProperty(space []byte, name []byte) Property {
	ns, _ := IdentifyNamespace(string(space))
	return Property{ns, IdentifyName(name)}
}

// Namespace returns the property's XMP Namespace
func (p Property) Namespace() Namespace {
	return p.ns
}

// Name returns the property's name
func (p Property) Name() Name {
	return p.name
}

func (p Property) String() string {
	return fmt.Sprintf("%s:%s", p.ns.Prefix(), p.name.String())
}
