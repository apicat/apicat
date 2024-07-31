package vector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
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

func NewWeaviate(cfg WeaviateOpt) (*weaviatedb, error) {
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
	if b, err := json.MarshalIndent(resp.Object, "", "  "); err == nil {
		fmt.Println(string(b))
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
			Properties: w.covertProperties(r.Properties),
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
		Properties: w.covertProperties(resp[0].Properties),
		Vector:     resp[0].Vector,
	}, nil
}

func (w *weaviatedb) CheckObjectExist(collectionName string, id string) (bool, error) {
	return w.client.Data().Checker().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).Do(w.ctx)
}

func (w *weaviatedb) UpdateObject(collectionName string, id string, object *ObjectData) error {
	return w.client.Data().Updater().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).WithProperties(object.Properties).WithVector(object.Vector).Do(w.ctx)
}

func (w *weaviatedb) DeleteObject(collectionName string, id string) error {
	return w.client.Data().Deleter().WithClassName(w.getValidCollectionName(collectionName)).WithID(id).Do(w.ctx)
}

func (w *weaviatedb) covertProperties(properties models.PropertySchema) map[string]interface{} {
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
