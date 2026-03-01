package report_entity

import (
	"context"
	"net/http"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// ReportReason 举报原因
type ReportReason struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ReasonMap 举报原因映射
var ReasonMap = map[string]*ReportReason{
	"malware":   {Key: "malware", Name: "恶意代码", Description: "脚本包含恶意代码"},
	"privacy":   {Key: "privacy", Name: "侵犯隐私", Description: "脚本侵犯用户隐私"},
	"copyright": {Key: "copyright", Name: "侵权/抄袭", Description: "脚本存在侵权或抄袭行为"},
	"spam":      {Key: "spam", Name: "垃圾信息/广告", Description: "脚本包含垃圾信息或广告"},
	"other":     {Key: "other", Name: "其他", Description: "其他原因"},
}

// ScriptReport 脚本举报
type ScriptReport struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key;autoIncrement"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);not null"`
	Reason     string `gorm:"column:reason;type:varchar(64);not null"`
	Content    string `gorm:"column:content;type:text"`
	Status     int32  `gorm:"column:status;type:tinyint(4);default:1;not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

// CheckOperate 检查是否可以操作
func (r *ScriptReport) CheckOperate(ctx context.Context) error {
	if r == nil {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ReportNotFound)
	}
	if r.Status == consts.DELETE {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ReportIsDelete)
	}
	return nil
}

// IsResolved 是否已解决
func (r *ScriptReport) IsResolved() bool {
	return r.Status == consts.AUDIT
}

// Resolve 解决举报
func (r *ScriptReport) Resolve(now time.Time) {
	r.Status = consts.AUDIT
	r.Updatetime = now.Unix()
}

// Reopen 重新打开举报
func (r *ScriptReport) Reopen(now time.Time) {
	r.Status = consts.ACTIVE
	r.Updatetime = now.Unix()
}

// ValidateReason 校验举报原因是否合法
func (r *ScriptReport) ValidateReason(ctx context.Context) error {
	if _, ok := ReasonMap[r.Reason]; !ok {
		return i18n.NewError(ctx, code.ReportReasonInvalid)
	}
	return nil
}
