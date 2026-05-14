// Package slogctx implements utilities for logging request-scoped data via [context.Context].
package slogctx

import (
	"context"
	"slices"
)

type ctxKey struct{}

// With returns a derived [context.Context] with the given arguments.
// If called with this context, a [slog.Logger] created with [NewHandler] will automatically record these arguments.
// The arguments are processed as if by [slog.Logger.Log].
func With(ctx context.Context, args ...any) context.Context {
	if a, ok := ctx.Value(ctxKey{}).([]any); ok {
		// Clip the slice to ensure a new backing array is allocated.
		return context.WithValue(ctx, ctxKey{}, append(slices.Clip(a), args...))
	}
	return context.WithValue(ctx, ctxKey{}, args)
}
