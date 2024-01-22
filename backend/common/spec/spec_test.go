package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParseSpec(t *testing.T) {
	raw, err := os.ReadFile("./testdata/spec.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	spec, err := ParseJSON(raw)
	if err != nil {
		t.Fatal("parse json error", err)
	}
	w, err := spec.ToJSON(JSONOption{Indent: "  "})
	if err != nil {
		t.Fatal("spec to json error", err)
	}
	t.Log(string(w))
}

func TestDereferenceSchema(t *testing.T) {

	ab, err := os.ReadFile("./testdata/self_to_self.json")
	if err != nil {
		t.Fatal("read file error", err)
	}

	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	parent := source.Definitions.Schemas.LookupID(2068)
	sub := source.Definitions.Schemas.LookupID(2332)

	parent.DereferenceSchema(sub)

	bs, _ := json.MarshalIndent(parent, "", " ")

	t.Log(string(bs))
}

func TestUnparkDereferenceSchema(t *testing.T) {
	ab, err := os.ReadFile("./testdata/ref_all_dereference.json")
	if err != nil {
		t.Fatal("read file error", err)
	}

	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	// a := source.Definitions.Schemas.LookupID(2336)
	// b := source.Definitions.Schemas.LookupID(2340)
	// c := source.Definitions.Schemas.LookupID(2341)
	s := `[{"id":2,"name":"schema-2","type":"schema","parentid":6,"schema":{"type":"object","x-apicat-orders":["fuck1","fuck2"],"properties":{"fuck1":{"type":"string","x-apicat-mock":"string"},"fuck2":{"$ref":"#/definitions/schemas/1"}}}},{"id":11,"name":"schema-1","type":"schema","schema":{"type":"object","x-apicat-orders":["fuck"],"properties":{"fuck":{"type":"string","x-apicat-mock":"string"}}}},{"id":1,"name":"schema-1","type":"schema","parentid":5,"schema":{"type":"object","x-apicat-orders":["1","2","3"],"properties":{"1":{"type":"string","x-apicat-mock":"string"},"2":{"type":"string","x-apicat-mock":"string"},"3":{"type":"string","x-apicat-mock":"string"}},"required":["1","2"]}},{"id":3,"name":"schema-3","type":"schema","schema":{"$ref":"#/definitions/schemas/1"}},{"id":4,"name":"schema-4","type":"schema","schema":{"$ref":"#/definitions/schemas/4"}}]`
	sub := Schemas{}
	err = json.Unmarshal([]byte(s), &sub)
	if err != nil {
		fmt.Println(err)
	}

	for _, cc := range source.Collections {
		cc.UnpackDereferenceSchema(sub)
	}

	// a.UnpackDereferenceSchema(sub)

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))
}

func TestRemoveSchema(t *testing.T) {
	ab, err := os.ReadFile("./testdata/self_to_self.json")
	if err != nil {
		t.Fatal("read file error", err)
	}

	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	parent := source.Definitions.Schemas.LookupID(2068)
	// sub := source.Definitions.Schemas.LookupID(2332)

	// parent.RemoveSchema(sub)

	// bs, _ := json.MarshalIndent(parent, "", " ")

	// fmt.Println(string(bs))

	for _, c := range source.Collections {
		c.RemoveSchema(parent)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))

}

func TestDereferenceSelf(t *testing.T) {

	ab, err := os.ReadFile("./testdata/self_to_self.json")
	if err != nil {
		t.Fatal("read file error", err)
	}

	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	onlySelf := source.Definitions.Schemas.LookupID(2332)

	onlySelf.dereferenceSelf()

	bs, _ := json.MarshalIndent(onlySelf, "", " ")

	t.Log(string(bs))

}

func TestResponseRef(t *testing.T) {

	ab, err := os.ReadFile("./testdata/response_ref.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	resp := source.Definitions.Responses.LookupID(378)
	// resp1 := source.Definitions.Responses.LookupID(378)
	// rr := []*HTTPResponseDefine{
	// 	resp, resp1,
	// }
	for _, c := range source.Collections {
		c.DereferenceResponse(resp)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))

}

func TestResponseRemove(t *testing.T) {
	ab, err := os.ReadFile("./testdata/response_ref.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	resp := source.Definitions.Responses.LookupID(378)

	for _, c := range source.Collections {
		c.RemoveResponse(resp.ID)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))

}

func TestDereferenceGlobalParameters(t *testing.T) {

	ab, err := os.ReadFile("./testdata/global_excepts.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	for _, c := range source.Collections {
		c.DereferenceGlobalParameters("header", source.Globals.Parameters.Header.LookupID(31))
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))

}

func TestAddParameters(t *testing.T) {

	ab, err := os.ReadFile("./testdata/global_excepts.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	for _, c := range source.Collections {
		c.AddParameters("header", source.Globals.Parameters.Header.LookupID(31))
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))

}

func TestOpenGlobalParameters(t *testing.T) {

	ab, err := os.ReadFile("./testdata/global_excepts.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}

	for _, c := range source.Collections {
		c.OpenGlobalParameters("path", source.Globals.Parameters.Path.LookupID(34).ID)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	t.Log(string(bs))

}

func TestSItemsListToTree(t *testing.T) {
	ab, err := os.ReadFile("./testdata/items_tree_export_openapi.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}
	ss := source.Definitions.Schemas
	ts := Schemas{}
	for _, v := range ss {
		ts = append(ts, v.ItemsTreeToList()...)
	}
	res := ts.ItemsListToTree()
	bs, _ := json.MarshalIndent(res, "", " ")

	t.Log(string(bs))
}

func TestSMakeSelfTree(t *testing.T) {
	ab, err := os.ReadFile("./testdata/items_tree_export_openapi.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}
	ss := source.Definitions.Schemas
	ts := Schemas{}
	for _, v := range ss {
		ts = append(ts, v.ItemsTreeToList()...)
	}
	s := ts.LookupID(2342)
	s2 := s.makeSelfTree(s.Schema.XCategory, map[string]*Schema{})
	bs, _ := json.MarshalIndent(s2, "", " ")

	t.Log(string(bs))
}

func TestRItemsListToTree(t *testing.T) {
	ab, err := os.ReadFile("./testdata/items_tree_export_openapi.json")
	if err != nil {
		t.Fatal("read file error", err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		t.Fatal("parse json error", err)
	}
	ss := source.Definitions.Responses
	ts := HTTPResponseDefines{}
	for _, v := range ss {
		ts = append(ts, v.ItemsTreeToList()...)
	}
	res := ts.ItemsListToTree()
	bs, _ := json.MarshalIndent(res, "", " ")

	t.Log(string(bs))
}
