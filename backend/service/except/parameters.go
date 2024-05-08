package except

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/global"
	referencerelationship "github.com/apicat/apicat/v2/backend/model/reference_relationship"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

// DerefParamExcept 解引用全局参数
func DerefParamExcept(ctx context.Context, p *global.GlobalParameter, deref bool) error {
	if deref {
		return unpackParamExcept(ctx, p)
	} else {
		return clearParamExcept(ctx, p)
	}
}

// clearParamExcept 清除全局参数排除
func clearParamExcept(ctx context.Context, p *global.GlobalParameter) error {
	// 查出所有排除该参数的文档
	pe := referencerelationship.ParameterExcept{ParameterID: p.ID}
	cIDs, err := pe.GetParameterExceptCIDs(ctx)
	if err != nil {
		return err
	}

	collections, err := collection.GetCollections(ctx, p.ProjectID, cIDs...)
	if err != nil {
		return err
	}

	// 删除文档中exceptParam数组中的id
	for _, c := range collections {
		if err := c.DelExceptParam(ctx, p, false); err != nil {
			return err
		}
	}

	// 清除排除关系(self -> collections)
	return deleteParamReference(ctx, p)
}

// unpackParamExcept 展开全局参数
func unpackParamExcept(ctx context.Context, p *global.GlobalParameter) error {
	// 查出所有文档
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
		if c.Type == collection.CategoryType {
			continue
		}

		if !arrutil.InArray(c.ID, cIDs) {
			// 未排除的文档在parameters中添加该全局参数
			if err := c.DelExceptParam(ctx, p, true); err != nil {
				return err
			}
		} else {
			// 排除的文档删除文档中exceptParam数组中的id
			if err := c.DelExceptParam(ctx, p, false); err != nil {
				return err
			}
		}
	}

	// 清除排除关系(self -> collections)
	return deleteParamReference(ctx, p)
}

// deleteParamReference 删除全局参数排除关系
func deleteParamReference(ctx context.Context, p *global.GlobalParameter) error {
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
