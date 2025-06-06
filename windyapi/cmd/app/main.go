package main

import (
	"Windy-API/api"
	"Windy-API/config"
	"Windy-API/persistence"
	"Windy-API/schema"
	"Windy-API/svc"
	"log"
	"os"
	"runtime"
)

func main() {
	svc.Logger().Info().Msg("Initiating Windy API Realtime Reports Application")

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

	db, err := schema.OpenDatabase(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = db.Close()
	}()

	validateConfigProperties(cfg)
	fileDirectories := []string{cfg.Process.RealTimeDirectory}

	if !validateDirectories(fileDirectories) {
		log.Fatal("Unable to proceed further as one of the file directories or file names do not exist")
	}
	p := persistence.Persistence{
		Cfg: cfg.Database,
		Db:  db,
	}

	a, err := api.NewApi(*cfg, p)
	if err != nil {
		log.Fatalf("failed to initialize api layer - %s", err.Error())
	}

	//Check if the weather_stations table exists, if it does not exist, then do not proceed further
	count, err := p.CheckForTableExistence(cfg.Database.Name, "weather_stations")
	if err == nil && count == 1 {
		countRow, innerErr := p.CheckRowExistence("weather_stations") //Make sure weather_stations is already populated with data
		if innerErr == nil && countRow == 0 {
			svc.Logger().Info().Msg("The weather_stations table is empty. Please populate the table with data before proceeding further")
			log.Fatal("The weather_stations table is empty. Please populate the table with data before proceeding further")
		}

	} else {
		svc.Logger().Info().Msg("The weather_stations table does not exist. Please create the table and populate the data before proceeding further")
		//log.Fatal("The weather_stations table does not exist. Please create the table and populate the data before proceeding further")
	}
	a.ProcessRequest()
}

func validateDirectories(fileDirectories []string) bool {
	for _, fileDirectory := range fileDirectories {
		if _, err := os.Stat(fileDirectory); os.IsNotExist(err) {
			println("File directory or filename does not exist: ", fileDirectory)
			return false
		}
	}
	return true
}

func validateConfigProperties(cfg *config.Config) {
	if cfg.Process.MaxRecordsToProcess <= 0 {
		cfg.Process.MaxRecordsToProcess = 100
	}
	if cfg.Process.BatchCount <= 0 {
		cfg.Process.BatchCount = 25000
	}
	if cfg.Process.ConcurrentRequests <= 0 {
		cfg.Process.ConcurrentRequests = 10
	}
	if cfg.Process.MaxWorkerPools <= 0 {
		cfg.Process.MaxWorkerPools = runtime.NumCPU()
	}

}
