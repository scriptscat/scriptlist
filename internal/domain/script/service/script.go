package service

import (
	"encoding/json"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/domain/script/repository"
	"github.com/scriptscat/scriptweb/internal/http/dto/request"
	"github.com/scriptscat/scriptweb/internal/pkg/cnt"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	"github.com/scriptscat/scriptweb/migrations"
	"github.com/scriptscat/scriptweb/pkg/utils"
	"gorm.io/gorm"
)

type Script interface {
	Search(search *repository.SearchList, page request.Pages) ([]*entity.Script, int64, error)
	UserScript(uid int64, self bool, page request.Pages) ([]*entity.Script, int64, error)
	Info(id int64) (*entity.Script, error)
	VersionList(id int64) ([]*entity.ScriptCode, error)
	GetCategory() ([]*entity.ScriptCategoryList, error)
	Download(id int64) error
	Update(id int64) error
	CreateScript(uid int64, req *request.CreateScript) error
}

type script struct {
	scriptRepo   repository.Script
	codeRepo     repository.ScriptCode
	categoryRepo repository.Category
	statisRepo   repository.Statistics
}

func NewScript(scriptRepo repository.Script, codeRepo repository.ScriptCode, categoryRepo repository.Category, statisRepo repository.Statistics, c *cron.Cron) Script {
	go migrations.DealMetaInfo()
	c.AddFunc("0/20 * * * *", func() {
		migrations.DealMetaInfo()
	})
	ret := &script{
		scriptRepo:   scriptRepo,
		codeRepo:     codeRepo,
		categoryRepo: categoryRepo,
		statisRepo:   statisRepo,
	}
	return ret
}

func (s *script) Search(search *repository.SearchList, page request.Pages) ([]*entity.Script, int64, error) {
	return s.scriptRepo.List(search, page)
}

func (s *script) UserScript(uid int64, self bool, page request.Pages) ([]*entity.Script, int64, error) {
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
	switch script.Status {
	case 1:
		return script, nil
	case 2:
		return nil, errs.ErrScriptAudit
	}
	return nil, errs.ErrScriptNotFound
}

func (s *script) VersionList(id int64) ([]*entity.ScriptCode, error) {
	return s.codeRepo.List(id, cnt.ACTIVE)
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

func (s *script) CreateScript(uid int64, req *request.CreateScript) error {
	script := &entity.Script{
		PostId:     0,
		UserId:     uid,
		Content:    req.Content,
		Type:       req.Type,
		Public:     req.Public,
		Unwell:     req.Unwell,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	code := &entity.ScriptCode{
		UserId:     uid,
		Version:    time.Now().Format("20060102150405"),
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
		version := ""
		if v, ok := metaJson["version"]; ok {
			version = v[0]
		}
		code.Code = req.Code
		code.Meta = meta
		code.MetaJson = string(metaJsonStr)
		code.Version = version
		if err := db.Db.Transaction(func(tx *gorm.DB) error {
			scriptRepo := repository.NewTxScript(tx)
			if err := scriptRepo.Save(script); err != nil {
				return err
			}
			codeRepo := repository.NewTxCode(tx)
			code.ScriptId = script.ID
			return codeRepo.Save(code)
		}); err != nil {
			return err
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
	default:
		return errs.NewBadRequestError(1010, "错误的类型")
	}
	return nil
}
