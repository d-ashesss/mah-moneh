// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	capital "github.com/d-ashesss/mah-moneh/internal/capital"

	mock "github.com/stretchr/testify/mock"

	users "github.com/d-ashesss/mah-moneh/internal/users"
)

// CapitalService is an autogenerated mock type for the CapitalService type
type CapitalService struct {
	mock.Mock
}

// GetCapital provides a mock function with given fields: ctx, u, month
func (_m *CapitalService) GetCapital(ctx context.Context, u *users.User, month string) (*capital.Capital, error) {
	ret := _m.Called(ctx, u, month)

	var r0 *capital.Capital
	if rf, ok := ret.Get(0).(func(context.Context, *users.User, string) *capital.Capital); ok {
		r0 = rf(ctx, u, month)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*capital.Capital)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *users.User, string) error); ok {
		r1 = rf(ctx, u, month)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCapitalService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCapitalService creates a new instance of CapitalService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCapitalService(t mockConstructorTestingTNewCapitalService) *CapitalService {
	mock := &CapitalService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
