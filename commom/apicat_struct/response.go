package apicat_struct

import "strconv"

type ResponseObjectList struct {
	List []*ResponseObject `json:"list"`
}

type ResponseObject struct {
	ID          uint       `json:"id,omitempty"`
	Name        string     `json:"name" binding:"required,lte=255"`
	Code        int        `json:"code" binding:"required"`
	Description string     `json:"description" binding:"required,lte=255"`
	Header      []*Header  `json:"header" binding:"omitempty,dive"`
	Content     BodyObject `json:"content" binding:"required"`
	Ref         string     `json:"$ref,omitempty" binding:"omitempty,lte=255"`
}

type Header struct {
	Name        string       `json:"name" binding:"required,lte=255"`
	Description string       `json:"description" binding:"lte=255"`
	Required    bool         `json:"required"`
	Schema      SchemaObject `json:"schema"`
}

func (rl *ResponseObjectList) Dereference(ro *ResponseObject) {
	for i, r := range rl.List {
		if r.Ref == "#/commons/responses/"+strconv.FormatUint(uint64(ro.ID), 10) {
			rl.List = append(rl.List[:i], rl.List[i+1:]...)
		}
	}
	rl.List = append(rl.List, ro)
}
