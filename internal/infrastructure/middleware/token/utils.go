package token

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/respond"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

func UserId(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(gin.H)["uid"].(string)), true
}

func UserInfo(ctx *gin.Context) (*respond.User, bool) {
	u, ok := ctx.Get(Userentity)
	if !ok {
		return nil, false
	}
	return u.(*respond.User), true
}

func Isadmin(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(gin.H)["uid"].(string)), false
}

func Authtoken(ctx *gin.Context) (*Token, bool) {
	t, ok := ctx.Get(AuthToken)
	if !ok {
		return nil, false
	}
	return t.(*Token), true
}
