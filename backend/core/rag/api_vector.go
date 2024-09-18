package rag

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/rag/utils"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type ApiVector struct {
	ctx                  context.Context
	projectID            string
	embeddingModel       model.Provider
	vectorDB             vector.VectorApi
	specDefinitions      *spec.Definitions
	specGlobalParameters *spec.GlobalParameters
}

func NewApiVector(ctx context.Context, projectID string) (*ApiVector, error) {
	embeddingModel, err := model.NewModel(config.GetModel().ToModuleStruct("embedding"))
	if err != nil {
		slog.ErrorContext(ctx, "model.NewModel", "err", err)
		return nil, err
	}

	vectorDB, err := vector.NewVector(config.GetVector().ToModuleStruct())
	if err != nil {
		slog.ErrorContext(ctx, "vector.NewVector", "err", err)
		return nil, err
	}

	if ok, err := vectorDB.CheckCollectionExist(projectID); err != nil {
		slog.ErrorContext(ctx, "vectorDB.CheckCollectionExist", "err", err)
		return nil, err
	} else if !ok {
		err = fmt.Errorf("vector db collection not exist, projectID: %s", projectID)
		slog.ErrorContext(ctx, "vectorDB.CheckCollectionExist", "err", err)
		return nil, err
	}

	return &ApiVector{
		ctx:            ctx,
		projectID:      projectID,
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
	}, nil
}

func (a *ApiVector) CreateLater(second int, cid uint) {
	initCache, err := newInitCache(a.projectID)
	if err != nil {
		slog.ErrorContext(a.ctx, "newInitCache", "err", err)
		return
	}

	status, err := initCache.GetStatus()
	if err != nil {
		slog.ErrorContext(a.ctx, "cacheWorker.GetStatus", "err", err)
		return
	}
	if status != "" {
		if err := initCache.SetCollectionLater(cid); err != nil {
			slog.ErrorContext(a.ctx, "cacheWorker.SetCollectionLater", "err", err)
			return
		}
	}

	go time.AfterFunc(time.Duration(second)*time.Second, func() {
		doc := &collection.Collection{ID: cid, ProjectID: a.projectID, Type: collection.HttpType}
		if exist, err := doc.Get(a.ctx); err != nil {
			slog.ErrorContext(a.ctx, "doc.Get", "err", err)
			return
		} else if !exist {
			slog.ErrorContext(a.ctx, "doc.Get", "err", fmt.Errorf("doc id:%d not exist", cid))
			return
		}
		if _, err := a.CreateNow(doc); err != nil {
			slog.ErrorContext(a.ctx, "a.Create", "err", err)
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
		doc.UpdateVectorID(a.ctx, vectorID)
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
		slog.ErrorContext(a.ctx, "ad.doc.ContentToSpec", "err", err)
		return nil, err
	}

	if err := specContent.DeepDerefAll(a.specGlobalParameters, a.specDefinitions); err != nil {
		slog.ErrorContext(a.ctx, "specContent.DeepDerefAll", "err", err)
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
	a.specDefinitions.Schemas, err = definition.GetDefinitionSchemasWithSpec(a.ctx, a.projectID)
	if err != nil {
		return err
	}
	a.specDefinitions.Responses, err = definition.GetDefinitionResponsesWithSpec(a.ctx, a.projectID)
	if err != nil {
		return err
	}

	a.specGlobalParameters, err = global.GetGlobalParametersWithSpec(a.ctx, a.projectID)
	if err != nil {
		return err
	}
	return nil
}
