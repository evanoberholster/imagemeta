package xmpns

import (
	"testing"
)

func TestProperty(t *testing.T) {
	p := NewProperty(XmpNamespace, RDF)

	if XmpNamespace != p.Namespace() {
		t.Errorf("Incorrect Property Namespace wanted %s got %s", XmpNamespace, p.Namespace())
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
	ns, _ := IdentifyNamespace("exif")

	if ns.Prefix() != "exif" {
		t.Errorf("Incorrect Name String wanted %s got %s", "exif", ns.Prefix())
	}
}

func TestIdentifyNamespace(t *testing.T) {
	tests := []struct {
		prefix   string
		expected Namespace
		found    bool
	}{
		{"xmp", XmpNamespace, true},
		{"xap", XmpNamespace, true},
		{"xmpBJ", XmpBJNamespace, true},
		{"xapBJ", XmpBJNamespace, true},
		{"xmpMM", XmpMMNamespace, true},
		{"xapMM", XmpMMNamespace, true},
		{"dc", DcNamespace, true},
		{"unknown", UnknownNamespace, true},
		{"nonexistent", UnknownNamespace, false},
	}

	for _, test := range tests {
		ns, found := IdentifyNamespace(test.prefix)
		if ns != test.expected || found != test.found {
			t.Errorf("IdentifyNamespace(%q) = %v, %v; want %v, %v", test.prefix, ns, found, test.expected, test.found)
		}
	}
}

func TestNamespace_Prefix(t *testing.T) {
	tests := []struct {
		ns       Namespace
		expected string
	}{
		{XmpNamespace, "xmp"},
		{XmpBJNamespace, "xmpBJ"},
		{XmpMMNamespace, "xmpMM"},
		{DcNamespace, "dc"},
		{UnknownNamespace, "unknown"},
		{Namespace(255), "unknown"},
	}

	for _, test := range tests {
		prefix := test.ns.Prefix()
		if prefix != test.expected {
			t.Errorf("Namespace(%v).Prefix() = %q; want %q", test.ns, prefix, test.expected)
		}
	}
}

func TestNamespace_String(t *testing.T) {
	tests := []struct {
		ns       Namespace
		expected string
	}{
		{XmpNamespace, "http://ns.adobe.com/xap/1.0/"},
		{DcNamespace, "http://purl.org/dc/elements/1.1/"},
		{XmpMMNamespace, "http://ns.adobe.com/xap/1.0/mm/"},
		{UnknownNamespace, "unknown"},
		{Namespace(255), "unknown"},
	}

	for _, test := range tests {
		space := test.ns.String()
		if space != test.expected {
			t.Errorf("Namespace(%v).Space() = %q; want %q", test.ns, space, test.expected)
		}
	}
}
