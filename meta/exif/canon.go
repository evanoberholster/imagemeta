package exif

import (
	"math/bits"
	"strings"

	"github.com/evanoberholster/imagemeta/imagetype"
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
		dst.ImageType = r.parseCanonString(t)
	case canon.CanonFirmwareVersion:
		dst.FirmwareVersion = canon.NormalizeFirmwareVersion(r.parseCanonString(t))
	case canon.CanonFocalLength:
		dst.CanonFocalLength = r.parseCanonFocalLength(t)
	case canon.CanonFlashInfo:
		dst.FlashInfo = r.parseCanonFlashInfo(t)
	case canon.CanonCameraInfo:
		dst.CameraInfo = r.parseCanonCameraInfo(t)
	case canon.FileNumber:
		dst.FileNumber = r.parseUint32(t)
	case canon.OwnerName:
		dst.OwnerName = r.parseCanonString(t)
	case canon.SerialNumber:
		dst.SerialNumber = r.parseUint32(t)
	case canon.CanonModelID:
		dst.ModelID = r.parseUint32(t)
	case canon.LensModel:
		dst.LensModel = r.parseCanonString(t)
	case canon.CanonInternalSerialNumber:
		dst.InternalSerialNumber = r.parseCanonString(t)
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
	case canon.AFPointsInFocus1D:
		dst.AFInfo = r.parseCanonAFPointsInFocus1D(t, dst.AFInfo)
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
	case canon.CanonColorTemperature:
		dst.ColorTemperature = r.parseUint16(t)
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
		dst.CustomPictureStyleFileName = r.parseCanonString(t)
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

func (r *Reader) parseCanonCameraInfo(t tag.Entry) canon.CameraInfo {
	layout := r.canonCameraInfoLayout(t)
	switch layout {
	case canon.CameraInfoLayout5D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout5D)
	case canon.CameraInfoLayout5DmkII:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout5DmkII)
	case canon.CameraInfoLayout5DmkIII:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout5DmkIII)
	case canon.CameraInfoLayout6D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout6D)
	case canon.CameraInfoLayout7D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout7D)
	case canon.CameraInfoLayout40D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout40D)
	case canon.CameraInfoLayout50D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout50D)
	case canon.CameraInfoLayout60D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout60D)
	case canon.CameraInfoLayout70D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout70D)
	case canon.CameraInfoLayout80D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout80D)
	case canon.CameraInfoLayout450D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout450D)
	case canon.CameraInfoLayout500D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout500D)
	case canon.CameraInfoLayout550D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout550D)
	case canon.CameraInfoLayout600D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout600D)
	case canon.CameraInfoLayout650D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout650D)
	case canon.CameraInfoLayout700D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout700D)
	case canon.CameraInfoLayout750D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout750D)
	case canon.CameraInfoLayout1000D:
		return r.parseCanonCameraInfoBytes(t, canon.CameraInfoSpecLayout1000D)
	case canon.CameraInfoLayoutPowerShot:
		return r.parseCanonCameraInfoPowerShot(t)
	case canon.CameraInfoLayoutPowerShot2:
		return r.parseCanonCameraInfoPowerShot2(t)
	case canon.CameraInfoLayoutUnknown32:
		return r.parseCanonCameraInfoUnknown32(t)
	case canon.CameraInfoLayoutR6:
		if v, ok := r.readCanonCameraInfoUint32At(t, 0x0af1); ok {
			return canon.CameraInfo{ShutterCount: v}
		}
	case canon.CameraInfoLayoutR6m2:
		if v, ok := r.readCanonCameraInfoUint32At(t, 0x0d29); ok {
			return canon.CameraInfo{ShutterCount: v}
		}
	case canon.CameraInfoLayoutR6m3:
		if v, ok := r.readCanonCameraInfoUint16At(t, 0x086d); ok {
			return canon.CameraInfo{ImageCount: uint32(v)}
		}
	default:
	}
	return canon.CameraInfo{}
}

