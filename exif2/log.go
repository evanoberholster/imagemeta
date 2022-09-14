package exif2

import (
	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/tag"
	"github.com/rs/zerolog"
)

func (ir *ifdReader) logInfo() bool {
	return ir.logger.GetLevel() <= zerolog.InfoLevel
}

func (ir *ifdReader) logWarn() bool {
	return ir.logger.GetLevel() <= zerolog.WarnLevel
}

func (ir *ifdReader) logError() bool {
	return ir.logger.GetLevel() <= zerolog.ErrorLevel
}

func (ir *ifdReader) logTagInfo(t tag.Tag) {
	if ir.logInfo() {
		ir.logger.Debug().Stringer("id", t.ID).Uint32("units", t.UnitCount).Str("tag", ifds.IfdType(t.Ifd).TagName(t.ID)).Stringer("offset", t.ValueOffset).Stringer("type", t.Type()).Send()
	}
}

func (ir *ifdReader) logTagWarn(t tag.Tag, msg string) {
	if ir.logWarn() {
		ir.logger.Warn().Stringer("id", t.ID).Uint32("units", t.UnitCount).Str("tag", ifds.IfdType(t.Ifd).TagName(t.ID)).Stringer("offset", t.ValueOffset).Uint32("readerOffset", ir.po).Stringer("type", t.Type()).Msg(msg)
	}
}

func (ir *ifdReader) logParseWarn(t tag.Tag, fnName string, msg string, err error) {
	if ir.logWarn() {
		l := ir.logger.Warn()
		if err != nil {
			l = l.Err(err)
		}
		l.Str("func", fnName).Stringer("id", t.ID).Uint32("units", t.UnitCount).Str("tag", ifds.IfdType(t.Ifd).TagName(t.ID)).Stringer("offset", t.ValueOffset).Uint32("readerOffset", ir.po).Stringer("type", t.Type()).Msg(msg)
	}
}
