package base

type IterationIDOption struct {
	IterationID string `uri:"iterationID" json:"iterationID" query:"iterationID" binding:"required,len=24"`
}

type IterationData struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
