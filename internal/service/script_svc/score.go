package script_svc

import (
	"context"
	"net/http"
	"time"

	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
	"go.uber.org/zap"
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
	ReplyScore(ctx context.Context, req *api.ReplyScoreRequest) (*api.ReplyScoreResponse, error)
}

type scoreSvc struct {
}

func (s *scoreSvc) ReplyScore(ctx context.Context, req *api.ReplyScoreRequest) (*api.ReplyScoreResponse, error) {
	commentID := req.CommentID //被评论的评分
	scriptId := req.ScriptId   //脚本id
	//判断脚本的状态，可用后下一步
	script, err2 := script_repo.Script().Find(ctx, scriptId)
	if err2 != nil {
		return nil, err2
	}
	if err := script.CheckOperate(ctx); err != nil {
		return nil, err
	}
	//判断评论的状态，可用后下一步
	score, err := script_repo.ScriptScore().Find(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if score == nil {
		return nil, httputils.NewError(http.StatusNotFound, -1, "无法找到评分信息")
	}
	reply, err := script_repo.ScriptScore().FindReplayByComment(ctx, commentID, scriptId)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		//不存在记录，创建一条记录
		err := script_repo.ScriptScore().CreateReplayByComment(ctx, &script_entity.ScriptScoreReply{
			CommentID:  commentID,
			ScriptID:   scriptId,
			Message:    req.Message,
			Createtime: time.Now().Unix(),
			Updatetime: time.Now().Unix(),
		})
		if err != nil {
			return nil, err
		}
		//给用户发一个信息
		if err := notice_svc.Notice().Send(ctx, score.UserID, notice_svc.ScriptScoreReplyTemplate,
			notice_svc.WithFrom(script.UserID), notice_svc.WithParams(&template.ScriptReplyScore{
				ScriptID: scriptId,
				Name:     script.Name,
				Content:  req.Message,
			})); err != nil {
			// 发送失败不影响主要流程, 只记录错误
			logger.Ctx(ctx).Error("作者回复通知失败", zap.Int64("script", scriptId), zap.Int64("commentID", commentID), zap.Error(err))
		}
		return &api.ReplyScoreResponse{}, nil
	}
	reply.Message = req.Message
	reply.Updatetime = time.Now().Unix()
	err = script_repo.ScriptScore().UpdateReplayByComment(ctx, reply)
	if err != nil {
		return nil, err
	}
	return &api.ReplyScoreResponse{}, nil
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

	//判断脚本的状态，可用后下一步
	script, err2 := script_repo.Script().Find(ctx, scriptId)
	if err2 != nil {
		return nil, err2
	}
	if err := script.CheckOperate(ctx); err != nil {
		return nil, err
	}

	//判断用户有没有评价过,查询score表，查看是否有记录
	score, err := script_repo.ScriptScore().FindByUser(ctx, uid, scriptId)
	if err != nil {
		return nil, err
	}
	if score == nil {
		//不存在记录，创建一条记录
		var InsertScore = &script_entity.ScriptScore{
			UserID:     uid,
			ScriptID:   scriptId,
			Score:      req.Score,
			Message:    req.Message,
			Createtime: time.Now().Unix(),
			Updatetime: time.Now().Unix(),
		}
		err := script_repo.ScriptScore().Create(ctx, InsertScore)
		if err != nil {
			return nil, err
		}
		//给脚本作者发一个信息
		if err := notice_svc.Notice().Send(ctx, script.UserID, notice_svc.ScriptScoreTemplate,
			notice_svc.WithFrom(uid), notice_svc.WithParams(&template.ScriptScore{
				ScriptID: scriptId,
				Name:     script.Name,
				Username: auth_svc.Auth().Get(ctx).Username,
				Score:    int(req.Score / 10),
			})); err != nil {
			// 发送失败不影响主要流程, 只记录错误
			logger.Ctx(ctx).Error("评分通知作者失败", zap.Int64("script", scriptId), zap.Int64("user", uid), zap.Error(err))
		}
		// 进入统计
		if err := script_repo.ScriptStatistics().IncrScore(ctx, scriptId, req.Score, 1); err != nil {
			logger.Ctx(ctx).Error("评分统计失败", zap.Int64("script", scriptId), zap.Int64("user", uid), zap.Error(err))
		}
		return &api.PutScoreResponse{ID: InsertScore.ID}, nil
	}
	// 存在记录,但是状态不是激活状态,可能已经被管理员删除了,禁止再次评论
	if score.State != consts.ACTIVE {
		return nil, i18n.NewError(ctx, code.ScriptScoreDeleted)
	}

	//更新用户的评价信息
	oldScore := score.Score
	score.Message = req.Message
	score.Score = req.Score
	score.Updatetime = time.Now().Unix()
	err = script_repo.ScriptScore().Update(ctx, score)
	if err != nil {
		return nil, err
	}
	// 更新统计信息,只是更新分数,不更改评分人数
	if err := script_repo.ScriptStatistics().IncrScore(ctx, scriptId, req.Score-oldScore, 0); err != nil {
		logger.Ctx(ctx).Error("评分统计失败", zap.Int64("script", scriptId), zap.Int64("user", uid), zap.Error(err))
	}
	return &api.PutScoreResponse{ID: score.ID}, nil
}

// ScoreList 获取脚本评分列表
func (s *scoreSvc) ScoreList(ctx context.Context, req *api.ScoreListRequest) (*api.ScoreListResponse, error) {
	//获取脚本的评分列表
	list, total, err := script_repo.ScriptScore().ScoreList(ctx, req.ScriptID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := make([]*api.Score, len(list))
	for i, v := range list {
		resp[i], _ = s.ToScore(ctx, v)
	}
	return &api.ScoreListResponse{
		PageResponse: httputils.PageResponse[*api.Score]{
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
	if ret == nil {
		return nil, i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ScriptScoreNotFound)
	}
	resp, err := s.ToScore(ctx, ret)
	if err != nil {
		return nil, err
	}
	return &api.SelfScoreResponse{
		Score: resp,
	}, nil
}

func (s *scoreSvc) ToScore(ctx context.Context, score *script_entity.ScriptScore) (*api.Score, error) {
	user, err := user_repo.User().Find(ctx, score.UserID)
	if err != nil {
		return nil, err
	}
	var message string
	comment, err := script_repo.ScriptScore().FindReplayByComment(ctx, score.ID, score.ScriptID)
	if err != nil {
		return nil, err
	}
	if comment != nil {
		message = comment.Message
	}
	return &api.Score{
		UserInfo:      user.UserInfo(),
		ID:            score.ID,
		ScriptID:      score.ScriptID,
		Score:         score.Score,
		Message:       score.Message,
		Createtime:    score.Createtime,
		Updatetime:    score.Updatetime,
		AuthorMessage: message,
		State:         score.State,
	}, nil
}

// DelScore 用于删除脚本的评价，注意，只有管理员才有权限删除评价
func (s *scoreSvc) DelScore(ctx context.Context, req *api.DelScoreRequest) (*api.DelScoreResponse, error) {
	score, err := script_repo.ScriptScore().Find(ctx, req.ScoreId)
	if err != nil {
		return nil, err
	}
	if score == nil {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptScoreNotFound)
	}
	if score.ScriptID != req.ScriptId {
		return nil, i18n.NewNotFoundError(ctx, code.ScriptScoreNotFound)
	}
	//删除评价
	err = script_repo.ScriptScore().Delete(ctx, req.ScoreId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
