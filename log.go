package imagemeta

import (
	"io"
	"os"

	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/evanoberholster/imagemeta/isobmff"
	"github.com/evanoberholster/imagemeta/jpeg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the logger
	logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)
)

func SetLogger(w io.Writer, level zerolog.Level) {
	logger = log.Output(w).Level(level)
	jpeg.Logger = logger
	exif2.Logger = logger
	isobmff.Logger = logger
}
