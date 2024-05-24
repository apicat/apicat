package except

import (
	"context"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/global"
	referencerelation "github.com/apicat/apicat/v2/backend/model/reference_relation"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
)

func DerefExceptParam(ctx context.Context, p *global.GlobalParameter, deref bool) error {
	pc := referencerelation.ExceptParamCollection{ExceptParamID: p.ID}
	cIDs, err := pc.GetCollectionIDs(ctx)
	if err != nil {
		return err
	}

	if deref {
		collections, err := collection.GetCollections(ctx, p.ProjectID)
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
	} else {
		if len(cIDs) > 0 {
			collections, err := collection.GetCollections(ctx, p.ProjectID, cIDs...)
			if err != nil {
				return err
			}

			for _, c := range collections {
				if err := c.DelExceptParam(ctx, p, false); err != nil {
					return err
				}
			}
		}
	}

	// 清除排除关系(self -> collections)
	return clearExceptCollectionToParam(ctx, p.ID)
}

func clearExceptCollectionToParam(ctx context.Context, pID uint) error {
	return referencerelation.DelExceptParamCollections(ctx, pID)
}
