// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "minireipaz/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// UserRedisRepository is an autogenerated mock type for the UserRedisRepository type
type UserRedisRepository struct {
	mock.Mock
}

// CheckLockExist provides a mock function with given fields: user
func (_m *UserRedisRepository) CheckLockExist(user *models.SyncUserRequest) (bool, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for CheckLockExist")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) (bool, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.SyncUserRequest) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckUserExist provides a mock function with given fields: user
func (_m *UserRedisRepository) CheckUserExist(user *models.SyncUserRequest) (bool, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for CheckUserExist")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) (bool, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.SyncUserRequest) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUser provides a mock function with given fields: user
func (_m *UserRedisRepository) InsertUser(user *models.SyncUserRequest) (bool, bool, bool, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for InsertUser")
	}

	var r0 bool
	var r1 bool
	var r2 bool
	var r3 error
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) (bool, bool, bool, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.SyncUserRequest) bool); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(*models.SyncUserRequest) bool); ok {
		r2 = rf(user)
	} else {
		r2 = ret.Get(2).(bool)
	}

	if rf, ok := ret.Get(3).(func(*models.SyncUserRequest) error); ok {
		r3 = rf(user)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// RemoveLock provides a mock function with given fields: user
func (_m *UserRedisRepository) RemoveLock(user *models.SyncUserRequest) bool {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for RemoveLock")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*models.SyncUserRequest) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewUserRedisRepository creates a new instance of UserRedisRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRedisRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRedisRepository {
	mock := &UserRedisRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
