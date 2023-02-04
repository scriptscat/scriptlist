package script_svc

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type ScriptSvc interface {
	// List 获取脚本列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Create 创建脚本/库
	Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error)
	// UpdateCode 更新脚本/库代码
	UpdateCode(ctx context.Context, req *api.UpdateCodeRequest) (*api.UpdateCodeResponse, error)
	// MigrateEs 全量迁移数据到es
	MigrateEs()
	// GetCode 获取脚本代码,version为latest时获取最新版本
	GetCode(ctx context.Context, id int64, version string) (*script_entity.Code, error)
	// Info 获取脚本信息
	Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error)
	// Code 获取脚本代码
	Code(ctx context.Context, req *api.CodeRequest) (*api.CodeResponse, error)
	// VersionList 获取版本列表
	VersionList(ctx context.Context, req *api.VersionListRequest) (*api.VersionListResponse, error)
	// VersionCode 获取指定版本代码
	VersionCode(ctx context.Context, req *api.VersionCodeRequest) (*api.VersionCodeResponse, error)
	// State 脚本关注等
	State(ctx context.Context, req *api.StateRequest) (*api.StateResponse, error)
	// Watch 关注脚本
	Watch(ctx context.Context, req *api.WatchRequest) (*api.WatchResponse, error)
	// GetSetting 获取脚本设置
	GetSetting(ctx context.Context, req *api.GetSettingRequest) (*api.GetSettingResponse, error)
	// UpdateSetting 更新脚本设置
	UpdateSetting(ctx context.Context, req *api.UpdateSettingRequest) (*api.UpdateSettingResponse, error)
	// SyncOnce 同步一次
	SyncOnce(ctx context.Context, script *script_entity.Script) error
	// Archive 归档脚本
	Archive(ctx context.Context, req *api.ArchiveRequest) (*api.ArchiveResponse, error)
	// Delete 删除脚本
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
	// ToScript 转换为script response结构
	ToScript(ctx context.Context, item *script_entity.Script, withcode bool, version string) (*api.Script, error)
}

type scriptSvc struct {
}

var defaultScript = &scriptSvc{}

func Script() ScriptSvc {
	return defaultScript
}

// List 获取脚本列表
func (s *scriptSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	resp, total, err := script_repo.Script().Search(ctx, &script_repo.SearchOptions{
		Keyword:  req.Keyword,
		Type:     req.ScriptType,
		Sort:     req.Sort,
		Category: make([]int64, 0),
	}, req.PageRequest)
	if err != nil {
		return nil, err
	}
	list := make([]*api.Script, 0)
	for _, item := range resp {
		data, err := s.ToScript(ctx, item, false, "")
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Script]{
			List:  list,
			Total: total,
		},
	}, nil
}

func (s *scriptSvc) ToScript(ctx context.Context, item *script_entity.Script, withcode bool, version string) (*api.Script, error) {
	data := &api.Script{
		ID:          item.ID,
		UserInfo:    user_entity.UserInfo{},
		PostID:      item.PostID,
		Name:        item.Name,
		Description: item.Description,
		Status:      item.Status,
		Type:        int(item.Type),
		Public:      int(item.Public),
		Unwell:      int(item.Unwell),
		Archive:     int(item.Archive),
		Createtime:  item.Createtime,
		Updatetime:  item.Updatetime,
	}
	user, err := user_repo.User().Find(ctx, item.UserID)
	if err != nil {
		logger.Ctx(ctx).Error("获取用户信息失败", zap.Error(err), zap.Int64("user_id", item.UserID))
	}
	data.UserInfo = user.UserInfo()
	// 评分统计信息
	statistics, err := script_repo.ScriptStatistics().FindByScriptID(ctx, item.ID)
	if err != nil {
		logger.Ctx(ctx).Error("获取统计信息失败", zap.Error(err), zap.Int64("script_id", item.ID))
	}
	if statistics != nil {
		data.Score = statistics.Score
		data.ScoreNum = statistics.ScoreCount
	}
	// 从平台统计拿数据,排序从脚本统计里拿数据
	num, err := statistics_repo.Statistics().TotalPv(ctx, item.ID, statistics_repo.DownloadStatistics)
	if err != nil {
		logger.Ctx(ctx).Error("获取统计信息失败", zap.Error(err), zap.Int64("script_id", item.ID))
	}
	data.TotalInstall = num
	num, err = statistics_repo.Statistics().DaysPvNum(ctx, item.ID, statistics_repo.DownloadStatistics, 1, time.Now())
	if err != nil {
		logger.Ctx(ctx).Error("获取统计信息失败", zap.Error(err), zap.Int64("script_id", item.ID))
	}
	data.TodayInstall = num
	// 脚本代码信息
	var scriptCode *script_entity.Code
	if version == "" {
		scriptCode, err = script_repo.ScriptCode().FindLatest(ctx, item.ID, withcode)
	} else {
		scriptCode, err = script_repo.ScriptCode().FindByVersion(ctx, item.ID, version, withcode)
	}
	if err != nil {
		logger.Ctx(ctx).Error("获取脚本代码失败", zap.Error(err), zap.Int64("script_id", item.ID))
	}
	if scriptCode == nil {
		return nil, i18n.NewError(ctx, code.ScriptNotFound)
	}
	data.Script = s.scriptCode(ctx, scriptCode)
	// 脚本分类信息
	list, err := script_repo.ScriptCategory().List(ctx, item.ID)
	if err != nil {
		logger.Ctx(ctx).Error("获取脚本分类失败", zap.Error(err), zap.Int64("script_id", item.ID))
	}
	data.Category = make([]*api.CategoryList, 0)
	for _, v := range list {
		category, err := script_repo.ScriptCategoryList().Find(ctx, v.CategoryID)
		if err != nil {
			logger.Ctx(ctx).Error("获取分类信息失败", zap.Error(err), zap.Int64("category_id", v.CategoryID))
		}
		if category != nil {
			data.Category = append(data.Category, &api.CategoryList{
				ID:   category.ID,
				Name: category.Name,
			})
		}
	}
	return data, nil
}

