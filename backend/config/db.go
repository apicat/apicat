package config

type DBFile struct {
	Driver   string `yaml:"driver" env:"APICAT_DB_DRIVER"`
	Path     string `yaml:"path" env:"APICAT_DB_PATH"`
	Host     string `yaml:"host" env:"APICAT_DB_HOST"`
	Port     string `yaml:"port" env:"APICAT_DB_PORT"`
	User     string `yaml:"user" env:"APICAT_DB_USER"`
	Password string `yaml:"password" env:"APICAT_DB_PASSWORD"`
	Dbname   string `yaml:"dbname" env:"APICAT_DB_NAME"`
}

type DB struct {
	Driver   ConfigItem `env:"APICAT_DB_DRIVER"`
	Path     ConfigItem `env:"APICAT_DB_PATH"`
	Host     ConfigItem `env:"APICAT_DB_HOST"`
	Port     ConfigItem `env:"APICAT_DB_PORT"`
	User     ConfigItem `env:"APICAT_DB_USER"`
	Password ConfigItem `env:"APICAT_DB_PASSWORD"`
	Dbname   ConfigItem `env:"APICAT_DB_NAME"`
}
