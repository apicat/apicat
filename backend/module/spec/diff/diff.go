package diff

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"

	"golang.org/x/exp/slices"
)

const (
	DIFF_NEW    = "+"
	DIFF_REMOVE = "-"
	DIFF_UPDATE = "!"
)

// Diff 比较两个接口的差异
// 返回对比后的两个接口 其中只有最新的那个 也就是target里边会通过x-apicat-diff标记是否有差异
// 差异并不包含排序
func Diff(original, target *spec.Collection) error {
	if original == nil || target == nil {
		return errors.New("source,target Collections length error")
	}

	copy_original := &spec.Collection{}
	b, err := json.Marshal(original)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, copy_original)
	if err != nil {
		return err
	}

	for _, originalNode := range copy_original.Content {
		for _, targetNode := range target.Content {
			if originalNode.Node.NodeType() == targetNode.Node.NodeType() {
				// assertion in three parts to diff
				switch originalNode.Node.NodeType() {
				case spec.NODE_HTTP_URL:
					originalUrl := originalNode.ToHttpUrl()
					targetUrl := targetNode.ToHttpUrl()
					if originalUrl.Attrs.Path != targetUrl.Attrs.Path {
						if originalUrl.Attrs.Path == "" {
							targetUrl.SetXDiff(DIFF_NEW)
						} else if targetUrl.Attrs.Path == "" {
							targetUrl.Attrs.Path = originalUrl.Attrs.Path
							targetUrl.SetXDiff(DIFF_REMOVE)
						} else {
							targetUrl.SetXDiff(DIFF_UPDATE)
						}
					}
				case spec.NODE_HTTP_REQUEST:
					originalReq := originalNode.ToHttpRequest()
					targetReq := targetNode.ToHttpRequest()
					compareRequest(originalReq, targetReq)
				case spec.NODE_HTTP_RESPONSE:
					originalRes := originalNode.ToHttpResponse()
					targetRes := targetNode.ToHttpResponse()
					compareResponse(originalRes, targetRes)
				default:
					return errors.New("node type error")
				}
			}
		}
	}
	return nil
}

func DiffModel(original, target *spec.DefinitionModel) error {
	if original == nil || target == nil {
		return errors.New("model is nil")
	}
	copy_original := &spec.DefinitionModel{}
	bb, err := json.Marshal(original)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bb, copy_original)
	if err != nil {
		return err
	}
	// if a.Name != b.Name {
	// 	b.XDiff = &DIFF_UPDATE
	// }
	compareJsonSchema(copy_original.Schema, target.Schema)
	return nil
}

func compareHttpParameters(a *spec.HTTPParameters, b *spec.HTTPParameters) {
	aMap := a.ToMap()
	bMap := b.ToMap()
	for _, typ := range spec.HttpParameterType {
		aParamList, aExist := aMap[typ]
		bParamList, bExist := bMap[typ]

		if !aExist && !bExist {
			continue
		}
		if aExist && !bExist {
			aParamList.SetXDiff(DIFF_REMOVE)
			bParamList = append(bParamList, aParamList...)
			continue
		}
		if !aExist && bExist {
			bParamList.SetXDiff(DIFF_NEW)
			continue
		}

		newv := compareParameterList(aParamList, bParamList)
		switch typ {
		case "path":
			b.Path = newv
		case "header":
			b.Header = newv
		case "query":
			b.Query = newv
		case "cookie":
			b.Cookie = newv
		}
	}
}

func compareParameterList(a, b spec.ParameterList) spec.ParameterList {
	names := map[string]struct{}{}
	for _, v := range a {
		names[v.Name] = struct{}{}
	}
	for _, v := range b {
		names[v.Name] = struct{}{}
	}

	for v := range names {
		ap := a.FindByName(v)
		bp := b.FindByName(v)

		if ap == nil && bp == nil {
			continue
		}
		if ap == nil && bp != nil {
			bp.SetXDiff(DIFF_NEW)
			continue
		}
		if ap != nil && bp == nil {
			ap.SetXDiff(DIFF_REMOVE)
			b = append(b, ap)
			continue
		}
		compareParameter(ap, bp)
	}
	return b
}

func compareParameter(a, b *spec.Parameter) bool {
	change := false
	if !compareParamBasicInfo(a, b, false) {
		b.SetXDiff(DIFF_UPDATE)
		change = true
	}
	sc := compareJsonSchema(a.Schema, b.Schema)
	return change || sc
}

