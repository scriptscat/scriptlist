package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	respond2 "github.com/scriptscat/scriptlist/internal/interfaces/api/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/internal/service/issue/application"
	entity2 "github.com/scriptscat/scriptlist/internal/service/issue/domain/entity"
	service4 "github.com/scriptscat/scriptlist/internal/service/notify/service"
	service2 "github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/user/service"
	"github.com/scriptscat/scriptlist/pkg/httputils"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

type ScriptIssue struct {
	scriptSvc     service2.Script
	userSvc       service.User
	notifySvc     service4.Sender
	issueSvc      application.Issue
	issueWatchSvc application.ScriptIssueWatch
}

func NewScriptIssue(scriptSvc service2.Script, userSvc service.User, notifySvc service4.Sender, issueSvc application.Issue, issueWatchSvc application.ScriptIssueWatch) *ScriptIssue {
	return &ScriptIssue{
		scriptSvc:     scriptSvc,
		userSvc:       userSvc,
		notifySvc:     notifySvc,
		issueSvc:      issueSvc,
		issueWatchSvc: issueWatchSvc,
	}
}

func (s *ScriptIssue) getScriptId(c *gin.Context) int64 {
	return utils.StringToInt64(c.Param("script"))
}

func (s *ScriptIssue) getIssueId(c *gin.Context) int64 {
	return utils.StringToInt64(c.Param("issue"))
}

func (s *ScriptIssue) list(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		script := s.getScriptId(c)
		page := &request.Pages{}
		if err := c.ShouldBind(page); err != nil {
			return err
		}
		var labels []string
		if c.Query("labels") != "" {
			labels = strings.Split(c.Query("labels"), ",")
		}
		list, total, err := s.issueSvc.List(script, c.Query("keyword"), labels, utils.StringToInt(c.Query("status")), page)
		if err != nil {
			return err
		}
		ret := make([]interface{}, len(list))
		for k, v := range list {
			u, _ := s.userSvc.UserInfo(v.UserID)
			ret[k] = respond2.ToIssue(u, v)
		}
		return &httputils.List{
			List:  ret,
			Total: total,
		}
	})
}

func (s *ScriptIssue) post(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		scriptId := s.getScriptId(c)
		user, _ := token.UserInfo(c)
		req := request.Issue{}
		if err := c.ShouldBind(&req); err != nil {
			return err
		}
		script, err := s.scriptSvc.Info(scriptId)
		if err != nil {
			return err
		}
		if script.Archive != 0 {
			return errs.ErrScriptArchived
		}
		var labels []string
		if req.Label != "" {
			labels = strings.Split(req.Label, ",")
		}
		issue, err := s.issueSvc.Issue(scriptId, user.UID, req.Title, req.Content, labels)
		if err != nil {
			return err
		}
		return respond2.ToIssue(user, issue)
	})
}

func (s *ScriptIssue) get(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		issue, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		u, _ := s.userSvc.UserInfo(issue.UserID)
		return respond2.ToIssue(u, issue)
	})
}

func (s *ScriptIssue) put(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := token.UserId(c)
		req := request.Issue{}
		if err := c.ShouldBind(&req); err != nil {
			return err
		}
		if _, _, err := s.isOperate(issueId, user); err != nil {
			return err
		}
		return s.issueSvc.UpdateIssue(issueId, user, req.Title, req.Content)
	})
}

func (s *ScriptIssue) del(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := token.UserId(c)
		if _, _, err := s.isOperate(issueId, user); err != nil {
			return err
		}
		return s.issueSvc.DelIssue(issueId, user)
	})
}

func (s *ScriptIssue) putLabels(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := token.UserId(c)
		if _, _, err := s.isOperate(issueId, user); err != nil {
			return err
		}
		labels := strings.Split(c.PostForm("labels"), ",")
		comment, err := s.issueSvc.Label(issueId, user, labels)
		if err != nil {
			return err
		}
		return comment
	})
}

func (s *ScriptIssue) isOperate(issueId, user int64) (*entity2.ScriptIssue, *entity.Script, error) {
	issue, err := s.issueSvc.GetIssue(issueId)
	if err != nil {
		return nil, nil, err
	}
	script, err := s.scriptSvc.Info(issue.ScriptID)
	if err != nil {
		return nil, nil, err
	}
	if script.Archive != 0 {
		return nil, nil, errs.ErrScriptArchived
	}
	if issue.UserID != user {
		if script.UserId != user {
			return nil, nil, errs.NewError(http.StatusForbidden, 1001, "没有权限删除")
		}
	}
	return issue, script, nil
}

