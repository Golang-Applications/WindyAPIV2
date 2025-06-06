package persistence

import (
	"Windy-API/config"
	"database/sql"
)

/*
func NewPersistence(cfg config.Database, db *sql.DB) (RealtimeForecastPersistence, error) {
	return &persistence{
		cfg: cfg,
		db:  db,
	}, nil
}

*/

type Persistence struct {
	Cfg config.Database
	Db  *sql.DB
}
