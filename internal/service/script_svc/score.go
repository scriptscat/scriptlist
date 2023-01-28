package script_svc

import (
	"context"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"time"
)

type ScoreSvc interface {
	// PutScore 脚本评分
	PutScore(ctx context.Context, req *api.PutScoreRequest) (*api.PutScoreResponse, error)
	// ScoreList 获取脚本评分列表
	ScoreList(ctx context.Context, req *api.ScoreListRequest) (*api.ScoreListResponse, error)
	// SelfScore 用于获取自己对脚本的评价
	SelfScore(ctx context.Context, req *api.SelfScoreRequest) (*api.SelfScoreResponse, error)
	// DelScore 用于删除脚本的评价，注意，只有管理员才有权限删除评价
	DelScore(ctx context.Context, req *api.DelScoreRequest) (*api.DelScoreResponse, error)
}

type scoreSvc struct {
}

var defaultScore = &scoreSvc{}

func Score() ScoreSvc {
	return defaultScore
}

// PutScore 脚本评分
func (s *scoreSvc) PutScore(ctx context.Context, req *api.PutScoreRequest) (*api.PutScoreResponse, error) {
	//评分模块的业务是这样的：
	//用户可以评分多次，评分后计算分数到统计表里（分数更新后需要重新计算分数），所以要判断之前是否有评分，另外必须验证了邮箱才能评分。
	//第一次评分完成后发送一封邮件通知给作者

	//判断用户邮箱是否验证
	//查询脚本是否存在
	//查询脚本状态，判断脚本是不是允许评分
	//判断用户有没有评价过,查询score表，查看是否有记录
	//更新用户对脚本的评价信息

	// 获取用户的id,获取用户邮箱验证状态，如果未验证，则抛出错误
	uid := auth_svc.Auth().Get(ctx).UID
	scriptId := req.ID

	if auth_svc.Auth().Get(ctx).EmailVerified == false {
		return nil, i18n.NewError(ctx, code.UserEmailNotVerified)
	}

	//判断脚本的状态，可用后下一步
	find, err2 := script_repo.Script().Find(ctx, scriptId)
	if err2 != nil {
		return nil, err2
	}
	if err := find.CheckOperate(ctx); err != nil {
		return nil, err
	}

	//获取用户的评分
	score := req.Score
	//判断用户有没有评价过,查询score表，查看是否有记录

	_, err := script_repo.ScriptScore().FindByUser(ctx, uid, scriptId)

	if err != nil {
		//不存在记录，创建一条记录
		err := script_repo.ScriptScore().Create(ctx, &entity.ScriptScore{
			UserID:     uid,
			ScriptID:   scriptId,
			Score:      score,
			Message:    req.Message,
			Createtime: time.Now().Unix(),
			Updatetime: time.Now().Unix(),
		})
		if err != nil {
			return nil, err
		}
		//给脚本作者发一个信息
		//等待一之处理
		//notice_svc.Notice().Send(ctx, uid)
	}

	//更新用户的评价信息
	err = script_repo.ScriptScore().Update(ctx, &entity.ScriptScore{Message: req.Message,
		Score: score})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ScoreList 获取脚本评分列表
func (s *scoreSvc) ScoreList(ctx context.Context, req *api.ScoreListRequest) (*api.ScoreListResponse, error) {
	//获取脚本的评分列表
	list, total, err := script_repo.ScriptScore().ScoreList(ctx, req.ScriptID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	//循环判断用户是否被封禁，过滤这些内容
	//resp := [...]*api.ScrScore{}
	resp := make([]*api.ScrScore, len(list))
	for i, v := range list {
		user, _ := user_repo.User().Find(ctx, v.UserID)

		info := user.UserInfo()
		//resp[i] = vo.ToScriptScore(info, v)
		resp[i] = &api.ScrScore{
			UserInfo: info,
			Score:    v,
		}
	}

	return &api.ScoreListResponse{
		PageResponse: httputils.PageResponse[*api.ScrScore]{
			List:  resp,
			Total: total,
		},
	}, nil

}

// SelfScore 用于获取自己对脚本的评价
func (s *scoreSvc) SelfScore(ctx context.Context, req *api.SelfScoreRequest) (*api.SelfScoreResponse, error) {
	//获取用户的id
	uid := auth_svc.Auth().Get(ctx).UID
	ret, err := script_repo.ScriptScore().FindByUser(ctx, uid, req.ScriptId)
	if err != nil {
		return nil, err
	}
	return &api.SelfScoreResponse{
		SelfScore: ret,
	}, nil

}

// DelScore 用于删除脚本的评价，注意，只有管理员才有权限删除评价
func (s *scoreSvc) DelScore(ctx context.Context, req *api.DelScoreRequest) (*api.DelScoreResponse, error) {
	admin := auth_svc.Auth().Get(ctx).AdminLevel.IsAdmin(model.SuperModerator)
	if !admin {
		return nil, i18n.NewError(ctx, code.UserNotPermission)
	}
	//删除评价
	err := script_repo.ScriptScore().Delete(ctx, req.ScoreId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
