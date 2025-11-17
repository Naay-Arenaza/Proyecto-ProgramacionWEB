package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	connStr  = "user=postgres password=12345 dbname=proyectos sslmode=disable"
)

func ConnectDB() (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, connStr)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
