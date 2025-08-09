package script_entity

import (
	"context"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// ScriptGroup 脚本组
type ScriptGroup struct {
	ID          int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID    int64  `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	Name        string `gorm:"column:name;type:varchar(255);not null"`
	Description string `gorm:"column:description;type:varchar(255);"`
	Status      int32  `gorm:"column:status;type:tinyint(4);not null"`
	Createtime  int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime  int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (g *ScriptGroup) Check(ctx context.Context) error {
	if g == nil {
		return i18n.NewNotFoundError(ctx, code.GroupNotFound)
	}
	if g.Status != consts.ACTIVE {
		return i18n.NewNotFoundError(ctx, code.GroupNotFound)
	}
	return nil
}
