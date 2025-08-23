package script_ctr

import (
	"context"
	"errors"
	"testing"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/cago-frame/cago/server/mux/muxtest"
	"github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	mock_script_repo "github.com/scriptscat/scriptlist/internal/repository/script_repo/mock"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	mock_user_repo "github.com/scriptscat/scriptlist/internal/repository/user_repo/mock"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	mock_auth_svc "github.com/scriptscat/scriptlist/internal/service/auth_svc/mock"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestScore_Router(t *testing.T) {
	// 初始化路由
	testMux := muxtest.NewTestMux(muxtest.WithBaseUrl(""))
	ctr := NewScore()
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
	mockScore := mock_script_repo.NewMockScriptScoreRepo(mockCtrl)
	script_repo.RegisterScriptScore(mockScore)
	mockAccess := mock_script_repo.NewMockScriptAccessRepo(mockCtrl)
	script_repo.RegisterScriptAccess(mockAccess)
	mockGroupMember := mock_script_repo.NewMockScriptGroupMemberRepo(mockCtrl)
	script_repo.RegisterScriptGroupMember(mockGroupMember)
	mockStatistics := mock_script_repo.NewMockScriptStatisticsRepo(mockCtrl)
	script_repo.RegisterScriptStatistics(mockStatistics)
	convey.Convey("删除评分", t, func() {
		convey.Convey("未登录", func() {
			mockAuth.U().NoLogin()
			err := testMux.Do(ctx, &script.DelScoreRequest{}, &script.DelScoreResponse{})
			convey.So(err, convey.ShouldBeError, "未登录")
		})
		convey.Convey("登录", func() {
			call := mockAuth.U().Get().Times(2)
			mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.Script{
				ID:     1,
				UserID: 1,
				Public: script_entity.PublicScript,
				Status: consts.ACTIVE,
			}, nil)
			convey.Convey("作者不能删除评分", func() {
				err := testMux.Do(ctx, &script.DelScoreRequest{ScriptId: 1, ScoreId: 1}, &script.DelScoreResponse{})
				convey.So(err, convey.ShouldBeError, "用户不允许操作")
			})
			convey.Convey("社区管理员可以删除评分", func() {
				call.Return(&model.AuthInfo{
					UID:        1,
					AdminLevel: model.Admin,
				})
				mockScore.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.ScriptScore{
					ID:       1,
					UserID:   1,
					ScriptID: 1,
					Score:    1,
				}, nil)
				mockScore.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil).Times(1)
				mockStatistics.EXPECT().IncrScore(gomock.Any(), int64(1), int64(1), -1).Return(nil)
				err := testMux.Do(ctx, &script.DelScoreRequest{ScriptId: 1, ScoreId: 1}, &script.DelScoreResponse{})
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
	convey.Convey("回复评分", t, func() {
		convey.Convey("未登录", func() {
			mockAuth.U().NoLogin()
			err := testMux.Do(ctx, &script.ReplyScoreRequest{}, &script.ReplyScoreResponse{})
			convey.So(err, convey.ShouldBeError, "未登录")
		})
		convey.Convey("非作者回复", func() {
			mockAuth.U().Get().Times(2) //登录检测查一次，权限查一次
			mockAccess.EXPECT().FindByLinkID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*script_entity.ScriptAccess{}, nil)
			mockScriptRepo.EXPECT().Find(gomock.Any(), int64(1)).Return(&script_entity.Script{
				ID:     1,
				UserID: 2,
				Status: consts.ACTIVE,
			}, nil)
			mockGroupMember.EXPECT().FindByUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*script_entity.ScriptGroupMember{}, nil)
			err := testMux.Do(ctx, &script.ReplyScoreRequest{ScriptId: 1, CommentID: 1, Message: "测试"}, &script.ReplyScoreResponse{})
			convey.So(err, convey.ShouldBeError, "用户不允许操作")
		})
		convey.Convey("作者回复", func() {
			//登录检测查一次，权限查一次
			mockAuth.U().Get().Times(2).Return(&model.AuthInfo{
				UID: 2,
			})
			mockScriptRepo.EXPECT().Find(gomock.Any(), int64(66)).Return(&script_entity.Script{
				ID:     1,
				UserID: 2,
				Status: consts.ACTIVE,
			}, nil).Times(2)
			mockScore.EXPECT().Find(gomock.Any(), int64(21)).Return(&script_entity.ScriptScore{
				ID:     21,
				UserID: 1,
			}, nil)
			convey.Convey("作者创建回复信息", func() {
				mockScore.EXPECT().FindReplayByComment(gomock.Any(), int64(21), int64(66)).Return(nil, nil)
				mockScore.EXPECT().CreateReplayByComment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, scoreReply *script_entity.ScriptScoreReply) error {
					if scoreReply.Message == "测试" && scoreReply.ScriptID == 66 && scoreReply.CommentID == 21 {
						return nil
					}
					return errors.New("数据不匹配")
				})
				mockUserRepo.EXPECT().Find(gomock.Any(), int64(2)).Return(nil, errors.New("发信用户搜索错误"))

				err := testMux.Do(ctx, &script.ReplyScoreRequest{ScriptId: 66, CommentID: 21, Message: "测试"}, &script.ReplyScoreResponse{})
				convey.So(err, convey.ShouldBeNil)
			})
			convey.Convey("作者修改回复信息", func() {
				mockScore.EXPECT().FindReplayByComment(gomock.Any(), int64(21), int64(66)).Return(&script_entity.ScriptScoreReply{
					ID:        12,
					CommentID: 21,
					ScriptID:  66,
					Message:   "初始化",
				}, nil)
				mockScore.EXPECT().UpdateReplayByComment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, reply *script_entity.ScriptScoreReply) error {
					if reply.Message == "测试" {
						return nil
					}
					return errors.New("数据不匹配")
				})
				err := testMux.Do(ctx, &script.ReplyScoreRequest{ScriptId: 66, CommentID: 21, Message: "测试"}, &script.ReplyScoreResponse{})
				convey.So(err, convey.ShouldBeNil)
			})
		})

	})

}
