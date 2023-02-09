package gray_control

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type And struct {
	controls []Control
}

func NewAnd(controls ...Control) *And {
	return &And{
		controls: controls,
	}
}

func (a *And) Match(ctx *gin.Context, target *script_entity.Code) (bool, error) {
	for _, v := range a.controls {
		ret, err := v.Match(ctx, target)
		if err != nil {
			return false, err
		}
		if !ret {
			return false, nil
		}
	}
	return true, nil
}

func (a *And) Append(control Control) *And {
	a.controls = append(a.controls, control)
	return a
}

type Or struct {
	controls []Control
}

func NewOr(controls ...Control) *Or {
	return &Or{controls: controls}
}

func (o *Or) Match(ctx *gin.Context, target *script_entity.Code) (bool, error) {
	for _, v := range o.controls {
		ret, err := v.Match(ctx, target)
		if err != nil {
			return false, err
		}
		if ret {
			return true, nil
		}
	}
	return false, nil
}

func (o *Or) Append(control Control) *Or {
	o.controls = append(o.controls, control)
	return o
}
