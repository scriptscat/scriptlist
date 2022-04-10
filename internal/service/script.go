package service

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	repository2 "github.com/scriptscat/scriptlist/internal/service/safe/domain/repository"
	service4 "github.com/scriptscat/scriptlist/internal/service/safe/service"
	"github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/vo"
	service3 "github.com/scriptscat/scriptlist/internal/service/statistics/service"
	"github.com/scriptscat/scriptlist/internal/service/user/service"
	"github.com/scriptscat/scriptlist/pkg/httputils"
)

type Script interface {
	GetScript(id int64, version string, withcode bool) (*vo.ScriptInfo, error)
	GetScriptList(search *repository.SearchList, page *request.Pages) (*httputils.List, error)
	GetScriptCodeList(id int64, page *request.Pages) (*httputils.List, error)
	GetLatestScriptCode(id int64, withcode bool) (*vo.ScriptCode, error)
	GetScriptCodeByVersion(id int64, version string, withcode bool) (*vo.ScriptCode, error)
	GetCategory() ([]*entity.ScriptCategoryList, error)
	AddScore(uid int64, id int64, score *request.Score) (bool, error)
	ScoreList(id int64, page *request.Pages) (*httputils.List, error)
	UserScore(uid, id int64) (*entity.ScriptScore, error)
	CreateScript(uid int64, req *request.CreateScript) (*entity.Script, error)
	UpdateScript(uid, id int64, req *request.UpdateScript) error
	UpdateScriptCode(uid, id int64, req *request.UpdateScriptCode) error
	SyncScript(uid, id int64) error
}

type script struct {
	userSvc   service.User
	scriptSvc application.Script
	scoreSvc  application.Score
	statisSvc service3.Statistics
	rateSvc   service4.Rate
}

func NewScript(userSvc service.User, scriptSvc application.Script, scoreSvc application.Score, statisSvc service3.Statistics, rateSvc service4.Rate, c *cron.Cron) Script {
	return &script{
		userSvc:   userSvc,
		scriptSvc: scriptSvc,
		scoreSvc:  scoreSvc,
		statisSvc: statisSvc,
		rateSvc:   rateSvc,
	}
}

