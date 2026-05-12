package slogattr_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"go-simpler.org/slogutil/slogattr"
)

func TestAll(t *testing.T) {
	replaceAttr := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		return a
	}

	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: replaceAttr})
	l := slog.New(h)

	l.Info("",
		slogattr.Error(errors.New("oops")),
		slogattr.Slice("slice", []int{1, 2, 3}),
		slogattr.Stringer("stringer", slog.LevelInfo),
		slogattr.LogValuer("log_valuer", point{1, 2}),
		slogattr.LogValuers("log_valuers", []point{{1, 2}, {3, 4}}),
		slogattr.Lazy("lazy", func() string { return "foo" + "bar" }),
	)

	var got bytes.Buffer
	json.Indent(&got, buf.Bytes(), "", "  ")

	const want = `{
  "error": "oops",
  "slice": [
    1,
    2,
    3
  ],
  "stringer": "INFO",
  "log_valuer": {
    "x": 1,
    "y": 2
  },
  "log_valuers": {
    "1": {
      "x": 1,
      "y": 2
    },
    "2": {
      "x": 3,
      "y": 4
    }
  },
  "lazy": "foobar"
}
`
	if got.String() != want {
		t.Errorf("\ngot: %s\nwant: %s", &got, want)
	}
}

type point struct{ x, y int }

func (p point) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("x", p.x),
		slog.Int("y", p.y),
	)
}
