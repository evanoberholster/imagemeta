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
	tagType = tagTypeFor(directory.Type, tagID, tagType)
	if !tag.NewEntry(tagID, tagType, unitCount, valueOffset, directory.Type, directory.Index, directory.ByteOrder).IsEmbedded() {
		valueOffset += directory.BaseOffset
	}

	entry := tag.NewEntry(tagID, tagType, unitCount, valueOffset, directory.Type, directory.Index, directory.ByteOrder)
	if !tagType.IsValid() {
		return entry, tag.ErrTagTypeNotValid
	}
	return entry, nil
}

// tagTypeFor resolves the effective tag type for parser dispatch.
func tagTypeFor(directoryType tag.IfdType, id tag.ID, typ tag.Type) tag.Type {
	if (typ.Is(tag.TypeLong) || typ.Is(tag.TypeUndefined)) && tagUsesIfdType(directoryType, id) {
		return tag.TypeIfd
	}
	return typ
}

func tagUsesIfdType(directoryType tag.IfdType, id tag.ID) bool {
	switch directoryType {
	case tag.IFD0:
		return id == tag.TagExifIFDPointer || id == tag.TagGPSIFDPointer
	case tag.ExifIFD:
		return id == tag.TagMakerNote
	default:
		return false
	}
}

// addTag appends a tag to parser state while preserving parse order constraints.
func (r *Reader) addTag(t tag.Entry) {
	if t.ValueOffset < r.po {
		if r.warnEnabled() {
			r.warnTagAdd(t, "ignoring reverse-offset tag")
		}
		return
	}
	if !r.state.addTag(t) && r.warnEnabled() {
		r.warnTagQueueFull(t)
	}
}

