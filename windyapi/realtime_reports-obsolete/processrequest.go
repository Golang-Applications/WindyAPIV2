package main

import (
	"Windy-API/realtime_reports-obsolete/pkg/model"
	"Windy-API/realtime_reports-obsolete/svc"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradhe/stopwatch"
	_ "github.com/fatih/color"
	"github.com/gofrs/uuid"
	_ "github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RealtimeParameters struct {
	Latitude    float64
	Longitude   float64
	Icao        string
	Station_ID  string
	HeaderID    uuid.UUID
	BodyRequest string
}

var windyAPIInput []RealtimeParameters
var wg sync.WaitGroup

func (app application) processrequests() {
	errorChan = make(chan error, 0)
	windyAPIJobs := make(chan Job, 0)
	windyAPIResults := make(chan Result, 0)

	defer close(windyAPIJobs)
	defer close(windyAPIResults)
	defer close(errorChan)

	start := startWatch()
	weatherResults := app.getWeatherStations(app.maxRecords)

	app.populateWindyAPIInput(weatherResults)

	println(len(windyAPIInput), " records to be processed")
	chunks, rowLength, colLength := app.chunkWindyAPIInput(windyAPIInput, app.concurrentRequests)
	windyAPIJobs = make(chan Job, rowLength*colLength)
	windyAPIResults = make(chan Result, rowLength*colLength)
	app.sendToWorkPool(windyAPIJobs, windyAPIResults)
	app.dispatchJobsFromWorkPool(chunks, windyAPIJobs)
	wg.Wait()
	app.processworkerResponse(windyAPIResults, len(chunks))
	elapsedSeconds := stopWatch(start)
	svc.Logger().Info().Msg(fmt.Sprintf("Total time taken to process %d records: %d second(s)", len(windyAPIInput), len(elapsedSeconds.String())-2))
}

// This method is used for testing Windy API without any worker pool.
// This method is also useful for unit testing/mocking Windy API.
func (app application) postRequestToWindyAPI(request string, windyAPIEndpoint string) (string, error, int) {
	client := http.Client{}
	req, _ := http.NewRequest("POST", windyAPIEndpoint, strings.NewReader(request))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		println(err.Error())
		return "", err, res.StatusCode
	}
	reader, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err, http.StatusBadRequest
	}
	defer func() {
		_ = res.Body.Close()
	}()
	return string(reader), nil, res.StatusCode
}
func (app application) getWindyRealtimeWeatherDtlsAsync(id int, jobs <-chan Job, results chan<- Result) {
	i := -1

	for job := range jobs {
		i++
		client := http.Client{}
		req, _ := http.NewRequest("POST", app.windyAPIEndpoint, strings.NewReader(job.Request.BodyRequest))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			svc.Logger().Error().Msg("Error while retrieving data for icao: " + windyAPIInput[i].Icao + " Error: " + err.Error())
			log.Println("Error while retrieving data for icao: ", windyAPIInput[i].Icao, " Error:", err.Error())
			errorChan <- errors.New(fmt.Sprint("Error while retrieving data for icao: ", windyAPIInput[i].Icao, err.Error()))
			continue
		}
		if res.StatusCode != http.StatusOK {
			log.Println("Data not available for station: ", windyAPIInput[i].Icao)
			errorChan <- errors.New(fmt.Sprint("Data not available for icao: ", windyAPIInput[i].Icao))
			continue
		}
		windyAPIResponse, err := convertWindyAPIResponseToJson(res.Body)
		if err != nil {
			log.Println("Error while saving data for icao: ", windyAPIInput[i].Icao, " Error:", err.Error())
			errorChan <- errors.New(fmt.Sprint("Error while saving data for Icao: ", windyAPIInput[i].Icao, " Error:", err.Error()))
			continue
		}
		results <- Result{
			ID:       id,
			Output:   windyAPIResponse,
			Icao:     job.Request.Icao,
			HeaderID: job.Request.HeaderID,
		}

	}

}

