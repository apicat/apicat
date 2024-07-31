package vector

type VectorApi interface {
	Check() error
	CreateCollection(name string, properties Properties) error
	CheckCollectionExist(name string) (bool, error)
	CreateObject(collectionName string, object *ObjectData) (string, error)
	CheckObjectExist(collectionName string, id string) (bool, error)
	UpdateObject(collectionName string, id string, object *ObjectData) error
	DeleteObject(collectionName string, id string) error
}
