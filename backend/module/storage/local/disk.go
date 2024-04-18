package local

import (
	"fmt"
	"os"
	"strings"
)

type Disk struct {
	Path string
	Url  string
}

var contentTypeMap = map[string]string{
	"image/gif":     "gif",
	"image/jpeg":    "jpg",
	"image/jpg":     "jpg",
	"image/png":     "png",
	"image/svg+xml": "svg",
	"image/webp":    "webp",
}

func NewDisk(cfg Disk) *Disk {
	return &cfg
}

func (d *Disk) Check() error {
	return mkdir(d.Path)
}

func (d *Disk) PutObject(key string, data []byte, contentType string) (string, error) {
	if strings.Index(key, "/") > 0 {
		dir := strings.Split(key, "/")
		dirPath := d.Path + "/" + strings.Join(dir[:len(dir)-1], "/")
		if err := mkdir(dirPath); err != nil {
			return "", err
		}
	}

	if _, ok := contentTypeMap[contentType]; !ok {
		return "", fmt.Errorf("not support content type: %s", contentType)
	}

	filePath := fmt.Sprintf("%s/%s", d.Path, key)
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", d.Url, key), nil
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
