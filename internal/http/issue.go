package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	entity2 "github.com/scriptscat/scriptlist/internal/domain/issue/entity"
	service3 "github.com/scriptscat/scriptlist/internal/domain/issue/service"
	service4 "github.com/scriptscat/scriptlist/internal/domain/notify/service"
	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	service2 "github.com/scriptscat/scriptlist/internal/domain/script/service"
	"github.com/scriptscat/scriptlist/internal/domain/user/service"
	request2 "github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

type ScriptIssue struct {
	scriptSvc     service2.Script
	userSvc       service.User
	notifySvc     service4.Sender
	issueSvc      service3.Issue
	issueWatchSvc service3.ScriptIssueWatch
}

func NewScriptIssue(scriptSvc service2.Script, userSvc service.User, notifySvc service4.Sender, issueSvc service3.Issue, issueWatchSvc service3.ScriptIssueWatch) *ScriptIssue {
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
	handle(c, func() interface{} {
		script := s.getScriptId(c)
		page := request2.Pages{}
		if err := c.ShouldBind(&page); err != nil {
			return err
		}
		var labels []string
		if c.Query("labels") != "" {
			labels = strings.Split(c.Query("labels"), ",")
		}
		list, err := s.issueSvc.List(script, c.Query("keyword"), labels, utils.StringToInt(c.Query("status")), page)
		if err != nil {
			return err
		}
		ret := make([]*respond.Issue, len(list))
		for k, v := range list {
			u, _ := s.userSvc.UserInfo(v.UserID)
			ret[k] = respond.ToIssue(u, v)
		}
		return ret
	})
}

func (s *ScriptIssue) post(c *gin.Context) {
	handle(c, func() interface{} {
		scriptId := s.getScriptId(c)
		user, _ := selfinfo(c)
		req := request2.Issue{}
		if err := c.ShouldBind(&req); err != nil {
			return err
		}
		_, err := s.scriptSvc.Info(scriptId)
		if err != nil {
			return err
		}
		var labels []string
		if req.Label != "" {
			labels = strings.Split(req.Label, ",")
		}
		issue, err := s.issueSvc.Issue(scriptId, user.UID, req.Title, req.Content, labels)
		if err != nil {
			return err
		}
		return respond.ToIssue(user, issue)
	})
}

func (s *ScriptIssue) get(c *gin.Context) {
	handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		issue, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		u, _ := s.userSvc.UserInfo(issue.UserID)
		return respond.ToIssue(u, issue)
	})
}

func (s *ScriptIssue) put(c *gin.Context) {
	handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := userId(c)
		req := request2.Issue{}
		if err := c.ShouldBind(&req); err != nil {
			return err
		}
		return s.issueSvc.UpdateIssue(issueId, user, req.Title, req.Content)
	})
}

func (s *ScriptIssue) del(c *gin.Context) {
	handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := userId(c)
		if _, _, err := s.isOperate(issueId, user); err != nil {
			return err
		}
		return s.issueSvc.DelIssue(issueId, user)
	})
}

func (s *ScriptIssue) putLabels(c *gin.Context) {
	handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := userId(c)
		if _, _, err := s.isOperate(issueId, user); err != nil {
			return err
		}
		labels := strings.Split(c.PostForm("labels"), ",")
		return s.issueSvc.Label(issueId, user, labels)
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
	if issue.UserID != user {
		if script.UserId != user {
			return nil, nil, errs.NewError(http.StatusForbidden, 1001, "没有权限删除")
		}
	}
	return issue, script, nil
}

func (s *ScriptIssue) open(open bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		handle(c, func() interface{} {
			issueId := s.getIssueId(c)
			uid, _ := userId(c)
			_, err := s.userSvc.UserInfo(uid)
			if err != nil {
				return err
			}
			_, _, err = s.isOperate(issueId, uid)
			if err != nil {
				return err
			}
			if open {
				err = s.issueSvc.Open(issueId, uid)
			} else {
				err = s.issueSvc.Close(issueId, uid)
			}
			if err != nil {
				return err
			}
			return nil
		})
	}
}

func (s *ScriptIssue) commentList(c *gin.Context) {
	handle(c, func() interface{} {
		issue := s.getIssueId(c)
		page := request2.Pages{}
		if err := c.ShouldBind(&page); err != nil {
			return err
		}
		list, err := s.issueSvc.CommentList(issue, page)
		if err != nil {
			return err
		}
		ret := make([]*respond.IssueComment, len(list))
		for k, v := range list {
			u, _ := s.userSvc.UserInfo(v.UserID)
			ret[k] = respond.ToIssueComment(u, v)
		}
		return ret
	})
}

func (s *ScriptIssue) getCommentId(c *gin.Context) int64 {
	return utils.StringToInt64(c.Param("comment"))
}

func (s *ScriptIssue) comment(c *gin.Context) {
	handle(c, func() interface{} {
		issueId := s.getIssueId(c)
		user, _ := selfinfo(c)
		issue, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		_, err = s.scriptSvc.Info(issue.ScriptID)
		if err != nil {
			return err
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
		return comment
	})
}

func (s *ScriptIssue) commentUpdate(c *gin.Context) {
	handle(c, func() interface{} {
		commentId := s.getCommentId(c)
		user, _ := selfinfo(c)
		return s.issueSvc.UpdateComment(commentId, user.UID, c.PostForm("content"))
	})
}

func (s *ScriptIssue) commentDel(c *gin.Context) {
	handle(c, func() interface{} {
		commentId := s.getCommentId(c)
		user, _ := selfinfo(c)
		comment, err := s.issueSvc.GetComment(commentId)
		if err != nil {
			return err
		}
		if comment.UserID != user.UID {
			issue, err := s.issueSvc.GetIssue(comment.IssueID)
			if err != nil {
				return err
			}
			script, err := s.scriptSvc.Info(issue.ScriptID)
			if err != nil {
				return err
			}
			if script.UserId != user.UID {
				return errs.NewError(http.StatusForbidden, 1000, "没有权限删除评论")
			}
		}
		return s.issueSvc.DelComment(commentId)
	})
}

func (s *ScriptIssue) iswatch(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := userId(c)
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
	handle(c, func() interface{} {
		uid, _ := userId(c)
		issueId := s.getIssueId(c)
		_, err := s.issueSvc.GetIssue(issueId)
		if err != nil {
			return err
		}
		return s.issueWatchSvc.Watch(issueId, uid)
	})
}

func (s *ScriptIssue) unwatch(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := userId(c)
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
	rg.POST("", userAuth(true), s.post)
	rg.GET("/:issue", s.get)
	rgg := rg.Group("/:issue", userAuth(true))
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
