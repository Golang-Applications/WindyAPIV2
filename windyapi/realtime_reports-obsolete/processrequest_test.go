package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSaveResponseAsJson(t *testing.T) {
	println("Inside the TestsaveResponseAsJson")
	windyAPIResponse := `{"ts":[1745506800000],"units":{"temp-surface":"K"},
	"temp-surface":[283.8457736850904],"temp-1000h": [291.624081511577],
	"temp-800h":[286.6305006494195],
	"temp-400h":[247.5879339877249],
	"temp-200h":[211.26076016257068],
	"dewpoint-surface": [281.266027704108],
	"dewpoint-1000h": [284.440791633449],
	"dewpoint-800h": [282.4366003987724],
	"dewpoint-400h": [247.5879339877249],
	"dewpoint-200h": [211.26076016257068],
	"past3hconvprecip-surface": [0],
	"past3hsnowprecip-surface": [0],
	"wind_u-surface": [0.814398444650396],
	"wind_u-1000h": [0.0870500308177414],
	"wind_u-800h": [0.0870500308177414],
	"wind_u-400h": [0.0870500308177414],
	"wind_u-200h": [0.0870500308177414],
	"wind_v-surface": [-0.500297031639831],
	"wind_v-1000h": [-0.500297031639831],
	"wind_v-800h": [-0.500297031639831],
	"wind_v-400h": [-0.500297031639831],
	"wind_v-200h": [-0.500297031639831],
	"gust-surface": [3.65900147271776],
	"cape-surface": [0, 202.508285775323],
	"ptype-surface": [0],
	"lclouds-surface": [40.1249427777946],
	"mclouds-surface": [0],
	"hclouds-surface": [0],
	"rh-surface": [92.985722952971],
	"rh-1000h": [92.985722952971],
	"rh-800h": [92.985722952971],
	"rh-400h": [92.985722952971],
	"rh-200h": [92.985722952971],
	"gh-surface": [null],
	"gh-1000h": [114.445183248037],
	"gh-800h": [114.445183248037],
	"gh-400h": [114.445183248037],
	"gh-200h": [114.445183248037],
	"pressure-surface": [101536.810383358]
}`
	err := app.saveResponseAsJson("ABCD", windyAPIResponse)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, true)
}

func TestBuildAPIRealtimeRequest(t *testing.T) {
	request :=
		`{"key":"WbtEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	windAPIReq := app.buildAPIRealtimeRequest(53.1900, -112.2500)
	assert.Equal(t, request, windAPIReq)
}

func TestGetWeatherStations(t *testing.T) {
	println("Inside the TestGetWeatherStations")
	rslt := app.getWeatherStations(app.maxRecords)
	resultWeatherStations = rslt
	assert.True(t, len(*resultWeatherStations) == 2)
}

func TestPopulateWindyAPIInput(t *testing.T) {
	if len(*resultWeatherStations) == 0 {
		assert.False(t, len(*resultWeatherStations) == 0)
	}
	app.populateWindyAPIInput(resultWeatherStations)
	assert.True(t, len(windyAPIInput) == 2)

}

func TestUUID(t *testing.T) {
	uuidV4 := newUuid()
	assert.True(t, len(uuidV4.String()) == 36)
}

func TestPostRequestToWindyAPI(t *testing.T) {

	windyAPIEndpoint := "https://api.windy.com/api/point-forecast/v2"

	// Define a HTTP handler.

	request := `{"key":"WbtEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	invalidRequest := `{"key:"WbtEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	emptyAPIKeyRequest := `{"key":"","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	invalidAPIKeyRequest := `{"key":"btEx7uGxTdmvdCEEmKTu6ahDVMAJMao","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`

	var windyAPITests = []struct {
		id          int
		name        string
		requestBody string
		endpoint    string
		expected    int
	}{
		{1, "Valid Endpoint", request, windyAPIEndpoint, http.StatusOK},
		{2, "Invalid Json", invalidRequest, windyAPIEndpoint, http.StatusBadRequest},
		{3, "Empty API Key", emptyAPIKeyRequest, windyAPIEndpoint, http.StatusBadRequest},
		{4, "Invalid API Key", invalidAPIKeyRequest, windyAPIEndpoint, http.StatusBadRequest},
		{5, "Invalid EndPoint", request, request, http.StatusNotFound},
	}
	var httpCode int
	for _, tt := range windyAPITests {
		switch tt.id {
		case 1:
			_, _, httpCode = app.postRequestToWindyAPI(request, tt.endpoint)
			break
		case 2:
			_, _, httpCode = app.postRequestToWindyAPI(invalidRequest, tt.endpoint)
			break
		case 3:
			_, _, httpCode = app.postRequestToWindyAPI(emptyAPIKeyRequest, tt.endpoint)
			break
		case 4:
			_, _, httpCode = app.postRequestToWindyAPI(invalidAPIKeyRequest, tt.endpoint)
			break
		case 5:
			invalidEndPpint := "https://api.windy.com/api/pointforecast/v2"
			_, _, httpCode = app.postRequestToWindyAPI(request, invalidEndPpint)
			break
		default:
		}

		if httpCode != tt.expected {
			t.Errorf("got %d, want %d", httpCode, tt.expected)
		}

	}

}
