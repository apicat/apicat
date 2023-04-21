package models

import (
	"fmt"
	"os"

	"github.com/apicat/apicat/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Init() {
	var err error

	switch config.SysConfig.DB.Driver {
	case "sqlite":
		if _, err = os.Stat(config.SysConfig.DB.Path); os.IsNotExist(err) {
			os.Mkdir(config.SysConfig.DB.Path, os.ModePerm)
		}
		Conn, err = gorm.Open(sqlite.Open(config.SysConfig.DB.Path+config.SysConfig.DB.Dbname+".db"), &gorm.Config{})
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			config.SysConfig.DB.User,
			config.SysConfig.DB.Password,
			config.SysConfig.DB.Host,
			config.SysConfig.DB.Port,
			config.SysConfig.DB.Dbname,
			config.SysConfig.DB.Charset,
		)
		Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		panic("There is no setting for the database driver type.")
	}

	if err != nil {
		panic("Failed to connect database.")
	}

	initTable()
}

func initTable() {
	if err := Conn.AutoMigrate(
		&Projects{},
		&Collections{},
		&CollectionHistories{},
		&Definitions{},
		&Servers{},
		&Commons{},
		&Tags{},
		&TagToCollections{},
		&GlobalParameters{},
	); err != nil {
		panic("Failed to create database table.")
	}
}
