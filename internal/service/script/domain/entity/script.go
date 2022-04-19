package entity

import (
	"net/http"
	"time"

	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	cnt2 "github.com/scriptscat/scriptlist/internal/service/user/cnt"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/vo"
)

const (
	USERSCRIPT_TYPE = iota + 1
	SUBSCRIBE_TYPE
	LIBRARY_TYPE
)

const (
	PUBLIC_SCRIPT = iota + 1
	PRIVATE_SCRIPT
)

type Script struct {
	ID          int64  `gorm:"column:id" json:"id" form:"id"`
	PostID      int64  `gorm:"column:post_id;index:post_id,unique" json:"post_id" form:"post_id"`
	UserID      int64  `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	Name        string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Description string `gorm:"column:description;type:text" json:"description" form:"description"`
	Content     string `gorm:"column:content;type:mediumtext" json:"content" form:"content"`
	Type        int    `gorm:"column:type;type:bigint;index:script_type;not null;default:1" json:"type"`
	Public      int    `gorm:"column:public;not null;default:1" json:"public"`
	// 不适内容
	Unwell        int    `gorm:"column:unwell;not null;default:2" json:"unwell"`
	SyncUrl       string `gorm:"column:sync_url;type:text;index:sync_url,length:128" json:"sync_url"`
	ContentUrl    string `gorm:"column:content_url;type:text;index:content_url,length:128" json:"content_url"`
	DefinitionUrl string `gorm:"column:definition_url;type:text;index:definition_url,length:128" json:"definition_url"`
	SyncMode      int    `gorm:"column:sync_mode;type:tinyint(2)"`
	// 归档
	Archive    int32 `gorm:"column:archive;type:tinyint(2)" json:"archive"`
	Status     int64 `gorm:"column:status" json:"status" form:"status"`
	Createtime int64 `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64 `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}

func (s *Script) TableName() string {
	return "cdb_tampermonkey_script"
}

func (s *Script) checkStatus() error {
	switch s.Status {
	case cnt.ACTIVE:
		return nil
	case cnt.AUDIT:
		return errs.ErrScriptAudit
	}
	return errs.ErrScriptNotFound
}

func (s *Script) checkUser(user *vo.User) error {
	if user.IsAdmin == cnt2.Admin || user.IsAdmin == cnt2.SuperModerator || user.UID == s.UserID {
		return nil
	}
	return errs.NewError(http.StatusUnauthorized, 1000, "没有权限操作")
}

func (s *Script) SetArchive(user *vo.User, archive int32) error {
	if err := s.checkUser(user); err != nil {
		return err
	}
	s.Archive = archive
	s.Updatetime = time.Now().Unix()
	return nil
}

func (s *Script) Delete(user *vo.User) error {
	if err := s.checkUser(user); err != nil {
		return err
	}
	s.Status = cnt.DELETE
	s.Updatetime = time.Now().Unix()
	return nil
}

func (s *Script) CreateScriptCode(uid int64, req *request.CreateScript) (*ScriptCode, error) {
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	if s.UserID != uid {
		return nil, errs.ErrScriptForbidden
	}
	if s.Archive != 0 {
		return nil, errs.NewError(http.StatusUnauthorized, 1000, "已归档的脚本不能操作")
	}
	s.Content = req.Content
	s.Public = req.Public
	s.Unwell = req.Unwell
	s.Updatetime = time.Now().Unix()
	return &ScriptCode{
		UserId:     uid,
		Version:    time.Now().Format("20060102150405"),
		Changelog:  req.Changelog,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}, nil
}

func (s *Script) AddScore(uid int64, score *ScriptScore, msg *request.Score) (*ScriptScore, error) {
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	if score == nil {
		score = &ScriptScore{
			UserId:     uid,
			ScriptId:   s.ID,
			State:      cnt.ACTIVE,
			Createtime: time.Now().Unix(),
		}
	} else {
		if score.State != cnt.ACTIVE {
			return nil, errs.NewBadRequestError(1000, "评分已被删除")
		}
	}
	score.Score = msg.Score
	score.Message = msg.Message
	score.Updatetime = time.Now().Unix()
	return score, nil
}

func (s *Script) DeleteScore(score *ScriptScore) error {
	if err := s.checkStatus(); err != nil {
		return err
	}
	if score.State != cnt.ACTIVE {
		return errs.NewBadRequestError(1000, "评分已被删除")
	}
	score.State = cnt.DELETE
	score.Updatetime = time.Now().Unix()
	return nil
}

func (s *Script) SetUnwell(user *vo.User) error {
	if err := s.checkUser(user); err != nil {
		return err
	}
	s.Unwell = 1
	s.Updatetime = time.Now().Unix()
	return nil
}

func (s *Script) SetUnpublic(user *vo.User) error {
	if err := s.checkUser(user); err != nil {
		return err
	}
	s.Public = 0
	s.Updatetime = time.Now().Unix()
	return nil
}

type ScriptCode struct {
	ID         int64  `gorm:"column:id" json:"id" form:"id"`
	UserId     int64  `gorm:"column:user_id;index:user_id" json:"user_id" form:"user_id"`
	ScriptId   int64  `gorm:"column:script_id;index:script_id" json:"script_id" form:"script_id"`
	Code       string `gorm:"column:code;type:mediumtext" json:"code" form:"code"`
	Meta       string `gorm:"column:meta;type:text" json:"meta" form:"meta"`
	MetaJson   string `gorm:"column:meta_json;type:text" json:"meta_json" form:"meta_json"`
	Version    string `gorm:"column:version;type:varchar(255)" json:"version" form:"version"`
	Changelog  string `gorm:"column:changelog;type:text" json:"changelog" form:"changelog"`
	Status     int64  `gorm:"column:status" json:"status" form:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint" json:"createtime" form:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint" json:"updatetime" form:"updatetime"`
}

func (s *ScriptCode) TableName() string {
	return "cdb_tampermonkey_script_code"
}

type LibDefinition struct {
	ID         int64  `gorm:"column:id" json:"id" form:"id"`
	UserId     int64  `gorm:"column:user_id;index:user_id;not null" json:"user_id" form:"user_id"`
	ScriptId   int64  `gorm:"column:script_id;index:script_id;not null" json:"script_id" form:"script_id"`
	CodeId     int64  `gorm:"column:code_id;index:code_id;not null" json:"code_id" form:"code_id"`
	Definition string `gorm:"column:definition;not null" json:"definition" form:"definition"`
	Createtime int64  `gorm:"column:createtime;type:bigint" json:"createtime" form:"createtime"`
}
