package exif

import (
	"math/bits"
	"strings"

	metacanon "github.com/evanoberholster/imagemeta/meta/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) warnCanonShortRead(t tag.Entry, parser string, got, want int) {
	if !r.warnEnabled() {
		return
	}
	r.warn().
		Str("parser", parser).
		Uint16("tagID", uint16(t.ID)).
		Str("tagName", t.Name()).
		Stringer("tagType", t.Type).
		Uint32("unitCount", t.UnitCount).
		Int("gotUnits", got).
		Int("wantUnits", want).
		Msg("canon maker-note payload too short")
}

func (r *Reader) parseCanonTag(t tag.Entry) bool {
	dst := &r.makerNoteInfo().Canon

	switch metacanon.MakerNoteTag(t.ID) {
	case metacanon.CanonImageType:
		dst.ImageType = r.parseStringAllowUndefined(t)
	case metacanon.CanonFirmwareVersion:
		dst.FirmwareVersion = r.parseStringAllowUndefined(t)
	case metacanon.CanonFocalLength:
		dst.FocalLength = r.parseCanonFocalLength(t)
	case metacanon.CanonFlashInfo:
		//dst.FlashInfo = r.parseCanonFlashInfo(t)
	case metacanon.FileNumber:
		dst.FileNumber = r.parseUint32(t)
	case metacanon.OwnerName:
		dst.OwnerName = r.parseStringAllowUndefined(t)
	case metacanon.SerialNumber:
		dst.SerialNumber = r.parseUint32(t)
	case metacanon.CanonCameraInfo:
		// intentionally not parsed
	case metacanon.CanonModelID:
		dst.ModelID = r.parseUint32(t)
	case metacanon.LensModel:
		dst.LensModel = r.parseStringAllowUndefined(t)
	case metacanon.CanonInternalSerialNumber:
		dst.InternalSerialNumber = r.parseStringAllowUndefined(t)
	case metacanon.CanonCameraSettings:
		dst.CameraSettings = r.parseCanonCameraSettings(t)
	case metacanon.CanonShotInfo:
		dst.ShotInfo = r.parseCanonShotInfo(t)
	case metacanon.CanonFileInfo:
		dst.FileInfo = r.parseCanonFileInfo(t)
	case metacanon.TimeInfo:
		dst.TimeInfo = r.parseCanonTimeInfo(t)
	case metacanon.BatteryType:
		dst.BatteryType = r.parseStringAllowUndefined(t)
	case metacanon.CanonAFInfo:
		dst.AFInfo = r.parseCanonAFInfo(t)
	case metacanon.CanonAFInfo2, metacanon.AFInfo3:
		dst.AFInfo = r.parseCanonAFInfo2(t)
	case metacanon.FaceDetect1:
		dst.FaceDetect1 = r.parseCanonFaceDetect1(t)
	case metacanon.FaceDetect2:
		dst.FaceDetect2 = r.parseCanonFaceDetect2(t)
	case metacanon.FaceDetect3:
		dst.FaceDetect3 = r.parseCanonFaceDetect3(t)
	case metacanon.ImageUniqueID:
		dst.ImageUniqueID = r.parseStringAllowUndefined(t)
	case metacanon.CanonCustomFunctions:
		// intentionally not parsed
	case metacanon.CanonAspectInfo:
		dst.AspectInfo = r.parseCanonAspectInfo(t)
	case metacanon.CanonProcessingInfo:
		dst.ProcessingInfo = r.parseCanonProcessingInfo(t)
	case metacanon.CanonColorSpace:
		dst.ColorSpace = r.parseUint16(t)
	case metacanon.CanonPreviewImageInfo:
		dst.PreviewImageInfo = r.parseCanonPreviewImageInfo(t)
	case metacanon.CanonSensorInfo:
		dst.SensorInfo = r.parseCanonSensorInfo(t)
	case metacanon.CanonPictureStyleUserDef:
		// intentionally not parsed
	case metacanon.CanonPictureStylePC:
		// intentionally not parsed
	case metacanon.CanonCustomPictureStyleFileName:
		dst.CustomPictureStyleFileName = r.parseStringAllowUndefined(t)
	case metacanon.CanonAFMicroAdj:
		dst.AFMicroAdj = r.parseCanonAFMicroAdj(t)
	case metacanon.CanonLightingOpt:
		dst.LightingOpt = r.parseCanonLightingOpt(t)
	case metacanon.CanonLensInfo:
		dst.LensInfo = r.parseCanonLensInfo(t)
	case metacanon.CanonMultiExp:
		dst.MultiExp = r.parseCanonMultiExp(t)
	case metacanon.CanonHDRInfo:
		dst.HDRInfo = r.parseCanonHDRInfo(t)
	case metacanon.CanonAFConfig:
		dst.AFConfig = r.parseCanonAFConfig(t)
	case metacanon.CanonRawBurstModeRoll:
		dst.RawBurstModeRoll = r.parseCanonRawBurstInfo(t)
	default:
		return false
	}
	return true
}

