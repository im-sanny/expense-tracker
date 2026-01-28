package handler

import (
	"encoding/json"

	"expense-tracker/internal/model"
	"net/http"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query(`SELECT id, date, amount, note FROM expenses`)
	if err != nil {
		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var expense []model.Expense
	for rows.Next() {
		var e model.Expense
		err := rows.Scan(&e.ID, &e.Date, &e.Amount, &e.Note)
		if err != nil {
			http.Error(w, "failed to read expense data", http.StatusInternalServerError)
			return
		}
		expense = append(expense, e)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "incomplete expense data", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(expense); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
