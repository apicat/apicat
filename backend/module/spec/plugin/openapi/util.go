package openapi

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

func jsonschemaIsRef(b *base.SchemaProxy) string {
	// check openapi 3.x
	if g := b.GoLow(); g != nil {
		if g.IsReference() {
			ref := g.GetReference()
			if strings.HasPrefix(ref, "#/definitions/") || strings.HasPrefix(ref, "#/components/schemas/") {
				// if len(mapping) > 0 {
				// 	id, ok := mapping[0][getRefName(ref)]
				// 	if ok {
				// 		refid := fmt.Sprintf("#/definitions/schemas/%d", id)
				// 		return &jsonschema.Schema{Reference: &refid}, nil
				// 	}
				// }
				refname := getRefName(ref)
				return refname
			}
		}
	}
	// check swagger 2.x
	if b.GetReference() != "" {
		refname := getRefName(b.GetReference())
		return refname
	}
	return ""
}

func jsonSchemaConverter(b *base.SchemaProxy) (*jsonschema.Schema, error) {
	if refname := jsonschemaIsRef(b); refname != "" {
		refid := fmt.Sprintf("#/definitions/schemas/%d", stringToUnid(refname))
		return &jsonschema.Schema{Reference: &refid}, nil
	}

	in := b.Schema()
	out := jsonschema.Schema{
		Title:         in.Title,
		Description:   in.Description,
		MultipleOf:    in.MultipleOf,
		Maximum:       in.Maximum,
		Minimum:       in.Minimum,
		MaxLength:     in.MaxLength,
		MinLength:     in.MinLength,
		Format:        in.Format,
		Pattern:       in.Pattern,
		MaxItems:      in.MaxItems,
		MinItems:      in.MinItems,
		UniqueItems:   in.UniqueItems,
		MaxProperties: in.MaxProperties,
		MinProperties: in.MinProperties,
		Nullable:      in.Nullable,
		ReadOnly:      in.ReadOnly,
		WriteOnly:     in.WriteOnly,
	}

	if len(in.Type) > 0 {
		out.Type = jsonschema.NewSchemaType(in.Type...)
	}

	if in.Default != nil && in.Default.Value != "" {
		out.Default = in.Default.Value
	}

	if in.Example != nil && in.Example.Value != "" {
		out.Examples = in.Example.Value
	}

	if in.ExclusiveMaximum != nil {
		em := &jsonschema.ValueOrBoolean[float64]{}
		if in.ExclusiveMaximum.IsA() {
			em.SetBoolean(in.ExclusiveMaximum.A)
		} else {
			em.SetValue(in.ExclusiveMaximum.B)
		}
		out.ExclusiveMaximum = em
	}

	if in.ExclusiveMinimum != nil {
		em := &jsonschema.ValueOrBoolean[float64]{}
		if in.ExclusiveMinimum.IsA() {
			em.SetBoolean(in.ExclusiveMinimum.A)
		} else {
			em.SetValue(in.ExclusiveMinimum.B)
		}
		out.ExclusiveMinimum = em
	}

	if in.Properties != nil {
		props := make(map[string]*jsonschema.Schema)
		names := make([]string, 0)
		for pair := range orderedmap.Iterate(context.Background(), in.Properties) {
			name := pair.Key()
			v := pair.Value()

			js, err := jsonSchemaConverter(v)
			if err != nil {
				return nil, err
			}
			props[name] = js
			names = append(names, name)
		}
		out.Properties = props
		out.XOrder = names
		out.Required = in.Required
	}

	if in.AdditionalProperties != nil {
		ap := &jsonschema.ValueOrBoolean[*jsonschema.Schema]{}
		if in.AdditionalProperties.IsA() {
			v, err := jsonSchemaConverter(in.AdditionalProperties.A)
			if err != nil {
				return nil, err
			}
			ap.SetValue(v)
		} else if in.AdditionalProperties.IsB() {
			ap.SetBoolean(in.AdditionalProperties.B)
		} else {
			return nil, fmt.Errorf("unsupport")
		}
		out.AdditionalProperties = ap
	}

	if in.Items != nil {
		items := &jsonschema.ValueOrBoolean[*jsonschema.Schema]{}
		if in.Items.IsA() {
			v, err := jsonSchemaConverter(in.Items.A)
			if err != nil {
				return nil, err
			}
			items.SetValue(v)
		} else {
			items.SetBoolean(in.Items.B)
		}
		out.Items = items
	}

	if in.Deprecated != nil && *in.Deprecated {
		out.Deprecated = in.Deprecated
	}

	if in.AllOf != nil {
		out.AllOf = make([]*jsonschema.Schema, 0)
		for _, v := range in.AllOf {
			js, err := jsonSchemaConverter(v)
			if err != nil {
				return nil, err
			}
			out.AllOf = append(out.AllOf, js)
		}

		out.MergeAllOf()
		if out.Type.First() != jsonschema.T_NULL && out.Type.First() != jsonschema.T_OBJ {
			// if the type is not object, we should take out the contents of allof, make the structure simple
			helper := jsonschema.NewMergeHelper(&out)
			if new := helper.Merge(out.AllOf); len(new) == 1 {
				out = *new[0]
			}
		}
	}

	if in.AnyOf != nil {
		out.AnyOf = make([]*jsonschema.Schema, 0)
		for _, v := range in.AnyOf {
			js, err := jsonSchemaConverter(v)
			if err != nil {
				return nil, err
			}
			out.AnyOf = append(out.AnyOf, js)
		}
	}

	if in.OneOf != nil {
		out.OneOf = make([]*jsonschema.Schema, 0)
		for _, v := range in.OneOf {
			js, err := jsonSchemaConverter(v)
			if err != nil {
				return nil, err
			}
			out.OneOf = append(out.OneOf, js)
		}
	}

	return &out, nil
}

