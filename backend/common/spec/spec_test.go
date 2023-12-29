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
		t.Log(err)
		t.Fail()
	}
	spec, err := ParseJSON(raw)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	w, err := spec.ToJSON(JSONOption{Indent: "  "})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(string(w))
}

func TestDereferenceSchema(t *testing.T) {

	ab, _ := os.ReadFile("../testdata/self_to_self.json")

	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	parent := source.Definitions.Schemas.LookupID(2068)
	sub := source.Definitions.Schemas.LookupID(2332)

	err = parent.DereferenceSchema(sub)
	if err != nil {
		fmt.Println(err)
		return
	}

	bs, _ := json.MarshalIndent(parent, "", " ")

	fmt.Println(string(bs))
}

func TestDereferenceSelf(t *testing.T) {

	ab, _ := os.ReadFile("../testdata/self_to_self.json")

	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	onlySelf := source.Definitions.Schemas.LookupID(2332)

	onlySelf.DereferenceSelf()

	bs, _ := json.MarshalIndent(onlySelf, "", " ")

	fmt.Println(string(bs))

}
