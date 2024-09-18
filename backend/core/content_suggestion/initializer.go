package content_suggestion

import (
	"context"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type VectorInitializer struct {
	ctx                   context.Context
	projectID             string
	initCache             *initCache
	apiVector             *ApiVector
	definitionModelVector *definitionModelVector
}

func NewVectorInitializer(ctx context.Context, projectID string) (*VectorInitializer, error) {
	initCache, err := newInitCache(projectID)
	if err != nil {
		return nil, err
	}

	return &VectorInitializer{
		projectID: projectID,
		initCache: initCache,
		ctx:       ctx,
	}, nil
}

func (vi *VectorInitializer) Run() {
	vectorDB, err := vector.NewVector(config.GetVector().ToModuleStruct())
	if err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.Run", "err", err)
		return
	}
	if exist, err := vectorDB.CheckCollectionExist(vi.projectID); err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.Run", "err", err)
		return
	} else {
		if exist {
			return
		}
	}
	if err := vectorDB.CreateCollection(vi.projectID, getAPIContentProperties()); err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.Run", "err", err)
		return
	}

	if status, err := vi.initCache.GetStatus(); err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.Run", "err", err)
	} else if status != "" {
		slog.DebugContext(vi.ctx, "VectorInitializer.Run", "The VectorInitializer is already running status:", status)
	} else {
		vi.createEmbeddings()
	}
}

func (vi *VectorInitializer) createEmbeddings() {
	collections, err := collection.GetCollections(vi.ctx, vi.projectID)
	if err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.collection.GetCollections", "err", err)
		return
	}
	for _, c := range collections {
		vi.initCache.SetStatus(INIT_RUNNING)
		if c.Type == collection.CategoryType {
			continue
		}
		vi.createCollectionEmbedding(c)
	}

	models, err := definition.GetDefinitionSchemas(vi.ctx, vi.projectID)
	if err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.definition.GetDefinitionSchemas", "err", err)
		return
	}
	for _, m := range models {
		vi.initCache.SetStatus(INIT_RUNNING)
		if m.Type == definition.SchemaCategory {
			continue
		}
		vi.createModelEmbedding(m)
	}

	vi.createCollectionLater()
	vi.createModelLater()

	if err := vi.initCache.Finished(); err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.vi.initCache.Finished", "err", err)
		return
	}
}

func (vi *VectorInitializer) createCollectionEmbedding(c *collection.Collection) {
	var err error
	if vi.apiVector == nil {
		if vi.apiVector, err = NewApiVector(vi.ctx, vi.projectID); err != nil {
			slog.ErrorContext(vi.ctx, "VectorInitializer.NewApiVector", "err", err)
			return
		}
	}

	if _, err := vi.apiVector.CreateNow(c); err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.vi.apiVector.CreateNow", "err", err)
	}
}

func (vi *VectorInitializer) createModelEmbedding(m *definition.DefinitionSchema) {
	var err error
	if vi.definitionModelVector == nil {
		if vi.definitionModelVector, err = NewDefinitionModelVector(vi.ctx, vi.projectID); err != nil {
			slog.ErrorContext(vi.ctx, "VectorInitializer.NewDefinitionModelVector", "err", err)
			return
		}
	}

	if _, err := vi.definitionModelVector.CreateNow(m); err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.vi.definitionModelVector.CreateNow", "err", err)
	}
}

func (vi *VectorInitializer) createCollectionLater() {
	apiVector, err := NewApiVector(vi.ctx, vi.projectID)
	if err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.NewApiVector", "err", err)
		return
	}

	for {
		id, err := vi.initCache.GetCollectionLater()
		if err != nil {
			slog.ErrorContext(vi.ctx, "VectorInitializer.vi.initCache.GetCollectionLater", "err", err)
			return
		}
		if id == 0 {
			return
		}

		c := &collection.Collection{ID: id, ProjectID: vi.projectID, Type: collection.HttpType}
		if exist, err := c.Get(vi.ctx); err != nil || !exist {
			if err != nil {
				slog.ErrorContext(vi.ctx, "VectorInitializer.c.Get", "err", err)
			}
			continue
		}
		if _, err := apiVector.CreateNow(c); err != nil {
			slog.ErrorContext(vi.ctx, "VectorInitializer.apiVector.CreateNow", "err", err)
		}
	}
}

func (vi *VectorInitializer) createModelLater() {
	definitionModelVector, err := NewDefinitionModelVector(vi.ctx, vi.projectID)
	if err != nil {
		slog.ErrorContext(vi.ctx, "VectorInitializer.NewDefinitionModelVector", "err", err)
		return
	}

	for {
		id, err := vi.initCache.GetModelLater()
		if err != nil {
			slog.ErrorContext(vi.ctx, "VectorInitializer.vi.initCache.GetModelLater", "err", err)
			return
		}
		if id == 0 {
			return
		}

		m := &definition.DefinitionSchema{ID: id, ProjectID: vi.projectID, Type: definition.SchemaSchema}
		if exist, err := m.Get(vi.ctx); err != nil || !exist {
			if err != nil {
				slog.ErrorContext(vi.ctx, "VectorInitializer.m.Get", "err", err)
			}
			continue
		}
		if _, err := definitionModelVector.CreateNow(m); err != nil {
			slog.ErrorContext(vi.ctx, "VectorInitializer.definitionModelVector.CreateNow", "err", err)
		}
	}
}
