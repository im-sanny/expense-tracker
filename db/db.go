package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	conStr := "postgres://postgres:360420@localhost:5432/etracker?sslmode=disable"

	DB, err = sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("db connection successful")
}
