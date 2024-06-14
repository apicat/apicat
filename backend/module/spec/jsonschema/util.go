package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type DerefHelper struct {
	RefMap map[int64]*Schema
}

func NewDerefHelper(refs map[int64]*Schema) *DerefHelper {
	return &DerefHelper{
		RefMap: refs,
	}
}

func (h *DerefHelper) DeepDeref(s *Schema) (Schema, error) {
	copySchema := &Schema{}
	if h == nil {
		return *copySchema, errors.New("schema id is 0")
	}

	tmp, err := json.Marshal(s)
	if err != nil {
		return *copySchema, errors.New("failed to marshal schema")
	}
	err = json.Unmarshal(tmp, copySchema)
	if err != nil {
		return *copySchema, errors.New("failed to unmarshal schema")
	}

	if copySchema.DeepRef() {
		if err := h.deref(copySchema, fmt.Sprintf("[%d]", copySchema.ID)); err != nil {
			return *copySchema, err
		}
	}
	return *copySchema, nil
}

func (h *DerefHelper) deref(s *Schema, path string) error {
	if h == nil {
		return errors.New("helper is nil")
	}

	if s.Ref() {
		refID, err := s.GetRefID()
		if err != nil {
			return err
		}

		ref, ok := h.RefMap[refID]
		if !ok {
			return fmt.Errorf("referenced schema id %d not found", refID)
		}

		// Check if the parent has a reference to this schema, avoid circular references
		if strings.Contains(path, fmt.Sprintf("[%d]", refID)) {
			*s = *NewSchema(T_OBJ)
			return nil
		}

		path = fmt.Sprintf("%s-[%d]", path, refID)
		// Dereference its reference schema
		if err := h.deref(ref, path); err != nil {
			return err
		}

		// Dereference itself
		if err := s.ReplaceRef(ref); err != nil {
			return err
		}
		return nil
	}

	if len(s.AllOf) > 0 {
		for _, v := range s.AllOf {
			if err := h.deref(v, path); err != nil {
				return err
			}
		}
	}
	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			if err := h.deref(v, path); err != nil {
				return err
			}
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			if err := h.deref(v, path); err != nil {
				return err
			}
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			if err := h.deref(v, path); err != nil {
				return err
			}
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		if err := h.deref(s.Items.Value(), path); err != nil {
			return err
		}
	}

	return nil
}
