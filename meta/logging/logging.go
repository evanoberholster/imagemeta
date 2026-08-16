package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var (
	// LevelTrace enables trace-level logging.
	LevelTrace slog.Level = -8
	// LevelDisabled disables all logging by setting a very high minimum level.
	LevelDisabled slog.Level = 12
	// Logger is a shared slog logger.
	//
	// Deprecated: use SetLogger and GetLogger instead of assigning directly.
	Logger = New(os.Stdout, LevelDisabled)
	// loggerPtr stores the effective process-wide logger used by meta packages.
	loggerPtr atomic.Pointer[slog.Logger]
)

func init() {
	loggerPtr.Store(Logger)
}

// SetLogger sets the global logger used by meta packages.
func SetLogger(l *slog.Logger) {
	l = normalizeLogger(l)
	Logger = l
	loggerPtr.Store(l)
}

// GetLogger returns the global logger used by meta packages.
func GetLogger() *slog.Logger {
	if l := loggerPtr.Load(); l != nil {
		return l
	}
	l := normalizeLogger(Logger)
	loggerPtr.Store(l)
	return l
}

// Mixin provides common logger state and level checks for parser types.
type Mixin struct {
	logger    *slog.Logger
	component string
}

// NewMixin initializes a reusable parser logger helper.
func NewMixin(l *slog.Logger) Mixin {
	return Mixin{logger: resolveLogger(l)}
}

// NewComponentMixin initializes a reusable parser logger helper with a stable component field.
func NewComponentMixin(l *slog.Logger, component string) Mixin {
	return Mixin{logger: resolveLogger(l), component: component}
}

// SetLogger replaces the logger used by this helper.
func (m *Mixin) SetLogger(l *slog.Logger) {
	m.logger = resolveLogger(l)
}

// Enabled reports whether events at level should be emitted.
func (m Mixin) Enabled(level slog.Level) bool {
	return LevelEnabled(m.eventLogger(), level)
}

// TraceEnabled reports whether trace-level callsite fields should be added.
func (m Mixin) TraceEnabled() bool {
	return TraceEnabled(m.eventLogger())
}

// Event builds an event at level and adds trace callsite context when enabled.
func (m Mixin) Event(level slog.Level, depth int) *Event {
	return ComponentEvent(m.eventLogger(), m.component, level, depth)
}

// DebugEnabled reports whether debug events should be emitted.
func (m Mixin) DebugEnabled() bool {
	return m.Enabled(slog.LevelDebug)
}

// InfoEnabled reports whether info events should be emitted.
func (m Mixin) InfoEnabled() bool {
	return m.Enabled(slog.LevelInfo)
}

// WarnEnabled reports whether warn events should be emitted.
func (m Mixin) WarnEnabled() bool {
	return m.Enabled(slog.LevelWarn)
}

// ErrorEnabled reports whether error events should be emitted.
func (m Mixin) ErrorEnabled() bool {
	return m.Enabled(slog.LevelError)
}

// Debug builds a debug-level event with optional trace caller context.
func (m Mixin) Debug(depth int) *Event {
	return m.Event(slog.LevelDebug, depth)
}

// Info builds an info-level event with optional trace caller context.
func (m Mixin) Info(depth int) *Event {
	return m.Event(slog.LevelInfo, depth)
}

// Warn builds a warn-level event with optional trace caller context.
func (m Mixin) Warn(depth int) *Event {
	return m.Event(slog.LevelWarn, depth)
}

// Error builds an error-level event with optional trace caller context.
func (m Mixin) Error(depth int) *Event {
	return m.Event(slog.LevelError, depth)
}

func (m Mixin) eventLogger() *slog.Logger {
	if m.logger != nil {
		return m.logger
	}
	return GetLogger()
}

// LevelEnabled reports whether logger emits events at level.
func LevelEnabled(l *slog.Logger, level slog.Level) bool {
	return resolveLogger(l).Enabled(context.Background(), level)
}

// TraceEnabled reports whether trace callsite fields should be added.
func TraceEnabled(l *slog.Logger) bool {
	return LevelEnabled(l, LevelTrace)
}

// Event builds an event at level and adds trace callsite context when enabled.
func NewEvent(l *slog.Logger, level slog.Level, depth int) *Event {
	ev := newEvent(resolveLogger(l), level, "")
	TraceCaller(l, ev, depth+1)
	return ev
}

