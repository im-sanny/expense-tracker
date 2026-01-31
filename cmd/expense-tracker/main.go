package main

import (
	"expense-tracker/internal/db"
	"expense-tracker/internal/handler"
	"expense-tracker/internal/repository"
	"log"
	"net/http"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	repo := repository.NewExpenseRepo(database)
	h := handler.NewHandler(repo)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /track", h.Get)
	mux.HandleFunc("GET /track/{id}", h.GetById)
	mux.HandleFunc("POST /track", h.Post)
	mux.HandleFunc("PUT /track/{id}", h.Put)
	mux.HandleFunc("PATCH /track/{id}", h.Patch)
	mux.HandleFunc("DELETE /track/{id}", h.Delete)

	log.Println("Server running on port :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
