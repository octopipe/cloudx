// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	manager "sigs.k8s.io/controller-runtime/pkg/manager"

	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Controller is an autogenerated mock type for the Controller type
type Controller struct {
	mock.Mock
}

// Reconcile provides a mock function with given fields: ctx, req
func (_m *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ret := _m.Called(ctx, req)

	var r0 reconcile.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, reconcile.Request) (reconcile.Result, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, reconcile.Request) reconcile.Result); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(reconcile.Result)
	}

	if rf, ok := ret.Get(1).(func(context.Context, reconcile.Request) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetupWithManager provides a mock function with given fields: mgr
func (_m *Controller) SetupWithManager(mgr manager.Manager) error {
	ret := _m.Called(mgr)

	var r0 error
	if rf, ok := ret.Get(0).(func(manager.Manager) error); ok {
		r0 = rf(mgr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewController creates a new instance of Controller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewController(t interface {
	mock.TestingT
	Cleanup(func())
}) *Controller {
	mock := &Controller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
