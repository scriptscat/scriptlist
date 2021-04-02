package repository

import (
	"context"
	"fmt"

	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/request"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"gorm.io/gorm"
)

type score struct {
}

func NewScore() Score {
	return &score{}
}

func (s *score) scorekey(scriptId int64, key string) string {
	return fmt.Sprintf("script:score:%d:%s", scriptId, key)
}

func (s *score) Save(score *entity.ScriptScore) error {
	if score.ID == 0 {
		if err := db.Redis.IncrBy(context.Background(), s.scorekey(score.ScriptId, "total"), 1).Err(); err != nil {
			return err
		}
		if err := db.Redis.IncrBy(context.Background(), s.scorekey(score.ScriptId, "score"), score.Score).Err(); err != nil {
			return err
		}
		return db.Db.Create(score).Error
	}
	old := &entity.ScriptScore{ID: score.ID}
	if err := db.Db.First(old).Error; err != nil {
		return err
	}
	if err := db.Redis.IncrBy(context.Background(), s.scorekey(score.ScriptId, "score"), score.Score-old.Score).Err(); err != nil {
		return err
	}
	return db.Db.Updates(score).Error
}

func (s *score) UserScore(uid, scriptId int64) (*entity.ScriptScore, error) {
	ret := &entity.ScriptScore{}
	if err := db.Db.Where("user_id=? and script_id=?", uid, scriptId).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *score) Avg(scriptId int64) (int64, error) {
	total, err := db.Redis.Get(context.Background(), s.scorekey(scriptId, "total")).Int64()
	if err != nil {
		return 0, err
	}
	score, err := db.Redis.Get(context.Background(), s.scorekey(scriptId, "score")).Int64()
	if err != nil {
		return 0, err
	}
	return score / total, nil
}

func (s *score) Count(scriptId int64) (int64, error) {
	return db.Redis.Get(context.Background(), s.scorekey(scriptId, "total")).Int64()
}

func (s *score) List(scriptId int64, page *request.Pages) ([]*entity.ScriptScore, int64, error) {
	list := make([]*entity.ScriptScore, 0)
	find := db.Db.Model(&entity.ScriptScore{}).Where("script_id=?", scriptId).Order("createtime desc")
	var num int64
	if err := find.Count(&num).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size()).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, num, nil

}
