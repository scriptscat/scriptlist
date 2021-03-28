package main

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/script_web/config"
	"github.com/scriptscat/script_web/internal/service"
	"github.com/scriptscat/script_web/internal/web"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
	"log"
	"strconv"
)

func main() {
	if err := config.Init("config.yaml"); err != nil {
		log.Fatal("config error: ", err)
	}

	gin.SetMode(config.AppConfig.Mode)

	db, err := gorm.Open(mysql.Open(config.AppConfig.MySQL.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("database error: ", err)
	}
	script := service.NewScript(db)

	r := gin.Default()

	web.Registry(r, web.NewScript(script))
	_ = r.Run(":" + strconv.Itoa(config.AppConfig.WebPort))
}
