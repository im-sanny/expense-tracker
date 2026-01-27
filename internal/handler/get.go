package handler

import (
	"encoding/json"

	"expense-tracker/internal/model"
	"log"
	"net/http"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var et []model.Expense

	rows, err := h.DB.Query(`SELECT id, date, amount, note FROM expenses`)
	if err != nil {
		http.Error(w, "Failed to query expense", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var e model.Expense
		err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.Note)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		et = append(et, e)
	}

	err1 := json.NewEncoder(w).Encode(et)
	if err != nil {
		log.Println(err1)
	}
}
