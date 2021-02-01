package xmp

import (
	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

// Tiff attributes of an XMP Packet.
//   xmlns:tiff="http://ns.adobe.com/tiff/1.0/"
// This implementation is incomplete and based on https://exiftool.org/TagNames/XMP.html#tiff
type Tiff struct {
	Make             string
	Model            string
	Software         string
	Copyright        []string
	ImageDescription []string
	ImageWidth       uint16
	ImageLength      uint16
	Orientation      uint8
}

func (t *Tiff) decode(p xmpns.Property, val []byte) error {
	switch p.Name() {
	case xmpns.Make:
		t.Make = parseString(val)
	case xmpns.Model:
		t.Model = parseString(val)
	case xmpns.ImageWidth:
		t.ImageWidth = uint16(parseUint(val))
	case xmpns.ImageLength:
		t.ImageLength = uint16(parseUint(val))
	case xmpns.Orientation:
		t.Orientation = uint8(parseUint(val))
	default:
		return ErrPropertyNotSet
	}
	return nil
}

// Orientation represents image orientation
type Orientation uint8

//1 = Horizontal (normal)
//2 = Mirror horizontal
//3 = Rotate 180
//4 = Mirror vertical
//5 = Mirror horizontal and rotate 270 CW
//6 = Rotate 90 CW
//7 = Mirror horizontal and rotate 90 CW
//8 = Rotate 270 CW
