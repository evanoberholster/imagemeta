package preview

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the logger
	Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel).With().Str("package", "preview").Logger()
)

func (pr *previewReader) logError(err error) *zerolog.Event {
	return pr.logger.Err(err)
}
