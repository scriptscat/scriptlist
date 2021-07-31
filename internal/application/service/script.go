package service

import (
	"github.com/golang/glog"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	service2 "github.com/scriptscat/scriptweb/internal/domain/script/service"
	service3 "github.com/scriptscat/scriptweb/internal/domain/statistics/service"
	"github.com/scriptscat/scriptweb/internal/domain/user/service"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/request"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/respond"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
)

type Script interface {
	GetScript(id int64) (*respond.ScriptInfo, error)
	GetScriptList(category []int64, domain, keyword, sort string, page request.Pages) (*respond.List, error)
	GetUserScript(uid int64, self bool, page request.Pages) (*respond.List, error)
	GetScriptCodeList(id int64) ([]*respond.ScriptCode, error)
	GetLatestScriptCode(id int64) (*respond.ScriptCodeInfo, error)
	GetScriptCodeByVersion(id int64, version string) (*respond.ScriptCodeInfo, error)
	GetCategory() ([]*entity.ScriptCategoryList, error)
	AddScore(uid int64, id int64, score *request.Score) error
	ScoreList(id int64, page *request.Pages) (*respond.List, error)
	UserScore(uid int64, id int64) (*entity.ScriptScore, error)
}

type script struct {
	userSvc   service.User
	scriptSvc service2.Script
	scoreSvc  service2.Score
	statisSvc service3.Statistics
}

func NewScript(userSvc service.User, scriptSvc service2.Script, scoreSvc service2.Score, statisSvc service3.Statistics) Script {
	return &script{
		userSvc:   userSvc,
		scriptSvc: scriptSvc,
		scoreSvc:  scoreSvc,
		statisSvc: statisSvc,
	}
}

func (s *script) GetScript(id int64) (*respond.ScriptInfo, error) {
	script, err := s.scriptSvc.Info(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userSvc.GetUser(script.UserId)
	if err != nil {
		return nil, err
	}
	latest, err := s.GetLatestScriptCode(id)
	if err != nil {
		return nil, err
	}
	return respond.ToScriptInfo(user, script, latest.ScriptCode), nil
}

func (s *script) GetLatestScriptCode(id int64) (*respond.ScriptCodeInfo, error) {
	codes, err := s.scriptSvc.VersionList(id)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return nil, errs.ErrScriptAudit
	}
	user, err := s.userSvc.GetUser(codes[0].UserId)
	if err != nil {
		return respond.ToScriptCodeInfo(user, codes[0]), err
	}
	return respond.ToScriptCodeInfo(user, codes[0]), nil
}

func (s *script) GetScriptList(category []int64, domain, keyword, sort string, page request.Pages) (*respond.List, error) {
	list, total, err := s.scriptSvc.Search(category, domain, keyword, sort, page)
	if err != nil {
		return nil, err
	}
	ret := make([]*respond.Script, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		latest, err := s.GetLatestScriptCode(v.ID)
		if err != nil {
			glog.Errorf("GetLatestScriptCode: %v", err)
		}
		if latest != nil {
			ret[i] = respond.ToScript(user, v, latest.ScriptCode)
			s.join(ret[i])
		}
	}
	return &respond.List{
		List:  ret,
		Total: total,
	}, nil
}

func (s *script) GetUserScript(uid int64, self bool, page request.Pages) (*respond.List, error) {
	list, total, err := s.scriptSvc.UserScript(uid, self, page)
	if err != nil {
		return nil, err
	}
	ret := make([]*respond.Script, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		latest, err := s.GetLatestScriptCode(v.ID)
		if err != nil {
			glog.Errorf("GetLatestScriptCode: %v", err)
		}
		if latest != nil {
			ret[i] = respond.ToScript(user, v, latest.ScriptCode)
			s.join(ret[i])
		}
	}
	return &respond.List{
		List:  ret,
		Total: total,
	}, nil
}

func (s *script) join(script *respond.Script) {
	// 统计
	script.TotalInstall, _ = s.statisSvc.TotalDownload(script.ID)
	script.TodayInstall, _ = s.statisSvc.TodayDownload(script.ID)
	// 评分
	script.Score, _ = s.scoreSvc.GetAvgScore(script.ID)
	script.ScoreNum, _ = s.scoreSvc.Count(script.ID)
}

func (s *script) GetScriptCodeList(id int64) ([]*respond.ScriptCode, error) {
	list, err := s.scriptSvc.VersionList(id)
	if err != nil {
		return nil, err
	}
	ret := make([]*respond.ScriptCode, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		ret[i] = respond.ToScriptCode(user, v)
	}
	return ret, nil
}

func (s *script) GetScriptCodeByVersion(id int64, version string) (*respond.ScriptCodeInfo, error) {
	list, err := s.scriptSvc.VersionList(id)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Version == version {
			user, _ := s.userSvc.GetUser(v.UserId)
			return respond.ToScriptCodeInfo(user, v), nil
		}
	}
	return nil, errs.ErrScriptCodeIsNil
}

func (s *script) GetCategory() ([]*entity.ScriptCategoryList, error) {
	return s.scriptSvc.GetCategory()
}

func (s *script) AddScore(uid int64, id int64, score *request.Score) error {
	if _, err := s.GetScript(id); err != nil {
		return err
	}
	return s.scoreSvc.AddScore(uid, id, score)
}

func (s *script) ScoreList(id int64, page *request.Pages) (*respond.List, error) {
	list, total, err := s.scoreSvc.ScoreList(id, page)
	if err != nil {
		return nil, err
	}
	resp := make([]*respond.ScriptScore, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		resp[i] = respond.ToScriptScore(user, v)
	}

	return &respond.List{
		List:  resp,
		Total: total,
	}, nil
}

func (s *script) UserScore(uid int64, id int64) (*entity.ScriptScore, error) {
	return s.scoreSvc.UserScore(uid, id)
}
