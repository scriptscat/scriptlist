package mock_auth_svc

import (
	"net/http"

	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model"
	"go.uber.org/mock/gomock"
)

func InitAuth(ctrl *gomock.Controller) *MockAuthSvc {
	mockAuth := NewMockAuthSvc(ctrl)
	mockAuth.EXPECT().RequireLogin(true).Return(gin.HandlerFunc(func(c *gin.Context) {
		if user := mockAuth.Get(c); user == nil {
			httputils.HandleResp(c, httputils.NewError(http.StatusUnauthorized, -1, "未登录"))
			return
		}
	})).AnyTimes()
	mockAuth.EXPECT().RequireLogin(false).Return(gin.HandlerFunc(func(c *gin.Context) {
		if user := mockAuth.Get(c); user == nil {
			return
		}
	})).AnyTimes()
	return mockAuth
}

type AuthUtil struct {
	m *MockAuthSvc
}

func (m *MockAuthSvc) U() *AuthUtil {
	return &AuthUtil{
		m,
	}
}

func (u *AuthUtil) NoLogin() *gomock.Call {
	return u.m.EXPECT().Get(gomock.Any()).Return(nil)
}

func (u *AuthUtil) Get(uids ...int64) *gomock.Call {
	uid := int64(1)
	if len(uids) > 0 {
		uid = uids[0]
	}
	return u.m.EXPECT().Get(gomock.Any()).Return(&model.AuthInfo{
		UID: uid,
	})

}
