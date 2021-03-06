// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import models "github.com/NodeFactoryIo/vedran/internal/models"

import time "time"

// RecordRepository is an autogenerated mock type for the RecordRepository type
type RecordRepository struct {
	mock.Mock
}

// CountFailedRequests provides a mock function with given fields:
func (_m *RecordRepository) CountFailedRequests() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountSuccessfulRequests provides a mock function with given fields:
func (_m *RecordRepository) CountSuccessfulRequests() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSuccessfulRecordsInsideInterval provides a mock function with given fields: nodeID, from, to
func (_m *RecordRepository) FindSuccessfulRecordsInsideInterval(nodeID string, from time.Time, to time.Time) ([]models.Record, error) {
	ret := _m.Called(nodeID, from, to)

	var r0 []models.Record
	if rf, ok := ret.Get(0).(func(string, time.Time, time.Time) []models.Record); ok {
		r0 = rf(nodeID, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Record)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, time.Time, time.Time) error); ok {
		r1 = rf(nodeID, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: record
func (_m *RecordRepository) Save(record *models.Record) error {
	ret := _m.Called(record)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Record) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
