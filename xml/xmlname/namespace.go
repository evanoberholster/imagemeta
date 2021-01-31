package xmlname

import "fmt"

var (
	// XMPRootNamespace is the root Namespace for an XMP file
	XMPRootNamespace = NewNamespace(X, XmpMeta)
	// RDFSeq is the rdf:Seq Namespace
	RDFSeq = NewNamespace(Rdf, Seq)
	// RDFAlt is the rdf:Alt Namespace
	RDFAlt = NewNamespace(Rdf, Alt)
	// RDFLi is the rdf:Li Namespace
	RDFLi = NewNamespace(Rdf, Li)
)

// Namespace is an NS with a TagName
type Namespace [2]uint8

// NewNamespace returns the correspoding Namespace for the given Space and Name.
func NewNamespace(space Space, name Name) Namespace {
	return Namespace{uint8(space), uint8(name)}
}

// Equals returns true if one namespace is equal to the other.
func (n Namespace) Equals(namespace Namespace) bool {
	return n[0] == namespace[0] && n[1] == namespace[1]
}

// IdentifyNamespace returns a Namespace correspondent to the "space" and "name" byte values.
func IdentifyNamespace(space []byte, name []byte) Namespace {
	return Namespace{uint8(IdentifySpace(space)), uint8(IdentifyName(name))}
}

// Space returns the XMP Namespace
func (n Namespace) Space() Space {
	return Space(n[0])
}

// Name returns the property name
func (n Namespace) Name() Name {
	return Name(n[1])
}

func (n Namespace) String() string {
	return fmt.Sprintf("%s:%s", n.Space().String(), n.Name().String())
}
