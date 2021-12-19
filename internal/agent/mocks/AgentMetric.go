// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	"github.com/stretchr/testify/mock"
)

// AgentMetric is an autogenerated mock type for the AgentMetric type
type AgentMetric struct {
	mock.Mock
}

// GetValue provides a mock function with given fields:
func (_m *AgentMetric) GetValue() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *AgentMetric) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Type provides a mock function with given fields:
func (_m *AgentMetric) Type() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
