package jpeg

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the logger
	Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel).With().Str("package", "jpeg").Logger()
)

func logInfo() bool {
	return Logger.GetLevel() <= zerolog.InfoLevel
}

func (jr *jpegReader) logMarker(str string) {
	if logInfo() {
		if len(str) == 0 {
			str = jr.marker.String()
		}
		Logger.Info().Str("marker", str).Int("length", int(jr.size)).Uint32("offset", uint32(jr.discarded)).Send()
	}
}
