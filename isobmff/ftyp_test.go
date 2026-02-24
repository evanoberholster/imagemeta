package isobmff

import (
	"bytes"
	"testing"
)

func TestBrandFromBufAdditionalBrands(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		want Brand
	}{
		{name: "heif", buf: []byte("heif"), want: brandHeif},
		{name: "avis", buf: []byte("avis"), want: brandAvis},
		{name: "3gp6", buf: []byte("3gp6"), want: brand3GP6},
		{name: "3g2a", buf: []byte("3g2a"), want: brand3G2A},
		{name: "M4V", buf: []byte("M4V "), want: brandM4V},
		{name: "mp71", buf: []byte("mp71"), want: brandMp71},
		{name: "dash", buf: []byte("dash"), want: brandDash},
		{name: "jxl", buf: []byte("jxl "), want: brandJxl},
		{name: "qt", buf: []byte("qt  "), want: brandQt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := brandFromBuf(tt.buf)
			if got != tt.want {
				t.Fatalf("brandFromBuf(%q) = %v, want %v", tt.buf, got, tt.want)
			}
			if got.String() != string(tt.buf) {
				t.Fatalf("brand string = %q, want %q", got.String(), string(tt.buf))
			}
		})
	}
}

func TestReadFTYPSkipsJXLSignatureBox(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x0C, // size
		'J', 'X', 'L', ' ', // type
		0x0D, 0x0A, 0x87, 0x0A, // JXL signature payload
		0x00, 0x00, 0x00, 0x10, // size
		'f', 't', 'y', 'p', // type
		'a', 'v', 'i', 'f', // major brand
		'0', '0', '0', '1', // minor version
	}

	r := NewReader(bytes.NewReader(data), nil, nil, nil)
	t.Cleanup(r.Close)

	if err := r.ReadFTYP(); err != nil {
		t.Fatalf("ReadFTYP() error = %v", err)
	}
	if r.ftyp.MajorBrand != brandAvif {
		t.Fatalf("MajorBrand = %v, want %v", r.ftyp.MajorBrand, brandAvif)
	}
}
