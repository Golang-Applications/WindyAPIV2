package config

import (
	"github.com/go-andiamo/cfgenv"
)

func Load() (*Config, error) {
	return cfgenv.LoadAs[Config]()
}

type Config struct {
	Version  string   `env:"optional,default=1.0.0" json:"-"`
	Build    string   `env:"optional,default=1.0.0" json:"-"`
	Database Database `env:"prefix=DATABASE" json:"database"`
	Process  Process  `env:"prefix=PROCESS" json:"process"`
	Request  Request  `env:"prefix=REQUEST" json:"requestParams"`
	WindyAPI WindyAPI `env:"prefix=WINDYAPI" json:"windyAPI"`
	Response Response `env:"prefix=RESPONSE" json:"response"`
}
type Database struct {
	Host     string `json:"host"`
	Port     int    `env:"optional,default=3306" json:"port"`
	Name     string `env:"optional,default=mmro" json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
	Timeout  int    `env:"optional,default=60" json:"timeout"`
}

type Process struct {
	RealTimeDirectory           string `env:"optional,default=/home/windyAPI" json:"realTimeDirectory"`
	MaxRecordsToProcess         int    `env:"optional,default=100" json:"maxRecordsToProcess"`
	HistoryRecords              bool   `env:"optional,default=true" json:"executeHistoryRecords"`
	BatchCount                  int    `env:"optional,default=30000" json:"batchCount"`
	HistoricalFileName          string `env:"optional,default=HistoricalWeatherData.json" json:"historicalFileName"`
	SaveHistoricalWeatherSQLLoc string `env:"optional,default=/home/windyAPI" json:"saveHistoricalWeatherSQLLoc"`
	ConcurrentRequests          int    `env:"optional,default=10" json:"concurrentRequests"`
	MaxWorkerPools              int    `env:"optional,default=1" json:"maxWorkerPools"`
}

type Request struct {
	Model      string   `json:"model"`
	Parameters []string `json:"parameters"`
	Levels     []string `json:"levels"`
	ApiKey     string   `json:"apiKey"`
}

type WindyAPI struct {
	Endpoint string `env:"optional,default=https://api.windy.com/api/v0.1/forecast" json:"-"`
}

type Response struct {
	SaveResponse bool `json:"saveFile"`
}
