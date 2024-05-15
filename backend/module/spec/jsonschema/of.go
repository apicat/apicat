package jsonschema

type AllOf []*Schema
type AnyOf []*Schema
type OneOf []*Schema

func (a AllOf) Merge() error {
	return nil
}
