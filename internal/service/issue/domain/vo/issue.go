package vo

import (
	"strings"

	"github.com/scriptscat/scriptlist/internal/service/issue/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/vo"
)

type Issue struct {
	*vo.User
	ID         int64    `json:"id"`
	ScriptID   int64    `json:"script_id"`
	UserID     int64    `json:"user_id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Labels     []string `json:"labels"`
	Status     int      `json:"status"`
	Createtime int64    `json:"createtime"`
	Updatetime int64    `json:"updatetime"`
}

func ToIssue(user *vo.User, issue *entity.ScriptIssue) *Issue {
	var labels []string
	if issue.Labels != "" {
		labels = strings.Split(issue.Labels, ",")
	}
	return &Issue{
		User:       user,
		ID:         issue.ID,
		ScriptID:   issue.ScriptID,
		UserID:     issue.UserID,
		Title:      issue.Title,
		Content:    issue.Content,
		Labels:     labels,
		Status:     issue.Status,
		Createtime: issue.Createtime,
		Updatetime: issue.Updatetime,
	}
}

type IssueComment struct {
	*vo.User
	ID         int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	IssueID    int64  `gorm:"column:issue_id;type:bigint(20);index:issue_id;NOT NULL" json:"issue_id"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"user_id"`
	Content    string `gorm:"column:content;type:text;NOT NULL" json:"content"`
	Type       int    `gorm:"column:type;type:tinyint(4);default:0;NOT NULL" json:"type"`
	Status     int    `gorm:"column:status;type:tinyint(4);default:0;NOT NULL" json:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)" json:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)" json:"updatetime"`
}

func ToIssueComment(user *vo.User, issue *entity.ScriptIssueComment) *IssueComment {
	return &IssueComment{
		User:       user,
		ID:         issue.ID,
		IssueID:    issue.IssueID,
		UserID:     issue.UserID,
		Content:    issue.Content,
		Type:       issue.Type,
		Status:     issue.Status,
		Createtime: issue.Createtime,
		Updatetime: issue.Updatetime,
	}
}

type IssueWatch struct {
	*vo.User
}

func ToIssueWatch(user *vo.User) *IssueWatch {
	return &IssueWatch{
		User: user,
	}
}
