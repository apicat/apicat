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
					diffRequest(originalNode.ToHttpRequest(), targetNode.ToHttpRequest())
				case spec.NODE_HTTP_RESPONSE:
					diffResponse(originalNode.ToHttpResponse(), targetNode.ToHttpResponse())
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
	diffJsonSchema(copy_original.Schema, target.Schema)
	return nil
}

func diffHttpParameters(a *spec.HTTPParameters, b *spec.HTTPParameters) {
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

		newv := diffParameterList(aParamList, bParamList)
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

func diffParameterList(a, b spec.ParameterList) spec.ParameterList {
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
		diffParameter(ap, bp)
	}
	return b
}

// Same returns true otherwise false
func diffParameter(a, b *spec.Parameter) bool {
	same := true
	if !diffParamBasicInfo(a, b, false) {
		b.SetXDiff(DIFF_UPDATE)
		same = false
	}
	ss := diffJsonSchema(a.Schema, b.Schema)
	return same && ss
}

func diffContent(a, b spec.HTTPBody) spec.HTTPBody {
	names := map[string]struct{}{}
	aBody, bBody := &spec.Body{}, &spec.Body{}
	// httpbody just have one value
	for contentType := range a {
		names[contentType] = struct{}{}
		aBody = a[contentType]
	}
	for contentType := range b {
		names[contentType] = struct{}{}
		bBody = b[contentType]
	}
	if len(names) != 1 {
		bBody.Schema.SetXDiff(DIFF_NEW)
	} else {
		diffJsonSchema(aBody.Schema, bBody.Schema)
	}
	return b
}

func diffRequest(a, b *spec.CollectionHttpRequest) {
	diffHttpParameters(a.Attrs.Parameters, b.Attrs.Parameters)
	b.Attrs.Content = diffContent(a.Attrs.Content, b.Attrs.Content)
}

func diffResponse(a, b *spec.CollectionHttpResponse) *spec.CollectionHttpResponse {
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
		bRes.Header = diffParameterList(aRes.Header, bRes.Header)
		bRes.Content = diffContent(aRes.Content, bRes.Content)
	}
	return b
}

// Same returns true otherwise false
func diffJsonSchema(a, b *jsonschema.Schema) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil && b != nil {
		b.SetXDiff(DIFF_NEW)
		return false
	}

	if a != nil && b == nil {
		a.SetXDiff(DIFF_REMOVE)
		b = a
		return false
	}

	if a.Reference != nil || b.Reference != nil {
		if a.Reference != nil && b.Reference != nil && a.Reference == b.Reference {
			return true
		}
		b.SetXDiff(DIFF_UPDATE)
		return false
	}

	if !diffThreeOf(a, b) {
		return false
	}

	if !slices.Equal(a.Type.List(), b.Type.List()) {
		b.SetXDiff(DIFF_UPDATE)
		return false
	}
	if len(a.Type.List()) == 0 {
		return true
	}

	aType := a.Type.First()
	bType := b.Type.First()
	// For array to object changes all are updated
	if aType != bType {
		b.SetXDiff(DIFF_UPDATE)
		return false
	}

	same := true
	if !diffJsonSchemaBasicInfo(a, b) {
		b.SetXDiff(DIFF_UPDATE)
		same = false
	}

	switch bType {
	case "object":
		nameMap := map[string]struct{}{}
		nameList := make([]string, 0)

		if len(a.XOrder) > 0 {
			for _, v := range a.XOrder {
				if _, ok := nameMap[v]; !ok {
					nameMap[v] = struct{}{}
					nameList = append(nameList, v)
				}
			}
		} else {
			for v := range a.Properties {
				if _, ok := nameMap[v]; !ok {
					nameMap[v] = struct{}{}
					nameList = append(nameList, v)
				}
			}
		}

		if len(b.XOrder) > 0 {
			for _, v := range b.XOrder {
				if _, ok := nameMap[v]; !ok {
					nameMap[v] = struct{}{}
					nameList = append(nameList, v)
				}
			}
		} else {
			for v := range b.Properties {
				if _, ok := nameMap[v]; !ok {
					nameMap[v] = struct{}{}
					nameList = append(nameList, v)
				}
			}
		}

		for _, v := range nameList {
			as, aExist := a.Properties[v]
			bs, bExist := b.Properties[v]
			if !aExist && bExist {
				bs.SetXDiff(DIFF_NEW)
				same = false
				continue
			}
			if aExist && !bExist {
				as.SetXDiff(DIFF_REMOVE)
				same = false
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
			ss := diffJsonSchema(as, bs)
			same = same && ss
		}
	case "array":
		if a.Items != nil && b.Items != nil {
			ss := diffJsonSchema(a.Items.Value(), b.Items.Value())
			same = same && ss
		}
	}
	return same
}

