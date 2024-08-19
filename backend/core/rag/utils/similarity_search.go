package utils

import (
	"context"
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/module/model"
	"github.com/apicat/apicat/v2/backend/module/vector"
)

type SimilaritySearch struct {
	collectionName   string
	text             string
	fields           []string
	additionalFields []string
	offset           int
	limit            int
	distance         float32
	certainty        float32
	embeddingModel   model.Provider
	vectorDB         vector.VectorApi
}

func NewSimilaritySearch(collectionName, text string) (*SimilaritySearch, error) {
	ctx := context.Background()

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

	return &SimilaritySearch{
		collectionName: collectionName,
		text:           text,
		embeddingModel: embeddingModel,
		vectorDB:       vectorDB,
	}, nil
}

func (s *SimilaritySearch) WithFields(fields []string) *SimilaritySearch {
	s.fields = fields
	return s
}

func (s *SimilaritySearch) WithAdditionalFields(fields []string) *SimilaritySearch {
	s.additionalFields = fields
	return s
}

func (s *SimilaritySearch) WithOffset(offset int) *SimilaritySearch {
	if offset > 0 {
		s.offset = offset
	}
	return s
}

func (s *SimilaritySearch) WithLimit(limit int) *SimilaritySearch {
	if limit > 0 {
		s.limit = limit
	}
	return s
}

func (s *SimilaritySearch) WithDistance(distance float32) *SimilaritySearch {
	if distance > 0 {
		s.distance = distance
	}
	return s
}

func (s *SimilaritySearch) WithCertainty(certainty float32) *SimilaritySearch {
	if certainty > 0 {
		s.certainty = certainty
	}
	return s
}

func (s *SimilaritySearch) Do() (string, error) {
	if s.text == "" {
		return "", errors.New("text is required")
	}

	v, err := s.embeddingModel.CreateEmbeddings(s.text)
	if err != nil {
		return "", err
	}

	opt := &vector.SearchOption{
		Vector:           v,
		Fields:           s.fields,
		AdditionalFields: s.additionalFields,
		Offset:           s.offset,
		Limit:            s.limit,
		Distance:         s.distance,
		Certainty:        s.certainty,
	}

	return s.vectorDB.SimilaritySearch(s.collectionName, opt)
}
