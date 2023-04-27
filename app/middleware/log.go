package middleware

import (
	"strings"
	"time"

	"github.com/apicat/apicat/commom/log"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/exp/slog"
)

func RequestIDLog(skip ...string) func(*gin.Context) {
	skipPaths := make(map[string]bool)
	for i := range skip {
		skipPaths[skip[i]] = true
	}
	return func(c *gin.Context) {
		start := time.Now()

		reqid := shortuuid.New()
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(log.ContextLogID(ctx, reqid))
		c.Header("x-apicat-requestid", reqid)

		c.Next()

		path := c.Request.URL.Path
		for k := range skipPaths {
			if strings.HasPrefix(path, k) {
				return
			}
		}

		if raw := c.Request.URL.RawQuery; raw != "" {
			path = path + "?" + raw
		}

		logattrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.Int("size", c.Writer.Size()),
			slog.Duration("latency", time.Since(start)),
		}

		lvl := slog.LevelInfo
		if len(c.Errors) > 0 {
			lvl = slog.LevelError
			logattrs = append(logattrs, slog.String("gin.err", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}
		slog.LogAttrs(c, lvl, "gin.access", logattrs...)
	}
}
