package relations

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/service/reference"
)

// ImportDefinitionSchemas 导入公共模型
func ImportDefinitionSchemas(ctx context.Context, projectID string, schemas spec.Schemas, tm *team.TeamMember, parentID uint) collection.VirtualIDToIDMap {
	res := collection.VirtualIDToIDMap{}
	if len(schemas) == 0 {
		return res
	}

	for i, schema := range schemas {
		record := &definition.DefinitionSchema{
			ProjectID:    projectID,
			ParentID:     parentID,
			Name:         schema.Name,
			DisplayOrder: uint(i),
		}

		if len(schema.Items) > 0 || schema.Type == definition.SchemaCategory {
			record.Type = definition.SchemaCategory
			if err := record.Create(ctx, tm); err == nil {
				res[schema.ID] = record.ID
				res.Merge(ImportDefinitionSchemas(ctx, projectID, schema.Items, tm, record.ID))
			}
		} else {
			if schemaStr, err := json.Marshal(schema.Schema); err == nil {
				record.Type = definition.SchemaSchema
				record.Description = schema.Description
				record.Schema = string(schemaStr)
				if err := record.Create(ctx, tm); err == nil {
					res[schema.ID] = record.ID
				}
			}
		}
	}

	list, err := definition.GetDefinitionSchemas(ctx, projectID)
	if err != nil {
		slog.ErrorContext(ctx, "GetDefinitionSchemas", "err", err)
		return res
	}

	for _, v := range list {
		schema := collection.ReplaceVirtualIDToID(v.Schema, res, "#/definitions/schemas/")
		v.Update(ctx, v.Name, v.Description, schema, tm.ID)

		if err := reference.UpdateSchemaRef(ctx, v); err != nil {
			slog.ErrorContext(ctx, "ImportDefinitionSchemas.UpdateSchemaRef", "err", err)
		}
	}

	return res
}