func (r *Reader) parseCanonUint16List(t tag.Entry, dst []uint16) int {
	if n := r.parseUint16List(t, dst); n > 0 {
		return n
	}
	return r.parseUndefinedUint16List(t, dst)
}

func (r *Reader) parseCanonInt32List(t tag.Entry, dst []int32) int {
	if n := r.parseInt32List(t, dst); n > 0 {
		return n
	}
	var u32 [2048]uint32
	if len(dst) > len(u32) {
		dst = dst[:len(u32)]
	}
	n := r.parseUint32List(t, u32[:len(dst)])
	for i := 0; i < n; i++ {
		dst[i] = int32(u32[i])
	}
	return n
}

func (r *Reader) parseCanonBlockPreview(t tag.Entry) metacanon.BlockPreview {
	dst := metacanon.BlockPreview{Size: t.Size()}
	if dst.Size == 0 {
		return dst
	}
	maxBytes := uint32(len(dst.Preview))
	if maxBytes > dst.Size {
		maxBytes = dst.Size
	}
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		n := int(maxBytes)
		copy(dst.Preview[:], r.state.buf[:n])
		dst.PreviewCount = uint8(n)
		return dst
	}
	buf, _, err := r.readTagBytes(t, maxBytes)
	if err != nil {
		if r.warnEnabled() {
			r.warn().
				Err(err).
				Str("parser", "parseCanonBlockPreview").
				Uint16("tagID", uint16(t.ID)).
				Str("tagName", t.Name()).
				Stringer("tagType", t.Type).
				Uint32("unitCount", t.UnitCount).
				Msg("failed reading canon maker-note payload")
		}
		return dst
	}
	if len(buf) == 0 {
		r.warnCanonShortRead(t, "parseCanonBlockPreview", 0, 1)
		return dst
	}
	n := len(buf)
	if n > len(dst.Preview) {
		n = len(dst.Preview)
	}
	copy(dst.Preview[:], buf[:n])
	dst.PreviewCount = uint8(n)
	return dst
}

func (r *Reader) parseCanonFlashInfo(t tag.Entry) metacanon.FlashInfo {
	return metacanon.FlashInfo{
		Raw: r.parseCanonBlockPreview(t),
	}
}

func (r *Reader) parseCanonPreviewImageInfo(t tag.Entry) metacanon.PreviewImageInfo {
	block := r.parseCanonBlockPreview(t)
	if block.PreviewCount < 24 {
		r.warnCanonShortRead(t, "parseCanonPreviewImageInfo", int(block.PreviewCount)/4, 6)
		return metacanon.PreviewImageInfo{}
	}
	return metacanon.PreviewImageInfo{
		PreviewQuality:     metacanon.Quality(int16(t.ByteOrder.Uint32(block.Preview[4:8]))),
		PreviewImageLength: t.ByteOrder.Uint32(block.Preview[8:12]),
		PreviewImageWidth:  t.ByteOrder.Uint32(block.Preview[12:16]),
		PreviewImageHeight: t.ByteOrder.Uint32(block.Preview[16:20]),
		PreviewImageStart:  t.ByteOrder.Uint32(block.Preview[20:24]),
	}
}

func (r *Reader) parseCanonSensorInfo(t tag.Entry) metacanon.SensorInfo {
	var raw [13]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 13 {
		r.warnCanonShortRead(t, "parseCanonSensorInfo", n, 13)
		return metacanon.SensorInfo{}
	}
	return metacanon.SensorInfo{
		SensorWidth:           int16(raw[1]),
		SensorHeight:          int16(raw[2]),
		SensorLeftBorder:      int16(raw[5]),
		SensorTopBorder:       int16(raw[6]),
		SensorRightBorder:     int16(raw[7]),
		SensorBottomBorder:    int16(raw[8]),
		BlackMaskLeftBorder:   int16(raw[9]),
		BlackMaskTopBorder:    int16(raw[10]),
		BlackMaskRightBorder:  int16(raw[11]),
		BlackMaskBottomBorder: int16(raw[12]),
	}
}

