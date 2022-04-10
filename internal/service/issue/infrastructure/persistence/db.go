package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/issue/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/issue/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/gorm"
)

type IssueRepositories struct {
	db *gorm.DB
	repository.Issue
	repository.IssueComment
	repository.IssueWatch
}

func NewRepositories(db *gorm.DB, redis *redis.Client) *IssueRepositories {
	return &IssueRepositories{
		db:           db,
		Issue:        NewIssue(db),
		IssueComment: NewComment(db),
		IssueWatch:   NewWatch(redis),
	}
}

func (r *IssueRepositories) AutoMigrate() error {
	return utils.ErrFunc(
		func() error {
			return r.db.AutoMigrate(&entity.ScriptIssue{})
		},
		func() error {
			return r.db.AutoMigrate(&entity.ScriptIssueComment{})
		},
	)
}
