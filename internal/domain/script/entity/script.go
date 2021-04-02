package entity

type Script struct {
	ID          int64  `gorm:"column:id" json:"id" form:"id"`
	PostId      int64  `gorm:"column:post_id;index:post_id,unique" json:"post_id" form:"post_id"`
	UserId      int64  `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	Name        string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Description string `gorm:"column:description;type:text" json:"description" form:"description"`
	Content     string `gorm:"column:content" json:"content" form:"content"`
	Status      int64  `gorm:"column:status" json:"status" form:"status"`
	Createtime  int64  `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime  int64  `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

func (s *Script) TableName() string {
	return "cdb_tampermonkey_script"
}

type ScriptCode struct {
	ID         int64  `gorm:"column:id" json:"id" form:"id"`
	UserId     int64  `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	ScriptId   int64  `gorm:"column:script_id;index:script_id" json:"script_id" form:"script_id"`
	Code       string `gorm:"column:code" json:"code" form:"code"`
	Meta       string `gorm:"column:meta" json:"meta" form:"meta"`
	MetaJson   string `gorm:"column:meta_json" json:"meta_json" form:"meta_json"`
	Version    string `gorm:"column:version;type:varchar(255)" json:"version" form:"version"`
	Changelog  string `gorm:"column:changelog;type:text" json:"changelog" form:"changelog"`
	Status     int64  `gorm:"column:status" json:"status" form:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint" json:"createtime" form:"createtime"`
}

func (s *ScriptCode) TableName() string {
	return "cdb_tampermonkey_script_code"
}
