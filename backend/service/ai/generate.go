package ai

import (
	"apicat-cloud/backend/config"
	"apicat-cloud/backend/model/collection"
	"apicat-cloud/backend/model/definition"
	"apicat-cloud/backend/module/cache"
	"apicat-cloud/backend/module/llm"
	llmcommon "apicat-cloud/backend/module/llm/common"
	"apicat-cloud/backend/module/spec"
	"apicat-cloud/backend/route/middleware/jwt"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"apicat-cloud/backend/module/spec/jsonschema"
	"apicat-cloud/backend/module/spec/plugin/openapi"

	"github.com/gin-gonic/gin"
)

var langMap = map[string]string{
	"zh":    "Chinese",
	"zh-CN": "Chinese",
	"en":    "English",
	"en-US": "English",
}

func APIGenerate(ctx *gin.Context, prompt string) (*collection.Collection, error) {
	content := createContent("api_generate.tmpl", contentData{
		Lang: langMap[jwt.GetUser(ctx).Language],
	})

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
	}

	systemPrompt, err := content.String()
	if err != nil {
		return nil, err
	}

	messages := []llmcommon.ChatCompletionMessage{
		{Role: a.ChatMessageRoleSystem(), Content: systemPrompt},
		{Role: a.ChatMessageRoleUser(), Content: prompt},
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   3000,
		Messages:    messages,
	})
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("API generate failed")
	}

	pattern := "(```yaml|```)([\\s\\S]*?)```"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(result, 1)
	if len(matches) == 0 || len(matches[0]) < 3 {
		return nil, fmt.Errorf("API generate failed, result: %s", result)
	}

	apiSpec, err := openapi.Decode([]byte(matches[0][2]))
	if err != nil {
		return nil, fmt.Errorf("openapi.Decode failed: %s content:\n%s", err.Error(), matches[0][2])
	}

	if len(apiSpec.Collections) == 0 {
		return nil, fmt.Errorf("Collection not found, original openapi content:\n%s", matches[0][2])
	}

	if err := apiSpec.Collections[0].DerefSchema(apiSpec.Definitions.Schemas...); err != nil {
		return nil, fmt.Errorf("DerefSchema failed: %s", err.Error())
	}
	if err := apiSpec.Collections[0].DerefResponse(apiSpec.Definitions.Responses...); err != nil {
		return nil, fmt.Errorf("DerefResponse failed: %s", err.Error())
	}

	r, err := json.Marshal(apiSpec.Collections[0].Content)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %s", err.Error())
	}

	return &collection.Collection{
		Title:   apiSpec.Collections[0].Title,
		Type:    string(apiSpec.Collections[0].Type),
		Content: string(r),
	}, nil
}

func APISummarize(ctx *gin.Context, collection *spec.Collection) (string, error) {
	if collection.Content == nil {
		return "", errors.New("Collection content is nil")
	}

	summary := make([]string, 0)
	if collection.Title != "" {
		summary = append(summary, fmt.Sprintf("This is an API about \"%s\"", collection.Title))
	}

	var (
		err      error
		url      *spec.HTTPURLNode
		request  *spec.HTTPRequestNode
		response *spec.HTTPResponsesNode
	)
	for _, node := range collection.Content {
		switch node.NodeType() {
		case spec.NAME_HTTP_URL:
			url, err = node.ToHTTPURLNode()
		case spec.NAME_HTTP_REQUEST:
			request, err = node.ToHTTPRequestNode()
		case spec.NAME_HTTP_RESPONSES:
			response, err = node.ToHTTPResponsesNode()
		}
		if err != nil {
			return "", err
		}
	}
	if url == nil || request == nil || response == nil {
		return "", errors.New("Incomplete API")
	}

	requestBaseInfo, err := apiRequestSummarize(url, request)
	if err != nil {
		return "", err
	}
	summary = append(summary, requestBaseInfo...)
	summary = append(summary, apiBodySummarize(&request.Content, true)...)
	summary = append(summary, apiResponseSummarize(response)...)

	return strings.Join(summary, "\n"), nil
}

