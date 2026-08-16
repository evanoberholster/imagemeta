package sony

import (
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// ParseCameraInfo3 parses maker-note tag 0x0010 CameraInfo3 (modern layout)
// from raw bytes. Expected payload length is at least 0x80 bytes.
func ParseCameraInfo3(raw []byte, bo utils.ByteOrder) SonyCameraInfo3 {
	var dst SonyCameraInfo3
	dst.LensSpec = DisplayText(bytesAt(raw, 0x0000, 8))
	dst.FocalLength = u16At(raw, bo, 0x000e)
	dst.FocalLengthTeleZoom = u16At(raw, bo, 0x0010)
	dst.FocusStatus = i16At(raw, bo, 0x0019)
	dst.AFPointSelected = u8At(raw, 0x001c)
	dst.FocusMode = u8At(raw, 0x001d)
	dst.AFPoint = u8At(raw, 0x0020)
	dst.AFStatusActiveSensor = i16At(raw, bo, 0x0021)
	decodeAFStatus15(&dst.AFStatus15, raw, bo, 0x0023)
	return dst
}

// ParseCameraInfo2 parses maker-note tag 0x0010 CameraInfo2 (legacy layout,
// models A700/A850/A900). Expected payload length is at least 0x80 bytes.
func ParseCameraInfo2(raw []byte, bo utils.ByteOrder) SonyCameraInfo2 {
	var dst SonyCameraInfo2
	dst.LensSpec = DisplayText(bytesAt(raw, 0x0000, 8))
	dst.AFPointSelected = u8At(raw, 0x0014)
	dst.FocusModeSetting = u8At(raw, 0x0015)
	dst.AFPoint = u8At(raw, 0x0018)
	dst.AFStatusActiveSensor = i16At(raw, bo, 0x001b)
	dst.AFStatusTopRight = i16At(raw, bo, 0x001d)
	dst.AFStatusBottomRight = i16At(raw, bo, 0x001f)
	dst.AFStatusBottom = i16At(raw, bo, 0x0021)
	dst.AFStatusMiddleHorizontal = i16At(raw, bo, 0x0023)
	dst.AFStatusCenterVertical = i16At(raw, bo, 0x0025)
	dst.AFStatusTop = i16At(raw, bo, 0x0027)
	dst.AFStatusTopLeft = i16At(raw, bo, 0x0029)
	dst.AFStatusBottomLeft = i16At(raw, bo, 0x002b)
	dst.AFStatusLeft = i16At(raw, bo, 0x002d)
	dst.AFStatusCenterHorizontal = i16At(raw, bo, 0x002f)
	dst.AFStatusRight = i16At(raw, bo, 0x0031)
	return dst
}

// ParseFocusInfo parses maker-note tag 0x0020 FocusInfo (legacy layout) from
// raw bytes. UnitCount must be 19154 or 19148.
func ParseFocusInfo(raw []byte, bo utils.ByteOrder) SonyFocusInfo {
	var dst SonyFocusInfo
	dst.DriveMode2 = u8At(raw, 0x000e)
	dst.Rotation = u8At(raw, 0x0010)
	dst.ImageStabilizationSetting = u8At(raw, 0x0014)
	dst.DynamicRangeOptimizerMode = u8At(raw, 0x0015)
	dst.BracketShotNumber = u8At(raw, 0x002b)
	dst.WhiteBalanceBracketing = u8At(raw, 0x002c)
	dst.BracketShotNumber2 = u8At(raw, 0x002d)
	dst.DynamicRangeOptimizerBracket = u8At(raw, 0x002e)
	dst.ExposureBracketShotNumber = u8At(raw, 0x002f)
	dst.ExposureProgram = u8At(raw, 0x003f)
	dst.CreativeStyle = u8At(raw, 0x0041)
	dst.ISOSetting = u8At(raw, 0x006d)
	dst.ISO = u8At(raw, 0x006f)
	dst.DynamicRangeOptimizerMode2 = u8At(raw, 0x0077)
	dst.DynamicRangeOptimizerLevel = u8At(raw, 0x0079)
	dst.FocusPosition = u8At(raw, 0x09bb)
	return dst
}

// ParseMoreSettings parses the nested 0x0001 sub-block inside Sony MoreInfo
// (tag 0x0020). Slice length should be 256 bytes.
func ParseMoreSettings(raw []byte, bo utils.ByteOrder) SonyMoreSettings {
	var dst SonyMoreSettings
	dst.DriveMode2 = u8At(raw, 0x0001)
	dst.ExposureProgram = u8At(raw, 0x0002)
	dst.MeteringMode = u8At(raw, 0x0003)
	dst.DynamicRangeOptimizerSetting = u8At(raw, 0x0004)
	dst.DynamicRangeOptimizerLevel = u8At(raw, 0x0005)
	dst.ColorSpace = u8At(raw, 0x0006)
	dst.CreativeStyleSetting = u8At(raw, 0x0007)
	dst.ContrastSetting = i8At(raw, 0x0008)
	dst.SaturationSetting = i8At(raw, 0x0009)
	dst.SharpnessSetting = i8At(raw, 0x000a)
	dst.WhiteBalanceSetting = u8At(raw, 0x000d)
	dst.ColorTemperatureSetting = u8At(raw, 0x000e)
	dst.ColorCompensationFilterSet = i8At(raw, 0x000f)
	dst.FlashMode = u8At(raw, 0x0010)
	dst.LongExposureNoiseReduction = u8At(raw, 0x0011)
	dst.HighISONoiseReduction = u8At(raw, 0x0012)
	dst.FocusMode = u8At(raw, 0x0013)
	dst.MultiFrameNoiseReduction = u8At(raw, 0x0015)
	dst.HDRSetting = u8At(raw, 0x0016)
	dst.HDRLevel = u8At(raw, 0x0017)
	dst.ViewingMode = u8At(raw, 0x0018)
	dst.FaceDetection = u8At(raw, 0x0019)
	dst.CustomWB_RBLevels[0] = u16RevAt(raw, bo, 0x001a)
	dst.CustomWB_RBLevels[1] = u16RevAt(raw, bo, 0x001c)
	dst.BrightnessValue = u8At(raw, 0x001d)
	dst.ExposureCompensationSet = u8At(raw, 0x001e)
	dst.FlashExposureCompSet = u8At(raw, 0x001f)
	dst.LiveViewAFMethod = u8At(raw, 0x0020)
	dst.ISO = u8At(raw, 0x0025)
	dst.FNumber = u8At(raw, 0x0026)
	dst.ExposureTime = u8At(raw, 0x0027)
	dst.FocalLength2 = u8At(raw, 0x0029)
	dst.ExposureCompensation2 = i16At(raw, bo, 0x002a)
	dst.FlashExposureCompSet2 = i16At(raw, bo, 0x002c)
	dst.Orientation2 = u8At(raw, 0x002e)
	dst.FocusPosition2 = u8At(raw, 0x002f)
	dst.FlashAction = u8At(raw, 0x0030)
	dst.FocusMode2 = u8At(raw, 0x0032)
	dst.FlashActionExternal = u8At(raw, 0x007c)
	dst.FlashStatus = u8At(raw, 0x0086)
	return dst
}

// ParseMoreInfo parses maker-note tag 0x0020 MoreInfo (modern layout) from
// raw bytes. UnitCount must be 20480.
func ParseMoreInfo(raw []byte, bo utils.ByteOrder) SonyMoreInfo {
	var dst SonyMoreInfo
	dst.MoreSettings = ParseMoreSettings(bytesAt(raw, 0x0001, 256), bo)
	dst.FaceInfo = ParseFaceInfo(bytesAt(raw, 0x0002, 256), bo)
	dst.ImageCount = u32At(raw, bo, 0x0201+0x011b)
	dst.ShutterCount = u32At(raw, bo, 0x0201+0x0125)
	dst.ShotNumberSincePowerUp = u32At(raw, bo, 0x0401+0x044e)
	return dst
}

// ParseCameraSettings parses maker-note tag 0x0114 CameraSettings (legacy
// 280-byte or 332-byte layout) from raw bytes. Dispatches to the correct
// internal parser based on payload length.
func ParseCameraSettings(raw []byte, bo utils.ByteOrder) SonyCameraSettings {
	if len(raw) >= 332 {
		return parseCameraSettings332(raw, bo)
	}
	return parseCameraSettings280(raw, bo)
}

func parseCameraSettings280(raw []byte, bo utils.ByteOrder) SonyCameraSettings {
	var dst SonyCameraSettings
	dst.ExposureTime = u16At(raw, bo, 0x0000)
	dst.FNumber = u16At(raw, bo, 0x0001)
	dst.HighSpeedSync = u16At(raw, bo, 0x0002)
	dst.ExposureCompensationSet = u16At(raw, bo, 0x0003)
	dst.DriveMode = u16At(raw, bo, 0x0004)
	dst.WhiteBalanceSetting = u16At(raw, bo, 0x0005)
	dst.WhiteBalanceFineTune = i16At(raw, bo, 0x0006)
	dst.ColorTemperatureSet = u16At(raw, bo, 0x0007)
	dst.ColorCompensationFilterSet = i16At(raw, bo, 0x0008)
	fillU16s(dst.CustomWB_RGBLevels[:], raw, bo, 0x0009, false)
	dst.ColorTemperatureCustom = u16At(raw, bo, 0x000c)
	dst.ColorCompensationFilterCustom = i16At(raw, bo, 0x000d)
	dst.WhiteBalance = u16At(raw, bo, 0x000f)
	dst.FocusModeSetting = u16At(raw, bo, 0x0010)
	dst.AFAreaMode = u16At(raw, bo, 0x0011)
	dst.AFPointSetting = u16At(raw, bo, 0x0012)
	dst.FlashMode = u16At(raw, bo, 0x0013)
	dst.FlashExposureCompSet = u16At(raw, bo, 0x0014)
	dst.MeteringMode = u16At(raw, bo, 0x0015)
	dst.ISOSetting = u16At(raw, bo, 0x0016)
	dst.DynamicRangeOptimizerMode = u16At(raw, bo, 0x0018)
	dst.DynamicRangeOptimizerLevel = u16At(raw, bo, 0x0019)
	dst.CreativeStyle = u16At(raw, bo, 0x001a)
	dst.ColorSpace = u16At(raw, bo, 0x001b)
	dst.Sharpness = i16At(raw, bo, 0x001c)
	dst.Contrast = i16At(raw, bo, 0x001d)
	dst.Saturation = i16At(raw, bo, 0x001e)
	dst.ZoneMatchingValue = u16At(raw, bo, 0x001f)
	dst.Brightness = i16At(raw, bo, 0x0022)
	dst.FlashControl = u16At(raw, bo, 0x0023)
	dst.PrioritySetupShutterRelease = u16At(raw, bo, 0x0028)
	dst.AFIlluminator = u16At(raw, bo, 0x0029)
	dst.AFWithShutter = u16At(raw, bo, 0x002a)
	dst.LongExposureNoiseReduction = u16At(raw, bo, 0x002b)
	dst.HighISONoiseReduction = u16At(raw, bo, 0x002c)
	dst.ImageStyle = u16At(raw, bo, 0x002d)
	dst.FocusModeSwitch = u16At(raw, bo, 0x002e)
	dst.ShutterSpeedSetting = u16At(raw, bo, 0x002f)
	dst.ApertureSetting = u16At(raw, bo, 0x0030)
	dst.ExposureProgram = u16At(raw, bo, 0x003c)
	dst.ImageStabilizationSetting = u16At(raw, bo, 0x003d)
	dst.FlashAction = u16At(raw, bo, 0x003e)
	dst.Rotation = u16At(raw, bo, 0x003f)
	dst.AELock = u16At(raw, bo, 0x0040)
	dst.FlashAction2 = u16At(raw, bo, 0x004c)
	dst.FocusMode = u16At(raw, bo, 0x004d)
	dst.BatteryState = u16At(raw, bo, 0x0050)
	dst.BatteryLevel = u16At(raw, bo, 0x0051)
	dst.FocusStatus = u16At(raw, bo, 0x0053)
	dst.SonyImageSize = u16At(raw, bo, 0x0054)
	dst.AspectRatio = u16At(raw, bo, 0x0055)
	dst.Quality = u16At(raw, bo, 0x0056)
	dst.ExposureLevelIncrements = u16At(raw, bo, 0x0058)
	dst.RedEyeReduction = u16At(raw, bo, 0x006a)
	return dst
}

func parseCameraSettings332(raw []byte, bo utils.ByteOrder) SonyCameraSettings {
	var dst SonyCameraSettings
	dst.ExposureTime = u16At(raw, bo, 0x0000)
	dst.FNumber = u16At(raw, bo, 0x0001)
	dst.HighSpeedSync = u16At(raw, bo, 0x0002)
	dst.ExposureCompensationSet = u16At(raw, bo, 0x0003)
	dst.WhiteBalanceSetting = u16At(raw, bo, 0x0004)
	dst.WhiteBalanceFineTune = i16At(raw, bo, 0x0005)
	dst.ColorTemperatureSet = u16At(raw, bo, 0x0006)
	dst.ColorCompensationFilterSet = i16At(raw, bo, 0x0007)
	fillU16s(dst.CustomWB_RGBLevels[:], raw, bo, 0x0008, false)
	dst.ColorTemperatureCustom = u16At(raw, bo, 0x000b)
	dst.ColorCompensationFilterCustom = i16At(raw, bo, 0x000c)
	dst.WhiteBalance = u16At(raw, bo, 0x000e)
	dst.FocusModeSetting = u16At(raw, bo, 0x000f)
	dst.AFAreaMode = u16At(raw, bo, 0x0010)
	dst.AFPointSetting = u16At(raw, bo, 0x0011)
	dst.FlashExposureCompSet = u16At(raw, bo, 0x0012)
	dst.MeteringMode = u16At(raw, bo, 0x0013)
	dst.ISOSetting = u16At(raw, bo, 0x0014)
	dst.DynamicRangeOptimizerMode = u16At(raw, bo, 0x0016)
	dst.DynamicRangeOptimizerLevel = u16At(raw, bo, 0x0017)
	dst.CreativeStyle = u16At(raw, bo, 0x0018)
	dst.Sharpness = i16At(raw, bo, 0x0019)
	dst.Contrast = i16At(raw, bo, 0x001a)
	dst.Saturation = i16At(raw, bo, 0x001b)
	dst.FlashControl = u16At(raw, bo, 0x001f)
	dst.LongExposureNoiseReduction = u16At(raw, bo, 0x0025)
	dst.HighISONoiseReduction = u16At(raw, bo, 0x0026)
	dst.ImageStyle = u16At(raw, bo, 0x0027)
	dst.ShutterSpeedSetting = u16At(raw, bo, 0x0028)
	dst.ApertureSetting = u16At(raw, bo, 0x0029)
	dst.ExposureProgram = u16At(raw, bo, 0x003c)
	dst.ImageStabilizationSetting = u16At(raw, bo, 0x003d)
	dst.FlashAction = u16At(raw, bo, 0x003e)
	dst.Rotation = u16At(raw, bo, 0x003f)
	dst.AELock = u16At(raw, bo, 0x0040)
	dst.FlashAction2 = u16At(raw, bo, 0x004c)
	dst.FocusMode = u16At(raw, bo, 0x004d)
	dst.FocusStatus = u16At(raw, bo, 0x0053)
	dst.SonyImageSize = u16At(raw, bo, 0x0054)
	dst.AspectRatio = u16At(raw, bo, 0x0055)
	dst.Quality = u16At(raw, bo, 0x0056)
	dst.ExposureLevelIncrements = u16At(raw, bo, 0x0058)
	dst.DriveMode = u16At(raw, bo, 0x007e)
	dst.FlashMode = u16At(raw, bo, 0x007f)
	dst.ColorSpace = u16At(raw, bo, 0x0083)
	return dst
}

// ParseCameraSettings3 parses maker-note tag 0x0114 CameraSettings3 (modern
// 1-byte layout, 1536 or 2048 byte payload) from raw bytes.
func ParseCameraSettings3(raw []byte, bo utils.ByteOrder) SonyCameraSettings3 {
	var dst SonyCameraSettings3
	dst.ShutterSpeedSetting = u8At(raw, 0x0000)
	dst.ApertureSetting = u8At(raw, 0x0001)
	dst.ISOSetting = u8At(raw, 0x0002)
	dst.ExposureCompensationSet = u8At(raw, 0x0003)
	dst.DriveModeSetting = u8At(raw, 0x0004)
	dst.ExposureProgram = u8At(raw, 0x0005)
	dst.FocusModeSetting = u8At(raw, 0x0006)
	dst.MeteringMode = u8At(raw, 0x0007)
	dst.SonyImageSize = u8At(raw, 0x0009)
	dst.AspectRatio = u8At(raw, 0x000a)
	dst.Quality = u8At(raw, 0x000b)
	dst.DynamicRangeOptimizerSetting = u8At(raw, 0x000c)
	dst.DynamicRangeOptimizerLevel = u8At(raw, 0x000d)
	dst.ColorSpace = u8At(raw, 0x000e)
	dst.CreativeStyleSetting = u8At(raw, 0x000f)
	dst.ContrastSetting = i8At(raw, 0x0010)
	dst.SaturationSetting = i8At(raw, 0x0011)
	dst.SharpnessSetting = i8At(raw, 0x0012)
	dst.WhiteBalanceSetting = u8At(raw, 0x0016)
	dst.ColorTemperatureSetting = u8At(raw, 0x0017)
	dst.ColorCompensationFilterSet = i8At(raw, 0x0018)
	dst.CustomWB_RGBLevels[0] = u16RevAt(raw, bo, 0x0019)
	dst.CustomWB_RGBLevels[1] = u16RevAt(raw, bo, 0x001b)
	dst.CustomWB_RGBLevels[2] = u16RevAt(raw, bo, 0x001d)
	dst.FlashMode = u8At(raw, 0x0020)
	dst.FlashControl = u8At(raw, 0x0021)
	dst.FlashExposureCompSet = u8At(raw, 0x0023)
	dst.AFAreaMode = u8At(raw, 0x0024)
	dst.LongExposureNoiseReduction = u8At(raw, 0x0025)
	dst.HighISONoiseReduction = u8At(raw, 0x0026)
	dst.SmileShutterMode = u8At(raw, 0x0027)
	dst.RedEyeReduction = u8At(raw, 0x0028)
	dst.HDRSetting = u8At(raw, 0x002d)
	dst.HDRLevel = u8At(raw, 0x002e)
	dst.ViewingMode = u8At(raw, 0x002f)
	dst.FaceDetection = u8At(raw, 0x0030)
	dst.SmileShutter = u8At(raw, 0x0031)
	dst.SweepPanoramaSize = u8At(raw, 0x0032)
	dst.SweepPanoramaDirection = u8At(raw, 0x0033)
	dst.DriveMode = u8At(raw, 0x0034)
	dst.MultiFrameNoiseReduction = u8At(raw, 0x0035)
	dst.LiveViewAFSetting = u8At(raw, 0x0036)
	dst.PanoramaSize3D = u8At(raw, 0x0038)
	dst.AFButtonPressed = u8At(raw, 0x0083)
	dst.LiveViewMetering = u8At(raw, 0x0084)
	dst.ViewingMode2 = u8At(raw, 0x0085)
	dst.AELock = u8At(raw, 0x0086)
	dst.FlashStatusBuiltIn = u8At(raw, 0x0087)
	dst.FlashStatusExternal = u8At(raw, 0x0088)
	dst.LiveViewFocusMode = u8At(raw, 0x008b)
	dst.LensMount = u8At(raw, 0x0099)
	dst.SequenceNumber = u8At(raw, 0x010c)
	if v := u32At(raw, bo, 0x0114); v != 0 {
		if folderNumber, ok := meta.SafecastUint32ToUint16((v & 0xffc000) >> 14); ok {
			dst.FolderNumber = folderNumber
		}
		if imageNumber, ok := meta.SafecastUint32ToUint16(v & 0x3fff); ok {
			dst.ImageNumber = imageNumber
		}
	}
	dst.ShotNumberSincePowerUp2 = u32At(raw, bo, 0x0200)
	return dst
}

// ParseShotInfo parses maker-note tag 0x3000 ShotInfo from raw bytes.
// Expected payload length is at least 0x40 bytes.
func ParseShotInfo(raw []byte, bo utils.ByteOrder) SonyShotInfo {
	var dst SonyShotInfo
	dst.FaceInfoOffset = u16At(raw, bo, 0x0002)
	dst.SonyDateTime = DisplayText(bytesAt(raw, 0x0006, 20))
	dst.SonyImageHeight = u16At(raw, bo, 0x001a)
	dst.SonyImageWidth = u16At(raw, bo, 0x001c)
	dst.FacesDetected = u16At(raw, bo, 0x0030)
	dst.FaceInfoLength = u16At(raw, bo, 0x0032)
	dst.MetaVersion = DisplayText(bytesAt(raw, 0x0034, 16))
	return dst
}

// ParseFaceInfo parses the face-detection sub-block inside Sony ShotInfo
// (tag 0x3000) from raw bytes.
func ParseFaceInfo(raw []byte, bo utils.ByteOrder) SonyFaceInfo {
	var dst SonyFaceInfo
	dst.FacesDetected = u16At(raw, bo, 0x0000)
	return dst
}

// ParseTag9400 parses maker-note tag 0x9400 from raw bytes.
func ParseTag9400(raw []byte, bo utils.ByteOrder) SonyTag9400 {
	var dst SonyTag9400
	dst.SequenceImageNumber = u32At(raw, bo, 0x0008)
	dst.SequenceFileNumber = u32At(raw, bo, 0x000c)
	dst.ReleaseMode2 = u8At(raw, 0x0010)
	dst.ShotNumberSincePowerUp = u32At(raw, bo, 0x001a)
	dst.SequenceLength = u8At(raw, 0x0022)
	dst.CameraOrientation = u8At(raw, 0x0028)
	dst.Quality2 = u8At(raw, 0x0029)
	dst.SonyImageHeight = u16At(raw, bo, 0x0044)
	dst.ModelReleaseYear = u8At(raw, 0x0052)
	return dst
}

// ParseTag9404 parses maker-note tag 0x9404 from raw bytes.
func ParseTag9404(raw []byte, bo utils.ByteOrder) SonyTag9404 {
	var dst SonyTag9404
	dst.ExposureProgram = u8At(raw, 0x000b)
	dst.IntelligentAuto = u8At(raw, 0x000d)
	return dst
}

// ParseTag9405 parses maker-note tag 0x9405 (lens correction parameters)
// from raw bytes.
func ParseTag9405(raw []byte, bo utils.ByteOrder) SonyTag9405 {
	var dst SonyTag9405
	dst.DistortionCorrParamsPresent = u8At(raw, 0x0600)
	dst.DistortionCorrection = u8At(raw, 0x0601)
	dst.LensFormat = u8At(raw, 0x0603)
	dst.LensMount = u8At(raw, 0x0604)
	dst.LensType = u16At(raw, bo, 0x0608)
	fillI16s(dst.VignettingCorrParams[:], raw, bo, 0x064a)
	fillI16s(dst.ChromaticAberrationCorrParams[:], raw, bo, 0x066a)
	fillI16s(dst.DistortionCorrParams[:], raw, bo, 0x06ca)
	return dst
}

// ParseTag9406 parses maker-note tag 0x9406 (battery information) from raw bytes.
func ParseTag9406(raw []byte, bo utils.ByteOrder) SonyTag9406 {
	var dst SonyTag9406
	dst.BatteryTemperature = u8At(raw, 0x0005)
	dst.BatteryLevelGrip1 = u8At(raw, 0x0006)
	dst.BatteryLevel = u8At(raw, 0x0007)
	dst.BatteryLevelGrip2 = u8At(raw, 0x0008)
	return dst
}

// ParseTag940A parses maker-note tag 0x940a (AF points selected) from raw bytes.
func ParseTag940A(raw []byte, bo utils.ByteOrder) SonyTag940A {
	var dst SonyTag940A
	dst.AFPointsSelected = u32At(raw, bo, 0x0004)
	return dst
}

// ParseTag940C parses maker-note tag 0x940c (lens mount and version info)
// from raw bytes.
func ParseTag940C(raw []byte, bo utils.ByteOrder) SonyTag940C {
	var dst SonyTag940C
	dst.LensMount2 = u8At(raw, 0x0008)
	dst.LensType3 = u16At(raw, bo, 0x0009)
	dst.CameraEMountVersion = u16At(raw, bo, 0x000b)
	dst.LensEMountVersion = u16At(raw, bo, 0x000d)
	dst.LensFirmwareVersion = u16At(raw, bo, 0x0014)
	return dst
}

// ParseTag2010 parses maker-note tag 0x2010 (modern settings block) from
// raw bytes. Expected payload length is at least 0x1220 bytes.
func ParseTag2010(raw []byte, bo utils.ByteOrder) SonyTag2010 {
	var dst SonyTag2010
	dst.SequenceImageNumber = u32At(raw, bo, 0x0000)
	dst.SequenceFileNumber = u32At(raw, bo, 0x0004)
	dst.ReleaseMode2 = u32At(raw, bo, 0x0008)
	dst.SonyDateTime = DisplayText(bytesAt(raw, 0x01b6, 7))
	dst.DynamicRangeOptimizer = u8At(raw, 0x0324)
	dst.ReleaseMode3 = u8At(raw, 0x1128)
	dst.SelfTimer = u8At(raw, 0x1134)
	dst.FlashMode = u8At(raw, 0x1138)
	dst.StopsAboveBaseISO = u16At(raw, bo, 0x113e)
	dst.BrightnessValue = u16At(raw, bo, 0x1140)
	dst.HDRSetting = u8At(raw, 0x1148)
	dst.ExposureCompensation = i16At(raw, bo, 0x114c)
	dst.PictureProfile = u8At(raw, 0x1162)
	dst.PictureEffect2 = u8At(raw, 0x1167)
	dst.Quality2 = u8At(raw, 0x1174)
	dst.MeteringMode = u8At(raw, 0x1178)
	dst.ExposureProgram = u8At(raw, 0x1179)
	fillU16s(dst.WB_RGBLevels[:], raw, bo, 0x1180, false)
	dst.SonyISO = u16At(raw, bo, 0x1218)
	return dst
}

// ParseTag202A parses maker-note tag 0x202a (focal plane AF points) from raw bytes.
func ParseTag202A(raw []byte, bo utils.ByteOrder) SonyTag202A {
	var dst SonyTag202A
	dst.FocalPlaneAFPointsUsed = u8At(raw, 0x0001)
	return dst
}

// ParseHiddenInfo parses maker-note tag 0x2044 (hidden data offset/length)
// from raw bytes.
func ParseHiddenInfo(raw []byte, bo utils.ByteOrder) SonyHiddenInfo {
	var dst SonyHiddenInfo
	dst.HiddenDataOffset = u32At(raw, bo, 0x0000)
	dst.HiddenDataLength = u32At(raw, bo, 0x0004)
	return dst
}

// ParseTag9050 parses maker-note tag 0x9050 (shutter/exposure/serial data)
// from raw bytes. Expected payload length is at least 0x200 bytes.
func ParseTag9050(raw []byte, bo utils.ByteOrder) SonyTag9050 {
	var dst SonyTag9050
	dst.SonyMaxAperture = u8At(raw, 0x0000)
	dst.SonyMinAperture = u8At(raw, 0x0001)
	dst.Shutter[0] = u16At(raw, bo, 0x0020)
	dst.Shutter[1] = u16At(raw, bo, 0x0022)
	dst.Shutter[2] = u16At(raw, bo, 0x0024)
	dst.FlashStatus = u8At(raw, 0x0031)
	dst.ShutterCount = u32At(raw, bo, 0x0032)
	dst.SonyExposureTime = u16At(raw, bo, 0x003a)
	dst.SonyFNumber = u16At(raw, bo, 0x003c)
	dst.ReleaseMode2 = u8At(raw, 0x003f)
	copy(dst.InternalSerialNumber[:], bytesAt(raw, 0x00f0, 5))
	dst.LensMount = u8At(raw, 0x0105)
	dst.LensFormat = u8At(raw, 0x0106)
	dst.LensType = u16At(raw, bo, 0x0109)
	dst.DistortionCorrParamsPresent = u8At(raw, 0x010b)
	dst.LensSpecFeatures = DisplayText(bytesAt(raw, 0x0115, 2))
	dst.ShutterCount3 = u32At(raw, bo, 0x01bd)
	return dst
}

// ParseTag9416 parses maker-note tag 0x9416 (detailed EXIF-like data)
// from raw bytes. Expected payload length is at least 0x0a00 bytes.
func ParseTag9416(raw []byte, bo utils.ByteOrder) SonyTag9416 {
	var dst SonyTag9416
	dst.SonyISO = u16At(raw, bo, 0x0004)
	dst.StopsAboveBaseISO = u16At(raw, bo, 0x0006)
	dst.SonyExposureTime2 = u16At(raw, bo, 0x000a)
	if v := u32At(raw, bo, 0x000c); v != 0 {
		dst.ExposureTime = meta.ExposureTime(float32(v) / 1000000.0)
	}
	dst.SonyFNumber2 = u16At(raw, bo, 0x0010)
	dst.SonyMaxApertureValue = u16At(raw, bo, 0x0012)
	dst.SequenceImageNumber = u32At(raw, bo, 0x001d)
	dst.ReleaseMode2 = u8At(raw, 0x002b)
	copy(dst.InternalSerialNumber[:], bytesAt(raw, 0x0038, 6))
	dst.ExposureProgram = u8At(raw, 0x0035)
	dst.CreativeStyle = u8At(raw, 0x0037)
	dst.LensMount = u8At(raw, 0x0048)
	dst.LensFormat = u8At(raw, 0x0049)
	dst.LensType2 = u16At(raw, bo, 0x004b)
	fillI16s(dst.DistortionCorrParams[:], raw, bo, 0x004f)
	dst.PictureProfile = u8At(raw, 0x0070)
	dst.FocalLength = u16At(raw, bo, 0x0071)
	dst.MinFocalLength = u16At(raw, bo, 0x0073)
	dst.MaxFocalLength = u16At(raw, bo, 0x0075)
	fillI16s(dst.VignettingCorrParams[:], raw, bo, 0x089d)
	dst.APSCSizeCapture = u8At(raw, 0x08e5)
	fillI16s(dst.ChromaticAberrationCorrParams[:], raw, bo, 0x0945)
	return dst
}

// ParseAFInfo parses maker-note tag 0x940e (AF information) from raw bytes.
// Expected payload length is at least 0x0180 bytes.
func ParseAFInfo(raw []byte, bo utils.ByteOrder) SonyAFInfo {
	var dst SonyAFInfo
	dst.AFType = u8At(raw, 0x0002)
	dst.AFStatusActiveSensor = i16At(raw, bo, 0x0004)
	dst.AFPoint = u8At(raw, 0x0007)
	dst.AFPointInFocus = u8At(raw, 0x0008)
	dst.AFPointAtShutterRelease = u8At(raw, 0x0009)
	dst.AFAreaMode = u8At(raw, 0x000a)
	dst.FocusMode = u8At(raw, 0x000b)
	decodeAFStatus15(&dst.AFStatus15, raw, bo, 0x0011)
	dst.AFPointsUsed = u32At(raw, bo, 0x016e)
	dst.AFMicroAdj = i8At(raw, 0x017d)
	dst.ExposureProgram = u8At(raw, 0x017e)
	return dst
}

// decodeAFStatus15 fills the 18-point AF status matrix from raw bytes
// starting at the specified offset.
func decodeAFStatus15(dst *SonyAFStatus15, raw []byte, bo utils.ByteOrder, off int) {
	dst.UpperLeft = i16At(raw, bo, off+0x00)
	dst.Left = i16At(raw, bo, off+0x02)
	dst.LowerLeft = i16At(raw, bo, off+0x04)
	dst.FarLeft = i16At(raw, bo, off+0x06)
	dst.TopHorizontal = i16At(raw, bo, off+0x08)
	dst.NearRight = i16At(raw, bo, off+0x0a)
	dst.CenterHorizontal = i16At(raw, bo, off+0x0c)
	dst.NearLeft = i16At(raw, bo, off+0x0e)
	dst.BottomHorizontal = i16At(raw, bo, off+0x10)
	dst.TopVertical = i16At(raw, bo, off+0x12)
	dst.CenterVertical = i16At(raw, bo, off+0x14)
	dst.BottomVertical = i16At(raw, bo, off+0x16)
	dst.FarRight = i16At(raw, bo, off+0x18)
	dst.UpperRight = i16At(raw, bo, off+0x1a)
	dst.Right = i16At(raw, bo, off+0x1c)
	dst.LowerRight = i16At(raw, bo, off+0x1e)
	dst.UpperMiddle = i16At(raw, bo, off+0x20)
	dst.LowerMiddle = i16At(raw, bo, off+0x22)
}
