package exif

import (
	"bytes"
	"math"
	"math/bits"
	"strings"

	"github.com/evanoberholster/imagemeta/meta"
	metacanon "github.com/evanoberholster/imagemeta/meta/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

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
		dst.LensModel = canonTerminateAtNUL(r.parseStringAllowUndefined(t))
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
		dst.BatteryType = r.parseCanonBatteryType(t)
	case metacanon.CanonAFInfo:
		candidate := r.parseCanonAFInfo(t)
		if canonShouldReplaceAFInfo(dst.AFInfo, candidate) {
			dst.AFInfo = candidate
		}
	case metacanon.CanonAFInfo2, metacanon.AFInfo3:
		candidate := r.parseCanonAFInfo2(t)
		if canonShouldReplaceAFInfo(dst.AFInfo, candidate) {
			dst.AFInfo = candidate
		}
	case metacanon.FaceDetect1:
		dst.FaceDetect1 = r.parseCanonFaceDetect1(t)
	case metacanon.FaceDetect2:
		dst.FaceDetect2 = r.parseCanonFaceDetect2(t)
	case metacanon.FaceDetect3:
		dst.FaceDetect3 = r.parseCanonFaceDetect3(t)
	case metacanon.ImageUniqueID:
		dst.ImageUniqueID = r.parseCanonImageUniqueID(t)
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
	switch t.Type {
	case tag.TypeShort, tag.TypeSignedShort:
		return r.parseCanonRawUint16List(t, dst, int(t.UnitCount))
	case tag.TypeUndefined:
		return r.parseCanonRawUint16List(t, dst, int(t.UnitCount/2))
	default:
		return 0
	}
}

