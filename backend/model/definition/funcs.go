package definition

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/model/user"
	"github.com/apicat/apicat/v2/backend/module/spec"

	"gorm.io/gorm"
)

func GetDefinitionResponses(ctx context.Context, projectID string, drIDs ...uint) ([]*DefinitionResponse, error) {
	var list []*DefinitionResponse
	tx := model.DB(ctx)
	if len(drIDs) > 0 {
		tx = tx.Where("id IN ?", drIDs)
	}
	err := tx.Where("project_id = ?", projectID).Order("display_order asc").Find(&list).Error
	return list, err
}

func GetDefinitionResponsesWithSpec(ctx context.Context, projectID string) (spec.DefinitionResponses, error) {
	var list []*DefinitionResponse
	err := model.DB(ctx).Where("project_id = ? AND type = ?", projectID, ResponseResponse).Find(&list).Error
	if err != nil {
		return nil, err
	}

	specResponses := make(spec.DefinitionResponses, 0)
	if len(list) > 0 {
		for _, r := range list {
			if specResponse, err := r.ToSpec(); err == nil {
				specResponses = append(specResponses, specResponse)
			} else {
				return nil, err
			}
		}
	}
	return specResponses, err
}

func ExportDefinitionResponses(ctx context.Context, p *project.Project) spec.DefinitionResponses {
	result := make(spec.DefinitionResponses, 0)

	response, err := GetDefinitionResponses(ctx, p.ID)
	if err != nil {
		slog.ErrorContext(ctx, "GetDefinitionResponses", "err", err)
		return result
	}

	return exportBuildDefinitionResponseTree(ctx, response, 0)
}

func exportBuildDefinitionResponseTree(ctx context.Context, responses []*DefinitionResponse, parentID uint) spec.DefinitionResponses {
	result := make(spec.DefinitionResponses, 0)

	for _, r := range responses {
		if r.ParentID == parentID {
			children := exportBuildDefinitionResponseTree(ctx, responses, r.ID)

			specresponse := &spec.DefinitionResponse{
				BasicResponse: spec.BasicResponse{
					ID:          int64(r.ID),
					Name:        r.Name,
					Description: r.Description,
				},
				Type:     string(r.Type),
				ParentId: int64(r.ParentID),
			}

			if r.Type == ResponseCategory {
				specresponse.Items = children
			} else {
				if err := json.Unmarshal([]byte(r.Header), &specresponse.Header); err != nil {
					continue
				}
				if err := json.Unmarshal([]byte(r.Content), &specresponse.Content); err != nil {
					continue
				}
			}

			result = append(result, specresponse)
		}
	}

	return result
}

func GetDefinitionSchemas(ctx context.Context, projectID string, dsIDs ...uint) ([]*DefinitionSchema, error) {
	var list []*DefinitionSchema
	tx := model.DB(ctx)
	if len(dsIDs) > 0 {
		tx = tx.Where("id IN ?", dsIDs)
	}
	err := tx.Where("project_id = ?", projectID).Order("display_order asc").Find(&list).Error
	return list, err
}

func GetDefinitionSchemasWithSpec(ctx context.Context, projectID string) (spec.DefinitionModels, error) {
	var list []*DefinitionSchema
	err := model.DB(ctx).Where("project_id = ? AND type = ?", projectID, SchemaSchema).Find(&list).Error
	if err != nil {
		return nil, err
	}

	specSchemas := make(spec.DefinitionModels, 0)
	if len(list) > 0 {
		for _, s := range list {
			if specSchema, err := s.ToSpec(); err == nil {
				specSchemas = append(specSchemas, specSchema)
			} else {
				return nil, err
			}
		}
	}
	return specSchemas, err
}

func GetDefinitionParameters(ctx context.Context, pID string) ([]*DefinitionParameter, error) {
	var list []*DefinitionParameter
	err := model.DB(ctx).Where("project_id = ?", pID).Find(&list).Error
	return list, err
}

func ExportDefinitionSchemas(ctx context.Context, p *project.Project) spec.DefinitionModels {
	result := make(spec.DefinitionModels, 0)

	schemas, err := GetDefinitionSchemas(ctx, p.ID)
	if err != nil {
		slog.ErrorContext(ctx, "GetDefinitionSchemas", "err", err)
		return result
	}

	return exportBuildDefinitionSchemaTree(ctx, schemas, 0)
}

func exportBuildDefinitionSchemaTree(ctx context.Context, schemas []*DefinitionSchema, parentID uint) spec.DefinitionModels {
	result := make(spec.DefinitionModels, 0)

	for _, s := range schemas {
		if s.ParentID == parentID {
			children := exportBuildDefinitionSchemaTree(ctx, schemas, s.ID)

			specschema, err := s.ToSpec()
			if err != nil {
				continue
			}
			if s.Type == SchemaCategory {
				specschema.Items = children
			}

			result = append(result, specschema)
		}
	}

	return result
}

// func ExportDefinitionParameters(ctx context.Context, projectID string) *spec.HTTPParameters {
// 	res := &spec.HTTPParameters{}
// 	res.Fill()

// 	parameters, err := GetDefinitionParameters(ctx, projectID)
// 	if err != nil {
// 		return res
// 	}

// 	for _, parameter := range parameters {
// 		schema := &jsonschema.Schema{}
// 		if err := json.Unmarshal([]byte(parameter.Schema), schema); err == nil {
// 			res.Add(string(parameter.In), &spec.Parameter{
// 				ID:       int64(parameter.ID),
// 				Name:     parameter.Name,
// 				Required: parameter.Required,
// 				Schema:   schema,
// 			})
// 		}
// 	}

// 	return res
// }

func GetDefinitionSchemaHistories(ctx context.Context, ds *DefinitionSchema, start, end time.Time) ([]*DefinitionSchemaHistory, error) {
	var list []*DefinitionSchemaHistory
	tx := model.DB(ctx)
	if !start.IsZero() && !end.IsZero() {
		tx = tx.Where("created_at BETWEEN ? AND ?", start, end)
	}
	return list, tx.Where("schema_id = ?", ds.ID).Order("created_at desc").Find(&list).Error
}

func MemberInfo(ctx context.Context, memberID uint, unscoped bool) (*team.TeamMember, error) {
	var (
		tm *team.TeamMember
		tx *gorm.DB
	)

	if unscoped {
		tx = model.DB(ctx).Unscoped()
	} else {
		tx = model.DB(ctx)
	}

	if err := tx.First(&tm, memberID).Error; err != nil {
		return nil, err
	} else {
		return tm, nil
	}
}

func UserInfo(ctx context.Context, memberID uint, unscoped bool) (*user.User, error) {
	if tm, err := MemberInfo(ctx, memberID, unscoped); err != nil {
		return nil, err
	} else {
		return tm.UserInfo(ctx, unscoped)
	}
}
