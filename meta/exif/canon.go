package exif

import (
	"bytes"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/canon"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) parseCanonTag(t tag.Entry) bool {
	if r.Exif.MakerNote.Canon == nil {
		r.Exif.MakerNote.Canon = &canon.Canon{}
	}
	dst := r.Exif.MakerNote.Canon
	switch canon.MakerNoteTag(t.ID) {
	case canon.CanonImageType:
		dst.ImageType = r.parseStringAllowUndefined(t)
	case canon.CanonFirmwareVersion:
		dst.FirmwareVersion = r.parseStringAllowUndefined(t)
	case canon.CanonFocalLength:
		dst.CanonFocalLength = r.parseCanonFocalLength(t)
	case canon.CanonFlashInfo:
		// dst.FlashInfo = r.parseCanonFlashInfo(t)
	case canon.CanonCameraInfo:
		// intentionally not parsed
	case canon.FileNumber:
		dst.FileNumber = r.parseUint32(t)
	case canon.OwnerName:
		dst.OwnerName = r.parseStringAllowUndefined(t)
	case canon.SerialNumber:
		dst.SerialNumber = r.parseUint32(t)
	case canon.CanonModelID:
		dst.ModelID = r.parseUint32(t)
	case canon.LensModel:
		dst.LensModel = canonTerminateAtNUL(r.parseStringAllowUndefined(t))
	case canon.CanonInternalSerialNumber:
		dst.InternalSerialNumber = r.parseStringAllowUndefined(t)
	case canon.CanonCameraSettings:
		dst.CanonCameraSettings = r.parseCanonCameraSettings(t)
	case canon.CanonShotInfo:
		dst.CanonShotInfo = r.parseCanonShotInfo(t)
	case canon.CanonFileInfo:
		dst.CanonFileInfo = r.parseCanonFileInfo(t)
	case canon.TimeInfo:
		dst.TimeInfo = r.parseCanonTimeInfo(t)
	case canon.BatteryType:
		dst.BatteryType = r.parseCanonBatteryType(t)
	case canon.CanonAFInfo:
		candidate := r.parseCanonAFInfo(t)
		if canonShouldReplaceAFInfo(dst.AFInfo, candidate) {
			dst.AFInfo = candidate
		}
	case canon.CanonAFInfo2, canon.AFInfo3:
		candidate := r.parseCanonAFInfo2(t)
		if canonShouldReplaceAFInfo(dst.AFInfo, candidate) {
			dst.AFInfo = candidate
		}
	case canon.FaceDetect1:
		dst.FaceDetect1 = r.parseCanonFaceDetect1(t)
	case canon.FaceDetect2:
		dst.FaceDetect2 = r.parseCanonFaceDetect2(t)
	case canon.FaceDetect3:
		dst.FaceDetect3 = r.parseCanonFaceDetect3(t)
	case canon.ImageUniqueID:
		dst.ImageUniqueID = r.parseCanonImageUniqueID(t)
	case canon.CanonCustomFunctions:
		// TODO(canon): Expand Canon maker-note parity with ExifTool's
		// CanonCustom tables and remaining Canon-specific fields.
		// intentionally not parsed
	case canon.CanonAspectInfo:
		dst.AspectInfo = r.parseCanonAspectInfo(t)
	case canon.CanonProcessingInfo:
		dst.ProcessingInfo = r.parseCanonProcessingInfo(t)
	case canon.CanonColorSpace:
		dst.ColorSpace = r.parseUint16(t)
	case canon.CanonPreviewImageInfo:
		dst.PreviewImageInfo = r.parseCanonPreviewImageInfo(t)
	case canon.CanonSensorInfo:
		dst.SensorInfo = r.parseCanonSensorInfo(t)
	case canon.CanonPictureStyleUserDef:
		// intentionally not parsed
	case canon.CanonPictureStylePC:
		// intentionally not parsed
	case canon.CanonCustomPictureStyleFileName:
		dst.CustomPictureStyleFileName = r.parseStringAllowUndefined(t)
	case canon.CanonAFMicroAdj:
		dst.AFMicroAdj = r.parseCanonAFMicroAdj(t)
	case canon.CanonLightingOpt:
		dst.LightingOpt = r.parseCanonLightingOpt(t)
	case canon.CanonLensInfo:
		dst.LensInfo = r.parseCanonLensInfo(t)
	case canon.CanonMultiExp:
		dst.MultiExp = r.parseCanonMultiExp(t)
	case canon.CanonHDRInfo:
		dst.HDRInfo = r.parseCanonHDRInfo(t)
	case canon.CanonAFConfig:
		dst.AFConfig = r.parseCanonAFConfig(t)
	case canon.CanonRawBurstModeRoll:
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

type canonSeq16 []uint16

func (s canonSeq16) present(seq int) bool {
	idx := seq - 1
	return idx >= 0 && idx < len(s)
}

func (s canonSeq16) u16(seq int) uint16 {
	if !s.present(seq) {
		return 0
	}
	return s[seq-1]
}

func (s canonSeq16) i16(seq int) int16 {
	return int16(s.u16(seq))
}

type canonFirstEntry1Int32 []int32

func (s canonFirstEntry1Int32) present(seq int) bool {
	return seq > 0 && seq < len(s)
}

func (s canonFirstEntry1Int32) i32(seq int) int32 {
	if !s.present(seq) {
		return 0
	}
	return s[seq]
}

func (s canonFirstEntry1Int32) u32(seq int) uint32 {
	return uint32(s.i32(seq))
}

func (r *Reader) parseCanonSizedUint16Payload(t tag.Entry, parser string, dst []uint16) (canonSeq16, bool) {
	n := r.parseCanonUint16List(t, dst)
	if n < 1 {
		r.warnCanonShortRead(t, parser, n, 1)
		return nil, false
	}
	if uint32(dst[0]) != t.Size() {
		r.warnCanonInvalidSize(t, parser, uint32(dst[0]))
		return nil, false
	}
	if n < 2 {
		return nil, false
	}
	return canonSeq16(dst[1:n]), true
}

func (r *Reader) parseCanonBlockPreview(t tag.Entry) canon.BlockPreview {
	dst := canon.BlockPreview{Size: t.Size()}
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
			r.tagLogContext(r.warn(), t).
				Err(err).
				Str("parser", "parseCanonBlockPreview").
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

func (r *Reader) parseCanonPreviewImageInfo(t tag.Entry) canon.PreviewImageInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 5 {
		r.warnCanonShortRead(t, "parseCanonPreviewImageInfo", n, 5)
		return canon.PreviewImageInfo{}
	}

	return canon.PreviewImageInfo{
		PreviewQuality:     canon.Quality(int16(raw[0])),
		PreviewImageLength: uint32(raw[1]),
		PreviewImageWidth:  uint32(raw[2]),
		PreviewImageHeight: uint32(raw[3]),
		PreviewImageStart:  uint32(raw[4]),
	}
}

func (r *Reader) parseCanonSensorInfo(t tag.Entry) canon.SensorInfo {
	var raw [13]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 13 {
		r.warnCanonShortRead(t, "parseCanonSensorInfo", n, 13)
		return canon.SensorInfo{}
	}
	return canon.SensorInfo{
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

func (r *Reader) parseCanonAFConfig(t tag.Entry) canon.AFConfig {
	var raw [25]int32
	n := r.parseCanonInt32List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonAFConfig", n, 2)
		return canon.AFConfig{}
	}
	settings := canonFirstEntry1Int32(raw[:n])
	dst := canon.AFConfig{
		AFConfigTool:              settings.u32(1) + 1,
		AFTrackingSensitivity:     settings.i32(2),
		AFAccelDecelTracking:      settings.i32(3),
		AFPointSwitching:          settings.i32(4),
		AIServoFirstImage:         settings.i32(5),
		AIServoSecondImage:        settings.i32(6),
		USMLensElectronicMF:       settings.i32(7),
		AFAssistBeam:              settings.i32(8),
		OneShotAFRelease:          settings.i32(9),
		AutoAFPointSelEOSiTRAF:    settings.i32(10),
		LensDriveWhenAFImpossible: settings.i32(11),
		SelectAFAreaSelectionMode: settings.u32(12),
		AFAreaSelectionMethod:     settings.i32(13),
		OrientationLinkedAF:       settings.i32(14),
		ManualAFPointSelPattern:   settings.i32(15),
		AFPointDisplayDuringFocus: settings.i32(16),
		VFDisplayIllumination:     settings.i32(17),
		AFStatusViewfinder:        settings.i32(18),
		InitialAFPointInServo:     settings.i32(19),
		SubjectToDetect:           settings.i32(20),
		EyeDetection:              settings.i32(24),
	}
	return dst
}

func (r *Reader) parseCanonRawBurstInfo(t tag.Entry) canon.RawBurstInfo {
	var raw [3]uint32
	if n := r.parseUint32List(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonRawBurstInfo", n, 3)
		return canon.RawBurstInfo{}
	}
	return canon.RawBurstInfo{
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

func (r *Reader) parseCanonFaceDetect1(t tag.Entry) canon.FaceDetect1Info {
	var raw [26]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 5 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect1", n, 5)
		return canon.FaceDetect1Info{}
	}
	dst := canon.FaceDetect1Info{
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
		dst.FacePositions[i] = canon.FacePosition{
			X: int16(raw[start]),
			Y: int16(raw[start+1]),
		}
	}
	return dst
}

func (r *Reader) parseCanonFaceDetect2(t tag.Entry) canon.FaceDetect2Info {
	var raw [8]byte
	if n := r.parseByteList(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect2", n, 3)
		return canon.FaceDetect2Info{}
	}
	return canon.FaceDetect2Info{
		FaceWidth:     raw[1],
		FacesDetected: raw[2],
	}
}

func (r *Reader) parseCanonFaceDetect3(t tag.Entry) canon.FaceDetect3Info {
	var raw [8]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect3", n, 4)
		return canon.FaceDetect3Info{}
	}
	return canon.FaceDetect3Info{
		FacesDetected: raw[3],
	}
}

// parseCanonFocalLength parses tag 0x0002 (CanonFocalLength).
func (r *Reader) parseCanonFocalLength(t tag.Entry) canon.FocalLengthInfo {
	var raw [8]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonFocalLength", n, 4)
		return canon.FocalLengthInfo{}
	}
	return canon.FocalLengthInfo{
		FocalType:       raw[0],
		FocalLength:     raw[1],
		FocalPlaneXSize: raw[2],
		FocalPlaneYSize: raw[3],
	}
}

// parseCanonAspectInfo parses tag 0x009a (AspectInfo).
func (r *Reader) parseCanonAspectInfo(t tag.Entry) canon.AspectInfo {
	var raw [8]uint32
	if n := r.parseUint32List(t, raw[:]); n < 5 {
		r.warnCanonShortRead(t, "parseCanonAspectInfo", n, 5)
		return canon.AspectInfo{}
	}
	return canon.AspectInfo{
		AspectRatio:        raw[0],
		CroppedImageWidth:  raw[1],
		CroppedImageHeight: raw[2],
		CroppedImageLeft:   raw[3],
		CroppedImageTop:    raw[4],
	}
}

// parseCanonProcessingInfo parses tag 0x00a0 (ProcessingInfo).
func (r *Reader) parseCanonProcessingInfo(t tag.Entry) canon.ProcessingInfo {
	var raw [24]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonProcessingInfo", n, 2)
		return canon.ProcessingInfo{}
	}
	return canon.ProcessingInfo{
		// ExifTool ProcessingInfo uses FIRST_ENTRY => 1, so raw[0] is the size word.
		// The payload length varies by model, so decode conditionally.
		ToneCurve:            canonI16At(raw[:], n, 1),
		Sharpness:            canonI16At(raw[:], n, 2),
		SharpnessFrequency:   canonI16At(raw[:], n, 3),
		SensorRedLevel:       canonI16At(raw[:], n, 4),
		SensorBlueLevel:      canonI16At(raw[:], n, 5),
		WhiteBalanceRed:      canonI16At(raw[:], n, 6),
		WhiteBalanceBlue:     canonI16At(raw[:], n, 7),
		WhiteBalance:         canonI16At(raw[:], n, 8),
		ColorTemperature:     canonI16At(raw[:], n, 9),
		PictureStyle:         canonI16At(raw[:], n, 10),
		DigitalGain:          canonI16At(raw[:], n, 11),
		WBShiftAB:            canonI16At(raw[:], n, 12),
		WBShiftGM:            canonI16At(raw[:], n, 13),
		UnsharpMaskFineness:  canonI16At(raw[:], n, 14),
		UnsharpMaskThreshold: canonI16At(raw[:], n, 15),
	}

}

// parseCanonAFMicroAdj parses tag 0x4013 (AFMicroAdj).
func (r *Reader) parseCanonAFMicroAdj(t tag.Entry) canon.AFMicroAdjInfo {
	var raw [8]int32
	n := r.parseCanonInt32List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonAFMicroAdj", n, 2)
		return canon.AFMicroAdjInfo{}
	}
	dst := canon.AFMicroAdjInfo{
		Mode: raw[1],
	}
	if n > 2 {
		dst.ValueNumerator = raw[2]
	}
	if n > 3 {
		dst.ValueDenominator = raw[3]
	}
	return dst
}

