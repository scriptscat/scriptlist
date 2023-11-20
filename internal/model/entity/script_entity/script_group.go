package script_entity

// ScriptGroup 脚本组
type ScriptGroup struct {
	ID          int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID    int64  `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	Name        string `gorm:"column:name;type:varchar(255);not null"`
	Description string `gorm:"column:description;type:varchar(255);not null"`
	Status      int32  `gorm:"column:status;type:tinyint(4);not null"`
	Createtime  int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime  int64  `gorm:"column:updatetime;type:bigint(20)"`
}
