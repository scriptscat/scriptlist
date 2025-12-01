package subscribe

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/pkg/utils"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc/template"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"go.uber.org/zap"
)

type Script struct {
}

func (s *Script) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptCreate(ctx, s.scriptCreate); err != nil {
		return err
	}
	if err := producer.SubscribeScriptCodeUpdate(ctx, s.scriptCodeUpdate); err != nil {
		return err
	}
	return nil
}

// 消费脚本创建消息,根据meta信息进行分类
func (s *Script) scriptCreate(ctx context.Context, script *script_entity.Script, codeId int64) error {
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", script.ID), zap.Int64("code_id", codeId))
	code, err := script_repo.ScriptCode().Find(ctx, codeId)
	if err != nil {
		return err
	}
	if code == nil {
		logger.Error("code不存在")
		return errors.New("code不存在")
	}
	if script.Type == script_entity.UserscriptType {
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
	}

	// 关注自己脚本
	if err := script_repo.ScriptWatch().Watch(ctx, script.ID, script.UserID, script_entity.ScriptWatchLevelIssueComment); err != nil {
		logger.Error("Watch", zap.Error(err))
	}

	return nil
}

// 消费脚本代码更新消息,发送邮件通知给关注了的用户
func (s *Script) scriptCodeUpdate(ctx context.Context, script *script_entity.Script, codeId int64) error {
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", script.ID), zap.Int64("code_id", codeId))
	code, err := script_repo.ScriptCode().Find(ctx, codeId)
	if err != nil {
		return err
	}
	if code == nil {
		logger.Error("code不存在")
		return errors.New("code不存在")
	}
	if script.Type == script_entity.UserscriptType {
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
		err := notification_svc.Notification().MultipleSend(ctx, uids, notification_entity.ScriptUpdateTemplate, notification_svc.WithParams(&template.ScriptUpdate{
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
		topDomain, domain := s.parseMatchDomain(v)
		if topDomain == "" {
			continue
		}
		domains[topDomain] = struct{}{}
		domains[domain] = struct{}{}
	}
	for _, v := range meta["include"] {
		topDomain, domain := s.parseMatchDomain(v)
		if topDomain == "" {
			continue
		}
		domains[topDomain] = struct{}{}
		domains[domain] = struct{}{}
	}
	list, err := script_repo.Domain().List(ctx, id)
	if err != nil {
		return err
	}
	domainMap := make(map[string]*script_entity.ScriptDomain)
	for _, v := range list {
		domainMap[v.Domain] = v
	}
	for domain := range domains {
		result, ok := domainMap[domain]
		if !ok {
			e := &script_entity.ScriptDomain{
				Domain:        domain,
				DomainReverse: utils.StringReverse(domain),
				ScriptID:      id,
				ScriptCodeID:  codeID,
				Status:        consts.ACTIVE,
				Createtime:    time.Now().Unix(),
			}
			if err := script_repo.Domain().Create(ctx, e); err != nil {
				logger.Ctx(ctx).Error("Create", zap.Error(err), zap.Int64("script_id", id), zap.String("domain", domain))
			}
		} else if result.Status != consts.ACTIVE {
			result.Status = consts.ACTIVE
			if err := script_repo.Domain().Update(ctx, result); err != nil {
				logger.Ctx(ctx).Error("Update", zap.Error(err), zap.Int64("script_id", id), zap.String("domain", domain))
			}
		}
		delete(domainMap, domain)
	}
	for _, v := range domainMap {
		if v.Status == consts.ACTIVE {
			if err := script_repo.Domain().Delete(ctx, v.ID); err != nil {
				logger.Ctx(ctx).Error("Delete", zap.Error(err), zap.Int64("script_id", id), zap.String("domain", v.Domain))
			}
		}
	}
	return nil
}

// 解析meta中的域名信息
// 返回格式为: 顶级域名, 原始域名
func (s *Script) parseMatchDomain(meta string) (string, string) {
	reg := regexp.MustCompile("(.+?://|^)(.+?)(/|$)")
	ret := reg.FindStringSubmatch(meta)
	if len(ret) == 0 || ret[2] == "" {
		return "", ""
	}
	if ret[2] == "*" {
		return "*", "*"
	}
	if ret[2][0] == '*' {
		ret[2] = ret[2][1:]
	}
	if ret[2][0] == '.' {
		ret[2] = ret[2][1:]
	}
	domain, err := publicsuffix.Domain(ret[2])
	if err != nil {
		return "", ""
	}
	if domain[0] == '*' {
		return "*", "*"
	}
	return domain, ret[2]
}
