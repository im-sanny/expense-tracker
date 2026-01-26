package main

import (
	"expense-tracker/db"
	"expense-tracker/handlers"
	"log"
	"net/http"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /track", handlers.Get)
	mux.HandleFunc("GET /track/{id}", handlers.GetById)
	mux.HandleFunc("POST /track", handlers.Post)
	mux.HandleFunc("PUT /track/{id}", handlers.Put)
	mux.HandleFunc("PATCH /track/{id}", handlers.Patch)
	mux.HandleFunc("DELETE /track/{id}", handlers.Delete)

	log.Println("Server running on port :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
