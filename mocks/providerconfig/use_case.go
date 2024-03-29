// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	pagination "github.com/octopipe/cloudx/internal/pagination"
	mock "github.com/stretchr/testify/mock"

	providerconfig "github.com/octopipe/cloudx/internal/providerconfig"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, providerConfig
func (_m *UseCase) Create(ctx context.Context, providerConfig providerconfig.ProviderConfig) (providerconfig.ProviderConfig, error) {
	ret := _m.Called(ctx, providerConfig)

	var r0 providerconfig.ProviderConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, providerconfig.ProviderConfig) (providerconfig.ProviderConfig, error)); ok {
		return rf(ctx, providerConfig)
	}
	if rf, ok := ret.Get(0).(func(context.Context, providerconfig.ProviderConfig) providerconfig.ProviderConfig); ok {
		r0 = rf(ctx, providerConfig)
	} else {
		r0 = ret.Get(0).(providerconfig.ProviderConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, providerconfig.ProviderConfig) error); ok {
		r1 = rf(ctx, providerConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, name, namespace
func (_m *UseCase) Delete(ctx context.Context, name string, namespace string) error {
	ret := _m.Called(ctx, name, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, name, namespace
func (_m *UseCase) Get(ctx context.Context, name string, namespace string) (providerconfig.ProviderConfig, error) {
	ret := _m.Called(ctx, name, namespace)

	var r0 providerconfig.ProviderConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (providerconfig.ProviderConfig, error)); ok {
		return rf(ctx, name, namespace)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) providerconfig.ProviderConfig); ok {
		r0 = rf(ctx, name, namespace)
	} else {
		r0 = ret.Get(0).(providerconfig.ProviderConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, namespace, chunkPagination
func (_m *UseCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[providerconfig.ProviderConfig], error) {
	ret := _m.Called(ctx, namespace, chunkPagination)

	var r0 pagination.ChunkingPaginationResponse[providerconfig.ProviderConfig]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[providerconfig.ProviderConfig], error)); ok {
		return rf(ctx, namespace, chunkPagination)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, pagination.ChunkingPaginationRequest) pagination.ChunkingPaginationResponse[providerconfig.ProviderConfig]); ok {
		r0 = rf(ctx, namespace, chunkPagination)
	} else {
		r0 = ret.Get(0).(pagination.ChunkingPaginationResponse[providerconfig.ProviderConfig])
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, pagination.ChunkingPaginationRequest) error); ok {
		r1 = rf(ctx, namespace, chunkPagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, providerConfig
func (_m *UseCase) Update(ctx context.Context, providerConfig providerconfig.ProviderConfig) (providerconfig.ProviderConfig, error) {
	ret := _m.Called(ctx, providerConfig)

	var r0 providerconfig.ProviderConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, providerconfig.ProviderConfig) (providerconfig.ProviderConfig, error)); ok {
		return rf(ctx, providerConfig)
	}
	if rf, ok := ret.Get(0).(func(context.Context, providerconfig.ProviderConfig) providerconfig.ProviderConfig); ok {
		r0 = rf(ctx, providerConfig)
	} else {
		r0 = ret.Get(0).(providerconfig.ProviderConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, providerconfig.ProviderConfig) error); ok {
		r1 = rf(ctx, providerConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
