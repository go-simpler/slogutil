// Package slogattr implements utilities for creating [slog.Attr].
package slogattr

import (
	"cmp"
	"fmt"
	"log/slog"
	"strconv"
)

// ErrorKey is the key used by the [Error] function.
// The associated value is an error.
const ErrorKey = "error"

// Error returns a [slog.Attr] for an error value.
func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any(ErrorKey, err)
}

// Slice returns a [slog.Attr] for a slice of [cmp.Ordered] values.
func Slice[T cmp.Ordered](key string, ts []T) slog.Attr {
	if len(ts) == 0 {
		return slog.Attr{}
	}
	return slog.Any(key, ts)
}

// Stringer returns a [slog.Attr] for an [fmt.Stringer] value.
func Stringer(key string, s fmt.Stringer) slog.Attr {
	return slog.Any(key, s)
}

// LogValuer returns a [slog.Attr] for a [slog.LogValuer] value.
func LogValuer(key string, v slog.LogValuer) slog.Attr {
	return slog.Any(key, v)
}

// LogValuers returns a [slog.Attr] for a slice of [slog.LogValuer] values.
// Given to [slog.JSONHandler], it will be written as
//
//	{"key":{"1":...,"2":...}}
//
// where ... is the [slog.LogValuer] itself.
func LogValuers[T slog.LogValuer](key string, ts []T) slog.Attr {
	// A slog.LogValuer in a slice is not supported by log/slog.
	// As a workaround, we convert the slice into a slog.GroupValue.
	// More:
	// - https://github.com/golang/go/issues/63204
	// - https://github.com/golang/go/issues/71088
	attrs := make([]slog.Attr, len(ts))
	for i := range ts {
		attrs[i] = slog.Any(strconv.Itoa(i+1), ts[i])
	}
	return slog.Attr{Key: key, Value: slog.GroupValue(attrs...)}
}

// Lazy returns a [slog.Attr] for a lazily evaluated [cmp.Ordered] value.
func Lazy[T cmp.Ordered](key string, fn func() T) slog.Attr {
	return slog.Any(key, lazy[T](fn))
}

type lazy[T cmp.Ordered] func() T

// LogValue implements [slog.LogValuer].
func (l lazy[T]) LogValue() slog.Value {
	return slog.AnyValue(l())
}
