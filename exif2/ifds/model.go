package ifds

import "github.com/evanoberholster/imagemeta/exif2/ifds/mknote/canon"

// CameraModel is a Camera Model found in Exif
type CameraModel uint32

const (
	CameraModelUnknown CameraModel = iota
)

func (cm CameraModel) String() string {
	if cm > 0x10000 {
		return canon.CameraModel(cm).String()
	}
	return ""
}
