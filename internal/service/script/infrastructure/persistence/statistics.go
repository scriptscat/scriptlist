package persistence

import (
	"time"

	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"gorm.io/gorm"
)

type statistics struct {
	db *gorm.DB
}

func NewStatistics(db *gorm.DB) repository.Statistics {
	return &statistics{db: db}
}

func (s *statistics) Download(id int64) error {
	date := time.Now().Format("2006-01-02")
	if s.db.Model(&entity.ScriptDateStatistics{}).Where("script_id=? and date=?", id, date).Update("download", gorm.Expr("download+1")).RowsAffected == 0 {
		if err := s.db.Save(&entity.ScriptDateStatistics{
			ScriptId: id,
			Date:     date,
			Download: 1,
		}).Error; err != nil {
			return err
		}
	}
	if s.db.Model(&entity.ScriptStatistics{}).Where("script_id=?", id).Update("download", gorm.Expr("download+1")).RowsAffected == 0 {
		return s.db.Save(&entity.ScriptStatistics{
			ScriptId: id,
			Download: 1,
		}).Error
	}
	return nil
}

func (s *statistics) Update(id int64) error {
	date := time.Now().Format("2006-01-02")
	if s.db.Model(&entity.ScriptDateStatistics{}).Where("script_id=? and date=?", id, date).Update("update", gorm.Expr("`update`+1")).RowsAffected == 0 {
		if err := s.db.Save(&entity.ScriptDateStatistics{
			ScriptId: id,
			Date:     date,
			Update:   1,
		}).Error; err != nil {
			return err
		}
	}
	if s.db.Model(&entity.ScriptStatistics{}).Where("script_id=?", id).Update("update", gorm.Expr("`update`+1")).RowsAffected == 0 {
		return s.db.Save(&entity.ScriptStatistics{
			ScriptId: id,
			Update:   1,
		}).Error
	}
	return nil
}
