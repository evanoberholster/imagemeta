package xmpns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProperty(t *testing.T) {
	p := NewProperty(XmpNS, RDF)

	assert.Equal(t, XmpNS, p.Namespace(), "Property Namespace")

	assert.Equal(t, RDF, p.Name(), "Property Name")

	p1 := IdentifyProperty([]byte("xmp"), []byte("RDF"))

	assert.Equal(t, p, p1, "Property")

	assert.Equal(t, "xmp:RDF", p.String(), "Property String")

	assert.True(t, p.Equals(p1), "Property Equals")

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
