package reference

import (
	"context"
	"regexp"
	"strconv"

	"github.com/apicat/apicat/v2/backend/model/definition"
)

func ParseRefSchemas(text string) []uint {
	// 定义正则表达式
	re := regexp.MustCompile(`"\$ref":"#/definitions/schemas/(\d+)"`)

	// 在字符串中查找匹配项 matches: [["$ref":"#/definitions/schemas/2050" 2050] ["$ref":"#/definitions/schemas/2051" 2051]]
	matches := re.FindAllStringSubmatch(text, -1)

	// 遍历匹配项
	list := make([]uint, 0)
	for _, match := range matches {
		if len(match) >= 2 {
			// 第一个匹配项是整个匹配，从第二个匹配项开始是捕获组
			refID, err := strconv.Atoi(match[1])
			if err == nil {
				list = append(list, uint(refID))
			}
		}
	}
	return list
}

// func ParseRefSchemasFromCollection(ctx context.Context, c *collection.Collection) ([]uint, error) {
// 	specC, err := c.ToSpec()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var list []uint

// 	// specC.
// }

// func ParseRefResponsesFromCollection(ctx context.Context, c *collection.Collection) ([]uint, error) {

// }

func ParseRefSchemasFromResponse(ctx context.Context, r *definition.DefinitionResponse) ([]uint, error) {
	specR, err := r.ToSpec()
	if err != nil {
		return nil, err
	}

	refSchemaIDs := specR.RefIDs()

	var list []uint
	for _, v := range refSchemaIDs {
		list = append(list, uint(v))
	}

	return list, nil
}

func ParseRefSchemasFromSchema(ctx context.Context, s *definition.DefinitionSchema) ([]uint, error) {
	specS, err := s.ToSpec()
	if err != nil {
		return nil, err
	}

	refSchemaIDs := specS.RefIDs()

	var list []uint
	for _, v := range refSchemaIDs {
		list = append(list, uint(v))
	}

	return list, nil
}

func ParseRefResponses(text string) []uint {
	// 定义正则表达式
	re := regexp.MustCompile(`"\$ref":"#/definitions/responses/(\d+)"`)

	// 在字符串中查找匹配项 matches: [["$ref":"#/definitions/responses/2050" 2050] ["$ref":"#/definitions/responses/2051" 2051]]
	matches := re.FindAllStringSubmatch(text, -1)

	// 遍历匹配项
	list := make([]uint, 0)
	for _, match := range matches {
		if len(match) >= 2 {
			// 第一个匹配项是整个匹配，从第二个匹配项开始是捕获组
			refID, err := strconv.Atoi(match[1])
			if err == nil {
				list = append(list, uint(refID))
			}
		}
	}
	return list
}
