package entity

import "gorm.io/datatypes"

type User struct {
	Uid                int64  `gorm:"column:uid" json:"uid" form:"uid"`
	Email              string `gorm:"column:email" json:"email" form:"email"`
	Username           string `gorm:"column:username" json:"username" form:"username"`
	Password           string `gorm:"column:password" json:"password" form:"password"`
	Status             int64  `gorm:"column:status" json:"status" form:"status"`
	Emailstatus        int64  `gorm:"column:emailstatus" json:"emailstatus" form:"emailstatus"`
	Avatarstatus       int64  `gorm:"column:avatarstatus" json:"avatarstatus" form:"avatarstatus"`
	Videophotostatus   int64  `gorm:"column:videophotostatus" json:"videophotostatus" form:"videophotostatus"`
	Adminid            int64  `gorm:"column:adminid" json:"adminid" form:"adminid"`
	Groupid            int64  `gorm:"column:groupid" json:"groupid" form:"groupid"`
	Groupexpiry        int64  `gorm:"column:groupexpiry" json:"groupexpiry" form:"groupexpiry"`
	Extgroupids        string `gorm:"column:extgroupids" json:"extgroupids" form:"extgroupids"`
	Regdate            int64  `gorm:"column:regdate" json:"regdate" form:"regdate"`
	Credits            int64  `gorm:"column:credits" json:"credits" form:"credits"`
	Notifysound        int64  `gorm:"column:notifysound" json:"notifysound" form:"notifysound"`
	Timeoffset         string `gorm:"column:timeoffset" json:"timeoffset" form:"timeoffset"`
	Newpm              int64  `gorm:"column:newpm" json:"newpm" form:"newpm"`
	Newprompt          int64  `gorm:"column:newprompt" json:"newprompt" form:"newprompt"`
	Accessmasks        int64  `gorm:"column:accessmasks" json:"accessmasks" form:"accessmasks"`
	Allowadmincp       int64  `gorm:"column:allowadmincp" json:"allowadmincp" form:"allowadmincp"`
	Onlyacceptfriendpm int64  `gorm:"column:onlyacceptfriendpm" json:"onlyacceptfriendpm" form:"onlyacceptfriendpm"`
	Conisbind          int64  `gorm:"column:conisbind" json:"conisbind" form:"conisbind"`
	Freeze             int64  `gorm:"column:freeze" json:"freeze" form:"freeze"`
}

type UserArchive User

func (u *User) TableName() string {
	return "pre_common_member"
}

func (u *UserArchive) TableName() string {
	return "pre_common_member"
}

type UserConfig struct {
	ID         int64             `gorm:"column:id" json:"id" form:"id"`
	Uid        int64             `gorm:"column:uid;index:user_id,unique" json:"uid" form:"uid"`
	Notify     datatypes.JSONMap `gorm:"column:notify" json:"notify" form:"uid"`
	Createtime int64             `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64             `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}
