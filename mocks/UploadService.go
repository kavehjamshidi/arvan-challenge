// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kavehjamshidi/arvan-challenge/domain"

	mock "github.com/stretchr/testify/mock"
)

// UploadService is an autogenerated mock type for the UploadService type
type UploadService struct {
	mock.Mock
}

// UploadFile provides a mock function with given fields: ctx, file
func (_m *UploadService) UploadFile(ctx context.Context, file domain.File) error {
	ret := _m.Called(ctx, file)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.File) error); ok {
		r0 = rf(ctx, file)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUploadService creates a new instance of UploadService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUploadService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UploadService {
	mock := &UploadService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}