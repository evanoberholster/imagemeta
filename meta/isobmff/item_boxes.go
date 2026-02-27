package isobmff

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog"
)

const mimeContentTypeMaxLen = 64

// readIinf parses the HEIF item info box and dispatches contained infe entries.
func (r *Reader) readIinf(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return err
	}

	var use32 bool
	switch b.flags.version() {
	case 0:
		use32 = false
	case 1:
		use32 = true
	default:
		return fmt.Errorf("readIinf: unsupported version %d", b.flags.version())
	}
	count, err := readUint16Or32(b, use32)
	if err != nil {
		return err
	}

	if logLevelInfo() {
		logInfo().Object("box", b).Uint32("count", count).Send()
	}

	var parsed uint32
	err = readContainerBoxes(b, func(inner *box) error {
		if inner.boxType == typeInfe {
			parsed++
			return r.readInfe(inner)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if logLevelDebug() && parsed != count {
		logDebug().Object("box", b).Uint32("declared", count).Uint32("parsed", parsed).Msg("iinf entry count mismatch")
	}
	return nil
}

// readInfe parses an item info entry and records IDs for Exif/XMP item payloads.
func (r *Reader) readInfe(b *box) (err error) {
	if err = b.readFlags(); err != nil {
		return err
	}

	var (
		id     itemID
		idSize uint8
	)
	switch b.flags.version() {
	case 2:
		idSize = 2
	case 3:
		idSize = 4
	default:
		if logLevelDebug() {
			logDebug().Object("box", b).Uint8("version", b.flags.version()).Msg("skipping unsupported infe version")
		}
		return nil
	}
	id, err = readItemIDBySize(b, idSize)
	if err != nil {
		return err
	}

	protectionIndex, err := b.readUint16()
	if err != nil {
		return err
	}
	itemFourCC, err := b.readFourCC()
	if err != nil {
		return err
	}
	var itemTypeBuf [4]byte
	bmffEndian.PutUint32(itemTypeBuf[:], itemFourCC)
	itemType := itemTypeFromBuf(itemTypeBuf[:])

	// item_name
	if err = b.discardCString(maxBoxStringLength); err != nil {
		return err
	}

	var contentType string
	switch itemType {
	case itemTypeMime:
		needContentType := r.hasGoal(metadataKindXMP) || logLevelDebug()
		if needContentType {
			var contentTypeFixed [mimeContentTypeMaxLen]byte
			buf, err := b.readCStringBytes(contentTypeFixed[:0], mimeContentTypeMaxLen)
			if err != nil {
				return err
			}
			if r.hasGoal(metadataKindXMP) && isXMPMIMETypeBytes(buf) {
				r.heic.xml.id = id
			}
			if logLevelDebug() {
				contentType = string(buf)
			}
		} else if err = b.discardCString(mimeContentTypeMaxLen); err != nil {
			return err
		}
	case itemTypeExif:
		r.heic.exif.id = id
	case itemTypeURI:
		// item_uri_type
		if err = b.discardCString(maxBoxStringLength); err != nil {
			return err
		}
	}

	if logLevelDebug() {
		ev := logDebug().
			Object("box", b).
			Object("flags", b.flags).
			Uint32("itemID", uint32(id)).
			Str("itemType", string(itemTypeBuf[:])).
			Uint16("idx", protectionIndex)
		if itemType == itemTypeMime {
			ev.Str("contentType", contentType)
		}
		ev.Send()
	}
	return nil
}

func isXMPMIMETypeBytes(contentType []byte) bool {
	ct := bytes.TrimSpace(contentType)
	switch {
	case asciiEqualFoldBytes(ct, []byte("application/rdf+xml")),
		asciiEqualFoldBytes(ct, []byte("application/xml")),
		asciiEqualFoldBytes(ct, []byte("text/xml")):
		return true
	}
	return asciiContainsFoldBytes(ct, []byte("xmp")) || asciiContainsFoldBytes(ct, []byte("rdf+xml"))
}

func asciiContainsFoldBytes(s, sub []byte) bool {
	if len(sub) == 0 {
		return true
	}
	if len(s) < len(sub) {
		return false
	}
	first := toASCIILowerByte(sub[0])
	end := len(s) - len(sub)
	for i := 0; i <= end; i++ {
		if toASCIILowerByte(s[i]) != first {
			continue
		}
		if asciiEqualFoldBytes(s[i:i+len(sub)], sub) {
			return true
		}
	}
	return false
}

func asciiEqualFoldBytes(a, b []byte) bool {
	for i := 0; i < len(b); i++ {
		if toASCIILowerByte(a[i]) != toASCIILowerByte(b[i]) {
			return false
		}
	}
	return true
}

func toASCIILowerByte(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

// itemType
type itemType uint8

// itemTypes
const (
	itemTypeUnknown itemType = iota
	itemTypeInfe
	itemTypeMime
	itemTypeURI
	itemTypeAv01
	itemTypeHvc1
	itemTypeGrid
	itemTypeExif
)

var (
	itemTypeHvc1FourCC = fourCCFromString("hvc1")
	itemTypeExifFourCC = fourCCFromString("Exif")
	itemTypeAv01FourCC = fourCCFromString("av01")
	itemTypeGridFourCC = fourCCFromString("grid")
	itemTypeInfeFourCC = fourCCFromString("infe")
	itemTypeMimeFourCC = fourCCFromString("mime")
	itemTypeURIFourCC  = fourCCFromString("uri ")
)

// itemTypeFromBuf maps the 4-byte infe item_type field to an internal itemType.
func itemTypeFromBuf(buf []byte) itemType {
	if len(buf) < 4 {
		return itemTypeUnknown
	}

	switch bmffEndian.Uint32(buf[:4]) {
	case itemTypeHvc1FourCC:
		return itemTypeHvc1
	case itemTypeExifFourCC:
		return itemTypeExif
	case itemTypeAv01FourCC:
		return itemTypeAv01
	case itemTypeGridFourCC:
		return itemTypeGrid
	case itemTypeInfeFourCC:
		return itemTypeInfe
	case itemTypeMimeFourCC:
		return itemTypeMime
	case itemTypeURIFourCC:
		return itemTypeURI
	default:
		return itemTypeUnknown
	}
}

// itemLocationBox is a "iloc" box
type itemLocationBox struct {
	count                                             uint32
	offsetSize, lengthSize, baseOffsetSize, indexSize uint8 // actually uint4
}

// MarshalZerologObject is a zerolog interface for logging
func (ilb itemLocationBox) MarshalZerologObject(e *zerolog.Event) {
	e.Uint32("items", ilb.count).Uint8("offsetSize", ilb.offsetSize).Uint8("lengthSize", ilb.lengthSize).Uint8("baseOffsetSize", ilb.baseOffsetSize).Uint8("indexSize", ilb.indexSize)
}

// ilocEntry is not a box
type ilocEntry struct {
	firstExtent        offsetLength
	id                 itemID
	baseOffset         uint64 // uint32 or uint64, depending on encoding
	count              uint16
	dataReferenceIndex uint16 // dri
	constructionMethod uint8  // cmeth actually uint4
}

// MarshalZerologObject is a zerolog interface for logging
func (ie ilocEntry) MarshalZerologObject(e *zerolog.Event) {
	e.Uint32("itemID", uint32(ie.id)).Object("extent", ie.firstExtent).Uint16("count", ie.count).Uint16("dri", ie.dataReferenceIndex).Uint8("cmeth", ie.constructionMethod)
}

// offsetLength contains an offset and length
type offsetLength struct {
	offset, length uint64
}

// MarshalZerologObject is a zerolog interface for logging
func (ol offsetLength) MarshalZerologObject(e *zerolog.Event) {
	e.Uint64("length", ol.length).Uint64("offset", ol.offset)
}

// readIloc parses item location extents and stores first-extent offsets for
// metadata items discovered in iinf/infe.
func (r *Reader) readIloc(b *box) (err error) {
	ilb, err := readIlocHeader(b)
	if err != nil {
		return err
	}

	for i := uint32(0); i < ilb.count; i++ {
		var ent ilocEntry
		var idSize uint8
		switch b.flags.version() {
		case 0, 1:
			idSize = 2
		case 2:
			idSize = 4
		default:
			return fmt.Errorf("readIloc: unsupported version %d", b.flags.version())
		}
		ent.id, err = readItemIDBySize(b, idSize)
		if err != nil {
			return err
		}

		if b.flags.version() > 0 { // versions 1 and 2
			cmeth, readErr := b.readUint16()
			if readErr != nil {
				return readErr
			}
			ent.constructionMethod = uint8(cmeth & 0x0f)
		}
		ent.dataReferenceIndex, err = b.readUint16()
		if err != nil {
			return err
		}

		ent.baseOffset, err = b.readUintN(ilb.baseOffsetSize)
		if err != nil {
			return err
		}
		ent.count, err = b.readUint16()
		if err != nil {
			return err
		}

		firstExtentResolved := false
		for j := 0; j < int(ent.count); j++ {
			if b.flags.version() > 0 && ilb.indexSize > 0 {
				if _, err = b.readUintN(ilb.indexSize); err != nil {
					return err
				}
			}

			extentOffset, readErr := b.readUintN(ilb.offsetSize)
			if readErr != nil {
				return readErr
			}
			extentLength, readErr := b.readUintN(ilb.lengthSize)
			if readErr != nil {
				return readErr
			}

			if j == 0 {
				offset, ok := r.resolveIlocExtentOffset(ent, extentOffset)
				if ok {
					ent.firstExtent = offsetLength{
						offset: offset,
						length: extentLength,
					}
					firstExtentResolved = true
				}
			}
		}
		if logLevelDebug() {
			logDebug().
				Uint32("itemID", uint32(ent.id)).
				Uint64("offset", ent.firstExtent.offset).
				Uint64("length", ent.firstExtent.length).
				Uint16("count", ent.count).
				Uint16("dri", ent.dataReferenceIndex).
				Uint8("cmeth", ent.constructionMethod).
				Send()
		}

		switch ent.id {
		case r.heic.exif.id:
			if firstExtentResolved {
				r.heic.exif.ol = ent.firstExtent
			}
		case r.heic.xml.id:
			if firstExtentResolved {
				r.heic.xml.ol = ent.firstExtent
			}
		}
	}
	return b.close()
}

func (r *Reader) resolveIlocExtentOffset(ent ilocEntry, extentOffset uint64) (uint64, bool) {
	if ent.dataReferenceIndex != 0 {
		if logLevelDebug() {
			logDebug().Uint16("dri", ent.dataReferenceIndex).Msg("skip iloc entry with external data reference")
		}
		return 0, false
	}
	if ent.baseOffset > ^uint64(0)-extentOffset {
		return 0, false
	}
	rel := ent.baseOffset + extentOffset

	switch ent.constructionMethod {
	case 0:
		return rel, true
	case 1:
		if r.heic.idatData.length == 0 {
			if logLevelDebug() {
				logDebug().Uint32("itemID", uint32(ent.id)).Msg("skip iloc idat method without idat payload")
			}
			return 0, false
		}
		if rel > r.heic.idatData.length {
			return 0, false
		}
		if r.heic.idatData.offset > ^uint64(0)-rel {
			return 0, false
		}
		return r.heic.idatData.offset + rel, true
	default:
		if logLevelDebug() {
			logDebug().Uint8("cmeth", ent.constructionMethod).Uint32("itemID", uint32(ent.id)).Msg("skip unsupported iloc construction method")
		}
		return 0, false
	}
}

// readIlocHeader parses iloc version/flags and field-size descriptors.
func readIlocHeader(b *box) (ilb itemLocationBox, err error) {
	if err = b.readFlags(); err != nil {
		return ilb, err
	}

	buf, err := b.Peek(2)
	if err != nil {
		return ilb, fmt.Errorf("readIlocHeader: %w", ErrBufLength)
	}
	ilb.offsetSize = buf[0] >> 4
	ilb.lengthSize = buf[0] & 15
	ilb.baseOffsetSize = buf[1] >> 4
	if b.flags.version() > 0 { // versions 1 and 2
		ilb.indexSize = buf[1] & 15
	}
	if _, err = b.Discard(2); err != nil {
		return ilb, err
	}

	switch b.flags.version() {
	case 0, 1:
		ilb.count, err = readUint16Or32(b, false)
		if err != nil {
			return ilb, err
		}
	case 2:
		ilb.count, err = readUint16Or32(b, true)
		if err != nil {
			return ilb, err
		}
	default:
		return ilb, fmt.Errorf("readIlocHeader: unsupported version %d", b.flags.version())
	}

	if logLevelInfo() {
		logInfoBox(b).Object("ItemLocation", ilb).Send()
	}
	return ilb, nil
}
