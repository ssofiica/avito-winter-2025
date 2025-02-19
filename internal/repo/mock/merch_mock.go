// Code generated by MockGen. DO NOT EDIT.
// Source: merch.go

// Package mock is a generated GoMock package.
package mock

import (
	entity "avito-winter-2025/internal/entity"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMerchInterface is a mock of MerchInterface interface.
type MockMerchInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMerchInterfaceMockRecorder
}

// MockMerchInterfaceMockRecorder is the mock recorder for MockMerchInterface.
type MockMerchInterfaceMockRecorder struct {
	mock *MockMerchInterface
}

// NewMockMerchInterface creates a new mock instance.
func NewMockMerchInterface(ctrl *gomock.Controller) *MockMerchInterface {
	mock := &MockMerchInterface{ctrl: ctrl}
	mock.recorder = &MockMerchInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerchInterface) EXPECT() *MockMerchInterfaceMockRecorder {
	return m.recorder
}

// Buy mocks base method.
func (m *MockMerchInterface) Buy(ctx context.Context, userId, merchId, cost uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Buy", ctx, userId, merchId, cost)
	ret0, _ := ret[0].(error)
	return ret0
}

// Buy indicates an expected call of Buy.
func (mr *MockMerchInterfaceMockRecorder) Buy(ctx, userId, merchId, cost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buy", reflect.TypeOf((*MockMerchInterface)(nil).Buy), ctx, userId, merchId, cost)
}

// GetByName mocks base method.
func (m *MockMerchInterface) GetByName(ctx context.Context, name string) (*entity.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name)
	ret0, _ := ret[0].(*entity.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockMerchInterfaceMockRecorder) GetByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockMerchInterface)(nil).GetByName), ctx, name)
}

// GetInventoryHistory mocks base method.
func (m *MockMerchInterface) GetInventoryHistory(ctx context.Context, id uint32) ([]entity.Inventory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInventoryHistory", ctx, id)
	ret0, _ := ret[0].([]entity.Inventory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInventoryHistory indicates an expected call of GetInventoryHistory.
func (mr *MockMerchInterfaceMockRecorder) GetInventoryHistory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInventoryHistory", reflect.TypeOf((*MockMerchInterface)(nil).GetInventoryHistory), ctx, id)
}
