package isobmff

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is a zerolog logger
	Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)
)

// logLevelInfo
func logLevelInfo() bool {
	return Logger.GetLevel() <= zerolog.InfoLevel
}

// logLevelDebug
func logLevelDebug() bool {
	return Logger.GetLevel() <= zerolog.DebugLevel
}

// logLevelError
func logLevelError() bool {
	return Logger.GetLevel() <= zerolog.ErrorLevel
}

// logLevelTrace
func logLevelTrace() bool {
	return Logger.GetLevel() == zerolog.TraceLevel
}

func logInfo() *zerolog.Event {
	ev := Logger.WithLevel(zerolog.InfoLevel)
	logTraceFunction(ev)
	return ev
}

func logDebug() *zerolog.Event {
	ev := Logger.WithLevel(zerolog.DebugLevel)
	logTraceFunction(ev)
	return ev
}

func logError() *zerolog.Event {
	ev := Logger.WithLevel(zerolog.ErrorLevel)
	logTraceFunction(ev)
	return ev
}
func logInfoBox(b *box) *zerolog.Event {
	ev := logInfo()
	if b != nil {
		b.log(ev)
	}
	return ev
}

func (b *box) log(ev *zerolog.Event) {
	ev.Str("BoxType", b.boxType.String()).Int64("offset", b.offset).Int64("size", b.size)
	if b.flags != 0 {
		ev.Object("flags", b.flags)
	}
}

func logTraceFunction(ev *zerolog.Event) {
	if logLevelTrace() {
		pc, _, _, ok := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			ev.Str("fn", details.Name())
		}
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
	e.Str("boxType", b.boxType.String()).Int64("offset", b.offset).Int64("size", b.size)
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
