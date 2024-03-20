package diff

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"

	"golang.org/x/exp/slices"
)

var (
	diffNew    = "+"
	diffRemove = "-"
	diffUpdate = "!"
	// diffType        = "type"
	// diffName        = "name"
	// diffDescription = "description"
	// diffRequired    = "required"
	// diffMock        = "mock"
	// diffDefault     = "default"
)

// Diff 比较两个接口的差异
// source,target 是完整的spec对象 因为需要解析schema等依赖
// spec.Collections 里面只能有一个接口
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

	for _, an := range copy_original.Content {
		for _, bn := range target.Content {
			if an.Node.NodeType() == bn.Node.NodeType() {
				// assertion in three parts to diff
				switch an.Node.NodeType() {
				case spec.NAME_HTTP_URL:
					au, err := an.ToHTTPURLNode()
					bu, err := bn.ToHTTPURLNode()
					if err != nil {
						return err
					}
					if au.Path != bu.Path {
						if au.Path == "" {
							bu.XDiff = &diffNew
						} else if bu.Path == "" {
							bu.Path = au.Path
							bu.XDiff = &diffRemove
						} else {
							bu.XDiff = &diffUpdate
						}
					}
				case spec.NAME_HTTP_REQUEST:
					ar, err := an.ToHTTPRequestNode()
					br, err := bn.ToHTTPRequestNode()
					if err != nil {
						return err
					}
					equalRequest(ar, br)
				case spec.NAME_HTTP_RESPONSES:
					ar, err := an.ToHTTPResponsesNode()
					br, err := bn.ToHTTPResponsesNode()
					if err != nil {
						return err
					}
					br.List = equalResponse(ar.List, br.List)
				default:
					return errors.New("node type error")
				}
			}
		}
	}
	return nil
}

func DiffSchema(original, target *spec.Schema) error {
	if original == nil || target == nil {
		return errors.New("schema is nil")
	}
	copy_original := &spec.Schema{}
	bb, err := json.Marshal(original)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bb, copy_original)
	if err != nil {
		return err
	}
	// if a.Name != b.Name {
	// 	b.XDiff = &diffUpdate
	// }
	equalJsonSchema(copy_original.Schema, target.Schema)
	return nil
}

func getMapOne(d map[string]map[string]spec.HTTPPart) (*spec.HTTPPart, *spec.HTTPURLNode) {
	for path, v := range d {
		for method, vv := range v {
			return &vv, &spec.HTTPURLNode{
				Method: method,
				Path:   path,
			}
		}
	}
	return nil, nil
}

