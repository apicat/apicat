package diff

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/apicat/apicat/backend/common/spec"
)

func TestDiff(t *testing.T) {
	a, _ := os.ReadFile("../testdata/specdiff_a.json")

	// b, _ := os.ReadFile("../testdata/specdiff_b.json")

	ab, _ := spec.ParseJSON(a)
	// bb, _ := spec.ParseJSON(b)

	cs := `{
		"type": "http",
		"id": 4216,
		"title": "Unnamed interface",
		"content": [
		  {
			"type": "apicat-http-url",
			"attrs": {
			  "path": "",
			  "method": "get"
			}
		  },
		  {
			"type": "apicat-http-request",
			"attrs": {
			  "globalExcepts": {
				"cookie": [],
				"header": [],
				"path": [],
				"query": []
			  },
			  "parameters": {
				"query": [],
				"path": [],
				"cookie": [],
				"header": []
			  }
			}
		  },
		  {
			"type": "apicat-http-response",
			"attrs": {
			  "list": [
				{
				  "code": 200,
				  "name": "Response Name",
				  "content": {
					"application/json": {
					  "schema": {
						"type": "object",
						"example": ""
					  }
					}
				  }
				}
			  ]
			}
		  }
		]
	  }`
	nullc := &spec.CollectItem{}
	err := json.Unmarshal([]byte(cs), nullc)
	if err != nil {
		t.Errorf("unmarshal error: %v", err)
	}

	collectitemB, err := Diff(nullc, ab.Collections[0])
	if err != nil {
		t.Errorf("diff error: %v", err)
	}
	// aaa, _ := json.MarshalIndent(collectitemA, "", " ")
	res, _ := json.MarshalIndent(collectitemB, "", " ")
	t.Log(string(res))
}

func TestSchemaDiff(t *testing.T) {
	sa := `{
		"id": 2360,
        "name": "tt",
		"schema": {
			"type": "object",
			"x-apicat-orders": [
			  "sex",
			  "money",
			  "test"
			],
			"properties": {
			  "money": {
				"type": "string",
				"x-apicat-mock": "string"
			  },
			  "sex": {
				"type": "string",
				"x-apicat-mock": "string"
			  },
			  "test": {
				"type": "object",
				"x-apicat-orders": [
				  "test_a"
				],
				"properties": {
				  "test_a": {
					"type": "string",
					"x-apicat-mock": "string"
				  }
				}
			  }
			},
			"example": ""
		  }
		}
		`
	sb := `{
		"id": 2361,
        "name": "tt2",
		"schema": {
			"type": "object",
			"x-apicat-orders": [
			  "name",
			  "money",
			  "test"
			],
			"properties": {
			  "money": {
				"type": "interger",
				"x-apicat-mock": "interger"
			  },
			  "name": {
				"type": "string",
				"x-apicat-mock": "string"
			  },
			  "test": {
				"type": "array"
			  }
			},
			"example": ""
		  }
		}
		`

	a := &spec.Schema{}
	_ = json.Unmarshal([]byte(sa), a)
	b := &spec.Schema{}
	_ = json.Unmarshal([]byte(sb), b)

	b, err := DiffSchema(a, b)
	if err != nil {
		t.Errorf("diffschema error: %v", err)
	}
	res, _ := json.MarshalIndent(b, "", " ")
	t.Log(string(res))
}

func TestGetMapOneCollectionMap(t *testing.T) {
	ab, _ := os.ReadFile("../testdata/self_to_self.json")

	source, err := spec.ParseJSON(ab)
	if err != nil {
		t.Errorf("parse source error: %v", err)
	}

	a, _ := getMapOne(source.CollectionsMap(true, 1))

	b, err := json.Marshal(a)

	if err != nil {
		t.Errorf("marshal error: %v", err)
	}
	t.Log(string(b))
}

func TestIsChangedBasic(t *testing.T) {
	as := `{
		"type": "http",
      "id": 4223,
      "title": "importent_info_changed",
	  "content": [{
		"type": "apicat-http-request",
          "attrs": {
		  "globalExcepts": {
			"cookie": [],
			"header": [],
			"path": [],
			"query": []
		  },
		  "parameters": {
			"query": [
				{
					"name": "query_1",
					"schema": {
					  "type": "string",
					  "x-apicat-mock": "string"
					}
				  },
			  {
				"name": "query_2",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  }
			],
			"path": [],
			"cookie": [
			  {
				"name": "cookie_1",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  },
			  {
				"name": "cookie_2",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  }
			],
			"header": [
			  {
				"name": "header_1",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "name",
				  "description": "这是名称",
				  "example": "李四"
				}
			  },
			  {
				"name": "header_2",
				"schema": {
				  "type": "integer",
				  "x-apicat-mock": "integer",
				  "description": "这是年龄",
				  "example": "18"
				}
			  }
			]
		  },
		  "content": {
			"application/json": {
			  "schema": {
				"type": "object",
				"x-apicat-orders": [
				  "json_1",
				  "json_obj"
				],
				"properties": {
				  "json_1": {
					"type": "string",
					"x-apicat-mock": "string"
				  },
				  "json_obj": {
					"type": "object",
					"x-apicat-orders": [
					  "name",
					  "age"
					],
					"properties": {
					  "age": {
						"type": "integer",
						"description": "这是年龄",
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
		}}]}
	  `
	bs := `{
		"type": "http",
		"id": 4223,
		"title": "importent_info_changed",
		"content": [{
			"type": "apicat-http-request",
          "attrs": {
		  "globalExcepts": {
			"cookie": [],
			"header": [],
			"path": [],
			"query": []
		  },
		  "parameters": {
			"query": [
			  {
				"name": "query_1",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  },
			  {
				"name": "query_2",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  }
			],
			"path": [],
			"cookie": [
			  {
				"name": "cookie_1",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  },
			  {
				"name": "cookie_2",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "string"
				}
			  }
			],
			"header": [
			  {
				"name": "header_1",
				"schema": {
				  "type": "string",
				  "x-apicat-mock": "name",
				  "description": "这是名称",
				  "example": "张三"
				}
			  },
			  {
				"name": "header_2",
				"schema": {
				  "type": "integer",
				  "x-apicat-mock": "integer",
				  "description": "这是年龄",
				  "example": "18"
				}
			  }
			]
		  },
		  "content": {
			"application/json": {
			  "schema": {
				"type": "object",
				"x-apicat-orders": [
				  "json_1",
				  "json_obj"
				],
				"properties": {
				  "json_1": {
					"type": "string",
					"x-apicat-mock": "string"
				  },
				  "json_obj": {
					"type": "object",
					"x-apicat-orders": [
					  "name",
					  "age"
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
		}}]}
	  `
	ar := &spec.CollectItem{}
	br := &spec.CollectItem{}
	err := json.Unmarshal([]byte(as), ar)
	if err != nil {
		t.Errorf("Marshal ar error: %s", err)
	}
	err = json.Unmarshal([]byte(bs), br)
	if err != nil {
		t.Errorf("Marshal br error: %s", err)
	}
	ok, err := IsChangedBasic(ar, br)
	if err != nil {
		t.Errorf("IsChangedBasic error: %s", err)
	}
	t.Log(ok)
}