func (r *Reader) parseCanonAFConfig(t tag.Entry) metacanon.AFConfig {
	var raw [25]int32
	n := r.parseCanonInt32List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonAFConfig", n, 2)
		return metacanon.AFConfig{}
	}
	dst := metacanon.AFConfig{
		AFConfigTool: uint32(raw[1]) + 1,
	}
	if n > 2 {
		dst.AFTrackingSensitivity = raw[2]
	}
	if n > 3 {
		dst.AFAccelDecelTracking = raw[3]
	}
	if n > 4 {
		dst.AFPointSwitching = raw[4]
	}
	if n > 5 {
		dst.AIServoFirstImage = raw[5]
	}
	if n > 6 {
		dst.AIServoSecondImage = raw[6]
	}
	if n > 7 {
		dst.USMLensElectronicMF = raw[7]
	}
	if n > 8 {
		dst.AFAssistBeam = raw[8]
	}
	if n > 9 {
		dst.OneShotAFRelease = raw[9]
	}
	if n > 10 {
		dst.AutoAFPointSelEOSiTRAF = raw[10]
	}
	if n > 11 {
		dst.LensDriveWhenAFImpossible = raw[11]
	}
	if n > 12 {
		dst.SelectAFAreaSelectionMode = uint32(raw[12])
	}
	if n > 13 {
		dst.AFAreaSelectionMethod = raw[13]
	}
	if n > 14 {
		dst.OrientationLinkedAF = raw[14]
	}
	if n > 15 {
		dst.ManualAFPointSelPattern = raw[15]
	}
	if n > 16 {
		dst.AFPointDisplayDuringFocus = raw[16]
	}
	if n > 17 {
		dst.VFDisplayIllumination = raw[17]
	}
	if n > 18 {
		dst.AFStatusViewfinder = raw[18]
	}
	if n > 19 {
		dst.InitialAFPointInServo = raw[19]
	}
	if n > 20 {
		dst.SubjectToDetect = raw[20]
	}
	if n > 24 {
		dst.EyeDetection = raw[24]
	}
	return dst
}

func (r *Reader) parseCanonRawBurstInfo(t tag.Entry) metacanon.RawBurstInfo {
	var raw [3]uint32
	n := r.parseUint32List(t, raw[:])
	if n < 3 {
		r.warnCanonShortRead(t, "parseCanonRawBurstInfo", n, 3)
		return metacanon.RawBurstInfo{}
	}
	return metacanon.RawBurstInfo{
		RawBurstImageNum:   raw[1],
		RawBurstImageCount: raw[2],
	}
}

func (r *Reader) parseCanonFaceDetect1(t tag.Entry) metacanon.FaceDetect1Info {
	var raw [26]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 3 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect1", n, 3)
		return metacanon.FaceDetect1Info{}
	}
	dst := metacanon.FaceDetect1Info{
		FacesDetected: raw[2],
	}
	if n < 5 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect1", n, 5)
		return dst
	}
	dst.FaceDetectFrameSize[0] = raw[3]
	dst.FaceDetectFrameSize[1] = raw[4]

	faceCount := int(dst.FacesDetected)
	if faceCount > len(dst.FacePositions) {
		faceCount = len(dst.FacePositions)
	}
	for i := 0; i < faceCount; i++ {
		start := 8 + i*2
		if start+1 >= n {
			r.warnCanonShortRead(t, "parseCanonFaceDetect1", n, start+2)
			break
		}
		dst.FacePositions[i] = metacanon.FacePosition{
			X: int16(raw[start]),
			Y: int16(raw[start+1]),
		}
	}
	return dst
}

func (r *Reader) parseCanonFaceDetect2(t tag.Entry) metacanon.FaceDetect2Info {
	var raw [8]byte
	n := r.parseByteList(t, raw[:])
	if n < 3 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect2", n, 3)
		return metacanon.FaceDetect2Info{}
	}
	return metacanon.FaceDetect2Info{
		FaceWidth:     raw[1],
		FacesDetected: raw[2],
	}
}

func (r *Reader) parseCanonFaceDetect3(t tag.Entry) metacanon.FaceDetect3Info {
	var raw [8]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 4 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect3", n, 4)
		return metacanon.FaceDetect3Info{}
	}
	return metacanon.FaceDetect3Info{
		FacesDetected: raw[3],
	}
}

// parseCanonFocalLength parses tag 0x0002 (CanonFocalLength).
func (r *Reader) parseCanonFocalLength(t tag.Entry) metacanon.FocalLengthInfo {
	var raw [8]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonFocalLength", n, 4)
		return metacanon.FocalLengthInfo{}
	}
	return metacanon.FocalLengthInfo{
		FocalType:       raw[0],
		FocalLength:     raw[1],
		FocalPlaneXSize: raw[2],
		FocalPlaneYSize: raw[3],
	}
}

// parseCanonAspectInfo parses tag 0x009a (AspectInfo).
func (r *Reader) parseCanonAspectInfo(t tag.Entry) metacanon.AspectInfo {
	var raw [8]uint32
	if n := r.parseUint32List(t, raw[:]); n < 5 {
		r.warnCanonShortRead(t, "parseCanonAspectInfo", n, 5)
		return metacanon.AspectInfo{}
	}
	return metacanon.AspectInfo{
		AspectRatio:        raw[0],
		CroppedImageWidth:  raw[1],
		CroppedImageHeight: raw[2],
		CroppedImageLeft:   raw[3],
		CroppedImageTop:    raw[4],
	}
}

