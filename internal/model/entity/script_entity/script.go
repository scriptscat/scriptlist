package script_entity

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
)

type Type int

const (
	UserscriptType Type = iota + 1 // 用户脚本
	SubscribeType                  // 订阅脚本
	LibraryType                    // 库
)

type Public int

const (
	PublicScript  Public = iota + 1 // 公开
	PrivateScript                   // 私有(半公开,只是不展示在列表中)
)

type UnwellContent int

const (
	Unwell UnwellContent = iota + 1 // 不适内容
	Well                            // 合适内容
)

type SyncMode int

const (
	SyncModeAuto   SyncMode = iota + 1 // 自动同步
	SyncModeManual                     // 手动同步
)

type ScriptArchive int

const (
	IsArchive ScriptArchive = iota + 1
	IsActive
)

type Script struct {
	ID            int64         `gorm:"column:id;type:bigint(20);not null;primary_key"`
	PostID        int64         `gorm:"column:post_id;type:bigint(20);index:post_id"`
	UserID        int64         `gorm:"column:user_id;type:bigint(20);index:user_id"`
	Name          string        `gorm:"column:name;type:varchar(255)"`
	Description   string        `gorm:"column:description;type:text"`
	Content       string        `gorm:"column:content;type:mediumtext"`
	Type          Type          `gorm:"column:type;type:bigint(20);default:1;not null;index:script_type"`
	Public        Public        `gorm:"column:public;type:bigint(20);default:1;not null"`
	Unwell        UnwellContent `gorm:"column:unwell;type:bigint(20);default:2;not null"`
	SyncUrl       string        `gorm:"column:sync_url;type:text;index:sync_url"`
	ContentUrl    string        `gorm:"column:content_url;type:text;index:content_url"`
	DefinitionUrl string        `gorm:"column:definition_url;type:text;index:definition_url"`
	SyncMode      SyncMode      `gorm:"column:sync_mode;type:tinyint(2)"`
	Archive       ScriptArchive `gorm:"column:archive;type:tinyint(2);default:2;not null"`
	Status        int64         `gorm:"column:status;type:bigint(20)"`
	Createtime    int64         `gorm:"column:createtime;type:bigint(20)"`
	Updatetime    int64         `gorm:"column:updatetime;type:bigint(20)"`
}

func (s *Script) TableName() string {
	return "cdb_tampermonkey_script"
}

// CheckOperate 检查是否可以操作
func (s *Script) CheckOperate(ctx context.Context) error {
	if s == nil {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ScriptNotFound)
	}
	if s.Status != consts.ACTIVE {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ScriptIsDelete)
	}
	return nil
}

// CheckPermission 检查操作权限
func (s *Script) CheckPermission(ctx context.Context, allowAdminLevel ...model.AdminLevel) error {
	if err := s.CheckOperate(ctx); err != nil {
		return err
	}
	user := auth_svc.Auth().Get(ctx)
	if s.UserID != user.UID {
		if len(allowAdminLevel) > 0 && user.AdminLevel.IsAdmin(allowAdminLevel[0]) {
			return nil
		}
		return i18n.NewErrorWithStatus(ctx, http.StatusForbidden, code.UserNotPermission)
	}
	return nil
}

// IsArchive 是否归档
func (s *Script) IsArchive(ctx context.Context) error {
	if err := s.CheckOperate(ctx); err != nil {
		return err
	}
	if s.Archive == IsArchive {
		return i18n.NewError(ctx, code.ScriptIsArchive)
	}
	return nil
}

type Code struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);index:user_id"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);index:script_id;index:script_version,unique"`
	Code       string `gorm:"column:code;type:mediumtext"`
	Meta       string `gorm:"column:meta;type:text"`
	MetaJson   string `gorm:"column:meta_json;type:text"`
	Version    string `gorm:"column:version;type:varchar(255);index:script_version,unique"`
	Changelog  string `gorm:"column:changelog;type:text"`
	Status     int64  `gorm:"column:status;type:tinyint(4)"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (s *Code) TableName() string {
	return "cdb_tampermonkey_script_code"
}

func (s *Code) UpdateCode(ctx context.Context, scriptCode string) (map[string][]string, error) {
	// 解析脚本的元数据
	scriptCodeStr, meta, err := parseCodeMeta(ctx, scriptCode)
	if err != nil {
		return nil, err
	}
	// 解析元数据
	metaJson := parseMetaToJson(meta)
	if len(metaJson["name"]) == 0 {
		return nil, i18n.NewError(ctx, code.ScriptNameIsEmpty)
	}
	if len(metaJson["description"]) == 0 {
		return nil, i18n.NewError(ctx, code.ScriptDescIsEmpty)
	}
	if len(metaJson["version"]) == 0 {
		return nil, i18n.NewError(ctx, code.ScriptVersionIsEmpty)
	}
	b, err := json.Marshal(metaJson)
	if err != nil {
		return nil, i18n.NewError(ctx, code.ScriptParseFailed)
	}
	s.Code = scriptCodeStr
	s.Meta = meta
	s.MetaJson = string(b)
	s.Version = metaJson["version"][0]
	return metaJson, nil
}

func (s *Code) Fields() string {
	return "id, user_id, script_id, meta, meta_json, version, changelog, status, createtime, updatetime"
}

// 解析脚本的元数据
func parseCodeMeta(ctx context.Context, scriptCode string) (string, string, error) {
	reg := regexp.MustCompile(`\/\/\s*==UserScript==([\s\S]+?)\/\/\s*==\/UserScript==`)
	ret := reg.FindString(scriptCode)
	if ret == "" {
		return "", "", i18n.NewError(ctx, code.ScriptParseFailed)
	}
	// 移除updateurl和downloadurl
	reg2 := regexp.MustCompile(`(?im)^//\s*@(updateurl|downloadurl)($|\s+(.+?)$)\s+`)
	ret = reg2.ReplaceAllString(ret, "")
	scriptCode = reg.ReplaceAllLiteralString(scriptCode, ret)
	return scriptCode, ret, nil
}

func parseMetaToJson(meta string) map[string][]string {
	reg := regexp.MustCompile(`(?im)^//\s*@(.+?)($|\s+(.+?)$)`)
	list := reg.FindAllStringSubmatch(meta, -1)
	ret := make(map[string][]string)
	for _, v := range list {
		v[1] = strings.ToLower(v[1])
		if _, ok := ret[v[1]]; !ok {
			ret[v[1]] = make([]string, 0)
		}
		ret[v[1]] = append(ret[v[1]], strings.TrimSpace(v[3]))
	}
	return ret
}

type LibDefinition struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);not null;index:user_id"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	CodeID     int64  `gorm:"column:code_id;type:bigint(20);not null;index:code_id"`
	Definition string `gorm:"column:definition;type:longtext;not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
}
