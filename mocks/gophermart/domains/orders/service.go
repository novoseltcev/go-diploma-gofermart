// Code generated by MockGen. DO NOT EDIT.
// Source: gophermart/domains/orders/service.go

// Package mock_orders is a generated GoMock package.
package mock_orders

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "github.com/novoseltcev/go-diploma-gofermart/gophermart/models"
	reflect "reflect"
)

// MockOrderStorager is a mock of OrderStorager interface
type MockOrderStorager struct {
	ctrl     *gomock.Controller
	recorder *MockOrderStoragerMockRecorder
}

// MockOrderStoragerMockRecorder is the mock recorder for MockOrderStorager
type MockOrderStoragerMockRecorder struct {
	mock *MockOrderStorager
}

// NewMockOrderStorager creates a new mock instance
func NewMockOrderStorager(ctrl *gomock.Controller) *MockOrderStorager {
	mock := &MockOrderStorager{ctrl: ctrl}
	mock.recorder = &MockOrderStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOrderStorager) EXPECT() *MockOrderStoragerMockRecorder {
	return m.recorder
}

// GetUserOrders mocks base method
func (m *MockOrderStorager) GetUserOrders(ctx context.Context, userID models.UserID) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", ctx, userID)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders
func (mr *MockOrderStoragerMockRecorder) GetUserOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockOrderStorager)(nil).GetUserOrders), ctx, userID)
}

// GetByNumber mocks base method
func (m *MockOrderStorager) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNumber", ctx, number)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNumber indicates an expected call of GetByNumber
func (mr *MockOrderStoragerMockRecorder) GetByNumber(ctx, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNumber", reflect.TypeOf((*MockOrderStorager)(nil).GetByNumber), ctx, number)
}

// Create mocks base method
func (m *MockOrderStorager) Create(ctx context.Context, userID models.UserID, number string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userID, number)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockOrderStoragerMockRecorder) Create(ctx, userID, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderStorager)(nil).Create), ctx, userID, number)
}
