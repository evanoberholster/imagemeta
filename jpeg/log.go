package jpeg

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DefaulLogger zerolog.Logger  = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel).With().Str("package", "jpeg").Logger()
	Logger       *zerolog.Logger = nil
)

func logInfo() bool {
	return Logger != nil
}

func logInfoMarker(markerStr string, length int, offset int) {
	Logger.Info().Str("marker", markerStr).Int("length", length).Uint32("offset", uint32(offset)).Send()
}
