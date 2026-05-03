package slogctx

import (
	"context"
	"sync"
)

type ctxKey struct{}

type payload struct {
	mu   sync.RWMutex
	args []any
}

// With returns a [context.Context] with the given arguments attached.
// If called with this context, a [slog.Logger] created with [NewHandler] will automatically record these arguments.
// The arguments are processed as if by [slog.Logger.Log].
func With(ctx context.Context, args ...any) context.Context {
	p, ok := ctx.Value(ctxKey{}).(*payload)
	if !ok {
		a := make([]any, 0, max(10, len(args)))
		a = append(a, args...)
		return context.WithValue(ctx, ctxKey{}, &payload{args: a})
	}

	p.mu.Lock()
	p.args = append(p.args, args...)
	p.mu.Unlock()

	return ctx
}
