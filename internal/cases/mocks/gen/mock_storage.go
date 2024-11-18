// Code generated by MockGen. DO NOT EDIT.
// Source: ./storage.go
//
// Generated by this command:
//
//	mockgen -source=./storage.go -destination=./mocks/gen/mock_storage.go
//

// Package mock_cases is a generated GoMock package.
package mock_cases

import (
	context "context"
	entities "kursy-kriptovalyut/internal/entities"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
	isgomock struct{}
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetActualCoins mocks base method.
func (m *MockStorage) GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActualCoins", ctx, titles)
	ret0, _ := ret[0].([]entities.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActualCoins indicates an expected call of GetActualCoins.
func (mr *MockStorageMockRecorder) GetActualCoins(ctx, titles any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActualCoins", reflect.TypeOf((*MockStorage)(nil).GetActualCoins), ctx, titles)
}

// GetAggregateCoins mocks base method.
func (m *MockStorage) GetAggregateCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregateCoins", ctx, titles)
	ret0, _ := ret[0].([]entities.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAggregateCoins indicates an expected call of GetAggregateCoins.
func (mr *MockStorageMockRecorder) GetAggregateCoins(ctx, titles any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregateCoins", reflect.TypeOf((*MockStorage)(nil).GetAggregateCoins), ctx, titles)
}

// GetCoinsList mocks base method.
func (m *MockStorage) GetCoinsList(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoinsList", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoinsList indicates an expected call of GetCoinsList.
func (mr *MockStorageMockRecorder) GetCoinsList(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoinsList", reflect.TypeOf((*MockStorage)(nil).GetCoinsList), ctx)
}

// Store mocks base method.
func (m *MockStorage) Store(ctx context.Context, coins []entities.Coin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, coins)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockStorageMockRecorder) Store(ctx, coins any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockStorage)(nil).Store), ctx, coins)
}
