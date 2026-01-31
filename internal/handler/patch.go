package handler

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/model"
	er "expense-tracker/pkg/errors"
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
	if err = json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := h.Repo.Patch(int64(id), &e)
	if err != nil {
		if errors.Is(err, er.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// "I use pointers so that only the fields actually sent in the PATCH request get updated. Fields not included are left unchanged."
