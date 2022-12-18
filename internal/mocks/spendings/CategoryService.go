// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	categories "github.com/d-ashesss/mah-moneh/internal/categories"

	mock "github.com/stretchr/testify/mock"

	users "github.com/d-ashesss/mah-moneh/internal/users"
)

// CategoryService is an autogenerated mock type for the CategoryService type
type CategoryService struct {
	mock.Mock
}

// GetUserCategories provides a mock function with given fields: ctx, u
func (_m *CategoryService) GetUserCategories(ctx context.Context, u *users.User) ([]*categories.Category, error) {
	ret := _m.Called(ctx, u)

	var r0 []*categories.Category
	if rf, ok := ret.Get(0).(func(context.Context, *users.User) []*categories.Category); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*categories.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *users.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCategoryService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCategoryService creates a new instance of CategoryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCategoryService(t mockConstructorTestingTNewCategoryService) *CategoryService {
	mock := &CategoryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}