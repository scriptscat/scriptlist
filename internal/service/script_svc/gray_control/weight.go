package gray_control

import (
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type Weight struct {
	weight    int
	weightDay int
}

func NewWeight(weight, weightDay int) Control {
	return &Weight{
		weight:    weight,
		weightDay: weightDay,
	}
}

func (w *Weight) Match(ctx *gin.Context, target *script_entity.Code) (bool, error) {
	weight, err := ctx.Cookie("gray_weight")
	if err != nil {
		return false, err
	}
	var n int
	if weight == "" {
		n = rand.Intn(100) + 1
		ctx.SetCookie("gray_weight", strconv.Itoa(n), 0, "/", "", false, true)
	}
	n, _ = strconv.Atoi(weight)
	return n <= w.weight, nil
}
