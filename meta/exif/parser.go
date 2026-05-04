package exif

import (
	"encoding/binary"
	"math"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/evanoberholster/imagemeta/meta/exif/makernote"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// tagFromBuffer decodes a tag entry from a raw TIFF directory buffer.
func tagFromBuffer(directory tag.Directory, buf []byte) (tag.Entry, error) {
	tagID := tag.ID(directory.ByteOrder.Uint16(buf[:2]))
	tagType := tag.Type(directory.ByteOrder.Uint16(buf[2:4]))
	unitCount := directory.ByteOrder.Uint32(buf[4:8])
	valueOffset := directory.ByteOrder.Uint32(buf[8:12])
	tagType = tag.NormalizeType(directory.Type, tagID, tagType)
	entry := tag.NewEntry(tagID, tagType, unitCount, valueOffset, directory.Type, directory.Index, directory.ByteOrder)

	if !entry.IsEmbedded() {
		entry.ValueOffset += directory.BaseOffset
	}
	if !tagType.IsValid() {
		return entry, tag.ErrTagTypeNotValid
	}
	return entry, nil
}

// addTag appends a tag to parser state while preserving parse order constraints.
func (r *Reader) addTag(t tag.Entry) {
	if t.IfdType != tag.MakerNoteIFD && t.ValueOffset < r.po && t.ID != tag.TagMakerNote {
		if r.WarnEnabled() {
			r.warnTagAdd(t, "ignoring reverse-offset tag")
		}
		return
	}
	if !r.state.addTag(t) && r.WarnEnabled() {
		r.warnTagQueueFull(t)
	}
}

// warnTagReadError logs tag payload read failures from value parsers.
func (r *Reader) warnTagReadError(t tag.Entry, err error, msg string) {
	if r.WarnEnabled() {
		r.tagLogContext(r.Warn(3), t).
			Err(err).
			Msg(msg)
	}
}

// warnTagAdd logs dropped-tag context while keeping warning-path overhead small.
func (r *Reader) warnTagAdd(t tag.Entry, msg string) {
	r.warnTagContext(t, msg, false)
}

// warnTagQueueFull logs tag-queue saturation details.
func (r *Reader) warnTagQueueFull(t tag.Entry) {
	r.warnTagContext(t, "tag queue full", true)
}

// parseSubIFDs parses the requested value from EXIF metadata.
func (r *Reader) parseSubIFDs(t tag.Entry) {
	switch t.Type {
	case tag.TypeLong, tag.TypeIfd:
	default:
		return
	}
	if r.state.len >= tagQueueMax {
		if r.WarnEnabled() {
			r.warnTagContext(t, "subifd queue capacity reached; skipping parse", true)
		}
		return
	}
	offsetRemaining := len(r.Exif.IFD0.subIFDOffsets) - int(r.Exif.IFD0.subIFDOffsetCount)
	if offsetRemaining <= 0 {
		if r.WarnEnabled() {
			r.warnTagContext(t, "subifd offset capacity reached; skipping parse", false)
		}
		return
	}

	maxEntries := int(t.UnitCount)
	queueRemaining := int(tagQueueMax - r.state.len)
	maxEntries = min(maxEntries, queueRemaining, offsetRemaining)
	if maxEntries <= 0 {
		return
	}

	if t.UnitCount == 1 {
		var offset uint32
		switch {
		case t.IsType(tag.TypeIfd):
			offset = t.ValueOffset
		case t.IsEmbedded():
			t.EmbeddedValue(r.state.buf[:4])
			offset = t.ByteOrder.Uint32(r.state.buf[:4])
		default:
			buf, _, err := r.readTagBytes(t, 4)
			if err != nil || len(buf) < 4 {
				return
			}
			offset = t.ByteOrder.Uint32(buf[:4])
		}
		if offset != 0 {
			r.appendSubIFDOffset(offset)
			r.addTag(tag.NewEntry(t.ID, tag.TypeIfd, 1, offset, tag.SubIFD0, 0, t.ByteOrder))
		}
		return
	}

	buf, _, err := r.readTagBytes(t, uint32(maxEntries*4))
	if err != nil {
		return
	}
	limit := min(maxEntries, len(buf)/4)
	for i := range limit {
		if r.state.len >= tagQueueMax {
			if r.WarnEnabled() {
				r.warnTagContext(t, "subifd queue capacity reached; stopping parse", true)
			}
			break
		}
		offset := t.ByteOrder.Uint32(buf[i*4 : i*4+4])
		if offset == 0 {
			continue
		}
		if int(r.Exif.IFD0.subIFDOffsetCount) >= len(r.Exif.IFD0.subIFDOffsets) {
			break
		}
		r.appendSubIFDOffset(offset)
		subType := tag.SubIFDTypeForIndex(i)
		r.addTag(tag.NewEntry(t.ID, tag.TypeIfd, 1, offset, subType, int8(i), t.ByteOrder))
	}
}

// appendSubIFDOffset stores a parsed SubIFD pointer in the bounded IFD0 list.
func (r *Reader) appendSubIFDOffset(offset uint32) {
	if offset == 0 || int(r.Exif.IFD0.subIFDOffsetCount) >= len(r.Exif.IFD0.subIFDOffsets) {
		return
	}
	r.Exif.IFD0.subIFDOffsets[r.Exif.IFD0.subIFDOffsetCount] = offset
	r.Exif.IFD0.subIFDOffsetCount++
}

// parseTag parses the requested value from EXIF metadata.
func (r *Reader) parseTag(t tag.Entry) {
	if r.DebugEnabled() {
		r.tagLogContext(r.Debug(3), t).Msg("parse tag")
	}
	switch t.IfdType {
	case tag.IFD0:
		if !r.parseIFD0Tag(t) {
			return
		}
	case tag.IFD1:
		if !r.parseImageIFDTag(t, r.Exif.IFD1) {
			return
		}
	case tag.IFD2:
		if !r.parseImageIFDTag(t, r.Exif.IFD2) {
			return
		}
	case tag.GPSIFD:
		if !r.parseGPSTag(t) {
			return
		}
	case tag.MakerNoteIFD:
		if !r.parseMakerNoteTag(t) {
			return
		}
		return
	case tag.ExifIFD:
		if !r.parseExifTag(t) {
			return
		}
	default:
		// SubIFD{0..7} tags are normalized through ExifIFD parsing semantics.
		if t.IfdType != tag.ExifIFD && !t.IfdType.IsSubIFD() {
			return
		}
	}
}

func (r *Reader) parseDNGTags() bool {
	switch r.Exif.ImageType {
	case imagetype.ImageDNG, imagetype.ImageTiff:
		return true
	default:
		return false
	}
}

// parseIFD0Tag parses IFD0 tags into typed model fields.
//
// Non-parsed IFD0 tags are documented in the explicit "intentionally non-parsed"
// case branch below. These tags are treated as handled for coverage/reporting
// parity but are not mapped into the Exif model.
func (r *Reader) parseIFD0Tag(t tag.Entry) bool {
	ifd0 := &r.Exif.IFD0
	exifIfd := &r.Exif.ExifIFD
	switch t.ID {
	// Text/date tags.
	case tag.TagDateTime:
		modifyDate := r.parseDate(t)
		ifd0.setDate(modifyDate)
		ifd0.markTagParsed(t.ID)
	case tag.TagDateTimeOriginal:
		dateTimeOriginal := r.parseDate(t)
		exifIfd.setDate(t.ID, dateTimeOriginal)
		exifIfd.markTagParsed(t.ID)
	case tag.TagMake:
		r.Exif.CameraMakeID, ifd0.Make = r.parseMakeTag(t)
	case tag.TagModel:
		ifd0.Model = r.parseString(t)
	case tag.TagArtist:
		ifd0.Artist = r.parseStringTrimRightSpaceNewline(t)
	case tag.TagCopyright:
		ifd0.Copyright = r.parseDisplayStringTrimRightSpaceNewline(t, 512)
	case tag.TagApplicationNotes:
		// TODO: TagApplicationNotes parsing intentionally disabled for now.
		// The payload is often large and not needed in the hot parse path.
	case tag.TagPrintIM:
		// ignore tag
		//ifd0.PrintIM = r.parseDisplayString(t, 512)
	case tag.TagImageDescription:
		ifd0.ImageDescription = r.parseString(t)
	case tag.TagSoftware:
		ifd0.Software = r.parseString(t)
	// Image/layout tags.
	case tag.TagSubfileType:
		ifd0.SubfileType = meta.SubfileType(r.parseUint32(t))
	case tag.TagPhotometricInterpretation:
		ifd0.PhotometricInterpretation = r.parseUint16(t)
	case tag.TagTileWidth:
		// ingnore tag
	case tag.TagTileLength:
		// ingnore tag
	case tag.TagTileOffsets:
		// ingnore tag
	case tag.TagTileByteCounts:
		// ingnore tag
	case tag.TagThumbnailOffset:
		ifd0.ThumbnailOffset = r.parseFirstUint32(t)
	case tag.TagThumbnailLength:
		ifd0.ThumbnailLength = r.parseFirstUint32(t)
	case tag.TagBitsPerSample:
		r.parseUint16List(t, ifd0.BitsPerSample[:])
	case tag.TagCompression:
		ifd0.Compression = meta.Compression(r.parseUint16(t))
	case tag.TagRowsPerStrip:
		// ingnore tag
	case tag.TagSubIFDs:
		r.parseSubIFDs(t)
	case tag.TagPlanarConfiguration:
		ifd0.PlanarConfiguration = r.parseUint16(t)
	case tag.TagXResolution:
		ifd0.XResolution = r.parseRationalValue(t).Float64()
	case tag.TagYResolution:
		ifd0.YResolution = r.parseRationalValue(t).Float64()
	case tag.TagResolutionUnit:
		ifd0.ResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	case tag.TagSamplesPerPixel:
		ifd0.SamplesPerPixel = r.parseUint16(t)
	case tag.TagImageWidth:
		ifd0.ImageWidth = r.parseUint32(t)
	case tag.TagImageLength:
		ifd0.ImageHeight = r.parseUint32(t)
	case tag.TagStripOffsets:
		ifd0.ImageOffset = r.parseFirstUint32(t)
	case tag.TagStripByteCounts:
		ifd0.ImageLength = r.parseFirstUint32(t)
	case tag.TagOrientation:
		ifd0.Orientation = meta.Orientation(r.parseUint16(t))
	// Color tags.
	case tag.TagYCbCrPositioning:
		ifd0.YCbCrPositioning = r.parseUint16(t)
	// Intentionally non-parsed IFD0 tags (recognized but not modeled).
	// Keep this list as the canonical location for IFD0 exclusions.
	case tag.TagCFARepeatPatternDim, tag.TagCFAPattern2, tag.TagReferenceBlackWhite, tag.TagTIFFEPStandardID,
		tag.TagCFAPlaneColor, tag.TagBlackLevelRepeatDim, tag.TagBlackLevel, tag.TagWhiteLevel,
		tag.TagColorMatrix1, tag.TagColorMatrix2, tag.TagAnalogBalance, tag.TagAsShotNeutral,
		tag.TagWhitePoint, tag.TagPrimaryChromaticities, tag.TagYCbCrCoefficients,
		tag.TagBaselineExposure, tag.TagBaselineNoise, tag.TagActiveArea, tag.TagDefaultScale,
		tag.TagDefaultCropOrigin, tag.TagDefaultCropSize, tag.TagDefaultUserCrop, tag.TagNewRawImageDigest,
		tag.TagCFALayout, tag.TagBayerGreenSplit, tag.TagBaselineSharpness, tag.TagLinearResponseLimit,
		tag.TagAntiAliasStrength, tag.TagShadowScale, tag.TagCalibrationIlluminant1, tag.TagCalibrationIlluminant2,
		tag.TagProfileEmbedPolicy, tag.TagNoiseProfile, tag.TagOpcodeList2, tag.TagLensSpecification,
		tag.TagIPTCNAA:
		return true
		// DNG extension tags.
	case tag.TagDNGAdobeData:
		r.dngMakerNote()
		r.parseDNGAdobeData(t)
	case tag.TagDNGVersion:
		dst := r.dngMakerNote()
		dst.DNGVersionCount = uint8(r.parseByteList(t, dst.DNGVersion[:]))
		if r.Exif.ImageType == imagetype.ImageTiff {
			r.Exif.ImageType = imagetype.ImageDNG
		}
	case tag.TagDNGBackwardVersion:
		dst := r.dngMakerNote()
		dst.DNGBackwardVersionCount = uint8(r.parseByteList(t, dst.DNGBackwardVersion[:]))
	case tag.TagUniqueCameraModel, tag.TagLocalizedCameraModel:
		dst := r.dngMakerNote()
		if dst.CameraModel != "" {
			return true
		}
		v := r.parseStringAllowUndefined(t)
		if v != "" {
			dst.CameraModel = v
		}
	case tag.TagOriginalRawFileName:
		dst := r.dngMakerNote()
		dst.OriginalRawFileName = r.parseStringAllowUndefined(t)
	case tag.TagProfileName:
		dst := r.dngMakerNote()
		dst.ProfileName = r.parseStringAllowUndefined(t)
	case tag.TagCameraSerial:
		if r.Exif.CameraSerial == "" {
			r.Exif.CameraSerial = r.parseString(t)
		}
	case tag.TagBestQualityScale:
		dst := r.dngMakerNote()
		dst.BestQualityScale = r.parseRationalValue(t)
	default:
		return false
	}
	return true
}

// parseMakeTag parses MakeTag and caches the maker-note make ID for subsequent maker-note parsing.
// The maker-note make ID is resolved from IFD0:Make
func (r *Reader) parseMakeTag(t tag.Entry) (makeID makernote.CameraMake, makeName string) {
	switch t.Type {
	case tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return makernote.CameraMakeUnknown, makeID.String()
	}

	var raw []byte
	if t.IsEmbedded() {
		size := t.Size()
		t.EmbeddedValue(r.state.buf[:4])
		raw = trimTrailingNULBytes(r.state.buf[:size])
	} else {
		buf, _, err := r.readTagBytes(t, uint32(len(r.state.buf)))
		if err != nil {
			return makernote.CameraMakeUnknown, makeID.String()
		}
		raw = trimTrailingNULBytes(buf)
	}

	if len(raw) == 0 {
		return makernote.CameraMakeUnknown, makeID.String()
	}

	if len(raw) <= 64 {
		makeID = makernote.IdentifyCameraMake(raw)
	} else {
		makeID = makernote.IdentifyCameraMakeString(string(raw))
	}

	if makeID == makernote.CameraMakeUnknown {
		makeName = string(raw)
	} else {
		makeName = makeID.String()
	}
	if r.InfoEnabled() {
		r.Info(3).
			Uint32("tagOffset", t.ValueOffset).
			Msg("parsed camera make")
	}
	return makeID, makeName
}

func (r *Reader) parseDNGAdobeData(t tag.Entry) {
	if t.Type != tag.TypeUndefined && t.Type != tag.TypeByte {
		return
	}
	if t.IsEmbedded() {
		return
	}
	if err := r.seekToTag(t); err != nil {
		return
	}

	size := t.Size()
	if size < 6 {
		_ = r.discard(int(size))
		return
	}

	header, err := r.fastRead(6)
	if err != nil {
		return
	}
	if len(header) < 6 || string(header[:6]) != "Adobe\x00" {
		remaining := int(size) - len(header)
		if remaining > 0 {
			_ = r.discard(remaining)
		}
		return
	}

	recordCount := 0
	recordBytesRemaining := size - 6
	for recordBytesRemaining >= 8 {
		recordHeader, readErr := r.fastRead(8)
		if readErr != nil || len(recordHeader) < 8 {
			return
		}
		recordBytesRemaining -= 8

		recordTag := string(recordHeader[:4])
		recordSize := binary.BigEndian.Uint32(recordHeader[4:8])
		if recordSize > recordBytesRemaining {
			_ = r.discard(int(recordBytesRemaining))
			break
		}

		recordCount++
		recordStart := r.po
		switch recordTag {
		case "MakN":
			r.parseDNGAdobeMakerNotes(recordStart, recordSize)
		default:
			_ = r.discard(int(recordSize))
		}

		recordEnd := recordStart + recordSize
		if r.po < recordEnd {
			_ = r.discard(int(recordEnd - r.po))
		}
		recordBytesRemaining -= recordSize

		if recordSize&1 != 0 {
			if recordBytesRemaining == 0 {
				break
			}
			_ = r.discard(1)
			recordBytesRemaining--
		}
	}

	if recordCount > 0xff {
		recordCount = 0xff
	}
	r.dngMakerNote().AdobeData.RecordCount = uint8(recordCount)
	if recordBytesRemaining > 0 {
		_ = r.discard(int(recordBytesRemaining))
	}
}

func (r *Reader) parseDNGAdobeMakerNotes(recordStart, recordSize uint32) {
	if recordSize < 6 {
		_ = r.discard(int(recordSize))
		return
	}

	header, err := r.fastRead(6)
	if err != nil || len(header) < 6 {
		return
	}

	byteOrder := utils.UnknownEndian
	switch {
	case header[0] == 'I' && header[1] == 'I':
		byteOrder = utils.LittleEndian
	case header[0] == 'M' && header[1] == 'M':
		byteOrder = utils.BigEndian
	}
	if byteOrder == utils.UnknownEndian {
		recordEnd := recordStart + recordSize
		if r.po < recordEnd {
			_ = r.discard(int(recordEnd - r.po))
		}
		return
	}

	originalOffset := binary.BigEndian.Uint32(header[2:6])
	headerLength := uint32(6)
	if recordSize >= 18 {
		prefix, peekErr := r.reader.Peek(12)
		if peekErr == nil && len(prefix) >= 4 &&
			prefix[0] == 0 && prefix[1] == 0 && prefix[2] == 0 && prefix[3] == 1 {
			if err := r.discard(12); err != nil {
				return
			}
			headerLength = 18
		}
	}
	if recordSize <= headerLength {
		return
	}

	dirStart := recordStart + headerLength
	if dirStart < originalOffset {
		recordEnd := recordStart + recordSize
		if r.po < recordEnd {
			_ = r.discard(int(recordEnd - r.po))
		}
		return
	}

	dng := r.dngMakerNote()
	dng.AdobeData.MakerNoteOriginalOffset = originalOffset
	dng.AdobeData.MakerNoteRecordLength = recordSize

	parent := tag.NewEntry(tag.TagMakerNote, tag.TypeUndefined, recordSize, recordStart, tag.ExifIFD, 0, byteOrder)
	child := tag.NewDirectory(byteOrder, tag.MakerNoteIFD, 0, dirStart, dirStart-originalOffset)
	queueStart := r.state.len
	if err := r.readMakerNoteDirectory(parent, child); err != nil {
		if r.WarnEnabled() {
			r.tagLogContext(r.Warn(3), parent).
				Err(err).
				Uint32("recordStart", recordStart).
				Uint32("recordSize", recordSize).
				Uint32("makerNoteOriginalOffset", originalOffset).
				Msg("failed parsing DNG Adobe maker notes")
		}
	}
	r.parseQueuedMakerNoteRange(queueStart)
}

func (r *Reader) parseQueuedMakerNoteRange(start uint32) {
	if start >= r.state.len {
		return
	}

	end := r.state.len
	if end-start > 1 {
		r.state.sortRange(start, end)
	}

	for i := start; i < end; i++ {
		t := r.state.tag[i]
		if t.IfdType != tag.MakerNoteIFD {
			continue
		}
		if t.IsIfd() {
			child := t.ChildDirectory()
			if child.Type == tag.MakerNoteIFD {
				if err := r.seekToTag(t); err != nil {
					continue
				}
				if err := r.readMakerNoteDirectory(t, child); err != nil {
					if r.WarnEnabled() {
						r.tagLogContext(r.Warn(3), t).
							Err(err).
							Str("childIFD", child.String()).
							Msg("failed parsing nested maker-note ifd")
					}
				}
			}
			continue
		}
		r.parseTag(t)
	}
	r.state.len = start
}

// parseExifTag parses ExifIFD/SubIFD tags into typed model fields.
//
// Non-parsed ExifIFD/SubIFD tags are currently handled by falling through to
// the default case (`return false`) when there is no modeled parser mapping.
//
// TODO(exiftool-coverage): Expand ExifIFD parity against ExifTool
// lib/Image/ExifTool/Exif.pm (%Image::ExifTool::Exif::Main) and
// https://exiftool.org/TagNames/EXIF.html (ExifIFD rows).
// Keep large/opaque payloads opt-in and avoid allocations in hot paths.
func (r *Reader) parseExifTag(t tag.Entry) bool {
	exifIfd := &r.Exif.ExifIFD
	ifd0 := &r.Exif.IFD0
	switch t.ID {
	// Time/date.
	case tag.TagDateTimeOriginal:
		exifIfd.setDate(t.ID, r.parseDate(t))
		exifIfd.markTagParsed(t.ID)
	case tag.TagDateTimeDigitized:
		exifIfd.setDate(t.ID, r.parseDate(t))
		exifIfd.markTagParsed(t.ID)
	case tag.TagSubSecTime:
		value, raw := r.parseSubSecTimeParts(t)
		ifd0.setSubSec(raw, value)
		ifd0.markTagParsed(t.ID)
	case tag.TagSubSecTimeOriginal:
		value, raw := r.parseSubSecTimeParts(t)
		exifIfd.setSubSec(t.ID, raw, value)
		exifIfd.markTagParsed(t.ID)
	case tag.TagSubSecTimeDigitized:
		value, raw := r.parseSubSecTimeParts(t)
		exifIfd.setSubSec(t.ID, raw, value)
		exifIfd.markTagParsed(t.ID)
	case tag.TagOffsetTime:
		ifd0.setOffset(r.parseOffsetTime(t))
		ifd0.markTagParsed(t.ID)
	case tag.TagOffsetTimeOriginal:
		exifIfd.setOffset(t.ID, r.parseOffsetTime(t))
		exifIfd.markTagParsed(t.ID)
	case tag.TagOffsetTimeDigitized:
		exifIfd.setOffset(t.ID, r.parseOffsetTime(t))
		exifIfd.markTagParsed(t.ID)

	// Text/meta.
	// TODO(exiftool-coverage): Add selected ExifIFD text fields with bounded
	// parsing and explicit truncation limits where needed:
	// ImageUniqueID(0xa420), CompositeImageCount(0xa461), and
	// CompositeImageExposureTimes(0xa462) display summaries.
	case tag.TagExifVersion:
		exifIfd.ExifVersion = r.parseExifVersion(t)
	case tag.TagLensMake:
		exifIfd.LensMake = r.parseString(t)
	case tag.TagLensModel:
		exifIfd.LensModel = r.parseString(t)
	case tag.TagLensSerialNumber:
		exifIfd.LensSerial = r.parseString(t)
	case tag.TagCameraOwnerName:
		exifIfd.CameraOwnerName = r.parseString(t)
	case tag.TagBodySerialNumber:
		exifIfd.BodySerialNumber = r.parseString(t)
	case tag.TagUserComment:
		exifIfd.UserComment = r.parseExifUserComment(t)
	case tag.TagFlashpixVersion:
		exifIfd.FlashpixVersion = r.parseStringAllowUndefined(t)
	case tag.TagDeviceSettingDescription:
		// Not parsed.
		// The payload is often large and not needed in the hot parse path.

	// Image characterization.
	// TODO(exiftool-coverage): Consider lightweight modeling for ExifIFD image
	// characterization tags used by ExifTool:
	// OECF(0x8828), SpatialFrequencyResponse(0x920c), SubjectLocation(0xa214),
	// and CFAPattern(0xa302). Prefer compact fixed-size summaries over raw blobs.
	case tag.TagPixelXDimension:
		exifIfd.PixelXDimension = r.parseUint32(t)
	case tag.TagPixelYDimension:
		exifIfd.PixelYDimension = r.parseUint32(t)
	case tag.TagInteropIFDPointer:
		exifIfd.interopIFDPointer = r.parseUint32(t)
	case tag.TagInteropIndex:
		exifIfd.InteropIndex = r.parseStringAllowUndefined(t)
	case tag.TagInteropVersion:
		exifIfd.InteropVersion = r.parseStringAllowUndefined(t)
	case tag.TagRelatedImageWidth:
		exifIfd.RelatedImageWidth = r.parseUint32(t)
	case tag.TagRelatedImageHeight:
		exifIfd.RelatedImageHeight = r.parseUint32(t)
	case tag.TagColorSpace:
		exifIfd.ColorSpace = r.parseUint16(t)
	case tag.TagLensSpecification:
		exifIfd.LensInfo = r.parseLensInfo(t)
	case tag.TagComponentsConfiguration:
		r.parseByteList(t, exifIfd.ComponentsConfiguration[:])
	case tag.TagCompressedBitsPerPixel:
		exifIfd.CompressedBitsPerPixel = r.parseRationalValue(t).Float64()
	case tag.TagFocalPlaneXResolution:
		exifIfd.FocalPlaneXResolution, exifIfd.focalPlaneXResolutionState = r.parseUnsignedRationalFloat64(t)
	case tag.TagFocalPlaneYResolution:
		exifIfd.FocalPlaneYResolution, exifIfd.focalPlaneYResolutionState = r.parseUnsignedRationalFloat64(t)
	case tag.TagFocalPlaneResolutionUnit:
		exifIfd.FocalPlaneResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	case tag.TagSubjectArea:
		r.parseUint16List(t, exifIfd.SubjectArea[:])
	case tag.TagGamma:
		exifIfd.Gamma = r.parseRationalValue(t).Float64()

	// Exposure/optics.
	// TODO(exiftool-coverage): Add ISO-family ExifIFD fields documented by
	// ExifTool for EXIF 2.31+ parity:
	// StandardOutputSensitivity(0x8831), ISOSpeed(0x8833),
	// ISOSpeedLatitudeyyy(0x8834), ISOSpeedLatitudezzz(0x8835).
	// Keep parse paths branch-light and avoid allocations.
	case tag.TagExposureTime:
		exifIfd.ExposureTime = r.parseExposureTime(t)
	case tag.TagShutterSpeedValue:
		exifIfd.ShutterSpeedValue = r.parseShutterSpeed(t)
	case tag.TagFNumber:
		exifIfd.FNumber = r.parseAperture(t)
	case tag.TagApertureValue:
		exifIfd.ApertureValue = r.parseApexAperture(t)
		if exifIfd.FNumber == 0 && apertureIsFinite(exifIfd.ApertureValue) {
			exifIfd.FNumber = exifIfd.ApertureValue
		}
	case tag.TagMaxApertureValue:
		exifIfd.MaxApertureValue = r.parseApexAperture(t)
	case tag.TagSubjectDistance:
		exifIfd.SubjectDistance, exifIfd.subjectDistanceState = r.parseUnsignedRationalFloat64(t)
	case tag.TagBrightnessValue:
		exifIfd.BrightnessValue = r.parseSignedRationalFloat32(t)
	case tag.TagExposureProgram:
		exifIfd.ExposureProgram = meta.ExposureProgram(r.parseUint16(t))
	case tag.TagSensitivityType:
		exifIfd.SensitivityType = r.parseUint16(t)
	case tag.TagRecommendedExposureIndex:
		exifIfd.RecommendedExposureIndex = r.parseUint32(t)
	case tag.TagExposureBiasValue:
		exifIfd.ExposureBias = r.parseExposureBias(t)
	case tag.TagExposureMode:
		exifIfd.ExposureMode = meta.ExposureMode(r.parseUint16(t))
	case tag.TagMeteringMode:
		exifIfd.MeteringMode = meta.MeteringMode(r.parseUint16(t))
	case tag.TagLightSource:
		exifIfd.LightSource = r.parseUint16(t)
	case tag.TagISOSpeedRatings:
		exifIfd.ISOSpeedRatings = r.parseUint32(t)
	case tag.TagFlash:
		exifIfd.Flash = meta.Flash(r.parseUint16(t))
	case tag.TagFocalLength:
		exifIfd.FocalLength = r.parseFocalLength(t)
	case tag.TagFocalLengthIn35mmFilm:
		exifIfd.FocalLengthIn35mmFormat = r.parseFocalLength(t)
	case tag.TagExposureIndex:
		exifIfd.ExposureIndex, exifIfd.exposureIndexState = r.parseUnsignedRationalFloat64(t)

	// Capture/mode.
	// TODO(exiftool-coverage): Evaluate adding environmental capture tags
	// surfaced by ExifTool:
	// AmbientTemperature(0x9400), Humidity(0x9401), Pressure(0x9402),
	// WaterDepth(0x9403), Acceleration(0x9404), CameraElevationAngle(0x9405).
	// These are niche but useful for action cameras and drone media.
	case tag.TagSensingMethod:
		exifIfd.SensingMethod = r.parseUint16(t)
	case tag.TagFileSource:
		exifIfd.FileSource = r.parseSceneType(t)
	case tag.TagSceneType:
		exifIfd.SceneType = r.parseSceneType(t)
	case tag.TagCustomRendered:
		exifIfd.CustomRendered = r.parseUint16(t)
	case tag.TagWhiteBalance:
		exifIfd.WhiteBalance = r.parseUint16(t)
	case tag.TagDigitalZoomRatio:
		exifIfd.DigitalZoomRatio = float32(r.parseRationalValue(t).Float64())
	case tag.TagSceneCaptureType:
		exifIfd.SceneCaptureType = r.parseUint16(t)
	case tag.TagGainControl:
		exifIfd.GainControl = r.parseUint16(t)
	case tag.TagContrast:
		exifIfd.Contrast = r.parseUint16(t)
	case tag.TagSaturation:
		exifIfd.Saturation = r.parseUint16(t)
	case tag.TagSharpness:
		exifIfd.Sharpness = r.parseUint16(t)
	case tag.TagSubjectDistanceRange:
		exifIfd.SubjectDistanceRange = r.parseUint16(t)
	case tag.TagCompositeImage:
		exifIfd.CompositeImage = r.parseUint16(t)
	default:
		return false
	}
	return true
}

// parseImageIFDTag parses the requested value from EXIF metadata.
func (r *Reader) parseImageIFDTag(t tag.Entry, dst *ImageIFD) bool {
	if dst == nil {
		dst = &ImageIFD{}
	}
	switch t.ID {
	case tag.TagSubfileType:
		dst.SubfileType = meta.SubfileType(r.parseUint32(t))
	case tag.TagBitsPerSample:
		// ingnore tag
	case tag.TagCompression:
		dst.Compression = meta.Compression(r.parseUint16(t))
	case tag.TagXResolution:
		dst.XResolution = r.parseRationalValue(t).Float64()
	case tag.TagYResolution:
		dst.YResolution = r.parseRationalValue(t).Float64()
	case tag.TagResolutionUnit:
		dst.ResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	case tag.TagImageWidth:
		dst.ImageWidth = r.parseUint32(t)
	case tag.TagImageLength:
		dst.ImageHeight = r.parseUint32(t)
	case tag.TagMake:
		_, dst.Make = r.parseMakeTag(t)
	case tag.TagModel:
		dst.Model = r.parseString(t)
	case tag.TagImageDescription:
		dst.ImageDescription = r.parseString(t)
	case tag.TagSoftware:
		dst.Software = r.parseString(t)
	case tag.TagDateTime:
		dst.setDate(r.parseDate(t))
		dst.markTagParsed(t.ID)
	case tag.TagStripOffsets, tag.TagThumbnailOffset:
		dst.ImageOffset = r.parseFirstUint32(t)
	case tag.TagStripByteCounts, tag.TagThumbnailLength:
		dst.ImageLength = r.parseFirstUint32(t)
	case tag.TagOrientation:
		dst.Orientation = meta.Orientation(r.parseUint16(t))
	default:
		return false
	}
	return true
}

// parseSubSecTime parses the requested value from EXIF metadata.
func (r *Reader) parseSubSecTime(t tag.Entry) uint16 {
	value, _ := r.parseSubSecTimeParts(t)
	return value
}

func (r *Reader) parseSubSecTimeParts(t tag.Entry) (uint16, string) {
	switch t.Type {
	case tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return 0, ""
	}
	var raw []byte
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		raw = trimASCIIWhitespace(trimNULBuffer(r.state.buf[:4]))
		return uint16(parseStrUint(raw)), string(raw)
	}
	buf, _, err := r.readTagBytes(t, 16)
	if err != nil {
		return 0, ""
	}
	raw = trimASCIIWhitespace(trimNULBuffer(buf))
	return uint16(parseStrUint(raw)), string(raw)
}

func apertureIsFinite(v meta.Aperture) bool {
	f := float64(v)
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}

func apexApertureToFNumber(v float64) meta.Aperture {
	if v == 0 {
		return 0
	}
	// Some cameras write a large sentinel for "infinite" aperture. Mirror
	// ExifTool behavior and preserve that as +Inf instead of a huge float.
	if math.Abs(v) > 1024 {
		return meta.Aperture(math.Inf(1))
	}
	return meta.Aperture(math.Exp2(v * 0.5))
}

func apexShutterSpeedToSeconds(v float64) meta.ShutterSpeed {
	if math.IsNaN(v) || math.Abs(v) >= 100 {
		return 0
	}
	return meta.ShutterSpeed(math.Exp2(-v))
}
