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
	if p, ok := ctx.Value(ctxKey{}).(*payload); ok {
		p.mu.RLock()
		r.Add(p.args...)
		p.mu.RUnlock()
	}
	return h.Handler.Handle(ctx, r)
}