// warnTagContext logs tag metadata with a caller-supplied message.
func (r *Reader) warnTagContext(t tag.Entry, msg string, includeQueueMax bool) {
	e := r.warn().
		Uint16("tagID", uint16(t.ID)).
		Str("tagName", tag.NameFor(t.IfdType, t.ID)).
		Uint16("tagType", uint16(t.Type)).
		Str("ifd", t.IfdType.String()).
		Int8("ifdIndex", t.IfdIndex).
		Uint32("units", t.UnitCount).
		Uint32("tagOffset", t.ValueOffset).
		Uint32("readerOffset", r.po)
	if includeQueueMax {
		e.Int("tagQueueMax", tagQueueMax)
	}
	e.Msg(msg)
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
		if r.warnEnabled() {
			r.warnTagContext(t, "subifd queue capacity reached; skipping parse", true)
		}
		return
	}
	offsetRemaining := len(r.Exif.IFD0.subIFDOffsets) - int(r.Exif.IFD0.subIFDOffsetCount)
	if offsetRemaining <= 0 {
		if r.warnEnabled() {
			r.warnTagContext(t, "subifd offset capacity reached; skipping parse", false)
		}
		return
	}

	maxEntries := int(t.UnitCount)
	queueRemaining := int(tagQueueMax - r.state.len)
	if maxEntries > queueRemaining {
		maxEntries = queueRemaining
	}
	if maxEntries > offsetRemaining {
		maxEntries = offsetRemaining
	}
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
			if r.warnEnabled() {
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
		subType := tag.SubIFD0
		if i < int(tag.SubIFD7-tag.SubIFD0)+1 {
			subType = tag.IfdType(uint8(tag.SubIFD0) + uint8(i))
		}
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
	switch t.IfdType {
	case tag.IFD0:
		if !r.parseIFD0Tag(t) {
			return
		}
	case tag.IFD1:
		if !r.parseImageIFDTag(t, &r.Exif.IFD1) {
			return
		}
	case tag.IFD2:
		if !r.parseImageIFDTag(t, &r.Exif.IFD2) {
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
	r.Exif.markTagParsed(uint16(t.ID))
}

// parseIFD0Tag parses IFD0 tags into typed model fields.
//
// Non-parsed IFD0 tags are documented in the explicit "intentionally non-parsed"
// case branch below. These tags are treated as handled for coverage/reporting
// parity but are not mapped into the Exif model.
func (r *Reader) parseIFD0Tag(t tag.Entry) bool {
	if r.parseIFD0TextTag(t) || r.parseIFD0ImageTag(t) || r.parseIFD0DNGTag(t) {
		return true
	}
	// Intentionally non-parsed IFD0 tags (recognized but not modeled).
	// Keep this list as the canonical location for IFD0 exclusions.
	switch t.ID {
	case tag.TagCFARepeatPatternDim, tag.TagCFAPattern2, tag.TagReferenceBlackWhite, tag.TagTIFFEPStandardID,
		tag.TagCFAPlaneColor, tag.TagBlackLevelRepeatDim, tag.TagBlackLevel, tag.TagWhiteLevel,
		tag.TagColorMatrix1, tag.TagColorMatrix2, tag.TagAnalogBalance, tag.TagAsShotNeutral,
		tag.TagBaselineExposure, tag.TagBaselineNoise, tag.TagActiveArea, tag.TagDefaultScale,
		tag.TagDefaultCropOrigin, tag.TagDefaultCropSize, tag.TagDefaultUserCrop, tag.TagNewRawImageDigest,
		tag.TagCFALayout, tag.TagBayerGreenSplit, tag.TagBaselineSharpness, tag.TagLinearResponseLimit,
		tag.TagAntiAliasStrength, tag.TagShadowScale, tag.TagCalibrationIlluminant1, tag.TagCalibrationIlluminant2,
		tag.TagProfileEmbedPolicy, tag.TagNoiseProfile, tag.TagOpcodeList2, tag.TagLensSpecification:
		return true
	default:
		return false
	}
}

// parseIFD0TextTag parses IFD0 tags with string or date types.
func (r *Reader) parseIFD0TextTag(t tag.Entry) bool {
	switch t.ID {
	case tag.TagDateTime:
		modifyDate := r.parseDate(t)
		r.Exif.Time.ModifyDate = modifyDate
		r.Exif.IFD0.ModifyDate = modifyDate
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagDateTimeOriginal:
		dateTimeOriginal := r.parseDate(t)
		r.Exif.Time.DateTimeOriginal = dateTimeOriginal
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagMake:
		r.Exif.CameraMakeID, r.Exif.IFD0.Make = r.parseMakeTag(t)
	case tag.TagModel:
		r.Exif.IFD0.Model = r.parseString(t)
	case tag.TagArtist:
		r.Exif.IFD0.Artist = r.parseStringTrimRightSpaceNewline(t)
	case tag.TagCopyright:
		r.Exif.IFD0.Copyright = r.parseDisplayStringTrimRightSpaceNewline(t, 512)
	case tag.TagApplicationNotes:
		// TODO: TagApplicationNotes parsing intentionally disabled for now.
		// The payload is often large and not needed in the hot parse path.
	case tag.TagPrintIM:
		// ignore tag
		//r.Exif.IFD0.PrintIM = r.parseDisplayString(t, 512)
	case tag.TagImageDescription:
		r.Exif.IFD0.ImageDescription = r.parseString(t)
	case tag.TagSoftware:
		r.Exif.IFD0.Software = r.parseString(t)
	default:
		return false
	}
	return true
}

// parseMakeTag parses MakeTag and caches the maker-note make ID for subsequent maker-note parsing.
// The maker-note make ID is resolved from IFD0:Make
func (r *Reader) parseMakeTag(t tag.Entry) (makeID makernote.CameraMake, Make string) {
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
		Make = string(raw)
		return
	}
	return makeID, makeID.String()
}

// parseIFD0ImageTag parses IFD0 image geometry and layout tags.
func (r *Reader) parseIFD0ImageTag(t tag.Entry) bool {
	switch t.ID {
	case tag.TagSubfileType:
		r.Exif.IFD0.SubfileType = meta.SubfileType(r.parseUint32(t))
	case tag.TagTileWidth:
		// ingnore tag
	case tag.TagTileLength:
		// ingnore tag
	case tag.TagTileOffsets:
		// ingnore tag
	case tag.TagTileByteCounts:
		// ingnore tag
	case tag.TagThumbnailOffset:
		r.Exif.IFD0.ThumbnailOffset = r.parseFirstUint32(t)
	case tag.TagThumbnailLength:
		r.Exif.IFD0.ThumbnailLength = r.parseFirstUint32(t)
	case tag.TagBitsPerSample:
		// ingnore tag
	case tag.TagCompression:
		r.Exif.IFD0.Compression = meta.Compression(r.parseUint16(t))
	case tag.TagRowsPerStrip:
		// ingnore tag
	case tag.TagSubIFDs:
		r.parseSubIFDs(t)
	case tag.TagPlanarConfiguration:
		// ingnore tag
	case tag.TagXResolution:
		r.Exif.IFD0.XResolution = r.parseRationalValue(t)
	case tag.TagYResolution:
		r.Exif.IFD0.YResolution = r.parseRationalValue(t)
	case tag.TagResolutionUnit:
		r.Exif.IFD0.ResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	case tag.TagImageWidth:
		r.Exif.IFD0.ImageWidth = r.parseUint32(t)
	case tag.TagImageLength:
		r.Exif.IFD0.ImageHeight = r.parseUint32(t)
	case tag.TagStripOffsets:
		r.Exif.IFD0.ImageOffset = r.parseFirstUint32(t)
	case tag.TagStripByteCounts:
		r.Exif.IFD0.ImageLength = r.parseFirstUint32(t)
	case tag.TagOrientation:
		r.Exif.IFD0.Orientation = meta.Orientation(r.parseUint16(t))
	default:
		return false
	}
	return true
}

// parseIFD0DNGTag parses IFD0 DNG extension fields.
func (r *Reader) parseIFD0DNGTag(t tag.Entry) bool {
	switch t.ID {
	case tag.TagDNGAdobeData:
		r.parseDNGAdobeData(t)
	case tag.TagDNGVersion:
		r.Exif.DNG.DNGVersionCount = uint8(r.parseByteList(t, r.Exif.DNG.DNGVersion[:]))
		if r.Exif.ImageType == imagetype.ImageTiff {
			r.Exif.ImageType = imagetype.ImageDNG
		}
	case tag.TagDNGBackwardVersion:
		r.Exif.DNG.DNGBackwardVersionCount = uint8(r.parseByteList(t, r.Exif.DNG.DNGBackwardVersion[:]))
	case tag.TagUniqueCameraModel, tag.TagLocalizedCameraModel:
		if r.Exif.DNG.CameraModel != "" {
			return true
		}
		v := r.parseStringAllowUndefined(t)
		if v != "" {
			r.Exif.DNG.CameraModel = v
		}
	case tag.TagOriginalRawFileName:
		r.Exif.DNG.OriginalRawFileName = r.parseStringAllowUndefined(t)
	case tag.TagProfileName:
		r.Exif.DNG.ProfileName = r.parseStringAllowUndefined(t)
	case tag.TagCameraSerial:
		if r.Exif.CameraSerial == "" {
			r.Exif.CameraSerial = r.parseString(t)
		}
	case tag.TagBestQualityScale:
		r.Exif.DNG.BestQualityScale = r.parseRationalValue(t)
	default:
		return false
	}
	return true
}

func (r *Reader) parseDNGAdobeData(t tag.Entry) {
	if t.Type != tag.TypeUndefined && t.Type != tag.TypeByte {
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
	r.Exif.DNG.AdobeData.RecordCount = uint8(recordCount)
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

	r.Exif.DNG.AdobeData.MakerNoteOriginalOffset = originalOffset
	r.Exif.DNG.AdobeData.MakerNoteRecordLength = recordSize

	parent := tag.NewEntry(tag.TagMakerNote, tag.TypeUndefined, recordSize, recordStart, tag.ExifIFD, 0, byteOrder)
	child := tag.NewDirectory(byteOrder, tag.MakerNoteIFD, 0, dirStart, dirStart-originalOffset)
	queueStart := r.state.len
	if err := r.readMakerNoteDirectory(parent, child); err != nil && r.warnEnabled() {
		r.warn().Err(err).Uint32("tagOffset", recordStart).Msg("failed parsing DNG Adobe maker notes")
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
				if err := r.readMakerNoteDirectory(t, child); err != nil && r.warnEnabled() {
					r.warn().Err(err).Uint16("tagID", uint16(t.ID)).Msg("failed parsing nested maker-note ifd")
				}
			}
			continue
		}
		r.parseTag(t)
	}
	r.state.len = start
}