// parseCanonLightingOpt parses tag 0x4018 (LightingOpt).
func (r *Reader) parseCanonLightingOpt(t tag.Entry) canon.LightingOptInfo {
	var raw [12]int32
	n := r.parseCanonInt32List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonLightingOpt", n, 2)
		return canon.LightingOptInfo{}
	}

	dst := canon.LightingOptInfo{
		// ExifTool LightingOpt table uses FIRST_ENTRY=1.
		PeripheralIlluminationCorr: raw[1],
	}
	if n > 2 {
		dst.AutoLightingOptimizer = raw[2]
	}
	if n > 3 {
		dst.HighlightTonePriority = raw[3]
	}
	if n > 4 {
		dst.LongExposureNoiseReduction = raw[4]
	}
	if n > 5 {
		dst.HighISONoiseReduction = raw[5]
	}
	if n > 10 {
		dst.DigitalLensOptimizer = raw[10]
	}
	if n > 11 {
		dst.DualPixelRaw = raw[11]
	}
	return dst
}

const canonLensInfoByteLength = 5

// parseCanonLensInfo parses tag 0x4019 (LensInfoForService).
func (r *Reader) parseCanonLensInfo(t tag.Entry) canon.LensInfoForService {
	dst := canon.LensInfoForService{}
	raw := r.parseOpaqueBytes(t, canonLensInfoByteLength)
	l := int(len(raw))
	if l != 5 {
		r.warnCanonShortRead(t, "parseCanonLensInfo", l, int(t.Size()))
		return canon.LensInfoForService{}
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
func (r *Reader) parseCanonMultiExp(t tag.Entry) canon.MultiExpInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonMultiExp", n, 4)
		return canon.MultiExpInfo{}
	}
	return canon.MultiExpInfo{
		// ExifTool MultiExp table uses FIRST_ENTRY=1.
		MultiExposure:        raw[1],
		MultiExposureControl: raw[2],
		MultiExposureShots:   raw[3],
	}
}

