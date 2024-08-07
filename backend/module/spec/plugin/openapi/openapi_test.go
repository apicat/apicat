package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/apicat/apicat/v2/backend/module/spec"
)

func TestDecode(t *testing.T) {
	fs := map[string]string{
		// "swagger": "../../testdata/swagger.json",
		// "openapi3.0": "../../testdata/openapi3.0.yaml",
		// "openapi3.1": "../../testdata/openapi3.1.yaml",
		"openpai": "../../testdata/openapi3-examples.json",
	}
	for k, v := range fs {
		raw, err := os.ReadFile(v)
		if err != nil {
			t.Fatal(k, err)
		}
		if x, err := Parse(raw); err != nil {
			t.Fatal(k, err)
		} else {
			d, _ := x.ToJSON(spec.JSONOption{Indent: ""})
			fmt.Println(string(d))
		}
	}
}

// func TestEncode(t *testing.T) {
// 	raw, err := os.ReadFile("../../testdata/spec.json")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	p, err := spec.NewSpecFromJson(raw)
// 	if err != nil {
// 		t.FailNow()
// 	}

// 	fmt.Printf("%+v\n", p)

// 	for _, v := range []string{"3.1.0"} {
// 		p, err := Encode(p, v)
// 		if err != nil {
// 			t.Fatal(v, err)
// 		}
// 		fmt.Println(string(p))
// 	}

// }

func TestToOpenapi(t *testing.T) {
	a, _ := os.ReadFile("../../testdata/items_tree_export_openapi.json")

	ab, _ := spec.NewSpecFromJson(a)

	b, err := Generate(ab, "3.1.0", "json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	// s, err := json.MarshalIndent(b, "", " ")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(s))
}

func TestToApiCat(t *testing.T) {
	a, _ := os.ReadFile("../../testdata/items_tree_import_openapi.json")

	b, err := Parse(a)
	if err != nil {
		fmt.Println(err)
	}

	bs, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bs))
}

func TestComponentsParamenters(t *testing.T) {
	a, _ := os.ReadFile("../../testdata/components_parameters.yml")

	b, err := Parse(a)
	if err != nil {
		fmt.Println(err)
	}

	bs, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bs))
}
