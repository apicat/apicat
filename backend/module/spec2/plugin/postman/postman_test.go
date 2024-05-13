package postman

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/twitter-postman.json")
	if err != nil {
		t.FailNow()
	}
	x, err := Import(raw)
	if err != nil {
		t.FailNow()
	}
	b, _ := json.MarshalIndent(x, "", "  ")
	fmt.Println(string(b))
}