func (s *ScriptIssue) open(open bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		httputils.Handle(c, func() interface{} {
			issueId := s.getIssueId(c)
			user, _ := token.UserInfo(c)
			_, err := s.userSvc.UserInfo(user.UID)
			if err != nil {
				return err
			}
			_, _, err = s.isOperate(issueId, user.UID)
			if err != nil {
				return err
			}
			var comment *entity2.ScriptIssueComment
			if open {
				comment, err = s.issueSvc.Open(issueId, user.UID)
			} else {
				comment, err = s.issueSvc.Close(issueId, user.UID)
			}
			if err != nil {
				return err
			}
			return respond2.ToIssueComment(user, comment)
		})
	}
}

func (s *ScriptIssue) commentList(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		issue := s.getIssueId(c)
		page := &request.Pages{}
		if err := c.ShouldBind(page); err != nil {
			return err
		}
		list, err := s.issueSvc.CommentList(issue, page)
		if err != nil {
			return err
		}
		ret := make([]interface{}, len(list))
		for k, v := range list {
			u, _ := s.userSvc.UserInfo(v.UserID)
			ret[k] = respond2.ToIssueComment(u, v)
		}
		return &httputils.List{
			List:  ret,
			Total: int64(len(ret)),
		}
	})
}

func (s *ScriptIssue) getCommentId(c *gin.Context) int64 {
	return utils.StringToInt64(c.Param("comment"))
}

func (s *ScriptIssue) comment(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := token.UserInfo(c)
		issue, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		script, err := s.scriptSvc.Info(issue.ScriptID)
		if err != nil {
			return err
		}
		if script.Archive != 0 {
			return errs.ErrScriptArchived
		}
		content := c.PostForm("content")
		comment, err := s.issueSvc.Comment(issueId, user.UID, content)
		if err != nil {
			return err
		}
		watch, err := s.issueWatchSvc.IsWatch(issueId, user.UID)
		if err == nil && watch == 0 {
			_ = s.issueWatchSvc.Watch(issueId, user.UID)
		}
		return respond2.ToIssueComment(user, comment)
	})
}

func (s *ScriptIssue) commentUpdate(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		commentId := s.getCommentId(c)
		user, _ := token.UserInfo(c)
		comment, err := s.issueSvc.GetComment(commentId)
		if err != nil {
			return err
		}
		issue, err := s.issueSvc.GetIssue(comment.IssueID)
		if err != nil {
			return err
		}
		script, err := s.scriptSvc.Info(issue.ScriptID)
		if err != nil {
			return err
		}
		if script.Archive != 0 {
			return errs.ErrScriptArchived
		}
		return s.issueSvc.UpdateComment(commentId, user.UID, c.PostForm("content"))
	})
}

func (s *ScriptIssue) commentDel(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		commentId := s.getCommentId(c)
		user, _ := token.UserInfo(c)
		comment, err := s.issueSvc.GetComment(commentId)
		if err != nil {
			return err
		}
		issue, err := s.issueSvc.GetIssue(comment.IssueID)
		if err != nil {
			return err
		}
		script, err := s.scriptSvc.Info(issue.ScriptID)
		if err != nil {
			return err
		}
		if script.Archive != 0 {
			return errs.ErrScriptArchived
		}
		return s.issueSvc.DelComment(commentId, user.UID)
	})
}

func (s *ScriptIssue) iswatch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := token.UserId(c)
		issueId := s.getIssueId(c)
		_, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		watch, err := s.issueWatchSvc.IsWatch(issueId, uid)
		if err != nil {
			return err
		}
		return gin.H{
			"watch": watch,
		}
	})
}

func (s *ScriptIssue) watch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := token.UserId(c)
		issueId := s.getIssueId(c)
		_, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		return s.issueWatchSvc.Watch(issueId, uid)
	})
}

func (s *ScriptIssue) unwatch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := token.UserId(c)
		issueId := s.getIssueId(c)
		_, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		return s.issueWatchSvc.Unwatch(issueId, uid)
	})
}

func (s *ScriptIssue) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/scripts/:script/issues")
	rg.GET("", s.list)
	rg.POST("", token.UserAuth(true), s.post)
	rg.GET("/:issue", s.get)
	rgg := rg.Group("/:issue", token.UserAuth(true))
	rgg.PUT("", s.put)
	rgg.DELETE("", s.del)
	rgg.PUT("/labels", s.putLabels)
	rgg.PUT("/open", s.open(true))
	rgg.PUT("/close", s.open(false))

	rg.GET("/:issue/comment", s.commentList)
	rgg.POST("/comment", s.comment)
	rgg.PUT("/comment/:comment", s.commentUpdate)
	rgg.DELETE("/comment/:comment", s.commentDel)

	rgg.GET("/watch", s.iswatch)
	rgg.POST("/watch", s.watch)
	rgg.DELETE("/watch", s.unwatch)

}
