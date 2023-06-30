package tracedlogging

import (
	"context"
	"log/slog"
)

type ctxTraceIdKey struct{}

type ContextHandler struct {
	handler slog.Handler
}

func NewContextHandler(h slog.Handler) *ContextHandler {
	if lh, ok := h.(*ContextHandler); ok {
		h = lh.handler
	}
	return &ContextHandler{h}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return h.handler.Handle(ctx, r)
	}
	if tid, ok := ctx.Value(ctxTraceIdKey{}).(string); ok {
		traceAttr := slog.Attr{
			Key:   "trace_id",
			Value: slog.StringValue(tid),
		}
		r.AddAttrs(traceAttr)
	}
	return h.handler.Handle(ctx, r)
}
func SetTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxTraceIdKey{}, id)
}
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewContextHandler(h.handler.WithAttrs(attrs))
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return NewContextHandler(h.handler.WithGroup(name))
}