// parseCanonProcessingInfo parses tag 0x00a0 (ProcessingInfo).
func (r *Reader) parseCanonProcessingInfo(t tag.Entry) metacanon.ProcessingInfo {
	var raw [24]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 14 {
		r.warnCanonShortRead(t, "parseCanonProcessingInfo", n, 14)
		return metacanon.ProcessingInfo{}
	}
	return metacanon.ProcessingInfo{
		ToneCurve:            int16(raw[0]), // ExifTool Processing table uses FIRST_ENTRY=1.
		Sharpness:            int16(raw[1]),
		SharpnessFrequency:   int16(raw[2]),
		SensorRedLevel:       int16(raw[3]),
		SensorBlueLevel:      int16(raw[4]),
		WhiteBalanceRed:      int16(raw[5]),
		WhiteBalanceBlue:     int16(raw[6]),
		WhiteBalance:         int16(raw[7]),
		ColorTemperature:     int16(raw[8]),
		PictureStyle:         int16(raw[9]),
		DigitalGain:          int16(raw[10]),
		WBShiftAB:            int16(raw[11]),
		WBShiftGM:            int16(raw[12]),
		UnsharpMaskFineness:  int16(raw[13]),
		UnsharpMaskThreshold: int16(raw[14]),
	}

}

func (r *Reader) parseCanonPictureStyle3(dst *[3]uint16, count *uint8, t tag.Entry) {
	var raw [3]uint16
	n := r.parseCanonUint16List(t, raw[:])
	*dst = [3]uint16{}
	*count = 0
	if n == 0 {
		if t.UnitCount > 0 {
			r.warnCanonShortRead(t, "parseCanonPictureStyle3", n, 1)
		}
		return
	}
	copy(dst[:], raw[:n])
	*count = uint8(n)
}

// parseCanonAFMicroAdj parses tag 0x4013 (AFMicroAdj).
func (r *Reader) parseCanonAFMicroAdj(t tag.Entry) metacanon.AFMicroAdjInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonAFMicroAdj", n, 3)
		return metacanon.AFMicroAdjInfo{}
	}
	// ExifTool AFMicroAdj table uses FIRST_ENTRY=1.
	return metacanon.AFMicroAdjInfo{
		Mode:             raw[0],
		ValueNumerator:   raw[1],
		ValueDenominator: raw[2],
	}
}

// parseCanonLightingOpt parses tag 0x4018 (LightingOpt).
func (r *Reader) parseCanonLightingOpt(t tag.Entry) metacanon.LightingOptInfo {
	var raw [11]int32
	n := r.parseCanonInt32List(t, raw[:])
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonLightingOpt", n, 1)
		return metacanon.LightingOptInfo{}
	}

	dst := metacanon.LightingOptInfo{
		// ExifTool LightingOpt table uses FIRST_ENTRY=1.
		PeripheralIlluminationCorr: raw[0],
	}
	if n > 1 {
		dst.AutoLightingOptimizer = raw[1]
	}
	if n > 2 {
		dst.HighlightTonePriority = raw[2]
	}
	if n > 3 {
		dst.LongExposureNoiseReduction = raw[3]
	}
	if n > 4 {
		dst.HighISONoiseReduction = raw[4]
	}
	if n > 9 {
		dst.DigitalLensOptimizer = raw[9]
	}
	if n > 10 {
		dst.DualPixelRaw = raw[10]
	}
	return dst
}

// parseCanonLensInfo parses tag 0x4019 (LensInfoForService).
func (r *Reader) parseCanonLensInfo(t tag.Entry) metacanon.LensInfoForService {
	dst := metacanon.LensInfoForService{}
	block := r.parseCanonBlockPreview(t)
	if block.PreviewCount == 0 {
		return dst
	}
	n := min(int(block.PreviewCount), len(dst.Raw))
	copy(dst.Raw[:], block.Preview[:n])
	dst.RawCount = uint8(n)
	// ExifTool ignores value if the first four bytes are all zero.
	if n >= 4 && dst.Raw[0] == 0 && dst.Raw[1] == 0 && dst.Raw[2] == 0 && dst.Raw[3] == 0 {
		return dst
	}
	dst.LensSerialNumber = canonHexBytes(dst.Raw[:n])
	return dst
}

// parseCanonMultiExp parses tag 0x4021 (MultiExp).
func (r *Reader) parseCanonMultiExp(t tag.Entry) metacanon.MultiExpInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonMultiExp", n, 3)
		return metacanon.MultiExpInfo{}
	}
	return metacanon.MultiExpInfo{
		// ExifTool MultiExp table uses FIRST_ENTRY=1.
		MultiExposure:        int32(raw[0]),
		MultiExposureControl: int32(raw[1]),
		MultiExposureShots:   int32(raw[2]),
	}
}

