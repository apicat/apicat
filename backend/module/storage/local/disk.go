package local

import (
	"fmt"
	"os"
	"strings"
)

type disk struct {
	path string
	url  string
}

var contentTypeMap = map[string]string{
	"image/gif":     "gif",
	"image/jpeg":    "jpg",
	"image/jpg":     "jpg",
	"image/png":     "png",
	"image/svg+xml": "svg",
	"image/webp":    "webp",
}

func NewDisk(cfg map[string]interface{}) (*disk, error) {
	for _, v := range []string{"Path", "Url"} {
		if _, ok := cfg[v]; !ok {
			return nil, fmt.Errorf("sendcloud config %s is required", v)
		}
	}

	return &disk{
		path: cfg["Path"].(string),
		url:  strings.TrimRight(cfg["Url"].(string), "/"),
	}, nil
}

func (d *disk) Check() error {
	return mkdir(d.path)
}

func (d *disk) PutObject(key string, data []byte, contentType string) (string, error) {
	if strings.Index(key, "/") > 0 {
		dir := strings.Split(key, "/")
		dirPath := d.path + "/" + strings.Join(dir[:len(dir)-1], "/")
		if err := mkdir(dirPath); err != nil {
			return "", err
		}
	}

	if _, ok := contentTypeMap[contentType]; !ok {
		return "", fmt.Errorf("not support content type: %s", contentType)
	}

	filePath := fmt.Sprintf("%s/%s", d.path, key)
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", d.url, key), nil
}

func mkdir(path string) error {
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path, 0755); err != nil {
				return err
			}
		}
	} else {
		if !info.IsDir() {
			return os.ErrNotExist
		}
	}
	return nil
}
