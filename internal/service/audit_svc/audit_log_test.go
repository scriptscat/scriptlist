package audit_svc

import (
	"context"
	"testing"

	api "github.com/scriptscat/scriptlist/internal/api/audit"
	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
	"github.com/scriptscat/scriptlist/internal/repository/audit_repo"
	mock_audit_repo "github.com/scriptscat/scriptlist/internal/repository/audit_repo/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuditLogSvc_List(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_audit_repo.NewMockAuditLogRepo(mockCtrl)
	audit_repo.RegisterAuditLog(mockRepo)

	ctx := context.Background()
	svc := &auditLogSvc{}

	t.Run("仅返回管理员删除的脚本", func(t *testing.T) {
		expectedLogs := []*audit_entity.AuditLog{
			{
				ID:         1,
				UserID:     10,
				Username:   "admin",
				Action:     audit_entity.ActionScriptDelete,
				TargetType: "script",
				TargetID:   100,
				TargetName: "test script",
				IsAdmin:    true,
				Reason:     "违规内容",
				Createtime: 1709280000,
			},
		}

		mockRepo.EXPECT().FindPage(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, opts *audit_repo.ListOptions) ([]*audit_entity.AuditLog, int64, error) {
				// 验证强制 is_admin=true 且 action=script_delete
				assert.NotNil(t, opts.IsAdmin)
				assert.True(t, *opts.IsAdmin)
				assert.Equal(t, "script_delete", opts.Action)
				return expectedLogs, int64(1), nil
			},
		)

		resp, err := svc.List(ctx, &api.ListRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(1), resp.Total)
		assert.Len(t, resp.List, 1)
		assert.Equal(t, int64(1), resp.List[0].ID)
		assert.Equal(t, "admin", resp.List[0].Username)
		assert.Equal(t, audit_entity.ActionScriptDelete, resp.List[0].Action)
		assert.True(t, resp.List[0].IsAdmin)
		assert.Equal(t, "违规内容", resp.List[0].Reason)
	})

	t.Run("空结果", func(t *testing.T) {
		mockRepo.EXPECT().FindPage(gomock.Any(), gomock.Any()).Return(
			[]*audit_entity.AuditLog{}, int64(0), nil,
		)

		resp, err := svc.List(ctx, &api.ListRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.List)
		assert.Equal(t, int64(0), resp.Total)
	})
}

func TestAuditLogSvc_ScriptList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_audit_repo.NewMockAuditLogRepo(mockCtrl)
	audit_repo.RegisterAuditLog(mockRepo)

	ctx := context.Background()
	svc := &auditLogSvc{}

	t.Run("返回单脚本所有日志", func(t *testing.T) {
		expectedLogs := []*audit_entity.AuditLog{
			{
				ID:         1,
				UserID:     10,
				Username:   "admin",
				Action:     audit_entity.ActionScriptDelete,
				TargetType: "script",
				TargetID:   100,
				IsAdmin:    true,
				Createtime: 1709280000,
			},
			{
				ID:         2,
				UserID:     20,
				Username:   "user1",
				Action:     audit_entity.ActionScriptUpdate,
				TargetType: "script",
				TargetID:   100,
				IsAdmin:    false,
				Createtime: 1709270000,
			},
		}

		mockRepo.EXPECT().FindPage(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, opts *audit_repo.ListOptions) ([]*audit_entity.AuditLog, int64, error) {
				// 验证按 target 筛选
				assert.Equal(t, "script", opts.TargetType)
				assert.Equal(t, int64(100), opts.TargetID)
				// 不强制 is_admin
				assert.Nil(t, opts.IsAdmin)
				return expectedLogs, int64(2), nil
			},
		)

		resp, err := svc.ScriptList(ctx, &api.ScriptListRequest{ID: 100})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(2), resp.Total)
		assert.Len(t, resp.List, 2)
		// 验证包含管理员和普通用户操作
		assert.True(t, resp.List[0].IsAdmin)
		assert.False(t, resp.List[1].IsAdmin)
	})
}

func TestToAPIItems(t *testing.T) {
	logs := []*audit_entity.AuditLog{
		{
			ID:         1,
			UserID:     10,
			Username:   "testuser",
			Action:     audit_entity.ActionScriptCreate,
			TargetType: "script",
			TargetID:   50,
			TargetName: "my script",
			IsAdmin:    false,
			Reason:     "",
			Createtime: 1709280000,
		},
	}

	items := toAPIItems(logs)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(1), items[0].ID)
	assert.Equal(t, int64(10), items[0].UserID)
	assert.Equal(t, "testuser", items[0].Username)
	assert.Equal(t, audit_entity.ActionScriptCreate, items[0].Action)
	assert.Equal(t, "script", items[0].TargetType)
	assert.Equal(t, int64(50), items[0].TargetID)
	assert.Equal(t, "my script", items[0].TargetName)
	assert.False(t, items[0].IsAdmin)
	assert.Equal(t, "", items[0].Reason)
	assert.Equal(t, int64(1709280000), items[0].Createtime)
}

func TestToAPIItems_EmptyList(t *testing.T) {
	items := toAPIItems([]*audit_entity.AuditLog{})
	assert.Empty(t, items)
	assert.NotNil(t, items)
}
