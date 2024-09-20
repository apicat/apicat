package response

type CollectionSuggestion struct {
	RequestID string `json:"requestID" binding:"required,gt=1"`
	Content   string `json:"content" binding:"required"`
}

type ModelSuggestion struct {
	RequestID string `json:"requestID" binding:"required,gt=1"`
	Schema    string `json:"schema" binding:"required"`
}
