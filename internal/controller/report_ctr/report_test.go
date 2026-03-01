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

func TestReport_CreateReport_InvalidReason(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	broker.SetBroker(event_bus.NewEvBusBroker())
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("非法举报原因应返回错误", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{UID: 1}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)

		err := testMux.Do(ctx, &api.CreateReportRequest{
			ScriptID: 1, Reason: "invalid_reason", Content: "test content",
		}, &api.CreateReportResponse{})
		convey.So(err, convey.ShouldBeError)
	})
}

func TestReport_CreateReport_Success(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	broker.SetBroker(event_bus.NewEvBusBroker())
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("成功创建举报", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{UID: 1}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		resp := &api.CreateReportResponse{}
		err := testMux.Do(ctx, &api.CreateReportRequest{
			ScriptID: 1, Reason: "malware", Content: "this script contains malware",
		}, resp)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestReport_List(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockUserRepo := mock_user_repo.NewMockUserRepo(mockCtrl)
	user_repo.RegisterUser(mockUserRepo)
	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	mockCommentRepo := mock_report_repo.NewMockScriptReportCommentRepo(mockCtrl)
	report_repo.RegisterScriptReportComment(mockCommentRepo)
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("获取举报列表", t, func() {
		mockAuth.U().NoLogin()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().FindPage(gomock.Any(), int64(1), int32(0), gomock.Any()).
			Return([]*report_entity.ScriptReport{
				{ID: 1, ScriptID: 1, UserID: 1, Reason: "malware", Status: consts.ACTIVE},
			}, int64(1), nil)
		mockUserRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&user_entity.User{
			ID: 1, Username: "testuser",
		}, nil)
		mockCommentRepo.EXPECT().CountByReport(gomock.Any(), int64(1)).Return(int64(3), nil)

		resp := &api.ListResponse{}
		err := testMux.Do(ctx, &api.ListRequest{ScriptID: 1}, resp)
		convey.So(err, convey.ShouldBeNil)
		convey.So(resp.Total, convey.ShouldEqual, 1)
		convey.So(resp.List[0].Reason, convey.ShouldEqual, "malware")
		convey.So(resp.List[0].CommentCount, convey.ShouldEqual, 3)
	})
}

func TestReport_GetReport_NotFound(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("举报不存在", t, func() {
		mockAuth.U().NoLogin()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(99)).Return(nil, nil)
		err := testMux.Do(ctx, &api.GetReportRequest{ScriptID: 1, ReportID: 99}, &api.GetReportResponse{})
		convey.So(err, convey.ShouldBeError)
	})
}

func TestReport_GetReport_Found(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockUserRepo := mock_user_repo.NewMockUserRepo(mockCtrl)
	user_repo.RegisterUser(mockUserRepo)
	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	mockCommentRepo := mock_report_repo.NewMockScriptReportCommentRepo(mockCtrl)
	report_repo.RegisterScriptReportComment(mockCommentRepo)
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("举报存在", t, func() {
		mockAuth.U().NoLogin()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&report_entity.ScriptReport{
			ID: 1, ScriptID: 1, UserID: 1, Reason: "spam", Content: "spam content", Status: consts.ACTIVE,
		}, nil)
		mockUserRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&user_entity.User{
			ID: 1, Username: "reporter",
		}, nil)
		mockCommentRepo.EXPECT().CountByReport(gomock.Any(), int64(1)).Return(int64(0), nil)

		resp := &api.GetReportResponse{}
		err := testMux.Do(ctx, &api.GetReportRequest{ScriptID: 1, ReportID: 1}, resp)
		convey.So(err, convey.ShouldBeNil)
		convey.So(resp.Content, convey.ShouldEqual, "spam content")
		convey.So(resp.Reason, convey.ShouldEqual, "spam")
	})
}

