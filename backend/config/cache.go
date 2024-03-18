package config

import (
	"encoding/json"
)

type Cache struct {
	Driver string `yaml:"Driver"`
	Redis  *Redis `yaml:"Redis"`
}

type Redis struct {
	Host     string `yaml:"Host"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

func GetCacheDefault() *Cache {
	return &Cache{
		Driver: "redis",
		Redis: &Redis{
			Host:     "apicat_redis:6379",
			Password: "",
			DB:       0,
		},
	}
}

func SetCache(cacheConfig *Cache) {
	globalConf.Cache = cacheConfig
}

func (c *Cache) ToMapInterface() map[string]interface{} {
	var (
		res      map[string]interface{}
		jsonByte []byte
	)
	jsonByte, _ = json.Marshal(c)
	_ = json.Unmarshal(jsonByte, &res)
	return res
}
