package api

import (
	"fmt"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type SetDBConfigData struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	DBName   string `json:"dbname" binding:"required"`
}

type DBConfigItemData struct {
	Value string `json:"value" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=value env"`
}

type DBConfigData struct {
	Host     DBConfigItemData `json:"host" binding:"required"`
	Port     DBConfigItemData `json:"port" binding:"required"`
	User     DBConfigItemData `json:"user" binding:"required"`
	Password DBConfigItemData `json:"password" binding:"required"`
	DBName   DBConfigItemData `json:"dbname" binding:"required"`
}

func dataStructProcess(field *config.ConfigItem) DBConfigItemData {
	if field.DataSource == "env" {
		return DBConfigItemData{
			Value: field.EnvName,
			Type:  field.DataSource,
		}
	} else {
		dataSource := "value"
		if field.DataSource != "" {
			dataSource = field.DataSource
		}
		return DBConfigItemData{
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
				EnvName:    fmt.Sprintf("${%s}", value),
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
	fmt.Printf("sysCfg: %+v\n", sysCfg)

	ctx.HTML(http.StatusOK, "db-config.tmpl", gin.H{
		"db_config": DBConfigData{
			Host:     dataStructProcess(&sysCfg.DB.Host),
			Port:     dataStructProcess(&sysCfg.DB.Port),
			User:     dataStructProcess(&sysCfg.DB.User),
			Password: dataStructProcess(&sysCfg.DB.Password),
			DBName:   dataStructProcess(&sysCfg.DB.Dbname),
		},
	})
}

func SetDBConfig(ctx *gin.Context) {
	data := DBConfigData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	sysCfg := config.GetSysConfig()
	fmt.Printf("sysCfg: %+v\n", sysCfg)

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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ENV.VarReadFailed"}),
		})
		return
	}

	sysCfg.DB.Host = hostField
	sysCfg.DB.Port = portField
	sysCfg.DB.User = userField
	sysCfg.DB.Password = passwordField
	sysCfg.DB.Dbname = dbNameField
	config.SetSysConfig(&sysCfg)

	models.Init()

	connStatus, err := models.DBConnStatus()
	var tm string
	if connStatus != 1 {
		switch connStatus {
		case 2:
			tm = translator.Trasnlate(ctx, &translator.TT{ID: "DB.ConnectFailed"})
		case 3:
			tm = translator.Trasnlate(ctx, &translator.TT{ID: "DB.NotFound"})
		default:
			tm = translator.Trasnlate(ctx, &translator.TT{ID: "DB.ConnectFailed"})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf(tm, err.Error()),
		})
		return
	}

	if err := config.SaveConfig(&sysCfg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Config.SaveFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
