package repository

import (
	"time"

	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"gorm.io/gorm"
)

type statistics struct {
}

func NewStatistics() Statistics {
	return &statistics{}
}

func (s *statistics) Download(id int64) error {
	date := time.Now().Format("2006-01-02")
	if db.Db.Model(&entity.ScriptDateStatistics{}).Where("script_id=? and date=?", id, date).Update("download", gorm.Expr("download+1")).RowsAffected == 0 {
		if err := db.Db.Save(&entity.ScriptDateStatistics{
			ScriptId: id,
			Date:     date,
			Download: 1,
		}).Error; err != nil {
			return err
		}
	}
	if db.Db.Model(&entity.ScriptStatistics{}).Where("script_id=?", id).Update("download", gorm.Expr("download+1")).RowsAffected == 0 {
		return db.Db.Save(&entity.ScriptStatistics{
			ScriptId: id,
			Download: 1,
		}).Error
	}
	return nil
}

func (s *statistics) Update(id int64) error {
	date := time.Now().Format("2006-01-02")
	if db.Db.Model(&entity.ScriptDateStatistics{}).Where("script_id=? and date=?", id, date).Update("update", gorm.Expr("`update`+1")).RowsAffected == 0 {
		if err := db.Db.Save(&entity.ScriptDateStatistics{
			ScriptId: id,
			Date:     date,
			Update:   1,
		}).Error; err != nil {
			return err
		}
	}
	if db.Db.Model(&entity.ScriptStatistics{}).Where("script_id=?", id).Update("update", gorm.Expr("`update`+1")).RowsAffected == 0 {
		return db.Db.Save(&entity.ScriptStatistics{
			ScriptId: id,
			Update:   1,
		}).Error
	}
	return nil
}
