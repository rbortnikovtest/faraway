// Code generated by MockGen. DO NOT EDIT.
// Source: puzzler.go
//
// Generated by this command:
//
//	mockgen -source=puzzler.go -destination=puzzler_mock.go -package=service
//

// Package service is a generated GoMock package.
package service

import (
	rsa "crypto/rsa"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockKeyRepository is a mock of KeyRepository interface.
type MockKeyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockKeyRepositoryMockRecorder
}

// MockKeyRepositoryMockRecorder is the mock recorder for MockKeyRepository.
type MockKeyRepositoryMockRecorder struct {
	mock *MockKeyRepository
}

// NewMockKeyRepository creates a new mock instance.
func NewMockKeyRepository(ctrl *gomock.Controller) *MockKeyRepository {
	mock := &MockKeyRepository{ctrl: ctrl}
	mock.recorder = &MockKeyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyRepository) EXPECT() *MockKeyRepositoryMockRecorder {
	return m.recorder
}

// GetPrivateKey mocks base method.
func (m *MockKeyRepository) GetPrivateKey() (*rsa.PrivateKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivateKey")
	ret0, _ := ret[0].(*rsa.PrivateKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrivateKey indicates an expected call of GetPrivateKey.
func (mr *MockKeyRepositoryMockRecorder) GetPrivateKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivateKey", reflect.TypeOf((*MockKeyRepository)(nil).GetPrivateKey))
}