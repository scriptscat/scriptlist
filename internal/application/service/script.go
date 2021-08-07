package service

import (
	"github.com/golang/glog"
	request2 "github.com/scriptscat/scriptweb/interfaces/dto/request"
	respond2 "github.com/scriptscat/scriptweb/interfaces/dto/respond"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	service2 "github.com/scriptscat/scriptweb/internal/domain/script/service"
	service3 "github.com/scriptscat/scriptweb/internal/domain/statistics/service"
	"github.com/scriptscat/scriptweb/internal/domain/user/service"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
)

type Script interface {
	GetScript(id int64, version string, withcode bool) (*respond2.ScriptInfo, error)
	GetScriptList(category []int64, domain, keyword, sort string, page request2.Pages) (*respond2.List, error)
	GetUserScript(uid int64, self bool, page request2.Pages) (*respond2.List, error)
	GetScriptCodeList(id int64) ([]*respond2.ScriptCode, error)
	GetLatestScriptCode(id int64) (*respond2.ScriptCode, error)
	GetScriptCodeByVersion(id int64, version string) (*respond2.ScriptCode, error)
	GetCategory() ([]*entity.ScriptCategoryList, error)
	AddScore(uid int64, id int64, score *request2.Score) error
	ScoreList(id int64, page *request2.Pages) (*respond2.List, error)
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

func (s *script) GetScript(id int64, version string, withcode bool) (*respond2.ScriptInfo, error) {
	script, err := s.scriptSvc.Info(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userSvc.GetUser(script.UserId)
	if err != nil {
		return nil, err
	}
	latest, err := s.GetScriptCodeByVersion(id, version)
	if err != nil {
		return nil, err
	}
	ret := respond2.ToScriptInfo(user, script, latest)
	if withcode {
		ret.Script.Script.Code = latest.Code
	}
	return ret, nil
}

func (s *script) GetLatestScriptCode(id int64) (*respond2.ScriptCode, error) {
	codes, err := s.scriptSvc.VersionList(id)
	if err != nil {
		return nil, err
	}
	if len(codes) == 0 {
		return nil, errs.ErrScriptAudit
	}
	user, err := s.userSvc.GetUser(codes[0].UserId)
	ret := respond2.ToScriptCode(user, codes[0])
	ret.Code = codes[0].Code
	return ret, err

}

func (s *script) GetScriptList(category []int64, domain, keyword, sort string, page request2.Pages) (*respond2.List, error) {
	list, total, err := s.scriptSvc.Search(category, domain, keyword, sort, page)
	if err != nil {
		return nil, err
	}
	ret := make([]*respond2.Script, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		latest, err := s.GetLatestScriptCode(v.ID)
		if err != nil {
			glog.Errorf("GetLatestScriptCode: %v", err)
		}
		if latest != nil {
			ret[i] = respond2.ToScript(user, v, latest)
			s.join(ret[i])
		}
	}
	return &respond2.List{
		List:  ret,
		Total: total,
	}, nil
}

func (s *script) GetUserScript(uid int64, self bool, page request2.Pages) (*respond2.List, error) {
	list, total, err := s.scriptSvc.UserScript(uid, self, page)
	if err != nil {
		return nil, err
	}
	ret := make([]*respond2.Script, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		latest, err := s.GetLatestScriptCode(v.ID)
		if err != nil {
			glog.Errorf("GetLatestScriptCode: %v", err)
		}
		if latest != nil {
			ret[i] = respond2.ToScript(user, v, latest)
			s.join(ret[i])
		}
	}
	return &respond2.List{
		List:  ret,
		Total: total,
	}, nil
}

func (s *script) join(script *respond2.Script) {
	// 统计
	script.TotalInstall, _ = s.statisSvc.TotalDownload(script.ID)
	script.TodayInstall, _ = s.statisSvc.TodayDownload(script.ID)
	// 评分
	script.Score, _ = s.scoreSvc.GetAvgScore(script.ID)
	script.ScoreNum, _ = s.scoreSvc.Count(script.ID)
}

func (s *script) GetScriptCodeList(id int64) ([]*respond2.ScriptCode, error) {
	list, err := s.scriptSvc.VersionList(id)
	if err != nil {
		return nil, err
	}
	ret := make([]*respond2.ScriptCode, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		ret[i] = respond2.ToScriptCode(user, v)
	}
	return ret, nil
}

func (s *script) GetScriptCodeByVersion(id int64, version string) (*respond2.ScriptCode, error) {
	if version == "" {
		return s.GetLatestScriptCode(id)
	}
	list, err := s.scriptSvc.VersionList(id)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Version == version {
			user, _ := s.userSvc.GetUser(v.UserId)
			ret := respond2.ToScriptCode(user, v)
			ret.Code = v.Code
			return ret, err
		}
	}
	return nil, errs.ErrScriptCodeIsNil
}

func (s *script) GetCategory() ([]*entity.ScriptCategoryList, error) {
	return s.scriptSvc.GetCategory()
}

func (s *script) AddScore(uid int64, id int64, score *request2.Score) error {
	if _, err := s.scriptSvc.Info(id); err != nil {
		return err
	}
	return s.scoreSvc.AddScore(uid, id, score)
}

func (s *script) ScoreList(id int64, page *request2.Pages) (*respond2.List, error) {
	list, total, err := s.scoreSvc.ScoreList(id, page)
	if err != nil {
		return nil, err
	}
	resp := make([]*respond2.ScriptScore, len(list))
	for i, v := range list {
		user, _ := s.userSvc.GetUser(v.UserId)
		resp[i] = respond2.ToScriptScore(user, v)
	}

	return &respond2.List{
		List:  resp,
		Total: total,
	}, nil
}

func (s *script) UserScore(uid int64, id int64) (*entity.ScriptScore, error) {
	return s.scoreSvc.UserScore(uid, id)
}
