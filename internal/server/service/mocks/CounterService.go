// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	metric "github.com/djokcik/praktikum-go-devops/internal/metric"
	mock "github.com/stretchr/testify/mock"
)

// CounterService is an autogenerated mock type for the CounterService type
type CounterService struct {
	mock.Mock
}

// AddValue provides a mock function with given fields: name, value
func (_m *CounterService) AddValue(name string, value metric.Counter) error {
	ret := _m.Called(name, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, metric.Counter) error); ok {
		r0 = rf(name, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetOne provides a mock function with given fields: name
func (_m *CounterService) GetOne(name string) (metric.Counter, error) {
	ret := _m.Called(name)

	var r0 metric.Counter
	if rf, ok := ret.Get(0).(func(string) metric.Counter); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(metric.Counter)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields:
func (_m *CounterService) List() ([]metric.Metric, error) {
	ret := _m.Called()

	var r0 []metric.Metric
	if rf, ok := ret.Get(0).(func() []metric.Metric); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]metric.Metric)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: name, value
func (_m *CounterService) Update(name string, value metric.Counter) (bool, error) {
	ret := _m.Called(name, value)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, metric.Counter) bool); ok {
		r0 = rf(name, value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, metric.Counter) error); ok {
		r1 = rf(name, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
