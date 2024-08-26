package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

func SchemaToText(k string, s *jsonschema.Schema) string {
	if k == "" || s == nil {
		return ""
	}

	if result := convert(k, s, 0); len(result) > 0 {
		return strings.Join(result, "\n")
	}
	return ""
}

func SchemaToTextList(k string, s *jsonschema.Schema) []string {
	if k == "" || s == nil {
		return nil
	}

	return convert(k, s, 0)
}

func convert(k string, s *jsonschema.Schema, depth int) (result []string) {
	if k == "" || s == nil {
		return result
	}

	result = append(result, genText(k, s.Description, depth))

	if len(s.AllOf) > 0 {
		for _, v := range s.AllOf {
			if v.Properties != nil {
				for _, p := range v.XOrder {
					if _, ok := v.Properties[p]; ok {
						result = append(result, convert(p, v.Properties[p], depth+1)...)
					}
				}
			}
			if refID, err := v.GetRefID(); err == nil && refID > 0 {
				refk := "ref" + strconv.FormatInt(refID, 10)
				result = append(result, convert(refk, v, depth+1)...)
			}
		}
	} else if len(s.AnyOf) > 0 {
		for i, v := range s.AnyOf {
			result = append(result, convert(fmt.Sprintf("type%d", i+1), v, depth+1)...)
		}
	} else if len(s.OneOf) > 0 {
		for i, v := range s.OneOf {
			result = append(result, convert(fmt.Sprintf("type%d", i), v, depth+1)...)
		}
	} else if s.Properties != nil {
		for _, v := range s.XOrder {
			if _, ok := s.Properties[v]; ok {
				result = append(result, convert(v, s.Properties[v], depth+1)...)
			}
		}
	} else if s.Items != nil && !s.Items.IsBool() {
		result = append(result, convert("items", s.Items.Value(), depth+1)...)
	}

	return result
}

func genText(name, description string, depth int) string {
	return fmt.Sprintf("%s%s: %s", strings.Repeat("  ", depth), name, description)
}
