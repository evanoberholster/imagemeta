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

func NewNamespace(space Space, name Name) Namespace {
	return Namespace{uint8(space), uint8(name)}
}

func (n Namespace) Equals(namespace Namespace) bool {
	return n[0] == namespace[0] && n[1] == namespace[1]
}

func IdentifyNamespace(space []byte, name []byte) Namespace {
	return Namespace{uint8(IdentifySpace(space)), uint8(IdentifyTagName(name))}
}

func (n Namespace) Space() Space {
	return Space(n[0])
}

func (n Namespace) Name() Name {
	return Name(n[1])
}

func (n Namespace) String() string {
	return fmt.Sprintf("%s:%s", n.Space().String(), n.Name().String())
}
