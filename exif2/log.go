package exif2

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger default is Panic Level to Console Writer
var Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)

func (ir *ifdReader) logLevelTrace() bool {
	return ir.logger.GetLevel() == zerolog.TraceLevel
}

func (ir *ifdReader) logLevelInfo() bool {
	return ir.logger.GetLevel() <= zerolog.InfoLevel
}

func (ir *ifdReader) logLevelDebug() bool {
	return ir.logger.GetLevel() <= zerolog.DebugLevel
}

func (ir *ifdReader) logLevelWarn() bool {
	return ir.logger.GetLevel() <= zerolog.WarnLevel
}

func (ir *ifdReader) logLevelError() bool {
	return ir.logger.GetLevel() <= zerolog.ErrorLevel
}

func (ir *ifdReader) logInfo() *zerolog.Event {
	e := ir.logger.Info()
	ir.logTraceFunction(e)
	return e
}

func (ir *ifdReader) logDebug() *zerolog.Event {
	e := ir.logger.Debug()
	ir.logTraceFunction(e)
	return e
}

func (ir *ifdReader) logWarn() *zerolog.Event {
	e := ir.logger.Warn()
	ir.logTraceFunction(e)
	return e
}

func (ir *ifdReader) logError(err error) *zerolog.Event {
	e := ir.logger.Err(err)
	ir.logTraceFunction(e)
	return e
}

func (ir *ifdReader) logTraceFunction(ev *zerolog.Event) {
	if ir.logLevelTrace() {
		pc, _, _, ok := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			ev.Str("fn", details.Name())
		}
	}
}
