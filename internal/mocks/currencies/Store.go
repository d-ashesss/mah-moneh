// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	currencies "github.com/d-ashesss/mah-moneh/internal/currencies"
	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// GetRate provides a mock function with given fields: base, target, month
func (_m *Store) GetRate(base string, target string, month string) (*currencies.Rate, error) {
	ret := _m.Called(base, target, month)

	var r0 *currencies.Rate
	if rf, ok := ret.Get(0).(func(string, string, string) *currencies.Rate); ok {
		r0 = rf(base, target, month)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*currencies.Rate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(base, target, month)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetRate provides a mock function with given fields: base, target, month, rate
func (_m *Store) SetRate(base string, target string, month string, rate float64) error {
	ret := _m.Called(base, target, month, rate)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, float64) error); ok {
		r0 = rf(base, target, month, rate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStore(t mockConstructorTestingTNewStore) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
