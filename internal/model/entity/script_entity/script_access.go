package script_entity

type ScriptAccess struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	LinkID     int64  `gorm:"column:link_id;type:bigint(20);not null"`             // 用户id或者用户组id
	Type       int32  `gorm:"column:type;type:tinyint(4);not null"`                // 1: 用户 2: 用户组
	Role       string `gorm:"column:access_permission;type:varchar(255);not null"` // 角色 访客: guest, 管理员: manager, 拥有者: owner
	Expiretime int64  `gorm:"column:expiretime;type:bigint(20)"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}
