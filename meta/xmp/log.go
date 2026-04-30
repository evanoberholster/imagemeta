package xmp

import (
	"log/slog"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

const componentName = "xmp"

func logWarn() *metalog.Event {
	return metalog.ComponentEvent(metalog.GetLogger(), componentName, slog.LevelWarn, 2)
}

func logWarnEnabled() bool {
	return metalog.LevelEnabled(metalog.GetLogger(), slog.LevelWarn)
}

func logPropertyParseWarning(err error, p property) {
	if !logWarnEnabled() {
		return
	}
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