func (r *Reader) canonCameraInfoLayout(t tag.Entry) canon.CameraInfoLayout {
	if layout, ok := canon.CameraInfoLayoutForModelID(r.canonModelID()); ok {
		return layout
	}
	if layout, ok := canon.CameraInfoLayoutForModelName(r.canonModelName()); ok {
		return layout
	}
	if t.Type != tag.TypeLong {
		return canon.CameraInfoLayoutUnknown
	}
	switch t.UnitCount {
	case 138, 148:
		return canon.CameraInfoLayoutPowerShot
	case 156, 162, 167, 171, 264:
		return canon.CameraInfoLayoutPowerShot2
	default:
		return canon.CameraInfoLayoutUnknown32
	}
}

func (r *Reader) parseCanonCameraInfoBytes(t tag.Entry, spec canon.CameraInfoSpec) canon.CameraInfo {
	buf, _, err := r.readTagBytes(t, t.Size())
	if err != nil || len(buf) == 0 {
		return canon.CameraInfo{}
	}
	return canon.CameraInfoDecode(buf, spec)
}

func (r *Reader) parseCanonCameraInfoPowerShot(t tag.Entry) canon.CameraInfo {
	return r.parseCanonCameraInfoTempFromCount(t,
		135, 138,
		145, 148,
	)
}

func (r *Reader) parseCanonCameraInfoPowerShot2(t tag.Entry) canon.CameraInfo {
	return r.parseCanonCameraInfoTempFromCount(t,
		153, 156,
		159, 162,
		164, 167,
		168, 171,
		261, 264,
	)
}

func (r *Reader) parseCanonCameraInfoUnknown32(t tag.Entry) canon.CameraInfo {
	count := int(t.UnitCount)
	var tempIndex int
	switch count {
	case 72:
		tempIndex = 71
	case 85:
		tempIndex = 83
	case 93, 94:
		tempIndex = 91
	case 96:
		tempIndex = 92
	case 104:
		tempIndex = 100
	default:
		if count > 400 {
			tempIndex = count - 3
		} else {
			return canon.CameraInfo{}
		}
	}
	if v, ok := r.readCanonCameraInfoInt32At(t, tempIndex*4); ok {
		return canon.CameraInfo{CameraTemperature: int16(v)}
	}
	return canon.CameraInfo{}
}

// parseCanonCameraInfoTempFromCount extracts camera temperature from a
// PowerShot int32 payload. Pairs of (tempIndex, unitCount) are checked.
func (r *Reader) parseCanonCameraInfoTempFromCount(t tag.Entry, pairs ...int) canon.CameraInfo {
	count := int(t.UnitCount)
	for i := 0; i+1 < len(pairs); i += 2 {
		if count == pairs[i+1] {
			if v, ok := r.readCanonCameraInfoInt32At(t, pairs[i]*4); ok {
				return canon.CameraInfo{CameraTemperature: int16(v)}
			}
			return canon.CameraInfo{}
		}
	}
	return canon.CameraInfo{}
}

func (r *Reader) readCanonCameraInfoUint32At(t tag.Entry, off int) (uint32, bool) {
	b, err := r.readCanonTagOffsetBytes(t, off, 4)
	if err != nil || len(b) < 4 {
		return 0, false
	}
	return canon.CIU32LEAt(b, 0), true
}

func (r *Reader) readCanonCameraInfoUint16At(t tag.Entry, off int) (uint16, bool) {
	b, err := r.readCanonTagOffsetBytes(t, off, 2)
	if err != nil || len(b) < 2 {
		return 0, false
	}
	return canon.CIU16LEAt(b, 0), true
}

func (r *Reader) readCanonCameraInfoInt32At(t tag.Entry, off int) (int32, bool) {
	v, ok := r.readCanonCameraInfoUint32At(t, off)
	if !ok {
		return 0, false
	}
	return int32(v), true
}

