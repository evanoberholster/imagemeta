package preview

import (
	"log/slog"
	"os"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

var (
	// Logger is the logger
	Logger = metalog.New(os.Stdout, metalog.LevelDisabled)
)

func (pr *previewReader) logError(err error) *metalog.Event {
	return metalog.ComponentEvent(pr.logger, "preview", slog.LevelError, 2).Err(err)
}
