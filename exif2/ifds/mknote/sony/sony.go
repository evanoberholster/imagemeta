package sony

// CameraModel is a Sony Camera Model found in Exif
type CameraModel uint32

// CameraModelFromString returns a sony camera model from the given string
func CameraModelFromString(str string) (CameraModel, bool) {
	if cm, ok := mapStringCameraModel[str]; ok {
		return cm, true
	}
	return SonyModelUnknown, false
}

var mapStringCameraModel = map[string]CameraModel{}

var mapCameraModelString = map[CameraModel]string{}

func (cm CameraModel) String() string {
	if str, ok := mapCameraModelString[cm]; ok {
		return str
	}
	return ""
}

// Canon Camera Models
const (
	SonyModelUnknown CameraModel = iota + 0x40000
)
