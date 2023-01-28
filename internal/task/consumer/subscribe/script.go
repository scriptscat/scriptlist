package subscribe

import (
	"context"
	"encoding/json"
	"regexp"
	"time"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"go.uber.org/zap"
)

type Script struct {
	// 分类id
	bgCategory   *script_entity.ScriptCategoryList
	cronCategory *script_entity.ScriptCategoryList
}

func (s *Script) Subscribe(ctx context.Context) error {
	var err error
	s.bgCategory, err = script_repo.ScriptCategoryList().FindByName(ctx, "后台脚本")
	if err != nil {
		return err
	}
	if s.bgCategory == nil {
		s.bgCategory = &script_entity.ScriptCategoryList{
			Name:       "后台脚本",
			Createtime: time.Now().Unix(),
		}
		if err := script_repo.ScriptCategoryList().Create(ctx, s.bgCategory); err != nil {
			return err
		}
	}
	s.cronCategory, err = script_repo.ScriptCategoryList().FindByName(ctx, "定时脚本")
	if err != nil {
		return err
	}
	if s.cronCategory == nil {
		s.cronCategory = &script_entity.ScriptCategoryList{
			Name:       "定时脚本",
			Createtime: time.Now().Unix(),
		}
		if err := script_repo.ScriptCategoryList().Create(ctx, s.cronCategory); err != nil {
			return err
		}
	}
	if err := producer.SubscribeScriptCreate(ctx, s.scriptCreate); err != nil {
		return err
	}
	if err := producer.SubscribeScriptCodeUpdate(ctx, s.scriptCodeUpdate); err != nil {
		return err
	}
	return nil
}

// 消费脚本创建消息,根据meta信息进行分类
func (s *Script) scriptCreate(ctx context.Context, script *script_entity.Script, code *script_entity.Code) error {
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", script.ID))
	// 根据meta信息, 将脚本分类到后台脚本, 定时脚本, 用户脚本
	metaJson := make(map[string][]string)
	if err := json.Unmarshal([]byte(code.MetaJson), &metaJson); err != nil {
		logger.Error("json.Unmarshal", zap.Error(err), zap.String("meta", code.MetaJson))
		return err
	}

	// 处理domain
	if err := s.saveDomain(ctx, script.ID, code.ID, metaJson); err != nil {
		logger.Error("saveDomain", zap.Error(err))
		return err
	}

	if len(metaJson["background"]) > 0 || len(metaJson["crontab"]) > 0 {
		// 后台脚本
		if err := script_repo.ScriptCategory().LinkCategory(ctx, script.ID, s.bgCategory.ID); err != nil {
			logger.Error("LinkCategory", zap.Error(err))
			return err
		}
	}
	if len(metaJson["crontab"]) > 0 {
		// 定时脚本
		if err := script_repo.ScriptCategory().LinkCategory(ctx, script.ID, s.cronCategory.ID); err != nil {
			logger.Error("LinkCategory", zap.Error(err))
			return err
		}
	}

	// 关注自己脚本
	if err := script_repo.ScriptWatch().Watch(ctx, script.ID, script.UserID, script_entity.ScriptWatchLevelIssueComment); err != nil {
		logger.Error("Watch", zap.Error(err))
	}

	return nil
}

// 消费脚本代码更新消息,发送邮件通知给关注了的用户
func (s *Script) scriptCodeUpdate(ctx context.Context, script *script_entity.Script, code *script_entity.Code) error {
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", script.ID))

	metaJson := make(map[string][]string)
	if err := json.Unmarshal([]byte(code.MetaJson), &metaJson); err != nil {
		logger.Error("json.Unmarshal", zap.Error(err), zap.String("meta", code.MetaJson))
		return err
	}
	// 处理domain
	if err := s.saveDomain(ctx, script.ID, code.ID, metaJson); err != nil {
		logger.Error("saveDomain", zap.Error(err))
		return err
	}
	logger.Info("update script code")

	list, err := script_repo.ScriptWatch().FindAll(ctx, script.ID, script_entity.ScriptWatchLevelVersion)
	if err != nil {
		logger.Error("获取关注列表失败", zap.Error(err))
	} else {
		uids := make([]int64, 0)
		for _, v := range list {
			uids = append(uids, v.UserID)
		}
		err := notice_svc.Notice().MultipleSend(ctx, uids, notice_svc.ScriptUpdateTemplate, notice_svc.WithParams(&template.ScriptUpdate{
			ID:      script.ID,
			Name:    script.Name,
			Version: code.Version,
		}))
		if err != nil {
			logger.Error("发送邮件失败", zap.Error(err))
		}
	}
	return nil
}

// 保存脚本相关域名
func (s *Script) saveDomain(ctx context.Context, id, codeID int64, meta map[string][]string) error {
	domains := make(map[string]struct{})
	for _, v := range meta["match"] {
		domain := s.parseMatchDomain(v)
		if domain == "" {
			continue
		}
		domains[domain] = struct{}{}
	}
	for _, v := range meta["include"] {
		domain := s.parseMatchDomain(v)
		if domain == "" {
			continue
		}
		domains[domain] = struct{}{}
	}
	for domain := range domains {
		result, err := script_repo.Domain().FindByDomain(ctx, id, domain)
		if err != nil {
			logger.Ctx(ctx).Error("FindByDomain", zap.Error(err), zap.Int64("script_id", id), zap.String("domain", domain))
			continue
		}
		if result == nil {
			e := &script_entity.ScriptDomain{
				Domain:        domain,
				DomainReverse: utils.StringReverse(domain),
				ScriptID:      id,
				ScriptCodeID:  codeID,
				Createtime:    time.Now().Unix(),
			}
			if err := script_repo.Domain().Create(ctx, e); err != nil {
				logger.Ctx(ctx).Error("Create", zap.Error(err), zap.Int64("script_id", id), zap.String("domain", domain))
			}
		}
	}
	return nil
}

func (s *Script) parseMatchDomain(meta string) string {
	reg := regexp.MustCompile("(.+?://|^)(.+?)(/|$)")
	ret := reg.FindStringSubmatch(meta)
	if len(ret) == 0 || ret[2] == "" {
		return ""
	}
	if ret[2] == "*" {
		return "*"
	}
	if ret[2][0] == '*' {
		ret[2] = ret[2][1:]
	}
	if ret[2][0] == '.' {
		ret[2] = ret[2][1:]
	}
	domain, err := publicsuffix.Domain(ret[2])
	if err != nil {
		return ""
	}
	if domain[0] == '*' {
		return "*"
	}
	return domain
}
