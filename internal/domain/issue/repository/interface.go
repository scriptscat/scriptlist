package repository

import (
	"github.com/scriptscat/scriptlist/internal/domain/issue/entity"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
)

type Issue interface {
	List(scriptId int64, keyword string, tag []string, status int, page *request.Pages) ([]*entity.ScriptIssue, int64, error)
	FindById(issue int64) (*entity.ScriptIssue, error)
	Save(issue *entity.ScriptIssue) error
}

type IssueComment interface {
	List(issue int64, status int, page *request.Pages) ([]*entity.ScriptIssueComment, error)
	FindById(comment int64) (*entity.ScriptIssueComment, error)
	Save(comment *entity.ScriptIssueComment) error
}

type Watch struct {
	UserId int64 `json:"user_id"`
}

type IssueWatch interface {
	List(issue int64) ([]*Watch, error)
	Num(issue int64) (int, error)
	Watch(issue, user int64) error
	Unwatch(issue, user int64) error
	IsWatch(issue, user int64) (int, error)
}
