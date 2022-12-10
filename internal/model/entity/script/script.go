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
	ID            int64         `gorm:"column:id;type:bigint(20);not null;primary_key"`
	PostID        int64         `gorm:"column:post_id;type:bigint(20);index:post_id,unique"`
	UserID        int64         `gorm:"column:user_id;type:bigint(20);index:user_id"`
	Name          string        `gorm:"column:name;type:varchar(255)"`
	Description   string        `gorm:"column:description;type:text"`
	Content       string        `gorm:"column:content;type:mediumtext"`
	Type          Type          `gorm:"column:type;type:bigint(20);default:1;not null;index:script_type"`
	Public        Public        `gorm:"column:public;type:bigint(20);default:1;not null"`
	Unwell        UnwellContent `gorm:"column:unwell;type:bigint(20);default:2;not null"`
	SyncUrl       string        `gorm:"column:sync_url;type:text;index:sync_url"`
	ContentUrl    string        `gorm:"column:content_url;type:text;index:content_url"`
	DefinitionUrl string        `gorm:"column:definition_url;type:text;index:definition_url"`
	SyncMode      SyncMode      `gorm:"column:sync_mode;type:tinyint(2)"`
	Archive       int32         `gorm:"column:archive;type:tinyint(2)"`
	Status        int64         `gorm:"column:status;type:bigint(20)"`
	Createtime    int64         `gorm:"column:createtime;type:bigint(20)"`
	Updatetime    int64         `gorm:"column:updatetime;type:bigint(20)"`
}

func (s *Script) TableName() string {
	return "cdb_tampermonkey_script"
}

type Code struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);index:user_id"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);index:script_id"`
	Code       string `gorm:"column:code;type:mediumtext"`
	Meta       string `gorm:"column:meta;type:text"`
	MetaJson   string `gorm:"column:meta_json;type:text"`
	Version    string `gorm:"column:version;type:varchar(255)"`
	Changelog  string `gorm:"column:changelog;type:text"`
	Status     int64  `gorm:"column:status;type:bigint(20)"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (s *Code) TableName() string {
	return "cdb_tampermonkey_script_code"
}

type LibDefinition struct {
	ID         int64  `gorm:"column:id"                                 json:"id"         form:"id"`
	UserId     int64  `gorm:"column:user_id;index:user_id;not null"     json:"user_id"    form:"user_id"`
	ScriptId   int64  `gorm:"column:script_id;index:script_id;not null" json:"script_id"  form:"script_id"`
	CodeId     int64  `gorm:"column:code_id;index:code_id;not null"     json:"code_id"    form:"code_id"`
	Definition string `gorm:"column:definition;not null"                json:"definition" form:"definition"`
	Createtime int64  `gorm:"column:createtime;type:bigint"             json:"createtime" form:"createtime"`
}
