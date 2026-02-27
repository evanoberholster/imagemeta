package xmpns

import (
	"testing"
)

func TestProperty(t *testing.T) {
	p := NewProperty(XmpNS, RDF)

	if XmpNS != p.Namespace() {
		t.Errorf("Incorrect Property Namespace wanted %s got %s", XmpNS, p.Namespace())
	}

	if RDF != p.Name() {
		t.Errorf("Incorrect Property Name wanted %s got %s", RDF, p.Name())
	}

	p1 := IdentifyProperty([]byte("xmp"), []byte("RDF"))

	if p1 != p {
		t.Errorf("Incorrect Property wanted %s got %s", p, p1)
	}

	str := "xmp:RDF"

	if str != p.String() || !p.Equals(p1) {
		t.Errorf("Incorrect Property String wanted %s got %s", str, p.String())
	}
}

func TestName(t *testing.T) {
	n := IdentifyName([]byte("subject"))

	if n.String() != "subject" {
		t.Errorf("Incorrect Name String wanted %s got %s", "subject", n.String())
	}
}

func TestNamespace(t *testing.T) {
	ns := IdentifyNamespace([]byte("exif"))

	if ns.String() != "exif" {
		t.Errorf("Incorrect Name String wanted %s got %s", "exif", ns.String())
	}
}
