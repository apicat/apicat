package models

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/apicat/apicat/config"
	"golang.org/x/exp/slog"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB

func Init() {
	var err error

	dbLogger := &tracelogger{}
	if strings.ToUpper(config.SysConfig.Log.Level) == "DEBUG" {
		dbLogger.lvl = logger.Info
	}

	switch config.SysConfig.DB.Driver {
	case "sqlite":
		if _, err = os.Stat(config.SysConfig.DB.Path); os.IsNotExist(err) {
			os.Mkdir(config.SysConfig.DB.Path, os.ModePerm)
		}
		Conn, err = gorm.Open(sqlite.Open(config.SysConfig.DB.Path+config.SysConfig.DB.Dbname+".db"), &gorm.Config{Logger: dbLogger})
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
		Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
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
		&CommonResponses{},
	); err != nil {
		panic("Failed to create database table.")
	}
}

// tracelogger 集成traceid
type tracelogger struct {
	lvl logger.LogLevel
}

func (l *tracelogger) LogMode(lvl logger.LogLevel) logger.Interface {
	newlog := *l
	newlog.lvl = lvl
	return &newlog
}
func (l *tracelogger) Info(ctx context.Context, s string, v ...interface{})  {}
func (l *tracelogger) Warn(ctx context.Context, s string, v ...interface{})  {}
func (l *tracelogger) Error(ctx context.Context, s string, v ...interface{}) {}
func (l *tracelogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.lvl <= logger.Silent {
		return
	}

	sql, rows := fc()
	dur := time.Since(begin)
	logattr := []any{
		slog.String("sql", sql),
		slog.Int64("rows", rows),
		slog.Duration("latency", dur),
	}
	switch {
	case err != nil && l.lvl >= logger.Error:
		slog.ErrorCtx(ctx, "gorm.trace err:"+err.Error(), logattr...)
	case dur >= time.Millisecond*500 && l.lvl >= logger.Warn:
		slog.WarnCtx(ctx, "gorm.trace slow sql", logattr...)
	case l.lvl == logger.Info:
		slog.InfoCtx(ctx, "gorm.trace", logattr...)
	}
}