func (s *script) GetScript(id int64, version string, withcode bool) (*vo.ScriptInfo, error) {
	script, err := s.scriptSvc.Info(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userSvc.UserInfo(script.UserId)
	if err != nil {
		return nil, err
	}
	latest, err := s.GetScriptCodeByVersion(id, version, withcode)
	if err != nil {
		return nil, err
	}

	ret := vo.ToScriptInfo(user, script, latest)
	s.join(ret.Script)
	return ret, nil
}

func (s *script) GetLatestScriptCode(id int64, withcode bool) (*vo.ScriptCode, error) {
	code, err := s.scriptSvc.GetLatestVersion(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userSvc.UserInfo(code.UserId)
	ret := vo.ToScriptCode(user, code)
	if withcode {
		ret.Meta = code.Meta
		ret.Code = code.Code
		if d, err := s.scriptSvc.GetCodeDefinition(code.ID); err == nil {
			ret.Definition = d.Definition
		}
	}
	return ret, err
}

func (s *script) GetScriptList(search *repository.SearchList, page *request.Pages) (*httputils.List, error) {
	list, total, err := s.scriptSvc.Search(search, page)
	if err != nil {
		return nil, err
	}
	ret := make([]interface{}, len(list))
	for i, v := range list {
		user, _ := s.userSvc.UserInfo(v.UserId)
		latest, err := s.GetLatestScriptCode(v.ID, false)
		if err != nil {
			glog.Errorf("GetLatestScriptCode: %v", err)
		}
		if latest != nil {
			item := vo.ToScript(user, v, latest)
			s.join(item)
			ret[i] = item
		}
	}
	return &httputils.List{
		List:  ret,
		Total: total,
	}, nil
}

func (s *script) join(script *vo.Script) {
	// 统计
	script.TotalInstall, _ = s.statisSvc.TotalDownload(script.ID)
	script.TodayInstall, _ = s.statisSvc.TodayDownload(script.ID)
	// 评分
	script.Score, _ = s.scoreSvc.GetAvgScore(script.ID)
	script.ScoreNum, _ = s.scoreSvc.Count(script.ID)
}

func (s *script) GetScriptCodeList(id int64, page *request.Pages) (*httputils.List, error) {
	list, num, err := s.scriptSvc.VersionList(id, page)
	if err != nil {
		return nil, err
	}
	ret := make([]interface{}, len(list))
	for i, v := range list {
		user, _ := s.userSvc.UserInfo(v.UserId)
		ret[i] = vo.ToScriptCode(user, v)
	}
	return &httputils.List{
		List:  ret,
		Total: num,
	}, nil
}

func (s *script) GetScriptCodeByVersion(id int64, version string, withcode bool) (*vo.ScriptCode, error) {
	if version == "" {
		return s.GetLatestScriptCode(id, withcode)
	}
	code, err := s.scriptSvc.GetScriptVersion(id, version)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, errs.ErrScriptCodeIsNil
	}
	user, _ := s.userSvc.UserInfo(code.UserId)
	ret := vo.ToScriptCode(user, code)
	if withcode {
		ret.Code = code.Code
		if d, err := s.scriptSvc.GetCodeDefinition(code.ID); err == nil {
			ret.Definition = d.Definition
		}
	}
	return ret, nil
}

func (s *script) GetCategory() ([]*entity.ScriptCategoryList, error) {
	return s.scriptSvc.GetCategory()
}

func (s *script) AddScore(uid int64, id int64, score *request.Score) (bool, error) {
	if _, err := s.scriptSvc.Info(id); err != nil {
		return false, err
	}
	return s.scoreSvc.AddScore(uid, id, score)
}

func (s *script) ScoreList(id int64, page *request.Pages) (*httputils.List, error) {
	list, total, err := s.scoreSvc.ScoreList(id, page)
	if err != nil {
		return nil, err
	}
	resp := make([]interface{}, len(list))
	for i, v := range list {
		user, _ := s.userSvc.UserInfo(v.UserId)
		resp[i] = vo.ToScriptScore(user, v)
	}

	return &httputils.List{
		List:  resp,
		Total: total,
	}, nil
}

func (s *script) UserScore(uid int64, id int64) (*entity.ScriptScore, error) {
	return s.scoreSvc.UserScore(uid, id)
}

func (s *script) CreateScript(uid int64, req *request.CreateScript) (*entity.Script, error) {
	var ret *entity.Script
	if err := s.rateSvc.Rate(&repository2.RateUserInfo{Uid: uid}, &repository2.RateRule{
		Name:     "post-script",
		Interval: 60,
	}, func() error {
		script, err := s.scriptSvc.CreateScript(uid, req)
		if err != nil {
			return err
		}
		ret = script
		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *script) UpdateScript(uid, id int64, req *request.UpdateScript) error {
	return s.rateSvc.Rate(&repository2.RateUserInfo{Uid: uid}, &repository2.RateRule{
		Name:     "update-script",
		Interval: 5,
	}, func() error {
		if err := s.scriptSvc.UpdateScript(uid, id, req); err != nil {
			return err
		}
		if req.SyncMode == application.SyncModeManual {
			return s.SyncScript(uid, id)
		}
		return nil
	})
}

func (s *script) UpdateScriptCode(uid, id int64, req *request.UpdateScriptCode) error {
	return s.rateSvc.Rate(&repository2.RateUserInfo{Uid: uid}, &repository2.RateRule{
		Name:     "update-script-code",
		Interval: 10,
	}, func() error {
		return s.scriptSvc.CreateScriptCode(uid, id, req)
	})
}

func (s *script) SyncScript(uid, id int64) error {
	script, err := s.scriptSvc.Info(id)
	if err != nil {
		return err
	}
	if script.SyncUrl == "" {
		return errs.NewBadRequestError(1000, "同步链接为空")
	}
	req := &request.UpdateScriptCode{
		Content:    script.Content,
		Definition: "",
		Changelog:  "",
		Public:     script.Public,
		Unwell:     script.Unwell,
	}
	req.Code, err = s.requestSyncUrl(script.SyncUrl)
	if err != nil {
		return err
	}
	if script.ContentUrl != "" {
		req.Content, err = s.requestSyncUrl(script.ContentUrl)
		if err != nil {
			return err
		}
	}
	if script.Type == entity.LIBRARY_TYPE && script.DefinitionUrl != "" {
		req.Definition, err = s.requestSyncUrl(script.DefinitionUrl)
		if err != nil {
			return err
		}
	}
	return s.scriptSvc.CreateScriptCode(uid, id, req)
}

func (s *script) requestSyncUrl(syncUrl string) (string, error) {
	c := http.Client{
		Timeout: time.Second * 10,
	}
	u, err := url.Parse(syncUrl)
	if err != nil {
		return "", errs.NewErrScriptSyncNetwork(syncUrl, err)
	}
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
	})
	if err != nil {
		return "", errs.NewErrScriptSyncNetwork(syncUrl, err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(b), nil
}
