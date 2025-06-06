package model

import "github.com/gofrs/uuid"

type WeatherForecast struct {
	HeaderID   uuid.UUID
	Station_ID string
	Icao       string
	Latitude   float64
	Longitude  float64
}
type RealtimeWeatherMapperToDB struct {
	ID                      uuid.UUID `json:"id"`
	Station_ID              uuid.UUID `json:"weather_station_id"`
	Icao                    string    `json:"icao"`
	TempSurface             float64   `json:"temp-surface"`
	Temp1000H               float64   `json:"temp-1000h"`
	Temp800H                float64   `json:"temp-800h"`
	Temp400H                float64   `json:"temp-400h"`
	Temp200H                float64   `json:"temp-200h"`
	DewpointSurface         float64   `json:"dewpoint-surface"`
	Dewpoint1000H           float64   `json:"dewpoint-1000h"`
	Dewpoint800H            float64   `json:"dewpoint-800h"`
	Dewpoint400H            float64   `json:"dewpoint-400h"`
	Dewpoint200H            float64   `json:"dewpoint-200h"`
	Past3HprecipSurface     float64   `json:"past3hprecip-surface"`
	Past3HconvprecipSurface float64   `json:"past3hconvprecip-surface"`
	Past3HsnowprecipSurface float64   `json:"past3hsnowprecip-surface"`
	WindUSurface            float64   `json:"wind_u-surface"`
	WindU1000H              float64   `json:"wind_u-1000h"`
	WindU800H               float64   `json:"wind_u-800h"`
	WindU400H               float64   `json:"wind_u-400h"`
	WindU200H               float64   `json:"wind_u-200h"`
	WindVSurface            float64   `json:"wind_v-surface"`
	WindV1000H              float64   `json:"wind_v-1000h"`
	WindV800H               float64   `json:"wind_v-800h"`
	WindV400H               float64   `json:"wind_v-400h"`
	WindV200H               float64   `json:"wind_v-200h"`
	GustSurface             float64   `json:"gust-surface"`
	CapeSurface             float64   `json:"cape-surface"`
	PtypeSurface            int       `json:"ptype-surface"`
	LcloudsSurface          float64   `json:"lclouds-surface"`
	McloudsSurface          float64   `json:"mclouds-surface"`
	HcloudsSurface          float64   `json:"hclouds-surface"`
	RhSurface               float64   `json:"rh-surface"`
	Rh1000H                 float64   `json:"rh-1000h"`
	Rh800H                  float64   `json:"rh-800h"`
	Rh400H                  float64   `json:"rh-400h"`
	Rh200H                  float64   `json:"rh-200h"`
	GhSurface               float64   `json:"gh-surface"`
	Gh1000H                 float64   `json:"gh-1000h"`
	Gh800H                  float64   `json:"gh-800h"`
	Gh400H                  float64   `json:"gh-400h"`
	Gh200H                  float64   `json:"gh-200h"`
	PressureSurface         float64   `json:"pressure-surface"`
	Created_By              string    `json:"created_by"`
}

