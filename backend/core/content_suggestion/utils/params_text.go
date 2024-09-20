package utils

import (
	"github.com/apicat/apicat/v2/backend/module/spec"
)

func ParamsToTextList(k string, p spec.ParameterList) (result []string) {
	if p == nil {
		return nil
	}

	result = append(result, genText(k, "", 0))
	for _, v := range p {
		result = append(result, convert(v.Name, v.Schema, 1)...)
	}
	return
}
