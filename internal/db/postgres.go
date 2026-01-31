package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	conStr := "postgres://postgres:360420@localhost:5432/etracker?sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("db connection successful")
	return db
}
