package config

import (
	"errors"
	"os"

	"github.com/apicat/apicat/v2/backend/module/vector"
)

type Vector struct {
	Driver   string
	Weaviate *Weaviate
}

type Weaviate struct {
	Host   string
	ApiKey string
}

func LoadVertorConfig() {
	globalConf.Vector = &Vector{}
	if v, exists := os.LookupEnv("VECTOR_DRIVER"); exists {
		switch v {
		case vector.WEAVIATE:
			globalConf.Vector.Driver = vector.WEAVIATE
			loadWeaviateConfig()
		}
	}
}

func loadWeaviateConfig() {
	globalConf.Vector.Weaviate = &Weaviate{}
	if v, exists := os.LookupEnv("WEAVIATE_HOST"); exists {
		globalConf.Vector.Weaviate.Host = v
	}
	if v, exists := os.LookupEnv("WEAVIATE_API_KEY"); exists {
		globalConf.Vector.Weaviate.ApiKey = v
	}
}

func CheckVectorConfig() error {
	if globalConf.Vector.Driver == "" {
		return errors.New("vector driver is empty")
	}
	switch globalConf.Vector.Driver {
	case vector.WEAVIATE:
		return checkWeaviateConfig()
	}
	return nil
}

func checkWeaviateConfig() error {
	if globalConf.Vector.Weaviate == nil {
		return errors.New("weaviate config is empty")
	}
	if globalConf.Vector.Weaviate.Host == "" {
		return errors.New("weaviate host is empty")
	}

	w, err := vector.NewVector(globalConf.Vector.ToModuleStruct())
	if err != nil {
		return err
	}
	if err := w.Check(); err != nil {
		return err
	}
	return nil
}

func (v *Vector) ToModuleStruct() vector.Vector {
	return vector.Vector{
		Driver: v.Driver,
		Weaviate: vector.WeaviateOpt{
			Host:   v.Weaviate.Host,
			ApiKey: v.Weaviate.ApiKey,
		},
	}
}

func GetVector() *Vector {
	return globalConf.Vector
}
