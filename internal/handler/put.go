package handler

import (
	"encoding/json"
	"expense-tracker/internal/model"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var e model.Expense
	e.Date = time.Now()
	err1 := json.NewDecoder(r.Body).Decode(&e)
	if err1 != nil {
		http.Error(w, "Failed to decode", http.StatusInternalServerError)
		return
	}

	err2 := h.DB.QueryRow(`
	UPDATE expenses SET
	amount=$1, note=$2 WHERE id=$3
	RETURNING id, date, amount, note`,
		e.Amount, e.Note, id).Scan(&e.ID, &e.Date, &e.Amount, &e.Note)

	if err2 != nil {
		http.Error(w, "Invalid query", http.StatusInternalServerError)
		return
	}

	err3 := json.NewEncoder(w).Encode(e)
	if err3 != nil {
		http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		return
	}

}
