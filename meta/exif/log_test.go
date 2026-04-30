package exif

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/evanoberholster/imagemeta/imagetype"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	metalog "github.com/evanoberholster/imagemeta/meta/logging"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func TestLoggerMixinEnabledChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     slog.Level
		trace     bool
		info      bool
		debug     bool
		warn      bool
		errorable bool
	}{
		{
			name:      "trace",
			level:     metalog.LevelTrace,
			trace:     true,
			info:      true,
			debug:     true,
			warn:      true,
			errorable: true,
		},
		{
			name:      "debug",
			level:     slog.LevelDebug,
			trace:     false,
			info:      true,
			debug:     true,
			warn:      true,
			errorable: true,
		},
		{
			name:      "info",
			level:     slog.LevelInfo,
			trace:     false,
			info:      true,
			debug:     false,
			warn:      true,
			errorable: true,
		},
		{
			name:      "warn",
			level:     slog.LevelWarn,
			trace:     false,
			info:      false,
			debug:     false,
			warn:      true,
			errorable: true,
		},
		{
			name:      "error",
			level:     slog.LevelError,
			trace:     false,
			info:      false,
			debug:     false,
			warn:      false,
			errorable: true,
		},
		{
			name:      "fatal",
			level:     slog.Level(12),
			trace:     false,
			info:      false,
			debug:     false,
			warn:      false,
			errorable: false,
		},
		{
			name:      "disabled",
			level:     metalog.LevelDisabled,
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
			l := metalog.New(io.Discard, tt.level)
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
			if got := m.logLevelDebug(); got != tt.debug {
				t.Fatalf("logLevelDebug() = %v, want %v", got, tt.debug)
			}
			if got := m.warnEnabled(); got != tt.warn {
				t.Fatalf("warnEnabled() = %v, want %v", got, tt.warn)
			}
			if got := m.logLevelWarn(); got != tt.warn {
				t.Fatalf("logLevelWarn() = %v, want %v", got, tt.warn)
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

func TestLoggerMixinSetLoggerRefreshesChecks(t *testing.T) {
	t.Parallel()

	m := newLoggerMixin(metalog.New(io.Discard, slog.LevelError))
	if !m.errorEnabled() || m.warnEnabled() {
		t.Fatalf("unexpected initial enabled states: error=%v warn=%v", m.errorEnabled(), m.warnEnabled())
	}

	m.setLogger(metalog.New(io.Discard, slog.LevelDebug))
	if !m.debugEnabled() || !m.warnEnabled() || !m.errorEnabled() {
		t.Fatalf("setLogger did not refresh checks: debug=%v warn=%v error=%v", m.debugEnabled(), m.warnEnabled(), m.errorEnabled())
	}
}

func TestLoggerMixinInfoEvent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	m := newLoggerMixin(metalog.New(&buf, slog.LevelInfo))
	if !m.infoEnabled() {
		t.Fatal("info logging should be enabled")
	}

	m.info().Str("phase", "decode").Msg("starting exif decode")

	out := buf.String()
	if !strings.Contains(out, `"level":"info"`) || !strings.Contains(out, `"phase":"decode"`) {
		t.Fatalf("info log output = %q", out)
	}
}

func TestTagLogContextIncludesDecodeFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	r := &Reader{loggerMixin: newLoggerMixin(metalog.New(&buf, slog.LevelWarn))}
	r.po = 42
	r.firstIFDOffset = 8
	r.exifLength = 512
	r.Exif.ImageType = imagetype.ImageJPEG
	r.Exif.IFD0.Make = "Canon"
	r.Exif.IFD0.Model = "EOS R5"

	entry := tag.NewEntry(tag.TagMakerNote, tag.TypeUndefined, 128, 256, tag.ExifIFD, 0, utils.LittleEndian)
	r.tagLogContext(r.warn(), entry).Msg("tag failed")

	out := buf.String()
	for _, want := range []string{
		`"level":"warn"`,
		`"readerOffset":42`,
		`"firstIFDOffset":8`,
		`"exifLength":512`,
		`"imageType":"image/jpeg"`,
		`"cameraMake":"Canon"`,
		`"cameraModel":"EOS R5"`,
		`"tagID":37500`,
		`"tagName":"MakerNote"`,
		`"tagType":7`,
		`"tagTypeName":"UNDEFINED"`,
		`"tagSize":128`,
		`"tagOffset":256`,
		`"tagEmbedded":false`,
		`"byteOrder":"LittleEndian"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
}

func TestRawTagHeaderLogContextIncludesInvalidHeaderFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	r := &Reader{loggerMixin: newLoggerMixin(metalog.New(&buf, slog.LevelWarn))}
	r.po = 14
	r.Exif.ImageType = imagetype.ImageTiff
	directory := tag.NewDirectory(utils.BigEndian, tag.IFD0, 0, 8, 0)

	r.rawTagHeaderLogContext(r.warn(), directory, 3, tag.TagModel, tag.Type(99), 2, 24).
		Msg("invalid exif tag header")

	out := buf.String()
	for _, want := range []string{
		`"ifdType":"IFD0"`,
		`"ifdIndex":0`,
		`"ifdOffset":8`,
		`"tagIndex":3`,
		`"tagID":272`,
		`"tagName":"Model"`,
		`"tagType":99`,
		`"tagTypeName":"UNKNOWN"`,
		`"units":2`,
		`"rawValueOffset":24`,
		`"byteOrder":"BigEndian"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
}
