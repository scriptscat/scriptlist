package gray_control

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type PreRelease struct {
	isPreRelease bool
}

func NewPreRelease(isPreRelease bool) Control {
	return &PreRelease{
		isPreRelease: isPreRelease,
	}
}

func (p *PreRelease) Match(ctx *gin.Context, target *script_entity.Code) (bool, error) {
	return p.isPreRelease, nil
}
