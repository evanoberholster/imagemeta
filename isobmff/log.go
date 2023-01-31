package isobmff

import (
	"os"
	"runtime"

	"github.com/pkg/errors"
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
	level := Logger.GetLevel()
	return level == zerolog.TraceLevel
}

func logErrorMsg(key string, format string, args ...interface{}) {
	Logger.Error().AnErr(key, errors.Errorf(format, args...)).Send()
}

func logInfoBox(b box) {
	ev := Logger.Info()
	b.log(ev)
	logTraceFunction(ev)
	ev.Send()
}

func logBoxExt(b *box, level zerolog.Level) *zerolog.Event {
	ev := Logger.WithLevel(level)
	if b != nil {
		b.log(ev)
	}
	logTraceFunction(ev)
	return ev
}

func (b *box) log(ev *zerolog.Event) {
	ev.Str("BoxType", b.boxType.String()).Int("offset", b.offset).Int64("size", b.size)
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
