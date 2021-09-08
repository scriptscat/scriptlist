package entity

const (
	USERSCRIPT_TYPE = iota + 1
	SUBSCRIBE_TYPE
	LIBRARY_TYPE
)

const (
	PUBLIC_SCRIPT = iota + 1
	PRIVATE_SCRIPT
)

type Script struct {
	ID          int64  `gorm:"column:id" json:"id" form:"id"`
	PostId      int64  `gorm:"column:post_id;index:post_id,unique" json:"post_id" form:"post_id"`
	UserId      int64  `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	Name        string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Description string `gorm:"column:description;type:text" json:"description" form:"description"`
	Content     string `gorm:"column:content" json:"content" form:"content"`
	Type        int    `gorm:"column:type;index:type;not null;default:1"`
	Public      int    `gorm:"column:public;not null;default:1"`
	Unwell      int    `gorm:"column:unwell;not null;default:2"`
	SyncUrl     string `gorm:"column:sync_url;"`
	ContentUrl  string `gorm:"column:content_url;"`
	SyncMode    int    `gorm:"column:sync_mode;"`
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
	Updatetime int64  `gorm:"column:updatetime;type:bigint" json:"updatetime" form:"updatetime"`
}

func (s *ScriptCode) TableName() string {
	return "cdb_tampermonkey_script_code"
}

type LibDefinition struct {
	ID         int64  `gorm:"column:id" json:"id" form:"id"`
	UserId     int64  `gorm:"column:user_id;index:user_id;not null" json:"user_id" form:"user_id"`
	ScriptId   int64  `gorm:"column:script_id;index:script_id;not null" json:"script_id" form:"script_id"`
	CodeId     int64  `gorm:"column:code_id;index:code_id;not null" json:"code_id" form:"code_id"`
	Definition string `gorm:"column:definition;not null" json:"definition" form:"definition"`
	Createtime int64  `gorm:"column:createtime;type:bigint" json:"createtime" form:"createtime"`
}