// parseCanonRawUint16List reads uint16 values in chunks to support large
// maker-note payloads that exceed readTagBytes/state.buf capacity.
func (r *Reader) parseCanonRawUint16List(t tag.Entry, dst []uint16, wordCount int) int {
	if len(dst) == 0 || wordCount <= 0 || t.UnitCount == 0 {
		return 0
	}
	if wordCount > len(dst) {
		wordCount = len(dst)
	}

	if t.IsEmbedded() {
		switch t.Type {
		case tag.TypeShort, tag.TypeSignedShort:
			return t.EmbeddedShorts(dst[:wordCount])
		}
		// UNDEFINED embedded payload is up to 4 bytes.
		t.EmbeddedValue(r.state.buf[:4])
		n := min(wordCount, 2)
		for i := range n {
			start := i * 2
			dst[i] = t.ByteOrder.Uint16(r.state.buf[start : start+2])
		}
		return n
	}

	if err := r.seekToTag(t); err != nil {
		return 0
	}

	remainingBytes := wordCount * 2
	readWords := 0

	for remainingBytes > 0 {
		chunkBytes := min(remainingBytes, len(r.state.buf))
		if chunkBytes&1 != 0 {
			chunkBytes--
		}
		if chunkBytes <= 0 {
			break
		}

		buf, err := r.fastRead(chunkBytes)
		if err != nil {
			break
		}
		gotWords := len(buf) / 2
		if gotWords == 0 {
			break
		}
		for i := range gotWords {
			start := i * 2
			dst[readWords+i] = t.ByteOrder.Uint16(buf[start : start+2])
		}
		readWords += gotWords
		remainingBytes -= gotWords * 2
	}

	remainingTagBytes := int(t.Size()) - (readWords * 2)
	if remainingTagBytes > 0 {
		if err := r.discard(remainingTagBytes); err != nil {
			return readWords
		}
	}

	return readWords
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

func (r *Reader) parseCanonPreviewImageInfo(t tag.Entry) metacanon.PreviewImageInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 5 {
		r.warnCanonShortRead(t, "parseCanonPreviewImageInfo", n, 5)
		return metacanon.PreviewImageInfo{}
	}

	return metacanon.PreviewImageInfo{
		PreviewQuality:     metacanon.Quality(int16(raw[0])),
		PreviewImageLength: uint32(raw[1]),
		PreviewImageWidth:  uint32(raw[2]),
		PreviewImageHeight: uint32(raw[3]),
		PreviewImageStart:  uint32(raw[4]),
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
	if n := r.parseUint32List(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonRawBurstInfo", n, 3)
		return metacanon.RawBurstInfo{}
	}
	return metacanon.RawBurstInfo{
		RawBurstImageNum:   raw[1],
		RawBurstImageCount: raw[2],
	}
}

// parseCanonImageUniqueID parses Canon maker-note tag 0x0028 into meta.UUID.
//
// ExifTool renders this value as hex text, but imagemeta stores it as a UUID.
func (r *Reader) parseCanonImageUniqueID(t tag.Entry) meta.UUID {
	buf := r.parseOpaqueBytes(t, canonUUIDBytesLength)
	if len(buf) != 16 {
		return meta.NilUUID
	}
	uuid, err := meta.UUIDFromBytes(buf)
	if err != nil {
		return meta.NilUUID
	}
	return uuid
}

func (r *Reader) parseCanonFaceDetect1(t tag.Entry) metacanon.FaceDetect1Info {
	var raw [26]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 5 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect1", n, 5)
		return metacanon.FaceDetect1Info{}
	}
	dst := metacanon.FaceDetect1Info{
		FacesDetected: raw[2],
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
	if n := r.parseByteList(t, raw[:]); n < 3 {
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
	if n := r.parseCanonUint16List(t, raw[:]); n < 4 {
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

const canonLensInfoByteLength = 5

// parseCanonLensInfo parses tag 0x4019 (LensInfoForService).
func (r *Reader) parseCanonLensInfo(t tag.Entry) metacanon.LensInfoForService {
	dst := metacanon.LensInfoForService{}
	raw := r.parseOpaqueBytes(t, canonLensInfoByteLength)
	l := int(len(raw))
	if l != 5 {
		r.warnCanonShortRead(t, "parseCanonLensInfo", l, int(t.Size()))
		return metacanon.LensInfoForService{}
	}
	n := min(l, canonLensInfoByteLength)
	copy(dst.Raw[:], raw[:n])
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
	var raw [53]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 1 {
		r.warnCanonShortRead(t, "parseCanonCameraSettings", n, 1)
		return metacanon.CameraSettings{}
	}
	declaredSizeBytes := uint32(raw[0])
	if declaredSizeBytes != t.Size() {
		if r.warnEnabled() {
			r.warn().
				Str("parser", "parseCanonCameraSettings").
				Uint16("tagID", uint16(t.ID)).
				Str("tagName", t.Name()).
				Stringer("tagType", t.Type).
				Uint32("unitCount", t.UnitCount).
				Uint32("declaredSizeBytes", declaredSizeBytes).
				Uint32("actualSizeBytes", t.Size()).
				Msg("invalid canon camera settings payload length")
		}
		return metacanon.CameraSettings{}
	}
	if n < 2 {
		return metacanon.CameraSettings{}
	}

	// Canon CameraSettings stores a 16-bit size word first. The remaining
	// words map directly to ExifTool's documented sequence numbers, so payload
	// sequence N lives at settings[N-1].
	settings := raw[1:n]

	var dst metacanon.CameraSettings
	dst.MacroMode = metacanon.MacroMode(settings[0]) // [1]

	if len(settings) > 1 {
		dst.SelfTimer = int16(settings[1]) // [2]
	}
	if len(settings) > 2 {
		dst.Quality = metacanon.Quality(int16(settings[2])) // [3]
	}
	if len(settings) > 3 {
		dst.CanonFlashMode = metacanon.CanonFlashMode(int16(settings[3])) // [4]
	}
	if len(settings) > 4 {
		dst.ContinuousDrive = metacanon.ContinuousDrive(int16(settings[4])) // [5]
	}
	if len(settings) > 6 {
		dst.FocusMode = metacanon.FocusMode(int16(settings[6])) // [7]
	}
	if len(settings) > 8 {
		dst.RecordMode = metacanon.RecordMode(int16(settings[8])) // [9]
	}
	if len(settings) > 9 {
		dst.CanonImageSize = metacanon.CanonImageSize(int16(settings[9])) // [10]
	}
	if len(settings) > 10 {
		dst.EasyMode = metacanon.EasyMode(int16(settings[10])) // [11]
	}
	if len(settings) > 11 {
		dst.DigitalZoom = metacanon.DigitalZoom(int16(settings[11])) // [12]
	}
	if len(settings) > 12 {
		dst.Contrast = int16(settings[12]) // [13]
	}
	if len(settings) > 13 {
		dst.Saturation = int16(settings[13]) // [14]
	}
	if len(settings) > 14 {
		dst.Sharpness = int16(settings[14]) // [15]
	}
	if len(settings) > 15 {
		dst.CameraISO = metacanon.CameraISO(int16(settings[15])) // [16]
	}
	if len(settings) > 16 {
		dst.MeteringMode = metacanon.MeteringMode(int16(settings[16])) // [17]
	}
	if len(settings) > 17 {
		dst.FocusRange = metacanon.FocusRange(int16(settings[17])) // [18]
	}
	if len(settings) > 18 {
		dst.AFPoint = settings[18] // [19]
	}
	if len(settings) > 19 {
		dst.CanonExposureMode = metacanon.ExposureMode(int16(settings[19])) // [20]
	}
	if len(settings) > 21 {
		dst.LensType = metacanon.CanonLensType(settings[21]) // [22]
	}
	if len(settings) > 22 {
		dst.MaxFocalLength = settings[22] // [23]
	}
	if len(settings) > 23 {
		dst.MinFocalLength = settings[23] // [24]
	}
	if len(settings) > 24 {
		dst.FocalUnits = settings[24] // [25]
	}
	if len(settings) > 25 {
		dst.MaxAperture = parseCanonMaxAperture(settings[25]) // [26]
	}
	if len(settings) > 26 {
		dst.MinAperture = parseCanonMaxAperture(settings[26]) // [27]
	}
	if len(settings) > 27 {
		dst.FlashModel = metacanon.FlashModel(int16(settings[27])) // [28]
	}
	if len(settings) > 28 {
		dst.FlashBits = settings[28] // [29]
	}
	if len(settings) > 31 {
		dst.FocusContinuous = metacanon.FocusContinuous(int16(settings[31])) // [32]
	}
	if len(settings) > 32 {
		dst.AESetting = metacanon.AESetting(int16(settings[32])) // [33]
	}
	if len(settings) > 33 {
		dst.ImageStabilization = metacanon.ImageStabilization(int16(settings[33])) // [34]
	}
	if len(settings) > 34 {
		dst.DisplayAperture = parseCanonDisplayAperture(settings[34]) // [35]
	}
	if len(settings) > 35 {
		dst.ZoomSourceWidth = settings[35] // [36]
	}
	if len(settings) > 36 {
		dst.ZoomTargetWidth = settings[36] // [37]
	}
	if len(settings) > 38 {
		dst.SpotMeteringMode = metacanon.SpotMeteringMode(int16(settings[38])) // [39]
	}
	if len(settings) > 39 {
		dst.PhotoEffect = metacanon.PhotoEffect(int16(settings[39])) // [40]
	}
	if len(settings) > 40 {
		dst.ManualFlashOutput = metacanon.ManualFlashOutput(int16(settings[40])) // [41]
	}
	if len(settings) > 41 {
		dst.ColorTone = int16(settings[41]) // [42]
	}
	if len(settings) > 45 {
		dst.SRAWQuality = metacanon.SRAWQuality(int16(settings[45])) // [46]
	}
	if len(settings) > 49 {
		dst.FocusBracketing = metacanon.FocusBracketing(int16(settings[49])) // [50]
	}
	if len(settings) > 50 {
		dst.Clarity = int16(settings[50]) // [51]
	}
	if len(settings) > 51 {
		dst.HDRPQ = metacanon.HDRPQ(settings[51]) // [52]
	}
	return dst
}

// parseCanonShotInfo parses tag 0x0004 (CanonShotInfo).
func (r *Reader) parseCanonShotInfo(t tag.Entry) metacanon.ShotInfo {
	var raw [64]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonShotInfo", n, 1)
		return metacanon.ShotInfo{}
	}
	declaredSizeBytes := uint32(raw[0])
	if declaredSizeBytes != t.Size() {
		if r.warnEnabled() {
			r.warn().
				Str("parser", "parseCanonShotInfo").
				Uint16("tagID", uint16(t.ID)).
				Str("tagName", t.Name()).
				Stringer("tagType", t.Type).
				Uint32("unitCount", t.UnitCount).
				Uint32("declaredSizeBytes", declaredSizeBytes).
				Uint32("actualSizeBytes", t.Size()).
				Msg("invalid canon shot info payload length")
		}
		return metacanon.ShotInfo{}
	}
	if n < 2 {
		return metacanon.ShotInfo{}
	}
	// Canon ShotInfo stores a 16-bit size word first. The remaining words map
	// directly to ExifTool's documented sequence numbers, so sequence N lives
	// at settings[N-1].
	settings := raw[1:n]
	var dst metacanon.ShotInfo

	dst.AutoISO = int16(settings[0]) // [1]
	dst.AutoISOValue = canonShotISO(dst.AutoISO)
	dst.BaseISO = int16(settings[1]) // [2]
	dst.BaseISOValue = canonShotISO(dst.BaseISO)
	dst.ActualISO = canonShotActualISO(dst.AutoISOValue, dst.BaseISOValue)

	if len(settings) > 2 {
		dst.MeasuredEV = int16(settings[2]) // [3]
	}
	if len(settings) > 3 {
		dst.TargetAperture = int16(settings[3]) // [4]
		dst.TargetApertureValue = canonShotAperture(dst.TargetAperture)
	}
	if len(settings) > 4 {
		dst.TargetExposureTime = int16(settings[4]) // [5]
		dst.TargetExposureTimeValue = canonShotExposureTime(dst.TargetExposureTime, false)
	}
	if len(settings) > 5 {
		dst.ExposureCompensation = int16(settings[5]) // [6]
	}
	if len(settings) > 6 {
		dst.WhiteBalance = metacanon.WhiteBalance(int16(settings[6])) // [7]
	}
	if len(settings) > 7 {
		dst.SlowShutter = metacanon.SlowShutter(int16(settings[7])) // [8]
	}
	if len(settings) > 8 {
		dst.SequenceNumber = int16(settings[8]) // [9]
	}
	if len(settings) > 9 {
		dst.OpticalZoomCode = int16(settings[9]) // [10]
	}
	if len(settings) > 11 {
		dst.CameraTemperature = int16(settings[11]) // [12]
		dst.CameraTemperatureC = canonShotCameraTemperature(dst.CameraTemperature, r.canonModelName())
	}
	if len(settings) > 12 {
		dst.FlashGuideNumber = int16(settings[12]) // [13]
		dst.FlashGuideNumberMeters = canonShotFlashGuideNumber(dst.FlashGuideNumber)
	}
	if len(settings) > 13 {
		dst.AFPointsInFocus = settings[13] // [14]
	}
	if len(settings) > 14 {
		dst.FlashExposureComp = int16(settings[14]) // [15]
	}
	if len(settings) > 15 {
		dst.AutoExposureBracketing = int16(settings[15]) // [16]
	}
	if len(settings) > 16 {
		dst.AEBBracketValue = int16(settings[16]) // [17]
	}
	if len(settings) > 17 {
		dst.ControlMode = int16(settings[17]) // [18]
	}
	if len(settings) > 19 && settings[18] != 0 {
		dst.FocusDistance = metacanon.NewFocusDistance(settings[18], settings[19]) // [19-20]
	}
	if len(settings) > 20 {
		dst.FNumber = int16(settings[20]) // [21]
		dst.FNumberValue = canonShotAperture(dst.FNumber)
	}
	if len(settings) > 21 {
		dst.ExposureTime = int16(settings[21]) // [22]
		dst.ExposureTimeValue = canonShotExposureTime(dst.ExposureTime, r.canonShotInfoLegacyExposureTime())
	}
	if len(settings) > 22 {
		dst.MeasuredEV2 = int16(settings[22]) // [23]
	}
	if len(settings) > 23 {
		dst.BulbDuration = int16(settings[23]) // [24]
	}
	if len(settings) > 25 {
		dst.CameraType = metacanon.CameraType(int16(settings[25])) // [26]
	}
	if len(settings) > 26 {
		dst.AutoRotate = metacanon.AutoRotate(int16(settings[26])) // [27]
	}
	if len(settings) > 27 {
		dst.NDFilter = metacanon.NDFilter(int16(settings[27])) // [28]
	}
	if len(settings) > 28 {
		dst.SelfTimer2 = int16(settings[28]) // [29]
	}
	if len(settings) > 32 {
		dst.FlashOutput = int16(settings[32]) // [33]
	}
	return dst
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

const canonBatteryTypePayloadSize = 76

const (
	canonUUIDBytesLength = 16
)

// parseCanonBatteryType parses Canon Camera:BatteryType (tag 0x0038) like ExifTool.
//
// ExifTool behavior:
//   - only valid when count == 76
//   - ignore first 4 bytes
//   - return bytes up to first NUL; empty => not present
func (r *Reader) parseCanonBatteryType(t tag.Entry) string {
	if t.Size() != canonBatteryTypePayloadSize {
		if r.warnEnabled() {
			r.warn().
				Str("parser", "parseCanonBatteryType").
				Uint16("tagID", uint16(t.ID)).
				Str("tagName", t.Name()).
				Stringer("tagType", t.Type).
				Uint32("unitCount", t.UnitCount).
				Uint32("sizeBytes", t.Size()).
				Msg("invalid canon battery type payload length")
		}
		return ""
	}
	raw, _, err := r.readTagBytes(t, canonBatteryTypePayloadSize)
	if err != nil || len(raw) < canonBatteryTypePayloadSize {
		r.warnCanonShortRead(t, "parseCanonBatteryType", len(raw), canonBatteryTypePayloadSize)
		return ""
	}
	payload := raw[4:] // skip 4-byte header
	i := bytes.IndexByte(payload, 0)
	if i < 0 {
		i = len(payload)
	}
	if i == 0 {
		return ""
	}
	return string(payload[:i])
}

// parseCanonAFInfo parses tag 0x0012 (AFInfo).
func (r *Reader) parseCanonAFInfo(t tag.Entry) metacanon.AFInfo {
	var wordsStack [2048]uint16
	words, truncated := canonAFWordsBuffer(wordsStack[:], t.UnitCount)
	if truncated {
		r.warnCanonTruncatedWords(t, "parseCanonAFInfo", len(words), int(t.UnitCount))
	}
	n := r.parseCanonUint16List(t, words)
	source := canonAFInfoSource(tag.ID(metacanon.CanonAFInfo))
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonAFInfo", n, 1)
		return metacanon.AFInfo{Source: source}
	}
	var dst metacanon.AFInfo
	fillCanonAFInfo(&dst, words[:n], r.canonModelName(), int(t.UnitCount))
	return dst
}

func fillCanonAFInfo(dst *metacanon.AFInfo, words []uint16, model string, afInfoCount int) {
	n := len(words)
	*dst = metacanon.AFInfo{
		Source:           metacanon.AFInfoSourceAFInfo,
		NumAFPoints:      canonU16At(words, n, 0),
		ValidAFPoints:    canonU16At(words, n, 1),
		CanonImageWidth:  canonU16At(words, n, 2),
		CanonImageHeight: canonU16At(words, n, 3),
		AFImageWidth:     canonU16At(words, n, 4),
		AFImageHeight:    canonU16At(words, n, 5),
		AFAreaWidth:      canonU16At(words, n, 6),
		AFAreaHeight:     canonU16At(words, n, 7),
	}

	isEOS := canonModelIsEOS(model)
	num := int(dst.NumAFPoints)
	if num <= 0 {
		return
	}

	xStart := 8
	yStart := xStart + num

	bitWords := canonBitWordCount(num)
	inFocusStart := yStart + num
	dst.AFPointsInFocusBits = canonDecodeBitWordsRange(words, n, inFocusStart, bitWords)

	if !isEOS {
		dst.PrimaryAFPoint = canonLegacyAFInfoPrimary(words, n, inFocusStart+bitWords, afInfoCount)
	}

	areas := canonDecodeUniformAFArea(
		words,
		n,
		xStart,
		yStart,
		num,
		int16(dst.AFAreaWidth),
		int16(dst.AFAreaHeight),
	)
	dst.AFArea = areas
	// AFInfo (0x0012) stores width/height/x/y directly in the AF area tuples.
	dst.AFPoints = areas
}

// parseCanonAFInfo2 parses tags 0x0026 and 0x003c (AFInfo2/AFInfo3).
func (r *Reader) parseCanonAFInfo2(t tag.Entry) metacanon.AFInfo {
	var wordsStack [2048]uint16
	words, truncated := canonAFWordsBuffer(wordsStack[:], t.UnitCount)
	if truncated {
		r.warnCanonTruncatedWords(t, "parseCanonAFInfo2", len(words), int(t.UnitCount))
	}
	n := r.parseCanonUint16List(t, words)
	source := canonAFInfoSource(t.ID)
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonAFInfo2", n, 1)
		return metacanon.AFInfo{Source: source}
	}
	model := r.canonModelName()
	isAFInfo3 := metacanon.MakerNoteTag(t.ID) == metacanon.AFInfo3
	dst := metacanon.AFInfo{
		Source:           source,
		AFAreaWidth:      0,
		AFAreaHeight:     0,
		AFAreaMode:       metacanon.AFAreaMode(canonU16At(words, n, 1)),
		NumAFPoints:      canonU16At(words, n, 2),
		ValidAFPoints:    canonU16At(words, n, 3),
		CanonImageWidth:  canonU16At(words, n, 4),
		CanonImageHeight: canonU16At(words, n, 5),
		AFImageWidth:     canonU16At(words, n, 6),
		AFImageHeight:    canonU16At(words, n, 7),
	}

	isEOS := canonModelIsEOS(model)
	num := int(dst.NumAFPoints)
	if num <= 0 {
		return dst
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
	areaCount := min(yLen, min(xLen, min(heightLen, widthLen)))
	var pts []metacanon.AFPoint

	if decodeCoords {
		if areaCount == 0 {
			dst.AFArea = nil
		} else {
			if decodePoints {
				combined := make([]metacanon.AFPoint, areaCount*2)
				dst.AFArea = combined[:areaCount]
				pts = combined[areaCount:]
			} else {
				dst.AFArea = make([]metacanon.AFPoint, areaCount)
			}
			for i := 0; i < len(dst.AFArea); i++ {
				dst.AFArea[i] = metacanon.NewAFPoint(
					int16(words[widthStart+i]),
					int16(words[heightStart+i]),
					int16(words[xStart+i]),
					int16(words[yStart+i]),
				)
			}
		}
	} else {
		dst.AFArea = nil
	}

	wantSelected := isEOS && decodeSelected
	if decodeInFocus || wantSelected {
		totalBits := 0
		if decodeInFocus {
			totalBits += canonCountBitWordsRange(words, n, bitsStart, maskWordCount)
		}
		if wantSelected {
			totalBits += canonCountBitWordsRange(words, n, selectedStart, maskWordCount)
		}
		combinedBits := make([]int, 0, totalBits)

		if decodeInFocus {
			startIdx := len(combinedBits)
			combinedBits = canonAppendBitWordsRange(combinedBits, words, n, bitsStart, maskWordCount)
			dst.AFPointsInFocusBits = combinedBits[startIdx:len(combinedBits)]
		} else {
			dst.AFPointsInFocusBits = nil
		}

		if wantSelected {
			// ExifTool only decodes AFPointsSelected for EOS models.
			startIdx := len(combinedBits)
			combinedBits = canonAppendBitWordsRange(combinedBits, words, n, selectedStart, maskWordCount)
			dst.AFPointsSelectedBits = combinedBits[startIdx:len(combinedBits)]
		} else {
			dst.AFPointsSelectedBits = nil
		}
	} else {
		dst.AFPointsInFocusBits = nil
		dst.AFPointsSelectedBits = nil
	}
	dst.PrimaryAFPoint = 0
	if !(isEOS && decodeSelected) && !isAFInfo3 {
		// Non-EOS AFInfo2 uses an unknown field of maskWordCount+1 at seq 13.
		dst.PrimaryAFPoint = canonU16At(words, n, selectedStart+maskWordCount+1)
	}

	if !decodePoints {
		dst.AFPoints = nil
		return dst
	}

	if areaCount <= 0 {
		dst.AFPoints = nil
		return dst
	}

	if pts == nil {
		pts = make([]metacanon.AFPoint, areaCount)
	}
	xAdjust := int16(dst.CanonImageWidth / 2)
	yAdjust := int16(dst.CanonImageHeight / 2)
	for i := 0; i < areaCount; i++ {
		var w, h, x, y int16
		if decodeCoords {
			area := dst.AFArea[i]
			w, h, x, y = area[0], area[1], area[2], area[3]
		} else {
			w = int16(words[widthStart+i])
			h = int16(words[heightStart+i])
			x = int16(words[xStart+i])
			y = int16(words[yStart+i])
		}
		x += xAdjust - (w / 2)
		y += yAdjust - (h / 2)
		pts[i] = metacanon.NewAFPoint(w, h, x, y)
	}
	dst.AFPoints = pts
	return dst
}

const canonAFWordsMax = 8192

func canonAFWordsBuffer(stack []uint16, unitCount uint32) ([]uint16, bool) {
	if unitCount == 0 {
		return stack[:0], false
	}
	wordCount := int(unitCount)
	truncated := false
	if unitCount > canonAFWordsMax {
		wordCount = canonAFWordsMax
		truncated = true
	}
	if wordCount <= len(stack) {
		return stack[:wordCount], truncated
	}
	return make([]uint16, wordCount), truncated
}

func (r *Reader) warnCanonTruncatedWords(t tag.Entry, parser string, got, want int) {
	if !r.warnEnabled() {
		return
	}
	r.warn().
		Str("parser", parser).
		Uint16("tagID", uint16(t.ID)).
		Str("tagName", t.Name()).
		Stringer("tagType", t.Type).
		Int("wordsDecoded", got).
		Int("wordsRequested", want).
		Msg("canon AF payload truncated to parser word cap")
}

func canonAFInfoSource(id tag.ID) metacanon.AFInfoSource {
	switch metacanon.MakerNoteTag(id) {
	case metacanon.CanonAFInfo:
		return metacanon.AFInfoSourceAFInfo
	case metacanon.CanonAFInfo2:
		return metacanon.AFInfoSourceAFInfo2
	case metacanon.AFInfo3:
		return metacanon.AFInfoSourceAFInfo3
	default:
		return metacanon.AFInfoSourceUnknown
	}
}

func canonShouldReplaceAFInfo(current, candidate metacanon.AFInfo) bool {
	curHas := canonAFInfoHasData(current)
	candHas := canonAFInfoHasData(candidate)
	switch {
	case candHas && !curHas:
		return true
	case !candHas && curHas:
		return false
	case !candHas && !curHas:
		return canonAFInfoSourcePriority(candidate.Source) > canonAFInfoSourcePriority(current.Source)
	}

	curScore := canonAFInfoQualityScore(current)
	candScore := canonAFInfoQualityScore(candidate)
	if candScore != curScore {
		return candScore > curScore
	}

	return canonAFInfoSourcePriority(candidate.Source) > canonAFInfoSourcePriority(current.Source)
}

func canonAFInfoHasData(v metacanon.AFInfo) bool {
	return v.NumAFPoints != 0 ||
		v.ValidAFPoints != 0 ||
		v.CanonImageWidth != 0 ||
		v.CanonImageHeight != 0 ||
		len(v.AFArea) != 0 ||
		len(v.AFPointsInFocusBits) != 0 ||
		len(v.AFPointsSelectedBits) != 0 ||
		v.PrimaryAFPoint != 0
}

func canonAFInfoQualityScore(v metacanon.AFInfo) int {
	score := int(v.NumAFPoints) + int(v.ValidAFPoints)
	score += len(v.AFArea)
	score += len(v.AFPoints)
	score += len(v.AFPointsInFocusBits)
	score += len(v.AFPointsSelectedBits)
	if v.CanonImageWidth != 0 && v.CanonImageHeight != 0 {
		score += 8
	}
	if v.AFImageWidth != 0 && v.AFImageHeight != 0 {
		score += 8
	}
	if v.AFAreaWidth != 0 || v.AFAreaHeight != 0 {
		score += 4
	}
	return score
}

func canonAFInfoSourcePriority(source metacanon.AFInfoSource) int {
	switch source {
	case metacanon.AFInfoSourceAFInfo2:
		return 3
	case metacanon.AFInfoSourceAFInfo3:
		return 2
	case metacanon.AFInfoSourceAFInfo:
		return 1
	default:
		return 0
	}
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

func canonDecodeUniformAFArea(vals []uint16, n, xStart, yStart, count int, w, h int16) []metacanon.AFPoint {
	pointCount := canonRangeLen(n, xStart, count)
	if yLen := canonRangeLen(n, yStart, count); yLen < pointCount {
		pointCount = yLen
	}
	if pointCount == 0 {
		return nil
	}

	areas := make([]metacanon.AFPoint, pointCount)
	for i := 0; i < pointCount; i++ {
		areas[i] = metacanon.NewAFPoint(w, h, int16(vals[xStart+i]), int16(vals[yStart+i]))
	}
	return areas
}

// canonLegacyAFInfoPrimary mirrors Canon.pm sequence handling for AFInfo:
// sequence 11 is either PrimaryAFPoint or an 8-word unknown block, and
// sequence 12 is always PrimaryAFPoint when enough payload remains.
func canonLegacyAFInfoPrimary(vals []uint16, n, seq11Start, afInfoCount int) uint16 {
	if afInfoCount == 36 {
		return canonU16At(vals, n, seq11Start+8)
	}
	if seq11Start+1 < n {
		return vals[seq11Start+1]
	}
	return canonU16At(vals, n, seq11Start)
}

func canonDecodeBitWordsRange(vals []uint16, n, start, count int) []int {
	capHint := canonCountBitWordsRange(vals, n, start, count)
	if capHint == 0 {
		return nil
	}
	out := make([]int, 0, capHint)
	return canonAppendBitWordsRange(out, vals, n, start, count)
}

func canonCountBitWordsRange(vals []uint16, n, start, count int) int {
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
	total := 0
	for i := start; i < end; i++ {
		total += bits.OnesCount16(vals[i])
	}
	return total
}

func canonAppendBitWordsRange(dst []int, vals []uint16, n, start, count int) []int {
	if count <= 0 || start < 0 || start >= n {
		return dst
	}
	end := start + count
	if end > n {
		end = n
	}
	if end <= start {
		return dst
	}

	base := 0
	for i := start; i < end; i++ {
		word := vals[i]
		for word != 0 {
			bit := bits.TrailingZeros16(word)
			dst = append(dst, base+bit)
			word &= word - 1
		}
		base += 16
	}
	return dst
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

func (r *Reader) canonShotInfoLegacyExposureTime() bool {
	model := r.canonModelName()
	if !strings.Contains(model, "EOS 20D") && !strings.Contains(model, "EOS 350D") {
		return false
	}
	return r.makerNoteInfo().Canon.CameraSettings.FocalUnits > 1
}

func canonShotISO(code int16) float32 {
	if code == 0 {
		return 100
	}
	return float32(100.0 * math.Exp2(float64(code-160)/32.0))
}

func canonShotActualISO(autoISO, baseISO float32) float32 {
	if autoISO <= 0 || baseISO <= 0 {
		return 0
	}
	return (autoISO * baseISO) / 100.0
}

func canonShotAperture(code int16) meta.Aperture {
	if code == 0 {
		return 0
	}
	return meta.Aperture(math.Exp2(canonEV(code) * 0.5))
}

func canonShotExposureTime(code int16, legacy20D350D bool) meta.ExposureTime {
	if code == 0 {
		return 0
	}
	if legacy20D350D {
		return meta.ExposureTime(math.Exp2(float64(code-640) / 32.0))
	}
	return meta.ExposureTime(math.Exp2(-canonEV(code)))
}

func canonShotCameraTemperature(raw int16, model string) int16 {
	if raw == 0 || !canonModelIsEOS(model) {
		return 0
	}
	return raw - 128
}

func canonShotFlashGuideNumber(raw int16) float32 {
	if raw < 0 {
		return 0
	}
	return float32(raw) / 32.0
}

// parseCanonMaxAperture converts Canon CameraSettings MaxAperture/MinAperture
// codes to f-numbers using ExifTool's CanonEv conversion.
//
// ExifTool Canon.pm:
//
//	ValueConv => exp(CanonEv($val)*log(2)/2)
func parseCanonMaxAperture(raw uint16) meta.Aperture {
	code := int16(raw)
	if code <= 0 {
		return 0
	}
	ev := canonEV(code)
	return meta.Aperture(math.Exp2(ev * 0.5))
}

// canonEV decodes Canon's hex-based EV codes (modulo 0x20).
func canonEV(code int16) float64 {
	val := int(code)
	sign := 1.0
	if val < 0 {
		val = -val
		sign = -1
	}

	frac := val & 0x1f
	base := val - frac
	fracEV := float64(frac)

	// ExifTool CanonEv special-cases Canon 1/3 and 2/3 encodings.
	switch frac {
	case 0x0c:
		fracEV = 32.0 / 3.0
	case 0x14:
		fracEV = 64.0 / 3.0
	}
	return sign * (float64(base) + fracEV) / 32.0
}

// parseCanonDisplayAperture converts DisplayAperture (sequence 35) as ExifTool:
// RawConv => '$val ? $val : undef', ValueConv => '$val / 10'.
func parseCanonDisplayAperture(raw uint16) meta.Aperture {
	if raw == 0 {
		return 0
	}
	return meta.Aperture(float32(raw) / 10.0)
}

func canonTerminateAtNUL(s string) string {
	start := 0
	for start < len(s) {
		switch s[start] {
		case ' ', '\t', '\n', '\r':
			start++
		default:
			goto findEnd
		}
	}

findEnd:
	end := len(s)
	for i := start; i < end; i++ {
		if s[i] == 0 {
			end = i
			break
		}
	}
	for end > start {
		switch s[end-1] {
		case ' ', '\t', '\n', '\r':
			end--
		default:
			return s[start:end]
		}
	}
	return s[start:end]
}

func canonHexBytes(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	const table = "0123456789abcdef"
	var out strings.Builder
	out.Grow(len(b) * 2)
	for i := range b {
		v := b[i]
		out.WriteByte(table[v>>4])
		out.WriteByte(table[v&0x0f])
	}
	return out.String()
}

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
