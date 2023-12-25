package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestDuff(t *testing.T) {
	ab, _ := os.ReadFile("../testdata/specdiff_a.json")

	bb, _ := os.ReadFile("../testdata/specdiff_b.json")
	collectitemB, err := Diff(ab, bb)
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

	b, err := DiffSchema([]byte(sa), []byte(sb))
	if err != nil {
		t.Log(err)
	}
	res, _ := json.MarshalIndent(b, "", " ")
	fmt.Println(string(res))
}
