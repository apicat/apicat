package vector

import "errors"

const (
	WEAVIATE = "weaviate"
)

type Vector struct {
	Driver   string
	Weaviate WeaviateOpt
}

func NewVector(cfg Vector) (VectorApi, error) {
	switch cfg.Driver {
	case WEAVIATE:
		return NewWeaviate(cfg.Weaviate)
	default:
		return nil, errors.New("vector driver not found")
	}
}
