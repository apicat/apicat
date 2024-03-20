package mock

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/apicat/apicat/backend/module/spec"

	"github.com/gin-gonic/gin"
)

func TestRun(t *testing.T) {
	Run("127.0.0.1:8001", WithApiUrl("http://127.0.0.1:8000"), WithApiPath("api/mock"))
}

func TestRenderData(t *testing.T) {
	m := &MockServer{}
	h := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(h)
	body := []byte(`{
		"code": 200,
		"name": "Response Name",
		"content": {
		  "application/json": {
			"schema": {
			  "type": "object",
			  "x-apicat-orders": [
				"obj1",
				"array_1"
			  ],
			  "properties": {
				"array_1": {
				  "type": "array",
				  "x-apicat-mock": "array",
				  "items": {
					"type": "string",
					"x-apicat-mock": "firstname"
				  }
				},
				"obj1": {
				  "type": "object",
				  "x-apicat-orders": [
					"age",
					"name"
				  ],
				  "properties": {
					"age": {
					  "type": "integer",
					  "x-apicat-mock": "integer"
					},
					"name": {
					  "type": "string",
					  "x-apicat-mock": "string"
					}
				  }
				}
			  },
			  "example": ""
			}
		  }
		}
	  }`)
	res := spec.HTTPResponse{}

	err := json.Unmarshal(body, &res)
	if err != nil {
		t.Error(err)
		return
	}
	m.renderMockResponse(c, res)
	fmt.Println(h.Body.String())
}
