package model

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var defaultDB *gorm.DB
var autoMigrateTables = []any{}

func DB(ctx context.Context) *gorm.DB {
	return defaultDB.WithContext(ctx)
}

func DBWithoutCtx() *gorm.DB {
	return defaultDB
}

// RegMigrate 注册需要自动同步表结构的对象
func RegMigrate(v ...any) {
	autoMigrateTables = append(autoMigrateTables, v...)
}

func Init() error {
	cfg := config.Get().Database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Database,
	)
	slog.Info("init database", "host", cfg.Host, "database", cfg.Database)
	dbLogger := &tracelogger{}
	if cfg.Debug {
		dbLogger.lvl = logger.Info
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return err
	}
	rawDB, err := db.DB()
	if err != nil {
		return err
	}
	if cfg.MaxIdleConns > 0 {
		rawDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		rawDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxIdleTime > 0 {
		rawDB.SetConnMaxLifetime(cfg.ConnMaxIdleTime)
	}
	if len(autoMigrateTables) > 0 {
		if err := db.AutoMigrate(autoMigrateTables...); err != nil {
			return err
		}
	}
	defaultDB = db
	return nil
}

type TimeModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
