package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	// The Develompment connection
	connStr := "host=localhost port=5432 user=postgres password=admin123 dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{
		DB: db,
	}, nil
}

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
