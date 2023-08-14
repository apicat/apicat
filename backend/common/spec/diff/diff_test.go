package diff

import (
	"fmt"
	"testing"

	"github.com/apicat/apicat/backend/common/spec"
)

func TestDuff(t *testing.T) {
	a := spec.Spec{
		// 依赖项
		Globals:     spec.Global{},
		Definitions: spec.Definitions{},
		// 集合里包含一个需要对比的接口
		Collections: []*spec.CollectItem{
			{Type: spec.ContentItemTypeHttp, Title: "aa"},
		},
	}

	b := spec.Spec{
		// 依赖项
		Globals:     spec.Global{},
		Definitions: spec.Definitions{},
		// 集合里包含一个需要对比的接口
		Collections: []*spec.CollectItem{
			{Type: spec.ContentItemTypeHttp, Title: "aa"},
		},
	}
	collectitemA, collectitemB := Diff(&a, &b, true)
	fmt.Println(collectitemA, collectitemB)
}
