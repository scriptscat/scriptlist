package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Score struct {
	ID         int64   `gorm:"column:id;type:bigint(20);not null;primary_key"`
	UserID     int64   `gorm:"column:user_id;type:bigint(20);index:user_id,unique;index:user_script,unique;index:user"`
	ScriptID   int64   `gorm:"column:script_id;type:bigint(20);index:user_id,unique;index:user_script,unique;index:script_id;index:script"`
	Score      float64 `gorm:"column:score;type:double"`
	Message    string  `gorm:"column:message;type:longtext"`
	Createtime int64   `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64   `gorm:"column:updatetime;type:bigint(20)"`
	State      int64   `gorm:"column:state;type:int(10);default:1"`
}

type ScrScore struct {
	user_entity.UserInfo
	Score *entity.ScriptScore
}

// PutScoreRequest 脚本评分
type PutScoreRequest struct {
	mux.Meta `path:"/scripts/:id/score" method:"PUT"`
	ID       int64   `uri:"id" binding:"required"`
	Message  string  `json:"message" binding:"required"`
	Score    float64 `json:"score" binding:"required,number,min=0,max=50"`
}

type PutScoreResponse struct {
}

// ScoreListRequest 获取脚本评分列表
type ScoreListRequest struct {
	mux.Meta              `path:"/scripts/:id/score" method:"GET"`
	ScriptID              int64 `uri:"id" binding:"required"`
	httputils.PageRequest `json:",inline"`
}

type ScoreListResponse struct {
	httputils.PageResponse[*ScrScore] `json:",inline"`
}

// SelfScoreRequest 用于获取自己对脚本的评价
type SelfScoreRequest struct {
	mux.Meta `path:"/scripts/:scriptId/score/self" method:"GET"`
	ScriptId int64 `uri:"scriptId" binding:"required"`
}
type SelfScoreResponse struct {
	SelfScore *entity.ScriptScore `json:",inline"`
}

// DelScoreRequest 用于删除脚本的评价，注意，只有管理员才有权限删除评价
type DelScoreRequest struct {
	mux.Meta `path:"/scripts/:scriptId/score/:scoreId" method:"DELETE"`
	ScriptId int64 `uri:"script" binding:"required"`
	ScoreId  int64 `uri:"scoreId" binding:"required"`
}

type DelScoreResponse struct {
}
