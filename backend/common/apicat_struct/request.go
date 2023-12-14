package apicat_struct

type RequestObject struct {
	GlobalExcepts GlobalExceptsObject `json:"globalExcepts"`
	Parameters    ParametersObject    `json:"parameters"`
	Content       any                 `json:"content"`
}