func SchemaGenerate(ctx *gin.Context, prompt string) (*definition.DefinitionSchema, error) {
	content := createContent("schema_generate.tmpl", contentData{
		Lang: langMap[jwt.GetUser(ctx).Language],
	})

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
	}

	systemPrompt, err := content.String()
	if err != nil {
		return nil, err
	}

	messages := []llmcommon.ChatCompletionMessage{
		{Role: a.ChatMessageRoleSystem(), Content: systemPrompt},
		{Role: a.ChatMessageRoleUser(), Content: prompt},
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   3000,
		Messages:    messages,
	})
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("Schema generate failed")
	}

	pattern := "(```json|```)([\\s\\S]*?)```"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(result, 1)
	if len(matches) == 0 || len(matches[0]) < 3 {
		return nil, fmt.Errorf("Schema generate failed, result: %s", result)
	}

	jsonschema := &jsonschema.Schema{}
	if err := json.Unmarshal([]byte(matches[0][2]), jsonschema); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %s", err.Error())
	}

	return &definition.DefinitionSchema{
		Name:        jsonschema.Title,
		Description: jsonschema.Description,
		Type:        "schema",
		Schema:      strings.TrimSpace(matches[0][2]),
	}, nil
}

type TSGenListOption struct {
	APISummary string
	TestCases  []string
	HasPrompt  bool
}

func TestCaseListGenerate(ctx *gin.Context, apiSummary string, testCases []string, prompt string) ([]string, error) {
	content := createContent("testcase_generate.tmpl", contentData{
		Lang: langMap[jwt.GetUser(ctx).Language],
		Data: &TSGenListOption{
			APISummary: apiSummary,
			TestCases:  testCases,
			HasPrompt:  prompt != "",
		},
	})

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
	}

	systemPrompt, err := content.String()
	if err != nil {
		return nil, err
	}

	messages := []llmcommon.ChatCompletionMessage{
		{Role: a.ChatMessageRoleSystem(), Content: systemPrompt},
	}
	if prompt != "" {
		messages = append(messages, llmcommon.ChatCompletionMessage{Role: a.ChatMessageRoleUser(), Content: prompt})
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   4000,
		Messages:    messages,
	})
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("empty content")
	}

	list := make([]string, 0)
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		slog.ErrorContext(ctx, "json.Unmarshal", "result", result)
		return nil, err
	}

	return list, nil
}

func GetGenerationKey(projectID string, collectionID uint) string {
	return fmt.Sprintf("ai_test_case_generate_%s_%d", projectID, collectionID)
}

func TestCaseGeneratingStatus(projectID string, collectionID uint) bool {
	c, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		return false
	}

	cacheKey := GetGenerationKey(projectID, collectionID)
	_, ok, _ := c.Get(cacheKey)
	if ok {
		return true
	}
	return false
}

func TestCaseGenerate(language, projectID, apiSummary string, collectionID uint, testCaseTitleList []string) {
	c, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		return
	}

	cacheKey := GetGenerationKey(projectID, collectionID)
	_, ok, _ := c.Get(cacheKey)
	if ok {
		return
	}

	for _, testCaseTitle := range testCaseTitleList {
		c.Set(cacheKey, "generating", time.Second*30)
		content, err := TestCaseDetailGenerate(language, apiSummary, testCaseTitle)
		if err != nil {
			slog.Error("TestCaseDetailGenerate", "err", err)
			continue
		}

		tc := &collection.TestCase{
			ProjectID:    projectID,
			CollectionID: collectionID,
			Title:        testCaseTitle,
			Content:      content,
		}
		if err := tc.CreateWithoutCtx(); err != nil {
			slog.Error("tc.CreateWithoutCtx", "err", err)
		}
	}
}

type TSGenDetailOption struct {
	APISummary      string
	TestCaseTitle   string
	TestCaseContent string
	Prompt          string
}

