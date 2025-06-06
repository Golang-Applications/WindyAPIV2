package api

import (
	"Windy-API/models"
	"Windy-API/svc"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type windyClient struct {
	apiEndpoint string
}

type Client interface {
	postRequestToWindyAPI() map[string]Result
}

func (c *windyClient) postRequestToWindyAPI() map[string]Result {

	if c.apiEndpoint == "" {
		log.Fatal("Empty windy API endpoint passed in the request")
	}
	if !isValidURL(c.apiEndpoint) {
		log.Fatal("Invalid Windy API endpoint passed in the request")
	}
	var wg sync.WaitGroup
	results := make(map[string]Result)
	var mu sync.Mutex // Mutex to protect the shared map

	for i, _ := range windyAPIInput {
		wg.Add(1)
		go func(icao string) {
			defer wg.Done()
			client := http.Client{}
			req, _ := http.NewRequest("POST", c.apiEndpoint, strings.NewReader(windyAPIInput[i].BodyRequest))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				svc.Logger().Error().Msg("Error while retrieving data for icao: " + windyAPIInput[i].Icao + " Error: " + err.Error())
				log.Println("Error while retrieving data for icao: ", windyAPIInput[i].Icao, " Error:", err.Error())
				return
			}
			if res.StatusCode != http.StatusOK {
				log.Println("Data not available for station: ", windyAPIInput[i].Icao)
				return
			}
			windyAPIResponse, err := convertWindyAPIResponseToJson(res.Body)
			if err != nil {
				log.Println("Error while saving data for icao: ", windyAPIInput[i].Icao, " Error:", err.Error())
				return
			}
			var resultBody model.Windy_Realtime_Report
			err = json.Unmarshal([]byte(windyAPIResponse), &resultBody)
			if err != nil {
				log.Println("Unable to unmarshal json for Icao code:", windyAPIInput[i].Icao, " Error :", err.Error())
				return
			}
			mu.Lock()
			rslt := Result{
				Icao:              windyAPIInput[i].Icao,
				HeaderID:          windyAPIInput[i].HeaderID,
				Output:            windyAPIResponse,
				WindyRealtimeRslt: resultBody,
			}
			results[windyAPIInput[i].Icao] = rslt
			mu.Unlock()
		}(windyAPIInput[i].Icao)
	}
	wg.Wait()
	return results
}

func isValidURL(endPointURL string) bool {
	_, err := url.ParseRequestURI(endPointURL)
	if err != nil {
		return false
	}
	u, err := url.Parse(endPointURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
