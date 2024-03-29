// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	v1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

// Pipeline is an autogenerated mock type for the Pipeline type
type Pipeline struct {
	mock.Mock
}

// Start provides a mock function with given fields: action, infra, taskStatusChan
func (_m *Pipeline) Start(action string, infra v1alpha1.Infra, taskStatusChan chan v1alpha1.TaskExecutionStatus) v1alpha1.ExecutionStatus {
	ret := _m.Called(action, infra, taskStatusChan)

	var r0 v1alpha1.ExecutionStatus
	if rf, ok := ret.Get(0).(func(string, v1alpha1.Infra, chan v1alpha1.TaskExecutionStatus) v1alpha1.ExecutionStatus); ok {
		r0 = rf(action, infra, taskStatusChan)
	} else {
		r0 = ret.Get(0).(v1alpha1.ExecutionStatus)
	}

	return r0
}

// NewPipeline creates a new instance of Pipeline. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipeline(t interface {
	mock.TestingT
	Cleanup(func())
}) *Pipeline {
	mock := &Pipeline{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
