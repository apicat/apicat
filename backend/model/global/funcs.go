package global

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"

	"gorm.io/gorm"
)

// GetGlobalParameters 获取项目的全局参数
func GetGlobalParameters(ctx context.Context, pID string) ([]*GlobalParameter, error) {
	var list []*GlobalParameter
	err := model.DB(ctx).Where("project_id = ?", pID).Order("display_order asc").Find(&list).Error
	return list, err
}

func GetGlobalParametersWithSpec(pID string) (*spec.GlobalParameters, error) {
	var list []*GlobalParameter
	err := model.DBWithoutCtx().Where("project_id = ?", pID).Order("display_order asc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	specParameters := spec.NewGlobalParameters()
	if len(list) > 0 {
		for _, gp := range list {
			if specParameter, err := gp.ToSpec(); err == nil {
				specParameters.Add(gp.In, specParameter)
			} else {
				return nil, err
			}
		}
	}
	return specParameters, nil
}

func ExportGlobalParameters(ctx context.Context, projectID string) *spec.GlobalParameters {
	res := spec.NewGlobalParameters()

	parameters, err := GetGlobalParameters(ctx, projectID)
	if err != nil {
		return res
	}

	for _, parameter := range parameters {
		if schema, err := jsonschema.NewSchemaFromJson(parameter.Schema); err == nil {
			res.Add(parameter.In, &spec.Parameter{
				ID:       int64(parameter.ID),
				Name:     parameter.Name,
				Required: parameter.Required,
				Schema:   schema,
			})
		}
	}

	return res
}

// SortGlobalParameters 排序
func SortGlobalParameters(ctx context.Context, pID, in string, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return model.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&GlobalParameter{}).Where("id = ? AND project_id = ? AND `in` = ?", id, pID, in).Update("display_order", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
