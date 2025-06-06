package Repository

import (
	"Windy-API/realtime_reports-obsolete/pkg/model"
	"database/sql"
)

type DatabaseRepo interface {
	GetDBConnection() *sql.DB
	GetWeatherStations(maxRecords int) (*[]model.WeatherForecast, error)
	InsertHourlyWeatherReport(sql string) (err error)
	CheckRowExistence() (int, error)
	InsertIntoDB(sql string) (err error)
	CheckForTableExistence(dbName string, tableName string) (int, error)
	InsertBatchStatements(sql string, pmtrsQuery string, pmtrsValues []model.RealtimeWeatherMapperToDB) (err error)
}
