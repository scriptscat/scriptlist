package entity

// ScriptCategoryList 拥有的分类列表
type ScriptCategoryList struct {
	ID   int64  `gorm:"column:id" json:"id" form:"id"`
	Name string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	// 本分类下脚本数量
	Num        int64 `gorm:"column:num" json:"num" form:"num"`
	Sort       int32 `gorm:"column:sort;type:int(10);index:category_sort" json:"sort" form:"sort"`
	Createtime int64 `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64 `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

// ScriptCategory 脚本分类
type ScriptCategory struct {
	ID         int64 `gorm:"column:id" json:"id" form:"id"`
	CategoryId int64 `gorm:"column:category_id;index:category_id;index:script_category,unique" json:"category_id" form:"category_id"`
	ScriptId   int64 `gorm:"column:script_id;index:script_id;index:script_category,unique" json:"script_id" form:"script_id"`
	Createtime int64 `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64 `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

// ScriptScore 脚本评分
type ScriptScore struct {
	ID       int64 `gorm:"column:id" json:"id" form:"id"`
	UserId   int64 `gorm:"column:user_id;index:user_script,unique;index:user" json:"user_id" form:"user_id"`
	ScriptId int64 `gorm:"column:script_id;index:user_script,unique;index:script" json:"script_id" form:"script_id"`
	// 评分,五星制,50
	Score int64 `gorm:"column:score" json:"score" form:"score"`
	// 评分原因
	Message    string `gorm:"column:message;type:text" json:"message" form:"message"`
	State      int32  `gorm:"column:state;type:int(10);default:1" json:"state" form:"state"`
	Createtime int64  `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64  `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

// ScriptStatistics 脚本总下载更新统计
type ScriptStatistics struct {
	ID         int64 `gorm:"column:id" json:"id" form:"id"`
	ScriptId   int64 `gorm:"column:script_id;index:script,unique" json:"script_id" form:"script_id"`
	Download   int64 `gorm:"column:download;default:0" json:"download" form:"download"`
	Update     int64 `gorm:"column:update;default:0" json:"update" form:"update"`
	Score      int64 `gorm:"column:score;default:0" json:"score" form:"score"`
	ScoreCount int64 `gorm:"column:score_count;default:0" json:"score_count" form:"score_count"`
}

// ScriptDateStatistics 脚本日下载更新统计
type ScriptDateStatistics struct {
	ID       int64  `gorm:"column:id" json:"id" form:"id"`
	ScriptId int64  `gorm:"column:script_id;index:script_date,unique;default:0" json:"script_id" form:"script_id"`
	Date     string `gorm:"type:varchar(255);column:date;index:script_date,unique;default:0" json:"date" form:"date"`
	Download int64  `gorm:"column:download;default:0" json:"download" form:"download"`
	Update   int64  `gorm:"column:update;default:0" json:"update" form:"update"`
}

// ScriptDomain 脚本域名
type ScriptDomain struct {
	ID            int64  `gorm:"column:id" json:"id" form:"id"`
	Domain        string `gorm:"column:domain;type:varchar(255);index:domain_script,unique" json:"domain" form:"domain"`
	DomainReverse string `gorm:"column:domain_reverse;type:varchar(255);index:domain_reverse" json:"domain_reverse" form:"domain_reverse"`
	ScriptId      int64  `gorm:"column:script_id;index:script_id;index:domain_script,unique" json:"script_id" form:"script_id"`
	ScriptCodeId  int64  `gorm:"column:script_code_id" json:"script_code_id" form:"script_code_id"`
	Createtime    int64  `gorm:"column:createtime" json:"createtime" form:"createtime"`
}
