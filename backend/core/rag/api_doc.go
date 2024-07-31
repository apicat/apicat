package rag

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/core/rag/utils"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type apiDoc struct {
	doc            *collection.Collection
	embeddingModel model.Provider
	vectorDB       vector.VectorApi
	ctx            context.Context
}

func NewApiDoc(ctx context.Context, doc *collection.Collection) (*apiDoc, error) {
	if doc.Type != collection.HttpType {
		return nil, errors.New("collection type error")
	}

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

	if ok, _ := vectorDB.CheckCollectionExist(doc.ProjectID); !ok {
		if err := vectorDB.CreateCollection(doc.ProjectID, getAPIContentProperties()); err != nil {
			slog.ErrorContext(ctx, "vectorDB.CreateCollection", "err", err)
			return nil, err
		}
	}

	return &apiDoc{
		doc:            doc,
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
		ctx:            ctx,
	}, nil
}

func (ad *apiDoc) Create() (string, error) {
	if ad.doc.Content == "" {
		return "", errors.New("content is empty")
	}

	embedding, err := ad.createEmbeddings()
	if err != nil {
		return "", err
	}

	properties := &apiContentProperty{
		CollectionID: vector.T_INT(ad.doc.ID),
		UpdatedAt:    vector.T_TEXT(ad.doc.UpdatedAt.Format("2006-01-02 15:04:05")),
	}
	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}
	return ad.vectorDB.CreateObject(ad.doc.ProjectID, data)
}

func (ad *apiDoc) Update() error {
	if ad.doc.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := ad.vectorDB.CheckObjectExist(ad.doc.ProjectID, ad.doc.VectorID); err != nil || !ok {
		slog.ErrorContext(ad.ctx, "ad.vectorDB.CheckObjectExist", "err", err)
		return errors.New("vector object not exist")
	}

	embedding, err := ad.createEmbeddings()
	if err != nil {
		return err
	}

	properties := &apiContentProperty{
		CollectionID: vector.T_INT(ad.doc.ID),
		UpdatedAt:    vector.T_TEXT(ad.doc.UpdatedAt.Format("2006-01-02 15:04:05")),
	}

	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}

	return ad.vectorDB.UpdateObject(ad.doc.ProjectID, ad.doc.VectorID, data)
}

func (ad *apiDoc) Delete() error {
	if ad.doc.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := ad.vectorDB.CheckObjectExist(ad.doc.ProjectID, ad.doc.VectorID); err != nil || !ok {
		slog.ErrorContext(ad.ctx, "ad.vectorDB.CheckObjectExist", "err", err)
		return errors.New("vector object not exist")
	}

	return ad.vectorDB.DeleteObject(ad.doc.ProjectID, ad.doc.VectorID)
}

func (ad *apiDoc) createEmbeddings() ([]float32, error) {
	specContent, err := ad.doc.ContentToSpec()
	if err != nil {
		slog.ErrorContext(ad.ctx, "ad.doc.ContentToSpec", "err", err)
		return nil, err
	}

	textList := make([]string, 0)
	textList = append(textList, ad.doc.Title)
	if ad.doc.Path != "" {
		textList = append(textList, ad.doc.Path)
	}

	specDefinitions := &spec.Definitions{}
	specDefinitions.Schemas, err = definition.GetDefinitionSchemasWithSpec(ad.ctx, ad.doc.ProjectID)
	if err != nil {
		slog.ErrorContext(ad.ctx, "definition.GetDefinitionSchemasWithSpec", "err", err)
		return nil, err
	}
	specDefinitions.Responses, err = definition.GetDefinitionResponsesWithSpec(ad.ctx, ad.doc.ProjectID)
	if err != nil {
		slog.ErrorContext(ad.ctx, "definition.GetDefinitionResponsesWithSpec", "err", err)
		return nil, err
	}

	specGlobalParameters, err := global.GetGlobalParametersWithSpec(ad.ctx, ad.doc.ProjectID)
	if err != nil {
		slog.ErrorContext(ad.ctx, "global.GetGlobalParametersWithSpec", "err", err)
		return nil, err
	}

	if err := specContent.DeepDerefAll(specGlobalParameters, specDefinitions); err != nil {
		slog.ErrorContext(ad.ctx, "specContent.DeepDerefAll", "err", err)
		return nil, err
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
	return ad.embeddingModel.CreateEmbeddings(text)
}
