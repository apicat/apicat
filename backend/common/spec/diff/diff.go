package diff

import (
	"encoding/json"
	"errors"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
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
func Diff(ac, bc []byte) (*spec.CollectItem, error) {

	source, err := spec.ParseJSON(ac)
	if err != nil {
		return nil, errors.New("source parse error")
	}
	target, err := spec.ParseJSON(bc)
	if err != nil {
		return nil, errors.New("target parse error")
	}

	if len(source.Collections) != 1 || len(target.Collections) != 1 {
		return nil, errors.New("source,target Collections length error")
	}
	a, au := getMapOne(source.CollectionsMap(true, 1))
	b, bu := getMapOne(target.CollectionsMap(true, 1))
	if au.Path != bu.Path {
		bu.XDiff = &diffUpdate
	}
	if a.Title != b.Title {
		b.XDiff = &diffUpdate
	}
	equalRequest(&a.HTTPRequestNode, &b.HTTPRequestNode)
	b.Responses = equalResponse(a.Responses, b.Responses)
	return b.ToCollectItem(*bu), nil
}

func DiffSchema(as, bs []byte) (*jsonschema.Schema, error) {

	a := &jsonschema.Schema{}
	err := json.Unmarshal(as, a)
	if err != nil {
		return nil, errors.New("source parse error")
	}
	b := &jsonschema.Schema{}
	err = json.Unmarshal(bs, b)
	if err != nil {
		return nil, errors.New("target parse error")
	}

	equalJsonSchema(a, b)
	return b, nil
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
	for _, v := range spec.HttpParameter {
		ap, a_has := a1[v]
		bp, b_has := b1[v]

		if !a_has && !b_has {
			continue
		}
		// TODO可能会考虑到排序的问题
		if a_has && !b_has {
			ap.SetXDiff(&diffRemove)
			bp = append(bp, ap...)
			continue
		}
		if !a_has && b_has {
			bp.SetXDiff(&diffNew)
			continue
		}

		newv := equalSchemas(ap, bp)
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
		as := a.Lookup(v)
		bs := b.Lookup(v)

		// TODO可能会考虑到排序的问题
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

func equalContent(a, b spec.HTTPBody) spec.HTTPBody {
	names := map[string]struct{}{}
	for v := range a {
		names[v] = struct{}{}
	}
	for v := range b {
		names[v] = struct{}{}
	}
	for v := range names {
		as, a_has := a[v]
		bs, b_has := b[v]
		if !a_has && b_has {
			bs.SetXDiff(&diffNew)
			continue
		}
		if a_has && !b_has {
			as.SetXDiff(&diffRemove)
			b[v] = as
			continue
		}
		equalSchema(as, bs)
	}
	return b
}

func equalRequest(a, b *spec.HTTPRequestNode) {
	equalParam(a.Parameters, &b.Parameters)
	b.Content = equalContent(a.Content, b.Content)
}

func equalResponse(a, b spec.HTTPResponses) spec.HTTPResponses {
	codes := map[int]struct{}{}
	for _, v := range a {
		codes[v.Code] = struct{}{}
	}
	for _, v := range b {
		codes[v.Code] = struct{}{}
	}
	aa := a.Map()
	bb := b.Map()
	for k := range codes {
		as, a_has := aa[k]
		bs, b_has := bb[k]
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
		bs.Header = equalSchemas(as.Header, bs.Header)
		bs.Content = equalContent(as.Content, bs.Content)
		// as,bs与b类型不同，需要特殊处理
	e:
		b.Add(k, &bs)
	}
	return b
}

func equalSchema(a, b *spec.Schema) {
	if !a.EqualNomal(b) {
		b.XDiff = &diffUpdate
	}
	equalJsonSchema(a.Schema, b.Schema)
}

func equalJsonSchema(a, b *jsonschema.Schema) {
	if !slices.Equal(a.Type.Value(), b.Type.Value()) {
		b.SetXDiff(&diffUpdate)
		return
	}
	at := a.Type.Value()[0]
	bt := b.Type.Value()[0]
	// For array to object changes all are updated
	if at != bt {
		b.SetXDiff(&diffUpdate)
		return
	}
	if !equalJsonSchemaNormal(a, b) {
		b.XDiff = &diffUpdate
	}
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
				b.Properties[v] = as
				continue
			}
			equalJsonSchema(as, bs)
		}
	case "array":
		if a.Items != nil && b.Items != nil {
			equalJsonSchema(a.Items.Value(), b.Items.Value())
		}
	}

}

func equalJsonSchemaNormal(a, b *jsonschema.Schema) bool {
	if a.Default != b.Default || a.Description != b.Description || a.Example != b.Example {
		return false
	}
	return true
}