// parseCanonHDRInfo parses tag 0x4025 (HDRInfo).
func (r *Reader) parseCanonHDRInfo(t tag.Entry) metacanon.HDRInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 2 {
		r.warnCanonShortRead(t, "parseCanonHDRInfo", n, 2)
		return metacanon.HDRInfo{}
	}
	return metacanon.HDRInfo{
		// ExifTool HDRInfo table uses FIRST_ENTRY=1.
		HDR:       int32(raw[0]),
		HDREffect: int32(raw[1]),
	}
}

// parseCanonCameraSettings parses tag 0x0001 (CanonCameraSettings).
func (r *Reader) parseCanonCameraSettings(t tag.Entry) metacanon.CameraSettings {
	var raw [64]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonCameraSettings", n, 2)
		return metacanon.CameraSettings{}
	}

	dst := metacanon.CameraSettings{
		MacroMode: metacanon.MacroMode(raw[0]),
		SelfTimer: int16(raw[1]),
	}
	if n > 2 {
		dst.Quality = metacanon.Quality(int16(raw[2]))
	}
	if n > 3 {
		dst.CanonFlashMode = metacanon.CanonFlashMode(int16(raw[3]))
	}
	if n > 4 {
		dst.ContinuousDrive = metacanon.ContinuousDrive(int16(raw[4]))
	}
	if n > 6 {
		dst.FocusMode = metacanon.FocusMode(int16(raw[6]))
	}
	if n > 8 {
		dst.RecordMode = metacanon.RecordMode(int16(raw[8]))
	}
	if n > 9 {
		dst.CanonImageSize = metacanon.CanonImageSize(int16(raw[9]))
	}
	if n > 10 {
		dst.EasyMode = metacanon.EasyMode(int16(raw[10]))
	}
	if n > 11 {
		dst.DigitalZoom = metacanon.DigitalZoom(int16(raw[11]))
	}
	if n > 12 {
		dst.Contrast = int16(raw[12])
	}
	if n > 13 {
		dst.Saturation = int16(raw[13])
	}
	if n > 14 {
		dst.Sharpness = int16(raw[14])
	}
	if n > 15 {
		dst.CameraISO = metacanon.CameraISO(int16(raw[15]))
	}
	if n > 16 {
		dst.MeteringMode = metacanon.MeteringMode(int16(raw[16]))
	}
	if n > 17 {
		dst.FocusRange = metacanon.FocusRange(int16(raw[17]))
	}
	if n > 18 {
		dst.AFPoint = raw[18]
	}
	if n > 19 {
		dst.CanonExposureMode = metacanon.ExposureMode(int16(raw[19]))
	}
	if n > 21 {
		dst.LensType = raw[21]
	}
	if n > 22 {
		dst.MaxFocalLength = raw[22]
	}
	if n > 23 {
		dst.MinFocalLength = raw[23]
	}
	if n > 24 {
		dst.FocalUnits = raw[24]
	}
	if n > 25 {
		dst.MaxAperture = int16(raw[25])
	}
	if n > 26 {
		dst.MinAperture = int16(raw[26])
	}
	if n > 27 {
		dst.FlashModel = metacanon.FlashModel(int16(raw[27]))
	}
	if n > 28 {
		dst.FlashBits = raw[28]
	}
	if n > 31 {
		dst.FocusContinuous = metacanon.FocusContinuous(int16(raw[31]))
	}
	if n > 32 {
		dst.AESetting = metacanon.AESetting(int16(raw[32]))
	}
	if n > 33 {
		dst.ImageStabilization = metacanon.ImageStabilization(int16(raw[33]))
	}
	if n > 34 {
		dst.DisplayAperture = raw[34]
	}
	if n > 35 {
		dst.ZoomSourceWidth = raw[35]
	}
	if n > 36 {
		dst.ZoomTargetWidth = raw[36]
	}
	if n > 38 {
		dst.SpotMeteringMode = metacanon.SpotMeteringMode(int16(raw[38]))
	}
	if n > 39 {
		dst.PhotoEffect = metacanon.PhotoEffect(int16(raw[39]))
	}
	if n > 40 {
		dst.ManualFlashOutput = metacanon.ManualFlashOutput(int16(raw[40]))
	}
	if n > 41 {
		dst.ColorTone = int16(raw[41])
	}
	if n > 45 {
		dst.SRAWQuality = metacanon.SRAWQuality(int16(raw[45]))
	}
	if n > 49 {
		dst.FocusBracketing = metacanon.FocusBracketing(int16(raw[49]))
	}
	if n > 50 {
		dst.Clarity = int16(raw[50])
	}
	if n > 51 {
		dst.HDRPQ = metacanon.HDRPQ(raw[51])
	}
	return dst
}

