package isobmff

import (
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/rs/zerolog"
)

const componentName = "isobmff"

// logLevelInfo
func logLevelInfo() bool {
	return metalog.LevelEnabled(metalog.Logger, zerolog.InfoLevel)
}

// logLevelDebug
func logLevelDebug() bool {
	return metalog.LevelEnabled(metalog.Logger, zerolog.DebugLevel)
}

// logLevelError
func logLevelError() bool {
	return metalog.LevelEnabled(metalog.Logger, zerolog.ErrorLevel)
}

func logInfo() *zerolog.Event {
	return metalog.ComponentEvent(metalog.Logger, componentName, zerolog.InfoLevel, 2)
}

func logDebug() *zerolog.Event {
	return metalog.ComponentEvent(metalog.Logger, componentName, zerolog.DebugLevel, 2)
}

func logError() *zerolog.Event {
	return metalog.ComponentEvent(metalog.Logger, componentName, zerolog.ErrorLevel, 2)
}
func logInfoBox(b *box) *zerolog.Event {
	ev := logInfo()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func logDebugBox(b *box) *zerolog.Event {
	ev := logDebug()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func logErrorBox(b *box) *zerolog.Event {
	ev := logError()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func (b *box) log(ev *zerolog.Event) {
	ev.Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int("size", b.size)
	if b.flags != 0 {
		ev.Object("flags", b.flags)
	}
}

// MarshalZerologObject is a zerolog interface for logging
func (e cctpEntry) MarshalZerologObject(ev *zerolog.Event) {
	ev.Uint32("size", e.size).
		Str("trackType", fourCCString(e.trackType)).
		Uint32("mediaType", e.mediaType).
		Uint32("unknown", e.unknown).
		Uint32("index", e.index)
}

// MarshalZerologObject is a zerolog interface for logging
func (c cctpBox) MarshalZerologArray(a *zerolog.Array) {
	for i := range c.entries {
		a.Object(c.entries[i])
	}
}

// MarshalZerologObject is a zerolog interface for logging
func (b box) MarshalZerologObject(e *zerolog.Event) {
	e.Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int("size", b.size)
	if b.flags != 0 {
		e.Object("flags", b.flags)
	}
}

// MarshalZerologObject is a zerolog interface for logging
func (f flags) MarshalZerologObject(e *zerolog.Event) {
	e.Uint8("version", f.version()).Uint32("flags", f.flags())
}

// MarshalZerologArray is a zerolog interface for logging.
func (ctbo ctboBox) MarshalZerologArray(a *zerolog.Array) {
	for i := 0; i < len(ctbo.items); i++ {
		a.Object(ctbo.items[i])
	}
}
