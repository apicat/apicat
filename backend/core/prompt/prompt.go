package prompt

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/apicat/apicat/v2/backend/module/model"
)

const (
	SYSTEM_PROMPT   = "___system:"
	USER_PROMPT     = "_____user:"
	ASISTENT_PROMPT = "_asistent:"
	// SYSTEM_PROMPT, USER_PROMPT, ASISTENT_PROMPT string lengths must be the same, and ROLE_PROMPT_LEN is the value of the length
	ROLE_PROMPT_LEN = 10
	PROMPT_END      = "----------------------------------------end"
)

//go:embed templates
var tplfs embed.FS

var templateFs = template.Must(template.New("prompt").ParseFS(tplfs, "templates/*.tmpl"))

type prompt struct {
	Lang            string
	SysmtemRole     string
	UserRole        string
	AssistantRole   string
	SystemPrompt    string
	UserPrompt      string
	AssistantPrompt string
	PromptEnd       string
	Context         any
}

type tpl struct {
	tplname string
	tpldata prompt
}

func NewPrompt(templateName string, lang string, context any) *tpl {
	p := prompt{
		Lang:            language(lang),
		SysmtemRole:     "system",
		UserRole:        "user",
		AssistantRole:   "assistant",
		SystemPrompt:    SYSTEM_PROMPT,
		UserPrompt:      USER_PROMPT,
		AssistantPrompt: ASISTENT_PROMPT,
		PromptEnd:       PROMPT_END,
		Context:         context,
	}

	return &tpl{
		tplname: templateName,
		tpldata: p,
	}
}

func (p *prompt) SetRole(role map[string]string) {
	if _, ok := role["system"]; ok {
		p.SysmtemRole = role["system"]
	}
	if _, ok := role["user"]; ok {
		p.UserRole = role["user"]
	}
	if _, ok := role["assistant"]; ok {
		p.AssistantRole = role["assistant"]
	}
}

func (t *tpl) Prompt() (model.ChatCompletionMessages, error) {
	content, err := t.split()
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("prompt content is empty")
	}

	var message model.ChatCompletionMessages
	for _, p := range content {
		p = strings.TrimSpace(p)
		if len(p) < ROLE_PROMPT_LEN {
			continue
		}

		if m, err := t.build(p); err != nil {
			return nil, err
		} else {
			message = append(message, m)
		}
	}
	return message, nil
}

func (t *tpl) string() (string, error) {
	var buf bytes.Buffer
	if err := templateFs.ExecuteTemplate(&buf, t.tplname, t.tpldata); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t *tpl) split() ([]string, error) {
	if content, err := t.string(); err != nil {
		return nil, err
	} else {
		return strings.Split(content, PROMPT_END), nil
	}
}

func (t *tpl) build(content string) (model.ChatCompletionMessage, error) {
	role := content[:ROLE_PROMPT_LEN]
	msg := content[ROLE_PROMPT_LEN:]
	roleName := ""

	switch role {
	case SYSTEM_PROMPT:
		roleName = t.tpldata.SysmtemRole
	case USER_PROMPT:
		roleName = t.tpldata.UserRole
	case ASISTENT_PROMPT:
		roleName = t.tpldata.AssistantRole
	default:
		return model.ChatCompletionMessage{}, fmt.Errorf("wrong role: %s", role)
	}

	return model.ChatCompletionMessage{
		Role:    roleName,
		Content: msg,
	}, nil
}