type Windy_Realtime_Report struct {
	Ts    []int64 `json:"ts"`
	Units struct {
		TempSurface             string `json:"temp-surface"`
		Temp1000H               string `json:"temp-1000h"`
		Temp800H                string `json:"temp-800h"`
		Temp400H                string `json:"temp-400h"`
		Temp200H                string `json:"temp-200h"`
		DewpointSurface         string `json:"dewpoint-surface"`
		Dewpoint1000H           string `json:"dewpoint-1000h"`
		Dewpoint800H            string `json:"dewpoint-800h"`
		Dewpoint400H            string `json:"dewpoint-400h"`
		Dewpoint200H            string `json:"dewpoint-200h"`
		Past3HprecipSurface     string `json:"past3hprecip-surface"`
		Past3HconvprecipSurface string `json:"past3hconvprecip-surface"`
		Past3HsnowprecipSurface string `json:"past3hsnowprecip-surface"`
		WindUSurface            string `json:"wind_u-surface"`
		WindU1000H              string `json:"wind_u-1000h"`
		WindU800H               string `json:"wind_u-800h"`
		WindU400H               string `json:"wind_u-400h"`
		WindU200H               string `json:"wind_u-200h"`
		WindVSurface            string `json:"wind_v-surface"`
		WindV1000H              string `json:"wind_v-1000h"`
		WindV800H               string `json:"wind_v-800h"`
		WindV400H               string `json:"wind_v-400h"`
		WindV200H               string `json:"wind_v-200h"`
		GustSurface             string `json:"gust-surface"`
		CapeSurface             string `json:"cape-surface"`
		PtypeSurface            any    `json:"ptype-surface"`
		LcloudsSurface          string `json:"lclouds-surface"`
		McloudsSurface          string `json:"mclouds-surface"`
		HcloudsSurface          string `json:"hclouds-surface"`
		RhSurface               string `json:"rh-surface"`
		Rh1000H                 string `json:"rh-1000h"`
		Rh800H                  string `json:"rh-800h"`
		Rh400H                  string `json:"rh-400h"`
		Rh200H                  string `json:"rh-200h"`
		GhSurface               string `json:"gh-surface"`
		Gh1000H                 string `json:"gh-1000h"`
		Gh800H                  string `json:"gh-800h"`
		Gh400H                  string `json:"gh-400h"`
		Gh200H                  string `json:"gh-200h"`
		PressureSurface         string `json:"pressure-surface"`
	} `json:"units"`
	TempSurface             []float64 `json:"temp-surface"`
	Temp1000H               []float64 `json:"temp-1000h"`
	Temp800H                []float64 `json:"temp-800h"`
	Temp400H                []float64 `json:"temp-400h"`
	Temp200H                []float64 `json:"temp-200h"`
	DewpointSurface         []float64 `json:"dewpoint-surface"`
	Dewpoint1000H           []float64 `json:"dewpoint-1000h"`
	Dewpoint800H            []float64 `json:"dewpoint-800h"`
	Dewpoint400H            []float64 `json:"dewpoint-400h"`
	Dewpoint200H            []float64 `json:"dewpoint-200h"`
	Past3HprecipSurface     []float64 `json:"past3hprecip-surface"`
	Past3HconvprecipSurface []float64 `json:"past3hconvprecip-surface"`
	Past3HsnowprecipSurface []float64 `json:"past3hsnowprecip-surface"`
	WindUSurface            []float64 `json:"wind_u-surface"`
	WindU1000H              []float64 `json:"wind_u-1000h"`
	WindU800H               []float64 `json:"wind_u-800h"`
	WindU400H               []float64 `json:"wind_u-400h"`
	WindU200H               []float64 `json:"wind_u-200h"`
	WindVSurface            []float64 `json:"wind_v-surface"`
	WindV1000H              []float64 `json:"wind_v-1000h"`
	WindV800H               []float64 `json:"wind_v-800h"`
	WindV400H               []float64 `json:"wind_v-400h"`
	WindV200H               []float64 `json:"wind_v-200h"`
	GustSurface             []float64 `json:"gust-surface"`
	CapeSurface             []float64 `json:"cape-surface"`
	PtypeSurface            []int     `json:"ptype-surface"`
	LcloudsSurface          []float64 `json:"lclouds-surface"`
	McloudsSurface          []float64 `json:"mclouds-surface"`
	HcloudsSurface          []float64 `json:"hclouds-surface"`
	RhSurface               []float64 `json:"rh-surface"`
	Rh1000H                 []float64 `json:"rh-1000h"`
	Rh800H                  []float64 `json:"rh-800h"`
	Rh400H                  []float64 `json:"rh-400h"`
	Rh200H                  []float64 `json:"rh-200h"`
	GhSurface               []float64 `json:"gh-surface"`
	Gh1000H                 []float64 `json:"gh-1000h"`
	Gh800H                  []float64 `json:"gh-800h"`
	Gh400H                  []float64 `json:"gh-400h"`
	Gh200H                  []float64 `json:"gh-200h"`
	PressureSurface         []float64 `json:"pressure-surface"`
	Warning                 string    `json:"warning"`
}