// Same returns true otherwise false
func diffJsonSchemaBasicInfo(a, b *jsonschema.Schema) bool {
	if fmt.Sprintf("%v", a.Default) != fmt.Sprintf("%v", b.Default) || a.Description != b.Description {
		return false
	}
	return true
}

// Same returns true otherwise false
func diffParamBasicInfo(a *spec.Parameter, b *spec.Parameter, onlyRequired bool) bool {
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

func diffThreeOf(a, b *jsonschema.Schema) bool {
	if !diffOf(a, b, "allOf") {
		return false
	}
	if !diffOf(a, b, "anyOf") {
		return false
	}
	if !diffOf(a, b, "oneOf") {
		return false
	}
	return true
}

func diffOf(a, b *jsonschema.Schema, of string) bool {
	var aCount, bCount int
	switch of {
	case "allOf":
		aCount = len(a.AllOf)
		bCount = len(b.AllOf)
	case "anyOf":
		aCount = len(a.AnyOf)
		bCount = len(b.AnyOf)
	case "oneOf":
		aCount = len(a.OneOf)
		bCount = len(b.OneOf)
	default:
		return false
	}

	if aCount == 0 && bCount == 0 {
		return true
	}
	if aCount == 0 && bCount > 0 {
		b.SetXDiff(DIFF_NEW)
		return false
	}
	if aCount > 0 && bCount == 0 {
		b.SetXDiff(DIFF_REMOVE)
		switch of {
		case "allOf":
			b.AllOf = a.AllOf
		case "anyOf":
			b.AnyOf = a.AnyOf
		case "oneOf":
			b.OneOf = a.OneOf
		}
		return false
	}

	same := true
	i := 0
	for {
		if i == aCount && i == bCount {
			break
		}

		if i == aCount && i < bCount {
			for j := i; j < bCount; j++ {
				switch of {
				case "allOf":
					b.AllOf[j].SetXDiff(DIFF_NEW)
				case "anyOf":
					b.AnyOf[j].SetXDiff(DIFF_NEW)
				case "oneOf":
					b.OneOf[j].SetXDiff(DIFF_NEW)
				}
			}
			same = false
			break
		}

		if i < aCount && i == bCount {
			for j := i; j < aCount; j++ {
				switch of {
				case "allOf":
					a.AllOf[j].SetXDiff(DIFF_REMOVE)
					b.AllOf = append(b.AllOf, a.AllOf[j])
				case "anyOf":
					a.AnyOf[j].SetXDiff(DIFF_REMOVE)
					b.AnyOf = append(b.AnyOf, a.AnyOf[j])
				case "oneOf":
					a.OneOf[j].SetXDiff(DIFF_REMOVE)
					b.OneOf = append(b.OneOf, a.OneOf[j])
				}
			}
			same = false
			break
		}

		switch of {
		case "allOf":
			if !diffJsonSchema(a.AllOf[i], b.AllOf[i]) {
				same = false
			}
		case "anyOf":
			if !diffJsonSchema(a.AnyOf[i], b.AnyOf[i]) {
				same = false
			}
		case "oneOf":
			if !diffJsonSchema(a.OneOf[i], b.OneOf[i]) {
				same = false
			}
		}
		i++
	}
	return same
}
