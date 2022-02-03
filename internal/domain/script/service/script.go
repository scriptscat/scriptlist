package service

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"github.com/scriptscat/scriptlist/internal/domain/script/broker"
	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	"github.com/scriptscat/scriptlist/internal/domain/script/repository"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/gorm"
)

const (
	SyncModeAuto   = 1
	SyncModeManual = 2
)

type Script interface {
	Search(search *repository.SearchList, page *request.Pages) ([]*entity.Script, int64, error)
	UserScript(uid int64, self bool, page *request.Pages) ([]*entity.Script, int64, error)
	Info(id int64) (*entity.Script, error)
	GetCode(id int64) (*entity.ScriptCode, error)
	GetScriptVersion(scriptId int64, version string) (*entity.ScriptCode, error)
	GetLatestVersion(scriptId int64) (*entity.ScriptCode, error)
	VersionList(id int64, page *request.Pages) ([]*entity.ScriptCode, int64, error)
	GetCategory() ([]*entity.ScriptCategoryList, error)
	Download(id int64) error
	Update(id int64) error
	CreateScript(uid int64, req *request.CreateScript) (*entity.Script, error)
	UpdateScript(uid, id int64, req *request.UpdateScript) error
	CreateScriptCode(uid, id int64, req *request.UpdateScriptCode) error
	GetCodeDefinition(codeid int64) (*entity.LibDefinition, error)
	FindSyncPrefix(uid int64, prefix string) ([]*entity.Script, error)
	FindSyncScript(page *request.Pages) ([]*entity.Script, error)
	HotKeyword() ([]redis.Z, error)
}

type script struct {
	scriptRepo   repository.Script
	codeRepo     repository.ScriptCode
	categoryRepo repository.Category
	statisRepo   repository.Statistics
	bgCategory   *entity.ScriptCategoryList
	cronCategory *entity.ScriptCategoryList
}

func NewScript(scriptRepo repository.Script, codeRepo repository.ScriptCode, categoryRepo repository.Category, statisRepo repository.Statistics) Script {
	ret := &script{
		scriptRepo:   scriptRepo,
		codeRepo:     codeRepo,
		categoryRepo: categoryRepo,
		statisRepo:   statisRepo,
	}
	ret.bgCategory = &entity.ScriptCategoryList{
		Name:       "后台脚本",
		Createtime: time.Now().Unix(),
	}
	ret.cronCategory = &entity.ScriptCategoryList{
		Name:       "定时脚本",
		Createtime: time.Now().Unix(),
	}
	if err := categoryRepo.Save(ret.bgCategory); err != nil {
		panic(err)
	}
	if err := categoryRepo.Save(ret.cronCategory); err != nil {
		panic(err)
	}

	return ret
}

func (s *script) Search(search *repository.SearchList, page *request.Pages) ([]*entity.Script, int64, error) {
	return s.scriptRepo.List(search, page)
}

func (s *script) UserScript(uid int64, self bool, page *request.Pages) ([]*entity.Script, int64, error) {
	return s.scriptRepo.List(&repository.SearchList{
		Uid:    uid,
		Self:   self,
		Status: cnt.ACTIVE,
	}, page)
}

func (s *script) Info(id int64) (*entity.Script, error) {
	script, err := s.scriptRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, errs.ErrScriptNotFound
	}
	switch script.Status {
	case cnt.ACTIVE:
		return script, nil
	case cnt.AUDIT:
		return nil, errs.ErrScriptAudit
	}
	return nil, errs.ErrScriptNotFound
}

func (s *script) VersionList(id int64, page *request.Pages) ([]*entity.ScriptCode, int64, error) {
	return s.codeRepo.List(id, cnt.ACTIVE, page)
}

func (s *script) GetCategory() ([]*entity.ScriptCategoryList, error) {
	return s.categoryRepo.List()
}

func (s *script) Download(id int64) error {
	return s.statisRepo.Download(id)
}

func (s *script) Update(id int64) error {
	return s.statisRepo.Update(id)
}

const (
	AUTO_SYNC_MODE = 1
	NONE_SYNC_MODE = 2
)

