package service

import (
	"time"

	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	"github.com/scriptscat/scriptlist/internal/domain/script/repository"
	request2 "github.com/scriptscat/scriptlist/internal/http/dto/request"
)

type Score interface {
	AddScore(uid, scriptId int64, msg *request2.Score) (bool, error)
	GetAvgScore(scriptId int64) (int64, error)
	Count(scriptId int64) (int64, error)
	UserScore(uid int64, scriptId int64) (*entity.ScriptScore, error)
	ScoreList(scriptId int64, page *request2.Pages) ([]*entity.ScriptScore, int64, error)
}

type score struct {
	repo repository.Score
}

func NewScore(repo repository.Score) Score {
	return &score{repo: repo}
}

func (s *score) AddScore(uid, scriptId int64, msg *request2.Score) (bool, error) {
	score, err := s.repo.UserScore(uid, scriptId)
	if err != nil {
		return false, err
	}
	exist := true
	if score == nil {
		score = &entity.ScriptScore{
			UserId:     uid,
			ScriptId:   scriptId,
			Createtime: time.Now().Unix(),
		}
		exist = false
	}
	score.Score = msg.Score
	score.Message = msg.Message
	score.Updatetime = time.Now().Unix()
	err = s.repo.Save(score)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s *score) GetAvgScore(scriptId int64) (int64, error) {
	return s.repo.Avg(scriptId)
}

func (s *score) Count(scriptId int64) (int64, error) {
	return s.repo.Count(scriptId)
}

func (s *score) UserScore(uid int64, scriptId int64) (*entity.ScriptScore, error) {
	return s.repo.UserScore(uid, scriptId)
}

func (s *score) ScoreList(scriptId int64, page *request2.Pages) ([]*entity.ScriptScore, int64, error) {
	return s.repo.List(scriptId, page)
}
