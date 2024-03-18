package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"apicat-cloud/backend/module/spec"

	"github.com/apicat/datagen"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type MockServer struct {
	ApiUrl  string
	ApiPath string
}

func NewMockServer(Options ...Option) *MockServer {
	m := &MockServer{
		ApiPath: "api/mock",
	}
	for _, Option := range Options {
		Option(m)
	}

	if !isValidURL(m.ApiUrl) {
		panic(fmt.Sprintf("mock: apiurl is not valid: -- %s --", m.ApiPath))
	}

	return m
}

func isValidURL(testURL string) bool {
	u, err := url.Parse(testURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Handler requests to process fake data
// [method] /mock/{id}/path*
func (m *MockServer) Handler(c *gin.Context) {
	id := c.Param("projectID")
	path := c.Param("path")
	if id == "" || path == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": "path or id is empty"})
		return
	}
	if c.Request.URL.Query().Encode() != "" {
		// this request has query
		i := strings.LastIndex(path, "?")
		if i != -1 {
			path = path[:i]
		}
	}

	// request to apicat
	// default is https://0.0.0.0:8000/mock/{id}/{path}
	targetURL := fmt.Sprintf("%s/%s/%s", m.ApiUrl, m.ApiPath, fmt.Sprintf("%s%s", id, path))

	// generate request
	req, err := http.NewRequest(c.Request.Method, targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get data"})
		return
	}

	// send request
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to get data: %s", err.Error())})
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Read data error"})
		return
	}
	resps := spec.HTTPResponses{}
	err = json.Unmarshal(body, &resps)
	if err != nil {
		c.JSON(http.StatusBadRequest, string(body))
		return
	}
	if len(resps) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not have this interface"})
		return
	}
	resp := resps[0]
	rs := resp.ItemsTreeToList()
	if len(rs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This interface is empty category"})
		return
	}
	resp.HTTPResponseDefine = *rs[0]
	m.renderMockResponse(c, *resp)
}

func (m *MockServer) renderMockResponse(c *gin.Context, res spec.HTTPResponse) {
	// find first contentType
	for k, v := range res.Content {
		b, _ := json.Marshal(v.Schema)
		responsedata, err := datagen.JSONSchemaGen(b, &datagen.GenOption{
			DatagenKey: "x-apicat-mock",
		})
		if err != nil {
			slog.ErrorCtx(c, "datagen jsonschema gen faild", slog.String("err", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if res.Header != nil {
			for _, h := range res.Header {
				if !h.Required {
					// random boolean, to generate not required
					if !datagen.Boolean() {
						continue
					}
				}
				hb, _ := json.Marshal(h.Schema)
				headerdata, err := datagen.JSONSchemaGen(hb, &datagen.GenOption{
					DatagenKey: "x-apicat-mock",
				})
				if err != nil {
					continue
				}
				// Allow access to this response header
				c.Writer.Header().Add("Access-Control-Expose-Headers", h.Name)
				c.Header(h.Name, fmt.Sprintf("%v", headerdata))
			}
		}
		c.Header("Content-Type", k)
		c.Writer.WriteHeader(res.Code)
		json.NewEncoder(c.Writer).Encode(responsedata) // nolint
		return
	}
	// no content?
	c.Writer.WriteHeader(res.Code)
}

// @addr is the ip and port of the service, default is 127.0.0.1:8001
// @apiOptions is the config options of the mock server
func Run(addr string, apiOptions ...Option) {

	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.ContextWithFallback = true

	mocksrv := NewMockServer(apiOptions...)

	if mocksrv == nil {
		panic("init mock error")
	}

	r.Use(gin.Recovery(), CORSMiddleware())

	r.Any("/mock/:projectID/*path", mocksrv.Handler)

	if addr == "" {
		panic(fmt.Sprintf("mock addr is empty: -- %s --", addr))
	}

	// if port is used, panic
	// if mock panic, no need to restart yet
	err := r.Run(addr)
	// mock server error
	if err != nil {
		panic(fmt.Sprintf("mock server error: %s", err))
	}
}

// CORSMiddleware Middleware for handling cross requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE,PATCH")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
