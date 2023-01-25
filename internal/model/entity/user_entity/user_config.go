package user_entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/codfrm/cago/pkg/utils"
)

type Notify struct {
	// 创建脚本
	CreateScript *bool `json:"create_script"`
	// 脚本更新
	ScriptUpdate *bool `json:"script_update"`
	// 脚本反馈
	ScriptIssue *bool `json:"script_issue"`
	// 脚本反馈评论
	ScriptIssueComment *bool `json:"script_issue_comment"`
	// 脚本评分
	Score *bool `json:"score"`
	// 艾特
	At *bool `json:"at"`
}

func (n *Notify) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, n)
	return err
}

func (n *Notify) Value() (driver.Value, error) {
	n.DefaultValue()
	return json.Marshal(n)
}

func (n *Notify) DefaultValue() {
	setTrue(n.ScriptIssue)
	setTrue(n.ScriptIssueComment)
	setTrue(n.ScriptUpdate)
	setTrue(n.At)
	setTrue(n.Score)
	setTrue(n.CreateScript)
}

func setTrue(b *bool) {
	t := true
	if b == nil {
		b = &t
	}
}

type UserConfig struct {
	ID  int64 `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Uid int64 `gorm:"column:uid;type:bigint(20);index:user_id"`
	// Webhook token
	Token      string  `gorm:"column:token;type:varchar(255);null;index:token,unique"`
	Notify     *Notify `gorm:"column:notify;type:json"`
	Createtime int64   `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64   `gorm:"column:updatetime;type:bigint(20)"`
}

func (u *UserConfig) GenToken() string {
	u.Token = utils.RandString(64, utils.Letter)
	return u.Token
}
