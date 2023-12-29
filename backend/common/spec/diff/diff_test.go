package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
)

func TestDuff(t *testing.T) {
	a, _ := os.ReadFile("../testdata/specdiff_a.json")

	b, _ := os.ReadFile("../testdata/specdiff_b.json")

	ab, _ := spec.ParseJSON(a)
	bb, _ := spec.ParseJSON(b)

	collectitemB, err := Diff(ab.Collections[0], bb.Collections[0])
	if err != nil {
		t.Log(err)
	}
	// aaa, _ := json.MarshalIndent(collectitemA, "", " ")
	res, _ := json.MarshalIndent(collectitemB, "", " ")
	fmt.Println(string(res))
}

func TestSchemaDiff(t *testing.T) {
	sa := `{
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
	  }`

	sb := `{
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
	  }`

	a := &jsonschema.Schema{}
	_ = json.Unmarshal([]byte(sa), a)
	b := &jsonschema.Schema{}
	_ = json.Unmarshal([]byte(sb), b)

	b, err := DiffSchema(a, b)
	if err != nil {
		t.Log(err)
	}
	res, _ := json.MarshalIndent(b, "", " ")
	fmt.Println(string(res))
}

func TestGetMapOneCollectionMap(t *testing.T) {
	ab, _ := os.ReadFile("../testdata/self_to_self.json")

	source, err := spec.ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	a, au := getMapOne(source.CollectionsMap(true, 1))

	fmt.Println(a)
	fmt.Println(au)

	b, err := json.Marshal(a)

	if err != nil {
		t.Log(err)
	}
	fmt.Println(string(b))
}