// parseCanonShotInfo parses tag 0x0004 (CanonShotInfo).
func (r *Reader) parseCanonShotInfo(t tag.Entry) metacanon.ShotInfo {
	var raw [64]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n == 0 {
		r.warnCanonShortRead(t, "parseCanonShotInfo", n, 1)
		return metacanon.ShotInfo{}
	}
	return metacanon.ShotInfo{
		AutoISO:                int16(raw[0]), // [1]
		BaseISO:                int16(raw[1]),
		MeasuredEV:             int16(raw[2]),
		TargetAperture:         int16(raw[3]),
		TargetExposureTime:     int16(raw[4]),
		ExposureCompensation:   int16(raw[5]),
		WhiteBalance:           int16(raw[6]),
		SlowShutter:            int16(raw[7]),
		SequenceNumber:         int16(raw[8]),
		OpticalZoomCode:        int16(raw[9]),
		CameraTemperature:      int16(raw[11]),
		FlashGuideNumber:       int16(raw[12]),
		AFPointsInFocus:        raw[13],
		FlashExposureComp:      int16(raw[14]),
		AutoExposureBracketing: int16(raw[15]),
		AEBBracketValue:        int16(raw[16]),
		ControlMode:            int16(raw[17]),
		FocusDistance:          metacanon.NewFocusDistance(raw[18], raw[19]),
		FNumber:                int16(raw[20]),
		ExposureTime:           int16(raw[21]),
		MeasuredEV2:            int16(raw[22]),
		BulbDuration:           int16(raw[23]),
		CameraType:             int16(raw[25]),
		AutoRotate:             int16(raw[26]),
		NDFilter:               int16(raw[27]),
		SelfTimer2:             int16(raw[28]),
		FlashOutput:            int16(raw[32]),
	}
}

// parseCanonFileInfo parses tag 0x0093 (CanonFileInfo).
func (r *Reader) parseCanonFileInfo(t tag.Entry) metacanon.FileInfo {
	var raw [64]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 60 {
		r.warnCanonShortRead(t, "parseCanonFileInfo", n, 60)
		return metacanon.FileInfo{}
	}

	// Tag 0x0093 index 1 is model-dependent (FileNumber or ShutterCount).
	// Preserve raw 32-bit representation for both fields.
	return metacanon.FileInfo{
		FileNumber:                  uint32(raw[0]) | (uint32(raw[1]) << 16),
		BracketMode:                 metacanon.BracketMode(int16(raw[2])),
		BracketValue:                int16(raw[3]),
		BracketShotNumber:           int16(raw[4]),
		RawJpgQuality:               metacanon.RawJpgQuality(raw[5]),
		RawJpgSize:                  metacanon.RawJpgSize(raw[6]),
		LongExposureNoiseReduction2: metacanon.OnOffAuto(raw[7]),
		WBBracketMode:               int16(raw[8]),
		WBBracketValueAB:            int16(raw[11]),
		WBBracketValueGM:            int16(raw[12]),
		FilterEffect:                metacanon.FilterEffect(raw[13]),
		ToningEffect:                metacanon.ToningEffect(raw[14]),
		MacroMagnification:          int16(raw[15]),
		LiveViewShooting:            metacanon.OnOffAuto(raw[18]),
		FocusDistance:               metacanon.NewFocusDistance(raw[19], raw[20]),
		ShutterMode:                 metacanon.ShutterMode(raw[22]),
		FlashExposureLock:           metacanon.OnOffAuto(raw[24]),
		AntiFlicker:                 metacanon.OnOffAuto(raw[31]),
		RFLensType:                  metacanon.CanonRFLensType(raw[60]),
	}
}

// parseCanonTimeInfo parses tag 0x0035 (TimeInfo).
func (r *Reader) parseCanonTimeInfo(t tag.Entry) metacanon.CanonTimeInfo {
	var raw [4]uint32
	if n := r.parseUint32List(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonTimeInfo", n, 3)
		return metacanon.CanonTimeInfo{}
	}
	return metacanon.CanonTimeInfo{
		TimeZone:        int32(raw[0]),
		TimeZoneCity:    metacanon.TimeZoneCity(int32(raw[1])),
		DaylightSavings: metacanon.DaylightSavings(int32(raw[2])),
	}
}

// parseCanonAFInfo parses tag 0x0012 (AFInfo).
func (r *Reader) parseCanonAFInfo(t tag.Entry) metacanon.AFInfo {
	var words [2048]uint16
	n := r.parseCanonUint16List(t, words[:])
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonAFInfo", n, 1)
		return metacanon.AFInfo{}
	}
	var dst metacanon.AFInfo
	fillCanonAFInfo(&dst, words[:n], r.canonModelName(), int(t.UnitCount))
	return dst
}

