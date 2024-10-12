package config

import (
	"errors"
	"os"

	"github.com/apicat/apicat/v2/backend/module/model"
)

type Model struct {
	LLMDriver       string
	EmbeddingDriver string
	OpenAI          *OpenAI
	AzureOpenAI     *AzureOpenAI
	Baichuan        *Baichuan
	Moonshot        *Moonshot
	DeepSeek        *DeepSeek
	VolcanoEngine   *VolcanoEngine
}

type OpenAI struct {
	ApiKey         string `json:"apiKey"`
	OrganizationID string `json:"organizationID"`
	ApiBase        string `json:"apiBase"`
	LLM            string `json:"llm"`
	Embedding      string `json:"embedding"`
}

type AzureOpenAI struct {
	ApiKey              string `json:"apiKey"`
	Endpoint            string `json:"endpoint"`
	LLM                 string `json:"llm"`
	LLMDeployName       string `json:"llmDeployName"`
	Embedding           string `json:"embedding"`
	EmbeddingDeployName string `json:"embeddingDeployName"`
}

type Baichuan struct {
	ApiKey    string `json:"apiKey"`
	LLM       string `json:"llm"`
	Embedding string `json:"embedding"`
}

type Moonshot struct {
	ApiKey    string `json:"apiKey"`
	LLM       string `json:"llm"`
	Embedding string `json:"embedding"`
}

type DeepSeek struct {
	ApiKey    string `json:"apiKey"`
	LLM       string `json:"llm"`
	Embedding string `json:"embedding"`
}

type VolcanoEngine struct {
	ApiKey              string `json:"apiKey"`
	Region              string `json:"region"`
	BaseUrl             string `json:"baseUrl"`
	LLM                 string `json:"llm"`
	LLMEndpointID       string `json:"llmEndpointID"`
	Embedding           string `json:"embedding"`
	EmbeddingEndpointID string `json:"embeddingEndpointID"`
}

func LoadModelConfig() {
	globalConf.Model = &Model{}

	if v, exists := os.LookupEnv("LLM_DRIVER"); exists {
		switch v {
		case model.OPENAI:
			globalConf.Model.LLMDriver = model.OPENAI
			loadOpenAIConfig()
		case model.AZURE_OPENAI:
			globalConf.Model.LLMDriver = model.AZURE_OPENAI
			loadAzureOpenAIConfig()
		case model.BAICHUAN:
			globalConf.Model.LLMDriver = model.BAICHUAN
			loadBaichuanConfig()
		case model.MOONSHOT:
			globalConf.Model.LLMDriver = model.MOONSHOT
			loadMoonshotConfig()
		case model.DEEPSEEK:
			globalConf.Model.LLMDriver = model.DEEPSEEK
			loadDeepSeekConfig()
		case model.VOLCANOENGINE:
			globalConf.Model.LLMDriver = model.VOLCANOENGINE
			loadVolcanoEngineConfig()
		}
	}

	if v, exists := os.LookupEnv("EMBEDDING_DRIVER"); exists {
		switch v {
		case model.OPENAI:
			globalConf.Model.EmbeddingDriver = model.OPENAI
			loadOpenAIConfig()
		case model.AZURE_OPENAI:
			globalConf.Model.EmbeddingDriver = model.AZURE_OPENAI
			loadAzureOpenAIConfig()
		case model.BAICHUAN:
			globalConf.Model.EmbeddingDriver = model.BAICHUAN
			loadBaichuanConfig()
		case model.MOONSHOT:
			globalConf.Model.EmbeddingDriver = model.MOONSHOT
			loadMoonshotConfig()
		case model.DEEPSEEK:
			globalConf.Model.EmbeddingDriver = model.DEEPSEEK
			loadDeepSeekConfig()
		case model.VOLCANOENGINE:
			globalConf.Model.EmbeddingDriver = model.VOLCANOENGINE
			loadVolcanoEngineConfig()
		}
	}
}

func loadOpenAIConfig() {
	globalConf.Model.OpenAI = &OpenAI{}
	if v, exists := os.LookupEnv("OPENAI_API_KEY"); exists {
		globalConf.Model.OpenAI.ApiKey = v
	}
	if v, exists := os.LookupEnv("OPENAI_ORGANIZATION_ID"); exists {
		globalConf.Model.OpenAI.OrganizationID = v
	}
	if v, exists := os.LookupEnv("OPENAI_API_BASE"); exists {
		globalConf.Model.OpenAI.ApiBase = v
	}
	if v, exists := os.LookupEnv("OPENAI_LLM"); exists {
		globalConf.Model.OpenAI.LLM = v
	}
	if v, exists := os.LookupEnv("OPENAI_EMBEDDING"); exists {
		globalConf.Model.OpenAI.Embedding = v
	}
}

