package handler

import (
	"database/sql"
	"encoding/json"
	"expense-tracker/internal/model"
	"net/http"
	"strconv"
)

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var e model.Expense
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow(`
	UPDATE expenses SET
	amount = COALESCE($1, amount),
	note = COALESCE($2, note)
	WHERE id=$3 RETURNING id, date, amount, note`,
		e.Amount, e.Note, id).Scan(
		&e.ID, &e.Date, &e.Amount, &e.Note)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// "I use pointers so that only the fields actually sent in the PATCH request get updated. Fields not included are left unchanged."