func TestReport_Resolve_NonAdmin(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockAccessRepo := mock_script_repo.NewMockScriptAccessRepo(mockCtrl)
	script_repo.RegisterScriptAccess(mockAccessRepo)
	mockGroupMember := mock_script_repo.NewMockScriptGroupMemberRepo(mockCtrl)
	script_repo.RegisterScriptGroupMember(mockGroupMember)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	broker.SetBroker(event_bus.NewEvBusBroker())
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("非管理员无权解决举报", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{UID: 1}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&report_entity.ScriptReport{
			ID: 1, ScriptID: 1, UserID: 3, Reason: "malware", Status: consts.ACTIVE,
		}, nil)
		mockAccessRepo.EXPECT().FindByLinkID(gomock.Any(), int64(1), int64(1), script_entity.AccessTypeUser).Return(nil, nil)
		mockGroupMember.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)

		err := testMux.Do(ctx, &api.ResolveRequest{ScriptID: 1, ReportID: 1, Close: true}, &api.ResolveResponse{})
		convey.So(err, convey.ShouldBeError)
	})
}

func TestReport_Resolve_Admin(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

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

	convey.Convey("管理员解决举报", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{
			UID: 1, AdminLevel: model.Admin,
		}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&report_entity.ScriptReport{
			ID: 1, ScriptID: 1, UserID: 3, Reason: "malware", Status: consts.ACTIVE,
		}, nil)
		mockReportRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockCommentRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockUserRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&user_entity.User{
			ID: 1, Username: "admin",
		}, nil)

		resp := &api.ResolveResponse{}
		err := testMux.Do(ctx, &api.ResolveRequest{ScriptID: 1, ReportID: 1, Close: true}, resp)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(resp.Comments), convey.ShouldBeGreaterThan, 0)
	})
}

func TestReport_Delete_NonAdmin(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockAccessRepo := mock_script_repo.NewMockScriptAccessRepo(mockCtrl)
	script_repo.RegisterScriptAccess(mockAccessRepo)
	mockGroupMember := mock_script_repo.NewMockScriptGroupMemberRepo(mockCtrl)
	script_repo.RegisterScriptGroupMember(mockGroupMember)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("非管理员无权删除举报", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{UID: 1}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&report_entity.ScriptReport{
			ID: 1, ScriptID: 1, UserID: 3, Status: consts.ACTIVE,
		}, nil)
		mockAccessRepo.EXPECT().FindByLinkID(gomock.Any(), int64(1), int64(1), script_entity.AccessTypeUser).Return(nil, nil)
		mockGroupMember.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)

		err := testMux.Do(ctx, &api.DeleteRequest{ScriptID: 1, ReportID: 1}, &api.DeleteResponse{})
		convey.So(err, convey.ShouldBeError)
	})
}

func TestReport_Delete_Admin(t *testing.T) {
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewReport()
	ctr.limit = limit.NewEmpty()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router)

	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	mockReportRepo := mock_report_repo.NewMockScriptReportRepo(mockCtrl)
	report_repo.RegisterScriptReport(mockReportRepo)
	ctx := context.Background()

	activeScript := &script_entity.Script{
		ID: 1, UserID: 2, Public: script_entity.PublicScript, Status: consts.ACTIVE,
	}

	convey.Convey("管理员删除举报", t, func() {
		mockAuth.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{
			UID: 1, AdminLevel: model.Admin,
		}).AnyTimes()
		mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(activeScript, nil)
		mockReportRepo.EXPECT().Find(gomock.Any(), int64(1), int64(1)).Return(&report_entity.ScriptReport{
			ID: 1, ScriptID: 1, UserID: 3, Status: consts.ACTIVE,
		}, nil)
		mockReportRepo.EXPECT().Delete(gomock.Any(), int64(1), int64(1)).Return(nil)

		err := testMux.Do(ctx, &api.DeleteRequest{ScriptID: 1, ReportID: 1}, &api.DeleteResponse{})
		convey.So(err, convey.ShouldBeNil)
	})
}
