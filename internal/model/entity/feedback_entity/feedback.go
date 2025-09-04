package feedback_entity

type Feedback struct {
	ID         int64  `gorm:"column:id;type:bigint;not null;primary_key"`
	Reason     string `gorm:"column:reason;type:varchar(255);not null"`
	Content    string `gorm:"column:Content;type:text"`
	ClientIp   string `gorm:"column:client_ip;type:varchar(255);not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint"`
}