// parseCanonHDRInfo parses tag 0x4025 (HDRInfo).
func (r *Reader) parseCanonHDRInfo(t tag.Entry) canon.HDRInfo {
	var raw [8]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonHDRInfo", n, 3)
		return canon.HDRInfo{}
	}
	return canon.HDRInfo{
		// ExifTool HDRInfo table uses FIRST_ENTRY=1.
		HDR:       raw[1],
		HDREffect: raw[2],
	}
}

// parseCanonCameraSettings parses tag 0x0001 (CanonCameraSettings).
func (r *Reader) parseCanonCameraSettings(t tag.Entry) canon.CameraSettings {
	var raw [53]uint16
	settings, ok := r.parseCanonSizedUint16Payload(t, "parseCanonCameraSettings", raw[:])
	if !ok {
		return canon.CameraSettings{}
	}

	return canon.CameraSettings{
		MacroMode:         canon.MacroMode(settings.u16(1)),
		SelfTimer:         settings.i16(2),
		Quality:           canon.Quality(settings.i16(3)),
		CanonFlashMode:    canon.CanonFlashMode(settings.i16(4)),
		ContinuousDrive:   canon.ContinuousDrive(settings.i16(5)),
		FocusMode:         canon.FocusMode(settings.i16(7)),
		RecordMode:        canon.RecordMode(settings.i16(9)),
		CanonImageSize:    canon.CanonImageSize(settings.i16(10)),
		EasyMode:          canon.EasyMode(settings.i16(11)),
		DigitalZoom:       canon.DigitalZoom(settings.i16(12)),
		Contrast:          settings.i16(13),
		Saturation:        settings.i16(14),
		Sharpness:         settings.i16(15),
		CameraISO:         canon.CameraISO(settings.i16(16)),
		MeteringMode:      canon.MeteringMode(settings.i16(17)),
		FocusRange:        canon.FocusRange(settings.i16(18)),
		AFPoint:           settings.u16(19),
		CanonExposureMode: canon.ExposureMode(settings.i16(20)),
		LensType:          canon.CanonLensType(settings.u16(22)),
		MaxFocalLength:    settings.u16(23),
		MinFocalLength:    settings.u16(24),
		FocalUnits:        settings.u16(25),
		MaxAperture:       parseCanonMaxAperture(settings.u16(26)),
		MinAperture:       parseCanonMaxAperture(settings.u16(27)),
		FlashModel:        canon.FlashModel(settings.i16(28)),
		FlashBits:         settings.u16(29),
		FocusContinuous:   canon.FocusContinuous(settings.i16(32)),
		AESetting:         canon.AESetting(settings.i16(33)),
		ImageStabilization: canon.ImageStabilization(
			settings.i16(34),
		),
		DisplayAperture:   parseCanonDisplayAperture(settings.u16(35)),
		ZoomSourceWidth:   settings.u16(36),
		ZoomTargetWidth:   settings.u16(37),
		SpotMeteringMode:  canon.SpotMeteringMode(settings.i16(39)),
		PhotoEffect:       canon.PhotoEffect(settings.i16(40)),
		ManualFlashOutput: canon.ManualFlashOutput(settings.i16(41)),
		ColorTone:         settings.i16(42),
		SRAWQuality:       canon.SRAWQuality(settings.i16(46)),
		FocusBracketing:   canon.FocusBracketing(settings.i16(50)),
		Clarity:           settings.i16(51),
		HDRPQ:             canon.HDRPQ(settings.u16(52)),
	}
}

