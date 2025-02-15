// Code generated by MockGen. DO NOT EDIT.
// Source: coin.go

// Package mock is a generated GoMock package.
package mock

import (
	entity "avito-winter-2025/internal/entity"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCoinInterface is a mock of CoinInterface interface.
type MockCoinInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCoinInterfaceMockRecorder
}

// MockCoinInterfaceMockRecorder is the mock recorder for MockCoinInterface.
type MockCoinInterfaceMockRecorder struct {
	mock *MockCoinInterface
}

// NewMockCoinInterface creates a new mock instance.
func NewMockCoinInterface(ctrl *gomock.Controller) *MockCoinInterface {
	mock := &MockCoinInterface{ctrl: ctrl}
	mock.recorder = &MockCoinInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoinInterface) EXPECT() *MockCoinInterfaceMockRecorder {
	return m.recorder
}

// CheckBalance mocks base method.
func (m *MockCoinInterface) CheckBalance(ctx context.Context, id uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckBalance", ctx, id)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckBalance indicates an expected call of CheckBalance.
func (mr *MockCoinInterfaceMockRecorder) CheckBalance(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckBalance", reflect.TypeOf((*MockCoinInterface)(nil).CheckBalance), ctx, id)
}

// GetCoinHistory mocks base method.
func (m *MockCoinInterface) GetCoinHistory(ctx context.Context, id uint32) ([]entity.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoinHistory", ctx, id)
	ret0, _ := ret[0].([]entity.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoinHistory indicates an expected call of GetCoinHistory.
func (mr *MockCoinInterfaceMockRecorder) GetCoinHistory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoinHistory", reflect.TypeOf((*MockCoinInterface)(nil).GetCoinHistory), ctx, id)
}

// SendCoin mocks base method.
func (m *MockCoinInterface) SendCoin(ctx context.Context, transaction entity.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoin", ctx, transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoin indicates an expected call of SendCoin.
func (mr *MockCoinInterfaceMockRecorder) SendCoin(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoin", reflect.TypeOf((*MockCoinInterface)(nil).SendCoin), ctx, transaction)
}
