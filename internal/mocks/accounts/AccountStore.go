// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	accounts "github.com/d-ashesss/mah-moneh/internal/accounts"

	mock "github.com/stretchr/testify/mock"

	users "github.com/d-ashesss/mah-moneh/internal/users"

	uuid "github.com/gofrs/uuid"
)

// AccountStore is an autogenerated mock type for the AccountStore type
type AccountStore struct {
	mock.Mock
}

// CreateAccount provides a mock function with given fields: ctx, acc
func (_m *AccountStore) CreateAccount(ctx context.Context, acc *accounts.Account) error {
	ret := _m.Called(ctx, acc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *accounts.Account) error); ok {
		r0 = rf(ctx, acc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAccount provides a mock function with given fields: ctx, acc
func (_m *AccountStore) DeleteAccount(ctx context.Context, acc *accounts.Account) error {
	ret := _m.Called(ctx, acc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *accounts.Account) error); ok {
		r0 = rf(ctx, acc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccount provides a mock function with given fields: ctx, UUID
func (_m *AccountStore) GetAccount(ctx context.Context, UUID uuid.UUID) (*accounts.Account, error) {
	ret := _m.Called(ctx, UUID)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *accounts.Account); ok {
		r0 = rf(ctx, UUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, UUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountAmounts provides a mock function with given fields: ctx, acc, month
func (_m *AccountStore) GetAccountAmounts(ctx context.Context, acc *accounts.Account, month string) (accounts.AmountCollection, error) {
	ret := _m.Called(ctx, acc, month)

	var r0 accounts.AmountCollection
	if rf, ok := ret.Get(0).(func(context.Context, *accounts.Account, string) accounts.AmountCollection); ok {
		r0 = rf(ctx, acc, month)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(accounts.AmountCollection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *accounts.Account, string) error); ok {
		r1 = rf(ctx, acc, month)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserAccounts provides a mock function with given fields: ctx, u
func (_m *AccountStore) GetUserAccounts(ctx context.Context, u *users.User) (accounts.AccountCollection, error) {
	ret := _m.Called(ctx, u)

	var r0 accounts.AccountCollection
	if rf, ok := ret.Get(0).(func(context.Context, *users.User) accounts.AccountCollection); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(accounts.AccountCollection)
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

// SetAccountAmount provides a mock function with given fields: ctx, acc, month, currency, amount
func (_m *AccountStore) SetAccountAmount(ctx context.Context, acc *accounts.Account, month string, currency string, amount float64) error {
	ret := _m.Called(ctx, acc, month, currency, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *accounts.Account, string, string, float64) error); ok {
		r0 = rf(ctx, acc, month, currency, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAccount provides a mock function with given fields: ctx, acc
func (_m *AccountStore) UpdateAccount(ctx context.Context, acc *accounts.Account) error {
	ret := _m.Called(ctx, acc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *accounts.Account) error); ok {
		r0 = rf(ctx, acc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAccountStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewAccountStore creates a new instance of AccountStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAccountStore(t mockConstructorTestingTNewAccountStore) *AccountStore {
	mock := &AccountStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
