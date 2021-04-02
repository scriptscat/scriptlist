package entity

import "github.com/scriptscat/scriptweb/internal/pkg/config"

// ScriptCategoryList 拥有的分类列表
type ScriptCategoryList struct {
	ID   int64  `gorm:"column:id" json:"id" form:"id"`
	Name string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	// 本分类下脚本数量
	Num        int64 `gorm:"column:num" json:"num" form:"num"`
	Createtime int64 `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64 `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

// ScriptCategory 脚本分类
type ScriptCategory struct {
	ID         int64 `gorm:"column:id" json:"id" form:"id"`
	CategoryId int64 `gorm:"column:category_id;index:category_id" json:"category_id" form:"category_id"`
	ScriptId   int64 `gorm:"column:script_id;index:script_id" json:"script_id" form:"script_id"`
	Createtime int64 `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64 `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

func (s *ScriptCategory) TableName() string {
	return config.AppConfig.Mysql.Prefix + "script_category"
}

// ScriptScore 脚本评分
type ScriptScore struct {
	ID       int64 `gorm:"column:id" json:"id" form:"id"`
	UserId   int64 `gorm:"column:user_id;index:user_script,unique;index:user" json:"user_id" form:"user_id"`
	ScriptId int64 `gorm:"column:script_id;index:user_script,unique;index:script" json:"script_id" form:"script_id"`
	// 评分,五星制,50
	Score int64 `gorm:"column:score" json:"score" form:"score"`
	// 评分原因
	Message    string `gorm:"column:message;type:varchar(255)" json:"message" form:"message"`
	Createtime int64  `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64  `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

// ScriptStatistics 脚本总下载更新统计
type ScriptStatistics struct {
	ID       int64 `gorm:"column:id" json:"id" form:"id"`
	ScriptId int64 `gorm:"column:script_id;index:script,unique" json:"script_id" form:"script_id"`
	Download int64 `gorm:"column:download" json:"download" form:"download"`
	Update   int64 `gorm:"column:update" json:"update" form:"update"`
}

// ScriptDateStatistics 脚本日下载更新统计
type ScriptDateStatistics struct {
	ID       int64  `gorm:"column:id" json:"id" form:"id"`
	ScriptId int64  `gorm:"column:script_id;index:script_date,unique" json:"script_id" form:"script_id"`
	Date     string `gorm:"column:date;index:script_date,unique" json:"date" form:"date"`
	Download int64  `gorm:"column:download" json:"download" form:"download"`
	Update   int64  `gorm:"column:update" json:"update" form:"update"`
}

// ScriptDomain 脚本域名
type ScriptDomain struct {
	ID            int64  `gorm:"column:id" json:"id" form:"id"`
	Domain        string `gorm:"column:Domain;type:varchar(255);index:domain_script,unique" json:"Domain" form:"Domain"`
	DomainReverse string `gorm:"column:domain_reverse;type:varchar(255);index:domain_reverse" json:"domain_reverse" form:"domain_reverse"`
	ScriptId      int64  `gorm:"column:script_id;index:script_id;index:domain_script,unique" json:"script_id" form:"script_id"`
	ScriptCodeId  int64  `gorm:"column:script_code_id" json:"script_code_id" form:"script_code_id"`
	Createtime    int64  `gorm:"column:createtime" json:"createtime" form:"createtime"`
}
