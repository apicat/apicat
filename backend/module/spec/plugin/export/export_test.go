package export

import (
	"fmt"
	"github.com/apicat/apicat/backend/module/spec"
	"os"
	"testing"
)

func TestMd(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/spec.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	s, err := spec.ParseJSON(raw)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	b, err := Markdown(s)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(string(b))
}

func TestHTML(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/spec.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	s, err := spec.ParseJSON(raw)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	b, err := HTML(s)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := os.WriteFile("aa.html", b, os.ModePerm); err != nil {
		t.Fail()
	}
}
