package entity

type ScriptStatistics struct {
	ID         int64 `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID   int64 `gorm:"column:script_id;type:bigint(20);index:script,unique"`
	Download   int64 `gorm:"column:download;type:bigint(20);default:0"`
	Update     int64 `gorm:"column:update;type:bigint(20);default:0"`
	Score      int64 `gorm:"column:score;type:bigint(20);default:0"`
	ScoreCount int64 `gorm:"column:score_count;type:bigint(20);default:0"`
}
