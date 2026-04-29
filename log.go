package imagemeta

import (
	"io"
	"os"

	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the logger
	logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)
)

func SetLogger(w io.Writer, level zerolog.Level) {
	logger = log.Output(w).Level(level)
	metalog.Logger = logger
}
