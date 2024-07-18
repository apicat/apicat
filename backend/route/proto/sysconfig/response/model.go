package response

type ModelConfigDetail struct {
	Driver string                 `json:"driver" binding:"required,gt=1"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

type ModelConfigList []*ModelConfigDetail

type DefaultModelDetail struct {
	Driver   string `json:"driver" binding:"required,gt=1"`
	Model    string `json:"model" binding:"required,gt=1"`
	Selected bool   `json:"selected" binding:"required"`
}

type DefaultModelList []*DefaultModelDetail

type DefaultModelMap map[string]DefaultModelList
