// Package server exposes the local workbench API and embedded React application.
package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apicat/apicat/v3/internal/model"
	"github.com/apicat/apicat/v3/internal/problem"
	"github.com/apicat/apicat/v3/internal/runner"
	"github.com/apicat/apicat/v3/internal/workspace"
	"github.com/apicat/apicat/v3/web"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Start(getDoc func() *model.APIDoc, specDir, host string, port int) error {
	handler, err := newHandler(getDoc, specDir, host)
	if err != nil {
		return err
	}
	addr := net.JoinHostPort(host, fmt.Sprint(port))
	fmt.Printf("Serving API workbench at http://%s\n", addr)
	server := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	return server.ListenAndServe()
}

func newHandler(getDoc func() *model.APIDoc, specDir, host string) (*echo.Echo, error) {
	assets, err := fs.Sub(web.Content, "dist")
	if err != nil {
		return nil, err
	}
	index, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		return nil, fmt.Errorf("build the frontend with npm --prefix web run build: %w", err)
	}
	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(key)
	store := workspace.New(specDir)
	e := echo.New()
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		code := echo.StatusCode(err)
		if code == 0 {
			code = 500
		}
		message := http.StatusText(code)
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
			message = fmt.Sprint(he.Message)
		}
		payload := map[string]any{"error": message}
		var localized *problem.Error
		if errors.As(err, &localized) {
			payload["messageKey"] = localized.Key
			payload["values"] = localized.Values
		}
		_ = c.JSON(code, payload)
	}
	// Every API is same-origin, including the session bootstrap. Host validation
	// prevents a rebinding website from obtaining the local execution capability.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			r := c.Request()
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("Referrer-Policy", "no-referrer")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			if strings.HasPrefix(r.URL.Path, "/api/") {
				c.Response().Header().Set("Cache-Control", "no-store")
				hostname := r.Host
				if h, _, err := net.SplitHostPort(r.Host); err == nil {
					hostname = h
				}
				allowed := hostname == "localhost" || hostname == host
				if ip := net.ParseIP(hostname); ip != nil {
					allowed = allowed || ip.IsLoopback() || host == "0.0.0.0" || host == "::"
				}
				if !allowed {
					return httpProblem(403, "不允许的主机名")
				}
				if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
					return httpProblem(403, "不允许跨站访问本地工作台")
				}
				if origin := r.Header.Get("Origin"); origin != "" {
					u, err := url.Parse(origin)
					if err != nil || u.Host != r.Host || (u.Scheme != "http" && u.Scheme != "https") {
						return httpProblem(403, "不允许跨域访问本地工作台")
					}
				}
				if r.Method != "GET" && r.Method != "HEAD" {
					if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Apicat-Token")), []byte(token)) != 1 {
						return httpProblem(403, "会话已过期，请刷新页面")
					}
					if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
						return httpProblem(415, "请使用 JSON 请求")
					}
					r.Body = http.MaxBytesReader(c.Response(), r.Body, 4<<20)
				}
			}
			return next(c)
		}
	})
	e.GET("/api/session", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"token": token, "specDir": specDir, "workspacePath": store.Path()})
	})
	e.GET("/api/document", func(c *echo.Context) error { return c.JSON(200, getDoc()) })
	e.GET("/api/workspace", func(c *echo.Context) error {
		state, err := store.Load()
		if err != nil {
			return httpProblem(500, "无法读取环境配置，请检查 workspace.json")
		}
		return c.JSON(200, state)
	})
	e.PUT("/api/workspace", func(c *echo.Context) error {
		var state workspace.State
		if err := decode(c, &state); err != nil {
			return err
		}
		if err := workspace.Validate(state); err != nil {
			return echo.NewHTTPError(400, err.Error()).Wrap(err)
		}
		saved, err := store.Save(state)
		if errors.Is(err, workspace.ErrConflict) {
			return echo.NewHTTPError(409, err.Error()).Wrap(err)
		}
		if err != nil {
			return httpProblem(500, "环境配置保存失败，请检查文件权限和可用空间")
		}
		return c.JSON(200, saved)
	})
	e.POST("/api/requests/execute", func(c *echo.Context) error {
		var request runner.Request
		if err := decode(c, &request); err != nil {
			return err
		}
		state, err := store.Load()
		if err != nil {
			return httpProblem(500, "无法读取环境配置")
		}
		response, err := runner.Execute(c.Request().Context(), request, state)
		if err != nil {
			return echo.NewHTTPError(400, err.Error()).Wrap(err)
		}
		return c.JSON(200, response)
	})
	e.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(assets))))
	// Explicit SPA routes keep missing assets/API URLs from returning index.html.
	for _, path := range []string{"/", "/endpoints", "/endpoints/:operationId", "/components", "/environments"} {
		e.GET(path, func(c *echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache")
			return c.HTMLBlob(200, index)
		})
	}
	return e, nil
}
func decode(c *echo.Context, value any) error {
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return httpProblem(400, "请求 JSON 格式或字段无效")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return httpProblem(400, "请求只能包含一个 JSON 对象")
	}
	return nil
}

func httpProblem(status int, key string) error {
	return echo.NewHTTPError(status, key).Wrap(problem.New(key, nil))
}
