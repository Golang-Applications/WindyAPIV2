package main

import (
	"Windy-API/config"
	"Windy-API/realtime_reports-obsolete/pkg/model"
	Repository "Windy-API/realtime_reports-obsolete/pkg/repository"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"time"
)

var app application

const (
	weatherStations  = "weather_stations"
	dbName           = "mmro"
	validDirectory   = "/Users/anandkumar/GolandProjects/awesomeProject/"
	invalidDirectory = "/Users/anandkumar/GolandProjects/awesomeProject/invalid_directory"
)

var resultWeatherStations *[]model.WeatherForecast

func TestMain(m *testing.M) {
	println("Inside the TestMain")
	err := godotenv.Load("./../local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Database.Timeout = 120
	conn, err := app.openDatabase(cfg.Database)
	if err != nil {
		log.Fatal("Error opening database connection: ", err.Error(), "")
	}
	app.DB = &Repository.RealtimeForecastDBRepo{conn, time.Second * time.Duration(120)}
	defer func() {
		_ = conn.Close()
	}()
	app.DSN = app.DsnConn()
	app.maxWorkerPool = 1
	app.maxRecords = 2
	app.batchCount = 2
	app.windyRequestModel = cfg.Request.Model
	app.windyRequestParams = cfg.Request.Parameters
	app.windyRequestLevels = cfg.Request.Levels
	app.windyAPIKey = cfg.Request.ApiKey

	os.Exit(m.Run())
}
