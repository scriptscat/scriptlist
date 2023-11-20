package script_entity

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type AccessPermission struct {
	allow map[string]bool
}

func (a *AccessPermission) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	s := strings.Split(string(bytes), ",")
	a.allow = make(map[string]bool)
	for _, v := range s {
		a.allow[v] = true
	}
	return nil
}

func (a *AccessPermission) Value() (driver.Value, error) {
	s := make([]string, 0)
	for k, v := range a.allow {
		if v {
			s = append(s, k)
		}
	}
	return strings.Join(s, ","), nil
}

// Read 读取权限
func (a *AccessPermission) Read() bool {
	return a.allow["read"]
}

// Write 写入权限
func (a *AccessPermission) Write() bool {
	return a.allow["write"]
}

// Partner 合伙人权限
func (a *AccessPermission) Partner() bool {
	return a.allow["partner"]
}

type ScriptAccess struct {
	ID               int64             `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID         int64             `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	LinkID           int64             `gorm:"column:link_id;type:bigint(20);not null"`
	Type             int32             `gorm:"column:type;type:tinyint(4);not null"`
	AccessPermission *AccessPermission `gorm:"column:access_permission;type:varchar(255);not null"`
	Expiretime       int64             `gorm:"column:expiretime;type:bigint(20)"`
	Createtime       int64             `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime       int64             `gorm:"column:updatetime;type:bigint(20)"`
}
