//go:generate catb --in_dir=backend/route/proto --out_dir=doc
package apicat

import (
	"fmt"
	"log"
	"net/url"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/logger"
	"github.com/apicat/apicat/v2/backend/module/mock"
	"github.com/apicat/apicat/v2/backend/module/storage"
	"github.com/apicat/apicat/v2/backend/route"
)

type App struct{}

func NewApp(conf string) *App {
	if err := config.Load(conf); err != nil {
		log.Printf("load config %s faild, use default config. err: %s", conf, err)
	}
	config.LoadFromEnv()
	return &App{}
}

func (a *App) Run() error {
	if err := model.Init(); err != nil {
		return fmt.Errorf("init %v", err)
	}
	sysconfig.Init()

	if err := runMock(); err != nil {
		return err
	}

	if err := logger.Init(config.Get().App.Debug, config.Get().Log); err != nil {
		return fmt.Errorf("init %v", err)
	}

	if err := cache.Init(config.Get().Cache.ToMapInterface()); err != nil {
		return fmt.Errorf("init %v", err)
	}

	if err := storage.Init(config.Get().Storage.ToMapInterface()); err != nil {
		return fmt.Errorf("init %v", err)
	}

	if err := route.Init(); err != nil {
		return fmt.Errorf("init %v", err)
	}
	return nil
}

func runMock() error {
	cfg := config.Get().App
	if cfg.AppUrl == "" || cfg.MockServerBind == "" {
		return fmt.Errorf("init mock err, cfg: %v", cfg)
	}

	// 尝试解析URL
	u, err := url.Parse(cfg.AppUrl)
	if err != nil {
		return fmt.Errorf("init mock err, cfg: %v", cfg)
	}

	// 检查协议是否是http或https
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("init mock err, cfg: %v", cfg)
	}

	go mock.Run(cfg.MockServerBind, mock.WithApiUrl(cfg.AppUrl))
	return nil
}
