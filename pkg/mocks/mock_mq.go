// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hitesh22rana/mq/pkg/mq (interfaces: MQ)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mq "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// MockMQ is a mock of MQ interface.
type MockMQ struct {
	ctrl     *gomock.Controller
	recorder *MockMQMockRecorder
}

// MockMQMockRecorder is the mock recorder for MockMQ.
type MockMQMockRecorder struct {
	mock *MockMQ
}

// NewMockMQ creates a new mock instance.
func NewMockMQ(ctrl *gomock.Controller) *MockMQ {
	mock := &MockMQ{ctrl: ctrl}
	mock.recorder = &MockMQMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMQ) EXPECT() *MockMQMockRecorder {
	return m.recorder
}

// CreateChannel mocks base method.
func (m *MockMQ) CreateChannel(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChannel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateChannel indicates an expected call of CreateChannel.
func (mr *MockMQMockRecorder) CreateChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChannel", reflect.TypeOf((*MockMQ)(nil).CreateChannel), arg0, arg1)
}

// Publish mocks base method.
func (m *MockMQ) Publish(arg0 context.Context, arg1 string, arg2 *mq.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockMQMockRecorder) Publish(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockMQ)(nil).Publish), arg0, arg1, arg2)
}

// Subscribe mocks base method.
func (m *MockMQ) Subscribe(arg0 context.Context, arg1 *mq.Subscriber, arg2 mq.Offset, arg3 uint64, arg4 string, arg5 chan<- *mq.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockMQMockRecorder) Subscribe(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockMQ)(nil).Subscribe), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UnSubscribe mocks base method.
func (m *MockMQ) UnSubscribe(arg0 context.Context, arg1 *mq.Subscriber, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnSubscribe", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnSubscribe indicates an expected call of UnSubscribe.
func (mr *MockMQMockRecorder) UnSubscribe(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnSubscribe", reflect.TypeOf((*MockMQ)(nil).UnSubscribe), arg0, arg1, arg2)
}
