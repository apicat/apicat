package apicat_struct

type SchemaObject struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

type Schema struct {
	Type        string `json:"type"`
	Example     string `json:"example"`
	Description string `json:"description"`
}
