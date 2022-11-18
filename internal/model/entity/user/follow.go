package user

type HomeFollow struct {
	// 关注用户
	Uid int64 `gorm:"column:uid" json:"uid"`
	// 关注用户名
	Username string `gorm:"column:username;type:char(15);NOT NULL" json:"username"`
	// 被关注用户
	Followuid int64 `gorm:"column:followuid;type:mediumint(8) unsigned;default:0;NOT NULL" json:"followuid"`
	// 被关注用户名
	Fusername string `gorm:"column:fusername;type:char(15);NOT NULL" json:"fusername"`
	// 备注名
	Bkname string `gorm:"column:bkname;type:varchar(255);NOT NULL" json:"bkname"`
	// 0正常 1特殊关注 -1不能再关注此人
	Status int `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`
	// 0单向 1互相
	Mutual int `gorm:"column:mutual;type:tinyint(1);default:0;NOT NULL" json:"mutual"`
	// 关注时间
	Dateline int64 `gorm:"column:dateline;type:int(10) unsigned;default:0;NOT NULL" json:"dateline"`
}

func (h *HomeFollow) TableName() string {
	return "pre_home_follow"
}
