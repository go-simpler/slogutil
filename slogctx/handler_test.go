package slogctx_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"go-simpler.org/slogutil/slogctx"
)

func TestHandler(t *testing.T) {
	replaceAttr := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
			return slog.Attr{}
		}
		return a
	}

	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: replaceAttr})
	l := slog.New(slogctx.NewHandler(h))

	ctx := t.Context()
	ctx = slogctx.With(ctx, "x", 1)
	foo(ctx, l)

	got := "\n" + buf.String()
	want := `
msg="hello from foo" x=1
msg="hello from bar" x=1 y=2
`
	if got != want {
		t.Errorf("\ngot: %s\nwant: %s", got, want)
	}
}

func foo(ctx context.Context, l *slog.Logger) {
	l.InfoContext(ctx, "hello from foo")
	ctx = slogctx.With(ctx, "y", 2)
	bar(ctx, l)
}

func bar(ctx context.Context, l *slog.Logger) {
	l.InfoContext(ctx, "hello from bar")
}

// goos: darwin
// goarch: arm64
// pkg: go-simpler.org/slogutil/slogctx
// cpu: Apple M1 Pro
// BenchmarkHandler/enabled-8              205089428                5.726 ns/op           0 B/op          0 allocs/op
// BenchmarkHandler/disabled-8             262692470                4.568 ns/op           0 B/op          0 allocs/op
func BenchmarkHandler(b *testing.B) {
	b.Run("enabled", func(b *testing.B) {
		benchmarkHandler(b, true)
	})
	b.Run("disabled", func(b *testing.B) {
		benchmarkHandler(b, false)
	})
}

func benchmarkHandler(b *testing.B, enabled bool) {
	b.Helper()
	b.ReportAllocs()

	ctx := b.Context()
	ctx = slogctx.With(ctx,
		"foo", 1,
		"bar", 2,
		"baz", 3,
	)

	h := slog.DiscardHandler
	if enabled {
		h = slogctx.NewHandler(h)
	}

	l := slog.New(h)
	for b.Loop() {
		l.InfoContext(ctx, "")
	}
}
