package audit_entity

// Action 操作类型
type Action string

const (
	ActionScriptCreate Action = "script_create" // 创建脚本
	ActionScriptUpdate Action = "script_update" // 更新脚本
	ActionScriptDelete Action = "script_delete" // 删除脚本
)

// AuditLog 审计日志实体
type AuditLog struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key;autoIncrement"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);not null;index:idx_user_id"`
	Username   string `gorm:"column:username;type:varchar(255);not null"`
	Action     Action `gorm:"column:action;type:varchar(64);not null;index:idx_action"`
	TargetType string `gorm:"column:target_type;type:varchar(64);not null;index:idx_target"`
	TargetID   int64  `gorm:"column:target_id;type:bigint(20);not null;index:idx_target"`
	TargetName string `gorm:"column:target_name;type:varchar(255);default:''"`
	IsAdmin    bool   `gorm:"column:is_admin;type:tinyint(1);not null;default:0;index:idx_admin"`
	Reason     string `gorm:"column:reason;type:varchar(1024);default:''"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null;index:idx_createtime"`
}

// IsAdminAction 是否是管理员操作
func (a *AuditLog) IsAdminAction() bool {
	return a.IsAdmin
}

// IsDeleteAction 是否是删除操作
func (a *AuditLog) IsDeleteAction() bool {
	return a.Action == ActionScriptDelete
}
