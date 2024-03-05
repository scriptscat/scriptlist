// Code generated by MockGen. DO NOT EDIT.
// Source: user.go
//
// Generated by this command:
//
//	mockgen -source=user.go -destination=mock/user.go
//

// Package mock_user_repo is a generated GoMock package.
package mock_user_repo

import (
	context "context"
	reflect "reflect"

	user_entity "github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockUserRepo) Find(ctx context.Context, id int64) (*user_entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, id)
	ret0, _ := ret[0].(*user_entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockUserRepoMockRecorder) Find(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserRepo)(nil).Find), ctx, id)
}

// FindByPrefix mocks base method.
func (m *MockUserRepo) FindByPrefix(ctx context.Context, query string) ([]*user_entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPrefix", ctx, query)
	ret0, _ := ret[0].([]*user_entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPrefix indicates an expected call of FindByPrefix.
func (mr *MockUserRepoMockRecorder) FindByPrefix(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPrefix", reflect.TypeOf((*MockUserRepo)(nil).FindByPrefix), ctx, query)
}