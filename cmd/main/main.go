package main

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptweb/internal/interfaces/apis"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"github.com/scriptscat/scriptweb/internal/pkg/migrate"
	"log"
)

func main() {
	if err := config.Init("config.yaml"); err != nil {
		log.Fatal("config error: ", err)
	}

	switch config.AppConfig.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	if err := db.Init(); err != nil {
		log.Fatal("database error: ", err)
	}
	if err := migrate.Migrate(); err != nil {
		log.Fatal("migrate error: ", err)
	}

	if err := apis.StartApi(); err != nil {
		log.Fatal("apis error: ", err)
	}

}
