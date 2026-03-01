package subscribe

import (
	"context"
	"testing"

	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/audit_repo"
	mock_audit_repo "github.com/scriptscat/scriptlist/internal/repository/audit_repo/mock"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuditLog_scriptDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_audit_repo.NewMockAuditLogRepo(mockCtrl)
	audit_repo.RegisterAuditLog(mockRepo)

	a := &AuditLog{}
	ctx := context.Background()

	t.Run("管理员删除脚本-带理由", func(t *testing.T) {
		msg := &producer.ScriptDeleteMsg{
			Script: &script_entity.Script{
				ID:   100,
				Name: "test script",
			},
			Operator: producer.Operator{
				OperatorUID:      1,
				OperatorUsername: "admin",
				IsAdmin:          true,
			},
			Reason: "违反使用规范",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, log *audit_entity.AuditLog) error {
				assert.Equal(t, int64(1), log.UserID)
				assert.Equal(t, "admin", log.Username)
				assert.Equal(t, audit_entity.ActionScriptDelete, log.Action)
				assert.Equal(t, "script", log.TargetType)
				assert.Equal(t, int64(100), log.TargetID)
				assert.Equal(t, "test script", log.TargetName)
				assert.True(t, log.IsAdmin)
				assert.Equal(t, "违反使用规范", log.Reason)
				assert.Greater(t, log.Createtime, int64(0))
				return nil
			},
		)

		err := a.scriptDelete(ctx, msg)
		assert.NoError(t, err)
	})

	t.Run("普通用户删除脚本-无理由", func(t *testing.T) {
		msg := &producer.ScriptDeleteMsg{
			Script: &script_entity.Script{
				ID:   200,
				Name: "my script",
			},
			Operator: producer.Operator{
				OperatorUID:      5,
				OperatorUsername: "user5",
				IsAdmin:          false,
			},
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, log *audit_entity.AuditLog) error {
				assert.Equal(t, int64(5), log.UserID)
				assert.Equal(t, "user5", log.Username)
				assert.False(t, log.IsAdmin)
				assert.Equal(t, "", log.Reason)
				return nil
			},
		)

		err := a.scriptDelete(ctx, msg)
		assert.NoError(t, err)
	})
}

func TestAuditLog_scriptCodeUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_audit_repo.NewMockAuditLogRepo(mockCtrl)
	audit_repo.RegisterAuditLog(mockRepo)

	a := &AuditLog{}
	ctx := context.Background()

	t.Run("管理员更新脚本代码", func(t *testing.T) {
		msg := &producer.ScriptCodeUpdateMsg{
			Script: &script_entity.Script{
				ID:   100,
				Name: "test script",
			},
			CodeID: 50,
			Operator: producer.Operator{
				OperatorUID:      2,
				OperatorUsername: "admin2",
				IsAdmin:          true,
			},
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, log *audit_entity.AuditLog) error {
				assert.Equal(t, int64(2), log.UserID)
				assert.Equal(t, "admin2", log.Username)
				assert.Equal(t, audit_entity.ActionScriptUpdate, log.Action)
				assert.Equal(t, "script", log.TargetType)
				assert.Equal(t, int64(100), log.TargetID)
				assert.True(t, log.IsAdmin)
				return nil
			},
		)

		err := a.scriptCodeUpdate(ctx, msg)
		assert.NoError(t, err)
	})
}

func TestAuditLog_scriptCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_audit_repo.NewMockAuditLogRepo(mockCtrl)
	audit_repo.RegisterAuditLog(mockRepo)

	a := &AuditLog{}
	ctx := context.Background()

	t.Run("创建脚本", func(t *testing.T) {
		msg := &producer.ScriptCreateMsg{
			Script: &script_entity.Script{
				ID:   300,
				Name: "new script",
			},
			CodeID: 60,
			Operator: producer.Operator{
				OperatorUID:      3,
				OperatorUsername: "user3",
				IsAdmin:          false,
			},
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, log *audit_entity.AuditLog) error {
				assert.Equal(t, int64(3), log.UserID)
				assert.Equal(t, "user3", log.Username)
				assert.Equal(t, audit_entity.ActionScriptCreate, log.Action)
				assert.Equal(t, "script", log.TargetType)
				assert.Equal(t, int64(300), log.TargetID)
				assert.Equal(t, "new script", log.TargetName)
				assert.False(t, log.IsAdmin)
				return nil
			},
		)

		err := a.scriptCreate(ctx, msg)
		assert.NoError(t, err)
	})
}
