package logger

import (
	"context"
	"io"
	"time"

	"log/slog"
)

func newTraceHandler(w io.Writer, addSource bool, lvl slog.Level) slog.Handler {
	sets := make(map[string]struct{})
	for _, v := range LogPrivacyAttrKey {
		sets[v] = struct{}{}
	}
	opt := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: addSource,
		ReplaceAttr: func(g []string, a slog.Attr) slog.Attr {
			// 修改默认日志时间格式
			if a.Key == slog.TimeKey {
				return slog.String(
					slog.TimeKey,
					a.Value.Time().Format(time.DateTime),
				)
			}
			// 模糊私密信息，防止信息泄露
			if _, ok := sets[a.Key]; ok {
				if a.Value.Kind() == slog.KindString {
					return slog.String(a.Key, replaceLogPrivacyAttrKey(a.Value.String()))
				}
			}
			return a
		},
	}
	return &logTraceHandle{TextHandler: slog.NewTextHandler(w, opt)}
}

var LogPrivacyAttrKey = []string{
	"password",
	"token",
	"secretkey",
	"accesskey",
}

func replaceLogPrivacyAttrKey(s string) string {
	p := []rune(s)
	n := len(p)
	if n == 0 {
		return s
	}
	if n < 3 {
		return string(p[0]) + "*"
	}
	start := n / 3
	for i := start; i < n-start; i++ {
		p[i] = '*'
	}
	return string(p)
}

type logTraceHandle struct {
	*slog.TextHandler
}

func (h *logTraceHandle) Handle(ctx context.Context, r slog.Record) error {
	guid := ctx.Value(ctxidKey)
	if guid != nil {
		r.Add("logid", guid.(string))
	}
	return h.TextHandler.Handle(ctx, r)
}

var ctxidKey = struct{}{}

func ContextLogID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxidKey, id)
}
