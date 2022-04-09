package api

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

func GetStatisticsToken(ctx *gin.Context) string {
	stk, _ := ctx.Cookie("_statistics")
	if stk == "" {
		stk = utils.RandString(32, 2)
		ctx.SetCookie("_statistics", stk, 3600*24*365, "/", "", false, true)
	}
	return stk
}
