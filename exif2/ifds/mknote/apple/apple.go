package apple

// CameraModel is an Apple Camera Model found in Exif
type CameraModel uint32

// Apple Camera Models
const (
	AppleModelUnknown CameraModel = iota + 0x20000
	iPhone
	iPhone3G
	iPhone3GS
	iPhone4
	iPhone4S
	iPhone5
	iPhone5c
	iPhone5s
	iPhone6
	iPhone6Plus
	iPhone6s
	iPhone6sPlus
	iPhoneSE //(1st generation)
	iPhone7
	iPhone7Plus
	iPhone8
	iPhone8Plus
	iPhoneX
	iPhoneXR
	iPhoneXS
	iPhoneXSMax
	iPhone11
	iPhone11Pro
	iPhone11ProMax
	iPhoneSE2nd // (2nd generation)
	iPhone12mini
	iPhone12
	iPhone12Pro
	iPhone12ProMax
	iPhone13mini
	iPhone13
	iPhone13Pro
	iPhone13ProMax
	iPhoneSE3rd // (3rd generation)
	iPhone14
	iPhone14Plus
	iPhone14Pro
	iPhone14ProMax

	//TODO: add iPad and iPod Touch Models
	iPodtouch
	iPad
	iPad2
	iPadAir
	iPadmini
)

// CameraModelFromString returns a apple camera model from the given string
func CameraModelFromString(str string) (CameraModel, bool) {
	if cm, ok := mapAppleCameraModel[str]; ok {
		return cm, true
	}
	return AppleModelUnknown, false
}

func (cm CameraModel) String() string {
	if str, ok := mapAppleCameraModelString[cm]; ok {
		return str
	}
	return ""
}

var mapAppleCameraModel = map[string]CameraModel{
	"UnknownAppleModel": AppleModelUnknown,
	"iPhone":            iPhone,
	"iPhone 3G":         iPhone3G,
	"iPhone 3GS":        iPhone3GS,
	"iPhone 4":          iPhone4,
	"iPhone 4S":         iPhone4S,
	"iPhone 5":          iPhone5,
	"iPhone 5c":         iPhone5c,
	"iPhone 5s":         iPhone5s,
	"iPhone 6":          iPhone6,
	"iPhone 6 Plus":     iPhone6Plus,
	"iPhone 6s":         iPhone6s,
	"iPhone 6s Plus":    iPhone6sPlus,
	"iPhone SE":         iPhoneSE,
	"iPhone 7":          iPhone7,
	"iPhone 7 Plus":     iPhone7Plus,
	"iPhone 8":          iPhone8,
	"iPhone 8 Plus":     iPhone8Plus,
	"iPhone X":          iPhoneX,
	"iPhone XR":         iPhoneXR,
	"iPhone XS":         iPhoneXS,
	"iPhone XS Max":     iPhoneXSMax,
	"iPhone 11":         iPhone11,
	"iPhone 11 Pro":     iPhone11Pro,
	"iPhone 11 Pro Max": iPhone11ProMax,
	"iPhone SE 2nd":     iPhoneSE2nd,
	"iPhone 12 mini":    iPhone12mini,
	"iPhone 12":         iPhone12,
	"iPhone 12 Pro":     iPhone12Pro,
	"iPhone 12 Pro Max": iPhone12ProMax,
	"iPhone 13 mini":    iPhone13mini,
	"iPhone 13":         iPhone13,
	"iPhone 13 Pro":     iPhone13Pro,
	"iPhone 13 Pro Max": iPhone13ProMax,
	"iPhone SE 3rd":     iPhoneSE3rd,
	"iPhone 14":         iPhone14,
	"iPhone 14 Plus":    iPhone14Plus,
	"iPhone 14 Pro":     iPhone14Pro,
	"iPhone 14 Pro Max": iPhone14ProMax,
	// iPad Models
	"iPod touch": iPodtouch,
	"iPad":       iPad,
	"iPad 2":     iPad2,
	"iPad Air":   iPadAir,
	"iPad mini":  iPadmini,
}

var mapAppleCameraModelString = map[CameraModel]string{
	iPhone:         "iPhone",
	iPhone3G:       "iPhone 3G",
	iPhone3GS:      "iPhone 3GS",
	iPhone4:        "iPhone 4",
	iPhone4S:       "iPhone 4S",
	iPhone5:        "iPhone 5",
	iPhone5c:       "iPhone 5c",
	iPhone5s:       "iPhone 5s",
	iPhone6:        "iPhone 6",
	iPhone6Plus:    "iPhone 6 Plus",
	iPhone6s:       "iPhone 6s",
	iPhone6sPlus:   "iPhone 6s Plus",
	iPhoneSE:       "iPhone SE",
	iPhone7:        "iPhone 7",
	iPhone7Plus:    "iPhone 7 Plus",
	iPhone8:        "iPhone 8",
	iPhone8Plus:    "iPhone 8 Plus",
	iPhoneX:        "iPhone X",
	iPhoneXR:       "iPhone XR",
	iPhoneXS:       "iPhone XS",
	iPhoneXSMax:    "iPhone XS Max",
	iPhone11:       "iPhone 11",
	iPhone11Pro:    "iPhone 11 Pro",
	iPhone11ProMax: "iPhone 11 Pro Max",
	iPhoneSE2nd:    "iPhone SE 2nd",
	iPhone12mini:   "iPhone 12 mini",
	iPhone12:       "iPhone 12",
	iPhone12Pro:    "iPhone 12 Pro",
	iPhone12ProMax: "iPhone 12 Pro Max",
	iPhone13mini:   "iPhone 13 mini",
	iPhone13:       "iPhone 13",
	iPhone13Pro:    "iPhone 13 Pro",
	iPhone13ProMax: "iPhone 13 Pro Max",
	iPhoneSE3rd:    "iPhone SE 3rd",
	iPhone14:       "iPhone 14",
	iPhone14Plus:   "iPhone 14 Plus",
	iPhone14Pro:    "iPhone 14 Pro",
	iPhone14ProMax: "iPhone 14 Pro Max",
	iPodtouch:      "iPod touch",
	iPad:           "iPad",
	iPad2:          "iPad 2",
	iPadAir:        "iPad Air",
	iPadmini:       "iPad mini",
}
