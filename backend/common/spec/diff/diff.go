package diff

import (
	"errors"
	"fmt"

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
func Diff(ref_obj, diff_obj *spec.CollectItem) (*spec.CollectItem, error) {

	if ref_obj == nil || diff_obj == nil {
		return nil, errors.New("source,target Collections length error")
	}

	for _, an := range ref_obj.Content {
		for _, bn := range diff_obj.Content {
			if an.Node.NodeType() == bn.Node.NodeType() {
				// assertion in three parts to diff
				switch an.Node.NodeType() {
				case "apicat-http-url":
					au, err := an.ToHTTPURLNode()
					bu, err := bn.ToHTTPURLNode()
					if err != nil {
						return nil, err
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
				case "apicat-http-request":
					ar, err := an.ToHTTPRequestNode()
					br, err := bn.ToHTTPRequestNode()
					if err != nil {
						return nil, err
					}
					equalRequest(ar, br)
				case "apicat-http-response":
					ar, err := an.ToHTTPResponsesNode()
					br, err := bn.ToHTTPResponsesNode()
					if err != nil {
						return nil, err
					}
					br.List = equalResponse(ar.List, br.List)
				default:
					return nil, errors.New("node type error")
				}
			}
		}
	}
	return diff_obj, nil
}

func DiffSchema(a, b *jsonschema.Schema) (*jsonschema.Schema, error) {

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
			if b == nil {
				b = make(spec.HTTPBody)
			}
			b[v] = as
			continue
		}
		if equalSchema(as, bs) {
			bs.XDiff = &diffUpdate
		}
	}
	return b
}

func equalRequest(a, b *spec.HTTPRequestNode) {
	equalParam(a.Parameters, &b.Parameters)
	b.Content = equalContent(a.Content, b.Content)
}

func equalResponse(a, b spec.HTTPResponses) spec.HTTPResponses {
	ids := map[string]struct{}{}
	for _, v := range a {
		ids[fmt.Sprintf("%d-%s", v.Code, v.Name)] = struct{}{}
	}
	for _, v := range b {
		ids[fmt.Sprintf("%d-%s", v.Code, v.Name)] = struct{}{}
	}
	// aa := a.Map()
	// bb := b.Map()
	for k := range ids {
		as := a.LookupID(k)
		a_has := as != nil
		bs := b.LookupID(k)
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
		bs.Header = equalSchemas(as.Header, bs.Header)
		bs.Content = equalContent(as.Content, bs.Content)
		// if bs is changed, goto e and add to result
	e:
		b.Add(bs.Code, k, &bs.HTTPResponseDefine)
	}
	return b
}

func equalSchema(a, b *spec.Schema) bool {
	change := false
	if !a.EqualNomal(b) {
		b.XDiff = &diffUpdate
		change = true
	}
	sc := equalJsonSchema(a.Schema, b.Schema)
	return change || sc
}

func equalJsonSchema(a, b *jsonschema.Schema) bool {
	if !slices.Equal(a.Type.Value(), b.Type.Value()) {
		b.SetXDiff(&diffUpdate)
		return true
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
	if a.Default != b.Default || a.Description != b.Description || a.Example != b.Example {
		return false
	}
	return true
}
