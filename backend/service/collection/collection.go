package collection

import (
	"context"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/core/content_suggestion"
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
