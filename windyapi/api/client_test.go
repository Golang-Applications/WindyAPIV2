package api

import (
	"Windy-API/config"
	"Windy-API/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestPostReqToWindyAPI(t *testing.T) {

	windyAPIEndpoint := "https://api.windy.com/api/point-forecast/v2"

	// Define a HTTP handler.

	request := `{"key":"WbtEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	//invalidRequest := `{"key:"WbtEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	//emptyAPIKeyRequest := `{"key":"","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	//invalidAPIKeyRequest := `{"key":"btEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`

	var windyAPITests = []struct {
		id             int
		name           string
		requestBody    string
		endpoint       string
		expectedLength int
	}{
		{1, "Valid Endpoint", request, windyAPIEndpoint, 1},
		{2, "Invalid Endpoint", request, windyAPIEndpoint, 0},
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error - %s", err.Error())
	}
	client := windyClient{
		apiEndpoint: cfg.WindyAPI.Endpoint,
	}

	var rslt map[string]Result
	for _, tt := range windyAPITests {
		switch tt.id {
		case 1:
			//testReq := a.buildAPIRealtimeRequest(53.19, -112.25)
			forecast := populateForecast(-34.90, 117.80, "YPAL", "94802")
			populateWindyAPIInput(forecast, cfg)
			rslt = client.postRequestToWindyAPI()
			assert.Equal(t, tt.expectedLength, len(rslt))
			break
		case 2:
			forecast := populateForecast(-34.90, 117.80, "YPAL", "94802")
			populateWindyAPIInput(forecast, cfg)
			client.apiEndpoint = "https://api.windy.com/api/point-forecast"
			rslt = client.postRequestToWindyAPI()
			assert.Equal(t, tt.expectedLength, len(rslt))
			break
		default:
		}

	}
}

func populateWindyAPIInput(forecast model.WeatherForecast, cfg *config.Config) {
	windyAPIInput = append(windyAPIInput, RealtimeParameters{
		Latitude:    forecast.Latitude,
		Longitude:   forecast.Longitude,
		Icao:        forecast.Icao,
		Station_ID:  forecast.Station_ID,
		HeaderID:    forecast.HeaderID,
		BodyRequest: buildAPIRealtimeRequest(cfg, forecast.Latitude, forecast.Longitude),
	})
}

func populateForecast(lat float64, lng float64, icao string, stationID string) model.WeatherForecast {
	var forecast model.WeatherForecast
	forecast.Icao = icao
	forecast.Longitude = lng
	forecast.Latitude = lat
	forecast.Station_ID = stationID
	forecast.HeaderID = newUuid()
	return forecast
}

func buildAPIRealtimeRequest(cfg *config.Config, lat float64, lng float64) string {

	mapRequest := make(map[string]any)
	mapRequest["lat"] = lat
	mapRequest["lon"] = lng
	mapRequest["model"] = cfg.Request.Model
	mapRequest["parameters"] = cfg.Request.Parameters
	mapRequest["levels"] = cfg.Request.Levels
	mapRequest["key"] = cfg.Request.ApiKey
	jsonRequest, _ := json.Marshal(mapRequest)
	return string(jsonRequest)
}
