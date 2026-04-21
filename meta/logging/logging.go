package logging

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is a zerolog logger
	Logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.PanicLevel)
)

// WithComponent returns a child logger annotated with a stable component name.
func WithComponent(l zerolog.Logger, component string) zerolog.Logger {
	if component == "" {
		return l
	}
	return l.With().Str("component", component).Logger()
}

// ComponentLogger returns a child logger from the shared package logger.
func ComponentLogger(component string) zerolog.Logger {
	return WithComponent(Logger, component)
}

// Mixin provides common logger state and level checks for parser types.
type Mixin struct {
	logger    zerolog.Logger
	component string
}

// NewMixin initializes a reusable parser logger helper.
func NewMixin(l zerolog.Logger) Mixin {
	return Mixin{logger: l}
}

// NewComponentMixin initializes a reusable parser logger helper with a stable component field.
func NewComponentMixin(l zerolog.Logger, component string) Mixin {
	return Mixin{logger: l, component: component}
}

// SetLogger replaces the logger used by this helper.
func (m *Mixin) SetLogger(l zerolog.Logger) {
	m.logger = l
}

// Level returns the configured zerolog level.
func (m Mixin) Level() zerolog.Level {
	return m.logger.GetLevel()
}

// Enabled reports whether events at level should be emitted.
func (m Mixin) Enabled(level zerolog.Level) bool {
	return LevelEnabled(m.logger, level)
}

// TraceEnabled reports whether trace-level callsite fields should be added.
func (m Mixin) TraceEnabled() bool {
	return TraceEnabled(m.logger)
}

// Event builds an event at level and adds trace callsite context when enabled.
func (m Mixin) Event(level zerolog.Level, depth int) *zerolog.Event {
	return ComponentEvent(m.logger, m.component, level, depth)
}

// LevelEnabled reports whether logger emits events at level.
func LevelEnabled(l zerolog.Logger, level zerolog.Level) bool {
	return l.GetLevel() <= level
}

// TraceEnabled reports whether trace callsite fields should be added.
func TraceEnabled(l zerolog.Logger) bool {
	return l.GetLevel() == zerolog.TraceLevel
}

// Event builds an event at level and adds trace callsite context when enabled.
func Event(l zerolog.Logger, level zerolog.Level, depth int) *zerolog.Event {
	ev := l.WithLevel(level)
	TraceCaller(l, ev, depth+1)
	return ev
}

// ComponentEvent builds an event at level and adds a stable component field when provided.
func ComponentEvent(l zerolog.Logger, component string, level zerolog.Level, depth int) *zerolog.Event {
	if component == "" {
		return Event(l, level, depth)
	}
	return Event(WithComponent(l, component), level, depth)
}

// TraceCaller annotates ev with caller function information when trace is enabled.
func TraceCaller(l zerolog.Logger, ev *zerolog.Event, depth int) {
	if !TraceEnabled(l) {
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