func (s *script) CreateScript(uid int64, req *request.CreateScript) (*entity.Script, error) {
	script := &entity.Script{
		UserId:     uid,
		Content:    req.Content,
		Type:       req.Type,
		Public:     req.Public,
		Unwell:     req.Unwell,
		Status:     cnt.ACTIVE,
		SyncMode:   NONE_SYNC_MODE,
		Createtime: time.Now().Unix(),
	}
	if err := s.createScriptCode(uid, script, req); err != nil {
		return nil, err
	}
	return script, broker.PublishEventScriptCreate(script.ID)
}

func (s *script) UpdateScript(uid, id int64, req *request.UpdateScript) error {
	script, err := s.Info(id)
	if err != nil {
		return err
	}
	if script.UserId != uid {
		return errs.ErrScriptForbidden
	}
	//script.Public = req.Public
	//script.Unwell = req.Unwell
	script.SyncUrl = req.SyncUrl
	script.ContentUrl = req.ContentUrl
	script.SyncMode = req.SyncMode
	switch script.Type {
	case entity.USERSCRIPT_TYPE, entity.SUBSCRIBE_TYPE:
	case entity.LIBRARY_TYPE:
		script.Name = req.Name
		script.Description = req.Description
		script.DefinitionUrl = req.DefinitionUrl
	default:
		return errors.New("错误的脚本类型")
	}
	return s.scriptRepo.Save(script)
}

func (s *script) CreateScriptCode(uid, id int64, req *request.UpdateScriptCode) error {
	script, err := s.Info(id)
	if err != nil {
		return err
	}
	if script.UserId != uid {
		return errs.ErrScriptForbidden
	}
	//if script.Type == entity.LIBRARY_TYPE {
	//	script.Name = req.Name
	//	script.Description = req.Description
	//}
	if err := s.createScriptCode(uid, script, &request.CreateScript{
		Content:     req.Content,
		Code:        req.Code,
		Name:        script.Name,
		Description: script.Description,
		Definition:  req.Definition,
		Type:        script.Type,
		Public:      req.Public,
		Unwell:      req.Unwell,
		Changelog:   req.Changelog,
	}); err != nil {
		return err
	}
	return nil
}

