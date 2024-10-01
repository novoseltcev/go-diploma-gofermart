// Code generated by MockGen. DO NOT EDIT.
// Source: gophermart/domains/users/service.go

// Package mock_users is a generated GoMock package.
package mock_users

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	reflect "reflect"
)

// MockUserStorager is a mock of UserStorager interface
type MockUserStorager struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoragerMockRecorder
}

// MockUserStoragerMockRecorder is the mock recorder for MockUserStorager
type MockUserStoragerMockRecorder struct {
	mock *MockUserStorager
}

// NewMockUserStorager creates a new mock instance
func NewMockUserStorager(ctrl *gomock.Controller) *MockUserStorager {
	mock := &MockUserStorager{ctrl: ctrl}
	mock.recorder = &MockUserStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserStorager) EXPECT() *MockUserStoragerMockRecorder {
	return m.recorder
}

// IsLoginExists mocks base method
func (m *MockUserStorager) IsLoginExists(ctx context.Context, login string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLoginExists", ctx, login)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLoginExists indicates an expected call of IsLoginExists
func (mr *MockUserStoragerMockRecorder) IsLoginExists(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLoginExists", reflect.TypeOf((*MockUserStorager)(nil).IsLoginExists), ctx, login)
}

// GetByLogin mocks base method
func (m *MockUserStorager) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLogin", ctx, login)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLogin indicates an expected call of GetByLogin
func (mr *MockUserStoragerMockRecorder) GetByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLogin", reflect.TypeOf((*MockUserStorager)(nil).GetByLogin), ctx, login)
}

// Create mocks base method
func (m *MockUserStorager) Create(ctx context.Context, login, password string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, login, password)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockUserStoragerMockRecorder) Create(ctx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserStorager)(nil).Create), ctx, login, password)
}
