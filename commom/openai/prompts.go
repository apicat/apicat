package openai

var createApiPrompt = `"""
Please use OpenAPI3.0.0 to describe the %s API, do not create schema.
"""`

var createSchemaPrompt = `"""
Please use JSON Schema to describe the %s.
The title field is required and must be in English.
"""`

var listApiBySchemaPrompt = `"""
Below I will provide a JSON Schema content named %s, please provide all related API descriptions, request paths and request methods according to this schema, and return the results in the form of an JSON array object.
For example:
[
	{"description": "create user", "method": "POST", "path": "/users"}
]
JSON Schema:
%s
"""`
