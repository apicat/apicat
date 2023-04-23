package apicat_struct

import "github.com/apicat/apicat/commom/spec"

type ResponseObjectList struct {
	List []*ResponseObject `json:"list"`
}

type ResponseObject struct {
	Name        string                 `json:"name"`
	Code        int                    `json:"code"`
	Description string                 `json:"description"`
	Header      []*Header              `json:"header,omitempty"`
	Content     map[string]spec.Schema `json:"content,omitempty"`
	Ref         string                 `json:"$ref,omitempty"`
}

type Header struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Example     string      `json:"example"`
	Default     string      `json:"default"`
	Required    bool        `json:"required"`
	Schema      spec.Schema `json:"schema"`
}
