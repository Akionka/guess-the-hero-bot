package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type TelegoContextHandler struct {
	slog.Handler
}

var _ slog.Handler = (*TelegoContextHandler)(nil)

func (h *TelegoContextHandler) Handle(ctx context.Context, r slog.Record) error {
	thCtx, ok := ctx.(*th.Context)
	if ok {
		r.AddAttrs(slog.Int("update_id", thCtx.UpdateID()))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *TelegoContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TelegoContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *TelegoContextHandler) WithGroup(name string) slog.Handler {
	return &TelegoContextHandler{Handler: h.Handler.WithGroup(name)}
}

type TelegoLogger struct {
	*slog.Logger
	LogErrors bool
	LogDebug  bool
}

func (s *TelegoLogger) Debugf(format string, args ...any) {
	if !s.LogDebug {
		return
	}
	s.Debug(fmt.Sprintf(format, args...), opAttr("slogLogger.Debugf"))
}

func (s *TelegoLogger) Errorf(format string, args ...any) {
	if !s.LogErrors {
		return
	}
	s.Error(fmt.Sprintf(format, args...), opAttr("slogLogger.Errorf"))
}

var _ telego.Logger = (*TelegoLogger)(nil)
