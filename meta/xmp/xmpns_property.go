package xmp

var (
	// XMPRootProperty is the root Property for an XMP file
	XMPRootProperty = NewProperty(XNS, XmpMeta)
	// RDFDescription is the rdf:Description Property
	RDFDescription = NewProperty(RdfNS, Description)
	// RDFSeq is the rdf:Seq Property
	RDFSeq = NewProperty(RdfNS, Seq)
	// RDFBag is the rdf:Seq Property
	RDFBag = NewProperty(RdfNS, Bag)
	// RDFAlt is the rdf:Alt Property
	RDFAlt = NewProperty(RdfNS, Alt)
	// RDFLi is the rdf:Li Property
	RDFLi = NewProperty(RdfNS, Li)
)

// Property is an XMP Namespace with the associated Name.
type Property uint32

// NewProperty returns the corresponding property for the given namespace and name.
func NewProperty(ns Namespace, name Name) Property {
	return Property((uint32(ns) << 16) | uint32(name))
}

// Equals returns true if one property is equal to the other.
func (p Property) Equals(p1 Property) bool {
	return p == p1
}

// IdentifyProperty returns a Property correspondent to the "space" and "name" byte values.
func IdentifyProperty(space []byte, name []byte) Property {
	return NewProperty(IdentifyNamespace(space), IdentifyName(name))
}

// Namespace returns the property's XMP namespace.
func (p Property) Namespace() Namespace {
	//nolint:gosec // upper 8 bits are masked before narrowing.
	return Namespace((uint32(p) >> 16) & 0xFF)
}

// Name returns the property's XMP local-name.
func (p Property) Name() Name {
	//nolint:gosec // lower 16 bits are masked before narrowing.
	return Name(uint32(p) & 0xFFFF)
}

// NameSpace is retained as a compatibility wrapper around Namespace.
func (p Property) NameSpace() Namespace {
	return p.Namespace()
}

// TagName is retained as a compatibility wrapper around Name.
func (p Property) TagName() Name {
	return p.Name()
}

func (p Property) String() string {
	return p.Namespace().String() + ":" + p.Name().String()
}
