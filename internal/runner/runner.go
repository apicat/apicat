// Package runner executes explicit HTTP requests. It never changes the spec.
package runner

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/apicat/apicat/v3/internal/workspace"

	"github.com/apicat/apicat/v3/internal/problem"
)

const MaxResponseBytes = 2 << 20

type Pair struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}
type Request struct {
	Method        string `json:"method"`
	URL           string `json:"url"`
	EnvironmentID string `json:"environmentId"`
	Headers       []Pair `json:"headers"`
	Query         []Pair `json:"query"`
	QueryResolved bool   `json:"queryResolved,omitempty"`
	Path          []Pair `json:"path"`
	Body          string `json:"body"`
	TimeoutMS     int    `json:"timeoutMs"`
}
type Response struct {
	Status     int         `json:"status"`
	StatusText string      `json:"statusText"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
	Encoding   string      `json:"encoding"`
	DurationMS int64       `json:"durationMs"`
	Size       int         `json:"size"`
	Truncated  bool        `json:"truncated"`
	URL        string      `json:"url"`
}

var variablePattern = regexp.MustCompile(`\{\{\s*([^{}]+?)\s*\}\}`)

func expand(value string, variables map[string]string) (string, error) {
	var missing string
	result := variablePattern.ReplaceAllStringFunc(value, func(token string) string {
		key := strings.TrimSpace(variablePattern.FindStringSubmatch(token)[1])
		v, ok := variables[key]
		if !ok {
			missing = key
		}
		return v
	})
	if missing != "" {
		return "", problem.New("变量 {name} 未定义或未启用", map[string]any{"name": missing})
	}
	return result, nil
}

// Build resolves workspace variables once; request drafts remain unmodified.
func Build(ctx context.Context, input Request, state workspace.State) (*http.Request, error) {
	variables := map[string]string{}
	for _, v := range state.Variables {
		if v.Enabled {
			variables[v.Key] = v.Value
		}
	}
	base := ""
	if input.EnvironmentID != "" {
		found := false
		for _, env := range state.Environments {
			if env.ID == input.EnvironmentID {
				found = true
				base = env.BaseURL
				for _, v := range env.Variables {
					if v.Enabled {
						variables[v.Key] = v.Value
					}
				}
			}
		}
		if !found {
			return nil, problem.New("所选环境不存在，请刷新环境配置", nil)
		}
	}
	target, err := expand(input.URL, variables)
	if err != nil {
		return nil, err
	}
	for _, p := range input.Path {
		if p.Enabled {
			value, err := expand(p.Value, variables)
			if err != nil {
				return nil, err
			}
			if value == "" && strings.Contains(target, "{"+p.Key+"}") {
				return nil, problem.New("路径参数 {name} 不能为空", map[string]any{"name": p.Key})
			}
			target = strings.ReplaceAll(target, "{"+p.Key+"}", url.PathEscape(value))
		}
	}
	if strings.ContainsAny(target, "{}") {
		return nil, problem.New("请填写所有路径参数或环境变量", nil)
	}
	u, err := url.Parse(target)
	if err != nil {
		return nil, problem.New("请求地址格式无效", nil)
	}
	if !u.IsAbs() {
		base, err = expand(base, variables)
		if err != nil {
			return nil, err
		}
		if base == "" {
			return nil, problem.New("请选择并配置环境地址，或输入完整的 HTTP(S) 地址", nil)
		}
		b, err := url.Parse(base)
		if err != nil || b.Host == "" || b.RawQuery != "" || b.Fragment != "" {
			return nil, problem.New("环境地址格式无效（不支持 query 或 fragment）", nil)
		}
		u, err = url.Parse(strings.TrimRight(base, "/") + "/" + strings.TrimLeft(target, "/"))
		if err != nil {
			return nil, err
		}
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil {
		return nil, problem.New("仅支持不含用户名密码的 HTTP(S) 地址", nil)
	}
	q := u.Query()
	for _, p := range input.Query {
		if p.Enabled && p.Key != "" {
			key, value := p.Key, p.Value
			if !input.QueryResolved {
				key, err = expand(key, variables)
				if err != nil {
					return nil, err
				}
				value, err = expand(value, variables)
				if err != nil {
					return nil, err
				}
			}
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()
	u.Fragment = ""
	method := strings.ToUpper(input.Method)
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "TRACE":
	default:
		return nil, problem.New("不支持的 HTTP 方法", nil)
	}
	body, err := expand(input.Body, variables)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	for _, p := range input.Headers {
		if p.Enabled && p.Key != "" {
			key, err := expand(p.Key, variables)
			if err != nil {
				return nil, err
			}
			value, err := expand(p.Value, variables)
			if err != nil {
				return nil, err
			}
			if strings.ContainsAny(key+value, "\r\n") {
				return nil, problem.New("请求头不能包含换行", nil)
			}
			switch strings.ToLower(key) {
			case "host":
				req.Host = value
			case "content-length", "transfer-encoding", "connection":
				return nil, problem.New("请求头 {name} 由 HTTP 客户端管理", map[string]any{"name": key})
			default:
				req.Header.Add(key, value)
			}
		}
	}
	return req, nil
}
func Execute(ctx context.Context, input Request, state workspace.State) (*Response, error) {
	timeout := input.TimeoutMS
	if timeout == 0 {
		timeout = 30000
	}
	if timeout < 100 || timeout > 120000 {
		return nil, problem.New("超时必须在 100–120000 毫秒之间", nil)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	req, err := Build(ctx, input, state)
	if err != nil {
		return nil, err
	}
	// No shared cookie jar, application-level retry loop, or automatic redirects.
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	started := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, problem.New("请求已取消或超时", nil)
		}
		// Do not expose resolved credentials or query secrets in transport errors.
		return nil, problem.New("请求失败，请检查目标地址、网络连接和 TLS 证书", nil)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBytes+1))
	if err != nil {
		return nil, problem.New("读取响应失败或请求超时", nil)
	}
	result := &Response{Status: resp.StatusCode, StatusText: http.StatusText(resp.StatusCode), Headers: resp.Header, DurationMS: time.Since(started).Milliseconds(), URL: req.URL.String(), Encoding: "text"}
	if len(data) > MaxResponseBytes {
		result.Truncated = true
		data = data[:MaxResponseBytes]
	}
	result.Size = len(data)
	if utf8.Valid(data) && !strings.ContainsRune(string(data), '\x00') {
		result.Body = string(data)
	} else {
		result.Encoding = "base64"
		result.Body = base64.StdEncoding.EncodeToString(data)
	}
	return result, nil
}