func fillCanonAFInfo(dst *metacanon.AFInfo, words []uint16, model string, afInfoCount int) {
	n := len(words)
	*dst = metacanon.AFInfo{}

	dst.AFAreaMode = 0
	dst.AFAreaWidths = nil
	dst.AFAreaHeights = nil
	dst.AFPointsSelectedBits = nil

	dst.NumAFPoints = canonU16At(words, n, 0)
	dst.ValidAFPoints = canonU16At(words, n, 1)
	dst.CanonImageWidth = canonU16At(words, n, 2)
	dst.CanonImageHeight = canonU16At(words, n, 3)
	dst.AFImageWidth = canonU16At(words, n, 4)
	dst.AFImageHeight = canonU16At(words, n, 5)
	dst.AFAreaWidth = canonU16At(words, n, 6)
	dst.AFAreaHeight = canonU16At(words, n, 7)

	isEOS := canonModelIsEOS(model)
	num := int(dst.NumAFPoints)
	if num <= 0 {
		return
	}

	xStart := 8
	yStart := xStart + num
	xVals := canonSignedRangeFromUint16(words, n, xStart, num)
	yVals := canonSignedRangeFromUint16(words, n, yStart, num)
	dst.AFAreaXPositions = xVals
	dst.AFAreaYPositions = yVals

	bitWords := canonBitWordCount(num)
	inFocusStart := yStart + num
	dst.AFPointsInFocusBits = canonDecodeBitWordsRange(words, n, inFocusStart, bitWords)

	if !isEOS {
		seq11 := inFocusStart + bitWords
		// ExifTool skips seq-11 PrimaryAFPoint when AFInfoCount==36.
		if afInfoCount != 36 {
			dst.PrimaryAFPoint = canonU16At(words, n, seq11)
		}
		// ExifTool also defines seq-12 PrimaryAFPoint after an 8-word unknown block.
		seq12 := seq11 + 8
		if afInfoCount != 36 {
			seq12++
		}
		if primary := canonU16At(words, n, seq12); primary != 0 {
			dst.PrimaryAFPoint = primary
		}
	}

	pointCount := num
	if len(xVals) < pointCount {
		pointCount = len(xVals)
	}
	if len(yVals) < pointCount {
		pointCount = len(yVals)
	}
	if pointCount == 0 {
		return
	}

	pts := make([]metacanon.AFPoint, pointCount)
	w := int16(dst.AFAreaWidth)
	h := int16(dst.AFAreaHeight)
	for i := 0; i < pointCount; i++ {
		pts[i] = metacanon.NewAFPoint(w, h, xVals[i], yVals[i])
	}
	dst.AFPoints = pts
}

