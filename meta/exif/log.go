package exif

import (
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

// readerLogContext adds decode-scoped fields useful when a stream cannot be
// identified by filename, such as CR3 item payloads or embedded maker notes.
func (r *Reader) readerLogContext(ev *metalog.Event) *metalog.Event {
	ev.Uint32("readerOffset", r.po).
		Uint32("firstIFDOffset", r.firstIFDOffset).
		Uint32("exifLength", r.exifLength).
		Str("imageType", r.Exif.ImageType.String())
	return ev
}

// directoryLogContext adds IFD identity and positioning fields to a log event.
func (r *Reader) directoryLogContext(ev *metalog.Event, d tag.Directory) *metalog.Event {
	r.readerLogContext(ev)
	return ev.
		Str("ifd", d.String()).
		Str("ifdType", d.Type.String()).
		Int8("ifdIndex", d.Index).
		Uint32("ifdOffset", d.Offset).
		Uint32("ifdBaseOffset", d.BaseOffset)
}

// tagLogContext adds decoded tag identity and positioning fields to a log event.
func (r *Reader) tagLogContext(ev *metalog.Event, t tag.Entry) *metalog.Event {
	r.readerLogContext(ev)
	return ev.
		Str("tagID", tag.HexUint16Upper(uint16(t.ID))).
		Str("tagName", t.Name()).
		Str("tagType", t.Type.String()).
		Str("ifd", t.IfdType.String()).
		Uint32("tagSize", t.Size()).
		Uint32("tagOffset", t.ValueOffset).
		Bool("tagEmbedded", t.IsEmbedded())
}

// rawTagHeaderLogContext adds raw directory-entry fields before tag type
// normalization. It is used when a tag header is invalid and no Entry exists.
func (r *Reader) rawTagHeaderLogContext(ev *metalog.Event, d tag.Directory, index int, id tag.ID, typ tag.Type, unitCount, valueOffset uint32) *metalog.Event {
	r.directoryLogContext(ev, d)
	return ev.
		Int("tagIndex", index).
		Str("tagID", tag.HexUint16Upper(uint16(id))).
		Str("tagName", tag.NameFor(d.Type, id)).
		Str("tagType", typ.String()).
		Uint32("units", unitCount).
		Uint32("rawValueOffset", valueOffset)
}

// infoDirectoryLogContext keeps successful directory progress logs compact.
func (r *Reader) infoDirectoryLogContext(ev *metalog.Event, d tag.Directory) *metalog.Event {
	ev.Str("ifd", d.Type.String()).
		Int8("ifdIndex", d.Index).
		Uint32("ifdOffset", d.Offset).
		Uint32("readerOffset", r.po)
	appendByteOrderIfSet(ev, d.ByteOrder)
	if d.BaseOffset != 0 {
		ev.Uint32("ifdBaseOffset", d.BaseOffset)
	}
	return ev
}

func byteOrderShort(order utils.ByteOrder) string {
	switch order {
	case utils.LittleEndian:
		return "LE"
	case utils.BigEndian:
		return "BE"
	default:
		return order.String()
	}
}

func appendByteOrderIfSet(ev *metalog.Event, order utils.ByteOrder) {
	switch order {
	case utils.LittleEndian, utils.BigEndian:
		ev.Str("byteOrder", byteOrderShort(order))
	}
}

// logTagContext logs tag metadata with a caller-supplied message.
func (r *Reader) warnTagContext(t tag.Entry, msg string, includeQueueMax bool) {
	e := r.tagLogContext(r.Warn(3), t)
	if includeQueueMax {
		e.Int("tagQueueMax", tagQueueMax)
	}
	e.Msg(msg)
}
