package logging

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestComponentMixinEventIncludesComponent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	m := NewComponentMixin(New(&buf, slog.LevelInfo), "unit")
	m.Event(slog.LevelInfo, 1).
		Str("k", "v").
		Msg("hello")

	out := buf.String()
	for _, want := range []string{
		`"component":"unit"`,
		`"k":"v"`,
		`"msg":"hello"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
}

func TestComponentMixinSetLoggerRefreshesWriter(t *testing.T) {
	t.Parallel()

	var buf1 bytes.Buffer
	var buf2 bytes.Buffer
	m := NewComponentMixin(New(&buf1, slog.LevelInfo), "unit")
	m.Event(slog.LevelInfo, 1).Msg("first")

	m.SetLogger(New(&buf2, slog.LevelInfo))
	m.Event(slog.LevelInfo, 1).Msg("second")

	if strings.Contains(buf1.String(), `"msg":"second"`) {
		t.Fatalf("first logger unexpectedly received second message: %q", buf1.String())
	}
	if !strings.Contains(buf2.String(), `"msg":"second"`) {
		t.Fatalf("second logger missing second message: %q", buf2.String())
	}
}

func TestEventStrsLogsStringArray(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	NewEvent(New(&buf, slog.LevelInfo), slog.LevelInfo, 1).
		Strs("brands", []string{"avif", "heic"}).
		Msg("brands")

	out := buf.String()
	if !strings.Contains(out, `"brands":["avif","heic"]`) {
		t.Fatalf("log output missing string array: %q", out)
	}
}

func TestMixinLevelHelpers(t *testing.T) {
	t.Parallel()

	m := NewComponentMixin(New(io.Discard, slog.LevelWarn), "unit")
	if m.DebugEnabled() {
		t.Fatal("debug should be disabled")
	}
	if m.InfoEnabled() {
		t.Fatal("info should be disabled")
	}
	if !m.WarnEnabled() {
		t.Fatal("warn should be enabled")
	}
	if !m.ErrorEnabled() {
		t.Fatal("error should be enabled")
	}
}

func TestSetLoggerUpdatesGlobalLogger(t *testing.T) {
	t.Parallel()

	oldLogger := GetLogger()
	t.Cleanup(func() { SetLogger(oldLogger) })

	var buf bytes.Buffer
	SetLogger(New(&buf, slog.LevelInfo))

	ComponentEvent(nil, "unit", slog.LevelInfo, 1).
		Str("k", "v").
		Msg("hello")

	out := buf.String()
	for _, want := range []string{
		`"component":"unit"`,
		`"k":"v"`,
		`"msg":"hello"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %s: %q", want, out)
		}
	}
}
