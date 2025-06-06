package main

import (
	"Windy-API/config"
	Repository "Windy-API/realtime_reports-obsolete/pkg/repository"
	"Windy-API/realtime_reports-obsolete/svc"
	"log"
	"os"
	"runtime"
	"time"
)

type application struct {
	DSN                  string
	host                 string
	port                 int
	user                 string
	password             string
	dbName               string
	loadIntoDB           bool
	realTimeDirectory    string
	DB                   Repository.DatabaseRepo
	maxRecords           int
	batchCount           int
	windyRequestModel    string
	windyRequestParams   []string
	windyRequestLevels   []string
	windyAPIKey          string
	windyAPIEndpoint     string
	saveWindyAPIResponse bool
	concurrentRequests   int
	maxWorkerPool        int
}

func main() {
	svc.Logger().Info().Msg("Initiating Windy API Realtime Reports Application")
	var app application

	svc.Logger().Info().Msg("Loading configuration")
	cfg, err := config.Load()
	if err != nil {
		svc.Logger().Err(err).Msg("Config load error")
		log.Fatalf("config load error - %s", err.Error())
	}
	svc.Init(*cfg)

	if cfg.Database.Timeout <= 0 {
		cfg.Database.Timeout = 120
	}
	conn, err := app.openDatabase(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = conn.Close()
	}()
	app.DB = &Repository.RealtimeForecastDBRepo{conn, time.Second * time.Duration(cfg.Database.Timeout)}
	app.initializeApplication(cfg)
	fileDirectories := []string{app.realTimeDirectory}

	if !validatefileDirectories(fileDirectories) {
		log.Fatal("Unable to proceed further as one of the file directories or file names do not exist")
	}
	//Check if the weather_stations table exists, if it does not exist, then do not proceed further
	count, err := app.DB.CheckForTableExistence(cfg.Database.Name, "weather_stations")
	if err == nil && count == 1 {
		countRow, innerErr := app.DB.CheckRowExistence() //Make sure weather_stations is already populated with data
		if innerErr == nil && countRow == 0 {
			svc.Logger().Info().Msg("The weather_stations table is empty. Please populate the table with data before proceeding further")
			log.Fatal("The weather_stations table is empty. Please populate the table with data before proceeding further")
		}

	} else {
		svc.Logger().Info().Msg("The weather_stations table does not exist. Please create the table and populate the data before proceeding further")
		log.Fatal("The weather_stations table does not exist. Please create the table and populate the data before proceeding further")
	}
	app.processrequests()
}

func validatefileDirectories(fileDirectories []string) bool {
	for _, fileDirectory := range fileDirectories {
		if _, err := os.Stat(fileDirectory); os.IsNotExist(err) {
			println("File directory or filename does not exist: ", fileDirectory)
			return false
		}
	}
	return true
}

func (app *application) initializeApplication(cfg *config.Config) {
	app.realTimeDirectory = cfg.Process.RealTimeDirectory
	app.maxRecords = cfg.Process.MaxRecordsToProcess
	app.windyRequestModel = cfg.Request.Model
	app.windyRequestParams = cfg.Request.Parameters
	app.windyRequestLevels = cfg.Request.Levels
	app.windyAPIKey = cfg.Request.ApiKey
	app.windyAPIEndpoint = cfg.WindyAPI.Endpoint
	app.saveWindyAPIResponse = cfg.Response.SaveResponse

	app.maxRecords = cfg.Process.MaxRecordsToProcess
	if app.maxRecords <= 0 {
		app.maxRecords = 100
	}
	app.batchCount = cfg.Process.BatchCount
	if app.batchCount <= 0 {
		app.batchCount = 25000
	}
	app.concurrentRequests = cfg.Process.ConcurrentRequests
	if app.concurrentRequests <= 0 {
		app.concurrentRequests = 10
	}
	app.maxWorkerPool = cfg.Process.MaxWorkerPools
	if app.maxWorkerPool <= 0 {
		app.maxWorkerPool = runtime.NumCPU()
	}

}
