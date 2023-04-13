package openapi

import (
	"fmt"
	"os"
	"testing"

	"github.com/apicat/apicat/commom/spec"
)

func TestDecode(t *testing.T) {
	fs := map[string]string{
		"swagger":    "../../testdata/swagger.json",
		"openapi3.0": "../../testdata/openapi3.0.yaml",
		"openapi3.1": "../../testdata/openapi3.1.yaml",
	}
	for k, v := range fs {
		raw, err := os.ReadFile(v)
		if err != nil {
			t.Fatal(k, err)
		}
		if _, err := Decode(raw); err != nil {
			t.Fatal(k, err)
		}
	}
}

func TestEncode(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/spec.json")
	if err != nil {
		t.Fail()
	}
	p, err := spec.ParseJSON(raw)
	if err != nil {
		t.Fail()
	}

	for _, v := range []string{"3.0"} {
		p, err := Encode(p, v)
		if err != nil {
			t.Fatal(v, err)
		}
		fmt.Println(string(p))
	}

}
