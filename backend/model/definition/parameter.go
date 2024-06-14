package definition

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
)

const (
	ParameterInHeader = "header"
	ParameterInCookie = "cookie"
	ParameterInQuery  = "query"
	ParameterInPath   = "path"
)

type DefinitionParameter struct {
	ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID string `gorm:"type:varchar(24);index;not null;comment:project id"`
	In        string `gorm:"type:varchar(32);not null;comment:param in:header,cookie,query,path"`
	Name      string `gorm:"type:varchar(255);not null;comment:param name"`
	Required  bool   `gorm:"type:tinyint;not null;comment:is required"`
	Schema    string `gorm:"type:mediumtext;comment:param schema"`
	model.TimeModel
}

// Create 创建全局参数
func (dp *DefinitionParameter) Create(ctx context.Context) error {
	return model.DB(ctx).Create(dp).Error
}

func (dp *DefinitionParameter) ToSpec() (*spec.Parameter, error) {
	p := &spec.Parameter{
		ID:       int64(dp.ID),
		Name:     dp.Name,
		Required: dp.Required,
	}

	if dp.Schema != "" {
		if err := json.Unmarshal([]byte(dp.Schema), &p.Schema); err != nil {
			return nil, err
		}
	}

	return p, nil
}
