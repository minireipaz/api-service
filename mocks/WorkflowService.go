// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "minireipaz/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// WorkflowService is an autogenerated mock type for the WorkflowService type
type WorkflowService struct {
	mock.Mock
}

// CreateWorkflow provides a mock function with given fields: workflowFrontend
func (_m *WorkflowService) CreateWorkflow(workflowFrontend *models.WorkflowFrontend) (bool, bool, *models.Workflow) {
	ret := _m.Called(workflowFrontend)

	if len(ret) == 0 {
		panic("no return value specified for CreateWorkflow")
	}

	var r0 bool
	var r1 bool
	var r2 *models.Workflow
	if rf, ok := ret.Get(0).(func(*models.WorkflowFrontend) (bool, bool, *models.Workflow)); ok {
		return rf(workflowFrontend)
	}
	if rf, ok := ret.Get(0).(func(*models.WorkflowFrontend) bool); ok {
		r0 = rf(workflowFrontend)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*models.WorkflowFrontend) bool); ok {
		r1 = rf(workflowFrontend)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(*models.WorkflowFrontend) *models.Workflow); ok {
		r2 = rf(workflowFrontend)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*models.Workflow)
		}
	}

	return r0, r1, r2
}

// GetAllWorkflows provides a mock function with given fields: userID
func (_m *WorkflowService) GetAllWorkflows(userID *string) ([]models.Workflow, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAllWorkflows")
	}

	var r0 []models.Workflow
	var r1 error
	if rf, ok := ret.Get(0).(func(*string) ([]models.Workflow, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(*string) []models.Workflow); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Workflow)
		}
	}

	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWorkflow provides a mock function with given fields: userID, workflowID
func (_m *WorkflowService) GetWorkflow(userID *string, workflowID *string) (*models.Workflow, bool) {
	ret := _m.Called(userID, workflowID)

	if len(ret) == 0 {
		panic("no return value specified for GetWorkflow")
	}

	var r0 *models.Workflow
	var r1 bool
	if rf, ok := ret.Get(0).(func(*string, *string) (*models.Workflow, bool)); ok {
		return rf(userID, workflowID)
	}
	if rf, ok := ret.Get(0).(func(*string, *string) *models.Workflow); ok {
		r0 = rf(userID, workflowID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Workflow)
		}
	}

	if rf, ok := ret.Get(1).(func(*string, *string) bool); ok {
		r1 = rf(userID, workflowID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// UpdateWorkflow provides a mock function with given fields: workflow
func (_m *WorkflowService) UpdateWorkflow(workflow *models.Workflow) (bool, bool) {
	ret := _m.Called(workflow)

	if len(ret) == 0 {
		panic("no return value specified for UpdateWorkflow")
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

// ValidateUserWorkflowUUID provides a mock function with given fields: worklfowID, name
func (_m *WorkflowService) ValidateUserWorkflowUUID(worklfowID *string, name *string) bool {
	ret := _m.Called(worklfowID, name)

	if len(ret) == 0 {
		panic("no return value specified for ValidateUserWorkflowUUID")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*string, *string) bool); ok {
		r0 = rf(worklfowID, name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ValidateWorkflowGlobalUUID provides a mock function with given fields: uuid
func (_m *WorkflowService) ValidateWorkflowGlobalUUID(uuid *string) bool {
	ret := _m.Called(uuid)

	if len(ret) == 0 {
		panic("no return value specified for ValidateWorkflowGlobalUUID")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*string) bool); ok {
		r0 = rf(uuid)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewWorkflowService creates a new instance of WorkflowService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWorkflowService(t interface {
	mock.TestingT
	Cleanup(func())
}) *WorkflowService {
	mock := &WorkflowService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
