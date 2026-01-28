package handler

import (
	"database/sql"
	"encoding/json"
	"expense-tracker/internal/model"

	"net/http"
	"strconv"
)

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var e model.Expense

	err = h.DB.QueryRow(`
	SELECT id, date, amount, note
	FROM expenses
	WHERE id=$1`, id).Scan(
		&e.ID, &e.Date, &e.Amount, &e.Note)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