func TestCaseDetailGenerate(language, apiSummary, testCaseTitle string) (string, error) {
	content := createContent("testcase_detail_generate.tmpl", contentData{
		Lang: langMap[language],
		Data: &TSGenDetailOption{
			APISummary:    apiSummary,
			TestCaseTitle: testCaseTitle,
		},
	})

	systemPrompt, err := content.String()
	if err != nil {
		return "", err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return "", err
	}

	messages := []llmcommon.ChatCompletionMessage{
		{Role: a.ChatMessageRoleSystem(), Content: systemPrompt},
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   2000,
		Messages:    messages,
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("Test case detail generate failed")
	}

	return result, nil
}

func TestCaseDetailRegenerate(testCase *collection.TestCase, language, apiSummary, prompt string) (string, error) {
	content := createContent("testcase_detail_regenerate.tmpl", contentData{
		Lang: langMap[language],
		Data: &TSGenDetailOption{
			APISummary: apiSummary,
		},
	})
	systemPrompt, err := content.String()
	if err != nil {
		return "", err
	}

	content = createContent("testcase_detail_regenerate_user.tmpl", contentData{
		Lang: langMap[language],
		Data: &TSGenDetailOption{
			TestCaseTitle:   testCase.Title,
			TestCaseContent: testCase.Content,
			Prompt:          prompt,
		},
	})
	userPrompt, err := content.String()
	if err != nil {
		return "", err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return "", err
	}

	messages := []llmcommon.ChatCompletionMessage{
		{Role: a.ChatMessageRoleSystem(), Content: systemPrompt},
		{Role: a.ChatMessageRoleUser(), Content: userPrompt},
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   4000,
		Messages:    messages,
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("test case detail regenerate failed")
	}

	pattern := `<NEWTESTCASE>([\s\S]*?)</NEWTESTCASE>`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(result, 1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", fmt.Errorf("API generate failed, result: %s", result)
	}

	return strings.TrimSpace(matches[0][1]), nil
}

func apiRequestSummarize(url *spec.HTTPURLNode, request *spec.HTTPRequestNode) ([]string, error) {
	summary := make([]string, 0)

	if url.Path == "" || url.Method == "" {
		return summary, errors.New("Incomplete url")
	}
	summary = append(summary, fmt.Sprintf("The request path of the API is \"%s\", and the request method is \"%s\".", url.Path, url.Method))

	if request.Parameters.Path != nil && len(request.Parameters.Path) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Path request parameters, namely: ", len(request.Parameters.Path)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Path)...)
	}
	if request.Parameters.Header != nil && len(request.Parameters.Header) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Header request parameters, namely: ", len(request.Parameters.Header)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Header)...)
	}
	if request.Parameters.Cookie != nil && len(request.Parameters.Cookie) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Cookie request parameters, namely: ", len(request.Parameters.Cookie)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Cookie)...)
	}
	if request.Parameters.Query != nil && len(request.Parameters.Query) > 0 {
		summary = append(summary, fmt.Sprintf("There are %d Query request parameters, namely: ", len(request.Parameters.Query)))
		summary = append(summary, apiParameterSummarize(&request.Parameters.Query)...)
	}

	return summary, nil
}

func apiParameterSummarize(parameters *spec.ParameterList) []string {
	summary := make([]string, 0)
	for i, p := range *parameters {
		if p.Name == "" {
			continue
		}

		content := fmt.Sprintf("%d. %s: %s", i+1, p.Name, p.Schema.Type.Value()[0])
		if p.Required {
			content += ", required"
		} else {
			content += ", optional"
		}
		if p.Schema.Description != "" {
			content += fmt.Sprintf(", %s", p.Schema.Description)
		}
		summary = append(summary, content)
	}

	return summary
}