func compareContent(a, b spec.HTTPBody) spec.HTTPBody {
	names := map[string]struct{}{}
	aBody, bBody := &spec.Body{}, &spec.Body{}
	// httpbody just have one value
	for v := range a {
		names[v] = struct{}{}
		aBody = a[v]
	}
	for v := range b {
		names[v] = struct{}{}
		bBody = b[v]
	}
	if len(names) != 1 {
		bBody.Schema.SetXDiff(DIFF_NEW)
	} else {
		compareJsonSchema(aBody.Schema, bBody.Schema)
	}
	return b
}

func compareRequest(a, b *spec.CollectionHttpRequest) {
	compareHttpParameters(a.Attrs.Parameters, b.Attrs.Parameters)
	b.Attrs.Content = compareContent(a.Attrs.Content, b.Attrs.Content)
}

func compareResponse(a, b *spec.CollectionHttpResponse) *spec.CollectionHttpResponse {
	codes := map[int]struct{}{}
	for _, v := range a.Attrs.List {
		codes[v.Code] = struct{}{}
	}
	for _, v := range b.Attrs.List {
		codes[v.Code] = struct{}{}
	}
	// aa := a.Map()
	// bb := b.Map()
	for code := range codes {
		aRes := a.Attrs.List.FindByCode(code)
		aExist := aRes != nil

		bRes := b.Attrs.List.FindByCode(code)
		bExist := bRes != nil

		if !aExist && bExist {
			bRes.SetXDiff(DIFF_NEW)
			b.Attrs.List.AddOrUpdate(bRes)
			continue
		}
		if aExist && !bExist {
			aRes.SetXDiff(DIFF_REMOVE)
			bRes = aRes
			b.Attrs.List.AddOrUpdate(bRes)
			continue
		}
		if bRes.Name != aRes.Name || bRes.Description != aRes.Description {
			bRes.SetXDiff(DIFF_UPDATE)
		}
		bRes.Header = compareParameterList(aRes.Header, bRes.Header)
		bRes.Content = compareContent(aRes.Content, bRes.Content)
	}
	return b
}

// Same returns true otherwise false
func compareJsonSchema(a, b *jsonschema.Schema) bool {
	if a == nil && b == nil {
		return false
	}

	if a == nil && b != nil {
		b.SetXDiff(DIFF_NEW)
		return true
	}

	if a != nil && b == nil {
		a.SetXDiff(DIFF_REMOVE)
		b = a
		return true
	}

	if a.Reference != nil || b.Reference != nil {
		if a.Reference != nil && b.Reference != nil && a.Reference == b.Reference {
			return false
		}
		b.SetXDiff(DIFF_UPDATE)
		return true
	}

	if !slices.Equal(a.Type.List(), b.Type.List()) {
		b.SetXDiff(DIFF_UPDATE)
		return true
	}
	if len(a.Type.List()) == 0 {
		return false
	}

	aType := a.Type.First()
	bType := b.Type.First()
	// For array to object changes all are updated
	if aType != bType {
		b.SetXDiff(DIFF_UPDATE)
		return true
	}

	if !compareJsonSchemaBasicInfo(a, b) {
		b.SetXDiff(DIFF_UPDATE)
	}

	change := false
	switch bType {
	case "object":
		names := map[string]struct{}{}
		for v := range a.Properties {
			names[v] = struct{}{}
		}
		for v := range b.Properties {
			names[v] = struct{}{}
		}

		for v := range names {
			as, aExist := a.Properties[v]
			bs, bExist := b.Properties[v]
			if !aExist && bExist {
				bs.SetXDiff(DIFF_NEW)
				continue
			}
			if aExist && !bExist {
				as.SetXDiff(DIFF_REMOVE)
				if b.Properties == nil {
					b.Properties = make(map[string]*jsonschema.Schema)
				}
				b.Properties[v] = as
				// add to xorder
				if b.XOrder == nil {
					b.XOrder = make([]string, 0)
				}
				b.XOrder = append(b.XOrder, v)
				continue
			}
			sc := compareJsonSchema(as, bs)
			change = change || sc
		}
	case "array":
		if a.Items != nil && b.Items != nil {
			sc := compareJsonSchema(a.Items.Value(), b.Items.Value())
			change = change || sc
		}
	}
	return change
}

// Same returns true otherwise false
func compareJsonSchemaBasicInfo(a, b *jsonschema.Schema) bool {
	if fmt.Sprintf("%v", a.Default) != fmt.Sprintf("%v", b.Default) || a.Description != b.Description {
		return false
	}
	return true
}

// Same returns true otherwise false
func compareParamBasicInfo(a *spec.Parameter, b *spec.Parameter, onlyRequired bool) bool {
	if onlyRequired {
		if a.Required != b.Required {
			return false
		}
	} else {
		if a.Description != b.Description || a.Required != b.Required {
			return false
		}
	}
	return true
}
