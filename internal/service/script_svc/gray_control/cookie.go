package gray_control

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type Cookie struct {
	regex string
}

func NewCookie(regex string) Control {
	return &Cookie{
		regex: regex,
	}
}

func (c *Cookie) Match(ctx *gin.Context, target *script_entity.Code) (bool, error) {
	cookie := ctx.GetHeader("cookie")
	return regexp.Match(c.regex, []byte(cookie))
}
