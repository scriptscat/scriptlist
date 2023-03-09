package statistics_entity

type StatisticsInfo struct {
	ID            int64    `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID      int64    `gorm:"column:script_id;type:bigint(20);not null;index:script_id,unique"`
	StatisticsKey string   `gorm:"column:statistics_key;type:varchar(128);index:statistics_key,unique"`
	Whitelist     []string `gorm:"column:whitelist;type:json;not null"`
	Status        int      `gorm:"column:status;type:tinyint(2);not null;default:1"`
	Createtime    int64    `gorm:"column:createtime;type:bigint(20);not null"`
}
