package xmp

import (
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/rs/zerolog"
)

func componentLogger() zerolog.Logger {
	return metalog.ComponentLogger("xmp")
}

func logWarn() *zerolog.Event {
	return metalog.Event(componentLogger(), zerolog.WarnLevel, 2)
}

func logPropertyParseWarning(err error, p property) {
	ev := logWarn().Err(err).
		Str("namespace", p.Namespace().String()).
		Str("property", p.Property().String()).
		Str("parent", p.Parent().String()).
		Int("regionIndex", p.RegionIndex()).
		Int("valueLength", len(p.Value()))
	if preview := p.valuePreview(96); preview != "" {
		ev.Str("valuePreview", preview)
	}
	ev.Msg("xmp property parse warning")
}
