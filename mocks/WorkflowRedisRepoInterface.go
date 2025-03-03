// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "minireipaz/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/google/uuid"
)

// WorkflowRedisRepoInterface is an autogenerated mock type for the WorkflowRedisRepoInterface type
type WorkflowRedisRepoInterface struct {
	mock.Mock
}

// AcquireLock provides a mock function with given fields: key, value, expiration
func (_m *WorkflowRedisRepoInterface) AcquireLock(key string, value string, expiration time.Duration) (bool, error) {
	ret := _m.Called(key, value, expiration)

	if len(ret) == 0 {
		panic("no return value specified for AcquireLock")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) (bool, error)); ok {
		return rf(key, value, expiration)
	}
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) bool); ok {
		r0 = rf(key, value, expiration)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string, time.Duration) error); ok {
		r1 = rf(key, value, expiration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: workflow
func (_m *WorkflowRedisRepoInterface) Create(workflow *models.Workflow) (bool, bool) {
	ret := _m.Called(workflow)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 bool
	var r1 bool
	if rf, ok := ret.Get(0).(func(*models.Workflow) (bool, bool)); ok {
		return rf(workflow)
	}
	if rf, ok := ret.Get(0).(func(*models.Workflow) bool); ok {
		r0 = rf(workflow)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.Workflow) bool); ok {
		r1 = rf(workflow)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetByUUID provides a mock function with given fields: id
func (_m *WorkflowRedisRepoInterface) GetByUUID(id uuid.UUID) (*models.Workflow, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByUUID")
	}

	var r0 *models.Workflow
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Workflow, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Workflow); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Workflow)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: workflow
func (_m *WorkflowRedisRepoInterface) Remove(workflow *models.Workflow) bool {
	ret := _m.Called(workflow)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*models.Workflow) bool); ok {
		r0 = rf(workflow)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RemoveLock provides a mock function with given fields: key
func (_m *WorkflowRedisRepoInterface) RemoveLock(key string) bool {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for RemoveLock")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Update provides a mock function with given fields: worflow
func (_m *WorkflowRedisRepoInterface) Update(worflow *models.Workflow) (bool, bool) {
	ret := _m.Called(worflow)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 bool
	var r1 bool
	if rf, ok := ret.Get(0).(func(*models.Workflow) (bool, bool)); ok {
		return rf(worflow)
	}
	if rf, ok := ret.Get(0).(func(*models.Workflow) bool); ok {
		r0 = rf(worflow)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.Workflow) bool); ok {
		r1 = rf(worflow)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ValidateUserWorkflowUUID provides a mock function with given fields: userID, name
func (_m *WorkflowRedisRepoInterface) ValidateUserWorkflowUUID(userID *string, name *string) bool {
	ret := _m.Called(userID, name)

	if len(ret) == 0 {
		panic("no return value specified for ValidateUserWorkflowUUID")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*string, *string) bool); ok {
		r0 = rf(userID, name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ValidateWorkflowGlobalUUID provides a mock function with given fields: _a0
func (_m *WorkflowRedisRepoInterface) ValidateWorkflowGlobalUUID(_a0 *string) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ValidateWorkflowGlobalUUID")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewWorkflowRedisRepoInterface creates a new instance of WorkflowRedisRepoInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWorkflowRedisRepoInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *WorkflowRedisRepoInterface {
	mock := &WorkflowRedisRepoInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
