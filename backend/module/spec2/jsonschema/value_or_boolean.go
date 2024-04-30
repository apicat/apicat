package jsonschema

import "encoding/json"

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
	return json.Unmarshal(raw, &v.value)
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
