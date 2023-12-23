package issue_ctr

import (
	"context"
	"testing"
	"time"

	"github.com/codfrm/cago/pkg/broker"
	"github.com/codfrm/cago/pkg/broker/event_bus"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/limit"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	mock_issue_repo "github.com/scriptscat/scriptlist/internal/repository/issue_repo/mock"

	"github.com/codfrm/cago/server/mux/muxtest"
	"github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	mock_script_repo "github.com/scriptscat/scriptlist/internal/repository/script_repo/mock"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	mock_user_repo "github.com/scriptscat/scriptlist/internal/repository/user_repo/mock"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	mock_auth_svc "github.com/scriptscat/scriptlist/internal/service/auth_svc/mock"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestIssue_Router(t *testing.T) {
	// 初始化路由
	testMux := muxtest.NewTestMux()
	ctr := NewIssue()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)
	// 注册相关repo
	mockUserRepo := mock_user_repo.NewMockUserRepo(mockCtrl)
	user_repo.RegisterUser(mockUserRepo)
	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	ctx := context.Background()
	mockIssue := mock_issue_repo.NewMockScriptIssueRepo(mockCtrl)
	issue_repo.RegisterScriptIssue(mockIssue)
	mockAccessRepo := mock_script_repo.NewMockScriptAccessRepo(mockCtrl)
	script_repo.RegisterScriptAccess(mockAccessRepo)
	mockGroupMember := mock_script_repo.NewMockScriptGroupMemberRepo(mockCtrl)
	script_repo.RegisterScriptGroupMember(mockGroupMember)
	mockComment := mock_issue_repo.NewMockScriptIssueCommentRepo(mockCtrl)
	issue_repo.RegisterScriptIssueComment(mockComment)
	broker.SetBroker(event_bus.NewEvBusBroker())
	convey.Convey("删除反馈", t, func() {
		mockAuth.U().Get()
		scriptCall := mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.Script{
			ID:     1,
			UserID: 2,
			Public: script_entity.PublicScript,
			Status: consts.ACTIVE,
		}, nil)
		convey.Convey("反馈不存在", func() {
			mockIssue.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
			err := testMux.Do(ctx, &issue.DeleteRequest{ScriptID: 1, IssueID: 1}, &issue.DeleteResponse{})
			convey.So(err, convey.ShouldBeError, "反馈不存在")
		})
		convey.Convey("存在", func() {
			issueCall := mockIssue.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&issue_entity.ScriptIssue{
				ID:       1,
				ScriptID: 1,
				UserID:   3,
				Status:   consts.ACTIVE,
			}, nil)
			convey.Convey("脚本归档", func() {
				scriptCall.Return(&script_entity.Script{
					ID:      1,
					UserID:  2,
					Public:  script_entity.PublicScript,
					Archive: script_entity.IsArchive,
					Status:  consts.ACTIVE,
				}, nil)
				err := testMux.Do(ctx, &issue.DeleteRequest{ScriptID: 1, IssueID: 1}, &issue.DeleteResponse{})
				convey.So(err, convey.ShouldBeError, "脚本已归档,无法进行此操作")
			})
			convey.Convey("自己删除", func() {
				mockAuth.U().Get()
				issueCall.Return(&issue_entity.ScriptIssue{
					ID:       1,
					ScriptID: 1,
					UserID:   1,
					Status:   consts.ACTIVE,
				}, nil)
				mockAccessRepo.EXPECT().FindByLinkID(gomock.Any(), int64(1), int64(1), script_entity.AccessTypeUser).Return(nil, nil)
				mockGroupMember.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
				err := testMux.Do(ctx, &issue.DeleteRequest{ScriptID: 1, IssueID: 1}, &issue.DeleteResponse{})
				convey.So(err, convey.ShouldBeError, "用户不允许操作")
			})
			convey.Convey("脚本所有者删除", func() {
				mockAuth.U().Get(2)
				mockIssue.EXPECT().Delete(gomock.Any(), int64(1), int64(1)).Return(nil)
				err := testMux.Do(ctx, &issue.DeleteRequest{ScriptID: 1, IssueID: 1}, &issue.DeleteResponse{})
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
	convey.Convey("打开评论", t, func() {
		mockAuth.U().Get()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.Script{
			ID:     1,
			UserID: 2,
			Public: script_entity.PublicScript,
			Status: consts.ACTIVE,
		}, nil)
		mockIssue.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&issue_entity.ScriptIssue{
			ID:       1,
			ScriptID: 1,
			UserID:   3,
			Status:   consts.ACTIVE,
		}, nil)
		convey.Convey("自己打开", func() {
			mockAuth.U().Get(3).Times(3)
			mockIssue.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			mockComment.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			mockUserRepo.EXPECT().Find(gomock.Any(), int64(3)).Return(&user_entity.User{
				ID: 3,
			}, nil)
			err := testMux.Do(ctx, &issue.OpenRequest{ScriptID: 1, IssueID: 1}, &issue.OpenResponse{})
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("非自己打开", func() {
			// 无权限
			mockAuth.U().Get().Times(2)
			mockGroupMember.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
			convey.Convey("无权限", func() {
				mockAccessRepo.EXPECT().FindByLinkID(gomock.Any(), int64(1), int64(1), script_entity.AccessTypeUser).Return(nil, nil)
				err := testMux.Do(ctx, &issue.OpenRequest{ScriptID: 1, IssueID: 1}, &issue.OpenResponse{})
				convey.So(err, convey.ShouldBeError, "用户不允许操作")
			})
			convey.Convey("管理员权限", func() {
				mockAuth.U().Get().Times(2)
				mockAccessRepo.EXPECT().FindByLinkID(gomock.Any(), int64(1), int64(1), script_entity.AccessTypeUser).
					Return([]*script_entity.ScriptAccess{{
						ID:         1,
						ScriptID:   1,
						LinkID:     1,
						Type:       1,
						Role:       script_entity.AccessRoleManager,
						Status:     consts.ACTIVE,
						Expiretime: time.Now().Add(time.Hour).Unix(),
					}}, nil)
				mockIssue.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				mockComment.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mockUserRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&user_entity.User{
					ID: 1,
				}, nil)
				err := testMux.Do(ctx, &issue.OpenRequest{ScriptID: 1, IssueID: 1}, &issue.OpenResponse{})
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}
