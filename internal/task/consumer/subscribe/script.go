package subscribe

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	script_repo2 "github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"go.uber.org/zap"
)

type Script struct {
	// 分类id
	bgCategory   *entity.ScriptCategoryList
	cronCategory *entity.ScriptCategoryList
}

func (s *Script) Subscribe(ctx context.Context, broker broker.Broker) error {
	var err error
	s.bgCategory, err = script_repo2.ScriptCategoryList().FindByName(ctx, "后台脚本")
	if err != nil {
		return err
	}
	if s.bgCategory == nil {
		s.bgCategory = &entity.ScriptCategoryList{
			Name:       "后台脚本",
			Createtime: time.Now().Unix(),
		}
		if err := script_repo2.ScriptCategoryList().Create(ctx, s.bgCategory); err != nil {
			return err
		}
	}
	s.cronCategory, err = script_repo2.ScriptCategoryList().FindByName(ctx, "定时脚本")
	if err != nil {
		return err
	}
	if s.cronCategory == nil {
		s.cronCategory = &entity.ScriptCategoryList{
			Name:       "定时脚本",
			Createtime: time.Now().Unix(),
		}
		if err := script_repo2.ScriptCategoryList().Create(ctx, s.cronCategory); err != nil {
			return err
		}
	}
	_, err = broker.Subscribe(ctx,
		producer.ScriptCreateTopic, s.scriptCreateHandler,
	)
	if err != nil {
		return err
	}
	_, err = broker.Subscribe(ctx, producer.ScriptCodeUpdateTopic, s.scriptCodeUpdate)
	return err
}

// 消费脚本创建消息,根据meta信息进行分类
func (s *Script) scriptCreateHandler(ctx context.Context, event broker.Event) error {
	msg, err := producer.ParseScriptCreateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).
			Error("json.Unmarshal", zap.Error(err), zap.String("body", string(event.Message().Body)))
		return err
	}
	if msg.Script == nil {
		return errors.New("script is nil")
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.Script.ID))

	// 根据meta信息, 将脚本分类到后台脚本, 定时脚本, 用户脚本
	metaJson := make(map[string][]string)
	if err := json.Unmarshal([]byte(msg.Code.MetaJson), &metaJson); err != nil {
		logger.Error("json.Unmarshal", zap.Error(err), zap.String("meta", msg.Code.MetaJson))
		return err
	}

	// 处理domain
	if err := s.saveDomain(ctx, msg.Script.ID, msg.Code.ID, metaJson); err != nil {
		logger.Error("saveDomain", zap.Error(err))
		return err
	}

	if len(metaJson["background"]) > 0 || len(metaJson["crontab"]) > 0 {
		// 后台脚本
		if err := script_repo2.ScriptCategory().LinkCategory(ctx, msg.Script.ID, s.bgCategory.ID); err != nil {
			logger.Error("LinkCategory", zap.Error(err))
			return err
		}
	}
	if len(metaJson["crontab"]) > 0 {
		// 定时脚本
		if err := script_repo2.ScriptCategory().LinkCategory(ctx, msg.Script.ID, s.cronCategory.ID); err != nil {
			logger.Error("LinkCategory", zap.Error(err))
			return err
		}
	}

	// 关注自己脚本
	if err := script_repo2.ScriptWatch().Watch(ctx, msg.Script.ID, msg.Script.UserID, entity.ScriptWatchLevelIssueComment); err != nil {
		logger.Error("Watch", zap.Error(err))
	}

	return nil
}

// 消费脚本代码更新消息,发送邮件通知给关注了的用户
func (s *Script) scriptCodeUpdate(ctx context.Context, event broker.Event) error {
	msg, err := producer.ParseScriptCodeUpdateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).
			Error("json.Unmarshal", zap.Error(err), zap.String("body", string(event.Message().Body)))
		return err
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.Script.ID))

	metaJson := make(map[string][]string)
	if err := json.Unmarshal([]byte(msg.Code.MetaJson), &metaJson); err != nil {
		logger.Error("json.Unmarshal", zap.Error(err), zap.String("meta", msg.Code.MetaJson))
		return err
	}
	// 处理domain
	if err := s.saveDomain(ctx, msg.Script.ID, msg.Code.ID, metaJson); err != nil {
		logger.Error("saveDomain", zap.Error(err))
		return err
	}
	logger.Info("update script code")

	list, err := script_repo2.ScriptWatch().ListAll(ctx, msg.Script.ID)
	if err != nil {
		logger.Error("获取关注列表失败", zap.Error(err))
	} else {
		uids := make([]int64, 0)
		for _, v := range list {
			uids = append(uids, v.UserID)
		}
		err := notice_svc.Notice().MultipleSend(ctx, uids, notice_svc.ScriptUpdateTemplate, notice_svc.WithParams(&template.ScriptUpdate{
			ID:      msg.Script.ID,
			Name:    msg.Script.Name,
			Version: msg.Code.Version,
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
		result, err := script_repo2.Domain().FindByDomain(ctx, id, domain)
		if err != nil {
			logger.Ctx(ctx).Error("FindByDomain", zap.Error(err), zap.Int64("script_id", id), zap.String("domain", domain))
			continue
		}
		if result == nil {
			e := &entity.ScriptDomain{
				Domain:        domain,
				DomainReverse: utils.StringReverse(domain),
				ScriptID:      id,
				ScriptCodeID:  codeID,
				Createtime:    time.Now().Unix(),
			}
			if err := script_repo2.Domain().Create(ctx, e); err != nil {
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
