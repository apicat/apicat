package content_suggestion

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/content_suggestion/utils"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type ApiVector struct {
	projectID            string
	embeddingModel       model.Provider
	vectorDB             vector.VectorApi
	specDefinitions      *spec.Definitions
	specGlobalParameters *spec.GlobalParameters
}

func NewApiVector(projectID string) (*ApiVector, error) {
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

	return &ApiVector{
		projectID:      projectID,
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
	}, nil
}

func (a *ApiVector) CreateLater(cid uint) {
	initCache, err := newInitCache(a.projectID)
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
		if err := initCache.SetCollectionLater(cid); err != nil {
			slog.Error("cacheWorker.SetCollectionLater", "err", err)
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
	k := fmt.Sprintf("api_vector_create_%s_%d", a.projectID, cid)
	if err := ca.LPush(k, cid); err != nil {
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

		doc := &collection.Collection{ID: cid, ProjectID: a.projectID, Type: collection.HttpType}
		if exist, err := doc.GetWithoutCtx(); err != nil {
			slog.Error("doc.Get", "err", err)
			return
		} else if !exist {
			slog.Error("doc.Get", "err", fmt.Errorf("doc id:%d not exist", cid))
			return
		}
		if vid, err := a.CreateNow(doc); err != nil {
			slog.Error("a.CreateNow", "err", err)
		} else {
			slog.Debug("api vector create success", "id", cid, "vid", vid)
		}
	})
}

func (a *ApiVector) CreateNow(doc *collection.Collection) (string, error) {
	if doc.VectorID != "" {
		if err := a.Update(doc); err != nil {
			return "", err
		}
		return doc.VectorID, nil
	}

	if doc.Type != collection.HttpType {
		return "", errors.New("collection type error")
	}
	if doc.ProjectID != a.projectID {
		return "", errors.New("project id error")
	}
	if doc.Content == "" {
		return "", errors.New("content is empty")
	}

	if err := a.getAllReferences(); err != nil {
		return "", err
	}

	embedding, err := a.createEmbeddings(doc)
	if err != nil {
		return "", err
	}

	properties := &apiContentProperty{
		CollectionID: vector.T_INT(doc.ID),
		UpdatedAt:    vector.T_TEXT(doc.UpdatedAt.Format("2006-01-02 15:04:05")),
	}
	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}

	if vectorID, err := a.vectorDB.CreateObject(a.projectID, data); err != nil {
		return "", err
	} else {
		doc.UpdateVectorID(vectorID)
		return vectorID, nil
	}
}

func (a *ApiVector) Update(doc *collection.Collection) error {
	if doc.Type != collection.HttpType {
		return errors.New("collection type error")
	}
	if doc.ProjectID != a.projectID {
		return errors.New("project id error")
	}
	if doc.VectorID == "" {
		return errors.New("vector id is empty")
	}
	if doc.Content == "" {
		return errors.New("content is empty")
	}

	if ok, err := a.vectorDB.CheckObjectExist(a.projectID, doc.VectorID); err != nil {
		return err
	} else if !ok {
		return errors.New("vector object not exist")
	}

	if err := a.getAllReferences(); err != nil {
		return err
	}

	embedding, err := a.createEmbeddings(doc)
	if err != nil {
		return err
	}

	properties := &apiContentProperty{
		CollectionID: vector.T_INT(doc.ID),
		UpdatedAt:    vector.T_TEXT(doc.UpdatedAt.Format("2006-01-02 15:04:05")),
	}

	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}

	return a.vectorDB.UpdateObject(a.projectID, doc.VectorID, data)
}

func (a *ApiVector) Delete(doc *collection.Collection) error {
	if doc.Type != collection.HttpType {
		return errors.New("collection type error")
	}
	if doc.ProjectID != a.projectID {
		return errors.New("project id error")
	}
	if doc.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := a.vectorDB.CheckObjectExist(a.projectID, doc.VectorID); err != nil {
		return err
	} else if !ok {
		return errors.New("vector object not exist")
	}

	return a.vectorDB.DeleteObject(a.projectID, doc.VectorID)
}

func (a *ApiVector) createEmbeddings(doc *collection.Collection) ([]float32, error) {
	specContent, err := doc.ContentToSpec()
	if err != nil {
		slog.Error("ad.doc.ContentToSpec", "err", err)
		return nil, err
	}

	if err := specContent.DeepDerefAll(a.specGlobalParameters, a.specDefinitions); err != nil {
		slog.Error("specContent.DeepDerefAll", "err", err)
		return nil, err
	}

	textList := make([]string, 0)
	textList = append(textList, doc.Title)
	if doc.Path != "" {
		textList = append(textList, doc.Path)
	}

	request := specContent.GetRequest()
	if request.Attrs.Parameters != nil {
		if len(request.Attrs.Parameters.Header) > 0 {
			textList = append(textList, utils.ParamsToTextList("header", request.Attrs.Parameters.Header)...)
		}
		if len(request.Attrs.Parameters.Cookie) > 0 {
			textList = append(textList, utils.ParamsToTextList("cookie", request.Attrs.Parameters.Cookie)...)
		}
		if len(request.Attrs.Parameters.Path) > 0 {
			textList = append(textList, utils.ParamsToTextList("path", request.Attrs.Parameters.Path)...)
		}
		if len(request.Attrs.Parameters.Query) > 0 {
			textList = append(textList, utils.ParamsToTextList("query", request.Attrs.Parameters.Query)...)
		}
	}
	if request.Attrs.Content != nil {
		textList = append(textList, utils.HTTPBodyToTextList(request.Attrs.Content)...)
	}

	response := specContent.GetResponse()
	if len(response.Attrs.List) > 0 {
		for _, v := range response.Attrs.List {
			textList = append(textList, utils.HTTPBodyToTextList(v.Content)...)
		}
	}

	text := strings.Join(textList, "\n")
	return a.embeddingModel.CreateEmbeddings(text)
}

func (a *ApiVector) getAllReferences() error {
	if a.specDefinitions != nil && a.specGlobalParameters != nil {
		return nil
	}

	var err error

	a.specDefinitions = &spec.Definitions{}
	a.specDefinitions.Schemas, err = definition.GetDefinitionSchemasWithSpec(a.projectID)
	if err != nil {
		return err
	}
	a.specDefinitions.Responses, err = definition.GetDefinitionResponsesWithSpec(a.projectID)
	if err != nil {
		return err
	}

	a.specGlobalParameters, err = global.GetGlobalParametersWithSpec(a.projectID)
	if err != nil {
		return err
	}
	return nil
}
