package jsonschema

import (
	"encoding/json"
)

type ValueOrBoolean[T any] struct {
	boolValue *bool
	value     T
}

func (v *ValueOrBoolean[T]) UnmarshalJSON(raw []byte) error {
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		v.boolValue = &b
		return nil
	}
	return json.Unmarshal(raw, v.value)
}

func (v ValueOrBoolean[T]) MarshalJSON() ([]byte, error) {
	var o any
	if v.boolValue != nil {
		o = *v.boolValue
	} else {
		o = v.value
	}
	return json.Marshal(o)
}

func (v *ValueOrBoolean[T]) SetValue(value T) {
	v.boolValue = nil
	v.value = value
}

func (v *ValueOrBoolean[T]) SetBoolean(b bool) {
	v.boolValue = &b
}

func (v *ValueOrBoolean[T]) Value() T {
	return v.value
}

func (v *ValueOrBoolean[T]) Bool() bool {
	if v.boolValue != nil {
		return *v.boolValue
	}
	return false
}

func (v *ValueOrBoolean[T]) IsBool() bool {
	return v.boolValue != nil
}

func CreateSliceOrOne[T any](v ...T) *SliceOrOneValue[T] {
	return &SliceOrOneValue[T]{
		value:   v,
		isSlice: len(v) > 1,
	}
}

type SliceOrOneValue[T any] struct {
	value   []T
	isSlice bool
}

func (s *SliceOrOneValue[T]) Value() []T {
	return s.value
}

func (s *SliceOrOneValue[T]) SetValue(v ...T) {
	if len(v) > 1 {
		s.isSlice = true
	}
	s.value = v
}

func (s *SliceOrOneValue[T]) UnmarshalJSON(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	if p := raw[0]; p == '[' {
		s.isSlice = true
		return json.Unmarshal(raw, &s.value)
	}
	var o T
	if err := json.Unmarshal(raw, &o); err != nil {
		return err
	}
	s.value = []T{o}
	return nil
}

func (s SliceOrOneValue[T]) MarshalJSON() ([]byte, error) {
	if len(s.value) == 0 {
		return []byte("null"), nil
	}
	var v any
	if s.isSlice {
		v = s.value
	} else {
		v = s.value[0]
	}
	return json.Marshal(v)
}
