package main

import (
	"github.com/stretchr/testify/mock"
)

func New() *MockRealtimeClient {
	return &MockRealtimeClient{}
}

type MockRealtimeClient struct {
	mock.Mock
}

func (m *MockRealtimeClient) getWindyRealtimeWeatherDtlsAsync(id int, jobs <-chan Job, results chan<- Result) {
	m.Called(id, jobs, results)
}
