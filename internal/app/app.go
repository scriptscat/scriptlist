package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/http"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"github.com/scriptscat/scriptlist/migrations"
)

func Run() {
	switch config.AppConfig.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	if err := db.Init(); err != nil {
		log.Fatal("database error: ", err)
	}
	if err := migrations.Migrate(); err != nil {
		log.Fatal("migrate error: ", err)
	}

	if err := http.StartApi(); err != nil {
		log.Fatal("apis error: ", err)
	}
}
