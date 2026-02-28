package xmp

import (
	"github.com/evanoberholster/imagemeta/meta"
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
	Compression      uint16
}

func (t *Tiff) parse(p property) error {
	switch p.Name() {
	case Make:
		t.Make = parseString(p.Value())
	case Model:
		t.Model = parseString(p.Value())
	case ImageWidth:
		t.ImageWidth = uint16(parseUint(p.Value()))
	case ImageLength:
		t.ImageLength = uint16(parseUint(p.Value()))
	case Orientation:
		t.Orientation = meta.Orientation(parseUint(p.Value()))
	case Compression:
		t.Compression = uint16(parseUint(p.Value()))
	default:
		return ErrPropertyNotSet
	}
	return nil
}
