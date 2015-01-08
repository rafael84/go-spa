package database

import "database/sql"

const (
	connectionString = `user=postgres 
		dbname=postgres 
		password=1234 
		host=127.0.0.1
		port=15432
		sslmode=disable`
)

func Connect() (*sql.DB, error) {
	return sql.Open("postgres", connectionString)
}
