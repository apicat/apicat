package vector

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateOpt struct {
	Host   string
	ApiKey string
}

type weaviatedb struct {
	ctx    context.Context
	client *weaviate.Client
}

func newWeaviate(cfg WeaviateOpt) (*weaviatedb, error) {
	c := weaviate.Config{
		Host:   cfg.Host,
		Scheme: "http",
	}

	if cfg.ApiKey != "" {
		c.AuthConfig = auth.ApiKey{Value: cfg.ApiKey}
	}

	if client, err := weaviate.NewClient(c); err == nil {
		return &weaviatedb{
			ctx:    context.Background(),
			client: client,
		}, nil
	} else {
		return nil, err
	}
}

func (w *weaviatedb) Check() error {
	if _, err := w.client.Misc().MetaGetter().Do(w.ctx); err != nil {
		return err
	}
	return nil
}

func (w *weaviatedb) CreateCollection(name string, properties Properties) error {
	props := make([]*models.Property, 0)
	for _, p := range properties {
		props = append(props, &models.Property{
			Name:        p.Name,
			DataType:    []string{p.DataType.TypeString()},
			Description: p.Description,
		})
	}

	c := &models.Class{
		Class:      w.getValidCollectionName(name),
		Properties: props,
	}
	return w.client.Schema().ClassCreator().WithClass(c).Do(w.ctx)
}

func (w *weaviatedb) CheckCollectionExist(name string) (bool, error) {
	if _, err := w.client.Schema().ClassGetter().WithClassName(w.getValidCollectionName(name)).Do(w.ctx); err != nil {
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (w *weaviatedb) DeleteCollection(name string) error {
	if err := w.client.Schema().ClassDeleter().WithClassName(w.getValidCollectionName(name)).Do(w.ctx); err != nil {
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode != http.StatusBadRequest {
			return err
		}
	}
	return nil
}

func (w *weaviatedb) CreateObject(collectionName string, object *ObjectData) (string, error) {
	resp, err := w.client.Data().Creator().WithClassName(w.getValidCollectionName(collectionName)).WithProperties(object.Properties).WithVector(object.Vector).Do(w.ctx)
	if err != nil {
		return "", err
	}
	return resp.Object.ID.String(), nil
}

func (w *weaviatedb) BatchCreateObject(collectionName string, data ObjectDataList) (ObjectDataList, error) {
	dataObjs := []models.PropertySchema{}
	for _, d := range data {
		dataObjs = append(dataObjs, d.Properties)
	}

	batcher := w.client.Batch().ObjectsBatcher()
	for i, dataObj := range dataObjs {
		batcher.WithObjects(&models.Object{
			Class:      w.getValidCollectionName(collectionName),
			Properties: dataObj,
			Vector:     data[i].Vector,
		})
	}
	resp, err := batcher.Do(w.ctx)
	if err != nil {
		return nil, err
	}

	result := make(ObjectDataList, 0)
	for _, r := range resp {
		result = append(result, &ObjectData{
			ID:         r.ID.String(),
			Properties: w.convertProperties(r.Properties),
			Vector:     r.Vector,
		})
	}
	return result, nil
}

func (w *weaviatedb) GetObjectByID(collectionName string, id string) (*ObjectData, error) {
	resp, err := w.client.Data().ObjectsGetter().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).Do(w.ctx)
	if err != nil {
		return nil, err
	}
	return &ObjectData{
		ID:         resp[0].ID.String(),
		Properties: w.convertProperties(resp[0].Properties),
		Vector:     resp[0].Vector,
	}, nil
}

func (w *weaviatedb) CheckObjectExist(collectionName string, id string) (bool, error) {
	if _, err := w.client.Data().Checker().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).Do(w.ctx); err != nil {
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (w *weaviatedb) UpdateObject(collectionName string, id string, object *ObjectData) error {
	return w.client.Data().Updater().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).WithProperties(object.Properties).WithVector(object.Vector).Do(w.ctx)
}

func (w *weaviatedb) DeleteObject(collectionName string, id string) error {
	return w.client.Data().Deleter().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).Do(w.ctx)
}