// parseCanonAFInfo2 parses tags 0x0026 and 0x003c (AFInfo2/AFInfo3).
func (r *Reader) parseCanonAFInfo2(t tag.Entry) metacanon.AFInfo {
	var words [2048]uint16
	n := r.parseCanonUint16List(t, words[:])
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonAFInfo2", n, 1)
		return metacanon.AFInfo{}
	}
	model := r.canonModelName()
	isAFInfo3 := metacanon.MakerNoteTag(t.ID) == metacanon.AFInfo3
	dst := metacanon.AFInfo{
		AFAreaWidth:      0,
		AFAreaHeight:     0,
		AFAreaMode:       metacanon.AFAreaMode(uint16(words[1])),
		NumAFPoints:      uint16(words[2]),
		ValidAFPoints:    uint16(words[3]),
		CanonImageWidth:  uint16(words[4]),
		CanonImageHeight: uint16(words[5]),
		AFImageWidth:     uint16(words[6]),
		AFImageHeight:    uint16(words[7]),
	}

	isEOS := canonModelIsEOS(model)
	num := int(dst.NumAFPoints)
	if num <= 0 {
		return metacanon.AFInfo{}
	}

	widthStart := 8
	heightStart := widthStart + num
	xStart := heightStart + num
	yStart := xStart + num
	bitsStart := yStart + num
	maskWordCount := canonBitWordCount(num)
	selectedStart := bitsStart + maskWordCount

	widthLen := canonRangeLen(n, widthStart, num)
	heightLen := canonRangeLen(n, heightStart, num)
	xLen := canonRangeLen(n, xStart, num)
	yLen := canonRangeLen(n, yStart, num)

	decodeCoords := r.afInfoDecodeOptions.has(AFInfoDecodeCoords)
	decodePoints := r.afInfoDecodeOptions.has(AFInfoDecodePoints)
	decodeInFocus := r.afInfoDecodeOptions.has(AFInfoDecodeInFocus)
	decodeSelected := r.afInfoDecodeOptions.has(AFInfoDecodeSelected)

	if decodeCoords {
		totalSigned := widthLen + heightLen + xLen + yLen
		if totalSigned == 0 {
			dst.AFAreaWidths = nil
			dst.AFAreaHeights = nil
			dst.AFAreaXPositions = nil
			dst.AFAreaYPositions = nil
		} else {
			coords := make([]int16, totalSigned)
			offset := 0

			dst.AFAreaWidths = coords[offset : offset+widthLen]
			for i := range widthLen {
				dst.AFAreaWidths[i] = int16(words[widthStart+i])
			}
			offset += widthLen

			dst.AFAreaHeights = coords[offset : offset+heightLen]
			for i := range heightLen {
				dst.AFAreaHeights[i] = int16(words[heightStart+i])
			}
			offset += heightLen

			dst.AFAreaXPositions = coords[offset : offset+xLen]
			for i := range xLen {
				dst.AFAreaXPositions[i] = int16(words[xStart+i])
			}
			offset += xLen

			dst.AFAreaYPositions = coords[offset : offset+yLen]
			for i := range yLen {
				dst.AFAreaYPositions[i] = int16(words[yStart+i])
			}
		}
	}

	if decodeInFocus {
		dst.AFPointsInFocusBits = canonDecodeBitWordsRange(words[:], n, bitsStart, maskWordCount)
	} else {
		dst.AFPointsInFocusBits = nil
	}
	dst.AFPointsSelectedBits = nil
	dst.PrimaryAFPoint = 0

	if isEOS && decodeSelected {
		// ExifTool only decodes AFPointsSelected for EOS models.
		dst.AFPointsSelectedBits = canonDecodeBitWordsRange(words[:], n, selectedStart, maskWordCount)
	} else if !isAFInfo3 {
		// Non-EOS AFInfo2 uses an unknown field of maskWordCount+1 at seq 13.
		dst.PrimaryAFPoint = canonU16At(words[:], n, selectedStart+maskWordCount+1)
	}

	if !decodePoints {
		dst.AFPoints = nil
		return dst
	}

	pointCount := min(yLen, min(xLen, min(heightLen, min(widthLen, num))))
	if pointCount <= 0 {
		dst.AFPoints = nil
		return dst
	}

	pts := make([]metacanon.AFPoint, pointCount)
	xAdjust := int16(dst.CanonImageWidth / 2)
	yAdjust := int16(dst.CanonImageHeight / 2)
	for i := 0; i < pointCount; i++ {
		w := int16(words[widthStart+i])
		h := int16(words[heightStart+i])
		x := int16(words[xStart+i]) + xAdjust - (w / 2)
		y := int16(words[yStart+i]) + yAdjust - (h / 2)
		pts[i] = metacanon.NewAFPoint(w, h, x, y)
	}
	dst.AFPoints = pts
	return dst
}

func canonU16At(vals []uint16, n, idx int) uint16 {
	if idx < 0 || idx >= n {
		return 0
	}
	return vals[idx]
}

func canonBitWordCount(pointCount int) int {
	if pointCount <= 0 {
		return 0
	}
	return (pointCount + 15) / 16
}

func canonRangeLen(n, start, count int) int {
	if count <= 0 || start < 0 || start >= n {
		return 0
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return 0
	}
	return end - start
}

func canonSignedRangeFromUint16(vals []uint16, n, start, count int) []int16 {
	if count <= 0 || start < 0 || start >= n {
		return nil
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return nil
	}
	out := make([]int16, end-start)
	for i := 0; i < len(out); i++ {
		out[i] = int16(vals[start+i])
	}
	return out
}

func canonDecodeBitWordsRange(vals []uint16, n, start, count int) []int {
	if count <= 0 || start < 0 || start >= n {
		return nil
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return nil
	}

	capHint := 0
	for i := start; i < end; i++ {
		capHint += bits.OnesCount16(vals[i])
	}
	out := make([]int, 0, capHint)
	base := 0
	for i := start; i < end; i++ {
		word := vals[i]
		for bit := 0; bit < 16; bit++ {
			if word&(1<<bit) != 0 {
				out = append(out, base+bit)
			}
		}
		base += 16
	}
	return out
}

func (r *Reader) canonModelName() string {
	if model := r.Exif.IFD0.Model; model != "" {
		return model
	}
	return r.makerNoteInfo().Canon.ImageType
}

func canonModelIsEOS(model string) bool {
	return strings.Contains(model, "EOS")
}

func canonHexBytes(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	const table = "0123456789abcdef"
	var out strings.Builder
	out.Grow(len(b) * 2)
	for i := 0; i < len(b); i++ {
		v := b[i]
		out.WriteByte(table[v>>4])
		out.WriteByte(table[v&0x0f])
	}
	return out.String()
}
