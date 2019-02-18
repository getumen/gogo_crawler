// Code generated by MockGen. DO NOT EDIT.
// Source: schedule_service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	models "github.com/getumen/gogo_crawler/domains/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockScheduleService is a mock of ScheduleService interface
type MockScheduleService struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleServiceMockRecorder
}

// MockScheduleServiceMockRecorder is the mock recorder for MockScheduleService
type MockScheduleServiceMockRecorder struct {
	mock *MockScheduleService
}

// NewMockScheduleService creates a new mock instance
func NewMockScheduleService(ctrl *gomock.Controller) *MockScheduleService {
	mock := &MockScheduleService{ctrl: ctrl}
	mock.recorder = &MockScheduleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScheduleService) EXPECT() *MockScheduleServiceMockRecorder {
	return m.recorder
}

// GenerateRequest mocks base method
func (m *MockScheduleService) GenerateRequest(ctx context.Context, namespace string, out chan<- *models.Request) {
	m.ctrl.Call(m, "GenerateRequest", ctx, namespace, out)
}

// GenerateRequest indicates an expected call of GenerateRequest
func (mr *MockScheduleServiceMockRecorder) GenerateRequest(ctx, namespace, out interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRequest", reflect.TypeOf((*MockScheduleService)(nil).GenerateRequest), ctx, namespace, out)
}

// ScheduleRequest mocks base method
func (m *MockScheduleService) ScheduleRequest(ctx context.Context, in <-chan *models.Response) {
	m.ctrl.Call(m, "ScheduleRequest", ctx, in)
}

// ScheduleRequest indicates an expected call of ScheduleRequest
func (mr *MockScheduleServiceMockRecorder) ScheduleRequest(ctx, in interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleRequest", reflect.TypeOf((*MockScheduleService)(nil).ScheduleRequest), ctx, in)
}

// ScheduleNewRequest mocks base method
func (m *MockScheduleService) ScheduleNewRequest(ctx context.Context, in <-chan *models.Request) {
	m.ctrl.Call(m, "ScheduleNewRequest", ctx, in)
}

// ScheduleNewRequest indicates an expected call of ScheduleNewRequest
func (mr *MockScheduleServiceMockRecorder) ScheduleNewRequest(ctx, in interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleNewRequest", reflect.TypeOf((*MockScheduleService)(nil).ScheduleNewRequest), ctx, in)
}

// MockScheduleRule is a mock of ScheduleRule interface
type MockScheduleRule struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleRuleMockRecorder
}

// MockScheduleRuleMockRecorder is the mock recorder for MockScheduleRule
type MockScheduleRuleMockRecorder struct {
	mock *MockScheduleRule
}

// NewMockScheduleRule creates a new mock instance
func NewMockScheduleRule(ctrl *gomock.Controller) *MockScheduleRule {
	mock := &MockScheduleRule{ctrl: ctrl}
	mock.recorder = &MockScheduleRuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScheduleRule) EXPECT() *MockScheduleRuleMockRecorder {
	return m.recorder
}

// UpdateStatsSuccess mocks base method
func (m *MockScheduleRule) UpdateStatsSuccess(request *models.Request) {
	m.ctrl.Call(m, "UpdateStatsSuccess", request)
}

// UpdateStatsSuccess indicates an expected call of UpdateStatsSuccess
func (mr *MockScheduleRuleMockRecorder) UpdateStatsSuccess(request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatsSuccess", reflect.TypeOf((*MockScheduleRule)(nil).UpdateStatsSuccess), request)
}

// UpdateStatsFail mocks base method
func (m *MockScheduleRule) UpdateStatsFail(request *models.Request) {
	m.ctrl.Call(m, "UpdateStatsFail", request)
}

// UpdateStatsFail indicates an expected call of UpdateStatsFail
func (mr *MockScheduleRuleMockRecorder) UpdateStatsFail(request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatsFail", reflect.TypeOf((*MockScheduleRule)(nil).UpdateStatsFail), request)
}

// ScheduleNextSuccess mocks base method
func (m *MockScheduleRule) ScheduleNextSuccess(request *models.Request) time.Time {
	ret := m.ctrl.Call(m, "ScheduleNextSuccess", request)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// ScheduleNextSuccess indicates an expected call of ScheduleNextSuccess
func (mr *MockScheduleRuleMockRecorder) ScheduleNextSuccess(request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleNextSuccess", reflect.TypeOf((*MockScheduleRule)(nil).ScheduleNextSuccess), request)
}

// ScheduleNextFail mocks base method
func (m *MockScheduleRule) ScheduleNextFail(request *models.Request) time.Time {
	ret := m.ctrl.Call(m, "ScheduleNextFail", request)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// ScheduleNextFail indicates an expected call of ScheduleNextFail
func (mr *MockScheduleRuleMockRecorder) ScheduleNextFail(request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleNextFail", reflect.TypeOf((*MockScheduleRule)(nil).ScheduleNextFail), request)
}
