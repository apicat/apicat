package content_suggestion

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/content_suggestion/utils"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type definitionModelVector struct {
	projectID      string
	embeddingModel model.Provider
	vectorDB       vector.VectorApi
	allModels      spec.DefinitionModels
}

func NewDefinitionModelVector(projectID string) (*definitionModelVector, error) {
	embeddingModel, err := model.NewModel(config.GetModel().ToModuleStruct("embedding"))
	if err != nil {
		slog.Error("model.NewModel", "err", err)
		return nil, err
	}

	vectorDB, err := vector.NewVector(config.GetVector().ToModuleStruct())
	if err != nil {
		slog.Error("vector.NewVector", "err", err)
		return nil, err
	}

	if ok, err := vectorDB.CheckCollectionExist(projectID); err != nil {
		slog.Error("vectorDB.CheckCollectionExist", "err", err)
		return nil, err
	} else if !ok {
		err = fmt.Errorf("vector db collection not exist, projectID: %s", projectID)
		slog.Error("vectorDB.CheckCollectionExist", "err", err)
		return nil, err
	}

	return &definitionModelVector{
		projectID:      projectID,
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
	}, nil
}

func (dm *definitionModelVector) CreateLater(mid uint) {
	initCache, err := newInitCache(dm.projectID)
	if err != nil {
		slog.Error("newInitCache", "err", err)
		return
	}

	status, err := initCache.GetStatus()
	if err != nil {
		slog.Error("cacheWorker.GetStatus", "err", err)
		return
	}
	if status != "" {
		if err := initCache.SetModelLater(mid); err != nil {
			slog.Error("cacheWorker.SetModelLater", "err", err)
			return
		}
	}

	ca, err := cache.NewCache(config.Get().Cache.ToModuleStruct())
	if err != nil {
		slog.Error("cache.NewCache", "err", err)
		return
	}
	if err := ca.Check(); err != nil {
		slog.Error("cache.Check", "err", err)
		return
	}
	k := fmt.Sprintf("definition_model_vector_create_%s_%d", dm.projectID, mid)
	if err := ca.LPush(k, mid); err != nil {
		slog.Error("cache.LPush", "err", err)
		return
	}
	ca.Expire(k, 10*time.Second)

	go time.AfterFunc(5*time.Second, func() {
		if _, ok, err := ca.RPop(k); err != nil || !ok {
			if err != nil {
				slog.Error("cache.RPop", "err", err)
			} else {
				slog.Debug("cache.RPop", "err", fmt.Errorf("cache.RPop %s not exist", k))
			}
			return
		}

		if len, err := ca.LLen(k); err != nil || len > 0 {
			if err != nil {
				slog.Error("cache.LLen", "err", err)
			} else {
				slog.Debug("cache.LLen", "err", fmt.Errorf("cache.LLen %s(%d) > 0", k, len))
			}
			return
		}

		m := &definition.DefinitionSchema{ID: mid, ProjectID: dm.projectID, Type: definition.SchemaSchema}
		if exist, err := m.GetWithoutCtx(); err != nil {
			slog.Error("m.Get", "err", err)
			return
		} else if !exist {
			slog.Error("m.Get", "err", fmt.Errorf("definition schema id:%d not exist", mid))
			return
		}
		if vid, err := dm.CreateNow(m); err != nil {
			slog.Error("dm.CreateNow", "err", err)
		} else {
			slog.Debug("definition model vector create success", "id", mid, "vid", vid)
		}
	})
}

func (dm *definitionModelVector) CreateNow(ds *definition.DefinitionSchema) (string, error) {
	if ds.VectorID != "" {
		if err := dm.Update(ds); err != nil {
			return "", err
		}
		return ds.VectorID, nil
	}

	if ds.Type != definition.SchemaSchema {
		return "", errors.New("definition type error")
	}
	if ds.ProjectID != dm.projectID {
		return "", errors.New("project id error")
	}
	if ds.Schema == "" {
		return "", errors.New("schema is empty")
	}

	if err := dm.getAllModels(); err != nil {
		return "", err
	}

	embedding, err := dm.createEmbeddings(ds.ID)
	if err != nil {
		slog.Error("dm.createEmbeddings", "err", err)
		return "", err
	}

	properties := &apiContentProperty{
		DefinitionModelID: vector.T_INT(ds.ID),
		UpdatedAt:         vector.T_TEXT(ds.UpdatedAt.Format("2006-01-02 15:04:05")),
	}
	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}

	if vectorID, err := dm.vectorDB.CreateObject(dm.projectID, data); err != nil {
		return "", err
	} else {
		ds.UpdateVectorID(vectorID)
		return vectorID, nil
	}
}

func (dm *definitionModelVector) Update(ds *definition.DefinitionSchema) error {
	if ds.Type != definition.SchemaSchema {
		return errors.New("definition type error")
	}
	if ds.ProjectID != dm.projectID {
		return errors.New("project id error")
	}
	if ds.Schema == "" {
		return errors.New("schema is empty")
	}
	if ds.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := dm.vectorDB.CheckObjectExist(dm.projectID, ds.VectorID); err != nil {
		slog.Error("dm.vectorDB.CheckObjectExist", "err", err)
		return err
	} else if !ok {
		return errors.New("vector object not exist")
	}

	if err := dm.getAllModels(); err != nil {
		return err
	}

	embedding, err := dm.createEmbeddings(ds.ID)
	if err != nil {
		return err
	}

	properties := &apiContentProperty{
		DefinitionModelID: vector.T_INT(ds.ID),
		UpdatedAt:         vector.T_TEXT(ds.UpdatedAt.Format("2006-01-02 15:04:05")),
	}

	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}

	return dm.vectorDB.UpdateObject(dm.projectID, ds.VectorID, data)
}

func (dm *definitionModelVector) Delete(ds *definition.DefinitionSchema) error {
	if ds.Type != definition.SchemaSchema {
		return errors.New("definition type error")
	}
	if ds.ProjectID != dm.projectID {
		return errors.New("project id error")
	}
	if ds.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := dm.vectorDB.CheckObjectExist(dm.projectID, ds.VectorID); err != nil {
		slog.Error("dm.vectorDB.CheckObjectExist", "err", err)
		return err
	} else if !ok {
		return nil
	}

	return dm.vectorDB.DeleteObject(dm.projectID, ds.VectorID)
}

func (dm *definitionModelVector) createEmbeddings(id uint) ([]float32, error) {
	textList := make([]string, 0)

	for _, m := range dm.allModels {
		if m.ID == int64(id) {
			textList = append(textList, m.Name)
			if m.Description != "" {
				textList = append(textList, m.Description)
			}
			textList = append(textList, utils.SchemaToTextList("root", m.Schema)...)
			text := strings.Join(textList, "\n")
			return dm.embeddingModel.CreateEmbeddings(text)
		}
	}
	return nil, errors.New("definition model not found")
}

func (dm *definitionModelVector) getAllModels() error {
	if dm.allModels != nil {
		return nil
	}

	var err error
	dm.allModels, err = definition.GetDefinitionSchemasWithSpec(dm.projectID)
	if err != nil {
		slog.Error("definition.GetDefinitionSchemasWithSpec", "err", err)
		return err
	}
	return nil
}
