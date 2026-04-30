package preview

import (
	"log/slog"
	"os"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

var (
	// Logger is the logger
	Logger = metalog.WithComponent(metalog.New(os.Stdout, metalog.LevelDisabled), "preview")
)

func (pr *previewReader) logError(err error) *metalog.Event {
	return metalog.NewEvent(pr.logger, slog.LevelError, 2).Err(err)
}
