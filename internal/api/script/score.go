package script

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Score struct {
	user_entity.UserInfo
	ID                      int64  `json:"id"`
	ScriptID                int64  `json:"script_id"`
	Score                   int64  `json:"score"`
	Message                 string `json:"message"`
	AuthorMessage           string `json:"author_message"`
	AuthorMessageCreatetime int64  `json:"author_message_createtime"`
	Createtime              int64  `json:"createtime"`
	Updatetime              int64  `json:"updatetime"`
	State                   int64  `json:"state"`
}

// PutScoreRequest 脚本评分
type PutScoreRequest struct {
	mux.Meta `path:"/scripts/:id/score" method:"PUT"`
	ID       int64  `uri:"id" binding:"required"`
	Message  string `json:"message"`
	Score    int64  `json:"score" binding:"required,number,oneof=10 20 30 40 50"`
}

type PutScoreResponse struct {
	ID int64 `json:"id"`
}

// ScoreListRequest 获取脚本评分列表
type ScoreListRequest struct {
	mux.Meta              `path:"/scripts/:id/score" method:"GET"`
	httputils.PageRequest `json:",inline"`
	ScriptID              int64 `uri:"id" binding:"required"`
}

type ScoreListResponse struct {
	httputils.PageResponse[*Score] `json:",inline"`
}

// SelfScoreRequest 用于获取自己对脚本的评价
type SelfScoreRequest struct {
	mux.Meta `path:"/scripts/:id/score/self" method:"GET"`
	ScriptId int64 `uri:"id" binding:"required"`
}

type SelfScoreResponse struct {
	*Score `json:",inline"`
}

// DelScoreRequest 用于删除脚本的评价，注意，只有管理员才有权限删除评价
type DelScoreRequest struct {
	mux.Meta `path:"/scripts/:id/score/:scoreId" method:"DELETE"`
	ScriptId int64 `uri:"id" binding:"required"`
	ScoreId  int64 `uri:"scoreId" binding:"required"`
}

type ReplyScoreRequest struct {
	mux.Meta  `path:"/scripts/:id/commentReply" method:"PUT"`
	ScriptId  int64  `uri:"id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	CommentID int64  `json:"commentID" binding:"required"`
}
type ReplyScoreResponse struct {
}
type DelScoreResponse struct {
}

// ScoreStateRequest 获取脚本评分状态
type ScoreStateRequest struct {
	mux.Meta `path:"/scripts/:id/score/state" method:"GET"`
	ScriptId int64 `uri:"id" binding:"required"`
}

type ScoreStateResponse struct {
	// 每个评分的数量
	ScoreGroup map[int64]int64 `json:"score_group"`
	// 评分人数
	ScoreUserCount int64 `json:"score_user_count"`
}
