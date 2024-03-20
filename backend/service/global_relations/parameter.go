package globalrelations

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/global"
	"github.com/apicat/apicat/backend/model/project"
	referencerelationship "github.com/apicat/apicat/backend/model/reference_relationship"
	"github.com/apicat/apicat/backend/module/array_operation"

	"github.com/apicat/apicat/backend/module/spec"

	"github.com/gin-gonic/gin"
)

// ReadExceptParameterReference 读取collection中排除的全局参数
func ReadExceptParameterReference(ctx context.Context, content string) []uint {
	list := make([]uint, 0)

	var specContent []*spec.NodeProxy
	if err := json.Unmarshal([]byte(content), &specContent); err != nil {
		return list
	}

	var request *spec.HTTPNode[spec.HTTPRequestNode]
	for _, i := range specContent {
		switch nx := i.Node.(type) {
		case *spec.HTTPNode[spec.HTTPRequestNode]:
			request = nx
		}
	}

	if request == nil {
		return list
	}

	for key, value := range request.Attrs.GlobalExcepts {
		if !array_operation.InArray[string](key, []string{string(global.ParameterInHeader), string(global.ParameterInPath), string(global.ParameterInHeader), string(global.ParameterInCookie)}) {
			continue
		}
		for _, v := range value {
			list = append(list, uint(v))
		}
	}

	return list
}

// dereferenceGlobalParameterInCollection 将全局参数展开后添加到request.parameters中
func dereferenceGlobalParameterInCollection(ctx *gin.Context, c *collection.Collection, gp *global.GlobalParameter) error {
	if c == nil || c.Type == collection.CategoryType {
		return nil
	}

	collection := &spec.Collection{
		ID:    c.ID,
		Title: c.Title,
		Type:  spec.CollectionType(c.Type),
	}
	if err := json.Unmarshal([]byte(c.Content), &collection.Content); err != nil {
		return err
	}

	parameter := &spec.Parameter{
		ID:       int64(gp.ID),
		Name:     gp.Name,
		Required: gp.Required,
	}
	if err := json.Unmarshal([]byte(gp.Schema), &parameter.Schema); err != nil {
		return err
	}

	collection.AddParameter(gp.In, parameter)
	content, err := json.Marshal(collection.Content)
	if err != nil {
		return err
	}
	c.Content = string(content)
	return nil
}

// removeExceptGlobalParameterInCollection 将全局参数从request.globalExcepts中删除
func removeExceptGlobalParameterInCollection(ctx *gin.Context, c *collection.Collection, gp *global.GlobalParameter) error {
	if c == nil || c.Type == collection.CategoryType {
		return nil
	}

	collection := &spec.Collection{
		ID:    c.ID,
		Title: c.Title,
		Type:  spec.CollectionType(c.Type),
	}
	if err := json.Unmarshal([]byte(c.Content), &collection.Content); err != nil {
		return err
	}

	collection.DelGlobalExceptID(gp.In, int64(gp.ID))
	content, err := json.Marshal(collection.Content)
	if err != nil {
		return err
	}
	c.Content = string(content)
	return nil
}

// RemoveExceptGlobalParameter 清除全局参数，将修改数据库中collection.content
func RemoveExceptGlobalParameter(ctx *gin.Context, gp *global.GlobalParameter) error {
	gpRecords, err := referencerelationship.GetParameterExceptsByParameter(ctx, gp.ProjectID, gp.ID)
	if err != nil {
		return err
	}
	if len(gpRecords) == 0 {
		return nil
	}

	exceptGPIDs := make([]uint, 0)
	var collectionIDs []uint
	for _, ref := range gpRecords {
		exceptGPIDs = append(exceptGPIDs, ref.ID)
		collectionIDs = append(collectionIDs, ref.ExceptCollectionID)
	}

	collections, err := collection.GetCollections(ctx, &project.Project{ID: gp.ProjectID}, collectionIDs...)
	if err != nil {
		return err
	}
	for _, c := range collections {
		if err := removeExceptGlobalParameterInCollection(ctx, c, gp); err != nil {
			return err
		}
		if err := c.Update(ctx, c.Title, c.Content, c.UpdatedBy); err != nil {
			return err
		}
	}

	return referencerelationship.BatchDeleteParameterExcept(ctx, exceptGPIDs...)
}

// UnpackExceptGlobalParameter 展开全局参数，将修改数据库中collection.content
func UnpackExceptGlobalParameter(ctx *gin.Context, gp *global.GlobalParameter) error {
	gpRecords, err := referencerelationship.GetParameterExceptsByParameter(ctx, gp.ProjectID, gp.ID)
	if err != nil {
		return err
	}
	exceptParamIDs := make([]uint, 0)
	exceptCollectionIDs := make([]uint, 0)
	for _, v := range gpRecords {
		exceptParamIDs = append(exceptParamIDs, v.ID)
		exceptCollectionIDs = append(exceptCollectionIDs, v.ExceptCollectionID)
	}

	collections, err := collection.GetCollections(ctx, &project.Project{ID: gp.ProjectID})
	if err != nil {
		return err
	}
	for _, c := range collections {
		if !array_operation.InArray(c.ID, exceptCollectionIDs) {
			// 将没有被排除的全局参数展开后添加到request.parameters中
			dereferenceGlobalParameterInCollection(ctx, c, gp)
		} else {
			// 将被排除的全局参数从request.globalExcepts中删除
			removeExceptGlobalParameterInCollection(ctx, c, gp)
		}
		if err := c.Update(ctx, c.Title, c.Content, c.UpdatedBy); err != nil {
			return err
		}
	}

	return referencerelationship.BatchDeleteParameterExcept(ctx, exceptParamIDs...)
}

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
