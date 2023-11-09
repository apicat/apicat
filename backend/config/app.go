package config

type AppFile struct {
	Name string `yaml:"name" env:"APICAT_APP_NAME"`
	Host string `yaml:"host" env:"APICAT_APP_HOST"`
	Port string `yaml:"port" env:"APICAT_APP_PORT"`
}

type App struct {
	Name ConfigItem `env:"APICAT_APP_NAME"`
	Host ConfigItem `env:"APICAT_APP_HOST"`
	Port ConfigItem `env:"APICAT_APP_PORT"`
}
