package gray_control

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type Weight struct {
	weight    int
	weightDay float64
}

func NewWeight(weight int, weightDay float64) Control {
	return &Weight{
		weight:    weight,
		weightDay: weightDay,
	}
}

func (w *Weight) Match(ctx *gin.Context, target *script_entity.Code) (bool, error) {
	weight, err := ctx.Cookie("gray_weight")
	if err != nil {
		if err != http.ErrNoCookie {
			return false, err
		}
	}
	var n int
	if weight == "" {
		n = rand.Intn(100)
		ctx.SetCookie("gray_weight", strconv.Itoa(n), 0, "/", "", false, true)
	}
	n, _ = strconv.Atoi(weight)
	return w.match(time.Now(), n, target.Createtime)
}

func (w *Weight) match(now time.Time, n int, createtime int64) (bool, error) {
	weight := w.weight
	if w.weightDay != 0 {
		// 不为0时,计算权重百分比
		wd := (now.Sub(time.Unix(createtime, 0)).Abs().Seconds() / 86400) / w.weightDay
		if wd < 1 {
			weight = int(float64(weight) * wd)
		}
	}
	// 如果不加一个随机数,那么权重低的总会更新
	x := int(createtime) + n
	return (x % 100) <= weight, nil
}
