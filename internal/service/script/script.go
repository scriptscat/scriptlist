package script

import (
	"context"
	"encoding/json"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/pkg/consts"
	"github.com/scriptscat/scriptlist/internal/repository"
	"github.com/scriptscat/scriptlist/internal/service/user"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type IScript interface {
	// List 获取脚本列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Create 创建脚本
	Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error)
}

type script struct {
}

var defaultScript = &script{}

func Script() IScript {
	return defaultScript
}

// List 获取脚本列表
func (s *script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
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
	if req.Type == entity.LibraryType {
		// 脚本引用库
		script.Name = req.Name
		script.Description = req.Description
		scriptCode.Code = req.Code
		scriptCode.Version = req.Version
	} else {
		// 解析脚本的元数据
		scriptCodeStr, meta, err := parseCodeMeta(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		// 解析元数据
		metaJson := parseMetaToJson(meta)
		if len(metaJson["name"]) == 0 {
			return nil, i18n.NewError(ctx, code.ScriptNameIsEmpty)
		}
		if len(metaJson["description"]) == 0 {
			return nil, i18n.NewError(ctx, code.ScriptDescIsEmpty)
		}
		if len(metaJson["version"]) == 0 {
			return nil, i18n.NewError(ctx, code.ScriptVersionIsEmpty)
		}
		script.Name = metaJson["name"][0]
		script.Description = metaJson["description"][0]
		scriptCode.Code = scriptCodeStr
		scriptCode.Meta = meta
		b, err := json.Marshal(metaJson)
		if err != nil {
			return nil, i18n.NewError(ctx, code.ScriptParseFailed)
		}
		scriptCode.MetaJson = string(b)
		scriptCode.Version = metaJson["version"][0]
	}

	// 保存数据库并发送消息
	if err := repository.Script().Create(ctx, script); err != nil {
		logger.Ctx(ctx).Error("script create failed", zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptCreateFailed,
		)
	}
	// 保存脚本代码
	scriptCode.ScriptID = script.ID
	if err := repository.ScriptCode().Create(ctx, scriptCode); err != nil {
		logger.Ctx(ctx).Error("script code create failed", zap.Error(err))
		return nil, i18n.NewInternalError(
			ctx,
			code.ScriptCreateFailed,
		)
	}

	if err := producer.PublishScriptCreate(ctx, script); err != nil {
		logger.Ctx(ctx).Error("publish script create failed", zap.Error(err))
		return nil, i18n.NewInternalError(ctx, code.ScriptCreateFailed)
	}
	return nil, nil
}