func (w *weaviatedb) SimilaritySearch(collectionName string, opt *SearchOption) (string, error) {
	if len(opt.Vector) == 0 {
		return "", errors.New("vector is required")
	}

	builder := w.client.GraphQL().Get().WithClassName(w.getValidCollectionName(collectionName))
	if len(opt.Fields) > 0 || len(opt.AdditionalFields) > 0 {
		if fields := w.getSearchFields(opt.Fields, opt.AdditionalFields); fields != nil {
			builder.WithFields(fields...)
		}
	}

	nearVector := w.client.GraphQL().NearVectorArgBuilder().WithVector(opt.Vector)
	if opt.Distance > 0 {
		nearVector.WithDistance(opt.Distance)
	}
	if opt.Certainty > 0 {
		nearVector.WithCertainty(opt.Certainty)
	}

	builder.WithNearVector(nearVector)

	if opt.Offset > 0 {
		builder.WithOffset(opt.Offset)
	}
	if opt.Limit > 0 {
		builder.WithLimit(opt.Limit)
	}
	if opt.WhereCondition != nil {
		if where, err := w.buildWhere(opt.WhereCondition); err == nil {
			builder.WithWhere(where)
		} else {
			return "", err
		}
	}

	result, err := builder.Do(w.ctx)
	if err != nil {
		return "", err
	}
	if len(result.Errors) > 0 {
		return "", errors.New(result.Errors[0].Message)
	}

	if _, ok := result.Data["Get"]; ok {
		if data, ok := result.Data["Get"].(map[string]interface{}); ok {
			if _, ok := data[w.getValidCollectionName(collectionName)]; ok {
				if b, err := json.Marshal(data[w.getValidCollectionName(collectionName)]); err == nil {
					return string(b), nil
				}
			}
		}
	}
	return "", nil
}

func (w *weaviatedb) convertProperties(properties models.PropertySchema) map[string]interface{} {
	propertiesMap := make(map[string]interface{})
	for k, v := range properties.(map[string]interface{}) {
		propertiesMap[k] = v
	}
	return propertiesMap
}

func (w *weaviatedb) getValidCollectionName(name string) string {
	if len(name) == 0 {
		return name
	}
	return "C" + name
}

func (w *weaviatedb) getSearchFields(fields, additionalFields []string) []graphql.Field {
	if len(fields) == 0 && len(additionalFields) == 0 {
		return nil
	}

	searchFields := make([]graphql.Field, 0)
	for _, f := range fields {
		searchFields = append(searchFields, graphql.Field{
			Name: f,
		})
	}

	if len(additionalFields) > 0 {
		additional := make([]graphql.Field, 0)
		for _, f := range additionalFields {
			additional = append(additional, graphql.Field{
				Name: f,
			})
		}
		searchFields = append(searchFields, graphql.Field{
			Name:   "_additional",
			Fields: additional,
		})
	}
	return searchFields
}

func (w *weaviatedb) buildWhere(condition []*WhereCondition) (*filters.WhereBuilder, error) {
	if len(condition) == 0 {
		return nil, errors.New("condition is required")
	}

	if len(condition) == 1 {
		return w.parseCondition(condition[0])
	}

	whereBuilder := filters.Where().WithOperator(filters.And)
	operands := make([]*filters.WhereBuilder, 0)
	for _, c := range condition {
		if operand, err := w.parseCondition(c); err == nil {
			operands = append(operands, operand)
		} else {
			return nil, err
		}
	}
	whereBuilder.WithOperands(operands)
	return whereBuilder, nil
}

func (w *weaviatedb) parseCondition(condition *WhereCondition) (*filters.WhereBuilder, error) {
	if condition == nil {
		return nil, errors.New("condition is required")
	}
	if condition.PropertyName == "" {
		return nil, errors.New("property name is required")
	}

	whereBuilder := filters.Where().WithPath([]string{condition.PropertyName})

	switch condition.Operator {
	case "=", "==", "eq":
		whereBuilder.WithOperator(filters.Equal)
	case "!=", "<>", "neq":
		whereBuilder.WithOperator(filters.NotEqual)
	case ">", "gt":
		whereBuilder.WithOperator(filters.GreaterThan)
	case ">=", "gte":
		whereBuilder.WithOperator(filters.GreaterThanEqual)
	case "<", "lt":
		whereBuilder.WithOperator(filters.LessThan)
	case "<=", "lte":
		whereBuilder.WithOperator(filters.LessThanEqual)
	default:
		return nil, errors.New("invalid operator")
	}

	switch condition.Value.TypeString() {
	case "text":
		whereBuilder.WithValueString(reflect.ValueOf(condition.Value).String())
	case "int":
		whereBuilder.WithValueInt(reflect.ValueOf(condition.Value).Int())
	case "boolean":
		whereBuilder.WithValueBoolean(reflect.ValueOf(condition.Value).Bool())
	case "number":
		whereBuilder.WithValueNumber(reflect.ValueOf(condition.Value).Float())
	}
	return whereBuilder, nil
}
