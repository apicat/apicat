package definition

import (
	"context"
	"encoding/json"

	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/module/spec"
)

const (
	ParameterInHeader = "header"
	ParameterInCookie = "cookie"
	ParameterInQuery  = "query"
	ParameterInPath   = "path"
)

type DefinitionParameter struct {
	ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID string `gorm:"type:varchar(24);index;not null;comment:项目id"`
	In        string `gorm:"type:varchar(32);not null;comment:位置:header,cookie,query,path"`
	Name      string `gorm:"type:varchar(255);not null;comment:参数名称"`
	Required  bool   `gorm:"type:tinyint;not null;comment:是否必传"`
	Schema    string `gorm:"type:mediumtext;comment:参数内容"`
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
