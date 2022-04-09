package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
	"github.com/scriptscat/scriptlist/internal/infrastructure/logs"
	"github.com/scriptscat/scriptlist/internal/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/internal/interfaces/api"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/database"
	"github.com/scriptscat/scriptlist/internal/pkg/kvdb"
)

func main() {
	cfg := "config.yaml"
	flag.StringVar(&cfg, "config", cfg, "配置文件")
	flag.Parse()
	if err := config.Init(cfg); err != nil {
		log.Fatal("config error: ", err)
	}
	logs.InitLogs()

	switch config.AppConfig.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewDatabase(config.AppConfig.Mysql, config.AppConfig.Mode == gin.DebugMode)
	if err != nil {
		log.Fatal("database error: ", err)
	}
	redis, err := kvdb.NewKvDb(config.AppConfig.Redis)
	if err != nil {
		log.Fatal("kvdb error: ", err)
	}
	cacheKv, err := kvdb.NewKvDb(config.AppConfig.Cache)
	if err != nil {
		log.Fatal("cache kvdb error: ", err)
	}
	cache := cache.NewRedisCache(cacheKv)
	repo := persistence.NewRepositories(db, redis, cache)
	if err := repo.Migrations(); err != nil {
		log.Fatal("database error: ", err)
	}

	if err := api.StartApi(repo); err != nil {
		log.Fatal("apis error: ", err)
	}

}
