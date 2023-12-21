package script_ctr

import (
	"context"
	"testing"

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
	"go.uber.org/mock/gomock"
)

func TestScore_Router(t *testing.T) {
	// 初始化路由
	testMux := muxtest.NewTestMux()
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
				err := testMux.Do(ctx, &script.DelScoreRequest{ScriptId: 1, ScoreId: 1}, &script.DelScoreResponse{})
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}
