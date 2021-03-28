package web

import "github.com/gin-gonic/gin"

type Service interface {
	Registry(r *gin.Engine)
}

func Registry(r *gin.Engine, registry ...Service) {
	for _, v := range registry {
		v.Registry(r)
	}
}
