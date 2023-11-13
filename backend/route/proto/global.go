package proto

type GlobalParameterDetails struct {
	ID       uint            `json:"id" binding:"required"`
	In       string          `json:"in" binding:"required,oneof=header query path cookie"`
	Name     string          `json:"name" binding:"required,lte=255"`
	Required bool            `json:"required"`
	Schema   ParameterSchema `json:"schema" binding:"required"`
}

type ParameterSchema struct {
	Type        string `json:"type" binding:"required,oneof=string number integer array"`
	Default     string `json:"default" binding:"omitempty,lte=255"`
	Example     string `json:"example" binding:"omitempty,lte=255"`
	Description string `json:"description" binding:"omitempty,lte=255"`
}

type GlobalParametersData struct {
	In       string          `json:"in" binding:"required,oneof=header query path cookie"`
	Name     string          `json:"name" binding:"required,lte=255"`
	Required bool            `json:"required"`
	Schema   ParameterSchema `json:"schema" binding:"required"`
}
