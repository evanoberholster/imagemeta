package imagemeta

import (
	"io"
	"log/slog"
	"os"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
)

var (
	// Logger is the logger
	logger = metalog.New(os.Stdout, metalog.LevelDisabled)
)

func SetLogger(w io.Writer, level slog.Level) {
	logger = metalog.New(w, level)
	metalog.SetLogger(logger)
}
