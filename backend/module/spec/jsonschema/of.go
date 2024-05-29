package jsonschema

import (
	"strconv"
)

type Of []*Schema

func (al Of) Merge() Of {
	helper := NewMergeHelper(&Schema{})
	helper.Merge(al)
	return helper.list
}

func (al *Of) DelRef(ref *Schema) {
	refid := strconv.FormatInt(ref.ID, 10)
	i := 0
	for i < len(*al) {
		if (*al)[i].IsRefID(refid) {
			*al = append((*al)[:i], (*al)[i+1:]...)
			return
		} else {
			(*al)[i].DelRef(ref)
			i++
		}
	}
}

type mergeHelper struct {
	result *Schema
	list   Of
}

func NewMergeHelper(s *Schema) *mergeHelper {
	m := &mergeHelper{
		result: s,
		list:   make(Of, 0),
	}
	m.list = append(m.list, m.result)
	return m
}

func (m *mergeHelper) Merge(froms Of) Of {
	for _, from := range froms {
		if from.CheckAllOf() {
			m.Merge(from.AllOf)
			continue
		}

		if from.Ref() {
			m.list = append(m.list, from)
			continue
		}

		if from.AnyOf != nil {
			for _, v := range from.AnyOf {
				v.MergeAllOf()
			}
		}

		if from.OneOf != nil {
			for _, v := range from.OneOf {
				v.MergeAllOf()
			}
		}

		if from.Properties != nil {
			for _, prop := range from.Properties {
				prop.MergeAllOf()
			}
		}
		if from.Items != nil && !from.Items.IsBool() {
			from.Items.Value().MergeAllOf()
		}

		m.merge(from)
	}
	return m.list
}

func (m *mergeHelper) merge(from *Schema) {
	m.mergeType(from)
	m.mergeProperties(from)
	m.mergeRequired(from)
	m.mergeOthers(from)
	m.mergeXOrder(from)
}

func (m *mergeHelper) mergeType(from *Schema) {
	if m.result.Type == nil && from.Type != nil {
		m.result.Type = from.Type
	}
}

func (m *mergeHelper) mergeProperties(from *Schema) {
	if from.Properties == nil {
		return
	}

	if m.result.Properties == nil {
		m.result.Properties = make(map[string]*Schema)
	}

	for name, prop := range from.Properties {
		m.result.Properties[name] = prop
	}
}

func (m *mergeHelper) mergeRequired(from *Schema) {
	if len(from.Required) > 0 {
		temp := make(map[string]bool)
		for _, v := range m.result.Required {
			temp[v] = true
		}

		for _, v := range from.Required {
			if _, ok := temp[v]; !ok {
				m.result.Required = append(m.result.Required, v)
				temp[v] = true
			}
		}
	}
}

func (m *mergeHelper) mergeXOrder(from *Schema) {
	if len(from.XOrder) > 0 {
		temp := make(map[string]bool)
		for _, v := range m.result.XOrder {
			temp[v] = true
		}

		for _, v := range from.XOrder {
			if _, ok := temp[v]; !ok {
				m.result.XOrder = append(m.result.XOrder, v)
				temp[v] = true
			}
		}
	}
}

func (m *mergeHelper) mergeOthers(from *Schema) {
	if len(m.result.Enum) == 0 && len(from.Enum) > 0 {
		m.result.Enum = from.Enum
	}
	if m.result.Pattern == "" && from.Pattern != "" {
		m.result.Pattern = from.Pattern
	}
	if m.result.MinLength == nil && from.MinLength != nil {
		m.result.MinLength = from.MinLength
	}
	if m.result.MaxLength == nil && from.MaxLength != nil {
		m.result.MaxLength = from.MaxLength
	}
	if m.result.ExclusiveMaximum == nil && from.ExclusiveMaximum != nil {
		m.result.ExclusiveMaximum = from.ExclusiveMaximum
	}
	if m.result.MultipleOf == nil && from.MultipleOf != nil {
		m.result.MultipleOf = from.MultipleOf
	}
	if m.result.ExclusiveMinimum == nil && from.ExclusiveMinimum != nil {
		m.result.ExclusiveMinimum = from.ExclusiveMinimum
	}
	if m.result.Maximum == nil && from.Maximum != nil {
		m.result.Maximum = from.Maximum
	}
	if m.result.Minimum == nil && from.Minimum != nil {
		m.result.Minimum = from.Minimum
	}
	if m.result.MaxProperties == nil && from.MaxProperties != nil {
		m.result.MaxProperties = from.MaxProperties
	}
	if m.result.MinProperties == nil && from.MinProperties != nil {
		m.result.MinProperties = from.MinProperties
	}
	if m.result.MaxItems == nil && from.MaxItems != nil {
		m.result.MaxItems = from.MaxItems
	}
	if m.result.MinItems == nil && from.MinItems != nil {
		m.result.MinItems = from.MinItems
	}
	if m.result.UniqueItems == nil && from.UniqueItems != nil {
		m.result.UniqueItems = from.UniqueItems
	}
	if m.result.Format == "" && from.Format != "" {
		m.result.Format = from.Format
	}
	if m.result.XMock == "" && from.XMock != "" {
		m.result.XMock = from.XMock
	}
}
