package script_entity

type ScriptScore struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);index:user_id,unique;index:user_script,unique;index:user"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);index:user_id,unique;index:user_script,unique;index:script_id;index:script"`
	Score      int64  `gorm:"column:score;type:double"`
	Message    string `gorm:"column:message;type:longtext"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
	State      int64  `gorm:"column:state;type:int(10);default:1"`
}
