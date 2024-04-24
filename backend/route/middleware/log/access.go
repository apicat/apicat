package log

import (
	"strings"
	"time"

	"github.com/apicat/apicat/v2/backend/utils/logger"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

const HeaderKey = "X-APICAT-REQUESTID"

func AccessLog(skip ...string) func(*gin.Context) {
	skipPaths := make(map[string]bool)
	for i := range skip {
		skipPaths[skip[i]] = true
	}
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		for k := range skipPaths {
			if strings.HasPrefix(path, k) {
				return
			}
		}

		start := time.Now()
		reqid := shortuuid.New()
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(logger.ContextLogID(ctx, reqid))
		c.Header(HeaderKey, reqid)

		c.Next()

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