// parseIFD0PanasonicRawTag parses Panasonic RW2/RWL root-IFD tags.
// func (r *Reader) parseIFD0PanasonicRawTag(t tag.Entry) bool {
// 	if r.Exif.ImageType != imagetype.ImagePanaRAW {
// 		return false
// 	}
// 	switch t.ID {
// 	case tag.TagPanasonicRawVersion:
// 		r.parseByteList(t, r.Exif.PanasonicRaw.Version[:])
// 	case tag.TagPanasonicSensorWidth:
// 		r.Exif.PanasonicRaw.SensorWidth = r.parseUint16(t)
// 	case tag.TagPanasonicSensorHeight:
// 		r.Exif.PanasonicRaw.SensorHeight = r.parseUint16(t)
// 	case tag.TagPanasonicBitsPerSample:
// 		r.Exif.PanasonicRaw.BitsPerSample = r.parseUint16(t)
// 	case tag.TagPanasonicCompression:
// 		r.Exif.PanasonicRaw.Compression = r.parseUint16(t)
// 	case tag.TagPanasonicISO:
// 		r.Exif.PanasonicRaw.ISO = uint32(r.parseUint16(t))
// 	case tag.TagPanasonicISOHighPrecision:
// 		r.Exif.PanasonicRaw.ISO = r.parseUint32(t)
// 	case tag.TagNoiseReductionParams:
// 		// Not parsed
// 	case tag.TagWBInfo2:
// 		// Not parsed
// 	case tag.TagPanasonicRawFormat:
// 		r.Exif.PanasonicRaw.RawFormat = r.parseUint16(t)
// 	case tag.TagJpgFromRaw:
// 		// TODO: parse JpgFromRaw payload into typed preview metadata if needed.
// 		// Keep offset/length only to avoid large allocations for embedded JPEGs.
// 		r.Exif.PanasonicRaw.JpgFromRawOffset = t.ValueOffset
// 		r.Exif.PanasonicRaw.JpgFromRawLength = t.UnitCount
// 	case tag.TagPanasonicRawDataOffset:
// 		r.Exif.PanasonicRaw.RawDataOffset = r.parseUint32(t)
// 	case tag.TagPanasonicDistortionInfo:
// 		// Not parsed
// 	case tag.TagPanasonicCropTop:
// 		r.Exif.PanasonicRaw.CropTop = r.parseUint16(t)
// 	case tag.TagPanasonicCropLeft:
// 		r.Exif.PanasonicRaw.CropLeft = r.parseUint16(t)
// 	case tag.TagPanasonicCropBottom:
// 		r.Exif.PanasonicRaw.CropBottom = r.parseUint16(t)
// 	case tag.TagPanasonicCropRight:
// 		r.Exif.PanasonicRaw.CropRight = r.parseUint16(t)
// 	case tag.TagPanasonicTitle:
// 		r.Exif.PanasonicRaw.Title = r.parseStringAllowUndefined(t)
// 	case tag.TagPanasonicTitle2:
// 		r.Exif.PanasonicRaw.Title2 = r.parseStringAllowUndefined(t)
// 	default:
// 		return false
// 	}
// 	return true
// }