func (r *Reader) readCanonTagOffsetBytes(t tag.Entry, off, n int) ([]byte, error) {
	size := int(t.Size())
	if off < 0 || n <= 0 || off > size-n {
		return nil, imagetype.ErrDataLength
	}
	if err := r.seekToTag(t); err != nil {
		return nil, err
	}
	if err := r.discard(off); err != nil {
		return nil, err
	}
	buf, err := r.fastRead(n)
	if err != nil {
		return nil, err
	}
	remaining := int(t.Size()) - off - len(buf)
	if remaining > 0 {
		if err := r.discard(remaining); err != nil {
			return nil, err
		}
	}
	return buf, nil
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

// parseCanonSizedUint16Payload validates Canon tables that begin with a byte
// size word, then returns a 1-based sequence view matching ExifTool indices.
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
		if r.WarnEnabled() {
			r.tagLogContext(r.Warn(3), t).
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
	n := r.parseCanonInt32List(t, raw[:])
	start := 0
	if n >= 6 && (raw[0] == int32(t.Size()) || (raw[0] > 5 && raw[1] >= 0 && raw[1] <= 5 && raw[2] > 0xffff)) {
		start = 1
	}
	if n-start < 5 {
		r.warnCanonShortRead(t, "parseCanonPreviewImageInfo", n, 5)
		return canon.PreviewImageInfo{}
	}

	dst := canon.PreviewImageInfo{
		PreviewQuality:     canon.Quality(int16(raw[start])),
		PreviewImageLength: uint32(raw[start+1]),
		PreviewImageWidth:  uint32(raw[start+2]),
		PreviewImageHeight: uint32(raw[start+3]),
		PreviewImageStart:  uint32(raw[start+4]),
	}
	if start == 1 && dst.PreviewImageStart != 0 {
		dst.PreviewImageStart += r.tiffHeaderOffset
	}
	return dst
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
	return canon.AFConfig{
		AFConfigTool:              uint32(raw[1]) + 1,
		AFTrackingSensitivity:     raw[2],
		AFAccelDecelTracking:      raw[3],
		AFPointSwitching:          raw[4],
		AIServoFirstImage:         raw[5],
		AIServoSecondImage:        raw[6],
		USMLensElectronicMF:       raw[7],
		AFAssistBeam:              raw[8],
		OneShotAFRelease:          raw[9],
		AutoAFPointSelEOSiTRAF:    raw[10],
		LensDriveWhenAFImpossible: raw[11],
		SelectAFAreaSelectionMode: uint32(raw[12]),
		AFAreaSelectionMethod:     raw[13],
		OrientationLinkedAF:       raw[14],
		ManualAFPointSelPattern:   raw[15],
		AFPointDisplayDuringFocus: raw[16],
		VFDisplayIllumination:     raw[17],
		AFStatusViewfinder:        raw[18],
		InitialAFPointInServo:     raw[19],
		SubjectToDetect:           raw[20],
		EyeDetection:              raw[24],
	}
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
// parseCanonImageUniqueID parses Canon maker-note tag 0x0028 into meta.UUID.
const canonUUIDBytesLength = 16

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
		FocalPlaneXSize: canon.FocalPlaneSizeMM(raw[2]),
		FocalPlaneYSize: canon.FocalPlaneSizeMM(raw[3]),
	}
}

