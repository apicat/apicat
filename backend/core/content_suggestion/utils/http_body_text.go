package utils

import (
	"github.com/apicat/apicat/v2/backend/module/spec"
)

func HTTPBodyToTextList(b spec.HTTPBody) (result []string) {
	if b == nil {
		return nil
	}

	skip := map[string]bool{
		"none":                     true,
		"raw":                      true,
		"application/octet-stream": true,
		"text/plain":               true,
		"text/html":                true,
	}

	for k, v := range b {
		if _, ok := skip[k]; ok {
			continue
		}
		result = append(result, convert("root", v.Schema, 0)...)
	}

	return
}
