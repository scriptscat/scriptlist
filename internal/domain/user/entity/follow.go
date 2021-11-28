package entity

type HomeFollow struct {
	Uid       uint   `gorm:"column:uid" json:"uid"`
	Username  string `gorm:"column:username;type:char(15);NOT NULL" json:"username"`
	Followuid uint   `gorm:"column:followuid;type:mediumint(8) unsigned;default:0;NOT NULL" json:"followuid"`
	Fusername string `gorm:"column:fusername;type:char(15);NOT NULL" json:"fusername"`
	Bkname    string `gorm:"column:bkname;type:varchar(255);NOT NULL" json:"bkname"`
	Status    int    `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`
	Mutual    int    `gorm:"column:mutual;type:tinyint(1);default:0;NOT NULL" json:"mutual"`
	Dateline  uint   `gorm:"column:dateline;type:int(10) unsigned;default:0;NOT NULL" json:"dateline"`
}

func (h *HomeFollow) TableName() string {
	return "pre_home_follow"
}