func loadAzureOpenAIConfig() {
	globalConf.Model.AzureOpenAI = &AzureOpenAI{}
	if v, exists := os.LookupEnv("AZURE_OPENAI_API_KEY"); exists {
		globalConf.Model.AzureOpenAI.ApiKey = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_ENDPOINT"); exists {
		globalConf.Model.AzureOpenAI.Endpoint = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_LLM"); exists {
		globalConf.Model.AzureOpenAI.LLM = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_LLM_DEPLOY_NAME"); exists {
		globalConf.Model.AzureOpenAI.LLMDeployName = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_EMBEDDING"); exists {
		globalConf.Model.AzureOpenAI.Embedding = v
	}
	if v, exists := os.LookupEnv("AZURE_OPENAI_EMBEDDING_DEPLOY_NAME"); exists {
		globalConf.Model.AzureOpenAI.EmbeddingDeployName = v
	}
}

func loadBaichuanConfig() {
	globalConf.Model.Baichuan = &Baichuan{}
	if v, exists := os.LookupEnv("BAICHUAN_API_KEY"); exists {
		globalConf.Model.Baichuan.ApiKey = v
	}
	if v, exists := os.LookupEnv("BAICHUAN_LLM"); exists {
		globalConf.Model.Baichuan.LLM = v
	}
	if v, exists := os.LookupEnv("BAICHUAN_EMBEDDING"); exists {
		globalConf.Model.Baichuan.Embedding = v
	}
}

func loadMoonshotConfig() {
	globalConf.Model.Moonshot = &Moonshot{}
	if v, exists := os.LookupEnv("MOONSHOT_API_KEY"); exists {
		globalConf.Model.Moonshot.ApiKey = v
	}
	if v, exists := os.LookupEnv("MOONSHOT_LLM"); exists {
		globalConf.Model.Moonshot.LLM = v
	}
	if v, exists := os.LookupEnv("MOONSHOT_EMBEDDING"); exists {
		globalConf.Model.Moonshot.Embedding = v
	}
}

func loadDeepSeekConfig() {
	globalConf.Model.DeepSeek = &DeepSeek{}
	if v, exists := os.LookupEnv("DEEPSEEK_API_KEY"); exists {
		globalConf.Model.DeepSeek.ApiKey = v
	}
	if v, exists := os.LookupEnv("DEEPSEEK_LLM"); exists {
		globalConf.Model.DeepSeek.LLM = v
	}
	if v, exists := os.LookupEnv("DEEPSEEK_EMBEDDING"); exists {
		globalConf.Model.DeepSeek.Embedding = v
	}
}

func loadVolcanoEngineConfig() {
	globalConf.Model.VolcanoEngine = &VolcanoEngine{}
	if v, exists := os.LookupEnv("VOLCANOENGINE_API_KEY"); exists {
		globalConf.Model.VolcanoEngine.ApiKey = v
	}
	if v, exists := os.LookupEnv("VOLCANOENGINE_REGION"); exists {
		globalConf.Model.VolcanoEngine.Region = v
	}
	if v, exists := os.LookupEnv("VOLCANOENGINE_BASE_URL"); exists {
		globalConf.Model.VolcanoEngine.BaseUrl = v
	}
	if v, exists := os.LookupEnv("VOLCANOENGINE_LLM"); exists {
		globalConf.Model.VolcanoEngine.LLM = v
	}
	if v, exists := os.LookupEnv("VOLCANOENGINE_LLM_ENDPOINT_ID"); exists {
		globalConf.Model.VolcanoEngine.LLMEndpointID = v
	}
	if v, exists := os.LookupEnv("VOLCANOENGINE_EMBEDDING"); exists {
		globalConf.Model.VolcanoEngine.Embedding = v
	}
	if v, exists := os.LookupEnv("VOLCANOENGINE_EMBEDDING_ENDPOINT_ID"); exists {
		globalConf.Model.VolcanoEngine.EmbeddingEndpointID = v
	}
}

