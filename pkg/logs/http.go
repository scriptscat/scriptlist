package logs

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GinLogger() []gin.HandlerFunc {
	var w io.Writer = &lumberjack.Logger{
		Filename:   "./logs/runtime/http.log",
		MaxSize:    2,
		MaxBackups: 30,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}
	if config.AppConfig.Mode == "debug" {
		gin.ForceConsoleColor()
		w = io.MultiWriter(w, os.Stdout)
	}
	return []gin.HandlerFunc{gin.LoggerWithWriter(w), gin.RecoveryWithWriter(w, func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": -1000, "msg": "系统错误"})
	})}
}