func getRefName(ref string) string {
	return ref[strings.LastIndex(ref, "/")+1:]
}

func toInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func stringToUnid(s string) int64 {
	n := len(s)
	x := int64(n * 10000)
	for i := 0; i < n; i++ {
		x += int64(s[i])
	}
	return x
}

func isGlobalParameter(ref string) bool {
	return strings.Contains(ref, "/x-apicat-global-parameters/")
}

func globalToLocalParameters(globalsParmaters *spec.GlobalParameters, isSwagger bool, skip map[string][]int64) []openAPIParamter {
	var outs []openAPIParamter
	skips := make(map[string]bool)
	for k, v := range skip {
		for _, x := range v {
			skips[fmt.Sprintf("%s|_%d", k, x)] = true
		}
	}

	for in, paramList := range globalsParmaters.ToMap() {
		for _, p := range paramList {
			if skips[fmt.Sprintf("%s|_%d", in, p.ID)] {
				continue
			}

			ref := fmt.Sprintf("%s-%s", in, p.Name)
			if isSwagger {
				ref = "#/x-apicat-global-parameters/" + ref
			} else {
				ref = "#/components/x-apicat-global-parameters/" + ref
			}
			outs = append(outs, openAPIParamter{
				Reference: ref,
			})
		}
	}
	return outs
}

func toParameter(p *spec.Parameter, in string, version string) openAPIParamter {
	if version[0] == '3' {
		return toParameter3(p, in)
	}
	return toParameter2(p, in)
}

func toParameter3(p *spec.Parameter, in string) openAPIParamter {
	return openAPIParamter{
		In:          in,
		Name:        p.Name,
		Required:    p.Required,
		Format:      p.Schema.Format,
		Example:     p.Schema.Examples,
		Description: p.Schema.Description,
		Schema:      p.Schema,
	}
}

