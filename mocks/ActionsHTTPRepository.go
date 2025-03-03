// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "minireipaz/pkg/domain/models"

	mock "github.com/stretchr/testify/mock"
)

// ActionsHTTPRepository is an autogenerated mock type for the ActionsHTTPRepository type
type ActionsHTTPRepository struct {
	mock.Mock
}

// PublishCommand provides a mock function with given fields: data, serviceUser
func (_m *ActionsHTTPRepository) PublishCommand(data *models.ActionsCommand, serviceUser *string) *models.ResponseGetGoogleSheetByID {
	ret := _m.Called(data, serviceUser)

	if len(ret) == 0 {
		panic("no return value specified for PublishCommand")
	}

	var r0 *models.ResponseGetGoogleSheetByID
	if rf, ok := ret.Get(0).(func(*models.ActionsCommand, *string) *models.ResponseGetGoogleSheetByID); ok {
		r0 = rf(data, serviceUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResponseGetGoogleSheetByID)
		}
	}

	return r0
}

// SendAction provides a mock function with given fields: newAction, actionUserToken
func (_m *ActionsHTTPRepository) SendAction(newAction *models.RequestGoogleAction, actionUserToken *string) bool {
	ret := _m.Called(newAction, actionUserToken)

	if len(ret) == 0 {
		panic("no return value specified for SendAction")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*models.RequestGoogleAction, *string) bool); ok {
		r0 = rf(newAction, actionUserToken)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewActionsHTTPRepository creates a new instance of ActionsHTTPRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewActionsHTTPRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ActionsHTTPRepository {
	mock := &ActionsHTTPRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
