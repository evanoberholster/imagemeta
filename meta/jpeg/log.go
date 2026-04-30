package jpeg

import (
	"log/slog"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

const componentName = "jpeg"

func logInfo() bool {
	return metalog.LevelEnabled(metalog.Logger, slog.LevelInfo)
}

func logDebug() bool {
	return metalog.LevelEnabled(metalog.Logger, slog.LevelDebug)
}

func logInfoEvent() *metalog.Event {
	return metalog.ComponentEvent(metalog.Logger, componentName, slog.LevelInfo, 2)
}

func logDebugEvent() *metalog.Event {
	return metalog.ComponentEvent(metalog.Logger, componentName, slog.LevelDebug, 2)
}

func (jr *jpegReader) logMarker(str string) {
	if logInfo() {
		if len(str) == 0 {
			str = jr.marker.String()
		}
		logInfoEvent().
			Str("marker", str).
			Int("length", int(jr.size)).
			Uint32("offset", uint32(jr.discarded)).
			Msg("read jpeg marker")
	}
}

func (jr *jpegReader) logDecodedItem(kind string, size int) {
	if !logInfo() {
		return
	}
	logInfoEvent().
		Str("metadataKind", kind).
		Str("marker", jr.marker.String()).
		Int("length", size).
		Uint32("offset", uint32(jr.discarded)).
		Msg("decoded metadata item")
}
