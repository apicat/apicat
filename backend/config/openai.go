package config

type OpenAIFile struct {
	Source   string `yaml:"source" env:"APICAT_OPENAI_SOURCE"`
	Key      string `yaml:"key" env:"APICAT_OPENAI_KEY"`
	Endpoint string `yaml:"endpoint" env:"APICAT_OPENAI_ENDPOINT"`
}

type OpenAI struct {
	Source   ConfigItem `env:"APICAT_OPENAI_SOURCE"`
	Key      ConfigItem `env:"APICAT_OPENAI_KEY"`
	Endpoint ConfigItem `env:"APICAT_OPENAI_ENDPOINT"`
}
