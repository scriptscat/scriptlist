package script

type Type int

const (
	UserscriptType Type = iota + 1 // 用户脚本
	SubscribeType                  // 订阅脚本
	LibraryType                    // 库
)

type Public int

const (
	PublicScript  Public = iota + 1 // 公开
	PrivateScript                   // 私有(半公开,只是不展示在列表中)
)

type UnwellContent int

const (
	Unwell UnwellContent = iota + 1 // 不适内容
	Well                            // 合适内容
)

type SyncMode int

const (
	SyncModeAuto   SyncMode = iota + 1 // 自动同步
	SyncModeManual                     // 手动同步
)

type Script struct {
	ID            int64         `gorm:"column:id" json:"id" form:"id"`
	PostID        int64         `gorm:"column:post_id;index:post_id,unique" json:"post_id" form:"post_id"`
	UserID        int64         `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	Name          string        `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Description   string        `gorm:"column:description;type:text" json:"description" form:"description"`
	Content       string        `gorm:"column:content;type:mediumtext" json:"content" form:"content"`
	Type          Type          `gorm:"column:type;type:bigint;index:script_type;not null;default:1" json:"type"`
	Public        Public        `gorm:"column:public;not null;default:1" json:"public"`
	Unwell        UnwellContent `gorm:"column:unwell;not null;default:2" json:"unwell"` // 不适内容
	SyncUrl       string        `gorm:"column:sync_url;type:text;index:sync_url,length:128" json:"sync_url"`
	ContentUrl    string        `gorm:"column:content_url;type:text;index:content_url,length:128" json:"content_url"`
	DefinitionUrl string        `gorm:"column:definition_url;type:text;index:definition_url,length:128" json:"definition_url"`
	SyncMode      SyncMode      `gorm:"column:sync_mode;type:tinyint(2)"`
	Archive       int32         `gorm:"column:archive;type:tinyint(2)" json:"archive"` // 归档
	Status        int64         `gorm:"column:status" json:"status" form:"status"`
	Createtime    int64         `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime    int64         `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

func (s *Script) TableName() string {
	return "cdb_tampermonkey_script"
}

type Code struct {
	ID         int64  `gorm:"column:id" json:"id" form:"id"`
	UserId     int64  `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	ScriptId   int64  `gorm:"column:script_id;index:script_id" json:"script_id" form:"script_id"`
	Code       string `gorm:"column:code;type:mediumtext" json:"code" form:"code"`
	Meta       string `gorm:"column:meta;type:text" json:"meta" form:"meta"`
	MetaJson   string `gorm:"column:meta_json;type:text" json:"meta_json" form:"meta_json"`
	Version    string `gorm:"column:version;type:varchar(255)" json:"version" form:"version"`
	Changelog  string `gorm:"column:changelog;type:text" json:"changelog" form:"changelog"`
	Status     int64  `gorm:"column:status" json:"status" form:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint" json:"createtime" form:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint" json:"updatetime" form:"updatetime"`
}

func (s *Code) TableName() string {
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
