// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"

// Certificate is an autogenerated mock type for the Certificate type
type Certificate struct {
	mock.Mock
}

// Approve provides a mock function with given fields: _a0, _a1
func (_m *Certificate) Approve(_a0 string, _a1 int64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Request provides a mock function with given fields: _a0, _a1
func (_m *Certificate) Request(_a0 string, _a1 []string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
