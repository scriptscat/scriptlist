package gray_control

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

// Control 控制策略
type Control interface {
	// Match 匹配
	Match(ctx *gin.Context, targetScript *script_entity.Code) (bool, error)
}
