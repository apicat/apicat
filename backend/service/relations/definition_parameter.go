package relations

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/module/spec"
)

// ImportDefinitionParameters 导入公共参数
func ImportDefinitionParameters(ctx context.Context, projectID string, parameters *spec.HTTPParameters) collection.VirtualIDToIDMap {
	res := collection.VirtualIDToIDMap{}

	if parameters.Header == nil && parameters.Cookie == nil && parameters.Query == nil && parameters.Path == nil {
		return res
	}

	var params spec.ParameterList
	parameterList := []string{definition.ParameterInHeader, definition.ParameterInCookie, definition.ParameterInQuery, definition.ParameterInPath}
	for _, key := range parameterList {
		switch key {
		case definition.ParameterInHeader:
			params = parameters.Header
		case definition.ParameterInCookie:
			params = parameters.Cookie
		case definition.ParameterInQuery:
			params = parameters.Query
		case definition.ParameterInPath:
			params = parameters.Path
		}

		for _, parameter := range params {
			if parameterStr, err := json.Marshal(parameter.Schema); err == nil {
				record := &definition.DefinitionParameter{
					ProjectID: projectID,
					In:        key,
					Name:      parameter.Name,
					Required:  parameter.Required,
					Schema:    string(parameterStr),
				}

				if model.DB(ctx).Create(record).Error == nil {
					res[parameter.ID] = record.ID
				}
			}
		}
	}

	return res
}
