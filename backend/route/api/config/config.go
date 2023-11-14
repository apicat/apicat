package config

import (
	"fmt"
	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model"
	"github.com/apicat/apicat/backend/route/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func dataStructProcess(field *config.ConfigItem) proto.DBConfigItemData {
	if field.DataSource == "env" {
		return proto.DBConfigItemData{
			Value: field.EnvName,
			Type:  field.DataSource,
		}
	} else {
		dataSource := "value"
		if field.DataSource != "" {
			dataSource = field.DataSource
		}
		return proto.DBConfigItemData{
			Value: field.Value,
			Type:  dataSource,
		}
	}
}

func generateConfigItem(value, dataSource string) (config.ConfigItem, error) {
	var field config.ConfigItem

	if dataSource == "env" {
		if ev, exist := os.LookupEnv(value); !exist {
			return field, fmt.Errorf("env var %s read failed", value)
		} else {
			field = config.ConfigItem{
				Value:      ev,
				DataSource: "env",
				EnvName:    value,
			}
		}
	} else {
		field = config.ConfigItem{
			Value:      value,
			DataSource: "value",
		}
	}

	return field, nil
}

func GetDBConfig(ctx *gin.Context) {
	sysCfg := config.GetSysConfig()

	ctx.HTML(http.StatusOK, "db-config.tmpl", gin.H{
		"db_config": proto.DBConfigData{
			Host:     dataStructProcess(&sysCfg.DB.Host),
			Port:     dataStructProcess(&sysCfg.DB.Port),
			User:     dataStructProcess(&sysCfg.DB.User),
			Password: dataStructProcess(&sysCfg.DB.Password),
			DBName:   dataStructProcess(&sysCfg.DB.Dbname),
		},
	})
}

func SetDBConfig(ctx *gin.Context) {
	data := proto.DBConfigData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	sysCfg := config.GetSysConfig()

	ok := true
	hostField, err := generateConfigItem(data.Host.Value, data.Host.Type)
	if err != nil {
		ok = false
	}
	portField, err := generateConfigItem(data.Port.Value, data.Port.Type)
	if err != nil {
		ok = false
	}
	userField, err := generateConfigItem(data.User.Value, data.User.Type)
	if err != nil {
		ok = false
	}
	passwordField, err := generateConfigItem(data.Password.Value, data.Password.Type)
	if err != nil {
		ok = false
	}
	dbNameField, err := generateConfigItem(data.DBName.Value, data.DBName.Type)
	if err != nil {
		ok = false
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ENV.VarReadFailed"}),
		})
		return
	}

	sysCfg.DB.Host = hostField
	sysCfg.DB.Port = portField
	sysCfg.DB.User = userField
	sysCfg.DB.Password = passwordField
	sysCfg.DB.Dbname = dbNameField
	config.SetSysConfig(&sysCfg)

	model.Init()

	connStatus, err := model.DBConnStatus()
	var tm string
	if connStatus != 1 {
		switch connStatus {
		case 2:
			tm = i18n.Trasnlate(ctx, &i18n.TT{ID: "DB.ConnectFailed"})
		case 3:
			tm = i18n.Trasnlate(ctx, &i18n.TT{ID: "DB.NotFound"})
		default:
			tm = i18n.Trasnlate(ctx, &i18n.TT{ID: "DB.ConnectFailed"})
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf(tm, err.Error()),
		})
		return
	}

	if err := config.SaveConfig(&sysCfg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Config.SaveFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
