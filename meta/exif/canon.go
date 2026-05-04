package exif

import (
	"math/bits"

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
		if canon.ShouldReplaceAFInfo(dst.AFInfo, candidate) {
			dst.AFInfo = candidate
		}
	case canon.AFPointsInFocus1D:
		dst.AFInfo = r.parseCanonAFPointsInFocus1D(t, dst.AFInfo)
	case canon.CanonAFInfo2, canon.AFInfo3:
		candidate := r.parseCanonAFInfo2(t)
		if canon.ShouldReplaceAFInfo(dst.AFInfo, candidate) {
			dst.AFInfo = candidate
		}
	case canon.FaceDetect1:
		dst.FaceDetect = r.parseCanonFaceDetect1(t)
	case canon.FaceDetect2:
		dst.FaceDetect = r.parseCanonFaceDetect2(t)
	case canon.FaceDetect3:
		dst.FaceDetect = r.parseCanonFaceDetect3(t)
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
	if spec, ok := canonCameraInfoSpecForLayout(layout); ok {
		return r.parseCanonCameraInfoBytes(t, spec)
	}
	switch layout {
	case canon.CameraInfoLayoutPowerShot, canon.CameraInfoLayoutPowerShot2:
		return r.parseCanonCameraInfoPowerShotTemp(t)
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

func canonCameraInfoSpecForLayout(layout canon.CameraInfoLayout) (canon.CameraInfoSpec, bool) {
	switch layout {
	case canon.CameraInfoLayout5D:
		return canon.CameraInfoSpecLayout5D, true
	case canon.CameraInfoLayout5DmkII:
		return canon.CameraInfoSpecLayout5DmkII, true
	case canon.CameraInfoLayout5DmkIII:
		return canon.CameraInfoSpecLayout5DmkIII, true
	case canon.CameraInfoLayout6D:
		return canon.CameraInfoSpecLayout6D, true
	case canon.CameraInfoLayout7D:
		return canon.CameraInfoSpecLayout7D, true
	case canon.CameraInfoLayout40D:
		return canon.CameraInfoSpecLayout40D, true
	case canon.CameraInfoLayout50D:
		return canon.CameraInfoSpecLayout50D, true
	case canon.CameraInfoLayout60D:
		return canon.CameraInfoSpecLayout60D, true
	case canon.CameraInfoLayout70D:
		return canon.CameraInfoSpecLayout70D, true
	case canon.CameraInfoLayout80D:
		return canon.CameraInfoSpecLayout80D, true
	case canon.CameraInfoLayout450D:
		return canon.CameraInfoSpecLayout450D, true
	case canon.CameraInfoLayout500D:
		return canon.CameraInfoSpecLayout500D, true
	case canon.CameraInfoLayout550D:
		return canon.CameraInfoSpecLayout550D, true
	case canon.CameraInfoLayout600D:
		return canon.CameraInfoSpecLayout600D, true
	case canon.CameraInfoLayout650D:
		return canon.CameraInfoSpecLayout650D, true
	case canon.CameraInfoLayout700D:
		return canon.CameraInfoSpecLayout700D, true
	case canon.CameraInfoLayout750D:
		return canon.CameraInfoSpecLayout750D, true
	case canon.CameraInfoLayout1000D:
		return canon.CameraInfoSpecLayout1000D, true
	case canon.CameraInfoLayout1D:
		return canon.CameraInfoSpecLayout1D, true
	case canon.CameraInfoLayout1DmkII:
		return canon.CameraInfoSpecLayout1DmkII, true
	case canon.CameraInfoLayout1DmkIIN:
		return canon.CameraInfoSpecLayout1DmkIIN, true
	case canon.CameraInfoLayout1DmkIII:
		return canon.CameraInfoSpecLayout1DmkIII, true
	case canon.CameraInfoLayout1DmkIV:
		return canon.CameraInfoSpecLayout1DmkIV, true
	case canon.CameraInfoLayout1DX:
		return canon.CameraInfoSpecLayout1DX, true
	default:
		return canon.CameraInfoSpec{}, false
	}
}

func (r *Reader) canonCameraInfoLayout(t tag.Entry) canon.CameraInfoLayout {
	if layout, ok := canon.CameraInfoLayoutForModelID(r.canonModelID()); ok {
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

// parseCanonCameraInfoPowerShotTemp extracts camera temperature from PowerShot
// int32u payloads where the last word is a known non-temperature field and the
// second-to-last word (at index count-3) is the temperature value.
//
// This matches ExifTool's CameraInfoPowerShot table where element [-3] (relative
// to CameraInfoCount) is CameraTemperature.
func (r *Reader) parseCanonCameraInfoPowerShotTemp(t tag.Entry) canon.CameraInfo {
	count := int(t.UnitCount)
	if count < 3 {
		return canon.CameraInfo{}
	}
	// Temperature is at index count-3 (ExifTool: element [-3]).
	tempOff := (count - 3) * 4
	if v, ok := r.readCanonCameraInfoInt32At(t, tempOff); ok {
		return canon.CameraInfo{CameraTemperature: int16(v)}
	}
	return canon.CameraInfo{}
}

func (r *Reader) parseCanonCameraInfoUnknown32(t tag.Entry) canon.CameraInfo {
	count := int(t.UnitCount)
	var tempOff int
	switch count {
	case 72:
		tempOff = 71 * 4
	case 85:
		tempOff = 83 * 4
	case 93, 94:
		tempOff = 91 * 4
	case 96:
		tempOff = 92 * 4
	case 104:
		tempOff = 100 * 4
	default:
		if count > 400 {
			tempOff = (count - 3) * 4
		} else {
			return canon.CameraInfo{}
		}
	}
	if v, ok := r.readCanonCameraInfoInt32At(t, tempOff); ok {
		return canon.CameraInfo{CameraTemperature: int16(v)}
	}
	return canon.CameraInfo{}
}

func (r *Reader) readCanonCameraInfoUint32At(t tag.Entry, off int) (uint32, bool) {
	b, err := r.readCanonTagOffsetBytes(t, off, 4)
	if err != nil || len(b) < 4 {
		return 0, false
	}
	return canon.U32LEAt(b, 0), true
}

func (r *Reader) readCanonCameraInfoUint16At(t tag.Entry, off int) (uint16, bool) {
	b, err := r.readCanonTagOffsetBytes(t, off, 2)
	if err != nil || len(b) < 2 {
		return 0, false
	}
	return canon.U16LEAt(b, 0), true
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
	if off < 0 || n <= 0 || off+n > size {
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

func (r *Reader) parseCanonUint16List(t tag.Entry, s canon.Seq16) int {
	switch t.Type {
	case tag.TypeShort, tag.TypeSignedShort:
		return r.parseCanonRawUint16List(t, s, int(t.UnitCount))
	case tag.TypeUndefined:
		return r.parseCanonRawUint16List(t, s, int(t.UnitCount/2))
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

// parseCanonSeq16 reads a uint16 payload, validates the size word, strips it,
// and returns a 1-based Seq16 view.
func (r *Reader) parseCanonSeq16(t tag.Entry, dst []uint16, parser string) canon.Seq16 {
	n := r.parseCanonUint16List(t, dst)
	if n < 2 {
		r.warnCanonShortRead(t, parser, n, 1)
		return nil
	}
	if uint32(dst[0]) != t.Size() {
		r.warnCanonInvalidSize(t, parser, uint32(dst[0]))
		return nil
	}
	return canon.Seq16(dst[1:n])
}

// parseCanonInt32List reads int32 values, falling back to uint32 on mismatch.
func (r *Reader) parseCanonInt32List(t tag.Entry, dst []int32) int {
	if n := r.parseInt32List(t, dst); n > 0 {
		return n
	}
	var u32 [2048]uint32
	if len(dst) > len(u32) {
		dst = dst[:len(u32)]
	}
	n := r.parseUint32List(t, u32[:len(dst)])
	for i := range n {
		dst[i] = int32(u32[i])
	}
	return n
}

// parseCanonSeq32 reads an int32 payload, strips the leading size word, and
// returns a 1-based Seq32 view.
func (r *Reader) parseCanonSeq32(t tag.Entry, dst []int32) canon.Seq32 {
	n := r.parseCanonInt32List(t, dst)
	if n < 2 {
		return nil
	}
	return canon.Seq32(dst[1:n])
}

func (r *Reader) parseCanonPreviewImageInfo(t tag.Entry) canon.PreviewImageInfo {
	var raw [8]int32
	n := r.parseCanonInt32List(t, raw[:])
	start := 0
	if n >= 6 && isPreviewImageInfoSizeWord(raw, int32(t.Size())) {
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

// isPreviewImageInfoSizeWord reports whether raw[0] is a size-header word that
// should be skipped when decoding PreviewImageInfo. ExifTool uses this heuristic
// when the first word matches the tag size, or when raw[0] is a small value
// followed by a plausible preview quality byte and an out-of-range width.
func isPreviewImageInfoSizeWord(raw [8]int32, tagSize int32) bool {
	if raw[0] == tagSize {
		return true
	}
	return raw[0] > 5 && raw[1] >= 0 && raw[1] <= 5 && raw[2] > 0xffff
}

func (r *Reader) parseCanonSensorInfo(t tag.Entry) canon.SensorInfo {
	var raw [13]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 13 {
		r.warnCanonShortRead(t, "parseCanonSensorInfo", n, 13)
		return canon.SensorInfo{}
	}
	return canon.Seq16(raw[:]).DecodeSensorInfo()
}

func (r *Reader) parseCanonAFConfig(t tag.Entry) canon.AFConfig {
	var raw [25]int32
	s := r.parseCanonSeq32(t, raw[:])
	if s == nil {
		return canon.AFConfig{}
	}
	return s.DecodeAFConfig()
}

func (r *Reader) parseCanonLightingOpt(t tag.Entry) canon.LightingOptInfo {
	var raw [12]int32
	s := r.parseCanonSeq32(t, raw[:])
	if s == nil {
		return canon.LightingOptInfo{}
	}
	return s.DecodeLightingOpt()
}

func (r *Reader) parseCanonMultiExp(t tag.Entry) canon.MultiExpInfo {
	var raw [8]int32
	s := r.parseCanonSeq32(t, raw[:])
	if s == nil {
		return canon.MultiExpInfo{}
	}
	return s.DecodeMultiExp()
}

func (r *Reader) parseCanonHDRInfo(t tag.Entry) canon.HDRInfo {
	var raw [8]int32
	s := r.parseCanonSeq32(t, raw[:])
	if s == nil {
		return canon.HDRInfo{}
	}
	return s.DecodeHDRInfo()
}

func (r *Reader) parseCanonAFMicroAdj(t tag.Entry) canon.AFMicroAdjInfo {
	var raw [8]int32
	s := r.parseCanonSeq32(t, raw[:])
	if s == nil {
		return canon.AFMicroAdjInfo{}
	}
	return s.DecodeAFMicroAdj()
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

func (r *Reader) parseCanonFaceDetect1(t tag.Entry) canon.FaceDetectInfo {
	var raw [26]uint16
	n := r.parseCanonUint16List(t, raw[:])
	if n < 5 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect1", n, 5)
		return canon.FaceDetectInfo{}
	}
	return canon.DecodeFaceDetect1(raw[:n])
}

func (r *Reader) parseCanonFaceDetect2(t tag.Entry) canon.FaceDetectInfo {
	var raw [8]byte
	if n := r.parseByteList(t, raw[:]); n < 3 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect2", n, 3)
		return canon.FaceDetectInfo{}
	}
	return canon.DecodeFaceDetect2(raw[:])
}

func (r *Reader) parseCanonFaceDetect3(t tag.Entry) canon.FaceDetectInfo {
	var raw [8]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonFaceDetect3", n, 4)
		return canon.FaceDetectInfo{}
	}
	return canon.DecodeFaceDetect3(raw[:])
}

// parseCanonFocalLength parses tag 0x0002 (CanonFocalLength).
func (r *Reader) parseCanonFocalLength(t tag.Entry) canon.FocalLengthInfo {
	var raw [8]uint16
	if n := r.parseCanonUint16List(t, raw[:]); n < 4 {
		r.warnCanonShortRead(t, "parseCanonFocalLength", n, 4)
		return canon.FocalLengthInfo{}
	}
	return canon.Seq16(raw[:]).DecodeFocalLength()
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
	return canon.Seq16(raw[1:n]).DecodeProcessingInfo()
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
	copy(dst.Raw[:], raw)
	dst.RawCount = canonLensInfoByteLength
	if dst.Raw[0] == 0 && dst.Raw[1] == 0 && dst.Raw[2] == 0 && dst.Raw[3] == 0 {
		return dst
	}
	dst.LensSerialNumber = canon.HexBytes(dst.Raw[:])
	return dst
}

// parseCanonCameraSettings parses tag 0x0001 (CanonCameraSettings).
func (r *Reader) parseCanonCameraSettings(t tag.Entry) canon.CameraSettings {
	var raw [53]uint16
	s := r.parseCanonSeq16(t, raw[:], "parseCanonCameraSettings")
	if s == nil {
		return canon.CameraSettings{}
	}
	return s.DecodeCameraSettings()
}

func (r *Reader) parseCanonShotInfo(t tag.Entry) canon.ShotInfo {
	var raw [64]uint16
	s := r.parseCanonSeq16(t, raw[:], "parseCanonShotInfo")
	if s == nil {
		return canon.ShotInfo{}
	}
	return s.DecodeShotInfo(canon.ShotInfoDecodeConfig{
		LegacyExposureTime: r.canonShotInfoLegacyExposureTime(),
		ModelID:            r.canonModelID(),
	})
}

// parseCanonFileInfo parses tag 0x0093 (CanonFileInfo).
func (r *Reader) parseCanonFileInfo(t tag.Entry) canon.FileInfo {
	var raw [64]uint16
	s := r.parseCanonSeq16(t, raw[:], "parseCanonFileInfo")
	if s == nil {
		return canon.FileInfo{}
	}
	return s.DecodeFileInfo(r.canonModelID())
}

// parseCanonTimeInfo parses tag 0x0035 (TimeInfo).
func (r *Reader) parseCanonTimeInfo(t tag.Entry) canon.CanonTimeInfo {
	var raw [4]int32
	s := r.parseCanonSeq32(t, raw[:])
	if s == nil {
		return canon.CanonTimeInfo{}
	}
	return s.DecodeTimeInfo()
}

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

func (r *Reader) parseCanonAFWords(t tag.Entry, parser string, source tag.ID, dst []uint16) ([]uint16, canon.AFInfoSource, bool) {
	words, truncated := canon.AFWordsBuffer(dst, t.UnitCount)
	if truncated {
		r.warnCanonTruncatedWords(t, parser, len(words), int(t.UnitCount))
	}
	n := r.parseCanonUint16List(t, words)
	src := canon.AFInfoSourceFromID(source)
	if n == 0 {
		r.warnCanonShortRead(t, parser, n, 1)
		return nil, src, false
	}
	return words[:n], src, true
}

// parseCanonAFInfo parses tag 0x0012 (AFInfo).
func (r *Reader) parseCanonAFInfo(t tag.Entry) canon.AFInfo {
	var wordsStack [2048]uint16
	words, source, ok := r.parseCanonAFWords(t, "parseCanonAFInfo", tag.ID(canon.CanonAFInfo), wordsStack[:])
	if !ok {
		return canon.AFInfo{Source: source}
	}
	return canon.DecodeAFInfo(words, canon.ModelIsEOS(r.canonModelID()), int(t.UnitCount))
}

// parseCanonAFInfo2 parses tags 0x0026 and 0x003c (AFInfo2/AFInfo3).
func (r *Reader) parseCanonAFInfo2(t tag.Entry) canon.AFInfo {
	var wordsStack [2048]uint16
	words, source, ok := r.parseCanonAFWords(t, "parseCanonAFInfo2", t.ID, wordsStack[:])
	if !ok {
		return canon.AFInfo{Source: source}
	}
	modelID := r.canonModelID()
	isAFInfo3 := canon.MakerNoteTag(t.ID) == canon.AFInfo3
	return canon.DecodeAFInfo2(words, canon.AFInfo2DecodeConfig{
		Source:         source,
		EOS:            canon.ModelIsEOS(modelID),
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
	var stack [512]byte
	buf := stack[:0]
	if len(raw) > len(stack) {
		buf = make([]byte, len(raw))
	} else {
		buf = stack[:len(raw)]
	}
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if ch >= 0x20 && ch <= 0x7e {
			buf[i] = ch
		} else {
			buf[i] = '.'
		}
	}

	start := 0
	end := len(buf)
	for start < end && buf[start] == ' ' {
		start++
	}
	for end > start && buf[end-1] == ' ' {
		end--
	}
	for start < end && buf[start] == '.' {
		start++
	}
	for end > start && buf[end-1] == '.' {
		end--
	}
	if start >= end {
		return ""
	}
	return string(buf[start:end])
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
