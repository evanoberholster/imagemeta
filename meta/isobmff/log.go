package isobmff

import (
	"log/slog"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

const componentName = "isobmff"

// logLevelInfo
func logLevelInfo() bool {
	return metalog.LevelEnabled(metalog.GetLogger(), slog.LevelInfo)
}

// logLevelDebug
func logLevelDebug() bool {
	return metalog.LevelEnabled(metalog.GetLogger(), slog.LevelDebug)
}

// logLevelError
func logLevelError() bool {
	return metalog.LevelEnabled(metalog.GetLogger(), slog.LevelError)
}

func logInfo() *metalog.Event {
	return metalog.ComponentEvent(metalog.GetLogger(), componentName, slog.LevelInfo, 2)
}

func logDebug() *metalog.Event {
	return metalog.ComponentEvent(metalog.GetLogger(), componentName, slog.LevelDebug, 2)
}

func logError() *metalog.Event {
	return metalog.ComponentEvent(metalog.GetLogger(), componentName, slog.LevelError, 2)
}
func logInfoBox(b *box) *metalog.Event {
	ev := logInfo()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func logDebugBox(b *box) *metalog.Event {
	ev := logDebug()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func logErrorBox(b *box) *metalog.Event {
	ev := logError()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func (b *box) log(ev *metalog.Event) {
	ev.Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int("size", b.size)
	if b.flags != 0 {
		ev.Object("flags", b.flags)
	}
}

// MarshalLogObject is a structured logging interface
func (e cctpEntry) MarshalLogObject(ev *metalog.Event) {
	ev.Uint32("size", e.size).
		Str("trackType", fourCCString(e.trackType)).
		Uint32("mediaType", e.mediaType).
		Uint32("unknown", e.unknown).
		Uint32("index", e.index)
}

// MarshalLogObject is a structured logging interface
func (c cctpBox) MarshalLogArray(a *metalog.Array) {
	for i := range c.entries {
		a.Object(c.entries[i])
	}
}

// MarshalLogObject is a structured logging interface
func (b box) MarshalLogObject(e *metalog.Event) {
	e.Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int("size", b.size)
	if b.flags != 0 {
		e.Object("flags", b.flags)
	}
}

// MarshalLogObject is a structured logging interface
func (f flags) MarshalLogObject(e *metalog.Event) {
	e.Uint8("version", f.version()).Uint32("flags", f.flags())
}

// MarshalLogArray is a structured logging interface.
func (ctbo ctboBox) MarshalLogArray(a *metalog.Array) {
	for i := 0; i < len(ctbo.items); i++ {
		a.Object(ctbo.items[i])
	}
}
