package exif

import (
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/rs/zerolog"
)

// loggerMixin provides common EXIF parser logging behavior and can be embedded
// into parser types to avoid repeating level checks and trace-callsite logic.
type loggerMixin struct {
	metalog.Mixin
}

// newLoggerMixin creates and initializes an internal helper value.
func newLoggerMixin(l zerolog.Logger) loggerMixin {
	return loggerMixin{Mixin: metalog.NewComponentMixin(l, "exif")}
}

// setLogger sets the internal state value used during parsing.
func (m *loggerMixin) setLogger(l zerolog.Logger) {
	m.SetLogger(l)
}

func (m loggerMixin) logLevel() zerolog.Level {
	return m.Level()
}

// logLevelDebug reports whether debug level logging is enabled.
func (m loggerMixin) logLevelDebug() bool {
	return m.logLevel() <= zerolog.DebugLevel
}

// logLevelWarn reports whether warn level logging is enabled.
func (m loggerMixin) logLevelWarn() bool {
	return m.logLevel() <= zerolog.WarnLevel
}

// traceEnabled reports whether trace logging is enabled.
func (m loggerMixin) traceEnabled() bool {
	return m.TraceEnabled()
}

// infoEnabled reports whether info logging is enabled.
func (m loggerMixin) infoEnabled() bool {
	return m.logLevel() <= zerolog.InfoLevel
}

// debugEnabled reports whether debug logging is enabled.
func (m loggerMixin) debugEnabled() bool {
	return m.logLevelDebug()
}

// warnEnabled reports whether warn logging is enabled.
func (m loggerMixin) warnEnabled() bool {
	return m.logLevelWarn()
}

// errorEnabled reports whether error logging is enabled.
func (m loggerMixin) errorEnabled() bool {
	return m.logLevel() <= zerolog.ErrorLevel
}

// errEnabled reports whether error logging is enabled.
func (m loggerMixin) errEnabled() bool {
	return m.errorEnabled()
}

// debug builds a debug-level log event with trace caller context when enabled.
func (m loggerMixin) debug() *zerolog.Event {
	return m.Event(zerolog.DebugLevel, 3)
}

// info builds an info-level log event with trace caller context when enabled.
func (m loggerMixin) info() *zerolog.Event {
	return m.Event(zerolog.InfoLevel, 3)
}

// warn builds a warn-level log event with trace caller context when enabled.
func (m loggerMixin) warn() *zerolog.Event {
	return m.Event(zerolog.WarnLevel, 3)
}

// readerLogContext adds decode-scoped fields useful when a stream cannot be
// identified by filename, such as CR3 item payloads or embedded maker notes.
func (r *Reader) readerLogContext(ev *zerolog.Event) *zerolog.Event {
	ev.Uint32("readerOffset", r.po).
		Uint32("firstIFDOffset", r.firstIFDOffset).
		Uint32("exifLength", r.exifLength).
		Str("imageType", r.Exif.ImageType.String())
	if r.Exif.IFD0.Make != "" {
		ev.Str("cameraMake", r.Exif.IFD0.Make)
	}
	if r.Exif.IFD0.Model != "" {
		ev.Str("cameraModel", r.Exif.IFD0.Model)
	}
	return ev
}

// directoryLogContext adds IFD identity and positioning fields to a log event.
func (r *Reader) directoryLogContext(ev *zerolog.Event, d tag.Directory) *zerolog.Event {
	r.readerLogContext(ev)
	return ev.
		Str("ifd", d.String()).
		Str("ifdType", d.Type.String()).
		Int8("ifdIndex", d.Index).
		Uint32("ifdOffset", d.Offset).
		Uint32("ifdBaseOffset", d.BaseOffset).
		Str("byteOrder", d.ByteOrder.String())
}

// tagLogContext adds decoded tag identity and positioning fields to a log event.
func (r *Reader) tagLogContext(ev *zerolog.Event, t tag.Entry) *zerolog.Event {
	r.readerLogContext(ev)
	return ev.
		Uint16("tagID", uint16(t.ID)).
		Str("tagName", t.Name()).
		Str("tagTypeName", t.Type.String()).
		Str("ifd", t.IfdType.String()).
		Uint32("tagSize", t.Size()).
		Uint32("tagOffset", t.ValueOffset).
		Bool("tagEmbedded", t.IsEmbedded()).
		Str("byteOrder", t.ByteOrder.String())
}

// rawTagHeaderLogContext adds raw directory-entry fields before tag type
// normalization. It is used when a tag header is invalid and no Entry exists.
func (r *Reader) rawTagHeaderLogContext(ev *zerolog.Event, d tag.Directory, index int, id tag.ID, typ tag.Type, unitCount, valueOffset uint32) *zerolog.Event {
	r.directoryLogContext(ev, d)
	return ev.
		Int("tagIndex", index).
		Uint16("tagID", uint16(id)).
		Str("tagName", tag.NameFor(d.Type, id)).
		Uint16("tagType", uint16(typ)).
		Str("tagTypeName", typ.String()).
		Uint32("units", unitCount).
		Uint32("rawValueOffset", valueOffset)
}

// infoDirectoryLogContext keeps successful directory progress logs compact.
func (r *Reader) infoDirectoryLogContext(ev *zerolog.Event, d tag.Directory) *zerolog.Event {
	ev.Str("ifd", d.Type.String()).
		Int8("ifdIndex", d.Index).
		Uint32("ifdOffset", d.Offset).
		Str("byteOrder", d.ByteOrder.String()).
		Uint32("readerOffset", r.po)
	if d.BaseOffset != 0 {
		ev.Uint32("ifdBaseOffset", d.BaseOffset)
	}
	return ev
}

// logTagContext logs tag metadata with a caller-supplied message.
func (r *Reader) warnTagContext(t tag.Entry, msg string, includeQueueMax bool) {
	e := r.tagLogContext(r.warn(), t)
	if includeQueueMax {
		e.Int("tagQueueMax", tagQueueMax)
	}
	e.Msg(msg)
}
