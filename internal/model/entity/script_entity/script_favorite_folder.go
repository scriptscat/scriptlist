package script_entity

type ScriptFavoriteFolder struct {
	ID          int64  `gorm:"column:id;type:bigint;not null;primary_key"`
	Name        string `gorm:"column:name;type:varchar(50);not null"`
	Description string `gorm:"column:description;type:varchar(200)"`
	UserID      int64  `gorm:"column:user_id;type:bigint;not null;index:idx_user_id"`
	Private     int32  `gorm:"column:private;type:tinyint;not null;default:1"` // 1私密 2公开
	Count       int64  `gorm:"column:count;type:bigint;default:0"`
	Sort        int64  `gorm:"column:sort;type:bigint;default:0"`
	Status      int32  `gorm:"column:status;type:tinyint;not null"`
	Createtime  int64  `gorm:"column:createtime;type:bigint"`
	Updatetime  int64  `gorm:"column:updatetime;type:bigint"`
}
