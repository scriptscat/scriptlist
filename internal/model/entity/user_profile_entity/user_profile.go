package user_profile_entity

// UserProfile 用户资料表
type UserProfile struct {
	ID          int64  `gorm:"column:id;type:bigint;not null;primary_key"`
	Username    string `gorm:"column:username;type:varchar(128);not null;index:uni_cm_user_profile_username,unique"`
	Nickname    string `gorm:"column:nickname;type:varchar(128);not null"`
	Description string `gorm:"column:description;type:varchar(512);not null"` // 个人简介
	Avatar      string `gorm:"column:avatar;type:varchar(255);not null"`
	Location    string `gorm:"column:location;type:varchar(128);not null"`    // 用户所在地
	Website     string `gorm:"column:website;type:varchar(255);not null"`     // 个人网站
	Email       string `gorm:"column:email;type:varchar(128);not null;"`      // 联系邮箱
	Status      int32  `gorm:"column:status;type:tinyint;default:1;not null"` // 用户状态 1:正常 2:封禁
	Createtime  int64  `gorm:"column:createtime;type:bigint;not null"`
	Updatetime  int64  `gorm:"column:updatetime;type:bigint;not null"`
}
