package main

type Client interface {
	getWindyRealtimeWeatherDtlsAsync(id int, jobs <-chan Job, results chan<- Result)
}
