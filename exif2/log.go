package exif2

import (
	"fmt"
	"runtime"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/rs/zerolog"
)

func (ir *ifdReader) logLevelTrace() bool {
	return ir.logger.GetLevel() == zerolog.TraceLevel
}

func (ir *ifdReader) logLevelInfo() bool {
	return ir.logger.GetLevel() <= zerolog.InfoLevel
}

func (ir *ifdReader) logLevelWarn() bool {
	return ir.logger.GetLevel() <= zerolog.WarnLevel
}

func (ir *ifdReader) logLevelError() bool {
	return ir.logger.GetLevel() <= zerolog.ErrorLevel
}

func (ir *ifdReader) logTagInfo(t tag.Tag) {
	if ir.logLevelInfo() {
		e := ir.logger.Debug()
		logTag(e, t).Send()
	}
}

//func (ir *ifdReader) logTagWarn(t tag.Tag, msg string) {
//	if ir.logLevelWarn() {
//		e := ir.logger.Warn().Uint32("readerOffset", ir.po)
//		logTag(e, t).Send()
//	}
//}

func (ir *ifdReader) logParseWarn(t tag.Tag, fnName string, msg string, err error) {
	if ir.logLevelWarn() {
		e := ir.logger.Warn()
		if err != nil {
			e = e.Err(err)
		}
		e.Str("func", fnName)
		logTag(e, t).Uint32("readerOffset", ir.po).Msg(msg)
	}
}

func (ir *ifdReader) logInfo() *zerolog.Event {
	e := ir.logger.Info()
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

func logTag(e *zerolog.Event, t tag.Tag) *zerolog.Event {
	return e.Stringer("id", t.ID).Uint32("units", t.UnitCount).Str("tag", ifds.IfdType(t.Ifd).TagName(t.ID)).Str("offset", fmt.Sprintf("0x%04x", t.ValueOffset)).Stringer("type", t.Type())
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
