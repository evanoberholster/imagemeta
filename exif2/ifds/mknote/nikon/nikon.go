package nikon

// IsNikonMkNoteHeaderBytes represents "Nikon" the first 5 bytes of the
func IsNikonMkNoteHeaderBytes(buf []byte) bool {
	return "Nikon" == string(buf[:5])
}

// CameraModel is a Nikon Camera Model found in Exif
type CameraModel uint32

func (cm CameraModel) String() string {
	if str, ok := mapCameraModelString[cm]; ok {
		return str
	}
	return ""
}

// CameraModelFromString returns a nikon camera model from the given string
func CameraModelFromString(str string) (CameraModel, bool) {
	if cm, ok := mapStringCameraModel[str]; ok {
		return cm, true
	}
	return NikonModelUnknown, false
}

// Nikon Camera Models
const (
	NikonModelUnknown CameraModel = iota + 0x30000
)

var mapCameraModelString = map[CameraModel]string{}

var mapStringCameraModel = map[string]CameraModel{}
