package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptweb/internal/http/dto/respond"
	"github.com/scriptscat/scriptweb/pkg/middleware/token"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

const SelfInfo = "self-info"

func userId(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(token.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(gin.H)["uid"].(string)), true
}

func selfinfo(ctx *gin.Context) (*respond.User, bool) {
	u, ok := ctx.Get(SelfInfo)
	if !ok {
		return nil, false
	}
	return u.(*respond.User), true
}

func isadmin(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(token.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(gin.H)["uid"].(string)), false
}

func authtoken(ctx *gin.Context) (*token.Token, bool) {
	t, ok := ctx.Get(token.AuthToken)
	if !ok {
		return nil, false
	}
	return t.(*token.Token), true
}
