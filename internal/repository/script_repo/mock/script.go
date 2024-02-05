// Code generated by MockGen. DO NOT EDIT.
// Source: script.go
//
// Generated by this command:
//
//	mockgen -source=script.go -destination=mock/script.go
//

// Package mock_script_repo is a generated GoMock package.
package mock_script_repo

import (
	context "context"
	reflect "reflect"

	httputils "github.com/codfrm/cago/pkg/utils/httputils"
	script_entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	script_repo "github.com/scriptscat/scriptlist/internal/repository/script_repo"
	gomock "go.uber.org/mock/gomock"
)

// MockScriptRepo is a mock of ScriptRepo interface.
type MockScriptRepo struct {
	ctrl     *gomock.Controller
	recorder *MockScriptRepoMockRecorder
}

// MockScriptRepoMockRecorder is the mock recorder for MockScriptRepo.
type MockScriptRepoMockRecorder struct {
	mock *MockScriptRepo
}

// NewMockScriptRepo creates a new mock instance.
func NewMockScriptRepo(ctrl *gomock.Controller) *MockScriptRepo {
	mock := &MockScriptRepo{ctrl: ctrl}
	mock.recorder = &MockScriptRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScriptRepo) EXPECT() *MockScriptRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockScriptRepo) Create(ctx context.Context, script *script_entity.Script) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, script)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockScriptRepoMockRecorder) Create(ctx, script any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScriptRepo)(nil).Create), ctx, script)
}

// Delete mocks base method.
func (m *MockScriptRepo) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockScriptRepoMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockScriptRepo)(nil).Delete), ctx, id)
}

// Find mocks base method.
func (m *MockScriptRepo) Find(ctx context.Context, id int64) (*script_entity.Script, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, id)
	ret0, _ := ret[0].(*script_entity.Script)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockScriptRepoMockRecorder) Find(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockScriptRepo)(nil).Find), ctx, id)
}

// FindSyncPrefix mocks base method.
func (m *MockScriptRepo) FindSyncPrefix(ctx context.Context, uid int64, prefix string) ([]*script_entity.Script, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSyncPrefix", ctx, uid, prefix)
	ret0, _ := ret[0].([]*script_entity.Script)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSyncPrefix indicates an expected call of FindSyncPrefix.
func (mr *MockScriptRepoMockRecorder) FindSyncPrefix(ctx, uid, prefix any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSyncPrefix", reflect.TypeOf((*MockScriptRepo)(nil).FindSyncPrefix), ctx, uid, prefix)
}

// FindSyncScript mocks base method.
func (m *MockScriptRepo) FindSyncScript(ctx context.Context, page httputils.PageRequest) ([]*script_entity.Script, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSyncScript", ctx, page)
	ret0, _ := ret[0].([]*script_entity.Script)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSyncScript indicates an expected call of FindSyncScript.
func (mr *MockScriptRepoMockRecorder) FindSyncScript(ctx, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSyncScript", reflect.TypeOf((*MockScriptRepo)(nil).FindSyncScript), ctx, page)
}

// Search mocks base method.
func (m *MockScriptRepo) Search(ctx context.Context, options *script_repo.SearchOptions, page httputils.PageRequest) ([]*script_entity.Script, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, options, page)
	ret0, _ := ret[0].([]*script_entity.Script)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Search indicates an expected call of Search.
func (mr *MockScriptRepoMockRecorder) Search(ctx, options, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockScriptRepo)(nil).Search), ctx, options, page)
}

// Update mocks base method.
func (m *MockScriptRepo) Update(ctx context.Context, script *script_entity.Script) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, script)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockScriptRepoMockRecorder) Update(ctx, script any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockScriptRepo)(nil).Update), ctx, script)
}