func (s *script) createScriptCode(uid int64, script *entity.Script, req *request.CreateScript) error {
	script.Content = req.Content
	script.Public = req.Public
	script.Unwell = req.Unwell
	script.Updatetime = time.Now().Unix()
	code := &entity.ScriptCode{
		UserId:     uid,
		Version:    time.Now().Format("20060102150405"),
		Changelog:  req.Changelog,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	switch req.Type {
	case entity.USERSCRIPT_TYPE, entity.SUBSCRIBE_TYPE:
		ncode := ""
		ncode, meta, scriptType, err := utils.GetCodeMeta(req.Code)
		if err != nil {
			return errs.NewBadRequestError(1000, err.Error())
		}
		req.Code = ncode
		if req.Type != scriptType {
			return errs.NewBadRequestError(1001, "脚本类型与脚本代码内容不等")
		}
		metaJson := utils.ParseMetaToJson(meta)
		if _, ok := metaJson["name"]; !ok {
			return errs.NewBadRequestError(1002, "脚本`name`不能为空")
		}
		if _, ok := metaJson["description"]; !ok {
			return errs.NewBadRequestError(1003, "脚本`description`不能为空")
		}
		script.Name = metaJson["name"][0]
		script.Description = metaJson["description"][0]
		metaJsonStr, err := json.Marshal(metaJson)
		if err != nil {
			return err
		}
		oldVersion := code.Version
		version := code.Version
		if v, ok := metaJson["version"]; ok {
			version = v[0]
		}
		code.Code = req.Code
		code.Meta = meta
		code.MetaJson = string(metaJsonStr)
		code.Version = version
		if script.ID != 0 {
			if ok, err := s.codeRepo.FindByVersion(script.ID, code.Version); err != nil {
				return err
			} else if ok != nil {
				//NOTE: 是不是应该优化一下对比功能
				if strings.ReplaceAll(ok.Code, "\r\n", "\n") != strings.ReplaceAll(code.Code, "\r\n", "\n") {
					return errs.ErrScriptCodeExist
				}
				code = ok
				oldVersion = code.Version
				script.Updatetime = ok.Createtime
				if req.Changelog != "" {
					code.Changelog = req.Changelog
					code.Updatetime = time.Now().Unix()
				}
			}
		}
		if err := db.Db.Transaction(func(tx *gorm.DB) error {
			scriptRepo := repository.NewTxScript(tx)
			if err := scriptRepo.Save(script); err != nil {
				return err
			}
			codeRepo := repository.NewTxCode(tx)
			code.ScriptId = script.ID
			if err := codeRepo.Save(code); err != nil {
				return err
			}
			categoryRepo := repository.NewTxCategory(tx)
			domains := make(map[string]struct{})
			if _, ok := metaJson["background"]; ok {
				_ = categoryRepo.LinkCategory(code.ScriptId, s.bgCategory.ID)
			}
			if _, ok := metaJson["crontab"]; ok {
				_ = categoryRepo.LinkCategory(code.ScriptId, s.bgCategory.ID)
				_ = categoryRepo.LinkCategory(code.ScriptId, s.cronCategory.ID)
			}
			for _, u := range metaJson["match"] {
				domain := utils.ParseMetaDomain(u)
				if domain != "" {
					domains[domain] = struct{}{}
				} else {
					glog.Warningf("deal meta url info: %d %s", code.ID, u)
				}
			}
			for _, u := range metaJson["include"] {
				domain := utils.ParseMetaDomain(u)
				if domain != "" {
					domains[domain] = struct{}{}
				} else {
					glog.Warningf("deal meta url info: %d %s", code.ID, u)
				}
			}
			for domain := range domains {
				if ret, err := codeRepo.FindScriptDomain(code.ScriptId, domain); err != nil {
					return err
				} else if ret != nil {
					continue
				}
				if err := codeRepo.SaveScriptDomain(&entity.ScriptDomain{
					Domain:        domain,
					DomainReverse: utils.StringReverse(domain),
					ScriptId:      code.ScriptId,
					ScriptCodeId:  code.ID,
					Createtime:    time.Now().Unix(),
				}); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		if version != oldVersion {
			_ = broker.PublishEventScriptVersionUpdate(script.ID, code.ID)
		}
	case entity.LIBRARY_TYPE:
		// 库的处理
		if req.Name == "" {
			return errs.NewBadRequestError(1004, "库的名称不能为空")
		}
		if req.Description == "" {
			return errs.NewBadRequestError(1005, "库的描述不能为空")
		}
		script.Name = req.Name
		script.Description = req.Description
		code.Code = req.Code
		if err := db.Db.Transaction(func(tx *gorm.DB) error {
			scriptRepo := repository.NewTxScript(tx)
			if err := scriptRepo.Save(script); err != nil {
				return err
			}
			codeRepo := repository.NewTxCode(tx)
			code.ScriptId = script.ID
			if err := codeRepo.Save(code); err != nil {
				return err
			}
			if req.Definition != "" {
				definition := &entity.LibDefinition{
					UserId:     uid,
					ScriptId:   script.ID,
					CodeId:     code.ID,
					Definition: req.Definition,
					Createtime: time.Now().Unix(),
				}
				if err := codeRepo.SaveDefinition(definition); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		_ = broker.PublishEventScriptVersionUpdate(script.ID, code.ID)
	default:
		return errs.NewBadRequestError(1010, "错误的类型")
	}
	return nil
}

func (s *script) GetCodeDefinition(codeid int64) (*entity.LibDefinition, error) {
	ret, err := s.codeRepo.FindDefinitionByCodeId(codeid)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, errs.ErrCodeDefinitionNotFound
	}
	return ret, nil
}

func (s *script) FindSyncPrefix(uid int64, prefix string) ([]*entity.Script, error) {
	return s.scriptRepo.FindSyncPrefix(uid, prefix)
}

func (s *script) FindSyncScript(page *request.Pages) ([]*entity.Script, error) {
	return s.scriptRepo.FindSyncScript(page)
}

func (s *script) GetCode(id int64) (*entity.ScriptCode, error) {
	return s.codeRepo.Find(id)
}

func (s *script) GetScriptVersion(scriptId int64, version string) (*entity.ScriptCode, error) {
	return s.codeRepo.FindByVersion(scriptId, version)
}

func (s *script) GetLatestVersion(scriptId int64) (*entity.ScriptCode, error) {
	codes, _, err := s.codeRepo.List(scriptId, cnt.ACTIVE, &request.Pages{})
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return nil, errs.ErrScriptAudit
	}
	return codes[0], nil
}

func (s *script) HotKeyword() ([]redis.Z, error) {
	return s.scriptRepo.HotKeyword()
}