func equalParam(a spec.HTTPParameters, b *spec.HTTPParameters) {
	a1 := a.Map()
	b1 := b.Map()
	for _, v := range spec.HttpParameterType {
		ap, a_has := a1[v]
		bp, b_has := b1[v]

		if !a_has && !b_has {
			continue
		}
		if a_has && !b_has {
			ap.SetXDiff(&diffRemove)
			bp = append(bp, ap...)
			continue
		}
		if !a_has && b_has {
			bp.SetXDiff(&diffNew)
			continue
		}

		newv := equalParameters(ap, bp)
		switch v {
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

func equalSchemas(a, b spec.Schemas) spec.Schemas {
	names := map[string]struct{}{}
	for _, v := range a {
		names[v.Name] = struct{}{}
	}
	for _, v := range b {
		names[v.Name] = struct{}{}
	}

	for v := range names {
		as := a.LookupByName(v)
		bs := b.LookupByName(v)

		if as == nil && bs == nil {
			continue
		}
		if as == nil && bs != nil {
			bs.SetXDiff(&diffNew)
			continue
		}
		if as != nil && bs == nil {
			as.SetXDiff(&diffRemove)
			b = append(b, as)
			continue
		}
		equalSchema(as, bs)
	}
	return b
}

func equalParameters(a, b spec.ParameterList) spec.ParameterList {
	names := map[string]struct{}{}
	for _, v := range a {
		names[v.Name] = struct{}{}
	}
	for _, v := range b {
		names[v.Name] = struct{}{}
	}

	for v := range names {
		ap := a.LookupByName(v)
		bp := b.LookupByName(v)

		if ap == nil && bp == nil {
			continue
		}
		if ap == nil && bp != nil {
			bp.SetXDiff(&diffNew)
			continue
		}
		if ap != nil && bp == nil {
			ap.SetXDiff(&diffRemove)
			b = append(b, ap)
			continue
		}
		equalParameter(ap, bp)
	}
	return b
}

func equalParameter(a, b *spec.Parameter) bool {
	change := false
	if !a.EqualNomal(b, false) {
		b.XDiff = &diffUpdate
		change = true
	}
	sc := equalJsonSchema(a.Schema, b.Schema)
	return change || sc
}

func equalContent(a, b spec.HTTPBody) spec.HTTPBody {
	names := map[string]struct{}{}
	as, bs := &spec.Schema{}, &spec.Schema{}
	// httpbody just have one value
	for v := range a {
		names[v] = struct{}{}
		as = a[v]
	}
	for v := range b {
		names[v] = struct{}{}
		bs = b[v]
	}
	if len(names) != 1 {
		bs.SetXDiff(&diffNew)
	} else {
		equalSchema(as, bs)
	}
	return b
}

func equalRequest(a, b *spec.HTTPRequestNode) {
	equalParam(a.Parameters, &b.Parameters)
	b.Content = equalContent(a.Content, b.Content)
}

func equalResponse(a, b spec.HTTPResponses) spec.HTTPResponses {
	ids := map[int]struct{}{}
	for _, v := range a {
		ids[v.Code] = struct{}{}
	}
	for _, v := range b {
		ids[v.Code] = struct{}{}
	}
	// aa := a.Map()
	// bb := b.Map()
	for k := range ids {
		as := a.LookupCode(k)
		a_has := as != nil
		bs := b.LookupCode(k)
		b_has := bs != nil
		if !a_has && b_has {
			bs.SetXDiff(&diffNew)
			goto e
		}
		if a_has && !b_has {
			as.SetXDiff(&diffRemove)
			bs = as
			goto e
		}
		if bs.Name != as.Name || bs.Description != as.Description {
			bs.XDiff = &diffUpdate
		}
		bs.Header = equalParameters(as.Header, bs.Header)
		bs.Content = equalContent(as.Content, bs.Content)
		// if bs is changed, goto e and add to result
	e:
		b.Merge(bs.Code, &bs.HTTPResponseDefine)
	}
	return b
}

func equalSchema(a, b *spec.Schema) bool {
	change := false
	if !a.EqualNomal(b) {
		change = true
	}
	sc := equalJsonSchema(a.Schema, b.Schema)
	return change || sc
}

func equalJsonSchema(a, b *jsonschema.Schema) bool {
	if a == nil && b == nil {
		return false
	}

	if a == nil && b != nil {
		b.SetXDiff(&diffNew)
		return true
	}

	if a != nil && b == nil {
		a.SetXDiff(&diffRemove)
		b = a
		return true
	}

	if a.Reference != nil || b.Reference != nil {
		if a.Reference != nil && b.Reference != nil && *a.Reference == *b.Reference {
			return false
		}
		b.SetXDiff(&diffUpdate)
		return true
	}

	if !slices.Equal(a.Type.Value(), b.Type.Value()) {
		b.SetXDiff(&diffUpdate)
		return true
	}
	if len(a.Type.Value()) == 0 {
		return false
	}
	at := a.Type.Value()[0]
	bt := b.Type.Value()[0]
	// For array to object changes all are updated
	if at != bt {
		b.SetXDiff(&diffUpdate)
		return true
	}

	if !equalJsonSchemaNormal(a, b) {
		b.XDiff = &diffUpdate
	}

	change := false
	switch bt {
	case "object":
		names := map[string]struct{}{}
		for v := range a.Properties {
			names[v] = struct{}{}
		}
		for v := range b.Properties {
			names[v] = struct{}{}
		}

		for v := range names {
			as, a_has := a.Properties[v]
			bs, b_has := b.Properties[v]
			if !a_has && b_has {
				bs.SetXDiff(&diffNew)
				continue
			}
			if a_has && !b_has {
				as.SetXDiff(&diffRemove)
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
			sc := equalJsonSchema(as, bs)
			change = change || sc
		}
	case "array":
		if a.Items != nil && b.Items != nil {
			sc := equalJsonSchema(a.Items.Value(), b.Items.Value())
			change = change || sc
		}
	}
	return change
}

func equalJsonSchemaNormal(a, b *jsonschema.Schema) bool {
	if fmt.Sprintf("%v", a.Default) != fmt.Sprintf("%v", b.Default) || a.Description != b.Description || fmt.Sprintf("%v", b.Example) != fmt.Sprintf("%v", b.Example) {
		return false
	}
	return true
}
