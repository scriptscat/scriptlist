package report_entity

import (
	"context"
	"testing"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/stretchr/testify/assert"
)

func TestScriptReport_CheckOperate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		report  *ScriptReport
		wantErr bool
	}{
		{
			name:    "nil报告应返回错误",
			report:  nil,
			wantErr: true,
		},
		{
			name:    "已删除报告应返回错误",
			report:  &ScriptReport{Status: consts.DELETE},
			wantErr: true,
		},
		{
			name:    "活跃报告应返回nil",
			report:  &ScriptReport{Status: consts.ACTIVE},
			wantErr: false,
		},
		{
			name:    "审核状态报告应返回nil",
			report:  &ScriptReport{Status: consts.AUDIT},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.report.CheckOperate(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScriptReport_IsResolved(t *testing.T) {
	tests := []struct {
		name   string
		report *ScriptReport
		want   bool
	}{
		{
			name:   "AUDIT状态应返回true",
			report: &ScriptReport{Status: consts.AUDIT},
			want:   true,
		},
		{
			name:   "ACTIVE状态应返回false",
			report: &ScriptReport{Status: consts.ACTIVE},
			want:   false,
		},
		{
			name:   "DELETE状态应返回false",
			report: &ScriptReport{Status: consts.DELETE},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.report.IsResolved())
		})
	}
}

func TestScriptReport_Resolve(t *testing.T) {
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	report := &ScriptReport{Status: consts.ACTIVE}

	report.Resolve(now)

	assert.Equal(t, int32(consts.AUDIT), report.Status)
	assert.Equal(t, now.Unix(), report.Updatetime)
}

func TestScriptReport_Reopen(t *testing.T) {
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	report := &ScriptReport{Status: consts.AUDIT}

	report.Reopen(now)

	assert.Equal(t, int32(consts.ACTIVE), report.Status)
	assert.Equal(t, now.Unix(), report.Updatetime)
}

func TestScriptReport_ValidateReason(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		reason  string
		wantErr bool
	}{
		{name: "malware是合法原因", reason: "malware", wantErr: false},
		{name: "privacy是合法原因", reason: "privacy", wantErr: false},
		{name: "copyright是合法原因", reason: "copyright", wantErr: false},
		{name: "spam是合法原因", reason: "spam", wantErr: false},
		{name: "other是合法原因", reason: "other", wantErr: false},
		{name: "空字符串是非法原因", reason: "", wantErr: true},
		{name: "unknown是非法原因", reason: "unknown", wantErr: true},
		{name: "arbitrary是非法原因", reason: "arbitrary", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &ScriptReport{Reason: tt.reason}
			err := report.ValidateReason(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScriptReportComment_CheckOperate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		comment *ScriptReportComment
		wantErr bool
	}{
		{
			name:    "nil评论应返回错误",
			comment: nil,
			wantErr: true,
		},
		{
			name:    "活跃评论应返回nil",
			comment: &ScriptReportComment{Status: consts.ACTIVE},
			wantErr: false,
		},
		{
			name:    "已删除评论应返回错误",
			comment: &ScriptReportComment{Status: consts.DELETE},
			wantErr: true,
		},
		{
			name:    "审核状态评论应返回错误",
			comment: &ScriptReportComment{Status: consts.AUDIT},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.CheckOperate(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
