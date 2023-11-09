package config

type LogFile struct {
	Path  string `yaml:"path" env:"APICAT_LOG_PATH"`
	Level string `yaml:"level" env:"APICAT_LOG_LEVEL"`
}

type Log struct {
	Path  ConfigItem `env:"APICAT_LOG_PATH"`
	Level ConfigItem `env:"APICAT_LOG_LEVEL"`
}