// parseCanonShotInfo parses tag 0x0004 (CanonShotInfo).
func (r *Reader) parseCanonShotInfo(t tag.Entry) canon.ShotInfo {
	var raw [64]uint16
	settings, ok := r.parseCanonSizedUint16Payload(t, "parseCanonShotInfo", raw[:])
	if !ok {
		return canon.ShotInfo{}
	}
	dst := canon.ShotInfo{
		AutoISO:                settings.i16(1),
		BaseISO:                settings.i16(2),
		MeasuredEV:             settings.i16(3),
		TargetAperture:         settings.i16(4),
		TargetExposureTime:     settings.i16(5),
		ExposureCompensation:   settings.i16(6),
		WhiteBalance:           canon.WhiteBalance(settings.i16(7)),
		SlowShutter:            canon.SlowShutter(settings.i16(8)),
		SequenceNumber:         settings.i16(9),
		OpticalZoomCode:        settings.i16(10),
		CameraTemperature:      settings.i16(12),
		FlashGuideNumber:       settings.i16(13),
		AFPointsInFocus:        settings.u16(14),
		FlashExposureComp:      settings.i16(15),
		AutoExposureBracketing: settings.i16(16),
		AEBBracketValue:        settings.i16(17),
		ControlMode:            settings.i16(18),
		FNumber:                settings.i16(21),
		ExposureTime:           settings.i16(22),
		MeasuredEV2:            settings.i16(23),
		BulbDuration:           settings.i16(24),
		CameraType:             canon.CameraType(settings.i16(26)),
		AutoRotate:             canon.AutoRotate(settings.i16(27)),
		NDFilter:               canon.NDFilter(settings.i16(28)),
		SelfTimer2:             settings.i16(29),
		FlashOutput:            settings.i16(33),
	}
	dst.AutoISOValue = canonShotISO(dst.AutoISO)
	dst.BaseISOValue = canonShotISO(dst.BaseISO)
	dst.ActualISO = canonShotActualISO(dst.AutoISOValue, dst.BaseISOValue)
	dst.TargetApertureValue = canonShotAperture(dst.TargetAperture)
	dst.TargetExposureTimeValue = canonShotExposureTime(dst.TargetExposureTime, false)
	dst.CameraTemperatureC = canonShotCameraTemperature(dst.CameraTemperature, r.canonModelName())
	dst.FlashGuideNumberMeters = canonShotFlashGuideNumber(dst.FlashGuideNumber)
	if settings.present(20) && settings.u16(19) != 0 {
		dst.FocusDistance = canon.NewFocusDistance(settings.u16(19), settings.u16(20))
	}
	dst.FNumberValue = canonShotAperture(dst.FNumber)
	dst.ExposureTimeValue = canonShotExposureTime(dst.ExposureTime, r.canonShotInfoLegacyExposureTime())
	return dst
}

