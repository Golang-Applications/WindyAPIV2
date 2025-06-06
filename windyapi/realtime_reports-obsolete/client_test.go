package main

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSomething(t *testing.T) {

	windyAPITestObj := new(MockRealtimeClient)
	windyAPITestObj.On("getWindyRealtimeWeatherDtlsAsync", 1, mock.Anything, mock.Anything).Return()
	windyAPITestObj.getWindyRealtimeWeatherDtlsAsync(1, nil, nil)
	windyAPITestObj.AssertExpectations(t)

}

func invokeMock(windyAPITestObj *MockRealtimeClient, id int, jobs <-chan Job, results chan<- Result) {
	windyAPITestObj.getWindyRealtimeWeatherDtlsAsync(id, jobs, results)
}
