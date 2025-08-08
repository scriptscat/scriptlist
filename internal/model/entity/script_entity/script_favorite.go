package script_entity

type ScriptFavorite struct {
	ID               int64 `gorm:"column:id;type:bigint;not null;primary_key"`
	UserID           int64 `gorm:"column:user_id;type:bigint;not null;index:idx_user_id;index:idx_user_script_id"`
	ScriptID         int64 `gorm:"column:script_id;type:bigint;not null;index:script_id;index:idx_user_script_id"`
	FavoriteFolderID int64 `gorm:"column:favorite_folder_id;type:bigint;not null;index:favorite_folder_id;index:idx_favorite_folder_id"`
	Status           int32 `gorm:"column:status;type:tinyint;not null"` // 1正常 2删除
	Createtime       int64 `gorm:"column:createtime;type:bigint"`
	Updatetime       int64 `gorm:"column:updatetime;type:bigint"`
}
