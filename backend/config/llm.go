package config

import "encoding/json"

type LLM struct {
	Driver      string       `yaml:"Driver"`
	OpenAI      *OpenAI      `yaml:"OpenAI"`
	AzureOpenAI *AzureOpenAI `yaml:"AzureOpenAI"`
}

type OpenAI struct {
	ApiKey         string `yaml:"ApiKey"`
	OrganizationID string `yaml:"OrganizationID"`
	ApiBase        string `yaml:"ApiBase"`
	LLMName        string `yaml:"LLMName"`
	EmbeddingName  string `yaml:"EmbeddingName"`
}

type AzureOpenAI struct {
	ApiKey        string `yaml:"ApiKey"`
	Endpoint      string `yaml:"Endpoint"`
	LLMName       string `yaml:"LLMName"`
	EmbeddingName string `yaml:"EmbeddingName"`
}

func SetLLM(c *LLM) {
	globalConf.LLM = c
}

func (l *LLM) ToMapInterface() map[string]interface{} {
	var (
		res      map[string]interface{}
		jsonByte []byte
	)
	jsonByte, _ = json.Marshal(l)
	_ = json.Unmarshal(jsonByte, &res)
	return res
}
