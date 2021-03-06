// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nsmak/bannersRotation/internal/app (interfaces: Storage)

// Package app_test is a generated GoMock package.
package app_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	app "github.com/nsmak/bannersRotation/internal/app"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddBannerToSlot mocks base method
func (m *MockStorage) AddBannerToSlot(arg0 context.Context, arg1, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBannerToSlot", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBannerToSlot indicates an expected call of AddBannerToSlot
func (mr *MockStorageMockRecorder) AddBannerToSlot(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBannerToSlot", reflect.TypeOf((*MockStorage)(nil).AddBannerToSlot), arg0, arg1, arg2)
}

// AddClickForBanner mocks base method
func (m *MockStorage) AddClickForBanner(arg0 context.Context, arg1, arg2, arg3 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddClickForBanner", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddClickForBanner indicates an expected call of AddClickForBanner
func (mr *MockStorageMockRecorder) AddClickForBanner(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClickForBanner", reflect.TypeOf((*MockStorage)(nil).AddClickForBanner), arg0, arg1, arg2, arg3)
}

// AddViewForBanner mocks base method
func (m *MockStorage) AddViewForBanner(arg0 context.Context, arg1, arg2, arg3 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddViewForBanner", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddViewForBanner indicates an expected call of AddViewForBanner
func (mr *MockStorageMockRecorder) AddViewForBanner(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddViewForBanner", reflect.TypeOf((*MockStorage)(nil).AddViewForBanner), arg0, arg1, arg2, arg3)
}

// BannersClickStatisticsFilterByDate mocks base method
func (m *MockStorage) BannersClickStatisticsFilterByDate(arg0 context.Context, arg1, arg2 int64) ([]app.BannerStatistic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BannersClickStatisticsFilterByDate", arg0, arg1, arg2)
	ret0, _ := ret[0].([]app.BannerStatistic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BannersClickStatisticsFilterByDate indicates an expected call of BannersClickStatisticsFilterByDate
func (mr *MockStorageMockRecorder) BannersClickStatisticsFilterByDate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BannersClickStatisticsFilterByDate", reflect.TypeOf((*MockStorage)(nil).BannersClickStatisticsFilterByDate), arg0, arg1, arg2)
}

// BannersShowStatisticsFilterByDate mocks base method
func (m *MockStorage) BannersShowStatisticsFilterByDate(arg0 context.Context, arg1, arg2 int64) ([]app.BannerStatistic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BannersShowStatisticsFilterByDate", arg0, arg1, arg2)
	ret0, _ := ret[0].([]app.BannerStatistic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BannersShowStatisticsFilterByDate indicates an expected call of BannersShowStatisticsFilterByDate
func (mr *MockStorageMockRecorder) BannersShowStatisticsFilterByDate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BannersShowStatisticsFilterByDate", reflect.TypeOf((*MockStorage)(nil).BannersShowStatisticsFilterByDate), arg0, arg1, arg2)
}

// BannersStatistics mocks base method
func (m *MockStorage) BannersStatistics(arg0 context.Context, arg1, arg2 int64) ([]app.BannerSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BannersStatistics", arg0, arg1, arg2)
	ret0, _ := ret[0].([]app.BannerSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BannersStatistics indicates an expected call of BannersStatistics
func (mr *MockStorageMockRecorder) BannersStatistics(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BannersStatistics", reflect.TypeOf((*MockStorage)(nil).BannersStatistics), arg0, arg1, arg2)
}

// RemoveBannerFromSlot mocks base method
func (m *MockStorage) RemoveBannerFromSlot(arg0 context.Context, arg1, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveBannerFromSlot", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveBannerFromSlot indicates an expected call of RemoveBannerFromSlot
func (mr *MockStorageMockRecorder) RemoveBannerFromSlot(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveBannerFromSlot", reflect.TypeOf((*MockStorage)(nil).RemoveBannerFromSlot), arg0, arg1, arg2)
}
