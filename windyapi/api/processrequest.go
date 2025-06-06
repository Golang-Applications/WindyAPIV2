package api

import (
	_ "Windy-API/config"
	"Windy-API/models"
	"Windy-API/svc"
	"encoding/json"
	"fmt"
	"github.com/bradhe/stopwatch"
	_ "github.com/fatih/color"
	"github.com/gofrs/uuid"
	_ "github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

type Job struct {
	ID      int
	Request RealtimeParameters
}

type Result struct {
	ID                int
	Output            string
	Icao              string
	HeaderID          uuid.UUID
	WindyRealtimeRslt model.Windy_Realtime_Report
}

var windyAPIInput []RealtimeParameters

func (a *api) processRequests() {
	client := windyClient{
		apiEndpoint: a.globalCfg.WindyAPI.Endpoint,
	}

	start := startWatch()
	weatherResults := a.getWeatherStations(a.globalCfg.Process.MaxRecordsToProcess)

	a.populateWindyAPIInput(weatherResults)
	response := client.postRequestToWindyAPI()

	a.processWindyResponse(response)
	elapsedSeconds := stopWatch(start)
	svc.Logger().Info().Msg(fmt.Sprintf("Total time taken to process %d records: %d second(s)", len(windyAPIInput), len(elapsedSeconds.String())-2))
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

func (a *api) saveResponseAsJson(icao string, windyAPIResponse string) error {
	fileNameWithPath := filepath.Join(a.globalCfg.Process.RealTimeDirectory, fmt.Sprintf("%s-%s.json", strconv.FormatInt(time.Now().Unix(), 10), icao))
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

func (a *api) buildAPIRealtimeRequest(lat float64, lng float64) string {
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
	mapRequest["model"] = a.globalCfg.Request.Model
	mapRequest["parameters"] = a.globalCfg.Request.Parameters
	mapRequest["levels"] = a.globalCfg.Request.Levels
	mapRequest["key"] = a.globalCfg.Request.ApiKey
	jsonRequest, _ := json.Marshal(mapRequest)
	return string(jsonRequest)
}

func (a *api) populateWindyAPIInput(weatherresults []model.WeatherForecast) {
	i := 0
	windyAPIInput = make([]RealtimeParameters, a.globalCfg.Process.MaxRecordsToProcess)
	for _, weatherHdr := range weatherresults {
		tmpInput := RealtimeParameters{
			Latitude:    weatherHdr.Latitude,
			Longitude:   weatherHdr.Longitude,
			Icao:        weatherHdr.Icao,
			Station_ID:  weatherHdr.Station_ID,
			HeaderID:    weatherHdr.HeaderID,
			BodyRequest: a.buildAPIRealtimeRequest(weatherHdr.Latitude, weatherHdr.Longitude),
		}
		windyAPIInput[i] = tmpInput
		i++
	}
}

func (a *api) getWeatherStations(maxRecords int) []model.WeatherForecast {
	weatherResults, err := a.dal.GetWeatherStations(maxRecords)
	if err != nil {
		svc.Logger().Error().Msg("Error while retrieving weather stations: " + err.Error())
		log.Fatalf("Error while retrieving weather stations: %s", err)
	}
	if weatherResults == nil {
		svc.Logger().Error().Msg("No Weather stations found in the database")
		log.Fatal("Stopping the process. No Weather stations found in the database")
	}
	svc.Logger().Info().Msg("Total records retrieved from the weather_stations table: " + strconv.Itoa(len(weatherResults)) + " records")
	return weatherResults
}

func startWatch() stopwatch.Watch {
	return stopwatch.Start()
}

func stopWatch(stopWatch stopwatch.Watch) time.Duration {
	stopWatch.Stop()
	return stopWatch.Seconds()
}

func (a *api) processWindyResponse(result map[string]Result) {
	for icao, _ := range result {
		pmtsValues, pmtsStrings := buildValues(result[icao].WindyRealtimeRslt, result[icao].Icao, result[icao].HeaderID)
		err := a.dal.InsertBatchStatements(buildWindyDetailSQL(), pmtsStrings, pmtsValues)
		if err != nil {
			fmt.Println("Error inserting batch statements: ", err)
			continue
		}
	}

}

func newUuid() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}
