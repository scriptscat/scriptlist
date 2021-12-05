package subscriber

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	service3 "github.com/scriptscat/scriptlist/internal/domain/user/service"
	mock_service "github.com/scriptscat/scriptlist/internal/domain/user/service/mock"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
)

func TestScriptSubscriber_parseContent(t *testing.T) {
	mockctl := gomock.NewController(t)
	mock := mock_service.NewMockUser(mockctl)
	mock.EXPECT().FindByUsername("CodFrm", true).Return(&respond.User{UID: 1, Username: "CodFrm"}, nil).AnyTimes()
	mock.EXPECT().FindByUsername(gomock.Any(), true).Return(nil, errs.ErrUserNotFound).AnyTimes()
	type fields struct {
		userSvc service3.User
	}
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*respond.User
		wantErr bool
	}{
		{
			name: "用户首部艾特", fields: fields{userSvc: mock}, args: args{content: "@CodFrm 艾特用户"}, want: []*respond.User{{
				UID:      1,
				Username: "CodFrm",
				Avatar:   "",
				IsAdmin:  0,
				Email:    "",
			}}, wantErr: false,
		}, {
			name: "用户中艾特", fields: fields{userSvc: mock}, args: args{content: "艾特用户@CodFrm 艾特用户"}, want: []*respond.User{{
				UID:      1,
				Username: "CodFrm",
				Avatar:   "",
				IsAdmin:  0,
				Email:    "",
			}}, wantErr: false,
		}, {
			name: "用户结尾艾特", fields: fields{userSvc: mock}, args: args{content: "艾特用户@CodFrm"}, want: []*respond.User{{
				UID:      1,
				Username: "CodFrm",
				Avatar:   "",
				IsAdmin:  0,
				Email:    "",
			}}, wantErr: false,
		}, {
			name: "未找到", fields: fields{userSvc: mock}, args: args{content: "艾特用户@CodFrm2"}, want: nil, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &ScriptSubscriber{
				userSvc: tt.fields.userSvc,
			}
			got, err := n.parseContent(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}