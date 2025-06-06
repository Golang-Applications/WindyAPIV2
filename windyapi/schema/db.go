package schema

import (
	"Windy-API/config"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func OpenDatabase(cfg config.Database) (*sql.DB, error) {
	db, err := openDatabase(cfg)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1049 {
		// database (catalog) does not exist - try creating it...
		return createDatabase(cfg)
	}
	return db, err
}

func openDatabase(cfg config.Database) (*sql.DB, error) {
	db, err := sql.Open("mysql", Dsn(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1049 {
		// database (catalog) does not exist - try creating it...
		return createDatabase(cfg)
	}
	return db, err
}

func createDatabase(cfg config.Database) (*sql.DB, error) {
	if db, err := sql.Open("mysql", Dsn(cfg.Host, cfg.Port, cfg.Username, cfg.Password, "")); err == nil {
		if _, err := db.Exec("CREATE SCHEMA " + cfg.Name); err == nil {
			if db, err := sql.Open("mysql", Dsn(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name)); err == nil {
				return db, db.Ping()
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func Dsn(host string, port int, username, password, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&multiStatements=true",
		username, password, host, port, dbName)
}
