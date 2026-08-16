package exif

import (
	"io"
	"slices"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote/sony"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

func (r *Reader) sonyMakerNoteHeaderLength(parent tag.Entry) int {
	header, ok := r.peekMakerNotePrefix(sonyMakerNoteHeaderLength)
	if !ok {
		return 0
	}
	if sonyMakerNoteStartsWithDirectory(parent, header) {
		return 0
	}
	return sonyMakerNoteHeaderLength
}

func sonyMakerNoteStartsWithDirectory(parent tag.Entry, header []byte) bool {
	if len(header) < 2 {
		return false
	}
	tagCount := int(parent.ByteOrder.Uint16(header[:2]))
	if tagCount <= 0 || tagCount > maxTagCount {
		return false
	}
	dirSize := 2 + tagCount*12
	parentSize := int(parent.Size())
	return parentSize == 0 || dirSize <= parentSize
}

func (r *Reader) readSonyMakerNoteDirectory(child tag.Directory, headerLen int) error {
	if err := r.discard(headerLen); err != nil {
		return err
	}
	tagCount, err := r.readUint16(child)
	if err != nil {
		return err
	}
	if tagCount > maxTagCount {
		return nil
	}

	total := int(tagCount) * 12
	if total < 0 || total > len(r.state.buf) {
		return imagetype.ErrDataLength
	}
	raw, err := r.fastRead(total)
	if err != nil {
		return err
	}
	if len(raw) < total {
		return io.ErrUnexpectedEOF
	}

	var entryStack [maxTagCount]tag.Entry
	entries := entryStack[:0]
	for pos := 0; pos < total; pos += 12 {
		t, err := tagFromBuffer(child, raw[pos:pos+12])
		if err != nil {
			continue
		}
		if t.IsEmbedded() {
			r.parseSonyTag(t)
		} else {
			entries = append(entries, t)
		}
	}
	if len(entries) > 0 {
		slices.SortFunc(entries, compareSonyTagEntry)
		for i := range entries {
			r.parseSonyTag(entries[i])
		}
	}
	r.promoteSonyDerivedFields()
	return nil
}

func compareSonyTagEntry(a, b tag.Entry) int {
	if a.ValueOffset < b.ValueOffset {
		return -1
	}
	if a.ValueOffset > b.ValueOffset {
		return 1
	}
	if a.ID < b.ID {
		return -1
	}
	if a.ID > b.ID {
		return 1
	}
	return 0
}

func (r *Reader) parseSonyTag(t tag.Entry) bool {
	if r.Exif.MakerNote.Sony == nil {
		r.Exif.MakerNote.Sony = &sony.Sony{}
	}
	dst := r.Exif.MakerNote.Sony
	switch t.ID {
	case sony.Rating:
		dst.Rating = r.parseMakerNoteUint32(t)
	case sony.Contrast:
		dst.Contrast = r.parseSonyInt32(t)
	case sony.Saturation:
		dst.Saturation = r.parseSonyInt32(t)
	case sony.Sharpness:
		dst.Sharpness = r.parseSonyInt32(t)
	case sony.CreativeStyle:
		dst.CreativeStyle = r.parseSonyText(t)
	case sony.DynamicRangeOptimizer:
		dst.DynamicRangeOptimizer = r.parseMakerNoteUint32(t)
	case sony.ImageStabilization:
		dst.ImageStabilization = r.parseMakerNoteUint32(t)
	case sony.ColorMode:
		dst.ColorMode = r.parseMakerNoteUint32(t)
	case sony.Quality:
		dst.Quality = r.parseMakerNoteUint32(t)
	case sony.Quality2:
		dst.Quality2 = r.parseSonyU16Pair(t)
	case sony.WhiteBalance:
		dst.WhiteBalance = r.parseMakerNoteUint32(t)
	case sony.WhiteBalanceFineTune:
		dst.WhiteBalanceFineTune = r.parseSonyInt32(t)
	case sony.FlashExposureComp:
		dst.FlashExposureComp = r.parseSonySignedRationalValue(t)
	case sony.Teleconverter:
		dst.Teleconverter = r.parseMakerNoteUint32(t)
	case sony.SonyModelID:
		value := r.parseMakerNoteUint32(t)
		converted, ok := meta.SafecastUint32ToUint16(value)
		if !ok {
			return true
		}
		dst.SonyModelID = converted
	case sony.LensType:
		dst.LensType = r.parseMakerNoteUint32(t)
	case sony.FileFormat:
		r.parseByteList(t, dst.FileFormat[:])
	case sony.ColorTemperature:
		dst.ColorTemperature = r.parseMakerNoteUint32(t)
	case sony.ColorCompensationFilter:
		dst.ColorCompensationFilter = r.parseSonyInt32(t)
	case sony.SceneMode:
		dst.SceneMode = r.parseMakerNoteUint32(t)
	case sony.ZoneMatching:
		dst.ZoneMatching = r.parseMakerNoteUint32(t)
	case sony.ExposureMode:
		dst.ExposureMode = r.parseSonyU16(t)
	case sony.FocusMode0xB042:
		dst.FocusMode0xB042 = r.parseSonyU16(t)
	case sony.AFAreaMode:
		dst.AFAreaMode = r.parseSonyU16(t)
	case sony.AFIlluminator:
		dst.AFIlluminator = r.parseSonyU16(t)
	case sony.JPEGQuality:
		dst.JPEGQuality = r.parseSonyU16(t)
	case sony.FlashLevel:
		dst.FlashLevel = r.parseMakerNoteInt16(t)
	case sony.ReleaseMode:
		dst.ReleaseMode = r.parseSonyU16(t)
	case sony.SequenceNumber:
		dst.SequenceNumber = r.parseSonyU16(t)
	case sony.AntiBlur:
		dst.AntiBlur = r.parseSonyU16(t)
	case sony.AFTracking:
		dst.AFTracking = r.parseSonyU8(t)
	case sony.DynamicRangeOptimizer0xB04F:
		dst.DynamicRangeOptimizer0xB04F = r.parseSonyU16(t)
	case sony.HighISONoiseReduction2:
		dst.HighISONoiseReduction2 = r.parseSonyU16(t)
	case sony.IntelligentAuto:
		dst.IntelligentAuto = r.parseSonyU16(t)
	case sony.WhiteBalance0xB054:
		dst.WhiteBalance0xB054 = r.parseSonyU16(t)
	case sony.FullImageSize:
		r.parseUint32List(t, dst.FullImageSize[:])
	case sony.PreviewImageSize:
		r.parseUint32List(t, dst.PreviewImageSize[:])
	case sony.LensSpec:
		dst.LensSpec = r.parseDisplayString(t, 8)
	case sony.CameraInfo:
		dst.CameraInfo2, dst.CameraInfo3 = r.parseSonyCameraInfo(t)
	case sony.FocusInfo:
		dst.FocusInfo, dst.MoreInfo = r.parseSonyFocusInfo(t)
	case sony.CameraSettings:
		r.parseSonyCameraSettings(t, dst)
	case sony.ShotInfo:
		dst.ShotInfo = r.parseSonyShotInfo(t)
	case sony.Tag9400a:
		dst.Tag9400 = r.parseSonyTag9400(t)
	case sony.Tag9404a:
		dst.Tag9404 = r.parseSonyTag9404(t)
	case sony.Tag9405a:
		dst.Tag9405 = r.parseSonyTag9405(t)
	case sony.Tag9406:
		dst.Tag9406 = r.parseSonyTag9406(t)
	case sony.Tag940a:
		dst.Tag940A = r.parseSonyTag940A(t)
	case sony.Tag940c:
		dst.Tag940C = r.parseSonyTag940C(t)
	case sony.Tag2010a:
		dst.Tag2010 = r.parseSonyTag2010(t)
	case sony.Tag202a:
		dst.Tag202A = r.parseSonyTag202A(t)
	case sony.HiddenInfo:
		dst.HiddenInfo = r.parseSonyHiddenInfo(t)
	case sony.Tag9050a:
		dst.Tag9050 = r.parseSonyTag9050(t)
	case sony.Sony0x9411, sony.Sony0x9416:
		dst.Tag9416 = r.parseSonyTag9416(t)
	case sony.AFInfo:
		dst.AFInfo = r.parseSonyAFInfo(t)
	case sony.Brightness:
		dst.Brightness = r.parseSonyInt32(t)
	case sony.LongExposureNoiseReduction:
		dst.LongExposureNoiseReduction = r.parseMakerNoteUint32(t)
	case sony.HighISONoiseReduction:
		dst.HighISONoiseReduction = r.parseSonyU16(t)
	case sony.HDR:
		dst.HDR = r.parseMakerNoteUint32(t)
	case sony.MultiFrameNoiseReduction:
		dst.MultiFrameNoiseReduction = r.parseMakerNoteUint32(t)
	case sony.PictureEffect:
		dst.PictureEffect = r.parseSonyU16(t)
	case sony.SoftSkinEffect:
		dst.SoftSkinEffect = r.parseMakerNoteUint32(t)
	case sony.VignettingCorrection:
		dst.VignettingCorrection = r.parseMakerNoteUint32(t)
	case sony.LateralChromaticAberration:
		dst.LateralChromaticAberration = r.parseMakerNoteUint32(t)
	case sony.DistortionCorrectionSetting:
		dst.DistortionCorrectionSetting = r.parseMakerNoteUint32(t)
	case sony.AutoPortraitFramed:
		dst.AutoPortraitFramed = r.parseSonyU16(t)
	case sony.FlashAction:
		dst.FlashAction = r.parseMakerNoteUint32(t)
	case sony.ElectronicFrontCurtainShutter:
		dst.ElectronicFrontCurtainShutter = r.parseMakerNoteUint32(t)
	case sony.FocusMode:
		dst.FocusMode = r.parseSonyU8(t)
	case sony.AFAreaModeSetting:
		dst.AFAreaModeSetting = r.parseSonyU8(t)
	case sony.AFPointSelected:
		dst.AFPointSelected = r.parseSonyU8(t)
	case sony.MultiFrameNREffect:
		dst.MultiFrameNREffect = r.parseMakerNoteUint32(t)
	case sony.RAWFileType:
		dst.RAWFileType = r.parseSonyU16(t)
	case sony.PrioritySetInAWB:
		dst.PrioritySetInAWB = r.parseSonyU8(t)
	case sony.MeteringMode2:
		dst.MeteringMode2 = r.parseSonyU16(t)
	case sony.Macro:
		dst.Macro = r.parseSonyU16(t)
	case sony.FocusMode0xB04E:
		dst.FocusMode0xB04E = r.parseSonyU16(t)
	case sony.WBShiftABGM:
		r.parseInt32List(t, dst.WBShiftABGM[:])
	case sony.FlexibleSpotPosition:
		r.parseUint16List(t, dst.FlexibleSpotPosition[:])
	case sony.WBShiftABGMPrecise:
		r.parseInt32List(t, dst.WBShiftABGMPrecise[:])
	case sony.FocusLocation:
		r.parseUint16List(t, dst.FocusLocation[:])
	case sony.VariableLowPassFilter:
		r.parseUint16List(t, dst.VariableLowPassFilter[:])
	case sony.ExposureStandardAdjustment:
		dst.ExposureStandardAdjustment = r.parseSonySignedRationalValue(t)
	case sony.SerialNumber:
		dst.SerialNumber = r.parseSonyText(t)
	case sony.Shadows:
		dst.Shadows = r.parseSonyInt32(t)
	case sony.Highlights:
		dst.Highlights = r.parseSonyInt32(t)
	case sony.Fade:
		dst.Fade = r.parseSonyInt32(t)
	case sony.SharpnessRange:
		dst.SharpnessRange = r.parseSonyInt32(t)
	case sony.Clarity:
		dst.Clarity = r.parseSonyInt32(t)
	case sony.FocusLocation2:
		r.parseUint16List(t, dst.FocusLocation2[:])
	default:
		return false
	}
	return true
}

func (r *Reader) parseSonyCameraInfo(t tag.Entry) (legacy sony.SonyCameraInfo2, modern sony.SonyCameraInfo3) {
	raw, ok := r.readSonyTagBytes(t, 0x80)
	if !ok {
		return legacy, modern
	}
	if sony.UsesCameraInfo3(r.Exif.IFD0.Model) && t.UnitCount > 6000 || t.UnitCount >= 15360 {
		modern = sony.ParseCameraInfo3(raw, t.ByteOrder)
	} else {
		legacy = sony.ParseCameraInfo2(raw, t.ByteOrder)
	}
	return legacy, modern
}

func (r *Reader) parseSonyFocusInfo(t tag.Entry) (legacy sony.SonyFocusInfo, modern sony.SonyMoreInfo) {
	raw, ok := r.readSonyTagBytes(t, 0x0a00)
	if !ok {
		return legacy, modern
	}
	switch t.UnitCount {
	case 19154, 19148:
		legacy = sony.ParseFocusInfo(raw, t.ByteOrder)
	case 20480:
		modern = sony.ParseMoreInfo(raw, t.ByteOrder)
	}
	return legacy, modern
}

func (r *Reader) parseSonyCameraSettings(t tag.Entry, dst *sony.Sony) {
	raw, ok := r.readSonyTagBytes(t, 0x0204)
	if !ok {
		return
	}
	switch t.UnitCount {
	case 280, 332:
		dst.CameraSettings = sony.ParseCameraSettings(raw, t.ByteOrder)
	case 1536, 2048:
		dst.CameraSettings3 = sony.ParseCameraSettings3(raw, t.ByteOrder)
	}
}

func (r *Reader) parseSonyShotInfo(t tag.Entry) sony.SonyShotInfo {
	raw, ok := r.readSonyTagBytes(t, 0x40)
	if !ok {
		return sony.SonyShotInfo{}
	}
	return sony.ParseShotInfo(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag9400(t tag.Entry) sony.SonyTag9400 {
	raw, ok := r.readSonyTagBytes(t, 0x60)
	if !ok {
		return sony.SonyTag9400{}
	}
	return sony.ParseTag9400(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag9404(t tag.Entry) sony.SonyTag9404 {
	raw, ok := r.readSonyTagBytes(t, 0x20)
	if !ok {
		return sony.SonyTag9404{}
	}
	return sony.ParseTag9404(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag9405(t tag.Entry) sony.SonyTag9405 {
	raw, ok := r.readSonyTagBytes(t, 0x076c)
	if !ok {
		return sony.SonyTag9405{}
	}
	return sony.ParseTag9405(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag9406(t tag.Entry) sony.SonyTag9406 {
	raw, ok := r.readSonyTagBytes(t, 0x10)
	if !ok {
		return sony.SonyTag9406{}
	}
	return sony.ParseTag9406(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag940A(t tag.Entry) sony.SonyTag940A {
	raw, ok := r.readSonyTagBytes(t, 0x20)
	if !ok {
		return sony.SonyTag940A{}
	}
	return sony.ParseTag940A(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag940C(t tag.Entry) sony.SonyTag940C {
	raw, ok := r.readSonyTagBytes(t, 0x20)
	if !ok {
		return sony.SonyTag940C{}
	}
	return sony.ParseTag940C(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag2010(t tag.Entry) sony.SonyTag2010 {
	raw, ok := r.readSonyTagBytes(t, 0x1220)
	if !ok {
		return sony.SonyTag2010{}
	}
	return sony.ParseTag2010(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag202A(t tag.Entry) sony.SonyTag202A {
	raw, ok := r.readSonyTagBytes(t, 0x10)
	if !ok {
		return sony.SonyTag202A{}
	}
	return sony.ParseTag202A(raw, t.ByteOrder)
}

func (r *Reader) parseSonyHiddenInfo(t tag.Entry) sony.SonyHiddenInfo {
	raw, ok := r.readSonyTagBytes(t, 0x08)
	if !ok {
		return sony.SonyHiddenInfo{}
	}
	return sony.ParseHiddenInfo(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag9050(t tag.Entry) sony.SonyTag9050 {
	raw, ok := r.readSonyTagBytes(t, 0x0200)
	if !ok {
		return sony.SonyTag9050{}
	}
	return sony.ParseTag9050(raw, t.ByteOrder)
}

func (r *Reader) parseSonyTag9416(t tag.Entry) sony.SonyTag9416 {
	raw, ok := r.readSonyTagBytes(t, 0x0a00)
	if !ok {
		return sony.SonyTag9416{}
	}
	return sony.ParseTag9416(raw, t.ByteOrder)
}

func (r *Reader) parseSonyAFInfo(t tag.Entry) sony.SonyAFInfo {
	raw, ok := r.readSonyTagBytes(t, 0x0180)
	if !ok {
		return sony.SonyAFInfo{}
	}
	return sony.ParseAFInfo(raw, t.ByteOrder)
}

func (r *Reader) parseSonyText(t tag.Entry) string {
	buf, err := r.readTagBytes(t, t.Size())
	if err != nil || len(buf) == 0 {
		return ""
	}
	return sony.DisplayText(buf)
}

func (r *Reader) parseSonyU16Pair(t tag.Entry) [2]uint16 {
	var dst [2]uint16
	r.parseUint16List(t, dst[:])
	return dst
}

func (r *Reader) parseSonyInt32(t tag.Entry) int32 {
	switch t.Type {
	case tag.TypeSignedLong:
		var dst [1]int32
		if n := r.parseInt32List(t, dst[:]); n > 0 {
			return dst[0]
		}
	case tag.TypeSignedShort, tag.TypeShort:
		return int32(r.parseMakerNoteInt16(t))
	case tag.TypeLong:
		value := r.parseMakerNoteUint32(t)
		converted, ok := meta.SafecastUint32ToInt32(value)
		if !ok {
			return 0
		}
		return converted
	}
	return 0
}

func (r *Reader) parseSonySignedRationalValue(t tag.Entry) float64 {
	var raw [2]int32
	if r.parseRationalSList(t, raw[:]) == 0 || raw[1] == 0 {
		return 0
	}
	return float64(raw[0]) / float64(raw[1])
}

func (r *Reader) parseSonyU16(t tag.Entry) uint16 {
	var raw [1]uint16
	if r.parseUint16List(t, raw[:]) == 0 {
		return 0
	}
	return raw[0]
}

func (r *Reader) parseSonyU8(t tag.Entry) uint8 {
	var raw [1]byte
	if r.parseByteList(t, raw[:]) == 0 {
		return 0
	}
	return raw[0]
}

func (r *Reader) readSonyTagBytes(t tag.Entry, maxBytes uint32) ([]byte, bool) {
	size := t.Size()
	if size == 0 {
		return nil, true
	}
	if maxBytes == 0 || maxBytes > size {
		maxBytes = size
	}
	if err := r.seekToTag(t); err != nil {
		return nil, false
	}
	if maxBytes > uint32(len(r.state.buf)) {
		maxBytes = uint32(len(r.state.buf))
	}
	buf, err := r.fastRead(int(maxBytes))
	if err != nil {
		return nil, false
	}
	remaining := int(size) - len(buf)
	if remaining > 0 {
		if err := r.discard(remaining); err != nil {
			return nil, false
		}
	}
	return buf, true
}

func (r *Reader) promoteSonyDerivedFields() {
	dst := r.Exif.MakerNote.Sony
	if dst == nil {
		return
	}

	if !isAsciiCreativeStyle(dst.CreativeStyle) {
		if v := sony.CreativeStyleValue(dst.CameraSettings3.CreativeStyleSetting).Name(); v != "" {
			dst.CreativeStyle = v
		} else if v := sony.CreativeStyleValue(dst.CameraSettings.CreativeStyle).Name(); v != "" {
			dst.CreativeStyle = v
		} else if v := sony.CreativeStyleValue(dst.Tag9416.CreativeStyle).Name(); v != "" {
			dst.CreativeStyle = v
		}
	}

	if dst.Quality == 0 {
		if v := dst.CameraSettings3.Quality; v != 0 {
			dst.Quality = uint32(v)
		} else if v := dst.CameraSettings.Quality; v != 0 {
			dst.Quality = uint32(v)
		} else if v := dst.Tag9400.Quality2; v != 0 {
			dst.Quality = uint32(v)
		}
	}

	if v := dst.Tag9050.LensType; v != 0 && v != 65535 {
		dst.LensType = uint32(v)
	} else if v := dst.Tag9416.LensType2; v != 0 && v != 65535 {
		dst.LensType = uint32(v)
	} else if v := dst.Tag940C.LensType3; v != 0 && v != 65535 {
		dst.LensType = uint32(v)
	}

	if dst.SonyModelID == 0 {
		dst.SonyModelID = sony.ModelIDFromModel(r.Exif.IFD0.Model)
	}
}

// isAsciiCreativeStyle reports whether s is a valid human-readable CreativeStyle
// string from tag 0xb020.  Returns true only when the raw string matches a known
// name; false for binary garbage, dots, or empty.
func isAsciiCreativeStyle(s string) bool {
	switch s {
	case "Standard", "Vivid", "Portrait", "Landscape", "Sunset",
		"Night View/Portrait", "B&W", "Adobe RGB",
		"Neutral", "Clear", "Deep", "Light",
		"Autumn Leaves", "Off", "Sepia":
		return true
	}
	return false
}
