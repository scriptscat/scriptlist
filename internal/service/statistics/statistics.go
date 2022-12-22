package statistics

import (
	"context"

	"github.com/codfrm/cago/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

const (
	ViewStatistics     = "view"
	DownloadStatistics = "download"
	UpdateStatistics   = "update"
)

// IStatistics 统计平台
type IStatistics interface {
	// ScriptRecord 脚本数据统计
	ScriptRecord(ctx context.Context, data *producer.ScriptStatisticsMsg) error
	// GetStatisticsToken 获取统计token
	GetStatisticsToken(ctx *gin.Context) string
}

type statistics struct {
}

var defaultStatistics = &statistics{}

func Statistics() IStatistics {
	return defaultStatistics
}

func (s *statistics) ScriptRecord(ctx context.Context, data *producer.ScriptStatisticsMsg) error {
	return producer.PublishScriptStatistics(ctx, data)
}

func (s *statistics) GetStatisticsToken(ctx *gin.Context) string {
	stk, _ := ctx.Cookie("_statistics")
	if stk == "" {
		stk = utils.RandString(32, 2)
		ctx.SetCookie("_statistics", stk, 3600*24*365, "/", "", false, true)
	}
	return stk
}
