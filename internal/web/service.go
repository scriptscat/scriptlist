package web

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/script_web/internal/service"
	"strings"
)

type Script struct {
	svc service.Script
}

func NewScript(svc service.Script) *Script {
	return &Script{
		svc: svc,
	}
}

func (s *Script) downloadScript(ctx *gin.Context) {

}

func (s *Script) getScriptMeta(ctx *gin.Context) {

}

func (s *Script) Registry(r *gin.Engine) {
	r.Use(func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.RequestURI, ".user.js") {
			s.downloadScript(ctx)
		} else if strings.HasSuffix(ctx.Request.RequestURI, ".meta.js") {
			s.getScriptMeta(ctx)
		}
	})
}
