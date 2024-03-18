package common

type Storage interface {
	Check() error
	PutObject(key string, data []byte, contentType string) (string, error)
}
