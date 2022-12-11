package script

import (
	"context"
	"strings"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/pkg/consts"
	script2 "github.com/scriptscat/scriptlist/internal/repository/script"
	"github.com/scriptscat/scriptlist/internal/service/user"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type IScript interface {
	// List 获取脚本列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Create 创建脚本/库
	Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error)
	// UpdateCode 更新脚本/库代码
	UpdateCode(ctx context.Context, req *api.UpdateCodeRequest) (*api.UpdateCodeResponse, error)
}

type script struct {
}

var defaultScript = &script{}

func Script() IScript {
	return defaultScript
}

// List 获取脚本列表
func (s *script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	if req.Keyword != "" {

	}
	return nil, nil
}

// Create 创建脚本
func (s *script) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	script := &entity.Script{
		UserID:     user.Auth().Get(ctx).UID,
		Content:    req.Content,
		Type:       req.Type,
		Public:     req.Public,
		Unwell:     req.Unwell,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	// 保存脚本代码
	scriptCode := &entity.Code{
		UserID:     user.Auth().Get(ctx).UID,
		Changelog:  req.Changelog,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: 0,
	}
	var definition *entity.LibDefinition
	if req.Type == entity.LibraryType {
		// 脚本引用库
		script.Name = req.Name
		script.Description = req.Description
		scriptCode.Code = req.Code
		scriptCode.Version = req.Version
		// 脚本定义
		if req.Definition != "" {
			definition = &entity.LibDefinition{
				UserID:     user.Auth().Get(ctx).UID,
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
	if err := script2.Script().Create(ctx, script); err != nil {
		logger.Ctx(ctx).Error("script create failed", zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptCreateFailed,
		)
	}
	// 保存脚本代码
	scriptCode.ScriptID = script.ID
	if err := script2.ScriptCode().Create(ctx, scriptCode); err != nil {
		logger.Ctx(ctx).Error("script code create failed", zap.Int64("script_id", script.ID), zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptCreateFailed,
		)
	}
	// 保存定义
	if definition != nil {
		definition.ScriptID = script.ID
		definition.CodeID = scriptCode.ID
		if err := script2.LibDefinition().Create(ctx, definition); err != nil {
			logger.Ctx(ctx).Error("script definition create failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptCreateFailed,
			)
		}
	}

	if err := producer.PublishScriptCreate(ctx, script, scriptCode); err != nil {
		logger.Ctx(ctx).Error("publish script create failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
		return nil, i18n.NewInternalError(ctx, code.ScriptCreateFailed)
	}
	return &api.CreateResponse{ID: script.ID}, nil
}

// UpdateCode 更新脚本/库代码
func (s *script) UpdateCode(ctx context.Context, req *api.UpdateCodeRequest) (*api.UpdateCodeResponse, error) {
	// 搜索到脚本
	script, err := script2.Script().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, i18n.NewError(ctx, code.ScriptNotFound)
	}
	if err := script.CheckPermission(ctx); err != nil {
		return nil, err
	}
	scriptCode := &entity.Code{
		UserID:     user.Auth().Get(ctx).UID,
		ScriptID:   script.ID,
		Changelog:  req.Changelog,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	var definition *entity.LibDefinition
	if script.Type == entity.LibraryType {
		oldVersion, err := script2.ScriptCode().FindByVersion(ctx, script.ID, req.Version)
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
			definition = &entity.LibDefinition{
				UserID:     user.Auth().Get(ctx).UID,
				Definition: req.Definition,
				Createtime: time.Now().Unix(),
			}
		}
	} else {
		metaJson, err := scriptCode.UpdateCode(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		oldVersion, err := script2.ScriptCode().FindByVersion(ctx, script.ID, metaJson["version"][0])
		if err != nil {
			return nil, err
		}
		if oldVersion != nil {
			// 如果脚本内容发生了改变但版本号没有发生改变
			if strings.ReplaceAll(oldVersion.Code, "\r\n", "\n") != strings.ReplaceAll(scriptCode.Code, "\r\n", "\n") {
				return nil, i18n.NewError(ctx, code.ScriptVersionExist)
			}
			scriptCode.ID = oldVersion.ID
		} else {
			script.Updatetime = time.Now().Unix()
		}
		// 更新名字和描述
		script.Name = metaJson["name"][0]
		script.Description = metaJson["description"][0]
	}

	// 保存数据库并发送消息
	if err := script2.Script().Update(ctx, script); err != nil {
		logger.Ctx(ctx).Error("script update failed", zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptUpdateFailed,
		)
	}
	// 根据id判断是新建还是更新
	if scriptCode.ID == 0 {
		if err := script2.ScriptCode().Create(ctx, scriptCode); err != nil {
			logger.Ctx(ctx).Error("script code create failed", zap.Int64("script_id", script.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptUpdateFailed,
			)
		}
	} else {
		if err := script2.ScriptCode().Update(ctx, scriptCode); err != nil {
			logger.Ctx(ctx).Error("script code update failed", zap.Int64("script_id", script.ID), zap.Error(err))
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
		if err := script2.LibDefinition().Create(ctx, definition); err != nil {
			logger.Ctx(ctx).Error("script definition create failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
			return nil, i18n.NewInternalError(
				ctx,
				code.ScriptUpdateFailed,
			)
		}
	}

	if err := producer.PublishScriptCodeUpdate(ctx, script, scriptCode); err != nil {
		logger.Ctx(ctx).Error("publish script code update failed", zap.Int64("script_id", script.ID), zap.Int64("code_id", scriptCode.ID), zap.Error(err))
		return nil, i18n.NewInternalError(ctx, code.ScriptUpdateFailed)
	}
	return &api.UpdateCodeResponse{}, nil
}
