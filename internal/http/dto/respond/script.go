package respond

import (
	"encoding/json"

	entity2 "github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/domain/user/entity"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

type Script struct {
	*User
	Script       *ScriptCode `json:"script"`
	ID           int64       `json:"id"`
	PostId       int64       `json:"post_id"`
	UserId       int64       `json:"user_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Status       int64       `json:"status"`
	Score        int64       `json:"score"`
	ScoreNum     int64       `json:"score_num"`
	TodayInstall int64       `json:"today_install"`
	TotalInstall int64       `json:"total_install"`
	Createtime   int64       `json:"createtime"`
	Updatetime   int64       `json:"updatetime"`
}

type ScriptInfo struct {
	*Script
	Content string `json:"content" form:"content"`
}

type ScriptCode struct {
	*User
	ID         int64       `json:"id" form:"id"`
	UserId     int64       `json:"user_id" form:"user_id"`
	Meta       string      `json:"meta,omitempty" form:"meta"`
	MetaJson   interface{} `json:"meta_json"`
	ScriptId   int64       `json:"script_id" form:"script_id"`
	Version    string      `json:"version" form:"version"`
	Changelog  string      `json:"changelog" form:"changelog"`
	Status     int64       `json:"status" form:"status"`
	Createtime int64       `json:"createtime" form:"createtime"`
	Code       string      `json:"code,omitempty" form:"code"`
}

type ScriptScore struct {
	*User
	ID       int64 `gorm:"column:id" json:"id" form:"id"`
	UserId   int64 `gorm:"column:user_id;index:user_script,unique;index:user" json:"user_id" form:"user_id"`
	ScriptId int64 `gorm:"column:script_id;index:user_script,unique;index:script" json:"script_id" form:"script_id"`
	// 评分,五星制,50
	Score int64 `gorm:"column:score" json:"score" form:"score"`
	// 评分原因
	Message    string `gorm:"column:message;type:varchar(255)" json:"message" form:"message"`
	Createtime int64  `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64  `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

func ToScriptScore(user *entity.User, score *entity2.ScriptScore) *ScriptScore {
	return &ScriptScore{
		User:       ToUser(user),
		ID:         score.ID,
		UserId:     score.UserId,
		ScriptId:   score.ScriptId,
		Score:      score.Score,
		Message:    score.Message,
		Createtime: score.Createtime,
		Updatetime: score.Updatetime,
	}
}

func ToScript(user *entity.User, scriptInfo *entity2.Script, script *ScriptCode) *Script {
	return &Script{
		User:        ToUser(user),
		ID:          scriptInfo.ID,
		PostId:      scriptInfo.PostId,
		UserId:      scriptInfo.UserId,
		Name:        scriptInfo.Name,
		Description: scriptInfo.Description,
		Script:      script,
		Status:      scriptInfo.Status,
		Createtime:  scriptInfo.Createtime,
		Updatetime:  scriptInfo.Updatetime,
	}
}

func ToScriptInfo(user *entity.User, script *entity2.Script, code *ScriptCode) *ScriptInfo {
	return &ScriptInfo{
		Script:  ToScript(user, script, code),
		Content: script.Content,
	}
}

func ToScriptCode(user *entity.User, code *entity2.ScriptCode) *ScriptCode {
	data := make(map[string]interface{})
	_ = json.Unmarshal([]byte(code.MetaJson), &data)
	domains := make(map[string]struct{})
	if _, ok := data["match"]; ok {
		for _, u := range data["match"].([]interface{}) {
			domain := utils.ParseMetaDomain(u.(string))
			if domain != "" {
				domains[domain] = struct{}{}
			}
		}
	}
	if _, ok := data["include"]; ok {
		for _, u := range data["include"].([]interface{}) {
			domain := utils.ParseMetaDomain(u.(string))
			if domain != "" {
				domains[domain] = struct{}{}
			}
		}
	}
	data["domains"] = domains
	return &ScriptCode{
		User:       ToUser(user),
		ID:         code.ID,
		UserId:     code.UserId,
		ScriptId:   code.ScriptId,
		Meta:       code.Meta,
		MetaJson:   data,
		Version:    code.Version,
		Changelog:  code.Changelog,
		Status:     code.Status,
		Createtime: code.Createtime,
	}
}
