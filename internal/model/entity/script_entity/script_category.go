package script_entity

type ScriptCategory struct {
	ID         int64 `gorm:"column:id;type:bigint(20);not null;primary_key"`
	CategoryID int64 `gorm:"column:category_id;type:bigint(20);index:script_category,unique;index:category_id"`
	ScriptID   int64 `gorm:"column:script_id;type:bigint(20);index:script_category,unique;index:script_id"`
	Createtime int64 `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64 `gorm:"column:updatetime;type:bigint(20)"`
}

type ScriptCategoryList struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Name       string `gorm:"column:name;type:varchar(255)"`
	Num        int64  `gorm:"column:num;type:bigint(20)"`
	Sort       int64  `gorm:"column:sort;type:int(10);index:category_sort"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}
