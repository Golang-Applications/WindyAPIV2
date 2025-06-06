package main

import (
	"Windy-API/realtime_reports-obsolete/pkg/model"
	"fmt"
	"github.com/gofrs/uuid"
)

func (app *application) processResponse(result model.Windy_Realtime_Report, icao string, headerID uuid.UUID) error {
	pmtsValues, pmtsStrings := app.buildValues(result, icao, headerID)
	err := app.DB.InsertBatchStatements(app.buildWindyDetailSQL(), pmtsStrings, pmtsValues)
	if err != nil {
		fmt.Println("Error inserting batch statements: ", err)
		return err
	}
	return nil
}

func (app *application) buildWindyDetailSQL() string {
	return `INSERT INTO weather_realtime_reports(id,
                                     weather_station_id,
                                     icao,
                                     temp_surface,
                                     temp_1000h,
                                     temp_800h,
                                     temp_400h,
                                     temp_200h,
                                     dewpoint_surface,
                                     dewpoint_1000h,
                                     dewpoint_800h,
                                     dewpoint_400h,
                                     dewpoint_200h,
                                     past3hprecip_surface,
                                     past3hconvprecip_surface,
                                     past3hsnowprecip_surface,
                                     wind_u_surface,
                                     wind_u_1000h,
                                     wind_u_800h,
                                     wind_u_400h,
                                     wind_u_200h,
                                     wind_v_surface,
                                     wind_v_1000h,
                                     wind_v_800h,
                                     wind_v_400h,
                                     wind_v_200h,
                                     gust_surface,
                                     cape_surface,
                                     ptype_surface,
                                     lclouds_surface,
                                     mclouds_surface,
                                     hclouds_surface,
                                     rh_surface,
                                     rh_1000h,
                                     rh_800h,
                                     rh_400h,
                                     rh_200h,
                                     gh_surface,
                                     gh_1000h,
                                     gh_800h,
                                     gh_400h,
                                     gh_200h,
                                     pressure_surface,
                                     created_by) values `
}

func (app *application) buildParameters() string {
	return `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
}

func (app *application) buildValues(result model.Windy_Realtime_Report, icao string, headerID uuid.UUID) ([]model.RealtimeWeatherMapperToDB, string) {
	parameterizedValues := make([]model.RealtimeWeatherMapperToDB, 0, len(result.Ts))
	parameterizedStrings := app.buildParameters()
	for i := 0; i < len(result.Ts); i++ {
		mapperToDB := model.RealtimeWeatherMapperToDB{
			ID:                      newUuid(),
			Station_ID:              headerID,
			Icao:                    icao,
			TempSurface:             result.TempSurface[i],
			Temp1000H:               result.Temp1000H[i],
			Temp800H:                result.Temp800H[i],
			Temp400H:                result.Temp400H[i],
			Temp200H:                result.Temp200H[i],
			DewpointSurface:         result.DewpointSurface[i],
			Dewpoint1000H:           result.Dewpoint1000H[i],
			Dewpoint800H:            result.Dewpoint800H[i],
			Dewpoint400H:            result.Dewpoint400H[i],
			Dewpoint200H:            result.Dewpoint200H[i],
			WindUSurface:            result.WindUSurface[i],
			Past3HconvprecipSurface: result.Past3HconvprecipSurface[i],
			Past3HprecipSurface:     result.Past3HprecipSurface[i],
			Past3HsnowprecipSurface: result.Past3HsnowprecipSurface[i],
			WindU1000H:              result.WindU1000H[i],
			WindU800H:               result.WindU800H[i],
			WindU400H:               result.WindU400H[i],
			WindU200H:               result.WindU200H[i],
			WindVSurface:            result.WindVSurface[i],
			WindV1000H:              result.WindV1000H[i],
			WindV800H:               result.WindV800H[i],
			WindV400H:               result.WindV400H[i],
			WindV200H:               result.WindV200H[i],
			GustSurface:             result.GustSurface[i],
			CapeSurface:             result.CapeSurface[i],
			PtypeSurface:            result.PtypeSurface[i],
			LcloudsSurface:          result.LcloudsSurface[i],
			McloudsSurface:          result.McloudsSurface[i],
			HcloudsSurface:          result.HcloudsSurface[i],
			RhSurface:               result.RhSurface[i],
			Rh1000H:                 result.Rh1000H[i],
			Rh800H:                  result.Rh800H[i],
			Rh400H:                  result.Rh400H[i],
			Rh200H:                  result.Rh200H[i],
			GhSurface:               result.GhSurface[i],
			Gh1000H:                 result.Gh1000H[i],
			Gh800H:                  result.Gh800H[i],
			Gh400H:                  result.Gh400H[i],
			Gh200H:                  result.Gh200H[i],
			PressureSurface:         result.PressureSurface[i],
			Created_By:              "00000000-0000-4000-0000-000000000000",
		}
		parameterizedValues = append(parameterizedValues, mapperToDB)
	}
	return parameterizedValues, parameterizedStrings
}