func CheckModelConfig() error {
	if globalConf.Model.LLMDriver != "" {
		switch globalConf.Model.LLMDriver {
		case model.OPENAI:
			if err := checkOpenAI("llm"); err != nil {
				return err
			}
		case model.AZURE_OPENAI:
			if err := checkAzureOpenAI("llm"); err != nil {
				return err
			}
		case model.BAICHUAN:
			if err := checkBaichuan("llm"); err != nil {
				return err
			}
		case model.MOONSHOT:
			if err := checkMoonshot("llm"); err != nil {
				return err
			}
		case model.DEEPSEEK:
			if err := checkDeepSeek("llm"); err != nil {
				return err
			}
		case model.VOLCANOENGINE:
			if err := checkVolcanoEngine("llm"); err != nil {
				return err
			}
		}
	}
	if globalConf.Model.EmbeddingDriver != "" {
		switch globalConf.Model.EmbeddingDriver {
		case model.OPENAI:
			if err := checkOpenAI("embedding"); err != nil {
				return err
			}
		case model.AZURE_OPENAI:
			if err := checkAzureOpenAI("embedding"); err != nil {
				return err
			}
		case model.BAICHUAN:
			if err := checkBaichuan("embedding"); err != nil {
				return err
			}
		case model.MOONSHOT:
			if err := checkMoonshot("embedding"); err != nil {
				return err
			}
		case model.DEEPSEEK:
			if err := checkDeepSeek("embedding"); err != nil {
				return err
			}
		case model.VOLCANOENGINE:
			if err := checkVolcanoEngine("embedding"); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkOpenAI(modelType string) error {
	if globalConf.Model.OpenAI == nil {
		return errors.New("openai config is empty")
	}
	if globalConf.Model.OpenAI.ApiKey == "" {
		return errors.New("openai api key is empty")
	}
	if modelType == "llm" && globalConf.Model.OpenAI.LLM == "" {
		return errors.New("openai llm is empty")
	}
	if modelType == "embedding" && globalConf.Model.OpenAI.Embedding == "" {
		return errors.New("openai embedding is empty")
	}
	if modelType == "llm" && !model.ModelAvailable(model.OPENAI, modelType, globalConf.Model.OpenAI.LLM) {
		return errors.New("llm model not supported")
	}
	if modelType == "embedding" && !model.ModelAvailable(model.OPENAI, modelType, globalConf.Model.OpenAI.Embedding) {
		return errors.New("embedding model not supported")
	}
	return nil
}

func checkAzureOpenAI(modelType string) error {
	if globalConf.Model.AzureOpenAI == nil {
		return errors.New("azure openai config is empty")
	}
	if globalConf.Model.AzureOpenAI.ApiKey == "" {
		return errors.New("azure openai api key is empty")
	}
	if globalConf.Model.AzureOpenAI.Endpoint == "" {
		return errors.New("azure openai endpoint is empty")
	}
	if modelType == "llm" && globalConf.Model.AzureOpenAI.LLM == "" {
		return errors.New("azure openai llm is empty")
	}
	if globalConf.Model.AzureOpenAI.LLM != "" && globalConf.Model.AzureOpenAI.LLMDeployName == "" {
		return errors.New("azure openai llm deploy name is empty")
	}
	if modelType == "embedding" && globalConf.Model.AzureOpenAI.Embedding == "" {
		return errors.New("azure openai embedding is empty")
	}
	if globalConf.Model.AzureOpenAI.Embedding != "" && globalConf.Model.AzureOpenAI.EmbeddingDeployName == "" {
		return errors.New("azure openai embedding deploy name is empty")
	}
	if modelType == "llm" && !model.ModelAvailable(model.AZURE_OPENAI, modelType, globalConf.Model.AzureOpenAI.LLM) {
		return errors.New("llm model not supported")
	}
	if modelType == "embedding" && !model.ModelAvailable(model.AZURE_OPENAI, modelType, globalConf.Model.AzureOpenAI.Embedding) {
		return errors.New("embedding model not supported")
	}
	return nil
}

func checkBaichuan(modelType string) error {
	if globalConf.Model.Baichuan == nil {
		return errors.New("baichuan config is empty")
	}
	if globalConf.Model.Baichuan.ApiKey == "" {
		return errors.New("baichuan api key is empty")
	}
	if modelType == "llm" && globalConf.Model.Baichuan.LLM == "" {
		return errors.New("baichuan llm is empty")
	}
	if modelType == "embedding" && globalConf.Model.Baichuan.Embedding == "" {
		return errors.New("baichuan embedding is empty")
	}
	if modelType == "llm" && !model.ModelAvailable(model.BAICHUAN, modelType, globalConf.Model.Baichuan.LLM) {
		return errors.New("llm model not supported")
	}
	if modelType == "embedding" && !model.ModelAvailable(model.BAICHUAN, modelType, globalConf.Model.Baichuan.Embedding) {
		return errors.New("embedding model not supported")
	}
	return nil
}

func checkMoonshot(modelType string) error {
	if globalConf.Model.Moonshot == nil {
		return errors.New("moonshot config is empty")
	}
	if globalConf.Model.Moonshot.ApiKey == "" {
		return errors.New("moonshot api key is empty")
	}
	if modelType == "llm" && globalConf.Model.Moonshot.LLM == "" {
		return errors.New("moonshot llm is empty")
	}
	if modelType == "embedding" && globalConf.Model.Moonshot.Embedding == "" {
		return errors.New("moonshot embedding is empty")
	}
	if modelType == "llm" && !model.ModelAvailable(model.MOONSHOT, modelType, globalConf.Model.Moonshot.LLM) {
		return errors.New("llm model not supported")
	}
	if modelType == "embedding" && !model.ModelAvailable(model.MOONSHOT, modelType, globalConf.Model.Moonshot.Embedding) {
		return errors.New("embedding model not supported")
	}
	return nil
}

func checkDeepSeek(modelType string) error {
	if globalConf.Model.DeepSeek == nil {
		return errors.New("deepseek config is empty")
	}
	if globalConf.Model.DeepSeek.ApiKey == "" {
		return errors.New("deepseek api key is empty")
	}
	if modelType == "llm" && globalConf.Model.DeepSeek.LLM == "" {
		return errors.New("deepseek llm is empty")
	}
	if modelType == "embedding" && globalConf.Model.DeepSeek.Embedding == "" {
		return errors.New("deepseek embedding is empty")
	}
	if modelType == "llm" && !model.ModelAvailable(model.DEEPSEEK, modelType, globalConf.Model.DeepSeek.LLM) {
		return errors.New("llm model not supported")
	}
	if modelType == "embedding" && !model.ModelAvailable(model.DEEPSEEK, modelType, globalConf.Model.DeepSeek.Embedding) {
		return errors.New("embedding model not supported")
	}
	return nil
}

func checkVolcanoEngine(modelType string) error {
	if globalConf.Model.VolcanoEngine == nil {
		return errors.New("volcano engine config is empty")
	}
	if globalConf.Model.VolcanoEngine.ApiKey == "" {
		return errors.New("volcano engine api key is empty")
	}
	if globalConf.Model.VolcanoEngine.Region == "" {
		return errors.New("volcano engine region is empty")
	}
	if globalConf.Model.VolcanoEngine.BaseUrl == "" {
		return errors.New("volcano engine base url is empty")
	}
	if modelType == "llm" && globalConf.Model.VolcanoEngine.LLM == "" {
		return errors.New("volcano engine llm is empty")
	}
	if modelType == "llm" && globalConf.Model.VolcanoEngine.LLMEndpointID == "" {
		return errors.New("volcano engine llm endpoint id is empty")
	}
	if modelType == "embedding" && globalConf.Model.VolcanoEngine.Embedding == "" {
		return errors.New("volcano engine embedding is empty")
	}
	if modelType == "embedding" && globalConf.Model.VolcanoEngine.EmbeddingEndpointID == "" {
		return errors.New("volcano engine embedding endpoint id is empty")
	}
	if modelType == "llm" && !model.ModelAvailable(model.VOLCANOENGINE, modelType, globalConf.Model.VolcanoEngine.LLM) {
		return errors.New("llm model not supported")
	}
	if modelType == "embedding" && !model.ModelAvailable(model.VOLCANOENGINE, modelType, globalConf.Model.VolcanoEngine.Embedding) {
		return errors.New("embedding model not supported")
	}
	return nil
}

func GetModel() *Model {
	return globalConf.Model
}

func SetModel(m *Model) {
	if m.LLMDriver != "" {
		globalConf.Model.LLMDriver = m.LLMDriver
		switch m.LLMDriver {
		case model.OPENAI:
			globalConf.Model.OpenAI = m.OpenAI
		case model.AZURE_OPENAI:
			globalConf.Model.AzureOpenAI = m.AzureOpenAI
		case model.BAICHUAN:
			globalConf.Model.Baichuan = m.Baichuan
		case model.MOONSHOT:
			globalConf.Model.Moonshot = m.Moonshot
		case model.DEEPSEEK:
			globalConf.Model.DeepSeek = m.DeepSeek
		case model.VOLCANOENGINE:
			globalConf.Model.VolcanoEngine = m.VolcanoEngine
		default:
			globalConf.Model.LLMDriver = ""
		}
	}
	if m.EmbeddingDriver != "" {
		globalConf.Model.EmbeddingDriver = m.EmbeddingDriver
		switch m.EmbeddingDriver {
		case model.OPENAI:
			globalConf.Model.OpenAI = m.OpenAI
		case model.AZURE_OPENAI:
			globalConf.Model.AzureOpenAI = m.AzureOpenAI
		case model.BAICHUAN:
			globalConf.Model.Baichuan = m.Baichuan
		case model.MOONSHOT:
			globalConf.Model.Moonshot = m.Moonshot
		case model.DEEPSEEK:
			globalConf.Model.DeepSeek = m.DeepSeek
		case model.VOLCANOENGINE:
			globalConf.Model.VolcanoEngine = m.VolcanoEngine
		default:
			globalConf.Model.EmbeddingDriver = ""
		}
	}
}

func (m *Model) ToModuleStruct(modelType string) model.Model {
	if m == nil {
		return model.Model{}
	}

	var driver string
	if modelType == "llm" {
		driver = m.LLMDriver
	} else if modelType == "embedding" {
		driver = m.EmbeddingDriver
	} else {
		return model.Model{}
	}

	switch driver {
	case model.OPENAI:
		return m.toOpenAICfg()
	case model.AZURE_OPENAI:
		return m.toAzureOpenAICfg()
	case model.BAICHUAN:
		return m.toBaichuanCfg()
	case model.MOONSHOT:
		return m.toMoonshotCfg()
	case model.DEEPSEEK:
		return m.toDeepSeekCfg()
	case model.VOLCANOENGINE:
		return m.toVolcanoEngineCfg()
	default:
		return model.Model{}
	}
}

func (m *Model) toOpenAICfg() model.Model {
	return model.Model{
		Driver: model.OPENAI,
		OpenAI: model.OpenAI{
			ApiKey:         m.OpenAI.ApiKey,
			OrganizationID: m.OpenAI.OrganizationID,
			ApiBase:        m.OpenAI.ApiBase,
			LLM:            m.OpenAI.LLM,
			Embedding:      m.OpenAI.Embedding,
		},
	}
}

func (m *Model) toAzureOpenAICfg() model.Model {
	return model.Model{
		Driver: model.AZURE_OPENAI,
		AzureOpenAI: model.AzureOpenAI{
			ApiKey:              m.AzureOpenAI.ApiKey,
			Endpoint:            m.AzureOpenAI.Endpoint,
			LLM:                 m.AzureOpenAI.LLM,
			LLMDeployName:       m.AzureOpenAI.LLMDeployName,
			Embedding:           m.AzureOpenAI.Embedding,
			EmbeddingDeployName: m.AzureOpenAI.EmbeddingDeployName,
		},
	}
}

func (m *Model) toBaichuanCfg() model.Model {
	return model.Model{
		Driver: model.BAICHUAN,
		Baichuan: model.Baichuan{
			ApiKey:    m.Baichuan.ApiKey,
			LLM:       m.Baichuan.LLM,
			Embedding: m.Baichuan.Embedding,
		},
	}
}

func (m *Model) toMoonshotCfg() model.Model {
	return model.Model{
		Driver: model.MOONSHOT,
		Moonshot: model.Moonshot{
			ApiKey:    m.Moonshot.ApiKey,
			LLM:       m.Moonshot.LLM,
			Embedding: m.Moonshot.Embedding,
		},
	}
}

func (m *Model) toDeepSeekCfg() model.Model {
	return model.Model{
		Driver: model.DEEPSEEK,
		DeepSeek: model.DeepSeek{
			ApiKey:    m.DeepSeek.ApiKey,
			LLM:       m.DeepSeek.LLM,
			Embedding: m.DeepSeek.Embedding,
		},
	}
}

func (m *Model) toVolcanoEngineCfg() model.Model {
	return model.Model{
		Driver: model.VOLCANOENGINE,
		VolcanoEngine: model.VolcanoEngine{
			ApiKey:              m.VolcanoEngine.ApiKey,
			Region:              m.VolcanoEngine.Region,
			BaseUrl:             m.VolcanoEngine.BaseUrl,
			LLM:                 m.VolcanoEngine.LLM,
			LLMEndpointID:       m.VolcanoEngine.LLMEndpointID,
			Embedding:           m.VolcanoEngine.Embedding,
			EmbeddingEndpointID: m.VolcanoEngine.EmbeddingEndpointID,
		},
	}
}
