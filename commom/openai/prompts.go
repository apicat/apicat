package openai

var listApiBySchemaPrompt = `"""
Below I will provide a JSON Schema content named %s, please provide all related API descriptions, request paths and request methods according to this schema, and return the results in the form of an JSON array object.
For example:
[
	{"description": "create user", "method": "POST", "path": "/users"}
]
JSON Schema:
%s
"""`