func apiResponseSummarize(responses *spec.HTTPResponsesNode) []string {
	summary := make([]string, 0)

	summary = append(summary, fmt.Sprintf("%d response situations are possible:", len(responses.List)))
	for i, response := range responses.List {
		desc := fmt.Sprintf("%d. Response status code %d: %s", i+1, response.Code, response.Name)
		if response.Description != "" {
			summary = append(summary, fmt.Sprintf("%s, %s", desc, response.Description))
		} else {
			summary = append(summary, desc)
		}

		if response.Header != nil && len(response.Header) > 0 {
			summary = append(summary, fmt.Sprintf("There are %d Header response parameters, namely: ", len(response.Header)))
			summary = append(summary, apiParameterSummarize(&response.Header)...)
		}

		if response.Content != nil && len(response.Content) > 0 {
			summary = append(summary, apiBodySummarize(&response.Content, false)...)
		}
	}

	return summary
}

func apiBodySummarize(body *spec.HTTPBody, isRequest bool) []string {
	summary := make([]string, 0)
	for contentType, schema := range *body {
		if contentType == "none" {
			continue
		}

		var contentTypeDesc string
		if isRequest {
			contentTypeDesc = fmt.Sprintf("The \"Content-Type\" in the request body is: %s", contentType)
		} else {
			contentTypeDesc = fmt.Sprintf("The \"Content-Type\" in the response body is: %s", contentType)
		}

		content := make(map[string]any)
		apiJsonSchemaSummarize("root", true, schema.Schema, content)

		switch contentType {
		case "multipart/form-data", "application/x-www-form-urlencoded":
			if len(content) > 0 {
				summary = append(summary, fmt.Sprintf("%s, containing the following %d parameters: ", contentTypeDesc, len(content)))
				i := 1
				for k, v := range content["root"].(map[string]any) {
					summary = append(summary, fmt.Sprintf("%d. %s: %s", i, k, v.(string)))
					i++
				}
			}
		case "application/json", "application/xm":
			if len(content) > 0 {
				summary = append(summary, fmt.Sprintf("%s, and the structure is as follows: ", contentTypeDesc))
				js, err := json.Marshal(content["root"])
				if err != nil {
					return summary
				}
				summary = append(summary, string(js))
			}
		case "raw", "application/octet-stream":
			if schema.Schema.Example != nil {
				if example, ok := schema.Schema.Example.(string); ok {
					summary = append(summary, fmt.Sprintf("%s. For example: %s", contentTypeDesc, example))
				}
			}
		}
	}

	return summary
}

// @title apiJsonSchemaSummarize
// @description 总结 JSON Schema 要表达的含义
// @param parameterName string 参数的名称
// @param required bool 参数是否必填
// @param schema *jsonschema 参数的 JSON Schema 描述
// @param result map[string]any 带结构的描述结果
// @return nil
func apiJsonSchemaSummarize(parameterName string, required bool, schema *jsonschema.Schema, result map[string]any) {
	if schema.Type == nil {
		return
	}

	switch schema.Type.Value()[0] {
	case "string", "integer", "number", "boolean":
		requiredStr := "optional"
		if required {
			requiredStr = "required"
		}
		desc := ""
		if schema.Description != "" {
			desc = fmt.Sprintf(", %s", schema.Description)
		}
		result[parameterName] = fmt.Sprintf("%s, %s%s", schema.Type.Value()[0], requiredStr, desc)
	case "object":
		if schema.Properties != nil && len(schema.Properties) > 0 {
			requiredParams := make(map[string]bool)
			if schema.Required != nil && len(schema.Required) > 0 {
				for _, v := range schema.Required {
					requiredParams[v] = true
				}
			}

			children := make(map[string]any)
			for k, v := range schema.Properties {
				_, r := requiredParams[k]
				apiJsonSchemaSummarize(k, r, v, children)
			}
			if len(children) > 0 {
				result[parameterName] = children
			}
		}
	case "array":
		if schema.Items != nil && !schema.Items.IsBool() {
			children := make(map[string]any)
			apiJsonSchemaSummarize("items", true, schema.Items.Value(), children)
			if len(children) > 0 {
				for _, v := range children {
					result[parameterName] = v
					break
				}
			}
		}
	}
}
