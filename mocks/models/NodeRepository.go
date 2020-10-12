// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	models "github.com/NodeFactoryIo/vedran/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// NodeRepository is an autogenerated mock type for the NodeRepository type
type NodeRepository struct {
	mock.Mock
}

// FindByID provides a mock function with given fields: ID
func (_m *NodeRepository) FindByID(ID string) (*models.Node, error) {
	ret := _m.Called(ID)

	var r0 *models.Node
	if rf, ok := ret.Get(0).(func(string) *models.Node); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *NodeRepository) GetAll() (*[]models.Node, error) {
	ret := _m.Called()

	var r0 *[]models.Node
	if rf, ok := ret.Get(0).(func() *[]models.Node); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.Node)
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

// IsNodeWhitelisted provides a mock function with given fields: ID
func (_m *NodeRepository) IsNodeWhitelisted(ID string) (bool, error) {
	ret := _m.Called(ID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(ID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: node
func (_m *NodeRepository) Save(node *models.Node) error {
	ret := _m.Called(node)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Node) error); ok {
		r0 = rf(node)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
