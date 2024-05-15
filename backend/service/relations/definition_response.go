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

// ImportDefinitionResponses 导入公共响应
func ImportDefinitionResponses(ctx context.Context, projectID string, responses spec.DefinitionResponses, tm *team.TeamMember, schemaVirtualIDToID collection.VirtualIDToIDMap, parentID uint) collection.VirtualIDToIDMap {
	res := collection.VirtualIDToIDMap{}
	if len(responses) == 0 {
		return res
	}

	var emptySlice []uint
	for i, response := range responses {
		record := &definition.DefinitionResponse{
			ProjectID:    projectID,
			ParentID:     parentID,
			Name:         response.Name,
			DisplayOrder: uint(i),
		}

		if len(response.Items) > 0 || response.Type == definition.ResponseCategory {
			record.Type = definition.ResponseCategory
			if err := record.Create(ctx, tm); err == nil {
				res[response.ID] = record.ID
				res.Merge(ImportDefinitionResponses(ctx, projectID, response.Items, tm, schemaVirtualIDToID, record.ID))
			}
		} else {
			var (
				header  string
				content string
			)
			if headerStr, err := json.Marshal(response.Header); err == nil {
				header = string(headerStr)
			}
			if contentStr, err := json.Marshal(response.Content); err == nil {
				content = collection.ReplaceVirtualIDToID(string(contentStr), schemaVirtualIDToID, "#/definitions/schemas/")
			}

			record.Type = definition.ResponseResponse
			record.Description = response.Description
			record.Header = header
			record.Content = content
			if err := record.Create(ctx, tm); err == nil {
				res[response.ID] = record.ID
			}

			if err := reference.UpdateResponseRef(ctx, record, emptySlice); err != nil {
				slog.ErrorContext(ctx, "ImportDefinitionResponses.UpdateResponseRef", "err", err)
			}
		}
	}

	return res
}
