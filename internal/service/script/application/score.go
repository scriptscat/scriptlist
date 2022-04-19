package application

import (
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
)

type Score interface {
	AddScore(uid, scriptId int64, msg *request.Score) (bool, error)
	GetAvgScore(scriptId int64) (int64, error)
	Count(scriptId int64) (int64, error)
	UserScore(uid int64, scriptId int64) (*entity.ScriptScore, error)
	ScoreList(scriptId int64, page *request.Pages) ([]*entity.ScriptScore, int64, error)
	Delete(scriptId, scoreId int64) error
}

type score struct {
	scriptRepo repository.Script
	repo       repository.Score
}

func NewScore(scriptRepo repository.Script, repo repository.Score) Score {
	return &score{scriptRepo: scriptRepo, repo: repo}
}

func (s *score) AddScore(uid, scriptId int64, msg *request.Score) (bool, error) {
	script, err := s.scriptRepo.Find(scriptId)
	if err != nil {
		return false, err
	}
	score, err := s.repo.UserScore(uid, scriptId)
	if err != nil {
		return false, err
	}
	saveScore, err := script.AddScore(uid, score, msg)
	if err != nil {
		return false, err
	}
	err = s.repo.Save(saveScore)
	if err != nil {
		return false, err
	}
	return score != nil, nil
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

func (s *score) ScoreList(scriptId int64, page *request.Pages) ([]*entity.ScriptScore, int64, error) {
	return s.repo.List(scriptId, page)
}

func (s *score) Delete(scriptId, scoreId int64) error {
	script, err := s.scriptRepo.Find(scriptId)
	if err != nil {
		return err
	}
	score, err := s.repo.Find(scoreId)
	if err != nil {
		return err
	}
	if err := script.DeleteScore(score); err != nil {
		return err
	}
	return s.repo.Save(score)
}
