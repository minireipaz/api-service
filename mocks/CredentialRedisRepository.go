// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "minireipaz/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// CredentialRedisRepository is an autogenerated mock type for the CredentialRedisRepository type
type CredentialRedisRepository struct {
	mock.Mock
}

// AddLock provides a mock function with given fields: sub
func (_m *CredentialRedisRepository) AddLock(sub *string) (bool, error) {
	ret := _m.Called(sub)

	if len(ret) == 0 {
		panic("no return value specified for AddLock")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*string) (bool, error)); ok {
		return rf(sub)
	}
	if rf, ok := ret.Get(0).(func(*string) bool); ok {
		r0 = rf(sub)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(sub)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveTemporalAuthURLData provides a mock function with given fields: currentCredential
func (_m *CredentialRedisRepository) SaveTemporalAuthURLData(currentCredential *models.RequestCreateCredential) (bool, error) {
	ret := _m.Called(currentCredential)

	if len(ret) == 0 {
		panic("no return value specified for SaveTemporalAuthURLData")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.RequestCreateCredential) (bool, error)); ok {
		return rf(currentCredential)
	}
	if rf, ok := ret.Get(0).(func(*models.RequestCreateCredential) bool); ok {
		r0 = rf(currentCredential)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.RequestCreateCredential) error); ok {
		r1 = rf(currentCredential)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCredentialRedisRepository creates a new instance of CredentialRedisRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCredentialRedisRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CredentialRedisRepository {
	mock := &CredentialRedisRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
