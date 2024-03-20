package sysconfig

import (
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/model/sysconfig"
)

func cfgFormat(cfg *sysconfig.Sysconfig) map[string]interface{} {
	var configMap map[string]interface{}
	json.Unmarshal([]byte(cfg.Config), &configMap)
	return configMap
}
