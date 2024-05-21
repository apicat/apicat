package base

type JsonSchemaOption struct {
	JsonSchema string `json:"jsonschema" binding:"required,gt=1"`
}
