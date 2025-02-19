// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserHTTPRepository is an autogenerated mock type for the UserHTTPRepository type
type UserHTTPRepository struct {
	mock.Mock
}

// NewUserHTTPRepository creates a new instance of UserHTTPRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserHTTPRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserHTTPRepository {
	mock := &UserHTTPRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
