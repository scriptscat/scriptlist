package persistence

import (
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/service/issue/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/issue/domain/repository"
	"gorm.io/gorm"
)

type comment struct {
	db *gorm.DB
}

func NewComment(db *gorm.DB) repository.IssueComment {
	return &comment{db: db}
}

func (c *comment) List(issue int64, status int, page *request.Pages) ([]*entity.ScriptIssueComment, error) {
	list := make([]*entity.ScriptIssueComment, 0)
	find := c.db.Model(&entity.ScriptIssueComment{}).Where("issue_id=? and status=?", issue, status).
		Limit(page.Size()).Offset((page.Page() - 1) * page.Size())
	if err := find.Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (c *comment) FindById(comment int64) (*entity.ScriptIssueComment, error) {
	ret := &entity.ScriptIssueComment{}
	if err := c.db.First(ret, "id=?", comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (c *comment) Save(comment *entity.ScriptIssueComment) error {
	return c.db.Save(comment).Error
}