func (s *scriptSvc) scriptCode(ctx context.Context, code *script_entity.Code) *api.Code {
	metaJson := make(map[string]interface{})
	if err := json.Unmarshal([]byte(code.MetaJson), &metaJson); err != nil {
		logger.Ctx(ctx).Error("json解析失败", zap.Error(err),
			zap.String("meta", code.MetaJson), zap.Int64("script_id", code.ScriptID), zap.Int64("code_id", code.ID))
	}
	ret := &api.Code{
		ID:         code.ID,
		MetaJson:   metaJson,
		ScriptID:   code.ScriptID,
		Version:    code.Version,
		Changelog:  code.Changelog,
		Status:     code.Status,
		Createtime: code.Createtime,
		Code:       code.Code,
	}
	return ret
}

// Create 创建脚本
func (s *scriptSvc) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	script := &script_entity.Script{
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    req.Content,
		Type:       req.Type,
		Public:     req.Public,
		Unwell:     req.Unwell,
		Status:     consts.ACTIVE,
		Archive:    script_entity.IsActive,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	// 保存脚本代码
	scriptCode := &script_entity.Code{
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Changelog:  req.Changelog,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: 0,
	}
	var definition *script_entity.LibDefinition
	if req.Type == script_entity.LibraryType {
		// 脚本引用库
		script.Name = req.Name
		script.Description = req.Description
		scriptCode.Code = req.Code
		scriptCode.Version = req.Version
		// 脚本定义
		if req.Definition != "" {
			definition = &script_entity.LibDefinition{
				UserID:     auth_svc.Auth().Get(ctx).UID,
				Definition: req.Definition,
				Createtime: time.Now().Unix(),
			}
		}
	} else {
		metaJson, err := scriptCode.UpdateCode(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		script.Name = metaJson["name"][0]
		script.Description = metaJson["description"][0]
	}

	// 保存数据库并发送消息
	if err := script_repo.Script().Create(ctx, script); err != nil {
		logger.Ctx(ctx).Error("scriptSvc create failed", zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptCreateFailed,
		)
	}
	// 保存脚本代码
	scriptCode.ScriptID = script.ID
	if err := script_repo.ScriptCode().Create(ctx, scriptCode); err != nil {
		logger.Ctx(ctx).Error("scriptSvc code create failed", zap.Int64("script_id", script.ID), zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptCreateFailed,
		)
	}
	// 保存定义
	if definition != nil {
		definition.ScriptID = script.ID
		definition.CodeID = scriptCode.ID
		if err := script_repo.LibDefinition().Create(ctx, definition); err != nil {
			logger.Ctx(ctx).Error("scriptSvc definition create failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptCreateFailed,
			)
		}
	}

	if err := producer.PublishScriptCreate(ctx, script, scriptCode); err != nil {
		logger.Ctx(ctx).Error("publish scriptSvc create failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
		return nil, i18n.NewInternalError(ctx, code.ScriptCreateFailed)
	}
	return &api.CreateResponse{ID: script.ID}, nil
}

// UpdateCode 更新脚本/库代码
func (s *scriptSvc) UpdateCode(ctx context.Context, req *api.UpdateCodeRequest) (*api.UpdateCodeResponse, error) {
	// 搜索到脚本
	script, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := script.CheckPermission(ctx); err != nil {
		return nil, err
	}
	if err := script.IsArchive(ctx); err != nil {
		return nil, err
	}
	scriptCode := &script_entity.Code{
		UserID:    auth_svc.Auth().Get(ctx).UID,
		ScriptID:  script.ID,
		Changelog: req.Changelog,
		Status:    consts.ACTIVE,
	}
	var definition *script_entity.LibDefinition
	if script.Type == script_entity.LibraryType {
		oldVersion, err := script_repo.ScriptCode().FindByVersion(ctx, script.ID, req.Version, true)
		if err != nil {
			return nil, err
		}
		if oldVersion != nil {
			return nil, i18n.NewError(ctx, code.ScriptVersionExist)
		}
		// 脚本引用库
		scriptCode.Code = req.Code
		scriptCode.Version = req.Version
		// 脚本定义
		if req.Definition != "" {
			definition = &script_entity.LibDefinition{
				UserID:     auth_svc.Auth().Get(ctx).UID,
				Definition: req.Definition,
				Createtime: time.Now().Unix(),
			}
		}
	} else {
		metaJson, err := scriptCode.UpdateCode(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		oldVersion, err := script_repo.ScriptCode().FindByVersion(ctx, script.ID, metaJson["version"][0], true)
		if err != nil {
			return nil, err
		}
		if oldVersion != nil {
			// 如果脚本内容发生了改变但版本号没有发生改变
			if strings.ReplaceAll(oldVersion.Code, "\r\n", "\n") != strings.ReplaceAll(scriptCode.Code, "\r\n", "\n") {
				return nil, i18n.NewError(ctx, code.ScriptVersionExist)
			}
			scriptCode.ID = oldVersion.ID
			scriptCode.Createtime = oldVersion.Createtime
		} else {
			script.Updatetime = time.Now().Unix()
			scriptCode.Createtime = time.Now().Unix()
		}
		// 更新名字和描述
		script.Name = metaJson["name"][0]
		script.Description = metaJson["description"][0]
	}
	script.Content = req.Content
	script.Public = req.Public
	script.Unwell = req.Unwell
	// 保存数据库并发送消息
	if err := script_repo.Script().Update(ctx, script); err != nil {
		logger.Ctx(ctx).Error("scriptSvc update failed", zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptUpdateFailed,
		)
	}
	// 根据id判断是新建还是更新
	if scriptCode.ID == 0 {
		if err := script_repo.ScriptCode().Create(ctx, scriptCode); err != nil {
			logger.Ctx(ctx).Error("scriptSvc code create failed", zap.Int64("script_id", script.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptUpdateFailed,
			)
		}
	} else {
		if err := script_repo.ScriptCode().Update(ctx, scriptCode); err != nil {
			logger.Ctx(ctx).Error("scriptSvc code update failed", zap.Int64("script_id", script.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptUpdateFailed,
			)
		}
	}

	// 保存定义
	if definition != nil {
		definition.ScriptID = script.ID
		definition.CodeID = scriptCode.ID
		if err := script_repo.LibDefinition().Create(ctx, definition); err != nil {
			logger.Ctx(ctx).Error("scriptSvc definition create failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptUpdateFailed,
			)
		}
	}

	if err := producer.PublishScriptCodeUpdate(ctx, script, scriptCode); err != nil {
		logger.Ctx(ctx).Error("publish scriptSvc code update failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
		return nil, i18n.NewInternalError(ctx, code.ScriptUpdateFailed)
	}
	return &api.UpdateCodeResponse{}, nil
}

// MigrateEs 全量迁移数据到es
func (s *scriptSvc) MigrateEs() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := 0
	for {
		if ok := func(ctx context.Context) bool {
			ctx, span := trace.Default().Tracer("MigrateEs").Start(context.Background(), "MigrateEs")
			defer func() {
				span.End()
				start += 20
			}()
			ctx = logger.ContextWithLogger(ctx, logger.Ctx(ctx).With(trace.LoggerLabel(ctx)...))
			list, err := script_repo.Migrate().List(ctx, start, 20)
			if err != nil {
				logger.Ctx(ctx).Error("获取迁移数据失败", zap.Error(err))
				return false
			}
			if len(list) == 0 {
				logger.Ctx(ctx).Info("迁移完成")
				return true
			}
			for _, item := range list {
				search, err := script_repo.Migrate().Convert(ctx, item)
				if err != nil {
					logger.Ctx(ctx).Error("转换数据失败", zap.Error(err))
					continue
				}
				if err := script_repo.Migrate().Save(ctx, search); err != nil {
					logger.Ctx(ctx).Error("保存数据失败", zap.Error(err))
					continue
				}
			}
			return false
		}(ctx); ok {
			break
		}
	}
}

// GetCode 获取脚本代码,version为latest时获取最新版本
func (s *scriptSvc) GetCode(ctx context.Context, id int64, version string) (*script_entity.Code, error) {
	if version == "latest" {
		return script_repo.ScriptCode().FindLatest(ctx, id, true)
	}
	return script_repo.ScriptCode().FindByVersion(ctx, id, version, true)
}

// Info 获取脚本信息
func (s *scriptSvc) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckOperate(ctx); err != nil {
		return nil, err
	}
	script, err := s.ToScript(ctx, m, false, "")
	if err != nil {
		return nil, err
	}
	return &api.InfoResponse{
		Script:  script,
		Content: m.Content,
	}, nil
}

// Code 获取脚本代码
func (s *scriptSvc) Code(ctx context.Context, req *api.CodeRequest) (*api.CodeResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckOperate(ctx); err != nil {
		return nil, err
	}
	script, err := s.ToScript(ctx, m, true, "")
	if err != nil {
		return nil, err
	}
	return &api.CodeResponse{
		Script:  script,
		Content: m.Content,
	}, nil
}

// VersionList 获取版本列表
func (s *scriptSvc) VersionList(ctx context.Context, req *api.VersionListRequest) (*api.VersionListResponse, error) {
	list, total, err := script_repo.ScriptCode().List(ctx, req.ID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	ret := &api.VersionListResponse{
		PageResponse: httputils.PageResponse[*api.Code]{
			Total: total,
			List:  make([]*api.Code, len(list)),
		},
	}
	for n, v := range list {
		ret.List[n] = s.scriptCode(ctx, v)
	}
	return ret, nil
}

// VersionCode 获取指定版本代码
func (s *scriptSvc) VersionCode(ctx context.Context, req *api.VersionCodeRequest) (*api.VersionCodeResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckOperate(ctx); err != nil {
		return nil, err
	}
	script, err := s.ToScript(ctx, m, true, req.Version)
	if err != nil {
		return nil, err
	}
	return &api.VersionCodeResponse{
		Script: script,
	}, nil
}

// State 获取脚本状态,脚本关注等
func (s *scriptSvc) State(ctx context.Context, req *api.StateRequest) (*api.StateResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckOperate(ctx); err != nil {
		return nil, err
	}
	user := auth_svc.Auth().Get(ctx)
	var level script_entity.ScriptWatchLevel
	if user != nil {
		watch, err := script_repo.ScriptWatch().FindByUser(ctx, req.ID, user.UID)
		if err != nil {
			return nil, err
		}
		if watch != nil {
			level = watch.Level
		}
	}
	return &api.StateResponse{
		Watch: level,
	}, nil
}

// Watch 关注脚本
func (s *scriptSvc) Watch(ctx context.Context, req *api.WatchRequest) (*api.WatchResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckOperate(ctx); err != nil {
		return nil, err
	}
	if err := script_repo.ScriptWatch().Watch(ctx, req.ID, auth_svc.Auth().Get(ctx).UID, req.Watch); err != nil {
		return nil, err
	}
	return nil, nil
}

// GetSetting 获取脚本设置
func (s *scriptSvc) GetSetting(ctx context.Context, req *api.GetSettingRequest) (*api.GetSettingResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckPermission(ctx); err != nil {
		return nil, err
	}
	return &api.GetSettingResponse{
		SyncUrl:       m.SyncUrl,
		ContentUrl:    m.ContentUrl,
		DefinitionUrl: m.DefinitionUrl,
		SyncMode:      m.SyncMode,
	}, nil
}

// UpdateSetting 更新脚本设置
func (s *scriptSvc) UpdateSetting(ctx context.Context, req *api.UpdateSettingRequest) (*api.UpdateSettingResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckPermission(ctx); err != nil {
		return nil, err
	}
	if err := m.IsArchive(ctx); err != nil {
		return nil, err
	}
	m.SyncUrl = req.SyncUrl
	m.ContentUrl = req.ContentUrl
	m.SyncMode = req.SyncMode
	switch m.Type {
	case script_entity.UserscriptType, script_entity.SubscribeType:
	case script_entity.LibraryType:
		m.Name = req.Name
		m.Description = req.Description
		m.DefinitionUrl = req.DefinitionUrl
	default:
		return nil, i18n.NewError(ctx, code.ScriptUpdateFailed)
	}
	if err := script_repo.Script().Update(ctx, m); err != nil {
		return nil, err
	}
	err = s.SyncOnce(ctx, m)
	if err != nil {
		return nil, err
	}
	return &api.UpdateSettingResponse{}, nil
}

func (s *scriptSvc) SyncOnce(ctx context.Context, script *script_entity.Script) error {
	if err := script.IsArchive(ctx); err != nil {
		return err
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", script.ID))
	// 读取代码
	codeContent, err := requestSyncUrl(ctx, script.SyncUrl)
	if err != nil {
		logger.Error("读取代码失败", zap.String("sync_url", script.SyncUrl), zap.Error(err))
		return err
	}
	// 查找最新的code
	code, err := script_repo.ScriptCode().FindLatest(ctx, script.ID, false)
	if err != nil {
		logger.Error("获取代码失败", zap.Error(err))
		return err
	}
	if code == nil {
		logger.Error("代码不存在")
		return err
	}
	oldVersion := code.Version
	if _, err := code.UpdateCode(ctx, codeContent); err != nil {
		logger.Error("解析代码失败", zap.String("sync_url", script.SyncUrl), zap.Error(err))
		return err
	}
	if oldVersion == code.Version {
		logger.Info("版本相同,略过", zap.String("sync_url", script.SyncUrl))
		return err
	}
	req := &api.UpdateCodeRequest{
		ID:         script.ID,
		Version:    "",
		Content:    "",
		Code:       codeContent,
		Definition: "",
		Changelog:  "该版本为系统自动同步更新",
		Public:     script.Public,
		Unwell:     script.Unwell,
	}
	// 读取content
	if script.ContentUrl != "" {
		content, err := requestSyncUrl(ctx, script.ContentUrl)
		if err != nil {
			logger.Error("读取content失败", zap.String("content_url", script.ContentUrl), zap.Error(err))
			req.Content = script.Content
		} else {
			req.Content = content
		}
	}
	if script.Type == script_entity.LibraryType {
		// 版本号,最后一位加一
		end := strings.LastIndex(code.Version, ".")
		if end == -1 {
			code.Version = code.Version + ".1"
		} else {
			ver, _ := strconv.Atoi(code.Version[end+1:])
			code.Version = code.Version[:end] + "." + strconv.Itoa(ver+1)
		}
	} else {
		req.Version = code.Version
	}
	if _, err := s.UpdateCode(ctx, req); err != nil {
		logger.Error("更新代码失败", zap.String("sync_url", script.SyncUrl), zap.Error(err))
		return err
	}
	return nil
}

// Archive 归档脚本
func (s *scriptSvc) Archive(ctx context.Context, req *api.ArchiveRequest) (*api.ArchiveResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckPermission(ctx); err != nil {
		return nil, err
	}
	if req.Archive {
		m.Archive = script_entity.IsArchive
	} else {
		m.Archive = script_entity.IsActive
	}
	if err := script_repo.Script().Update(ctx, m); err != nil {
		return nil, err
	}
	return &api.ArchiveResponse{}, nil
}

// Delete 删除脚本
func (s *scriptSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	m, err := script_repo.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := m.CheckPermission(ctx, model.SuperModerator); err != nil {
		return nil, err
	}
	m.Status = consts.DELETE
	if err := script_repo.Script().Update(ctx, m); err != nil {
		return nil, err
	}
	if err := producer.PublishScriptDelete(ctx, m); err != nil {
		logger.Ctx(ctx).Error("发布删除脚本消息失败", zap.Error(err))
	}
	return &api.DeleteResponse{}, nil
}
