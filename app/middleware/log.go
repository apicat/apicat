package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/apicat/apicat/common/log"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/exp/slog"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type CreateProject struct {
	Title string `json:"title"`
}

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

		// Read the Body content
		var bodyByReq []byte
		var request string
		if c.Request.Body != nil {
			bodyByReq, _ = io.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByReq))

		blw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw

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

		if path == "/api/projects/" && c.Request.Method == "POST" {
			createProject := CreateProject{}
			if err := json.Unmarshal(bodyByReq, &createProject); err == nil {
				requestJson, _ := json.Marshal(createProject)
				request = string(requestJson)
			}
		} else {
			request = string(bodyByReq)
		}

		logattrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
			slog.Int("size", c.Writer.Size()),
			slog.Duration("latency", time.Since(start)),
			slog.String("request", request),
			slog.String("response", blw.body.String()),
		}

		lvl := slog.LevelInfo
		if len(c.Errors) > 0 {
			lvl = slog.LevelError
			logattrs = append(logattrs, slog.String("gin.err", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}
		slog.LogAttrs(c, lvl, "gin.access", logattrs...)
	}
}
