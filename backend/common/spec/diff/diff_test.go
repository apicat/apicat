package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/apicat/apicat/backend/common/spec"
)

func TestDuff(t *testing.T) {
	ab, _ := os.ReadFile("../testdata/specdiff_a.json")
	a, _ := spec.ParseJSON(ab)

	bb, _ := os.ReadFile("../testdata/specdiff_b.json")
	b, _ := spec.ParseJSON(bb)
	_, collectitemB := Diff(a, b)
	// aaa, _ := json.MarshalIndent(collectitemA, "", " ")
	bbb, _ := json.MarshalIndent(collectitemB, "", " ")
	fmt.Println(string(bbb))
}
