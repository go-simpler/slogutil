package slogctx_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"go-simpler.org/slogctx"
)

func TestHandler(t *testing.T) {
	replaceAttr := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	}

	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: replaceAttr})
	l := slog.New(slogctx.NewHandler(h))

	ctx := t.Context()
	ctx = slogctx.With(ctx)

	foo(ctx, l)
	l.InfoContext(ctx, "got foo bar")

	got := "\n" + buf.String()
	want := `
level=INFO msg="adding foo"
level=INFO msg="adding bar" foo=1
level=INFO msg="got foo bar" foo=1 bar=2
`
	if got != want {
		t.Errorf("\ngot: %s\nwant: %s", got, want)
	}
}

func foo(ctx context.Context, l *slog.Logger) {
	l.InfoContext(ctx, "adding foo")
	ctx = slogctx.With(ctx, "foo", 1)
	bar(ctx, l)
}

func bar(ctx context.Context, l *slog.Logger) {
	l.InfoContext(ctx, "adding bar")
	ctx = slogctx.With(ctx, "bar", 2)
}

// goos: darwin
// goarch: arm64
// pkg: go-simpler.org/slogctx
// cpu: Apple M1 Pro
// BenchmarkHandler/enabled-8              204999273                5.743 ns/op           0 B/op          0 allocs/op
// BenchmarkHandler/disabled-8             261334262                4.591 ns/op           0 B/op          0 allocs/op
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
