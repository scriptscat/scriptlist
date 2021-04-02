package entity

type StatisticsDownload struct {
	ID           int64  `gorm:"column:id" json:"id" form:"id"`
	UserId       int64  `gorm:"column:user_id" json:"user_id" form:"user_id"`
	ScriptId     int64  `gorm:"column:script_id;index:script_id;index:script_time" json:"script_id" form:"script_id"`
	ScriptCodeId int64  `gorm:"column:script_code_id;index:script_code_id;index:script_code_time" json:"script_code_id" form:"script_code_id"`
	Ip           string `gorm:"column:ip" json:"ip" form:"ip"`
	Ua           string `gorm:"column:ua" json:"ua" form:"ua"`
	Createtime   int64  `gorm:"column:createtime;index:script_time;index:script_code_time" json:"createtime" form:"createtime"`
}

type StatisticsUpdate StatisticsDownload