// parseCanonFileInfo parses tag 0x0093 (CanonFileInfo).
func (r *Reader) parseCanonFileInfo(t tag.Entry) canon.FileInfo {
	var raw [64]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 2 {
		r.warnCanonShortRead(t, "parseCanonFileInfo", n, 2)
		return canon.FileInfo{}
	}

	// Tag 0x0093 index 1 is model-dependent (FileNumber or ShutterCount).
	// Preserve raw 32-bit representation for both fields.
	return canon.FileInfo{
		FileNumber:                  uint32(canonU16At(raw[:], n, 1)) | (uint32(canonU16At(raw[:], n, 2)) << 16),
		BracketMode:                 canon.BracketMode(canonI16At(raw[:], n, 3)),
		BracketValue:                canonI16At(raw[:], n, 4),
		BracketShotNumber:           canonI16At(raw[:], n, 5),
		RawJpgQuality:               canon.RawJpgQuality(canonU16At(raw[:], n, 6)),
		RawJpgSize:                  canon.RawJpgSize(canonU16At(raw[:], n, 7)),
		LongExposureNoiseReduction2: canon.OnOffAuto(canonU16At(raw[:], n, 8)),
		WBBracketMode:               canonI16At(raw[:], n, 9),
		WBBracketValueAB:            canonI16At(raw[:], n, 12),
		WBBracketValueGM:            canonI16At(raw[:], n, 13),
		FilterEffect:                canon.FilterEffect(canonU16At(raw[:], n, 14)),
		ToningEffect:                canon.ToningEffect(canonU16At(raw[:], n, 15)),
		MacroMagnification:          canonI16At(raw[:], n, 16),
		LiveViewShooting:            canon.OnOffAuto(canonU16At(raw[:], n, 19)),
		FocusDistance:               canon.NewFocusDistance(canonU16At(raw[:], n, 20), canonU16At(raw[:], n, 21)),
		ShutterMode:                 canon.ShutterMode(canonU16At(raw[:], n, 23)),
		FlashExposureLock:           canon.OnOffAuto(canonU16At(raw[:], n, 25)),
		AntiFlicker:                 canon.OnOffAuto(canonU16At(raw[:], n, 32)),
		RFLensType:                  canon.CanonRFLensType(canonU16At(raw[:], n, 61)),
	}
}

