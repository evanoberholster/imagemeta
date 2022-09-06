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
		ir.logger.Debug().Stringer("id", t.ID).Uint32("units", t.UnitCount).Str("tag", ifds.IfdType(t.Ifd).TagName(t.ID)).Uint32("offset", t.ValueOffset).Stringer("type", t.Type()).Send()
	}
}

func (ir *ifdReader) logTagWarn(t tag.Tag, msg string) {
	if ir.logWarn() {
		ir.logger.Warn().Stringer("id", t.ID).Uint32("units", t.UnitCount).Str("tag", ifds.IfdType(t.Ifd).TagName(t.ID)).Uint32("offset", t.ValueOffset).Uint32("reader.offset", ir.po).Stringer("type", t.Type()).Msg(msg)
	}
}
