// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"

	categories "github.com/d-ashesss/mah-moneh/internal/categories"

	mock "github.com/stretchr/testify/mock"

	users "github.com/d-ashesss/mah-moneh/internal/users"

	uuid "github.com/gofrs/uuid"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// DeleteCategory provides a mock function with given fields: ctx, cat
func (_m *Store) DeleteCategory(ctx context.Context, cat *categories.Category) error {
	ret := _m.Called(ctx, cat)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *categories.Category) error); ok {
		r0 = rf(ctx, cat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCategory provides a mock function with given fields: ctx, _a1
func (_m *Store) GetCategory(ctx context.Context, _a1 uuid.UUID) (*categories.Category, error) {
	ret := _m.Called(ctx, _a1)

	var r0 *categories.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*categories.Category, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *categories.Category); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*categories.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserCategories provides a mock function with given fields: ctx, u
func (_m *Store) GetUserCategories(ctx context.Context, u *users.User) ([]*categories.Category, error) {
	ret := _m.Called(ctx, u)

	var r0 []*categories.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *users.User) ([]*categories.Category, error)); ok {
		return rf(ctx, u)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *users.User) []*categories.Category); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*categories.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *users.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveCategory provides a mock function with given fields: ctx, cat
func (_m *Store) SaveCategory(ctx context.Context, cat *categories.Category) error {
	ret := _m.Called(ctx, cat)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *categories.Category) error); ok {
		r0 = rf(ctx, cat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
