package isobmff

import (
	"fmt"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
)

// Canon CR3-specific UUID box identifiers.
//
// CR3 stores key metadata and preview payloads in UUID boxes rather than only
// in generic HEIF item containers.
var (
	// cr3MetaBoxUUID is the uuid that corresponds with Canon CR3 Metadata.
	cr3MetaBoxUUID = meta.UUIDFromString("85c0b687-820f-11e0-8111-f4ce462b6a48")

	// cr3XPacketUUID is the uuid that corresponds with Canon CR3 xpacket data
	cr3XPacketUUID = meta.UUIDFromString("be7acfcb-97a9-42e8-9c71-999491e3afac")

	// cr3PreviewUUID is the uuid that corresponds with Canon CR3 Preview Image.
	cr3PreviewUUID = meta.UUIDFromString("eaf42b5e-1c98-4b88-b9fb-b7dc406e4d16")
)

// readUUIDBox routes Canon UUID payloads to CR3-specific readers.
//
// This handles:
// - XPacket/XMP payloads
// - Canon metadata container (CMT*/CTBO/CCTP/THMB)
// - PRVW JPEG preview payload
func (r *Reader) readUUIDBox(b *box) error {
	if !b.isType(typeUUID) {
		return fmt.Errorf("Box %s: %w", b.boxType, ErrWrongBoxType)
	}
	uuid, err := b.readUUID()
	if err != nil {
		return err
	}
	if logLevelInfo() {
		logInfoBox(b).Str("uuid", uuid.String()).Send()
	}
	switch uuid {
	case cr3XPacketUUID:
		header, evalErr := evaluateXPacketHeader(b)
		if evalErr != nil {
			return evalErr
		}
		if logLevelInfo() {
			logInfoBox(b).
				Bool("hasXPacketPI", header.HasXPacketPI).
				Bool("hasXMPMeta", header.HasXMPMeta).
				Uint32("xpacketLength", header.Length).
				Send()
		}
		if r.xmpReader != nil {
			callbackErr := r.xmpReader(newLimitedReader(b, b.remain), header)
			if callbackErr != nil {
				if err = handleCallbackError(b, callbackErr); err != nil {
					return err
				}
			} else {
				r.setHave(metadataKindXMP, true)
			}
		}
	case cr3MetaBoxUUID:
		if err = r.readCrxMoovBox(b); err != nil {
			return err
		}
	case cr3PreviewUUID:
		if err = r.readPreview(b); err != nil {
			return err
		}
	default:
		if logLevelDebug() {
			logDebug().Object("box", b).Send()
		}
	}
	return b.close()
}

