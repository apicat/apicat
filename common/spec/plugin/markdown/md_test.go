package markdown

import (
	"fmt"
	"os"
	"testing"

	"github.com/apicat/apicat/common/spec"
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
	b, err := Encode(s)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(string(b))
}
