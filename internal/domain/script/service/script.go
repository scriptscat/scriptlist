package service

import (
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptweb/interfaces/dto/request"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/domain/script/repository"
	"github.com/scriptscat/scriptweb/internal/pkg/cnt"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	"github.com/scriptscat/scriptweb/internal/pkg/migrate"
)

type Script interface {
	Search(category []int64, domain, keyword, sort string, page request.Pages) ([]*entity.Script, int64, error)
	UserScript(uid int64, self bool, page request.Pages) ([]*entity.Script, int64, error)
	Info(id int64) (*entity.Script, error)
	VersionList(id int64) ([]*entity.ScriptCode, error)
	GetCategory() ([]*entity.ScriptCategoryList, error)
	Download(id int64) error
	Update(id int64) error
}

type script struct {
	scriptRepo   repository.Script
	codeRepo     repository.ScriptCode
	categoryRepo repository.Category
	statisRepo   repository.Statistics
}

func NewScript(scriptRepo repository.Script, codeRepo repository.ScriptCode, categoryRepo repository.Category, statisRepo repository.Statistics, c *cron.Cron) Script {
	go migrate.DealMetaInfo()
	c.AddFunc("0/20 * * * *", func() {
		migrate.DealMetaInfo()
	})
	ret := &script{
		scriptRepo:   scriptRepo,
		codeRepo:     codeRepo,
		categoryRepo: categoryRepo,
		statisRepo:   statisRepo,
	}
	return ret
}

func (s *script) Search(category []int64, domain, keyword, sort string, page request.Pages) ([]*entity.Script, int64, error) {
	return s.scriptRepo.List(&repository.SearchList{
		Category: category,
		Domain:   domain,
		Sort:     sort,
		Status:   cnt.ACTIVE,
		Keyword:  keyword,
	}, page)
}

func (s *script) UserScript(uid int64, self bool, page request.Pages) ([]*entity.Script, int64, error) {
	var status int64 = cnt.ACTIVE
	if self {
		status = cnt.UNKNOWN
	}
	return s.scriptRepo.List(&repository.SearchList{
		Uid:    uid,
		Status: status,
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
