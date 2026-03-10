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
	logger      zerolog.Logger
	enabledMask uint8
}

const (
	logMaskTrace uint8 = 1 << iota
	logMaskInfo
	logMaskDebug
	logMaskWarn
	logMaskError
)

// newLoggerMixin creates and initializes an internal helper value.
func newLoggerMixin(l zerolog.Logger) loggerMixin {
	m := loggerMixin{logger: l}
	m.refreshEnabledMask()
	return m
}

// setLogger sets the internal state value used during parsing.
func (m *loggerMixin) setLogger(l zerolog.Logger) {
	m.logger = l
	m.refreshEnabledMask()
}

// refreshEnabledMask precomputes level checks into a compact bitmask.
func (m *loggerMixin) refreshEnabledMask() {
	level := m.logger.GetLevel()
	var mask uint8
	if level == zerolog.TraceLevel {
		mask |= logMaskTrace
	}
	if level <= zerolog.InfoLevel {
		mask |= logMaskInfo
	}
	if level <= zerolog.DebugLevel {
		mask |= logMaskDebug
	}
	if level <= zerolog.WarnLevel {
		mask |= logMaskWarn
	}
	if level <= zerolog.ErrorLevel {
		mask |= logMaskError
	}
	m.enabledMask = mask
}

// traceEnabled reports whether trace logging is enabled.
func (m loggerMixin) traceEnabled() bool {
	return m.enabledMask&logMaskTrace != 0
}

// infoEnabled reports whether info logging is enabled.
func (m loggerMixin) infoEnabled() bool {
	return m.enabledMask&logMaskInfo != 0
}

// debugEnabled reports whether debug logging is enabled.
func (m loggerMixin) debugEnabled() bool {
	return m.enabledMask&logMaskDebug != 0
}

// warnEnabled reports whether warn logging is enabled.
func (m loggerMixin) warnEnabled() bool {
	return m.enabledMask&logMaskWarn != 0
}

// errorEnabled reports whether error logging is enabled.
func (m loggerMixin) errorEnabled() bool {
	return m.enabledMask&logMaskError != 0
}

// errEnabled reports whether error logging is enabled.
func (m loggerMixin) errEnabled() bool {
	return m.errorEnabled()
}

// info builds an info-level log event with trace caller context when enabled.
func (m loggerMixin) info() *zerolog.Event {
	ev := m.logger.WithLevel(zerolog.InfoLevel)
	m.traceCaller(ev, 3)
	return ev
}

// debug builds a debug-level log event with trace caller context when enabled.
func (m loggerMixin) debug() *zerolog.Event {
	ev := m.logger.WithLevel(zerolog.DebugLevel)
	m.traceCaller(ev, 3)
	return ev
}

// warn builds a warn-level log event with trace caller context when enabled.
func (m loggerMixin) warn() *zerolog.Event {
	ev := m.logger.WithLevel(zerolog.WarnLevel)
	m.traceCaller(ev, 3)
	return ev
}

// err builds an error-level log event with the provided error attached.
func (m loggerMixin) err(err error) *zerolog.Event {
	ev := m.logger.Err(err)
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