func (r *Reader) parseCanonFlashInfo(t tag.Entry) canon.FlashInfo {
	return canon.FlashInfo{Raw: r.parseCanonBlockPreview(t)}
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
	if dst.Raw[0] == 0 && dst.Raw[1] == 0 && dst.Raw[2] == 0 && dst.Raw[3] == 0 {
		return dst
	}
	dst.LensSerialNumber = canon.HexBytes(dst.Raw[:n])
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

// Canon CameraSettings uses ExifTool FIRST_ENTRY=1 indexing. The payload's
// size word is stripped by parseCanonSizedUint16Payload before these positions
// are read.
const (
	cameraSettingsMacroMode          = 1
	cameraSettingsSelfTimer          = 2
	cameraSettingsQuality            = 3
	cameraSettingsFlashMode          = 4
	cameraSettingsContinuousDrive    = 5
	cameraSettingsFocusMode          = 7
	cameraSettingsRecordMode         = 9
	cameraSettingsImageSize          = 10
	cameraSettingsEasyMode           = 11
	cameraSettingsDigitalZoom        = 12
	cameraSettingsContrast           = 13
	cameraSettingsSaturation         = 14
	cameraSettingsSharpness          = 15
	cameraSettingsCameraISO          = 16
	cameraSettingsMeteringMode       = 17
	cameraSettingsFocusRange         = 18
	cameraSettingsAFPoint            = 19
	cameraSettingsExposureMode       = 20
	cameraSettingsLensType           = 22
	cameraSettingsMaxFocalLength     = 23
	cameraSettingsMinFocalLength     = 24
	cameraSettingsFocalUnits         = 25
	cameraSettingsMaxAperture        = 26
	cameraSettingsMinAperture        = 27
	cameraSettingsFlashModel         = 28
	cameraSettingsFlashBits          = 29
	cameraSettingsFocusContinuous    = 32
	cameraSettingsAESetting          = 33
	cameraSettingsImageStabilization = 34
	cameraSettingsDisplayAperture    = 35
	cameraSettingsZoomSourceWidth    = 36
	cameraSettingsZoomTargetWidth    = 37
	cameraSettingsSpotMeteringMode   = 39
	cameraSettingsPhotoEffect        = 40
	cameraSettingsManualFlashOutput  = 41
	cameraSettingsColorTone          = 42
	cameraSettingsSRAWQuality        = 46
	cameraSettingsFocusBracketing    = 50
	cameraSettingsClarity            = 51
	cameraSettingsHDRPQ              = 52
)

// parseCanonCameraSettings parses tag 0x0001 (CanonCameraSettings).
func (r *Reader) parseCanonCameraSettings(t tag.Entry) canon.CameraSettings {
	var raw [53]uint16
	settings, ok := r.parseCanonSizedUint16Payload(t, "parseCanonCameraSettings", raw[:])
	if !ok {
		return canon.CameraSettings{}
	}

	return canon.CameraSettings{
		MacroMode:          canon.MacroMode(settings.u16(cameraSettingsMacroMode)),
		SelfTimer:          settings.i16(cameraSettingsSelfTimer),
		Quality:            canon.Quality(settings.i16(cameraSettingsQuality)),
		CanonFlashMode:     canon.CanonFlashMode(settings.i16(cameraSettingsFlashMode)),
		ContinuousDrive:    canon.ContinuousDrive(settings.i16(cameraSettingsContinuousDrive)),
		FocusMode:          canon.FocusMode(settings.i16(cameraSettingsFocusMode)),
		RecordMode:         canon.RecordMode(settings.i16(cameraSettingsRecordMode)),
		CanonImageSize:     canon.CanonImageSize(settings.i16(cameraSettingsImageSize)),
		EasyMode:           canon.EasyMode(settings.i16(cameraSettingsEasyMode)),
		DigitalZoom:        canon.DigitalZoom(settings.i16(cameraSettingsDigitalZoom)),
		Contrast:           canon.CameraSettingValue(settings.i16(cameraSettingsContrast)),
		Saturation:         canon.CameraSettingValue(settings.i16(cameraSettingsSaturation)),
		Sharpness:          canon.CameraSettingValue(settings.i16(cameraSettingsSharpness)),
		CameraISO:          canonCameraSettingISO(settings.i16(cameraSettingsCameraISO)),
		MeteringMode:       canon.MeteringMode(settings.i16(cameraSettingsMeteringMode)),
		FocusRange:         canon.FocusRange(settings.i16(cameraSettingsFocusRange)),
		AFPoint:            settings.u16(cameraSettingsAFPoint),
		CanonExposureMode:  canon.ExposureMode(settings.i16(cameraSettingsExposureMode)),
		LensType:           canon.CanonLensType(settings.u16(cameraSettingsLensType)),
		MaxFocalLength:     settings.u16(cameraSettingsMaxFocalLength),
		MinFocalLength:     settings.u16(cameraSettingsMinFocalLength),
		FocalUnits:         settings.u16(cameraSettingsFocalUnits),
		MaxAperture:        canon.MaxApertureFromCode(settings.u16(cameraSettingsMaxAperture)),
		MinAperture:        canon.MaxApertureFromCode(settings.u16(cameraSettingsMinAperture)),
		FlashModel:         canon.FlashModel(settings.i16(cameraSettingsFlashModel)),
		FlashBits:          settings.u16(cameraSettingsFlashBits),
		FocusContinuous:    canon.FocusContinuous(settings.i16(cameraSettingsFocusContinuous)),
		AESetting:          canon.AESetting(settings.i16(cameraSettingsAESetting)),
		ImageStabilization: canon.ImageStabilization(settings.i16(cameraSettingsImageStabilization)),
		DisplayAperture:    canon.DisplayApertureFromCode(settings.u16(cameraSettingsDisplayAperture)),
		ZoomSourceWidth:    settings.u16(cameraSettingsZoomSourceWidth),
		ZoomTargetWidth:    settings.u16(cameraSettingsZoomTargetWidth),
		SpotMeteringMode:   canon.SpotMeteringMode(settings.i16(cameraSettingsSpotMeteringMode)),
		PhotoEffect:        canon.PhotoEffect(settings.i16(cameraSettingsPhotoEffect)),
		ManualFlashOutput:  canon.ManualFlashOutput(settings.i16(cameraSettingsManualFlashOutput)),
		ColorTone:          canon.CameraSettingValue(settings.i16(cameraSettingsColorTone)),
		SRAWQuality:        canon.SRAWQuality(settings.i16(cameraSettingsSRAWQuality)),
		FocusBracketing:    canon.FocusBracketing(settings.i16(cameraSettingsFocusBracketing)),
		Clarity:            canon.ClarityValue(settings.i16(cameraSettingsClarity)),
		HDRPQ:              canon.HDRPQ(settings.u16(cameraSettingsHDRPQ)),
	}
}

const (
	shotInfoAutoISO              = 1
	shotInfoBaseISO              = 2
	shotInfoMeasuredEV           = 3
	shotInfoTargetAperture       = 4
	shotInfoTargetExposureTime   = 5
	shotInfoExposureCompensation = 6
	shotInfoWhiteBalance         = 7
	shotInfoSlowShutter          = 8
	shotInfoSequenceNumber       = 9
	shotInfoOpticalZoomCode      = 10
	shotInfoCameraTemperature    = 12
	shotInfoFlashGuideNumber     = 13
	shotInfoAFPointsInFocus      = 14
	shotInfoFlashExposureComp    = 15
	shotInfoAEB                  = 16
	shotInfoAEBBracketValue      = 17
	shotInfoControlMode          = 18
	shotInfoFocusDistanceUpper   = 19
	shotInfoFocusDistanceLower   = 20
	shotInfoFNumber              = 21
	shotInfoExposureTime         = 22
	shotInfoMeasuredEV2          = 23
	shotInfoBulbDuration         = 24
	shotInfoCameraType           = 26
	shotInfoAutoRotate           = 27
	shotInfoNDFilter             = 28
	shotInfoSelfTimer2           = 29
	shotInfoFlashOutput          = 33
)

// parseCanonShotInfo parses tag 0x0004 (CanonShotInfo).
func (r *Reader) parseCanonShotInfo(t tag.Entry) canon.ShotInfo {
	var raw [64]uint16
	settings, ok := r.parseCanonSizedUint16Payload(t, "parseCanonShotInfo", raw[:])
	if !ok {
		return canon.ShotInfo{}
	}
	autoISO := settings.i16(shotInfoAutoISO)
	baseISO := settings.i16(shotInfoBaseISO)
	measuredEV := settings.i16(shotInfoMeasuredEV)
	targetAperture := settings.i16(shotInfoTargetAperture)
	targetExposureTime := settings.i16(shotInfoTargetExposureTime)
	exposureCompensation := settings.i16(shotInfoExposureCompensation)
	cameraTemperature := settings.i16(shotInfoCameraTemperature)
	flashGuideNumber := settings.i16(shotInfoFlashGuideNumber)
	fNumber := settings.i16(shotInfoFNumber)
	exposureTime := settings.i16(shotInfoExposureTime)
	measuredEV2 := settings.i16(shotInfoMeasuredEV2)
	modelID := r.canonModelID()

	dst := canon.ShotInfo{
		AutoISO:                canon.ShotISO(autoISO),
		BaseISO:                canon.ShotISO(baseISO),
		MeasuredEV:             canon.ShotMeasuredEV(measuredEV),
		TargetAperture:         canon.ShotAperture(targetAperture),
		TargetExposureTime:     canon.ShotExposureTime(targetExposureTime, false),
		ExposureCompensation:   canon.ShotExposureCompensation(exposureCompensation),
		WhiteBalance:           canon.WhiteBalance(settings.i16(shotInfoWhiteBalance)),
		SlowShutter:            canon.SlowShutter(settings.i16(shotInfoSlowShutter)),
		SequenceNumber:         settings.i16(shotInfoSequenceNumber),
		OpticalZoomCode:        settings.i16(shotInfoOpticalZoomCode),
		CameraTemperature:      canonShotCameraTemperature(cameraTemperature, modelID),
		FlashGuideNumber:       canon.ShotFlashGuideNumber(flashGuideNumber),
		AFPointsInFocus:        settings.u16(shotInfoAFPointsInFocus),
		FlashExposureComp:      settings.i16(shotInfoFlashExposureComp),
		AutoExposureBracketing: settings.i16(shotInfoAEB),
		AEBBracketValue:        settings.i16(shotInfoAEBBracketValue),
		ControlMode:            settings.i16(shotInfoControlMode),
		FNumber:                canon.ShotAperture(fNumber),
		ExposureTime:           canon.ShotExposureTime(exposureTime, r.canonShotInfoLegacyExposureTime()),
		MeasuredEV2:            canon.ShotMeasuredEV2(measuredEV2),
		BulbDuration:           settings.i16(shotInfoBulbDuration),
		CameraType:             canon.CameraType(settings.i16(shotInfoCameraType)),
		AutoRotate:             canon.AutoRotate(settings.i16(shotInfoAutoRotate)),
		NDFilter:               canon.NDFilter(settings.i16(shotInfoNDFilter)),
		SelfTimer2:             settings.i16(shotInfoSelfTimer2),
		FlashOutput:            settings.i16(shotInfoFlashOutput),
	}
	dst.ActualISO = canon.ShotActualISO(dst.AutoISO, dst.BaseISO)
	focusUpper := settings.u16(shotInfoFocusDistanceUpper)
	if settings.present(shotInfoFocusDistanceLower) && focusUpper != 0 {
		dst.FocusDistance = canon.NewFocusDistance(focusUpper, settings.u16(shotInfoFocusDistanceLower))
	}
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
	dst := canon.FileInfo{
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
	if canonModelUsesLegacyShutterCount(r.canonModelID()) {
		dst.ShutterCount = uint32(canonU16At(raw[:], n, 2))
		dst.FileNumber = 0
	}
	return dst
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

// parseCanonBatteryType parses Canon Camera:BatteryType (tag 0x0038) like ExifTool.
func (r *Reader) parseCanonBatteryType(t tag.Entry) string {
	raw, _, err := r.readTagBytes(t, canon.BatteryTypePayloadSize)
	if err != nil || len(raw) < canon.BatteryTypePayloadSize {
		r.warnCanonShortRead(t, "parseCanonBatteryType", len(raw), canon.BatteryTypePayloadSize)
		return ""
	}
	return canon.ParseBatteryType(raw[canon.BatteryTypeHeaderLen:])
}

func (r *Reader) parseCanonAFPointsInFocus1D(t tag.Entry, current canon.AFInfo) canon.AFInfo {
	var raw [4]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 1 {
		r.warnCanonShortRead(t, "parseCanonAFPointsInFocus1D", n, 1)
		return current
	}
	if current.Source == canon.AFInfoSourceUnknown {
		current.Source = canon.AFInfoSourceAFInfo
	}
	current.AFPointsInFocusBits = canonBitsetWords(raw[:n])
	return current
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
	return canon.DecodeAFInfo(words[:n], canonModelIsEOS(r.canonModelID()), int(t.UnitCount))
}

func fillCanonAFInfo(dst *canon.AFInfo, words []uint16, modelID canon.CanonCameraModel, afInfoCount int) {
	*dst = canon.DecodeAFInfo(words, canonModelIsEOS(modelID), afInfoCount)
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
	modelID := r.canonModelID()
	isAFInfo3 := canon.MakerNoteTag(t.ID) == canon.AFInfo3
	return canon.DecodeAFInfo2(words[:n], canon.AFInfo2DecodeConfig{
		Source:         source,
		EOS:            canonModelIsEOS(modelID),
		AFInfo3:        isAFInfo3,
		DecodeCoords:   r.afInfoDecodeOptions.has(AFInfoDecodeCoords),
		DecodePoints:   r.afInfoDecodeOptions.has(AFInfoDecodePoints),
		DecodeInFocus:  r.afInfoDecodeOptions.has(AFInfoDecodeInFocus),
		DecodeSelected: r.afInfoDecodeOptions.has(AFInfoDecodeSelected),
	})
}

func (r *Reader) parseCanonString(t tag.Entry) string {
	raw := r.parseOpaqueBytes(t, min(t.Size(), 512))
	if len(raw) == 0 {
		return ""
	}
	return canonString(raw)
}

func canonString(raw []byte) string {
	raw = trimAtNUL(raw)
	if len(raw) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(raw))
	for _, ch := range raw {
		if ch >= 0x20 && ch <= 0x7e {
			b.WriteByte(ch)
		} else {
			b.WriteByte('.')
		}
	}
	s := b.String()
	s = strings.TrimSpace(s)
	return strings.Trim(s, ".")
}

func canonBitsetWords(vals []uint16) []int {
	n := 0
	for _, w := range vals {
		n += bits.OnesCount16(w)
	}
	if n == 0 {
		return nil
	}
	out := make([]int, 0, n)
	base := 0
	for _, word := range vals {
		for word != 0 {
			bit := bits.TrailingZeros16(word)
			out = append(out, base+bit)
			word &= word - 1
		}
		base += 16
	}
	return out
}

func (r *Reader) warnCanonTruncatedWords(t tag.Entry, parser string, got, want int) {
	if !r.WarnEnabled() {
		return
	}
	r.tagLogContext(r.Warn(3), t).
		Str("parser", parser).
		Int("wordsDecoded", got).
		Int("wordsRequested", want).
		Msg("canon AF payload truncated to parser word cap")
}

func (r *Reader) warnCanonShortRead(t tag.Entry, parser string, got, want int) {
	if !r.WarnEnabled() {
		return
	}
	r.tagLogContext(r.Warn(3), t).
		Str("parser", parser).
		Int("gotUnits", got).
		Int("wantUnits", want).
		Msg("canon maker-note payload too short")
}

func (r *Reader) warnCanonInvalidSize(t tag.Entry, parser string, declaredSizeBytes uint32) {
	if !r.WarnEnabled() {
		return
	}
	r.tagLogContext(r.Warn(3), t).
		Str("parser", parser).
		Uint32("declaredSizeBytes", declaredSizeBytes).
		Uint32("actualSizeBytes", t.Size()).
		Msg("invalid canon maker-note payload length")
}
