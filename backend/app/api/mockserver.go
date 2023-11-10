package api

import (
	"encoding/json"
	"fmt"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/definition"
	"github.com/apicat/apicat/backend/model/project"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/datagen"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type MockServer struct {
	cache sync.Map
}

func NewMockServer() *MockServer {
	return &MockServer{
		cache: sync.Map{},
	}
}

// Handler requests to process fake data
// [method] /mock/{id}/path*
func (m *MockServer) Handler(c *gin.Context) {
	p := &project.Projects{}
	slog.InfoCtx(c, "mock", slog.String("path", c.Param("path")))
	if err := p.Get(c.Param("id")); err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	routes := m.getRequestRoutesSchemaOrCache(p.ID)
	part := m.matchRoute(c, routes)
	if part == nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	rescode, _ := strconv.Atoi(c.Query("mock_response_code"))
	var index int
	for i, v := range part.Responses {
		if rescode > 0 {
			if v.Code == rescode {
				index = i
				break
			}
		} else {
			if v.Code == 200 {
				index = i
				break
			}
			// match 201,204....
			if v.Code > 200 && v.Code < 300 {
				index = i
			}
		}
	}
	if len(part.Responses) > 0 {
		res := part.Responses[index]
		slog.InfoCtx(c, "schema response", slog.String("name", res.Name))
		m.renderMockResponse(c, res)
	}
}

func (m *MockServer) ClearCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			p, ok := ctx.Get("CurrentProject")
			if ok {
				ctx.Next()
				m.cache.Delete(p.(*project.Projects).ID)
			}
		}
	}
}

func (m *MockServer) getRequestRoutesSchemaOrCache(id uint) map[string]map[string]spec.HTTPPart {
	cm, ok := m.cache.Load(id)
	if ok {
		return cm.(map[string]map[string]spec.HTTPPart)
	}
	specObj := &spec.Spec{}
	specObj.Definitions.Schemas = definition.DefinitionSchemasExport(id)
	specObj.Definitions.Parameters = definition.DefinitionParametersExport(id)
	specObj.Definitions.Responses = definition.DefinitionResponsesExport(id)
	specObj.Collections = collection.CollectionsExport(id)
	newcm := specObj.CollectionsMap(true, 3)
	m.cache.Store(id, newcm)
	return newcm
}

func (m *MockServer) matchRoute(c *gin.Context, routes map[string]map[string]spec.HTTPPart) *spec.HTTPPart {
	p := strings.Split(c.Param("path"), "/")
	matched := map[string]struct {
		vars int
		data spec.HTTPPart
	}{}
	for path, methods := range routes {
		rp := strings.Split(path, "/")
		if len(rp) != len(p) {
			continue
		}
		// match path
		var flag bool
		var hasVar int
		for k, v := range rp {
			if v != p[k] {
				if len(v) > 0 && v[0] == '{' {
					hasVar++
				} else {
					flag = true
					break
				}
			}
		}
		if flag {
			continue
		}
		// match method
		h, ok := methods[strings.ToLower(c.Request.Method)]
		if ok {
			if hasVar > 0 {
				matched[path] = struct {
					vars int
					data spec.HTTPPart
				}{hasVar, h}
			} else {
				slog.InfoCtx(c, "find route", slog.String("path", path), slog.String("mockpath", c.Param("path")))
				return &h
			}
		}
	}
	for path, v := range matched {
		slog.InfoCtx(c, "find route", slog.String("path", path), slog.String("mockpath", c.Param("path")))
		return &v.data
	}
	return nil
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
					// 如果非必填 则不一定返回他
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
