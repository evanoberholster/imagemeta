package isobmff

import (
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	exiftag "github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// ReadMetadata advances metadata extraction by scanning top-level BMFF boxes.
//
// The reader can be called repeatedly; once configured goals are satisfied it
// returns io.EOF to signal completion.
func (r *Reader) ReadMetadata() (err error) {
	if r.goalsInitialized && r.stopAfterMetadata {
		return io.EOF
	}

	for {
		b, readErr := r.readBox()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return io.EOF
			}
			return fmt.Errorf("ReadMetadata: %w", readErr)
		}
		keepScanning, boxErr := r.readMetadataBox(&b)
		if boxErr != nil && logLevelError() {
			logError().Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int("size", b.size).Err(boxErr).Send()
		}
		if boxErr == nil && r.goalsInitialized && r.hasMetadataGoals() && r.metadataGoalsSatisfied() {
			r.stopAfterMetadata = true
		}
		if keepScanning {
			continue
		}
		return boxErr
	}
}

func (r *Reader) readMetadataBox(b *box) (keepScanning bool, err error) {
	switch b.boxType {
	case typeMdat:
		return false, r.readMdat(b)
	case typeExif:
		return false, r.readExif(b)
	case typeMeta:
		return false, r.readMeta(b)
	case typeMoov:
		return false, r.readMoovBox(b)
	case typeUUID:
		return false, r.readUUIDBox(b)
	case typeJXL, typeJumb, typeJxlc, typeJxll, typeJxlp:
		// JPEG XL container boxes can appear before metadata boxes.
		// Skip and continue scanning.
		err = b.close()
		return err == nil, err
	default:
		if logLevelInfo() {
			logInfo().Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int("size", b.size).Send()
		}
		return false, b.close()
	}
}

