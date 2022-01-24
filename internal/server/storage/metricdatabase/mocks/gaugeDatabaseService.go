// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	metric "github.com/djokcik/praktikum-go-devops/internal/metric"

	mock "github.com/stretchr/testify/mock"
)

// GaugeDatabaseService is an autogenerated mock type for the gaugeDatabaseService type
type GaugeDatabaseService struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, name
func (_m *GaugeDatabaseService) Get(ctx context.Context, name string) (metric.Gauge, error) {
	ret := _m.Called(ctx, name)

	var r0 metric.Gauge
	if rf, ok := ret.Get(0).(func(context.Context, string) metric.Gauge); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(metric.Gauge)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx
func (_m *GaugeDatabaseService) List(ctx context.Context) ([]metric.Metric, error) {
	ret := _m.Called(ctx)

	var r0 []metric.Metric
	if rf, ok := ret.Get(0).(func(context.Context) []metric.Metric); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]metric.Metric)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, name, value
func (_m *GaugeDatabaseService) Update(ctx context.Context, name string, value metric.Gauge) error {
	ret := _m.Called(ctx, name, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metric.Gauge) error); ok {
		r0 = rf(ctx, name, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
