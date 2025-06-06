package main

import (
	"Windy-API/realtime_reports-obsolete/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
)

type Job struct {
	ID      int
	Request RealtimeParameters
}

type Result struct {
	ID       int
	Output   string
	Icao     string
	HeaderID uuid.UUID
}

var errorChan chan error

func (app *application) processworkerResponse(windyAPIResults chan Result, chunks int) {
	var windyAPIResponse string
	i := 0
	for i = 0; i < chunks; i++ {
		select {
		case result := <-windyAPIResults:
			windyAPIResponse = result.Output
			if app.saveWindyAPIResponse {
				saveJsonErr := app.saveResponseAsJson(result.Icao, windyAPIResponse)
				if saveJsonErr != nil {
					log.Println("Error while saving response as json: ", "Error:", saveJsonErr.Error())
					continue
				}
			}
			var resultBody model.Windy_Realtime_Report
			err := json.Unmarshal([]byte(windyAPIResponse), &resultBody)
			if err != nil {
				log.Println("Unable to unmarshal json for Icao code:", result.Icao, " Error :", err.Error())
				continue
			}
			_ = app.processResponse(resultBody, result.Icao, result.HeaderID)
		case err := <-errorChan:
			fmt.Println(err)
			continue
		}
	}
}
