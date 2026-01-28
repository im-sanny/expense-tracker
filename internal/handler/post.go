package handler

import (
	"encoding/json"
	"expense-tracker/internal/model"
	"net/http"
	"time"
)

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var e model.Expense
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	e.Date = time.Now()
	err = h.DB.QueryRow(`
	INSERT INTO
	expenses(date, amount, note)
	VALUES($1, $2, $3)
	RETURNING id, date, amount, note`,
		e.Date, e.Amount, e.Note).Scan(
		&e.ID, &e.Date, &e.Amount, &e.Note)

	if err != nil {
		http.Error(w, "failed to insert data in expense", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
