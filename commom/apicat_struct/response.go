package apicat_struct

type ResponseObject struct {
	Name        string      `json:"name" binding:"required,lte=255"`
	Code        int         `json:"code" binding:"required"`
	Description string      `json:"description" binding:"required,lte=255"`
	Header      []*Header   `json:"header" binding:"omitempty,dive"`
	Content     *BodyObject `json:"content" binding:"required"`
}

type Header struct {
	Name        string       `json:"name" binding:"required,lte=255"`
	Description string       `json:"description" binding:"lte=255"`
	Required    bool         `json:"required"`
	Schema      SchemaObject `json:"schema"`
}
