// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// NanoIDGenerator is an autogenerated mock type for the NanoIDGenerator type
type NanoIDGenerator struct {
	mock.Mock
}

// Generate provides a mock function with no fields
func (_m *NanoIDGenerator) Generate() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Generate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewNanoIDGenerator creates a new instance of NanoIDGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNanoIDGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *NanoIDGenerator {
	mock := &NanoIDGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
