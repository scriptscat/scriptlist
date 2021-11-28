package repository

import (
	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
)

type SearchList struct {
	Uid                   int64
	Self                  bool
	Category              []int64
	Status                int64
	Keyword, Sort, Domain string
}

type Script interface {
	Find(id int64) (*entity.Script, error)
	Save(script *entity.Script) error
	List(search *SearchList, page request.Pages) ([]*entity.Script, int64, error)
	FindSyncPrefix(uid int64, prefix string) ([]*entity.Script, error)
	FindSyncScript(page request.Pages) ([]*entity.Script, error)
}

type ScriptCode interface {
	Find(id int64) (*entity.ScriptCode, error)
	Save(script *entity.ScriptCode) error
	FindByVersion(scriptId int64, version string) (*entity.ScriptCode, error)
	List(script, status int64) ([]*entity.ScriptCode, error)
	SaveDefinition(definition *entity.LibDefinition) error
	SaveScriptDomain(domain *entity.ScriptDomain) error
	FindScriptDomain(scriptId int64, domain string) (*entity.ScriptDomain, error)
	FindDefinitionByCodeId(codeid int64) (*entity.LibDefinition, error)
}

type Score interface {
	Save(score *entity.ScriptScore) error
	UserScore(uid, scriptId int64) (*entity.ScriptScore, error)
	Avg(scriptId int64) (int64, error)
	Count(scriptId int64) (int64, error)
	List(scriptId int64, page *request.Pages) ([]*entity.ScriptScore, int64, error)
}

type Category interface {
	List() ([]*entity.ScriptCategoryList, error)
	LinkCategory(script, category int64) error
	Save(category *entity.ScriptCategoryList) error
}

type Statistics interface {
	Download(id int64) error
	Update(id int64) error
}

type Watch struct {
	UserId int64 `json:"user_id"`
	// Watch级别,0:未监听 1:版本更新监听 2:新建issue监听 2:评论都监听
	Level int `json:"level"`
}

type ScriptWatch interface {
	List(script int64) ([]*Watch, error)
	Num(script int64) (int, error)
	Watch(script, user int64, level int) error
	Unwatch(script, user int64) error
	IsWatch(script, user int64) (int, error)
}
