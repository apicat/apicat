package relations

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/spec"
)

// ImportGlobalParameters 导入全局参数
func ImportGlobalParameters(ctx context.Context, projectID string, parameters *spec.HTTPParameters) map[int64]uint {
	res := collection.VirtualIDToIDMap{}

	if parameters.Header == nil && parameters.Cookie == nil && parameters.Query == nil && parameters.Path == nil {
		return res
	}

	var params spec.ParameterList
	parameterList := []string{global.ParameterInHeader, global.ParameterInCookie, global.ParameterInQuery, global.ParameterInPath}
	for _, key := range parameterList {
		switch key {
		case global.ParameterInHeader:
			params = parameters.Header
		case global.ParameterInCookie:
			params = parameters.Cookie
		case global.ParameterInQuery:
			params = parameters.Query
		case global.ParameterInPath:
			params = parameters.Path
		}

		for _, parameter := range params {
			if parameterStr, err := json.Marshal(parameter.Schema); err == nil {
				record := &global.GlobalParameter{
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
