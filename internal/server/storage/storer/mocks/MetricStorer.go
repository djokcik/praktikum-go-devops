// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MetricStorer is an autogenerated mock type for the MetricStorer type
type MetricStorer struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MetricStorer) Close() {
	_m.Called()
}

// RestoreDBValue provides a mock function with given fields: ctx
func (_m *MetricStorer) RestoreDBValue(ctx context.Context) {
	_m.Called(ctx)
}

// SaveDBValue provides a mock function with given fields: ctx
func (_m *MetricStorer) SaveDBValue(ctx context.Context) {
	_m.Called(ctx)
}