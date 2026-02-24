package isobmff

import (
	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/rs/zerolog"
)

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
		if err != nil {
			if logLevelError() {
				logError().Str("boxType", inner.boxType.String()).Int("offset", inner.offset).Int64("size", inner.size).Err(err).Send()
			}
			return err
		}
		if err = inner.close(); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	if r.goalTHMB && !sawTHMB {
		// Some CR3 variants do not include a THMB box.
		r.goalTHMB = false
	}
	return b.close()
}

// CMT Box

// readCMTBox reads an ISOBMFF Box "CMT1","CMT2","CMT3", or "CMT4" from CR3.
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
	r.haveExif = true
	return nil
}

// CNCV Box

// CNCVBox is Canon Compressor Version box
// CaNon Codec Version?
type CNCVBox struct {
	//format [9]byte
	//version [6]uint8
	version [30]byte
}

func readCNCVBox(b *box) (err error) {
	if !b.isType(typeCNCV) {
		return ErrWrongBoxType
	}
	if !logLevelInfo() {
		return nil
	}
	var cncv CNCVBox
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

// CCTP Box

type CCTPBox struct {
	count   uint32
	entries []CCTPEntry
}

type CCTPEntry struct {
	size      uint32
	trackType uint32
	mediaType uint32
	unknown   uint32
	index     uint32
}

func (e CCTPEntry) MarshalZerologObject(ev *zerolog.Event) {
	ev.Uint32("size", e.size).
		Str("trackType", fourCCString(e.trackType)).
		Uint32("mediaType", e.mediaType).
		Uint32("unknown", e.unknown).
		Uint32("index", e.index)
}

func (c CCTPBox) MarshalZerologArray(a *zerolog.Array) {
	for i := range c.entries {
		a.Object(c.entries[i])
	}
}

func readCCTPBox(b *box) (err error) {
	if !b.isType(typeCCTP) {
		return ErrWrongBoxType
	}
	if !logLevelInfo() {
		return nil
	}
	var cctp CCTPBox
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
	cctp.entries = make([]CCTPEntry, 0, entryCount)
	for i := 0; i < entryCount; i++ {
		var ent CCTPEntry
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

// CTBO Box

func readCTBOBox(b *box) (err error) {
	if !b.isType(typeCTBO) {
		return ErrWrongBoxType
	}
	if !logLevelInfo() {
		return nil
	}
	var ctbo CTBOBox
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

// CTBOBox is a Canon tracks base offsets box.
type CTBOBox struct {
	items []offsetLength
	count uint32
}

// MarshalZerologArray is a zerolog interface for logging.
func (ctbo CTBOBox) MarshalZerologArray(a *zerolog.Array) {
	for i := 0; i < len(ctbo.items); i++ {
		a.Object(ctbo.items[i])
	}
}

// THMB Box

type THMBBox struct {
	Width  uint16
	Height uint16
	Size   uint32
}

func (r *Reader) readTHMBBox(b *box) (err error) {
	if !b.isType(typeTHMB) && !b.isType(typeThmb) {
		return ErrWrongBoxType
	}
	if r.previewImageReader == nil {
		return nil
	}
	var thmb THMBBox
	thmb, err = parseTHMBBox(b)
	if err != nil {
		return err
	}
	if thmb.Size > uint32(b.remain) {
		return ErrRemainLengthInsufficient
	}

	if thmb.Size > 0 {
		payloadOffset := b.offset + int(b.size) - b.remain
		payload := box{
			reader:  b.reader,
			outer:   b,
			boxType: b.boxType,
			offset:  payloadOffset,
			size:    int64(thmb.Size),
			remain:  int(thmb.Size),
		}
		header := meta.PreviewHeader{
			Size:      thmb.Size,
			Width:     thmb.Width,
			Height:    thmb.Height,
			ImageType: imagetype.ImageJPEG,
			Source:    meta.PreviewSourceTHMB,
		}
		callbackErr := r.previewImageReader(newLimitedReader(&payload, payload.remain), header)
		if callbackErr != nil {
			if err = handleCallbackError(&payload, callbackErr); err != nil {
				return err
			}
		} else {
			r.haveTHMB = true
		}
		if closeErr := payload.close(); closeErr != nil {
			return closeErr
		}
	}

	return nil
}

func parseTHMBBox(b *box) (thmb THMBBox, err error) {
	if b.remain < 16 {
		return thmb, ErrBufLength
	}
	buf, err := b.Peek(16)
	if err != nil {
		return thmb, err
	}

	thmb.Width = bmffEndian.Uint16(buf[4:6])
	thmb.Height = bmffEndian.Uint16(buf[6:8])
	thmb.Size = bmffEndian.Uint32(buf[8:12])

	if _, err = b.Discard(16); err != nil {
		return thmb, err
	}
	return thmb, nil
}

func fourCCFromUint32(v uint32) [4]byte {
	return [4]byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

func fourCCString(v uint32) string {
	buf := fourCCFromUint32(v)
	return string(buf[:])
}

// Trak
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
