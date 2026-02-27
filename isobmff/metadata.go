package isobmff

import (
	"errors"
	"fmt"
	"io"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
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
			if readErr == io.EOF {
				return io.EOF
			}
			return fmt.Errorf("ReadMetadata: %w", readErr)
		}
		switch b.boxType {
		case typeMdat:
			err = r.readMdat(&b)
		case typeExif:
			err = r.readExif(&b)
		case typeMeta:
			err = r.readMeta(&b)
		case typeMoov:
			err = r.readMoovBox(&b)
		case typeUUID:
			err = r.readUUIDBox(&b)
		case typeJXL, typeJumb, typeJxlc, typeJxll, typeJxlp:
			// JPEG XL container boxes can appear before metadata boxes.
			// Skip and continue scanning.
			err = b.close()
			if err == nil {
				continue
			}
		default:
			if logLevelInfo() {
				logInfo().Str("boxType", b.boxType.String()).Int("offset", b.offset).Int64("size", b.size).Send()
			}
			err = b.close()
		}
		if err != nil && logLevelError() {
			logError().Str("boxType", b.boxType.String()).Int("offset", b.offset).Int64("size", b.size).Err(err).Send()
		}
		if err == nil && r.goalsInitialized && r.hasMetadataGoals() && r.metadataGoalsSatisfied() {
			r.stopAfterMetadata = true
		}
		return err
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
		mdatItemPreview
	)
	type mdatItem struct {
		kind          mdatItemKind
		offset        offsetLength
		itemType      boxType
		previewHeader meta.PreviewHeader
	}

	var items [3]mdatItem
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
	if header, ol, ok := r.selectMdatPreviewCandidate(); ok {
		items[itemCount] = mdatItem{
			kind:          mdatItemPreview,
			offset:        ol,
			itemType:      typeMdat,
			previewHeader: header,
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
				logDebug().Object("box", b).Err(openErr).Uint64("offset", items[i].offset.offset).Uint64("length", items[i].offset.length).Msg("skip unresolved mdat extent")
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

			header, headerErr := readExifHeader(&inner, ifds.IFD0, imageType)
			if headerErr != nil {
				if logLevelDebug() {
					logDebug().Object("box", inner).Err(headerErr).Msg("skip invalid mdat exif header")
				}
				if closeErr := inner.close(); closeErr != nil {
					return closeErr
				}
				continue
			}
			if r.exifReader != nil {
				callbackErr := r.exifReader(newLimitedReader(&inner, inner.remain), header)
				if callbackErr != nil {
					if err = handleCallbackError(&inner, callbackErr); err != nil {
						return err
					}
				} else {
					r.setHave(metadataKindExif, true)
				}
			}
		case mdatItemXMP:
			if r.xmpReader != nil {
				header, headerErr := evaluateXPacketHeader(&inner)
				if headerErr != nil {
					if closeErr := inner.close(); closeErr != nil {
						return closeErr
					}
					return headerErr
				}
				callbackErr := r.xmpReader(newLimitedReader(&inner, inner.remain), header)
				if callbackErr != nil {
					if err = handleCallbackError(&inner, callbackErr); err != nil {
						return err
					}
				} else {
					r.setHave(metadataKindXMP, true)
				}
			}
		case mdatItemPreview:
			if r.previewImageReader != nil {
				callbackErr := r.previewImageReader(newLimitedReader(&inner, inner.remain), items[i].previewHeader)
				if callbackErr != nil {
					if err = handleCallbackError(&inner, callbackErr); err != nil {
						return err
					}
				} else {
					r.setHave(metadataKindPRVW, true)
				}
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

func (r *Reader) selectMdatPreviewCandidate() (meta.PreviewHeader, offsetLength, bool) {
	if r.ftyp.MajorBrand == brandCrx || !r.hasGoal(metadataKindPRVW) || r.hasHave(metadataKindPRVW) || r.previewImageReader == nil {
		return meta.PreviewHeader{}, offsetLength{}, false
	}

	bestID := invalidItemID
	bestRank := 99
	bestLength := uint64(^uint64(0))
	bestOffset := offsetLength{}

	consider := func(id itemID, rank int) {
		if id == invalidItemID || id == r.heic.exif.id || id == r.heic.xml.id {
			return
		}
		ol, ok := r.lookupItemLocation(id)
		if !ok || ol.length == 0 {
			return
		}
		if !r.itemIsPreviewEligible(id) {
			return
		}
		if rank > bestRank {
			return
		}
		if rank == bestRank && ol.length >= bestLength {
			return
		}
		bestID = id
		bestRank = rank
		bestLength = ol.length
		bestOffset = ol
	}

	for i := range r.heic.references {
		ref := r.heic.references[i]
		if ref.referenceType == typeThmb && r.heic.pitm != invalidItemID && ref.toID == r.heic.pitm {
			consider(ref.fromID, 0)
		}
	}
	for i := range r.heic.references {
		ref := r.heic.references[i]
		if ref.referenceType == typeThmb {
			consider(ref.fromID, 1)
		}
	}
	for i := range r.heic.locations {
		id := r.heic.locations[i].id
		if r.heic.pitm != invalidItemID && id == r.heic.pitm {
			continue
		}
		consider(id, 2)
	}
	if r.heic.pitm != invalidItemID {
		consider(r.heic.pitm, 3)
	}

	if bestID == invalidItemID {
		return meta.PreviewHeader{}, offsetLength{}, false
	}

	imageType := r.previewItemImageType(bestID)
	width, height := r.itemPreviewDimensions(bestID)
	size := uint32(bestOffset.length)
	if bestOffset.length > uint64(^uint32(0)) {
		size = ^uint32(0)
	}
	header := meta.PreviewHeader{
		Size:      size,
		Width:     width,
		Height:    height,
		ImageType: imageType,
		Source:    meta.PreviewSourcePRVW,
	}
	return header, bestOffset, true
}

func (r *Reader) itemIsPreviewEligible(id itemID) bool {
	info, ok := r.lookupItemInfo(id)
	if !ok {
		return true
	}
	switch info.itemType {
	case itemTypeExif, itemTypeURI:
		return false
	case itemTypeMime:
		if info.mimeType == "" || isXMPMIMEType(info.mimeType) {
			return false
		}
		return asciiContainsFold(info.mimeType, "image/")
	case itemTypeHvc1, itemTypeAv01:
		return true
	default:
		return false
	}
}

func (r *Reader) previewItemImageType(id itemID) imagetype.ImageType {
	info, ok := r.lookupItemInfo(id)
	if !ok {
		return r.metadataImageType()
	}
	switch info.itemType {
	case itemTypeHvc1:
		return imagetype.ImageHEIC
	case itemTypeAv01:
		return imagetype.ImageAVIF
	case itemTypeMime:
		if asciiContainsFold(info.mimeType, "jpeg") || asciiContainsFold(info.mimeType, "jpg") {
			return imagetype.ImageJPEG
		}
		if asciiContainsFold(info.mimeType, "avif") {
			return imagetype.ImageAVIF
		}
		if asciiContainsFold(info.mimeType, "heic") || asciiContainsFold(info.mimeType, "heif") {
			return imagetype.ImageHEIC
		}
	}
	return r.metadataImageType()
}

func (r *Reader) itemPreviewDimensions(id itemID) (uint16, uint16) {
	for i := range r.heic.propertyLinks {
		link := r.heic.propertyLinks[i]
		if link.itemID != id || link.propertyIndex == 0 {
			continue
		}
		idx := int(link.propertyIndex - 1)
		if idx < 0 || idx >= len(r.heic.properties) {
			continue
		}
		prop := r.heic.properties[idx]
		if prop.boxType != typeIspe {
			continue
		}
		w := prop.width
		h := prop.height
		if w > uint32(^uint16(0)) {
			w = uint32(^uint16(0))
		}
		if h > uint32(^uint16(0)) {
			h = uint32(^uint16(0))
		}
		return uint16(w), uint16(h)
	}
	return 0, 0
}

// readExif parses a top-level Exif box payload and streams it to the Exif callback.
func (r *Reader) readExif(b *box) (err error) {
	if !b.isType(typeExif) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}

	if err = seekExifTIFFHeader(b); err != nil {
		return fmt.Errorf("readExif: %w", err)
	}

	header, err := readExifHeader(b, ifds.IFD0, r.metadataImageType())
	if err != nil {
		return fmt.Errorf("readExif: %w", err)
	}

	if r.exifReader != nil {
		callbackErr := r.exifReader(newLimitedReader(b, b.remain), header)
		if callbackErr != nil {
			if err = handleCallbackError(b, callbackErr); err != nil {
				return err
			}
		} else {
			r.setHave(metadataKindExif, true)
		}
	}

	return b.close()
}

// newExifBox creates a bounded view into an mdat payload using iloc offsets.
func (r *Reader) newExifBox(b *box) (inner box, err error) {
	return newMdatExtentBox(b, boxPayloadOffset(b), r.heic.exif.ol, typeExif)
}

func boxPayloadOffset(b *box) uint64 {
	return uint64(b.offset + int(b.size) - b.remain)
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
	discardBytes := targetOffset - currentOffset
	if discardBytes > uint64(b.remain) {
		return inner, ErrRemainLengthInsufficient
	}
	if _, err = b.Discard(int(discardBytes)); err != nil {
		return inner, err
	}

	if ol.length > uint64(b.remain) {
		return inner, ErrRemainLengthInsufficient
	}
	maxInt := int64(^uint(0) >> 1)
	if ol.length > uint64(maxInt) {
		return inner, errLargeBox
	}

	size := int64(ol.length)

	inner = box{
		reader:  b.reader,
		outer:   b,
		boxType: innerType,
		offset:  int(targetOffset),
		size:    size,
		remain:  int(size),
	}
	return inner, nil
}

// readExifHeader parses byte-order and IFD0 offset from the TIFF header prefix.
func readExifHeader(b *box, firstIfd ifds.IfdType, it imagetype.ImageType) (header meta.ExifHeader, err error) {
	buf, err := b.Peek(16)
	if err != nil {
		err = fmt.Errorf("readExifHeader: %w", err)
		return
	}
	endian := utils.BinaryOrder(buf[:4])
	if endian == utils.UnknownEndian {
		return header, ErrBufLength
	}
	header = meta.NewExifHeader(endian, endian.Uint32(buf[4:8]), 0, uint32(b.remain), it)
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

		// Some payloads are wrapped in a local Exif box header.
		if hasEmbeddedExifBoxHeader(buf) {
			boxSize := bmffEndian.Uint32(buf[:4])
			if boxSize < 8 || int(boxSize) > b.remain {
				return ErrBufLength
			}
			if _, err = b.Discard(8); err != nil {
				return err
			}
			continue
		}

		if hasTIFFHeader(buf[:4]) {
			return nil
		}

		// HEIF/JXL style Exif payloads can start with a 4-byte TIFF header offset.
		offset := int(bmffEndian.Uint32(buf[:4]))
		if offset < 0 || offset > b.remain-8 {
			return nil
		}
		_, err = b.Discard(4 + offset)
		return err
	}
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

// readMeta parses HEIF/JXL metadata containers and records item references
// needed to locate Exif/XMP payloads in mdat.
func (r *Reader) readMeta(b *box) (err error) {
	if !b.isType(typeMeta) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	if err = b.readFlags(); err != nil {
		return err
	}
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeUUID:
			err = r.readUUIDBox(&inner)
		case typeHdlr:
			_, err = readHdlr(&inner)
		case typePitm:
			r.heic.pitm, err = readPitm(&inner)
		case typeIinf:
			err = r.readIinf(&inner)
		case typeIref:
			err = r.readIref(&inner)
		case typeIprp:
			err = r.readIprp(&inner)
		case typeIdat:
			r.heic.idatData = offsetLength{
				offset: boxPayloadOffset(&inner),
				length: uint64(inner.remain),
			}
			if r.ftyp.MajorBrand == brandCrx {
				r.heic.idat, err = readIdat(&inner)
			}
		case typeIloc:
			err = r.readIloc(&inner)
		default:
			if logLevelInfo() {
				logInfo().Str("boxType", inner.boxType.String()).Int("offset", inner.offset).Int64("size", inner.size).Send()
			}
		}
		if err = finalizeInnerBox(&inner, err); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return b.close()
}

// readMoovBox reads an 'moov' box from a BMFF file.
func (r *Reader) readMoovBox(b *box) (err error) {
	if !b.isType(typeMoov) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
		switch inner.boxType {
		case typeUUID:
			err = r.readUUIDBox(&inner)
		case typeTrak:
			err = readCrxTrakBox(&inner)
		default:
			if logLevelInfo() {
				logInfo().Str("boxType", inner.boxType.String()).Int("offset", inner.offset).Int64("size", inner.size).Send()
			}
		}
		if err = finalizeInnerBox(&inner, err); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return b.close()
}

// finalizeInnerBox applies common child-box lifecycle handling:
// propagate parse errors, then close the child box payload.
func finalizeInnerBox(inner *box, parseErr error) error {
	if parseErr != nil {
		if logLevelError() && inner != nil {
			logError().Str("boxType", inner.boxType.String()).Int("offset", inner.offset).Int64("size", inner.size).Err(parseErr).Send()
		}
		return parseErr
	}
	if inner == nil {
		return nil
	}
	if err := inner.close(); err != nil {
		if logLevelError() {
			logError().Str("boxType", inner.boxType.String()).Int("offset", inner.offset).Int64("size", inner.size).Err(err).Send()
		}
		return err
	}
	return nil
}

func newLimitedReader(r io.Reader, limit int) io.Reader {
	if limit <= 0 {
		return io.LimitReader(r, 0)
	}
	return &io.LimitedReader{R: r, N: int64(limit)}
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
