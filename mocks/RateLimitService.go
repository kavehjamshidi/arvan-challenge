// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// RateLimitService is an autogenerated mock type for the RateLimitService type
type RateLimitService struct {
	mock.Mock
}

// CheckRateLimit provides a mock function with given fields: ctx, userID
func (_m *RateLimitService) CheckRateLimit(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRateLimitService creates a new instance of RateLimitService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRateLimitService(t interface {
	mock.TestingT
	Cleanup(func())
}) *RateLimitService {
	mock := &RateLimitService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}