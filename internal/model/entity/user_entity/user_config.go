package user_entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cago-frame/cago/pkg/utils"
)

// Notify 通知配置
// 0: 未设置 1: 开启 2: 关闭
type Notify struct {
	// 创建脚本
	CreateScript int `json:"create_script"`
	// 脚本更新
	ScriptUpdate int `json:"script_update"`
	// 脚本反馈
	ScriptIssue int `json:"script_issue"`
	// 脚本反馈评论
	ScriptIssueComment int `json:"script_issue_comment"`
	// 脚本评分
	Score int `json:"score"`
	// 艾特
	At int `json:"at"`
	// 脚本举报
	ScriptReport int `json:"script_report"`
	// 脚本举报评论
	ScriptReportComment int `json:"script_report_comment"`
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
	return json.Marshal(n)
}

// IsEnabled 判断通知是否开启 (0=默认开启, 1=开启, 2=关闭)
func (n *Notify) IsEnabled(v int) bool {
	return v != 2
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
