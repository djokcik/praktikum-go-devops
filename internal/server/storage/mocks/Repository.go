// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	storage "github.com/djokcik/praktikum-go-devops/internal/server/storage"
	model "github.com/djokcik/praktikum-go-devops/internal/server/storage/model"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Configure provides a mock function with given fields: db
func (_m *Repository) Configure(db *model.Database) {
	_m.Called(db)
}

// Get provides a mock function with given fields: filter
func (_m *Repository) Get(filter *storage.GetRepositoryFilter) (interface{}, error) {
	ret := _m.Called(filter)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(*storage.GetRepositoryFilter) interface{}); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*storage.GetRepositoryFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: filter
func (_m *Repository) List(filter *storage.ListRepositoryFilter) (interface{}, error) {
	ret := _m.Called(filter)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(*storage.ListRepositoryFilter) interface{}); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*storage.ListRepositoryFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, entity
func (_m *Repository) Update(id interface{}, entity interface{}) (bool, error) {
	ret := _m.Called(id, entity)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) bool); ok {
		r0 = rf(id, entity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, interface{}) error); ok {
		r1 = rf(id, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
