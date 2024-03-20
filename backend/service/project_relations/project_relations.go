package project_relations

import (
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/module/spec"

	"github.com/gin-gonic/gin"
)

// 构造apicat协议结构，填充公共部分数据
func NewApiCatPopulatePublicData(ctx *gin.Context, p *project.Project) *spec.Spec {
	apicatData := spec.NewSpec()
	apicatData.Info = &spec.Info{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}

	// 填充Server数据
	apicatData.Servers = project.ExportServers(ctx, p.ID)

	// 填充GlobalParameter数据
	apicatData.Globals.Parameters = global.ExportGlobalParameters(ctx, p.ID)

	// 填充DefinitionSchema数据
	apicatData.Definitions.Schemas = definition.ExportDefinitionSchemas(ctx, p)
	// 填充DefinitionParameter数据
	apicatData.Definitions.Parameters = definition.ExportDefinitionParameters(ctx, p.ID)
	// 填充DefinitionResponse数据
	apicatData.Definitions.Responses = definition.ExportDefinitionResponses(ctx, p)

	return apicatData
}

func SpecFillInfo(ctx *gin.Context, s *spec.Spec, p *project.Project) {
	s.Info = &spec.Info{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
	}
}

func SpecFillServers(ctx *gin.Context, s *spec.Spec, pID string) {
	s.Servers = project.ExportServers(ctx, pID)
}

func SpecFillGlobals(ctx *gin.Context, s *spec.Spec, pID string) {
	SpecFillGlobalParameters(ctx, s, pID)
}

func SpecFillGlobalParameters(ctx *gin.Context, s *spec.Spec, pID string) {
	s.Globals.Parameters = global.ExportGlobalParameters(ctx, pID)
}

func SpecFillDefinitions(ctx *gin.Context, s *spec.Spec, pID string) {
	SpecFillDefinitionSchemas(ctx, s, pID)
	SpecFillDefinitionParameters(ctx, s, pID)
	SpecFillDefinitionResponses(ctx, s, pID)
}

func SpecFillDefinitionSchemas(ctx *gin.Context, s *spec.Spec, pID string) {
	s.Definitions.Schemas = definition.ExportDefinitionSchemas(ctx, &project.Project{ID: pID})
}

func SpecFillDefinitionParameters(ctx *gin.Context, s *spec.Spec, pID string) {
	s.Definitions.Parameters = definition.ExportDefinitionParameters(ctx, pID)
}

func SpecFillDefinitionResponses(ctx *gin.Context, s *spec.Spec, pID string) {
	s.Definitions.Responses = definition.ExportDefinitionResponses(ctx, &project.Project{ID: pID})
}

func SpecFillCollections(ctx *gin.Context, s *spec.Spec, pID string) {
	s.Collections = collection.ExportCollections(ctx, pID)
}
