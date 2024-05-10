// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/wlachs/blog/internal/jwt (interfaces: TokenUtils)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTokenUtils is a mock of TokenUtils interface.
type MockTokenUtils struct {
	ctrl     *gomock.Controller
	recorder *MockTokenUtilsMockRecorder
}

// MockTokenUtilsMockRecorder is the mock recorder for MockTokenUtils.
type MockTokenUtilsMockRecorder struct {
	mock *MockTokenUtils
}

// NewMockTokenUtils creates a new mock instance.
func NewMockTokenUtils(ctrl *gomock.Controller) *MockTokenUtils {
	mock := &MockTokenUtils{ctrl: ctrl}
	mock.recorder = &MockTokenUtilsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenUtils) EXPECT() *MockTokenUtilsMockRecorder {
	return m.recorder
}

// GenerateJWT mocks base method.
func (m *MockTokenUtils) GenerateJWT(userName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateJWT", userName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateJWT indicates an expected call of GenerateJWT.
func (mr *MockTokenUtilsMockRecorder) GenerateJWT(userName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateJWT", reflect.TypeOf((*MockTokenUtils)(nil).GenerateJWT), userName)
}

// ParseJWT mocks base method.
func (m *MockTokenUtils) ParseJWT(t string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseJWT", t)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseJWT indicates an expected call of ParseJWT.
func (mr *MockTokenUtilsMockRecorder) ParseJWT(t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseJWT", reflect.TypeOf((*MockTokenUtils)(nil).ParseJWT), t)
}