// parseExifTag parses ExifIFD/SubIFD tags into typed model fields.
//
// Non-parsed ExifIFD/SubIFD tags are currently handled by falling through to
// the default case (`return false`) when there is no modeled parser mapping.
func (r *Reader) parseExifTag(t tag.Entry) bool {
	switch t.ID {
	case tag.TagDateTimeOriginal:
		r.Exif.Time.DateTimeOriginal = r.parseDate(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagDateTimeDigitized:
		r.Exif.Time.CreateDate = r.parseDate(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagSubSecTime:
		r.Exif.Time.SubSecTime = r.parseSubSecTime(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagSubSecTimeOriginal:
		r.Exif.Time.SubSecTimeOriginal = r.parseSubSecTime(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagSubSecTimeDigitized:
		r.Exif.Time.SubSecTimeDigitized = r.parseSubSecTime(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagOffsetTime:
		r.Exif.Time.OffsetTime = r.parseOffsetTime(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagOffsetTimeOriginal:
		r.Exif.Time.OffsetTimeOriginal = r.parseOffsetTime(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagOffsetTimeDigitized:
		r.Exif.Time.OffsetTimeDigitized = r.parseOffsetTime(t)
		r.Exif.Time.markTagParsed(t.ID)
	case tag.TagExifVersion:
		r.Exif.ExifIFD.ExifVersion = r.parseStringAllowUndefined(t)
	case tag.TagLensMake:
		r.Exif.ExifIFD.LensMake = r.parseString(t)
	case tag.TagLensModel:
		r.Exif.ExifIFD.LensModel = r.parseString(t)
	case tag.TagLensSerialNumber:
		r.Exif.ExifIFD.LensSerial = r.parseString(t)
	case tag.TagCameraOwnerName:
		r.Exif.ExifIFD.CameraOwnerName = r.parseString(t)
		if r.Exif.IFD0.Artist == "" {
			r.Exif.IFD0.Artist = r.Exif.ExifIFD.CameraOwnerName
		}
	case tag.TagBodySerialNumber:
		r.Exif.ExifIFD.BodySerialNumber = r.parseString(t)
		if r.Exif.CameraSerial == "" {
			r.Exif.CameraSerial = r.Exif.ExifIFD.BodySerialNumber
		}
	case tag.TagUserComment:
		r.Exif.ExifIFD.UserComment = r.parseExifUserComment(t)
	case tag.TagFlashpixVersion:
		r.Exif.ExifIFD.FlashpixVersion = r.parseStringAllowUndefined(t)
	case tag.TagDeviceSettingDescription:
		r.Exif.ExifIFD.DeviceSettingDescription = r.parseStringAllowUndefined(t)
	case tag.TagPixelXDimension:
		r.Exif.ExifIFD.PixelXDimension = r.parseUint32(t)
		if r.Exif.IFD0.ImageWidth == 0 {
			r.Exif.IFD0.ImageWidth = r.Exif.ExifIFD.PixelXDimension
		}
	case tag.TagPixelYDimension:
		r.Exif.ExifIFD.PixelYDimension = r.parseUint32(t)
		if r.Exif.IFD0.ImageHeight == 0 {
			r.Exif.IFD0.ImageHeight = r.Exif.ExifIFD.PixelYDimension
		}
	case tag.TagInteropIFDPointer:
		r.Exif.ExifIFD.InteropIFDPointer = r.parseUint32(t)
	case tag.TagColorSpace:
		r.Exif.ExifIFD.ColorSpace = r.parseUint16(t)
	case tag.TagLensSpecification:
		r.Exif.ExifIFD.LensInfo = r.parseLensInfo(t)
	case tag.TagComponentsConfiguration:
		r.parseByteList(t, r.Exif.ExifIFD.ComponentsConfiguration[:])
	case tag.TagCompressedBitsPerPixel:
		r.Exif.ExifIFD.CompressedBitsPerPixel = r.parseRationalValue(t)
	case tag.TagFocalPlaneXResolution:
		r.Exif.ExifIFD.FocalPlaneXResolution = r.parseRationalValue(t)
	case tag.TagFocalPlaneYResolution:
		r.Exif.ExifIFD.FocalPlaneYResolution = r.parseRationalValue(t)
	case tag.TagFocalPlaneResolutionUnit:
		r.Exif.ExifIFD.FocalPlaneResolutionUnit = meta.ResolutionUnit(r.parseUint16(t))
	case tag.TagSubjectArea:
		r.parseUint16List(t, r.Exif.ExifIFD.SubjectArea[:])
	case tag.TagExposureTime:
		r.Exif.ExifIFD.ExposureTime = r.parseExposureTime(t)
	case tag.TagShutterSpeedValue:
		r.Exif.ExifIFD.ShutterSpeedValue = r.parseShutterSpeed(t)
	case tag.TagFNumber:
		r.Exif.ExifIFD.FNumber = r.parseAperture(t)
	case tag.TagApertureValue:
		r.Exif.ExifIFD.ApertureValue = r.parseApexAperture(t)
		if r.Exif.ExifIFD.FNumber == 0 && apertureIsFinite(r.Exif.ExifIFD.ApertureValue) {
			r.Exif.ExifIFD.FNumber = apertureValueToFNumber(r.Exif.ExifIFD.ApertureValue)
		}
	case tag.TagMaxApertureValue:
		r.Exif.ExifIFD.MaxApertureValue = r.parseApexAperture(t)
	case tag.TagSubjectDistance:
		r.Exif.ExifIFD.SubjectDistance = r.parseRationalValue(t)
	case tag.TagBrightnessValue:
		r.Exif.ExifIFD.BrightnessValue = r.parseSignedRationalFloat32(t)
	case tag.TagExposureProgram:
		r.Exif.ExifIFD.ExposureProgram = meta.ExposureProgram(r.parseUint16(t))
	case tag.TagSensitivityType:
		r.Exif.ExifIFD.SensitivityType = r.parseUint16(t)
	case tag.TagRecommendedExposureIndex:
		r.Exif.ExifIFD.RecommendedExposureIndex = r.parseUint32(t)
	case tag.TagExposureBiasValue:
		r.Exif.ExifIFD.ExposureBias = r.parseExposureBias(t)
	case tag.TagExposureMode:
		r.Exif.ExifIFD.ExposureMode = meta.ExposureMode(r.parseUint16(t))
	case tag.TagMeteringMode:
		r.Exif.ExifIFD.MeteringMode = meta.MeteringMode(r.parseUint16(t))
	case tag.TagLightSource:
		r.Exif.ExifIFD.LightSource = r.parseUint16(t)
	case tag.TagISOSpeedRatings:
		r.Exif.ExifIFD.ISOSpeedRatings = r.parseUint32(t)
	case tag.TagFlash:
		r.Exif.ExifIFD.Flash = meta.Flash(r.parseUint16(t))
	case tag.TagFocalLength:
		r.Exif.ExifIFD.FocalLength = r.parseFocalLength(t)
	case tag.TagFocalLengthIn35mmFilm:
		r.Exif.ExifIFD.FocalLengthIn35mmFormat = r.parseFocalLength(t)
	case tag.TagExposureIndex:
		r.Exif.ExifIFD.ExposureIndex = r.parseRationalValue(t)
	case tag.TagSensingMethod:
		r.Exif.ExifIFD.SensingMethod = r.parseUint16(t)
	case tag.TagFileSource:
		r.Exif.ExifIFD.FileSource = r.parseSceneType(t)
	case tag.TagSceneType:
		r.Exif.ExifIFD.SceneType = r.parseSceneType(t)
	case tag.TagCustomRendered:
		r.Exif.ExifIFD.CustomRendered = r.parseUint16(t)
	case tag.TagWhiteBalance:
		r.Exif.ExifIFD.WhiteBalance = r.parseUint16(t)
	case tag.TagDigitalZoomRatio:
		r.Exif.ExifIFD.DigitalZoomRatio = r.parseRationalValue(t)
	case tag.TagSceneCaptureType:
		r.Exif.ExifIFD.SceneCaptureType = r.parseUint16(t)
	case tag.TagGainControl:
		r.Exif.ExifIFD.GainControl = r.parseUint16(t)
	case tag.TagContrast:
		r.Exif.ExifIFD.Contrast = r.parseUint16(t)
	case tag.TagSaturation:
		r.Exif.ExifIFD.Saturation = r.parseUint16(t)
	case tag.TagSharpness:
		r.Exif.ExifIFD.Sharpness = r.parseUint16(t)
	case tag.TagSubjectDistanceRange:
		r.Exif.ExifIFD.SubjectDistanceRange = r.parseUint16(t)
	case tag.TagCompositeImage:
		r.Exif.ExifIFD.CompositeImage = r.parseUint16(t)
	default:
		return false
	}
	return true
}

// parseGPSTag parses GPS IFD tags into typed model fields.
//
// Non-parsed GPS tags are currently handled by falling through to
// the default path (`return false`) when there is no modeled parser mapping.
func (r *Reader) parseGPSTag(t tag.Entry) bool {
	switch t.ID {
	case tag.TagGPSVersionID:
		r.parseByteList(t, r.Exif.GPS.versionID[:])
	case tag.TagGPSDifferential:
		r.Exif.GPS.differential = r.parseUint16(t)
	case tag.TagGPSAltitudeRef:
		r.Exif.GPS.altitudeRef = r.parseGPSRef(t)
	case tag.TagGPSLatitudeRef:
		r.Exif.GPS.latitudeRef = r.parseGPSRef(t)
	case tag.TagGPSLongitudeRef:
		r.Exif.GPS.longitudeRef = r.parseGPSRef(t)
	case tag.TagGPSDestLatitudeRef:
		r.Exif.GPS.destLatitudeRef = r.parseGPSRef(t)
	case tag.TagGPSDestLongitudeRef:
		r.Exif.GPS.destLongitudeRef = r.parseGPSRef(t)
	case tag.TagGPSSpeedRef:
		r.Exif.GPS.speedRef = r.parseGPSRef(t)
	case tag.TagGPSTrackRef:
		r.Exif.GPS.trackRef = r.parseGPSRef(t)
	case tag.TagGPSImgDirectionRef:
		r.Exif.GPS.imgDirectionRef = r.parseGPSRef(t)
	case tag.TagGPSDestBearingRef:
		r.Exif.GPS.destBearingRef = r.parseGPSRef(t)
	case tag.TagGPSDestDistanceRef:
		r.Exif.GPS.destDistanceRef = r.parseGPSRef(t)
	case tag.TagGPSAltitude:
		r.Exif.GPS.altitude = r.parseGPSAltitude(t)
	case tag.TagGPSLatitude:
		r.Exif.GPS.latitude = r.parseGPSCoord(t)
	case tag.TagGPSLongitude:
		r.Exif.GPS.longitude = r.parseGPSCoord(t)
	case tag.TagGPSDestLatitude:
		r.Exif.GPS.destLatitude = r.parseGPSCoord(t)
	case tag.TagGPSDestLongitude:
		r.Exif.GPS.destLongitude = r.parseGPSCoord(t)
	case tag.TagGPSSatellites:
		r.Exif.GPS.satellites = r.parseString(t)
	case tag.TagGPSStatus:
		r.Exif.GPS.status = r.parseString(t)
	case tag.TagGPSMeasureMode:
		r.Exif.GPS.measureMode = r.parseString(t)
	case tag.TagGPSMapDatum:
		r.Exif.GPS.mapDatum = r.parseString(t)
	case tag.TagGPSDOP:
		r.Exif.GPS.dop = r.parseRationalValue(t)
	case tag.TagGPSSpeed:
		r.Exif.GPS.speed = r.parseRationalValue(t)
	case tag.TagGPSTrack:
		r.Exif.GPS.track = r.parseRationalValue(t)
	case tag.TagGPSImgDirection:
		r.Exif.GPS.imgDirection = r.parseRationalValue(t)
	case tag.TagGPSDestBearing:
		r.Exif.GPS.destBearing = r.parseRationalValue(t)
	case tag.TagGPSDestDistance:
		r.Exif.GPS.destDistance = r.parseRationalValue(t)
	case tag.TagGPSHPositioningError:
		r.Exif.GPS.hPositioningError = r.parseRationalValue(t)
	case tag.TagGPSTimeStamp:
		r.Exif.GPS.setTime(r.parseGPSTimeStamp(t))
	case tag.TagGPSDateStamp:
		r.Exif.GPS.setDate(r.parseGPSDateStamp(t))
	default:
		return false
	}
	return true
}

// parseImageIFDTag parses the requested value from EXIF metadata.
func (r *Reader) parseImageIFDTag(t tag.Entry, dst *ImageIFD) bool {
	switch t.ID {
	case tag.TagSubfileType:
		dst.SubfileType = meta.SubfileType(r.parseUint32(t))
	case tag.TagBitsPerSample:
		// ingnore tag
	case tag.TagCompression:
		dst.Compression = meta.Compression(r.parseUint16(t))
	case tag.TagXResolution:
		dst.XResolution = r.parseRationalValue(t)
	case tag.TagYResolution:
		dst.YResolution = r.parseRationalValue(t)
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
		dst.ModifyDate = r.parseDate(t)
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
	switch t.Type {
	case tag.TypeASCII, tag.TypeASCIINoNul:
	default:
		return 0
	}
	if t.IsEmbedded() {
		t.EmbeddedValue(r.state.buf[:4])
		return uint16(parseStrUint(trimNULBuffer(r.state.buf[:4])))
	}
	buf, _, err := r.readTagBytes(t, 16)
	if err != nil {
		return 0
	}
	return uint16(parseStrUint(trimNULBuffer(buf)))
}

// apertureValueToFNumber converts APEX aperture values into F-number approximations.
func apertureValueToFNumber(v meta.Aperture) meta.Aperture {
	if v == 0 {
		return 0
	}
	return v
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
