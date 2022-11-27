package entity

import "gorm.io/datatypes"

type UserConfig struct {
	ID         int64             `gorm:"column:id" json:"id" form:"id"`
	Uid        int64             `gorm:"column:uid;index:user_id,unique" json:"uid" form:"uid"`
	Notify     datatypes.JSONMap `gorm:"column:notify" json:"notify" form:"uid"`
	Createtime int64             `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64             `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}
