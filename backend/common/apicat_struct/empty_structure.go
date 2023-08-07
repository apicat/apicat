package apicat_struct

func TypeEmptyStructure() map[string]string {
	TypeEmptyStructure := map[string]string{
		"string":   `{"type":"string"}`,
		"integer":  `{"type":"string"}`,
		"number":   `{"type":"string"}`,
		"boolean":  `{"type":"boolean"}`,
		"array":    `{"type":"array","items":{"type":"string"}}`,
		"object":   `{"type":"object","properties":{},"example":"","required":[],"x-apicat-orders":[]}`,
		"response": `{"name":"Response Name","code":200,"description":"","content":{"application/json":{"schema":{"type":"object","properties":{},"required":[],"x-apicat-orders":[],"example":""}}}}`,
	}

	return TypeEmptyStructure
}
