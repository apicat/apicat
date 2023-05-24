package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/apicat/apicat/common/spec"
	"github.com/apicat/apicat/models"
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
	p := &models.Projects{}
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
				m.cache.Delete(p.(*models.Projects).ID)
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
	specObj.Definitions.Schemas = models.DefinitionSchemasExport(id)
	specObj.Definitions.Parameters = models.DefinitionParametersExport(id)
	specObj.Definitions.Responses = models.DefinitionResponsesExport(id)
	newcm := specObj.CollectionsMap(true)
	m.cache.Store(id, newcm)
	return newcm
}

func (m *MockServer) matchRoute(c *gin.Context, routes map[string]map[string]spec.HTTPPart) *spec.HTTPPart {
	p := strings.Split(c.Param("path"), "/")
	for path, methods := range routes {
		rp := strings.Split(path, "/")
		if len(rp) != len(p) {
			continue
		}
		// match path
		var flag bool
		for k, v := range rp {
			if v != p[k] && v[0] != '{' {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		// match method
		h, ok := methods[c.Request.Method]
		if ok {
			slog.InfoCtx(c, "find route", slog.String("path", path))
			return &h
		}
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
		c.Writer.Header().Set("Content-Type", k)
		c.Writer.WriteHeader(res.Code)
		c.Writer.Write(responsedata)
		return
	}
	// no content?
	c.Writer.WriteHeader(res.Code)
}
