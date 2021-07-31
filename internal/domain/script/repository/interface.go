package repository

import (
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/request"
)

type SearchList struct {
	Uid                   int64
	Category              []int64
	Status                int64
	Keyword, Sort, Domain string
}

type Script interface {
	Find(id int64) (*entity.Script, error)
	Save(script *entity.Script) error
	List(search *SearchList, page request.Pages) ([]*entity.Script, int64, error)
}

type ScriptCode interface {
	Find(id int64) (*entity.ScriptCode, error)
	Save(script *entity.ScriptCode) error
	List(script, status int64) ([]*entity.ScriptCode, error)
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
