package reference

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/global"
	"github.com/apicat/apicat/v2/backend/module/spec"
	arrutil "github.com/apicat/apicat/v2/backend/utils/array"
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

// ParseExceptParameterFromCollection 读取collection中排除的全局参数
func ParseExceptParameterFromCollection(c *collection.Collection) []uint {
	list := make([]uint, 0)

	var specContent spec.CollectionNodes
	if err := json.Unmarshal([]byte(c.Content), &specContent); err != nil {
		return list
	}

	var request *spec.CollectionHttpRequest
	for _, i := range specContent {
		switch i.NodeType() {
		case spec.NODE_HTTP_REQUEST:
			request = i.ToHttpRequest()
		}
	}

	if request == nil {
		return list
	}

	globalExceptKey := []string{
		string(global.ParameterInHeader),
		string(global.ParameterInPath),
		string(global.ParameterInHeader),
		string(global.ParameterInCookie),
	}
	for key, value := range request.Attrs.GlobalExcepts.ToMap() {
		if !arrutil.InArray(key, globalExceptKey) {
			continue
		}
		for _, v := range value {
			list = append(list, uint(v))
		}
	}

	return list
}
