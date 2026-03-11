package exif

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the package logger.
	Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel).With().Str("package", "exif").Logger()
)

// loggerMixin provides common EXIF parser logging behavior and can be embedded
// into parser types to avoid repeating level checks and trace-callsite logic.
type loggerMixin struct {
	logger zerolog.Logger
}

// newLoggerMixin creates and initializes an internal helper value.
func newLoggerMixin(l zerolog.Logger) loggerMixin {
	return loggerMixin{logger: l}
}

// setLogger sets the internal state value used during parsing.
func (m *loggerMixin) setLogger(l zerolog.Logger) {
	m.logger = l
}

func (m loggerMixin) logLevel() zerolog.Level {
	return m.logger.GetLevel()
}

// logLevelDebug reports whether debug level logging is enabled.
func (m loggerMixin) logLevelDebug() bool {
	return m.logLevel() <= zerolog.DebugLevel
}

// logLevelWarn reports whether warn level logging is enabled.
func (m loggerMixin) logLevelWarn() bool {
	return m.logLevel() <= zerolog.WarnLevel
}

// traceEnabled reports whether trace logging is enabled.
func (m loggerMixin) traceEnabled() bool {
	return m.logLevel() == zerolog.TraceLevel
}

// infoEnabled reports whether info logging is enabled.
func (m loggerMixin) infoEnabled() bool {
	return m.logLevel() <= zerolog.InfoLevel
}

// debugEnabled reports whether debug logging is enabled.
func (m loggerMixin) debugEnabled() bool {
	return m.logLevelDebug()
}

// warnEnabled reports whether warn logging is enabled.
func (m loggerMixin) warnEnabled() bool {
	return m.logLevelWarn()
}

// errorEnabled reports whether error logging is enabled.
func (m loggerMixin) errorEnabled() bool {
	return m.logLevel() <= zerolog.ErrorLevel
}

// errEnabled reports whether error logging is enabled.
func (m loggerMixin) errEnabled() bool {
	return m.errorEnabled()
}

// debug builds a debug-level log event with trace caller context when enabled.
func (m loggerMixin) debug() *zerolog.Event {
	ev := m.logger.Debug()
	m.traceCaller(ev, 3)
	return ev
}

// warn builds a warn-level log event with trace caller context when enabled.
func (m loggerMixin) warn() *zerolog.Event {
	ev := m.logger.Warn()
	m.traceCaller(ev, 3)
	return ev
}

// traceCaller annotates the event with caller function information at trace level.
func (m loggerMixin) traceCaller(ev *zerolog.Event, depth int) {
	if !m.traceEnabled() {
		return
	}
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return
	}
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		ev.Str("fn", fn.Name())
	}
}
