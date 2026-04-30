package xmp

import (
	"log/slog"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

func componentLogger() *slog.Logger {
	return metalog.ComponentLogger("xmp")
}

func logWarn() *metalog.Event {
	return metalog.NewEvent(componentLogger(), slog.LevelWarn, 2)
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
