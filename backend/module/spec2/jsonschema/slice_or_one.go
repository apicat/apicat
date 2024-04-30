package jsonschema

import "encoding/json"

type SliceOrOneValue[T any] struct {
	value   []T
	isSlice bool
}

func NewSliceOrOne[T any](v ...T) *SliceOrOneValue[T] {
	return &SliceOrOneValue[T]{
		value:   v,
		isSlice: len(v) > 1,
	}
}

func (s *SliceOrOneValue[T]) Value() []T {
	if s == nil {
		return []T{}
	}
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