func toParameter2(p *spec.Parameter, in string) openAPIParamter {
	typ := jsonschema.T_NULL
	if num := len(p.Schema.Type.List()); num > 0 {
		typ = p.Schema.Type.First()
	}

	if in == "cookie" {
		in = "header"
	}
	return openAPIParamter{
		In:          in,
		Type:        typ,
		Name:        p.Name,
		Required:    p.Required,
		Format:      p.Schema.Format,
		Default:     p.Schema.Default,
		Description: p.Schema.Description,
		Schema:      p.Schema,
	}
}

func convertJsonSchemaRef(v *jsonschema.Schema, version string, mapping map[int64]string) *jsonschema.Schema {
	sh := *v
	if s, ok := sh.Examples.(string); ok && s == "" {
		sh.Examples = nil
	}

	if sh.Reference != nil {
		if id := toInt64(getRefName(*sh.Reference)); id > 0 {
			var ref string
			name_id := fmt.Sprintf("%s-%d", mapping[id], id)
			if version[0] == '2' {
				ref = fmt.Sprintf("#/definitions/%s", name_id)
			} else {
				ref = fmt.Sprintf("#/components/schemas/%s", name_id)
			}
			return &jsonschema.Schema{Reference: &ref}
		}
	}

	if len(sh.AllOf) > 0 {
		for i, v := range sh.AllOf {
			sh.AllOf[i] = convertJsonSchemaRef(v, version, mapping)
		}
	}
	if len(sh.AnyOf) > 0 {
		for i, v := range sh.AnyOf {
			sh.AnyOf[i] = convertJsonSchemaRef(v, version, mapping)
		}
	}
	if len(sh.OneOf) > 0 {
		for i, v := range sh.OneOf {
			sh.OneOf[i] = convertJsonSchemaRef(v, version, mapping)
		}
	}
	if sh.Properties != nil {
		for k, v := range sh.Properties {
			sh.Properties[k] = convertJsonSchemaRef(v, version, mapping)
		}
	}
	if sh.Items != nil {
		if !sh.Items.IsBool() {
			sh.Items.SetValue(convertJsonSchemaRef(sh.Items.Value(), version, mapping))
		}
	}
	if sh.AdditionalProperties != nil {
		if !sh.AdditionalProperties.IsBool() {
			sh.AdditionalProperties.SetValue(convertJsonSchemaRef(sh.AdditionalProperties.Value(), version, mapping))
		}
	}
	return &sh
}

func deepGetHttpCollection(in *spec.Collections) map[string]map[string]specPathItem {
	paths := make(map[string]map[string]specPathItem)

	for _, collection := range *in {
		if collection.Type == spec.TYPE_CATEGORY && len(collection.Items) > 0 {
			childrenPaths := deepGetHttpCollection(&collection.Items)
			for path, methods := range childrenPaths {
				for method, item := range methods {
					if _, ok := paths[path]; !ok {
						paths[path] = map[string]specPathItem{
							method: item,
						}
					} else {
						paths[path][method] = item
					}
				}
			}
		}

		if collection.Type != spec.TYPE_HTTP {
			continue
		}

		item := specPathItem{
			Title:      collection.Title,
			OperatorID: fmt.Sprintf("%d", collection.ID),
			Tags:       collection.Tags,
		}

		var info spec.CollectionHttpUrl
		for _, node := range collection.Content {
			switch node.NodeType() {
			case spec.NODE_HTTP_URL:
				info = *node.ToHttpUrl()
			case spec.NODE_HTTP_REQUEST:
				item.Req = *node.ToHttpRequest()
			case spec.NODE_HTTP_RESPONSE:
				item.Res = *node.ToHttpResponse()
			}
		}

		if _, ok := paths[info.Attrs.Path]; !ok {
			paths[info.Attrs.Path] = map[string]specPathItem{
				info.Attrs.Method: item,
			}
		} else {
			paths[info.Attrs.Path][info.Attrs.Method] = item
		}
	}
	return paths
}
