package openapi

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/apicat/apicat/backend/common/spec"
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
		if x, err := Decode(raw); err != nil {
			t.Fatal(k, err)
		} else {
			d, _ := x.ToJSON(spec.JSONOption{Indent: ""})
			t.Log(string(d))
		}
	}
}

// func TestEncode(t *testing.T) {
// 	raw, err := os.ReadFile("../../testdata/spec.json")
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	p, err := spec.ParseJSON(raw)
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
	a, err := os.ReadFile("../../testdata/items_tree_export_openapi.json")
	if err != nil {
		t.Fatal("read file error", err)
	}

	ab, err := spec.ParseJSON(a)
	if err != nil {
		t.Fatal("parse json spec error", err)
	}

	b, err := Encode(ab, "3.1.0")
	if err != nil {
		t.Fatal("encode spec error", err)
	}
	t.Log(string(b))
}

func TestToApiCat(t *testing.T) {
	a, err := os.ReadFile("../../testdata/items_tree_import_openapi.json")
	if err != nil {
		t.Fatal("read file error", err)
	}

	b, err := Decode(a)
	if err != nil {
		t.Fatal("decode spec error", err)
	}

	bs, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		t.Fatal("marshal error", err)
	}

	t.Log(string(bs))
}

func TestComponentsParamenters(t *testing.T) {
	a, _ := os.ReadFile("../../testdata/components_parameters.yml")

	b, err := Decode(a)
	if err != nil {
		t.Fatal("decode spec error", err)
	}

	bs, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		t.Fatal("marshal error", err)
	}

	t.Log(string(bs))

}
func TestCollectionDereferenceSchema(t *testing.T) {
	ab, err := os.ReadFile("../../testdata/openapi3.1.yaml")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := Decode(ab)
	if err != nil {
		t.Fatal("decode spec error", err)
	}
	s := source.Definitions.Schemas.Lookup("User")
	if s == nil {
		t.Fatal("not found this schema : ", s)
	}
	for _, c := range source.Collections {
		c.DereferenceSchema(s)
	}

	// b, err := Encode(source, "3.1.0")
	// if err != nil {
	// 	t.Fatal("encode spec error", err)
	// }

	// if collection.id is 0, openapi.path.OperatorID will be error

	// bs, err := json.MarshalIndent(source, "", " ")
	// if err != nil {
	// 	t.Fatal("marshal error", err)
	// }

	// t.Log(string(bs))

}

func BenchmarkCollectionDereferenceSchema(b *testing.B) {
	ab, err := os.ReadFile("../../testdata/openapi3.1.yaml")
	if err != nil {
		b.Fatal("read file error", err)
	}
	source, err := Decode(ab)
	if err != nil {
		b.Fatal("decode spec error", err)
	}
	s := source.Definitions.Schemas.Lookup("User")
	if s == nil {
		b.Fatal("not found this schema : ", s)
	}
	for i := 0; i < b.N; i++ {
		for _, c := range source.Collections {
			c.DereferenceSchema(s)
		}
	}
}
