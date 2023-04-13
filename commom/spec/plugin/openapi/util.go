package openapi

import (
	"fmt"

	"github.com/apicat/apicat/commom/spec/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func jsonSchemaConverter(in *base.Schema) (*jsonschema.Schema, error) {
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
			js, err := jsonSchemaConverter(v.Schema())
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
			v, err := jsonSchemaConverter(addprop.Schema())
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
			v, err := jsonSchemaConverter(in.Items.A.Schema())
			if err != nil {
				return nil, err
			}
			items.SetValue(v)
		} else {
			items.SetBoolean(in.Items.B)
		}
		out.Items = items
	}

	if in.Deprecated != nil && *in.Deprecated == true {
		out.Deprecated = true
	}

	return &out, nil
}
