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
	// cr3MetaBoxUUID identifies Canon's CR3 metadata UUID payload.
	cr3MetaBoxUUID = meta.UUIDFromString("85c0b687-820f-11e0-8111-f4ce462b6a48")

	// cr3XPacketUUID identifies Canon's CR3 XPacket/XMP UUID payload.
	cr3XPacketUUID = meta.UUIDFromString("be7acfcb-97a9-42e8-9c71-999491e3afac")

	// cr3PreviewUUID identifies Canon's CR3 PRVW UUID payload.
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
		if err = r.callXMPReader(b, header); err != nil {
			return err
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
	err = readContainerBoxes(b, func(inner *box) error {
		switch inner.boxType {
		case typeCCTP:
			return readCCTPBox(inner)
		case typeCNCV:
			return readCNCVBox(inner)
		case typeCTBO:
			return readCTBOBox(inner)
		case typeCMT1:
			return r.readCMTBox(inner, ifds.IFD0)
		case typeCMT2:
			return r.readCMTBox(inner, ifds.ExifIFD)
		case typeCMT3:
			return r.readCMTBox(inner, ifds.MknoteIFD)
		case typeCMT4:
			return r.readCMTBox(inner, ifds.GPSIFD)
		case typeTHMB, typeThmb:
			sawTHMB = true
			return r.readTHMBBox(inner)
		default:
			if logLevelDebug() {
				logDebug().Str("boxType", inner.boxType.String()).Int64("offset", inner.offset).Int64("size", inner.size).Send()
			}
			return nil
		}
	})
	if err != nil {
		return err
	}
	if r.hasGoal(metadataKindTHMB) && !sawTHMB {
		// Some CR3 variants do not include a THMB box.
		r.setGoal(metadataKindTHMB, false)
	}
	return nil
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
	return r.callExifReader(b, header)
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
	if b.remain < int64(len(cncv.version)) {
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
	entryCount := int64(cctp.count)
	maxEntries := b.remain / 24
	if entryCount > maxEntries {
		entryCount = maxEntries
	}
	if entryCount > int64(^uint(0)>>1) {
		entryCount = int64(^uint(0) >> 1)
	}
	cctp.entries = make([]cctpEntry, 0, int(entryCount))
	for i := int64(0); i < entryCount; i++ {
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
	itemCount := int64(ctbo.count)
	maxItems := b.remain / 20
	if itemCount > maxItems {
		itemCount = maxItems
	}
	if itemCount > int64(^uint(0)>>1) {
		itemCount = int64(^uint(0) >> 1)
	}
	ctbo.items = make([]offsetLength, 0, int(itemCount))
	for i := int64(0); i < itemCount; i++ {
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
	if !r.hasGoal(metadataKindTHMB) {
		return nil
	}
	if r.previewImageReader == nil {
		return nil
	}
	var thmb thmbBox
	thmb, err = parseTHMBBox(b)
	if err != nil {
		return err
	}
	header := meta.PreviewHeader{
		Size:      thmb.size,
		Width:     thmb.width,
		Height:    thmb.height,
		ImageType: imagetype.ImageJPEG,
		Source:    meta.PreviewSourceTHMB,
	}
	if err = r.emitPreviewPayload(b, thmb.size, header, metadataKindTHMB); err != nil {
		return err
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