// ComponentEvent builds an event at level and adds a stable component field when provided.
func ComponentEvent(l *slog.Logger, component string, level slog.Level, depth int) *Event {
	ev := newEvent(resolveLogger(l), level, component)
	TraceCaller(l, ev, depth+1)
	return ev
}

// TraceCaller annotates ev with caller function information when trace is enabled.
func TraceCaller(l *slog.Logger, ev *Event, depth int) {
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

// New creates a JSON slog logger with lower-case level names.
func New(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				if lv, ok := a.Value.Any().(slog.Level); ok {
					a.Value = slog.StringValue(strings.ToLower(lv.String()))
				}
			case slog.MessageKey:
				a.Key = "msg"
			}
			return a
		},
	}))
}

func resolveLogger(l *slog.Logger) *slog.Logger {
	if l != nil {
		return l
	}
	return GetLogger()
}

func normalizeLogger(l *slog.Logger) *slog.Logger {
	if l != nil {
		return l
	}
	if current := loggerPtr.Load(); current != nil {
		return current
	}
	if Logger != nil {
		return Logger
	}
	return New(os.Stdout, LevelDisabled)
}

// Event is a chainable slog event builder compatible with the existing log call style.
type Event struct {
	logger *slog.Logger
	level  slog.Level
	attrs  []slog.Attr
	buf    [16]slog.Attr
}

func newEvent(l *slog.Logger, level slog.Level, component string) *Event {
	e := &Event{logger: l, level: level}
	e.attrs = e.buf[:0]
	if component != "" {
		e.attrs = append(e.attrs, slog.String("component", component))
	}
	return e
}

func (e *Event) addAttr(a slog.Attr) *Event {
	e.attrs = append(e.attrs, a)
	return e
}

func (e *Event) Str(key, value string) *Event {
	return e.addAttr(slog.String(key, value))
}

func (e *Event) Strs(key string, values []string) *Event {
	return e.addAttr(slog.Any(key, values))
}

func (e *Event) Stringer(key string, value fmt.Stringer) *Event {
	if value == nil {
		return e.addAttr(slog.Any(key, nil))
	}
	return e.addAttr(slog.String(key, value.String()))
}

func (e *Event) Int(key string, value int) *Event {
	return e.addAttr(slog.Int(key, value))
}

func (e *Event) Int8(key string, value int8) *Event {
	return e.addAttr(slog.Int64(key, int64(value)))
}

func (e *Event) Int64(key string, value int64) *Event {
	return e.addAttr(slog.Int64(key, value))
}

func (e *Event) Uint8(key string, value uint8) *Event {
	return e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *Event) Uint16(key string, value uint16) *Event {
	return e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *Event) Uint32(key string, value uint32) *Event {
	return e.addAttr(slog.Uint64(key, uint64(value)))
}

func (e *Event) Uint64(key string, value uint64) *Event {
	return e.addAttr(slog.Uint64(key, value))
}

func (e *Event) Bool(key string, value bool) *Event {
	return e.addAttr(slog.Bool(key, value))
}

func (e *Event) Dur(key string, value time.Duration) *Event {
	return e.addAttr(slog.Duration(key, value))
}

func (e *Event) Err(err error) *Event {
	if err != nil {
		return e.addAttr(slog.Any("error", err))
	}
	return e
}

func (e *Event) Object(key string, value any) *Event {
	return e.addAttr(slog.Any(key, marshalObjectAny(value)))
}

func (e *Event) Array(key string, value any) *Event {
	if m, ok := value.(interface{ MarshalLogArray(*Array) }); ok {
		arr := &Array{}
		m.MarshalLogArray(arr)
		return e.addAttr(slog.Any(key, arr.values))
	}
	return e.addAttr(slog.Any(key, value))
}

func (e *Event) Msg(msg string) {
	e.logger.LogAttrs(context.Background(), e.level, msg, e.attrs...)
}

func marshalObjectAny(value any) any {
	if value == nil {
		return nil
	}
	if m, ok := value.(interface{ MarshalLogObject(*Event) }); ok {
		ev := &Event{}
		ev.attrs = ev.buf[:0]
		m.MarshalLogObject(ev)
		return logObject{attrs: ev.attrs}
	}
	return value
}

type logObject struct {
	attrs []slog.Attr
}

func (o logObject) LogValue() slog.Value {
	return slog.GroupValue(o.attrs...)
}

// Array is a chainable helper for structured array log fields.
type Array struct {
	values []any
}

func (a *Array) Object(value any) *Array {
	a.values = append(a.values, marshalObjectAny(value))
	return a
}
