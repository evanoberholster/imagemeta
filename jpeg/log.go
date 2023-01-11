package jpeg

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// DefaultLoggger is the Default Logger, logs only Panic to the console
	DefaulLogger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel).With().Str("package", "jpeg").Logger()
	logLevel     zerolog.Level
	logger       *zerolog.Logger = nil
)

func logInfo() bool {
	return logLevel >= zerolog.InfoLevel
}

// SetLogger sets the package logger
func SetLogger(l *zerolog.Logger) {
	logLevel = l.GetLevel()
	logger = l
}

// Logger returns the package logger
func Logger() *zerolog.Logger {
	return logger
}

func (jr *jpegReader) logMarker(str string) {
	if logInfo() {
		if len(str) == 0 {
			str = jr.marker.String()
		}
		logger.Info().Str("marker", str).Int("length", int(jr.size)).Uint32("offset", uint32(jr.discarded)).Send()
	}
}
