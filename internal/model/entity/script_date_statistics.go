package entity

type ScriptDateStatistics struct {
	ID       int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID int64  `gorm:"column:script_id;type:bigint(20);default:0;index:script_date,unique"`
	Date     string `gorm:"column:date;type:varchar(255);default:0"`
	Download int64  `gorm:"column:download;type:bigint(20);default:0"`
	Update   int64  `gorm:"column:update;type:bigint(20);default:0"`
}
