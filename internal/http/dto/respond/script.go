package respond

import (
	"encoding/json"

	entity2 "github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

type Script struct {
	*User
	Script       *ScriptCode `json:"script"`
	ID           int64       `json:"id"`
	PostId       int64       `json:"post_id"`
	UserId       int64       `json:"user_id"`
	IsManager    bool        `json:"is_manager"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Status       int64       `json:"status"`
	Score        int64       `json:"score"`
	ScoreNum     int64       `json:"score_num"`
	Type         int         `json:"type"`
	Public       int         `json:"public"`
	Unwell       int         `json:"unwell"`
	TodayInstall int64       `json:"today_install"`
	TotalInstall int64       `json:"total_install"`
	Createtime   int64       `json:"createtime"`
	Updatetime   int64       `json:"updatetime"`
}

type ScriptSetting struct {
	SyncUrl       string `json:"sync_url"`
	ContentUrl    string `json:"content_url"`
	DefinitionUrl string `json:"definition_url"`
	SyncMode      int    `json:"sync_mode"`
}

type ScriptInfo struct {
	*Script
	Setting *ScriptSetting `json:"setting,omitempty"`
	Content string         `json:"content" form:"content"`
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
	Definition string      `json:"definition,omitempty"`
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

func ToScriptScore(user *User, score *entity2.ScriptScore) *ScriptScore {
	return &ScriptScore{
		User:       user,
		ID:         score.ID,
		UserId:     score.UserId,
		ScriptId:   score.ScriptId,
		Score:      score.Score,
		Message:    score.Message,
		Createtime: score.Createtime,
		Updatetime: score.Updatetime,
	}
}

func ToScript(user *User, script *entity2.Script, code *ScriptCode) *Script {
	ret := &Script{
		User:        user,
		Script:      code,
		ID:          script.ID,
		PostId:      script.PostId,
		UserId:      script.UserId,
		Name:        script.Name,
		Description: script.Description,
		Status:      script.Status,
		Type:        script.Type,
		Public:      script.Public,
		Unwell:      script.Unwell,
		Createtime:  script.Createtime,
		Updatetime:  script.Updatetime,
	}
	if user != nil {
		ret.IsManager = user.UID == script.UserId
	}
	return ret
}

func ToScriptInfo(user *User, script *entity2.Script, code *ScriptCode) *ScriptInfo {
	ret := &ScriptInfo{
		Script:  ToScript(user, script, code),
		Content: script.Content,
	}
	if user.UID == script.UserId {
		ret.Setting = &ScriptSetting{
			SyncUrl:       script.SyncUrl,
			ContentUrl:    script.ContentUrl,
			DefinitionUrl: script.DefinitionUrl,
			SyncMode:      script.SyncMode,
		}
	}
	return ret
}

func ToScriptCode(user *User, code *entity2.ScriptCode) *ScriptCode {
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
		User:       user,
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
