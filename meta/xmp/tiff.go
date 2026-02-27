package xmp

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/xmp/xmpns"
)

// Tiff attributes of an XMP Packet.
//
//	xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
//
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#tiff
type Tiff struct {
	Make             string // Camera Make
	Model            string // Camera Model
	Software         string
	Copyright        []string
	ImageDescription []string
	ImageWidth       uint16
	ImageLength      uint16
	Orientation      meta.Orientation
}

func (t *Tiff) parse(p property) error {
	switch p.Name() {
	case xmpns.Make:
		t.Make = parseString(p.Value())
	case xmpns.Model:
		t.Model = parseString(p.Value())
	case xmpns.ImageWidth:
		t.ImageWidth = uint16(parseUint(p.Value()))
	case xmpns.ImageLength:
		t.ImageLength = uint16(parseUint(p.Value()))
	case xmpns.Orientation:
		t.Orientation = meta.Orientation(parseUint(p.Value()))
	default:
		return ErrPropertyNotSet
	}
	return nil
}
