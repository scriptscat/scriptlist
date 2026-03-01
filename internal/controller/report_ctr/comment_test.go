package report_ctr

import (
	"context"
	"testing"

	"github.com/cago-frame/cago/pkg/broker"
	"github.com/cago-frame/cago/pkg/broker/event_bus"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/limit"
	"github.com/cago-frame/cago/server/mux/muxtest"
	api "github.com/scriptscat/scriptlist/internal/api/report"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/report_repo"
	mock_report_repo "github.com/scriptscat/scriptlist/internal/repository/report_repo/mock"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	mock_script_repo "github.com/scriptscat/scriptlist/internal/repository/script_repo/mock"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	mock_user_repo "github.com/scriptscat/scriptlist/internal/repository/user_repo/mock"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	mock_auth_svc "github.com/scriptscat/scriptlist/internal/service/auth_svc/mock"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestReportComment_List(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	commentCtr := NewReportComment()
	commentCtr.limit = limit.NewPeriodLimit(0, 999, nil, "")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	commentCtr.Router(testMux.Router)

	mockUserRepo := mock_user_repo.NewMockUserRepo(mockCtrl)
	user_repo.RegisterUser(mockUserRepo)
	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	mockCommentRepo := mock_report_repo.NewMockScriptReportCommentRepo(mockCtrl)
	report_repo.RegisterScriptReportComment(mockCommentRepo)
	broker.SetBroker(event_bus.NewEvBusBroker())
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}
	activeReport := &report_entity.ScriptReport{
		ID: 1, ScriptID: 1, UserID: 3, Reason: "malware", Status: consts.ACTIVE,
	}

	convey.Convey("获取评论列表", t, func() {
		mockAuth.U().NoLogin()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(activeReport, nil)
		mockCommentRepo.EXPECT().FindAll(gomock.Any(), int64(1)).Return([]*report_entity.ScriptReportComment{
			{ID: 1, ReportID: 1, UserID: 1, Content: "comment 1", Type: report_entity.CommentTypeComment, Status: consts.ACTIVE},
			{ID: 2, ReportID: 1, UserID: 2, Content: "comment 2", Type: report_entity.CommentTypeComment, Status: consts.ACTIVE},
		}, nil)
		mockUserRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&user_entity.User{
			ID: 1, Username: "user1",
		}, nil)
		mockUserRepo.EXPECT().Find(gomock.Any(), int64(2)).Return(&user_entity.User{
			ID: 2, Username: "user2",
		}, nil)

		resp := &api.ListCommentResponse{}
		err := testMux.Do(ctx, &api.ListCommentRequest{ScriptID: 1, ReportID: 1}, resp)
		convey.So(err, convey.ShouldBeNil)
		convey.So(resp.Total, convey.ShouldEqual, 2)
	})
}

func TestReportComment_Delete_NotFound(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	commentCtr := NewReportComment()
	commentCtr.limit = limit.NewPeriodLimit(0, 999, nil, "")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	commentCtr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	mockCommentRepo := mock_report_repo.NewMockScriptReportCommentRepo(mockCtrl)
	report_repo.RegisterScriptReportComment(mockCommentRepo)
	broker.SetBroker(event_bus.NewEvBusBroker())
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}
	activeReport := &report_entity.ScriptReport{
		ID: 1, ScriptID: 1, UserID: 3, Reason: "malware", Status: consts.ACTIVE,
	}

	convey.Convey("评论不存在", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{
			UID: 1, AdminLevel: model.Admin,
		}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(activeReport, nil)
		mockCommentRepo.EXPECT().Find(gomock.Any(), int64(1), int64(99)).Return(nil, nil)

		err := testMux.Do(ctx, &api.DeleteCommentRequest{ScriptID: 1, ReportID: 1, CommentID: 99}, &api.DeleteCommentResponse{})
		convey.So(err, convey.ShouldBeError)
	})
}

func TestReportComment_Delete_Admin(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	commentCtr := NewReportComment()
	commentCtr.limit = limit.NewPeriodLimit(0, 999, nil, "")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	commentCtr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	mockCommentRepo := mock_report_repo.NewMockScriptReportCommentRepo(mockCtrl)
	report_repo.RegisterScriptReportComment(mockCommentRepo)
	broker.SetBroker(event_bus.NewEvBusBroker())
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}
	activeReport := &report_entity.ScriptReport{
		ID: 1, ScriptID: 1, UserID: 3, Reason: "malware", Status: consts.ACTIVE,
	}

	convey.Convey("管理员删除评论", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{
			UID: 1, AdminLevel: model.Admin,
		}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(activeReport, nil)
		mockCommentRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&report_entity.ScriptReportComment{
			ID: 1, ReportID: 1, UserID: 3, Status: consts.ACTIVE,
		}, nil)
		mockCommentRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

		err := testMux.Do(ctx, &api.DeleteCommentRequest{ScriptID: 1, ReportID: 1, CommentID: 1}, &api.DeleteCommentResponse{})
		convey.So(err, convey.ShouldBeNil)
	})
}
