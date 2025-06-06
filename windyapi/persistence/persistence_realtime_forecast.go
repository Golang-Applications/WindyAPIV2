package persistence

import (
	"Windy-API/models"
	"context"
	"fmt"
	"log"
	"time"
)

func (p *Persistence) GetWeatherStations(maxRecords int) ([]model.WeatherForecast, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Cfg.Timeout))
	defer cancel()

	var forecasts []model.WeatherForecast
	query := `select hdr.id, hdr.station_id, hdr.icao,hdr.latitude,hdr.longitude from weather_stations hdr LIMIT ?`

	rows, err := p.Db.QueryContext(ctx, query, maxRecords)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var forecast model.WeatherForecast
		err = rows.Scan(
			&forecast.HeaderID,
			&forecast.Station_ID,
			&forecast.Icao,
			&forecast.Latitude,
			&forecast.Longitude,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		forecasts = append(forecasts, forecast)
	}
	return forecasts, nil
}

func (p *Persistence) InsertHourlyWeatherReport(sql string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Cfg.Timeout))
	defer cancel()

	stmtInserts, err := p.Db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmtInserts.Close()
	}()
	_, err = stmtInserts.Exec()
	//_, err = repo.DB.ExecContext(ctx, sql)
	if err == nil {
		return nil
	}
	log.Println(err.Error())
	return err
}

func (p *Persistence) CheckRowExistence(tableName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Cfg.Timeout))
	defer cancel()
	var count int
	sqlStmt := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s LIMIT 1) as is_row_exists", tableName)
	err := p.Db.QueryRowContext(ctx, sqlStmt).Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *Persistence) CheckForTableExistence(dbName string, tableName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Cfg.Timeout))
	defer cancel()
	var count int
	sqlStmt := "select count(table_name) AS count from information_schema.tables where table_schema=? and table_name=?"
	err := p.Db.QueryRowContext(ctx, sqlStmt, dbName, tableName).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *Persistence) InsertBatchStatements(sql string, pmtrsQuery string, pmtrsValues []model.RealtimeWeatherMapperToDB) (batchError error) {
	var execErr error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Cfg.Timeout))
	defer cancel()
	tx, execErr := p.Db.BeginTx(ctx, nil)
	if execErr != nil {
		return execErr
	}
	defer func() {
		if execErr != nil {
			_ = tx.Rollback()
		} else {
			execErr = tx.Commit()
		}

	}()
	sql += fmt.Sprintf("%s", pmtrsQuery)
	prepareStmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer func() {
		_ = prepareStmt.Close()
	}()

	for i := 0; i < len(pmtrsValues); i++ {
		_, execErr = prepareStmt.ExecContext(ctx, pmtrsValues[i].ID,
			pmtrsValues[i].Station_ID,
			pmtrsValues[i].Icao,
			pmtrsValues[i].TempSurface,
			pmtrsValues[i].Temp1000H,
			pmtrsValues[i].Temp800H,
			pmtrsValues[i].Temp400H,
			pmtrsValues[i].Temp200H,
			pmtrsValues[i].DewpointSurface,
			pmtrsValues[i].Dewpoint1000H,
			pmtrsValues[i].Dewpoint800H,
			pmtrsValues[i].Dewpoint400H,
			pmtrsValues[i].Dewpoint200H,
			pmtrsValues[i].Past3HprecipSurface,
			pmtrsValues[i].Past3HconvprecipSurface,
			pmtrsValues[i].Past3HsnowprecipSurface,
			pmtrsValues[i].WindUSurface,
			pmtrsValues[i].WindU1000H,
			pmtrsValues[i].WindU800H,
			pmtrsValues[i].WindU400H,
			pmtrsValues[i].WindU200H,
			pmtrsValues[i].WindVSurface,
			pmtrsValues[i].WindV1000H,
			pmtrsValues[i].WindV800H,
			pmtrsValues[i].WindV400H,
			pmtrsValues[i].WindV200H,
			pmtrsValues[i].GustSurface,
			pmtrsValues[i].CapeSurface,
			pmtrsValues[i].PtypeSurface,
			pmtrsValues[i].LcloudsSurface,
			pmtrsValues[i].McloudsSurface,
			pmtrsValues[i].HcloudsSurface,
			pmtrsValues[i].RhSurface,
			pmtrsValues[i].Rh1000H,
			pmtrsValues[i].Rh800H,
			pmtrsValues[i].Rh400H,
			pmtrsValues[i].Rh200H,
			pmtrsValues[i].GhSurface,
			pmtrsValues[i].Gh1000H,
			pmtrsValues[i].Gh800H,
			pmtrsValues[i].Gh400H,
			pmtrsValues[i].Gh200H,
			pmtrsValues[i].PressureSurface,
			pmtrsValues[i].Created_By,
		)
	}
	return execErr
}
