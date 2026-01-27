package handler

import (
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
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var e model.Expense

	err1 := h.DB.QueryRow(`
	SELECT id, date, amount, note
	FROM expenses
	WHERE id=$1`, id).Scan(&e.ID, &e.Date, &e.Amount, &e.Note)
	if err1 != nil {
		http.Error(w, "Failed to query", http.StatusInternalServerError)
		return
	}

	err2 := json.NewEncoder(w).Encode(e)
	if err2 != nil {
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
	}
}