func (r *Reader) readMdat(b *box) (err error) {
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}

	type mdatItemKind uint8
	const (
		mdatItemExif mdatItemKind = iota + 1
		mdatItemXMP
	)
	type mdatItem struct {
		kind     mdatItemKind
		offset   offsetLength
		itemType boxType
	}

	var items [2]mdatItem
	itemCount := 0
	if r.hasGoal(metadataKindExif) && !r.hasHave(metadataKindExif) && r.heic.exif.ol.length > 0 {
		items[itemCount] = mdatItem{
			kind:     mdatItemExif,
			offset:   r.heic.exif.ol,
			itemType: typeExif,
		}
		itemCount++
	}
	if r.hasGoal(metadataKindXMP) && !r.hasHave(metadataKindXMP) && r.heic.xml.ol.length > 0 {
		items[itemCount] = mdatItem{
			kind:     mdatItemXMP,
			offset:   r.heic.xml.ol,
			itemType: typeUUID,
		}
		itemCount++
	}
	if itemCount == 0 {
		return b.close()
	}

	payloadStart := boxPayloadOffset(b)
	for i := 0; i < itemCount-1; i++ {
		for j := i + 1; j < itemCount; j++ {
			offI, errI := resolveMdatExtentOffset(payloadStart, items[i].offset.offset)
			offJ, errJ := resolveMdatExtentOffset(payloadStart, items[j].offset.offset)
			if errI == nil && errJ == nil && offJ < offI {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	imageType := r.metadataImageType()
	for i := 0; i < itemCount; i++ {
		inner, openErr := newMdatExtentBox(b, payloadStart, items[i].offset, items[i].itemType)
		if openErr != nil {
			if logLevelDebug() {
				logDebug().Object("box", b).Err(openErr).Uint64("offset", items[i].offset.offset).Int("length", items[i].offset.length).Msg("skip unresolved mdat extent")
			}
			continue
		}

		switch items[i].kind {
		case mdatItemExif:
			if err = seekExifTIFFHeader(&inner); err != nil {
				if logLevelDebug() {
					logDebug().Object("box", inner).Err(err).Msg("skip non-TIFF mdat exif candidate")
				}
				if closeErr := inner.close(); closeErr != nil {
					return closeErr
				}
				continue
			}

			header, headerErr := readExifHeader(&inner, exiftag.IFD0, imageType)
			if headerErr != nil {
				if logLevelDebug() {
					logDebug().Object("box", inner).Err(headerErr).Msg("skip invalid mdat exif header")
				}
				if closeErr := inner.close(); closeErr != nil {
					return closeErr
				}
				continue
			}
			if err = r.callExifReader(&inner, header); err != nil {
				return err
			}
		case mdatItemXMP:
			header, headerErr := evaluateXPacketHeader(&inner)
			if headerErr != nil {
				if closeErr := inner.close(); closeErr != nil {
					return closeErr
				}
				return headerErr
			}
			if err = r.callXMPReader(&inner, header); err != nil {
				return err
			}
		}

		if logLevelInfo() {
			logInfo().Object("box", inner).Int("remain", inner.remain).Send()
		}
		if closeErr := inner.close(); closeErr != nil {
			return closeErr
		}
	}
	return b.close()
}

// readExif parses a top-level Exif box payload and streams it to the Exif callback.
func (r *Reader) readExif(b *box) (err error) {
	if !b.isType(typeExif) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}

	if err = seekExifTIFFHeader(b); err != nil {
		return fmt.Errorf("readExif: %w", err)
	}

	header, err := readExifHeader(b, exiftag.IFD0, r.metadataImageType())
	if err != nil {
		return fmt.Errorf("readExif: %w", err)
	}

	if err = r.callExifReader(b, header); err != nil {
		return err
	}

	return b.close()
}

func boxPayloadOffset(b *box) uint64 {
	payloadOffset := b.offset + int64(b.size-b.remain)
	if payloadOffset <= 0 {
		return 0
	}
	return uint64(payloadOffset)
}

func resolveMdatExtentOffset(payloadStart, rawOffset uint64) (uint64, error) {
	if rawOffset >= payloadStart {
		return rawOffset, nil
	}
	if rawOffset > ^uint64(0)-payloadStart {
		return 0, ErrBufLength
	}
	return payloadStart + rawOffset, nil
}

func newMdatExtentBox(b *box, payloadStart uint64, ol offsetLength, innerType boxType) (inner box, err error) {
	if ol.length == 0 {
		return inner, ErrBufLength
	}

	targetOffset, err := resolveMdatExtentOffset(payloadStart, ol.offset)
	if err != nil {
		return inner, err
	}
	currentOffset := boxPayloadOffset(b)
	if targetOffset < currentOffset {
		return inner, ErrRemainLengthInsufficient
	}
	discardBytes, err := uint64ToInt(targetOffset - currentOffset)
	if err != nil {
		return inner, err
	}
	if discardBytes > b.remain {
		return inner, ErrRemainLengthInsufficient
	}
	if err = discardBoxBytes(b, discardBytes); err != nil {
		return inner, err
	}

	if ol.length > b.remain {
		return inner, ErrRemainLengthInsufficient
	}
	targetOffset64, err := uint64ToInt64(targetOffset)
	if err != nil {
		return inner, err
	}

	inner = box{
		reader:  b.reader,
		outer:   b,
		boxType: innerType,
		offset:  targetOffset64,
		size:    ol.length,
		remain:  ol.length,
	}
	return inner, nil
}

// readExifHeader parses byte-order and IFD0 offset from the TIFF header prefix.
func readExifHeader(b *box, firstIfd exiftag.IfdType, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	buf, err := b.Peek(8)
	if err != nil {
		err = fmt.Errorf("readExifHeader: %w", err)
		return
	}
	endian := utils.BinaryOrder(buf[:4])
	if endian == utils.UnknownEndian {
		return header, ErrBufLength
	}
	header = meta.NewExifHeader(endian, endian.Uint32(buf[4:8]), 0, clampIntToUint32(b.remain), it)
	header.FirstIfd = firstIfd
	if logLevelInfo() {
		logInfo().Object("box", b).Object("header", header).Send()
	}
	_, err = b.Discard(8)
	return header, err
}

// seekExifTIFFHeader advances through common Exif wrappers until TIFF bytes.
func seekExifTIFFHeader(b *box) error {
	for {
		if b.remain < 8 {
			return ErrBufLength
		}
		buf, err := b.Peek(8)
		if err != nil {
			return err
		}

		// Some Exif payloads include the APP1 prefix before TIFF data.
		if hasExifAPP1Prefix(buf) {
			if _, err = b.Discard(6); err != nil {
				return err
			}
			continue
		}

		if hasTIFFHeader(buf[:4]) {
			return nil
		}

		// HEIF/JXL style Exif payloads can start with a 4-byte TIFF header offset.
		if ok, offsetErr := discardTIFFOffsetPrefix(b, bmffEndian.Uint32(buf[:4])); offsetErr != nil {
			return offsetErr
		} else if ok {
			return nil
		}

		// Some payloads are wrapped in a local Exif box header.
		if hasEmbeddedExifBoxHeader(buf) {
			boxSize := bmffEndian.Uint32(buf[:4])
			if boxSize < 8 {
				return ErrBufLength
			}
			boxSizeInt, convErr := uint64ToInt(uint64(boxSize))
			if convErr != nil || boxSizeInt > b.remain {
				return ErrBufLength
			}
			if _, err = b.Discard(8); err != nil {
				return err
			}
			continue
		}

		if ok, scanErr := seekTIFFHeaderByScan(b, 64); scanErr != nil {
			return scanErr
		} else if ok {
			return nil
		}
		return ErrBufLength
	}
}

func discardTIFFOffsetPrefix(b *box, offset uint32) (bool, error) {
	skip := uint64(4) + uint64(offset)
	need := skip + 4
	if need > uint64(b.remain) {
		return false, nil
	}
	needInt, err := uint64ToInt(need)
	if err != nil {
		return false, nil
	}
	skipInt, err := uint64ToInt(skip)
	if err != nil {
		return false, nil
	}
	buf, err := b.Peek(needInt)
	if err != nil {
		return false, err
	}
	if !hasTIFFHeader(buf[skipInt : skipInt+4]) {
		return false, nil
	}
	_, err = b.Discard(skipInt)
	if err != nil {
		return false, err
	}
	return true, nil
}

func seekTIFFHeaderByScan(b *box, maxScan int) (bool, error) {
	if b.remain < 4 || maxScan <= 0 {
		return false, nil
	}
	limit := maxScan
	if limit > b.remain-4 {
		limit = b.remain - 4
	}
	if limit <= 0 {
		return false, nil
	}

	buf, err := b.Peek(limit + 4)
	if err != nil {
		return false, err
	}
	for i := 1; i <= limit; i++ {
		if !hasTIFFHeader(buf[i : i+4]) {
			continue
		}
		_, err = b.Discard(i)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func hasEmbeddedExifBoxHeader(buf []byte) bool {
	return len(buf) >= 8 &&
		buf[4] == 'E' &&
		buf[5] == 'x' &&
		buf[6] == 'i' &&
		buf[7] == 'f'
}

func hasExifAPP1Prefix(buf []byte) bool {
	return len(buf) >= 6 &&
		buf[0] == 'E' &&
		buf[1] == 'x' &&
		buf[2] == 'i' &&
		buf[3] == 'f' &&
		buf[4] == 0x00 &&
		buf[5] == 0x00
}

func hasTIFFHeader(buf []byte) bool {
	return len(buf) >= 4 &&
		((buf[0] == 'I' && buf[1] == 'I' && buf[2] == 0x2A && buf[3] == 0x00) ||
			(buf[0] == 'M' && buf[1] == 'M' && buf[2] == 0x00 && buf[3] == 0x2A))
}

func (r *Reader) metadataImageType() imagetype.ImageType {
	if it, ok := imageTypeFromBrand(r.ftyp.MajorBrand); ok {
		return it
	}
	for _, compatibleBrand := range r.ftyp.Compatible {
		if it, ok := imageTypeFromBrand(compatibleBrand); ok {
			return it
		}
	}
	return imagetype.ImageHEIF
}

func imageTypeFromBrand(brand brand) (imagetype.ImageType, bool) {
	switch brand {
	case brandJxl:
		return imagetype.ImageJXL, true
	case brandAvif, brandAvis:
		return imagetype.ImageAVIF, true
	case brandHeic, brandHeim, brandHeis, brandHeix, brandHevc, brandHevm, brandHevs, brandHevx:
		return imagetype.ImageHEIC, true
	case brandHeif, brandMiaf, brandMif1, brandMif2, brandMsf1:
		return imagetype.ImageHEIF, true
	case brandCrx:
		return imagetype.ImageCR3, true
	default:
		return imagetype.ImageHEIF, false
	}
}

// readMeta parses HEIF/JXL/CR3 metadata containers and records the item
// metadata needed to locate Exif/XMP payloads in mdat.
func (r *Reader) readMeta(b *box) (err error) {
	if !b.isType(typeMeta) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	if err = b.readFlags(); err != nil {
		return err
	}
	parseCR3ItemGraph := r.ftyp.MajorBrand == brandCrx
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	err = readContainerBoxes(b, func(inner *box) error {
		switch inner.boxType {
		case typeUUID:
			return r.readUUIDBox(inner)
		case typeHdlr:
			if parseCR3ItemGraph {
				_, err = readHdlr(inner)
				return err
			}
			return nil
		case typePitm:
			if parseCR3ItemGraph {
				r.heic.pitm, err = readPitm(inner)
				return err
			}
			return nil
		case typeIinf:
			return r.readIinf(inner)
		case typeIref:
			if parseCR3ItemGraph {
				return r.readIref(inner)
			}
			return nil
		case typeIprp:
			if parseCR3ItemGraph {
				return r.readIprp(inner)
			}
			return nil
		case typeIdat:
			r.heic.idatData = offsetLength{
				offset: boxPayloadOffset(inner),
				length: inner.remain,
			}
			return nil
		case typeIloc:
			return r.readIloc(inner)
		default:
			if logLevelInfo() {
				logInfo().Str("boxType", inner.boxType.String()).Int64("offset", inner.offset).Int("size", inner.size).Send()
			}
			return nil
		}
	})
	return err
}

// readMoovBox reads an 'moov' box from a BMFF file.
func (r *Reader) readMoovBox(b *box) (err error) {
	if !b.isType(typeMoov) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	err = readContainerBoxes(b, func(inner *box) error {
		switch inner.boxType {
		case typeUUID:
			return r.readUUIDBox(inner)
		case typeTrak:
			return readCrxTrakBox(inner)
		default:
			if logLevelInfo() {
				logInfo().Str("boxType", inner.boxType.String()).Int64("offset", inner.offset).Int("size", inner.size).Send()
			}
			return nil
		}
	})
	return err
}

// finalizeInnerBox applies common child-box lifecycle handling:
// propagate parse errors, then close the child box payload.
func finalizeInnerBox(inner *box, parseErr error) error {
	if parseErr != nil {
		if logLevelError() && inner != nil {
			logError().Str("boxType", inner.boxType.String()).Int64("offset", inner.offset).Int("size", inner.size).Err(parseErr).Send()
		}
		return parseErr
	}
	if inner == nil {
		return nil
	}
	if err := inner.close(); err != nil {
		if logLevelError() {
			logError().Str("boxType", inner.boxType.String()).Int64("offset", inner.offset).Int("size", inner.size).Err(err).Send()
		}
		return err
	}
	return nil
}

func handleCallbackError(b *box, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	if logLevelError() {
		if b == nil {
			logError().Err(err).Send()
		} else {
			logError().Object("box", b).Err(err).Send()
		}
	}
	return nil
}

func clampIntToUint32(v int) uint32 {
	if v <= 0 {
		return 0
	}
	if uint64(v) > uint64(^uint32(0)) {
		return ^uint32(0)
	}
	//nolint:gosec // G115: value is clamped to uint32 bounds above.
	return uint32(v)
}

func uint64ToInt64(v uint64) (int64, error) {
	if v > uint64(maxIntValue) {
		return 0, errLargeBox
	}
	return int64(v), nil
}

func uint64ToInt(v uint64) (int, error) {
	if v > uint64(maxIntValue) {
		return 0, errLargeBox
	}
	return int(v), nil
}
