package except

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/global"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

func DeleteParamReference(ctx context.Context, p *global.GlobalParameter) error {
	pe := referencerelationship.ParameterExcept{ParameterID: p.ID}
	refs, err := pe.GetParameterExcepts(ctx)
	if err != nil {
		return err
	}

	var ids []uint
	for _, item := range refs {
		ids = append(ids, item.ID)
	}

	return referencerelationship.BatchDeleteParameterExcept(ctx, ids...)
}

// 直接删除全局参数
// 查出所有排除该参数的文档，删除文档中exceptParam数组中的id
func ClearParamExcept(ctx context.Context, p *global.GlobalParameter) error {
	pe := referencerelationship.ParameterExcept{ParameterID: p.ID}
	cIDs, err := pe.GetParameterExceptCIDs(ctx)
	if err != nil {
		return err
	}

	collections, err := collection.GetCollections(ctx, p.ProjectID, cIDs...)
	if err != nil {
		return err
	}

	for _, c := range collections {
		if err := c.DelExceptParam(ctx, p, false); err != nil {
			return err
		}
	}

	return DeleteParamReference(ctx, p)
}

// 展开全局参数
// 查出所有排除该参数的文档，删除文档中exceptParam数组中的id，其他未排除的文档在parameters中添加该全局参数
func UnpackParamExcept(ctx context.Context, p *global.GlobalParameter) error {
	collections, err := collection.GetCollections(ctx, p.ProjectID)
	if err != nil {
		return err
	}

	pe := referencerelationship.ParameterExcept{ParameterID: p.ID}
	cIDs, err := pe.GetParameterExceptCIDs(ctx)
	if err != nil {
		return err
	}

	for _, c := range collections {
		if !arrutil.InArray(c.ID, cIDs) {
			if err := c.DelExceptParam(ctx, p, true); err != nil {
				return err
			}
		} else {
			if err := c.DelExceptParam(ctx, p, false); err != nil {
				return err
			}
		}
	}

	return DeleteParamReference(ctx, p)
}
