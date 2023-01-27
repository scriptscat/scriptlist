package issue_entity

type ScriptIssueWatch struct {
	ID         int64 `gorm:"column:id;type:bigint(20);not null;primary_key"`
	IssueID    int64 `gorm:"column:issue_id;type:bigint(20);not null;index:user_issue,unique;index:issue"`
	UserID     int64 `gorm:"column:user_id;type:bigint(20);not null;index:user_issue,unique"`
	Status     int32 `gorm:"column:status;type:tinyint(4);default:1;not null"`
	Createtime int64 `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime int64 `gorm:"column:updatetime;type:bigint(20)"`
}
