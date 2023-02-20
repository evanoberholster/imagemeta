package ifds

import (
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/apple"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/canon"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/nikon"
	"github.com/evanoberholster/imagemeta/exif2/ifds/mknote/sony"
)

// CameraModel is a Camera Model found in Exif
type CameraModel uint32

const (
	CameraModelUnknown CameraModel = iota
	CanonModelUnknown  CameraModel = 0x10000
	AppleModelUnknown  CameraModel = 0x20000
	NikonModelUnknown  CameraModel = 0x30000
	SonyModelUnknown   CameraModel = 0x40000
)

func (cm CameraModel) String() string {
	switch cm / 0x10000 {
	case 1: // Canon 0x10000
		return canon.CameraModel(cm).String()
	case 2: // Apple 0x20000
		return apple.CameraModel(cm).String()
	case 3: // Nikon 0x30000
		return nikon.CameraModel(cm).String()
	case 4: // Sony 0x40000
		return sony.CameraModel(cm).String()
	default:
		return ""
	}
}
