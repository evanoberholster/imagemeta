package exif

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
)

func TestLoggerMixinEnabledMask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     zerolog.Level
		trace     bool
		info      bool
		debug     bool
		warn      bool
		errorable bool
	}{
		{
			name:      "trace",
			level:     zerolog.TraceLevel,
			trace:     true,
			info:      true,
			debug:     true,
			warn:      true,
			errorable: true,
		},
		{
			name:      "debug",
			level:     zerolog.DebugLevel,
			trace:     false,
			info:      true,
			debug:     true,
			warn:      true,
			errorable: true,
		},
		{
			name:      "info",
			level:     zerolog.InfoLevel,
			trace:     false,
			info:      true,
			debug:     false,
			warn:      true,
			errorable: true,
		},
		{
			name:      "warn",
			level:     zerolog.WarnLevel,
			trace:     false,
			info:      false,
			debug:     false,
			warn:      true,
			errorable: true,
		},
		{
			name:      "error",
			level:     zerolog.ErrorLevel,
			trace:     false,
			info:      false,
			debug:     false,
			warn:      false,
			errorable: true,
		},
		{
			name:      "fatal",
			level:     zerolog.FatalLevel,
			trace:     false,
			info:      false,
			debug:     false,
			warn:      false,
			errorable: false,
		},
		{
			name:      "disabled",
			level:     zerolog.Disabled,
			trace:     false,
			info:      false,
			debug:     false,
			warn:      false,
			errorable: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := zerolog.New(io.Discard).Level(tt.level)
			m := newLoggerMixin(l)

			if got := m.traceEnabled(); got != tt.trace {
				t.Fatalf("traceEnabled() = %v, want %v", got, tt.trace)
			}
			if got := m.infoEnabled(); got != tt.info {
				t.Fatalf("infoEnabled() = %v, want %v", got, tt.info)
			}
			if got := m.debugEnabled(); got != tt.debug {
				t.Fatalf("debugEnabled() = %v, want %v", got, tt.debug)
			}
			if got := m.warnEnabled(); got != tt.warn {
				t.Fatalf("warnEnabled() = %v, want %v", got, tt.warn)
			}
			if got := m.errorEnabled(); got != tt.errorable {
				t.Fatalf("errorEnabled() = %v, want %v", got, tt.errorable)
			}
			if got := m.errEnabled(); got != tt.errorable {
				t.Fatalf("errEnabled() = %v, want %v", got, tt.errorable)
			}
		})
	}
}

func TestLoggerMixinSetLoggerRefreshesMask(t *testing.T) {
	t.Parallel()

	m := newLoggerMixin(zerolog.New(io.Discard).Level(zerolog.ErrorLevel))
	if !m.errorEnabled() || m.warnEnabled() {
		t.Fatalf("unexpected initial enabled states: error=%v warn=%v", m.errorEnabled(), m.warnEnabled())
	}

	m.setLogger(zerolog.New(io.Discard).Level(zerolog.DebugLevel))
	if !m.debugEnabled() || !m.warnEnabled() || !m.errorEnabled() {
		t.Fatalf("setLogger did not refresh mask: debug=%v warn=%v error=%v", m.debugEnabled(), m.warnEnabled(), m.errorEnabled())
	}
}
