package openapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/backend/module/spec/jsonschema"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func jsonschemaIsRef(b *base.SchemaProxy) *string {
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
				return &refname
			}
		}
	}
	// check swagger 2.x
	if b.GetReference() != "" {
		refname := getRefName(b.GetReference())
		return &refname
	}
	return nil
}

func jsonSchemaConverter(b *base.SchemaProxy) (*jsonschema.Schema, error) {
	if refname := jsonschemaIsRef(b); refname != nil {
		refid := fmt.Sprintf("#/definitions/schemas/%d", stringToUnid(*refname))
		return &jsonschema.Schema{Reference: &refid}, nil
	}
	in := b.Schema()
	var t jsonschema.SliceOrOneValue[string]
	t.SetValue(in.Type...)
	out := jsonschema.Schema{
		Type:          &t,
		Title:         in.Title,
		Description:   in.Description,
		MultipleOf:    in.MultipleOf,
		Maximum:       in.Maximum,
		Minimum:       in.MinItems,
		MaxLength:     in.MaxLength,
		MinLength:     in.MinLength,
		Format:        in.Format,
		Pattern:       in.Pattern,
		MaxItems:      in.MaxItems,
		MinItems:      in.MinItems,
		UniqueItems:   in.UniqueItems,
		MaxProperties: in.MaxProperties,
		MinProperties: in.MinProperties,
		Default:       in.Default,
		Nullable:      in.Nullable,
		ReadOnly:      in.ReadOnly,
		WriteOnly:     in.WriteOnly,
		Example:       in.Example,
	}

	if in.ExclusiveMaximum != nil {
		em := &jsonschema.ValueOrBoolean[int64]{}
		if in.ExclusiveMaximum.IsA() {
			em.SetBoolean(in.ExclusiveMaximum.A)
		} else {
			em.SetValue(in.ExclusiveMaximum.B)
		}
		out.ExclusiveMaximum = em
	}

	if in.ExclusiveMinimum != nil {
		em := &jsonschema.ValueOrBoolean[int64]{}
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
		for name, v := range in.Properties {
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
		switch addprop := in.AdditionalProperties.(type) {
		case *base.SchemaProxy:
			v, err := jsonSchemaConverter(addprop)
			if err != nil {
				return nil, err
			}
			ap.SetValue(v)
		case bool:
			ap.SetBoolean(addprop)
		default:
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
		out.Deprecated = true
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