// readCrxMoovBox parses Canon's metadata UUID payload and dispatches its
// proprietary child boxes.
func (r *Reader) readCrxMoovBox(b *box) (err error) {
	sawTHMB := false
	var inner box
	var ok bool
	for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {

		switch inner.boxType {
		case typeCCTP:
			err = readCCTPBox(&inner)
		case typeCNCV:
			err = readCNCVBox(&inner)
		case typeCTBO:
			err = readCTBOBox(&inner)
		case typeCMT1:
			err = r.readCMTBox(&inner, ifds.IFD0)
		case typeCMT2:
			err = r.readCMTBox(&inner, ifds.ExifIFD)
		case typeCMT3:
			err = r.readCMTBox(&inner, ifds.MknoteIFD)
		case typeCMT4:
			err = r.readCMTBox(&inner, ifds.GPSIFD)
		case typeTHMB, typeThmb:
			sawTHMB = true
			err = r.readTHMBBox(&inner)
		default:
			if logLevelDebug() {
				logDebug().Str("boxType", inner.boxType.String()).Int("offset", inner.offset).Int64("size", inner.size).Send()
			}
		}
		if err = finalizeInnerBox(&inner, err); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	if r.hasGoal(metadataKindTHMB) && !sawTHMB {
		// Some CR3 variants do not include a THMB box.
		r.setGoal(metadataKindTHMB, false)
	}
	return b.close()
}

// readCMTBox reads an ISOBMFF Box "CMT1","CMT2","CMT3", or "CMT4" from CR3.
//
// Each CMT box contains TIFF/Exif-like IFD data for a specific IFD family.
func (r *Reader) readCMTBox(b *box, ifdType ifds.IfdType) (err error) {
	if r.exifReader == nil {
		return nil
	}
	header, err := readExifHeader(b, ifdType, imagetype.ImageCR3)
	if err != nil {
		return err
	}
	callbackErr := r.exifReader(newLimitedReader(b, b.remain), header)
	if callbackErr != nil {
		return handleCallbackError(b, callbackErr)
	}
	r.setHave(metadataKindExif, true)
	return nil
}

// cncvBox is Canon's compressor/codec version payload.
type cncvBox struct {
	version [30]byte
}

// readCNCVBox reads and logs Canon codec version bytes.
// This box is informational and not required for Exif/XMP extraction.
func readCNCVBox(b *box) (err error) {
	if !b.isType(typeCNCV) {
		return ErrWrongBoxType
	}
	if !logLevelInfo() {
		return nil
	}
	var cncv cncvBox
	if b.remain < len(cncv.version) {
		return ErrBufLength
	}
	buf, err := b.Peek(30)
	if err != nil {
		return err
	}
	copy(cncv.version[:], buf[:30])
	logInfo().Object("box", b).Str("CNCV", string(cncv.version[:])).Send()
	return nil
}

// cctpBox describes Canon track metadata entries.
type cctpBox struct {
	count   uint32
	entries []cctpEntry
}

// cctpEntry is one CCTP track descriptor record.
type cctpEntry struct {
	size      uint32
	trackType uint32
	mediaType uint32
	unknown   uint32
	index     uint32
}

// readCCTPBox reads Canon CCTP entries for diagnostics/logging.
func readCCTPBox(b *box) (err error) {
	if !b.isType(typeCCTP) {
		return ErrWrongBoxType
	}
	if !logLevelInfo() {
		return nil
	}
	var cctp cctpBox
	if err = b.readFlags(); err != nil {
		return err
	}
	if cctp.count, err = b.readUint32(); err != nil {
		return err
	}
	entryCount := int(cctp.count)
	maxEntries := b.remain / 24
	if entryCount > maxEntries {
		entryCount = maxEntries
	}
	cctp.entries = make([]cctpEntry, 0, entryCount)
	for i := 0; i < entryCount; i++ {
		var ent cctpEntry
		if ent.size, err = b.readUint32(); err != nil {
			return err
		}
		if ent.trackType, err = b.readFourCC(); err != nil {
			return err
		}
		if ent.mediaType, err = b.readUint32(); err != nil {
			return err
		}
		if ent.unknown, err = b.readUint32(); err != nil {
			return err
		}
		if ent.index, err = b.readUint32(); err != nil {
			return err
		}
		cctp.entries = append(cctp.entries, ent)
	}
	logInfo().Object("box", b).Array("tracks", cctp).Send()
	return nil
}

// readCTBOBox reads Canon track base-offset records.
// Offsets are used for Canon-internal metadata organization.
func readCTBOBox(b *box) (err error) {
	if !b.isType(typeCTBO) {
		return ErrWrongBoxType
	}
	if !logLevelInfo() {
		return nil
	}
	var ctbo ctboBox
	if ctbo.count, err = b.readUint32(); err != nil {
		return err
	}
	itemCount := int(ctbo.count)
	maxItems := b.remain / 20
	if itemCount > maxItems {
		itemCount = maxItems
	}
	ctbo.items = make([]offsetLength, 0, itemCount)
	for i := 0; i < itemCount; i++ {
		_, readErr := b.readUint32() // item index (1-based)
		if readErr != nil {
			return readErr
		}
		var ent offsetLength
		if ent.offset, err = b.readUintN(8); err != nil {
			return err
		}
		if ent.length, err = b.readUintN(8); err != nil {
			return err
		}
		ctbo.items = append(ctbo.items, ent)
	}
	if len(ctbo.items) > int(ctbo.count) {
		ctbo.items = ctbo.items[:ctbo.count]
	}
	logInfo().Object("box", b).Array("items", ctbo).Send()
	return nil
}

// ctboBox is a Canon tracks base offsets box.
type ctboBox struct {
	items []offsetLength
	count uint32
}

// thmbBox is Canon thumbnail metadata (JPEG dimensions and payload size).
type thmbBox struct {
	width  uint16
	height uint16
	size   uint32
}

// readTHMBBox extracts the THMB JPEG payload and streams it to the preview callback.
func (r *Reader) readTHMBBox(b *box) (err error) {
	if !b.isType(typeTHMB) && !b.isType(typeThmb) {
		return ErrWrongBoxType
	}
	if r.previewImageReader == nil {
		return nil
	}
	var thmb thmbBox
	thmb, err = parseTHMBBox(b)
	if err != nil {
		return err
	}
	if thmb.size > uint32(b.remain) {
		return ErrRemainLengthInsufficient
	}

	if thmb.size > 0 {
		payloadOffset := b.offset + int(b.size) - b.remain
		payload := box{
			reader:  b.reader,
			outer:   b,
			boxType: b.boxType,
			offset:  payloadOffset,
			size:    int64(thmb.size),
			remain:  int(thmb.size),
		}
		header := meta.PreviewHeader{
			Size:      thmb.size,
			Width:     thmb.width,
			Height:    thmb.height,
			ImageType: imagetype.ImageJPEG,
			Source:    meta.PreviewSourceTHMB,
		}
		callbackErr := r.previewImageReader(newLimitedReader(&payload, payload.remain), header)
		if callbackErr != nil {
			if err = handleCallbackError(&payload, callbackErr); err != nil {
				return err
			}
		} else {
			r.setHave(metadataKindTHMB, true)
		}
		if closeErr := payload.close(); closeErr != nil {
			return closeErr
		}
	}

	return nil
}

// parseTHMBBox parses THMB header fields and advances to the payload.
func parseTHMBBox(b *box) (thmb thmbBox, err error) {
	if b.remain < 16 {
		return thmb, ErrBufLength
	}
	buf, err := b.Peek(16)
	if err != nil {
		return thmb, err
	}

	thmb.width = bmffEndian.Uint16(buf[4:6])
	thmb.height = bmffEndian.Uint16(buf[6:8])
	thmb.size = bmffEndian.Uint32(buf[8:12])

	if _, err = b.Discard(16); err != nil {
		return thmb, err
	}
	return thmb, nil
}

// fourCCFromUint32 converts a packed FourCC value to bytes.
func fourCCFromUint32(v uint32) [4]byte {
	return [4]byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

// fourCCString converts a packed FourCC to a string for logs/debugging.
func fourCCString(v uint32) string {
	buf := fourCCFromUint32(v)
	return string(buf[:])
}

// readCrxTrakBox is a placeholder for optional CR3 /trak parsing.
// The current metadata pipeline intentionally avoids track/sample-table parsing.
func readCrxTrakBox(b *box) (err error) {
	//var inner box
	//var ok bool
	////for inner, ok, err = b.readInnerBox(); err == nil && ok; inner, ok, err = b.readInnerBox() {
	//	switch inner.boxType {
	//	case typeMdia:
	//		err = readCrxMdia(&inner)
	//	case typeHdlr:
	//
	//	case typeStsd:
	//
	//	case typeStsz:
	//
	//	case typeCo64:
	//
	//	}
	//	if logLevelInfo() {
	//		logInfoBox(inner)
	//	}
	//	inner.close()
	//}
	if logLevelInfo() {
		logInfo().Object("box", b).Send()
	}
	return nil
}
