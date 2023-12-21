package script_ctr

import (
	"context"
	"testing"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/server/mux/muxtest"
	"github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	mock_script_repo "github.com/scriptscat/scriptlist/internal/repository/script_repo/mock"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	mock_user_repo "github.com/scriptscat/scriptlist/internal/repository/user_repo/mock"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	mock_auth_svc "github.com/scriptscat/scriptlist/internal/service/auth_svc/mock"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestScript_Router(t *testing.T) {
	// 初始化路由
	testMux := muxtest.NewTestMux()
	ctr := NewScript()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuth := mock_auth_svc.InitAuth(mockCtrl)
	auth_svc.RegisterAuth(mockAuth)
	ctr.Router(testMux.Router, testMux.Router)
	// 注册相关repo
	mockUserRepo := mock_user_repo.NewMockUserRepo(mockCtrl)
	user_repo.RegisterUser(mockUserRepo)
	mockScriptRepo := mock_script_repo.NewMockScriptRepo(mockCtrl)
	script_repo.RegisterScript(mockScriptRepo)
	ctx := context.Background()
	convey.Convey("测试路由访问权限", t, func() {
		convey.Convey("是否登录", func() {
			convey.Convey("未登录", func() {
				convey.Convey("无需登录", func() {
					mockScriptRepo.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), nil).Times(1)
					err := testMux.Do(ctx, &script.ListRequest{}, &script.ListResponse{})
					assert.NoError(t, err)
				})
				convey.Convey("强制登录", func() {
					mockAuth.U().NoLogin()
					err := testMux.Do(ctx, &script.CreateRequest{}, &script.CreateResponse{})
					assert.EqualError(t, err, "未登录")
				})
			})
			convey.Convey("已登录", func() {
				mockAuth.U().Get()
				err := testMux.Do(ctx, &script.CreateRequest{}, &script.CreateResponse{})
				assert.EqualError(t, err, "脚本详细描述为必填字段")
			})
		})
		convey.Convey("脚本状态", func() {
			convey.Convey("脚本不存在", func() {
				mockAuth.U().NoLogin()
				mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(nil, nil)
				err := testMux.Do(ctx, &script.CodeRequest{ID: 1}, &script.CodeRequest{})
				convey.So(err, convey.ShouldBeError, "脚本不存在")
			})
			convey.Convey("脚本被删除", func() {
				mockAuth.U().NoLogin()
				mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.Script{
					ID:     1,
					Status: consts.DELETE,
				}, nil)
				err := testMux.Do(ctx, &script.CodeRequest{ID: 1}, &script.CodeRequest{})
				convey.So(err, convey.ShouldBeError, "脚本被删除")
			})
		})
		mockAccess := mock_script_repo.NewMockScriptAccessRepo(mockCtrl)
		script_repo.RegisterScriptAccess(mockAccess)
		mockMember := mock_script_repo.NewMockScriptGroupMemberRepo(mockCtrl)
		script_repo.RegisterScriptGroupMember(mockMember)
		mockCodeRepo := mock_script_repo.NewMockScriptCodeRepo(mockCtrl)
		script_repo.RegisterScriptCode(mockCodeRepo)
		convey.Convey("私有脚本", func() {
			scriptCall := mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.Script{
				ID:     1,
				UserID: 2,
				Public: script_entity.PrivateScript,
				Status: consts.ACTIVE,
			}, nil)
			convey.Convey("未登录", func() {
				mockAuth.U().NoLogin().Times(2)
				err := testMux.Do(ctx, &script.CodeRequest{ID: 1}, &script.CodeRequest{})
				convey.So(err, convey.ShouldBeError, "用户未登录")
			})
			convey.Convey("登录无权限", func() {
				mockAuth.U().Get().Times(3)
				mockAccess.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
				mockMember.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
				err := testMux.Do(ctx, &script.CodeRequest{ID: 1}, &script.CodeRequest{})
				convey.So(err, convey.ShouldBeError, "用户不允许操作")
			})
			convey.Convey("拥有权限", func() {
				call := mockAuth.U().Get().Times(3)
				convey.Convey("权限到期", func() {
					mockAccess.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).
						Return(&script_entity.ScriptAccess{
							ID:         1,
							ScriptID:   1,
							LinkID:     1,
							Type:       1,
							Role:       "guest",
							Expiretime: time.Now().Add(-time.Hour).Unix(),
						}, nil)
					convey.Convey("访客可以访问", func() {
						err := testMux.Do(ctx, &script.VersionListRequest{ID: 1}, &script.VersionListRequest{})
						convey.So(err, convey.ShouldBeError, "用户不允许操作")
					})
				})
				convey.Convey("访客权限", func() {
					mockAccess.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).
						Return(&script_entity.ScriptAccess{
							ID:         1,
							ScriptID:   1,
							LinkID:     1,
							Type:       1,
							Role:       "guest",
							Expiretime: time.Now().Add(time.Hour).Unix(),
						}, nil)
					convey.Convey("访客可以访问", func() {
						mockCodeRepo.EXPECT().List(gomock.Any(), int64(1), gomock.Any()).Return(nil, int64(0), nil)
						err := testMux.Do(ctx, &script.VersionListRequest{ID: 1}, &script.VersionListRequest{})
						convey.So(err, convey.ShouldBeNil)
					})
					convey.Convey("访客不可以访问", func() {
						err := testMux.Do(ctx, &script.GetSettingRequest{ID: 1}, &script.GetSettingResponse{})
						convey.So(err, convey.ShouldBeError, "用户不允许操作")
					})
				})
				convey.Convey("脚本管理员", func() {
					mockAccess.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).
						Return(&script_entity.ScriptAccess{
							ID:         1,
							ScriptID:   1,
							LinkID:     1,
							Type:       1,
							Role:       "manager",
							Expiretime: time.Now().Add(time.Hour).Unix(),
						}, nil)
					convey.Convey("管理员不允许归档与删除", func() {
						err := testMux.Do(ctx, &script.ArchiveRequest{ID: 1, Archive: true}, &script.ArchiveResponse{})
						convey.So(err, convey.ShouldBeError, "用户不允许操作")
					})
				})
				convey.Convey("脚本拥有者", func() {
					call.Return(&model.AuthInfo{
						UID: 2,
					})
					convey.Convey("脚本拥有者可以删除和归档", func() {
						mockScriptRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
						err := testMux.Do(ctx, &script.ArchiveRequest{ID: 1, Archive: true}, &script.ArchiveResponse{})
						convey.So(err, convey.ShouldBeNil)
					})
				})
				convey.Convey("版主", func() {
					call.Return(&model.AuthInfo{
						UID:        1,
						AdminLevel: model.Moderator,
					})
					mockAccess.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
					mockMember.EXPECT().FindByUserId(gomock.Any(), int64(1), int64(1)).Return(nil, nil)
					convey.Convey("版主不可以删除和归档", func() {
						err := testMux.Do(ctx, &script.ArchiveRequest{ID: 1, Archive: true}, &script.ArchiveResponse{})
						convey.So(err, convey.ShouldBeError, "用户不允许操作")
					})
				})
				convey.Convey("超级管理员", func() {
					call.Return(&model.AuthInfo{
						UID:        1,
						AdminLevel: model.Admin,
					})
					convey.Convey("超级管理员可以删除和归档", func() {
						mockScriptRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
						err := testMux.Do(ctx, &script.ArchiveRequest{ID: 1, Archive: true}, &script.ArchiveResponse{})
						convey.So(err, convey.ShouldBeNil)
					})
					convey.Convey("归档不可操作", func() {
						scriptCall.Return(&script_entity.Script{
							ID:      1,
							UserID:  1,
							Archive: script_entity.IsArchive,
							Public:  script_entity.PrivateScript,
							Status:  consts.ACTIVE,
						}, nil)
						err := testMux.Do(ctx, &script.UpdateSettingRequest{ID: 1}, &script.UpdateSettingResponse{})
						convey.So(err, convey.ShouldBeError, "脚本已归档,无法进行此操作")
					})
				})
			})
		})
	})
}
