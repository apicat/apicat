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

func GetDBConfig(ctx *gin.Context) {
	sysCfg := config.GetSysConfig()

	dataStructureProcess := func(field *config.ConfigItem) DBConfigItemData {
		if field.DataSource == "env" {
			return DBConfigItemData{
				Value: field.EnvName,
				Type:  field.DataSource,
			}
		} else {
			return DBConfigItemData{
				Value: field.Value,
				Type:  field.DataSource,
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"db_config": DBConfigData{
			Host:     dataStructureProcess(&sysCfg.DB.Host),
			Port:     dataStructureProcess(&sysCfg.DB.Port),
			User:     dataStructureProcess(&sysCfg.DB.User),
			Password: dataStructureProcess(&sysCfg.DB.Password),
			DBName:   dataStructureProcess(&sysCfg.DB.Dbname),
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
	generateConfigItem := func(field *config.ConfigItem, value, dataSource string) {
		if dataSource == "env" {
			if ev, exist := os.LookupEnv(value); exist {
				field = &config.ConfigItem{
					Value:      ev,
					DataSource: "env",
					EnvName:    fmt.Sprintf("${%s}", value),
				}
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": translator.Trasnlate(ctx, &translator.TT{ID: "ENV.VarReadFailed"}),
				})
				return
			}

		} else {
			field = &config.ConfigItem{
				Value:      value,
				DataSource: "value",
			}
		}
	}

	generateConfigItem(&sysCfg.DB.Host, data.Host.Value, data.Host.Type)
	generateConfigItem(&sysCfg.DB.Port, data.Port.Value, data.Port.Type)
	generateConfigItem(&sysCfg.DB.User, data.User.Value, data.User.Type)
	generateConfigItem(&sysCfg.DB.Password, data.Password.Value, data.Password.Type)
	generateConfigItem(&sysCfg.DB.Dbname, data.DBName.Value, data.DBName.Type)

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
