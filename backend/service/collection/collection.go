package collection

import (
	"context"
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/content_suggestion"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type CollectionService struct {
	ctx context.Context
}

func NewCollectionService(ctx context.Context) *CollectionService {
	return &CollectionService{ctx: ctx}
}

func (cs *CollectionService) CreateVector(projectID string, cids ...uint) {
	if v, err := content_suggestion.NewApiVector(projectID); err != nil {
		slog.ErrorContext(cs.ctx, "content_suggestion.NewApiVector", "err", err)
	} else {
		for _, cid := range cids {
			v.CreateLater(cid)
		}
	}
}

func (cs *CollectionService) DelVector(c *collection.Collection) error {
	if c == nil {
		return errors.New("collection is nil")
	}
	if c.ProjectID == "" {
		return errors.New("project id is empty")
	}
	if c.VectorID != "" {
		vectorDB, err := vector.NewVector(config.GetVector().ToModuleStruct())
		if err != nil {
			slog.ErrorContext(cs.ctx, "vector.NewVector", "err", err)
			return err
		}
		return vectorDB.DeleteObject(c.ProjectID, c.VectorID)
	}
	return nil
}
