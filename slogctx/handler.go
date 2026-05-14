package slogctx

import (
	"context"
	"log/slog"
)

type handler struct {
	slog.Handler
}

// NewHandler creates a [slog.Handler] that automatically records arguments attached to the context via [With].
func NewHandler(h slog.Handler) slog.Handler {
	return &handler{Handler: h}
}

// Handle implements [slog.Handler].
func (h *handler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic // hugeParam: can't change the signature.
	if args, ok := ctx.Value(ctxKey{}).([]any); ok {
		r.Add(args...)
	}
	return h.Handler.Handle(ctx, r)
}
