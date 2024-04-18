package ai

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/route/middleware/jwt"

	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/llm"
	llmcommon "github.com/apicat/apicat/v2/backend/module/llm/common"
	"github.com/gin-gonic/gin"
)

type TSGenListOption struct {
	APISummary string
	TestCases  []string
	Prompt     string
}

type TSGenDetailOption struct {
	APISummary      string
	TestCaseTitle   string
	TestCaseContent string
	Prompt          string
}

type TestCase struct {
	Purpose     string `xml:"purpose"`
	Type        string `xml:"type"`
	Description string `xml:"description"`
	Steps       string `xml:"steps"`
	Input       string `xml:"input"`
	Output      string `xml:"output"`
}

func TestCaseListGenerate(ctx *gin.Context, apiSummary string, testCases []string, prompt string) ([]string, error) {
	var tpl *tpl
	if len(testCases) == 0 {
		tpl = NewTpl("testcase_generate.tmpl", jwt.GetUser(ctx).Language, apiSummary)
	} else {
		tpl = NewTpl("testcase_more_generate.tmpl", jwt.GetUser(ctx).Language, TSGenListOption{
			APISummary: apiSummary,
			TestCases:  testCases,
			Prompt:     prompt,
		})
	}

	messages, err := tpl.Prompt()
	if err != nil {
		return nil, err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
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

	result = strings.TrimSuffix(result, "```")
	list := make([]string, 0)
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		slog.ErrorContext(ctx, "json.Unmarshal", "result", result)
		return nil, err
	}
	return list, nil
}

func TestCaseGenerate(language, projectID, apiSummary string, collectionID uint, testCaseTitleList []string) {
	c, err := cache.NewCache(config.Get().Cache.ToCfg())
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
		result, err := TestCaseDetailGenerate(language, apiSummary, testCaseTitle)
		if err != nil {
			slog.Error("TestCaseDetailGenerate", "err", err)
			continue
		}

		content := fmt.Sprintf(
			"`%s`<br>\n>%s\n### %s\n%s\n### %s\n%s\n### %s\n%s",
			result.Type,
			result.Description,
			i18n.NewTran("testCase.Steps").TranslateIn(language),
			result.Steps,
			i18n.NewTran("testCase.Input").TranslateIn(language),
			result.Input,
			i18n.NewTran("testCase.Output").TranslateIn(language),
			result.Output,
		)

		tc := &collection.TestCase{
			ProjectID:    projectID,
			CollectionID: collectionID,
			Title:        result.Purpose,
			Content:      content,
		}
		if err := tc.CreateWithoutCtx(); err != nil {
			slog.Error("tc.CreateWithoutCtx", "err", err)
		}
	}
}

func TestCaseGeneratingStatus(projectID string, collectionID uint) bool {
	c, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		return false
	}

	cacheKey := GetGenerationKey(projectID, collectionID)
	_, ok, _ := c.Get(cacheKey)
	return ok
}

func GetGenerationKey(projectID string, collectionID uint) string {
	return fmt.Sprintf("ai_test_case_generate_%s_%d", projectID, collectionID)
}

func TestCaseDetailGenerate(language, apiSummary, testCaseTitle string) (*TestCase, error) {
	tpl := NewTpl("testcase_detail_generate.tmpl", langMap[language], TSGenDetailOption{
		APISummary:    apiSummary,
		TestCaseTitle: testCaseTitle,
	})
	messages, err := tpl.Prompt()
	if err != nil {
		return nil, err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
	}

	result, err := a.ChatCompletionRequest(&llmcommon.ChatCompletionRequest{
		Temperature: 0.3,
		MaxTokens:   2000,
		Messages:    messages,
	})
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("test case detail generate failed")
	}
	result = strings.ReplaceAll(result, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")

	var ts *TestCase
	if err := xml.Unmarshal([]byte(result), &ts); err != nil {
		slog.Error("xml.Unmarshal", "result", result)
		return nil, err
	}

	v := reflect.ValueOf(ts)
	value := v.Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := value.Type().Field(i).Name
		if field.Kind() == reflect.String && field.String() == "" {
			return nil, fmt.Errorf("missing element: %s, original dats: %s", fieldName, result)
		}
	}
	ts.removeIndent()

	return ts, nil
}

func TestCaseDetailRegenerate(testCase *collection.TestCase, language, apiSummary, prompt string) (*TestCase, error) {
	tpl := NewTpl("testcase_detail_regenerate.tmpl", langMap[language], TSGenDetailOption{
		APISummary:      apiSummary,
		TestCaseTitle:   testCase.Title,
		TestCaseContent: testCase.Content,
		Prompt:          prompt,
	})
	messages, err := tpl.Prompt()
	if err != nil {
		return nil, err
	}

	a, err := llm.NewLLM(config.Get().LLM.ToMapInterface())
	if err != nil {
		return nil, err
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
		return nil, errors.New("test case detail regenerate failed")
	}
	result = strings.ReplaceAll(result, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")

	var ts *TestCase
	if err := xml.Unmarshal([]byte(result), &ts); err != nil {
		slog.Error("xml.Unmarshal", "result", result)
		return nil, err
	}

	v := reflect.ValueOf(ts)
	value := v.Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := value.Type().Field(i).Name
		if field.Kind() == reflect.String && field.String() == "" {
			return nil, fmt.Errorf("missing element: %s, original dats: %s", fieldName, result)
		}
	}
	ts.removeIndent()
	return ts, nil
}

func (t *TestCase) removeIndent() {
	t.Purpose = removeIndent(t.Purpose)
	t.Type = removeIndent(t.Type)
	t.Description = removeIndent(t.Description)
	t.Steps = removeIndent(t.Steps)
	t.Input = removeIndent(t.Input)
	t.Output = removeIndent(t.Output)
}

func removeIndent(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return s
	}

	numSpaces := 0
	newLines := make([]string, 0)
	for _, line := range lines {
		if len(newLines) == 0 {
			trimmed := strings.TrimLeft(line, " ")
			if trimmed == "" {
				continue
			}
			numSpaces = len(line) - len(trimmed)
		}

		if len(line) < numSpaces {
			continue
		}
		newLines = append(newLines, line[numSpaces:])
	}
	return strings.Join(newLines, "\n")
}