// parseCanonTimeInfo parses tag 0x0035 (TimeInfo).
func (r *Reader) parseCanonTimeInfo(t tag.Entry) canon.CanonTimeInfo {
	var raw [4]int32
	if n := r.parseCanonInt32List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonTimeInfo", n, 4)
		return canon.CanonTimeInfo{}
	}
	return canon.CanonTimeInfo{
		TimeZone:        raw[1],
		TimeZoneCity:    canon.TimeZoneCity(raw[2]),
		DaylightSavings: canon.DaylightSavings(raw[3]),
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
			r.tagLogContext(r.warn(), t).
				Str("parser", "parseCanonBatteryType").
				Uint32("sizeBytes", t.Size()).
				Uint32("wantSizeBytes", canonBatteryTypePayloadSize).
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
func (r *Reader) parseCanonAFInfo(t tag.Entry) canon.AFInfo {
	var wordsStack [2048]uint16
	words, truncated := canonAFWordsBuffer(wordsStack[:], t.UnitCount)
	if truncated {
		r.warnCanonTruncatedWords(t, "parseCanonAFInfo", len(words), int(t.UnitCount))
	}
	n := r.parseCanonUint16List(t, words)
	source := canonAFInfoSource(tag.ID(canon.CanonAFInfo))
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonAFInfo", n, 1)
		return canon.AFInfo{Source: source}
	}
	return canon.DecodeAFInfo(words[:n], canonModelIsEOS(r.canonModelName()), int(t.UnitCount))
}

