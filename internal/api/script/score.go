package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Score struct {
	user_entity.UserInfo
	ID         int64  `json:"id"`
	ScriptID   int64  `json:"script_id"`
	Score      int64  `json:"score"`
	Message    string `json:"message"`
	Createtime int64  `json:"createtime"`
	Updatetime int64  `json:"updatetime"`
	State      int64  `json:"state"`
}

// PutScoreRequest 脚本评分
type PutScoreRequest struct {
	mux.Meta `path:"/scripts/:id/score" method:"PUT"`
	ID       int64  `uri:"id" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Score    int64  `json:"score" binding:"required,number,min=0,max=50"`
}

type PutScoreResponse struct {
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

type DelScoreResponse struct {
}
