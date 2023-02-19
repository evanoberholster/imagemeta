package mknote

// Apple Camera Models
const (
	UnknownAppleModel CameraModel = iota + 0x20000
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

	//TODO: add iPad Models
)

var mapAppleCameraModel = map[string]CameraModel{}
