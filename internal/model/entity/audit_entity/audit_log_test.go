package audit_entity

import (
	"testing"
)

func TestAuditLog_IsAdminAction(t *testing.T) {
	tests := []struct {
		name    string
		isAdmin bool
		want    bool
	}{
		{"管理员操作", true, true},
		{"普通用户操作", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuditLog{IsAdmin: tt.isAdmin}
			if got := a.IsAdminAction(); got != tt.want {
				t.Errorf("IsAdminAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuditLog_IsDeleteAction(t *testing.T) {
	tests := []struct {
		name   string
		action Action
		want   bool
	}{
		{"删除操作", ActionScriptDelete, true},
		{"更新操作", ActionScriptUpdate, false},
		{"创建操作", ActionScriptCreate, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuditLog{Action: tt.action}
			if got := a.IsDeleteAction(); got != tt.want {
				t.Errorf("IsDeleteAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
