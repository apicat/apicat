package rag

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/vector"
)

type apiContentProperty struct {
	CollectionID      vector.T_INT  `json:"collection_id,omitempty"`
	DefinitionModelID vector.T_INT  `json:"definition_model_id,omitempty"`
	UpdatedAt         vector.T_TEXT `json:"updated_at,omitempty"`
}

func (acp *apiContentProperty) ToMapInterface() (result map[string]interface{}) {
	b, _ := json.Marshal(acp)
	json.Unmarshal(b, &result)
	return
}

func getAPIContentProperties() vector.Properties {
	p := apiContentProperty{
		CollectionID:      0,
		DefinitionModelID: 0,
		UpdatedAt:         "",
	}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	properties := make(vector.Properties, 0)
	for i := 0; i < t.NumField(); i++ {
		properties = append(properties, &vector.Property{
			Name:     strings.Split(t.Field(i).Tag.Get("json"), ",")[0],
			DataType: v.Field(i).Interface().(vector.DataType),
		})
	}
	return properties
}