func fillCanonAFInfo(dst *canon.AFInfo, words []uint16, model string, afInfoCount int) {
	*dst = canon.DecodeAFInfo(words, canonModelIsEOS(model), afInfoCount)
}

// parseCanonAFInfo2 parses tags 0x0026 and 0x003c (AFInfo2/AFInfo3).
func (r *Reader) parseCanonAFInfo2(t tag.Entry) canon.AFInfo {
	var wordsStack [2048]uint16
	words, truncated := canonAFWordsBuffer(wordsStack[:], t.UnitCount)
	if truncated {
		r.warnCanonTruncatedWords(t, "parseCanonAFInfo2", len(words), int(t.UnitCount))
	}
	n := r.parseCanonUint16List(t, words)
	source := canonAFInfoSource(t.ID)
	if n == 0 {
		r.warnCanonShortRead(t, "parseCanonAFInfo2", n, 1)
		return canon.AFInfo{Source: source}
	}
	model := r.canonModelName()
	isAFInfo3 := canon.MakerNoteTag(t.ID) == canon.AFInfo3
	return canon.DecodeAFInfo2(words[:n], canon.AFInfo2DecodeConfig{
		Source:         source,
		EOS:            canonModelIsEOS(model),
		AFInfo3:        isAFInfo3,
		DecodeCoords:   r.afInfoDecodeOptions.has(AFInfoDecodeCoords),
		DecodePoints:   r.afInfoDecodeOptions.has(AFInfoDecodePoints),
		DecodeInFocus:  r.afInfoDecodeOptions.has(AFInfoDecodeInFocus),
		DecodeSelected: r.afInfoDecodeOptions.has(AFInfoDecodeSelected),
	})
}

func (r *Reader) warnCanonTruncatedWords(t tag.Entry, parser string, got, want int) {
	if !r.warnEnabled() {
		return
	}
	r.tagLogContext(r.warn(), t).
		Str("parser", parser).
		Int("wordsDecoded", got).
		Int("wordsRequested", want).
		Msg("canon AF payload truncated to parser word cap")
}

func (r *Reader) warnCanonShortRead(t tag.Entry, parser string, got, want int) {
	if !r.warnEnabled() {
		return
	}
	r.tagLogContext(r.warn(), t).
		Str("parser", parser).
		Int("gotUnits", got).
		Int("wantUnits", want).
		Msg("canon maker-note payload too short")
}

func (r *Reader) warnCanonInvalidSize(t tag.Entry, parser string, declaredSizeBytes uint32) {
	if !r.warnEnabled() {
		return
	}
	r.tagLogContext(r.warn(), t).
		Str("parser", parser).
		Uint32("declaredSizeBytes", declaredSizeBytes).
		Uint32("actualSizeBytes", t.Size()).
		Msg("invalid canon maker-note payload length")
}
