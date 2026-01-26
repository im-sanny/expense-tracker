package main

import (
	"expense-tracker/db"
	"log"
	"net/http"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()

	log.Println("Server running on port :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
