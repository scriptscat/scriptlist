package script_ctr

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/limit"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"github.com/scriptscat/scriptlist/internal/service/statistics_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type Script struct {
	limit *limit.PeriodLimit
}

func NewScript() *Script {
	return &Script{
		limit: limit.NewPeriodLimit(
			300, 10, redis.Default(), "limit:create:script",
		),
	}
}

// List 获取脚本列表
func (s *Script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return script_svc.Script().List(ctx, req)
}

// Create 创建脚本/库
func (s *Script) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	resp, err := s.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return script_svc.Script().Create(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CreateResponse), nil
}

// UpdateCode 更新脚本/库代码
func (s *Script) UpdateCode(ctx context.Context, req *api.UpdateCodeRequest) (*api.UpdateCodeResponse, error) {
	resp, err := s.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return script_svc.Script().UpdateCode(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.UpdateCodeResponse), nil
}

// MigrateEs 全量迁移数据到es
func (s *Script) MigrateEs(ctx context.Context, req *api.MigrateEsRequest) (*api.MigrateEsResponse, error) {
	if auth_svc.Auth().Get(ctx).AdminLevel != model.Admin {
		return nil, httputils.NewError(http.StatusForbidden, -1, "无权限")
	}
	go script_svc.Script().MigrateEs()
	return &api.MigrateEsResponse{}, nil
}

func (s *Script) Download() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.Path, ".user.js") || strings.HasSuffix(ctx.Request.URL.Path, ".user.sub.js") {
			s.downloadScript(ctx)
		} else if strings.HasSuffix(ctx.Request.URL.Path, ".meta.js") {
			s.getScriptMeta(ctx)
		} else {
			ctx.AbortWithStatus(http.StatusNotFound)
		}
	}
}

func (s *Script) getScriptID(ctx *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, httputils.NewError(http.StatusNotFound, -1, "id不能为空")
	}
	return id, nil
}

func (s *Script) downloadScript(ctx *gin.Context) {
	id, err := s.getScriptID(ctx)
	if id == 0 {
		httputils.HandleResp(ctx, err)
		return
	}
	version := ctx.Query("version")
	if version == "" {
		version = "latest"
	}
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	// 获取脚本
	code, err := script_svc.Script().GetCode(ctx, id, version)
	if err != nil {
		httputils.HandleResp(ctx, err)
		return
	}
	record := &producer.ScriptStatisticsMsg{
		ScriptID:        code.ScriptID,
		ScriptCodeID:    code.ID,
		UserID:          0,
		IP:              ctx.ClientIP(),
		UA:              ua,
		StatisticsToken: statistics_svc.Statistics().GetStatisticsToken(ctx),
		Download:        statistics_repo.DownloadStatistics,
		Time:            time.Now(),
	}
	user := auth_svc.Auth().Get(ctx)
	if user != nil {
		record.UserID = user.UID
	}
	err = statistics_svc.Statistics().ScriptRecord(ctx, record)
	if err != nil {
		logger.Ctx(ctx).Error("脚本下载统计记录失败", zap.Any("record", record), zap.Error(err))
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Code)
}

func (s *Script) getScriptMeta(ctx *gin.Context) {
	id, err := s.getScriptID(ctx)
	if err != nil {
		httputils.HandleResp(ctx, err)
		return
	}
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	// 获取脚本
	code, err := script_svc.Script().GetCode(ctx, id, "latest")
	if err != nil {
		httputils.HandleResp(ctx, err)
		return
	}
	record := &producer.ScriptStatisticsMsg{
		ScriptID:        code.ScriptID,
		ScriptCodeID:    code.ID,
		UserID:          0,
		IP:              ctx.ClientIP(),
		UA:              ua,
		StatisticsToken: statistics_svc.Statistics().GetStatisticsToken(ctx),
		Download:        statistics_repo.UpdateStatistics,
		Time:            time.Now(),
	}
	user := auth_svc.Auth().Get(ctx)
	if user != nil {
		record.UserID = user.UID
	}
	err = statistics_svc.Statistics().ScriptRecord(ctx, record)
	if err != nil {
		logger.Ctx(ctx).Error("脚本下载统计记录失败", zap.Any("record", record), zap.Error(err))
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Meta)
}

// Info 获取脚本信息
func (s *Script) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return script_svc.Script().Info(ctx, req)
}

// Code 获取脚本代码
func (s *Script) Code(ctx context.Context, req *api.CodeRequest) (*api.CodeResponse, error) {
	return script_svc.Script().Code(ctx, req)
}

// VersionList 获取版本列表
func (s *Script) VersionList(ctx context.Context, req *api.VersionListRequest) (*api.VersionListResponse, error) {
	return script_svc.Script().VersionList(ctx, req)
}

// VersionCode 获取指定版本代码
func (s *Script) VersionCode(ctx context.Context, req *api.VersionCodeRequest) (*api.VersionCodeResponse, error) {
	return script_svc.Script().VersionCode(ctx, req)
}

// State 获取脚本状态,脚本关注等
func (s *Script) State(ctx context.Context, req *api.StateRequest) (*api.StateResponse, error) {
	return script_svc.Script().State(ctx, req)
}

// Watch 关注脚本
func (s *Script) Watch(ctx context.Context, req *api.WatchRequest) (*api.WatchResponse, error) {
	return script_svc.Script().Watch(ctx, req)
}

// GetSetting 获取脚本设置
func (s *Script) GetSetting(ctx context.Context, req *api.GetSettingRequest) (*api.GetSettingResponse, error) {
	return script_svc.Script().GetSetting(ctx, req)
}

var whiteList = []string{
	"github.com",
	"github.io",
	"raw.githubusercontent.com",
	"gitlab.com",
	"greasyfork.org",
	"scriptcat.org",
	"zhaojiaoben.cn",
	"gitee.com",
	"jsdelivr.net",
}

// UpdateSetting 更新脚本设置
func (s *Script) UpdateSetting(ctx context.Context, req *api.UpdateSettingRequest) (*api.UpdateSettingResponse, error) {
	// 允许域名白名单
	if req.SyncUrl != "" {
		u, err := url.Parse(req.SyncUrl)
		if err != nil {
			return nil, err
		}
		var flag bool
		for _, v := range whiteList {
			if strings.Contains(u.Host, v) {
				flag = true
			}
		}
		if !flag {
			return nil, i18n.NewError(ctx, code.ScriptNotAllowUrl)
		}
	}
	return script_svc.Script().UpdateSetting(ctx, req)
}

// Archive 归档脚本
func (s *Script) Archive(ctx context.Context, req *api.ArchiveRequest) (*api.ArchiveResponse, error) {
	return script_svc.Script().Archive(ctx, req)
}

// Delete 删除脚本
func (s *Script) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return script_svc.Script().Delete(ctx, req)
}
