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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := metalog.New(io.Discard, tt.level)
			m := metalog.NewComponentMixin(l, "exif")

			if got := m.TraceEnabled(); got != tt.trace {
				t.Fatalf("TraceEnabled() = %v, want %v", got, tt.trace)
			}
			if got := m.InfoEnabled(); got != tt.info {
				t.Fatalf("InfoEnabled() = %v, want %v", got, tt.info)
			}
			if got := m.DebugEnabled(); got != tt.debug {
				t.Fatalf("DebugEnabled() = %v, want %v", got, tt.debug)
			}
			if got := m.WarnEnabled(); got != tt.warn {
				t.Fatalf("WarnEnabled() = %v, want %v", got, tt.warn)
			}
			if got := m.ErrorEnabled(); got != tt.errorable {
				t.Fatalf("ErrorEnabled() = %v, want %v", got, tt.errorable)
			}
		})
	}
}

func TestLoggerMixinSetLoggerRefreshesChecks(t *testing.T) {
	t.Parallel()

	m := metalog.NewComponentMixin(metalog.New(io.Discard, slog.LevelError), "exif")
	if !m.ErrorEnabled() || m.WarnEnabled() {
		t.Fatalf("unexpected initial enabled states: error=%v warn=%v", m.ErrorEnabled(), m.WarnEnabled())
	}

	m.SetLogger(metalog.New(io.Discard, slog.LevelDebug))
	if !m.DebugEnabled() || !m.WarnEnabled() || !m.ErrorEnabled() {
		t.Fatalf("SetLogger did not refresh checks: debug=%v warn=%v error=%v", m.DebugEnabled(), m.WarnEnabled(), m.ErrorEnabled())
	}
}

func TestLoggerMixinInfoEvent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	m := metalog.NewComponentMixin(metalog.New(&buf, slog.LevelInfo), "exif")
	if !m.InfoEnabled() {
		t.Fatal("info logging should be enabled")
	}

	m.Info(3).Str("phase", "decode").Msg("starting exif decode")

	out := buf.String()
	if !strings.Contains(out, `"level":"info"`) || !strings.Contains(out, `"phase":"decode"`) {
		t.Fatalf("info log output = %q", out)
	}
}

func TestTagLogContextIncludesDecodeFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	r := &Reader{Mixin: metalog.NewComponentMixin(metalog.New(&buf, slog.LevelWarn), "exif")}
	r.po = 42
	r.firstIFDOffset = 8
	r.exifLength = 512
	r.Exif.ImageType = imagetype.ImageJPEG

	entry := tag.NewEntry(tag.TagMakerNote, tag.TypeUndefined, 128, 256, tag.ExifIFD, 0, utils.LittleEndian)
	r.tagLogContext(r.Warn(3), entry).Msg("tag failed")

	out := buf.String()
	for _, want := range []string{
		`"level":"warn"`,
		`"readerOffset":42`,
		`"firstIFDOffset":8`,
		`"exifLength":512`,
		`"imageType":"image/jpeg"`,
		`"tagID":"0x927C"`,
		`"tagName":"MakerNote"`,
		`"tagType":"UNDEFINED"`,
		`"tagSize":128`,
		`"tagOffset":256`,
		`"tagEmbedded":false`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
}

func TestRawTagHeaderLogContextIncludesInvalidHeaderFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	r := &Reader{Mixin: metalog.NewComponentMixin(metalog.New(&buf, slog.LevelWarn), "exif")}
	r.po = 14
	r.Exif.ImageType = imagetype.ImageTiff
	directory := tag.NewDirectory(utils.BigEndian, tag.IFD0, 0, 8, 0)

	r.rawTagHeaderLogContext(r.Warn(3), directory, 3, tag.TagModel, tag.Type(99), 2, 24).
		Msg("invalid exif tag header")

	out := buf.String()
	for _, want := range []string{
		`"ifdType":"IFD0"`,
		`"ifdIndex":0`,
		`"ifdOffset":8`,
		`"tagIndex":3`,
		`"tagID":"0x0110"`,
		`"tagName":"Model"`,
		`"tagType":"UNKNOWN"`,
		`"units":2`,
		`"rawValueOffset":24`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
}

func TestInfoDirectoryLogContextIncludesByteOrderWhenSet(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	r := &Reader{Mixin: metalog.NewComponentMixin(metalog.New(&buf, slog.LevelInfo), "exif")}
	r.po = 31
	directory := tag.NewDirectory(utils.LittleEndian, tag.IFD0, 0, 8, 0)

	r.infoDirectoryLogContext(r.Info(3), directory).Msg("directory lifecycle")

	out := buf.String()
	if !strings.Contains(out, `"byteOrder":"LE"`) {
		t.Fatalf("log output missing byteOrder for known endianness: %q", out)
	}
}

func TestInfoDirectoryLogContextOmitsByteOrderWhenUnknown(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	r := &Reader{Mixin: metalog.NewComponentMixin(metalog.New(&buf, slog.LevelInfo), "exif")}
	r.po = 31
	directory := tag.NewDirectory(utils.UnknownEndian, tag.IFD0, 0, 8, 0)

	r.infoDirectoryLogContext(r.Info(3), directory).Msg("directory lifecycle")

	out := buf.String()
	if strings.Contains(out, `"byteOrder":`) {
		t.Fatalf("log output unexpectedly includes byteOrder for unknown endianness: %q", out)
	}
}
