package rag

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model/definition"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type definitionModel struct {
	schema         *definition.DefinitionSchema
	embeddingModel model.Provider
	vectorDB       vector.VectorApi
	ctx            context.Context
}

func NewDefinitionModel(ctx context.Context, ds *definition.DefinitionSchema) (*definitionModel, error) {
	if ds.Type != definition.SchemaSchema {
		return nil, errors.New("definition type error")
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

	if ok, _ := vectorDB.CheckCollectionExist(ds.ProjectID); !ok {
		if err := vectorDB.CreateCollection(ds.ProjectID, getAPIContentProperties()); err != nil {
			slog.ErrorContext(ctx, "vectorDB.CreateCollection", "err", err)
			return nil, err
		}
	}

	return &definitionModel{
		schema:         ds,
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
		ctx:            ctx,
	}, nil
}

func (dm *definitionModel) Create() (string, error) {
	if dm.schema.Schema == "" {
		return "", errors.New("schema is empty")
	}

	embedding, err := dm.createEmbeddings()
	if err != nil {
		slog.ErrorContext(dm.ctx, "dm.createEmbeddings", "err", err)
		return "", err
	}

	properties := &apiContentProperty{
		DefinitionModelID: vector.T_INT(dm.schema.ID),
		UpdatedAt:         vector.T_TEXT(dm.schema.UpdatedAt.Format("2006-01-02 15:04:05")),
	}
	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}
	return dm.vectorDB.CreateObject(dm.schema.ProjectID, data)
}

func (dm *definitionModel) Update() error {
	if dm.schema.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := dm.vectorDB.CheckObjectExist(dm.schema.ProjectID, dm.schema.VectorID); err != nil || !ok {
		slog.ErrorContext(dm.ctx, "dm.vectorDB.CheckObjectExist", "err", err)
		return errors.New("vector object not exist")
	}

	embedding, err := dm.createEmbeddings()
	if err != nil {
		return err
	}

	properties := &apiContentProperty{
		DefinitionModelID: vector.T_INT(dm.schema.ID),
		UpdatedAt:         vector.T_TEXT(dm.schema.UpdatedAt.Format("2006-01-02 15:04:05")),
	}

	data := &vector.ObjectData{
		Properties: properties.ToMapInterface(),
		Vector:     embedding,
	}

	return dm.vectorDB.UpdateObject(dm.schema.ProjectID, dm.schema.VectorID, data)
}

func (dm *definitionModel) Delete() error {
	if dm.schema.VectorID == "" {
		return errors.New("vector id is empty")
	}

	if ok, err := dm.vectorDB.CheckObjectExist(dm.schema.ProjectID, dm.schema.VectorID); err != nil || !ok {
		slog.ErrorContext(dm.ctx, "dm.vectorDB.CheckObjectExist", "err", err)
		return errors.New("vector object not exist")
	}

	return dm.vectorDB.DeleteObject(dm.schema.ProjectID, dm.schema.VectorID)
}

func (dm *definitionModel) createEmbeddings() ([]float32, error) {
	specSchema, err := dm.schema.ToSpec()
	if err != nil {
		return nil, err
	}

	textList := make([]string, 0)
	textList = append(textList, dm.schema.Name)
	if dm.schema.Description != "" {
		textList = append(textList, dm.schema.Description)
	}

	if specSchema.Schema.DeepRef() {
		if allModels, err := definition.GetDefinitionSchemasWithSpec(dm.ctx, dm.schema.ProjectID); err != nil {
			slog.ErrorContext(dm.ctx, "definition.GetDefinitionSchemasWithSpec", "err", err)
			return nil, err
		} else {
			specSchema.DeepDeref(allModels)
			textList = append(textList, SchemaToTextList("root", specSchema.Schema)...)
		}
	}

	text := strings.Join(textList, "\n")
	return dm.embeddingModel.CreateEmbeddings(text)
}
