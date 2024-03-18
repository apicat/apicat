package sysconfig

import (
	"apicat-cloud/backend/model/sysconfig"
	"encoding/json"
)

func cfgFormat(cfg *sysconfig.Sysconfig) map[string]interface{} {
	var configMap map[string]interface{}
	json.Unmarshal([]byte(cfg.Config), &configMap)
	return configMap
}