func convertWindyAPIResponseToJson(body io.ReadCloser) (string, error) {
	reader, err := io.ReadAll(body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer func() {
		_ = body.Close()
	}()
	return string(reader), nil
}

func (app application) saveResponseAsJson(icao string, windyAPIResponse string) error {
	fileNameWithPath := filepath.Join(app.realTimeDirectory, fmt.Sprintf("%s-%s.json", strconv.FormatInt(time.Now().Unix(), 10), icao))
	f, err := os.OpenFile(fileNameWithPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		println(err.Error())
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, _ = f.WriteString(windyAPIResponse)
	return nil
}

func newUuid() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func (app application) buildAPIRealtimeRequest(lat float64, lng float64) string {
	//latitude float64, longitude float64
	/*
		{
		    "lat": 53.1900,
		    "lon": -112.2500,
		    "model": "gfs",
		    "parameters": ["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"],
		    "levels": ["surface", "1000h", "800h","400h","200h"],
		    "key": "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv"
		}
	*/

	mapRequest := make(map[string]any)
	mapRequest["lat"] = lat
	mapRequest["lon"] = lng
	mapRequest["model"] = app.windyRequestModel
	mapRequest["parameters"] = app.windyRequestParams
	mapRequest["levels"] = app.windyRequestLevels
	mapRequest["key"] = app.windyAPIKey
	jsonRequest, _ := json.Marshal(mapRequest)
	return string(jsonRequest)
}

func (app application) populateWindyAPIInput(weatherresults *[]model.WeatherForecast) {
	i := 0
	windyAPIInput = make([]RealtimeParameters, app.maxRecords)
	for _, weatherHdr := range *weatherresults {
		tmpInput := RealtimeParameters{
			Latitude:    weatherHdr.Latitude,
			Longitude:   weatherHdr.Longitude,
			Icao:        weatherHdr.Icao,
			Station_ID:  weatherHdr.Station_ID,
			HeaderID:    weatherHdr.HeaderID,
			BodyRequest: app.buildAPIRealtimeRequest(weatherHdr.Latitude, weatherHdr.Longitude),
		}
		windyAPIInput[i] = tmpInput
		i++
	}
}

func (app application) sendToWorkPool(windyAPIJobs <-chan Job, windyAPIResults chan<- Result) {
	for i := 0; i < app.maxWorkerPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			go app.getWindyRealtimeWeatherDtlsAsync(i, windyAPIJobs, windyAPIResults)
		}()
	}
}

func (app application) dispatchJobsFromWorkPool(chunks [][]RealtimeParameters, windyAPIJobs chan Job) {
	count := 0
	println("In the dispatchJobsFromWorkPool function:")
	for i := 0; i < len(chunks); i++ {
		for j, _ := range chunks[i] {
			windyAPIJobs <- Job{
				ID:      count,
				Request: chunks[i][j],
			}
			println(chunks[i][j].Icao)
			count++
		}
	}
	//println("Total jobs sent to work pool: ", count)
}

func (app application) getWeatherStations(maxRecords int) *[]model.WeatherForecast {
	weatherResults, err := app.DB.GetWeatherStations(maxRecords)
	if err != nil {
		svc.Logger().Error().Msg("Error while retrieving weather stations: " + err.Error())
		log.Fatalf("Error while retrieving weather stations: %s", err.Error())
	}
	if weatherResults == nil {
		svc.Logger().Error().Msg("No Weather stations found in the database")
		log.Fatal("Stopping the process. No Weather stations found in the database")
	}
	svc.Logger().Info().Msg("Total records retrieved from the weather_stations table: " + strconv.Itoa(len(*weatherResults)) + " records")
	return weatherResults
}

func (app application) chunkWindyAPIInput(windyAPIInput []RealtimeParameters, chunkSize int) ([][]RealtimeParameters, int, int) {
	var chunks [][]RealtimeParameters
	var colLength int
	for i := 0; i < len(windyAPIInput); i += chunkSize {
		end := i + chunkSize
		if end > len(windyAPIInput) {
			end = len(windyAPIInput)
		}
		chunks = append(chunks, windyAPIInput[i:end])
		colLength++
	}
	return chunks, len(windyAPIInput), colLength
}

func startWatch() stopwatch.Watch {
	return stopwatch.Start()
}

func stopWatch(stopWatch stopwatch.Watch) time.Duration {
	stopWatch.Stop()
	return stopWatch.Seconds()
}
