package models

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/apicat/apicat/backend/config"
	"golang.org/x/exp/slog"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// Conn 是一个指向数据库连接的指针
	Conn *gorm.DB

	// connStatus 数据库连接状态，非成功状态跳转到配置页，
	connStatus     uint = connOK
	connOK         uint = 1
	connFail       uint = 2
	connDBNotFound uint = 3
	connErr        error
)

func Init() {
	var err error

	dbLogger := &tracelogger{}
	if strings.ToUpper(config.GetSysConfig().Log.Level.Value) == "DEBUG" {
		dbLogger.lvl = logger.Info
	}

	switch config.GetSysConfig().DB.Driver.Value {
	case "sqlite":
		if _, err = os.Stat(config.GetSysConfig().DB.Path.Value); os.IsNotExist(err) {
			os.Mkdir(config.GetSysConfig().DB.Path.Value, os.ModePerm)
		}
		Conn, err = gorm.Open(sqlite.Open(config.GetSysConfig().DB.Path.Value+config.GetSysConfig().DB.Dbname.Value+".db"), &gorm.Config{Logger: dbLogger})
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.GetSysConfig().DB.User.Value,
			config.GetSysConfig().DB.Password.Value,
			config.GetSysConfig().DB.Host.Value,
			config.GetSysConfig().DB.Port.Value,
			config.GetSysConfig().DB.Dbname.Value,
		)
		Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
	default:
		panic("There is no setting for the database driver type.")
	}

	if err != nil {
		// 检查是数据库无法连接还是数据库不存在，数据库不存在则创建
		// 创建数据库连接字符串
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/",
			config.GetSysConfig().DB.User.Value,
			config.GetSysConfig().DB.Password.Value,
			config.GetSysConfig().DB.Host.Value,
			config.GetSysConfig().DB.Port.Value,
		)

		// 尝试连接到MySQL数据库
		Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			connStatus = connFail
			connErr = err
			slog.Error("failed to connect to database", slog.String("err", err.Error()))
			return
		}

		if err := pingConn(); err != nil {
			connStatus = connFail
			connErr = err
			return
		}

		if err := connectDB(config.GetSysConfig().DB.Dbname.Value); err != nil {
			connStatus = connDBNotFound
			connErr = err
			return
		}
	}

	connStatus = connOK
	initTable()
}

// DBConnStatus status 返回数据库连接状态，1-成功，2-失败，3-数据库不存在 err 连接遇到错误时返回错误信息
func DBConnStatus() (status uint, err error) {
	return connStatus, connErr
}

func pingConn() error {
	sqlDB, err := Conn.DB()
	if err != nil {
		slog.Error("failed to obtain database connection object", slog.String("err", err.Error()))
		return err
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		slog.Error("failed to ping database", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// 连接数据库，不存在则创建。调用时Conn应已与mysql建立连接，但未指定数据库
func connectDB(dbName string) error {
	var count int64
	if result := Conn.Raw("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Count(&count); result.Error != nil {
		slog.Error("failed to check database", slog.String("err", result.Error.Error()))
		return result.Error
	}

	if count == 0 {
		// 当数据库连接正常但数据库不存在时，帮用户创建数据库
		if err := Conn.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;").Error; err != nil {
			slog.Error("failed to create database", slog.String("err", err.Error()))
			return err
		}
	}

	// 连接到新数据库
	dbLogger := &tracelogger{}
	if strings.ToUpper(config.GetSysConfig().Log.Level.Value) == "DEBUG" {
		dbLogger.lvl = logger.Info
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetSysConfig().DB.User.Value,
		config.GetSysConfig().DB.Password.Value,
		config.GetSysConfig().DB.Host.Value,
		config.GetSysConfig().DB.Port.Value,
		config.GetSysConfig().DB.Dbname.Value,
	)
	var err error
	Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		panic("连接到新数据库时出错：" + err.Error())
	}

	return nil
}

func initTable() {
	if err := Conn.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
		&Projects{},
		&Collections{},
		&CollectionHistories{},
		&Servers{},
		&Tags{},
		&TagToCollections{},
		&GlobalParameters{},
		&DefinitionSchemas{},
		&DefinitionResponses{},
		&DefinitionParameters{},
		&Users{},
		&ProjectMembers{},
		&ShareTmpTokens{},
		&DefinitionSchemaHistories{},
		&Iterations{},
		&IterationApis{},
		&ProjectGroups{},
	); err != nil {
		panic(err.Error())
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
