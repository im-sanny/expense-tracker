package handler

import (
	"encoding/json"
	"errors"
	"expense-tracker/pkg/apperrors"

	"net/http"
	"strconv"
)

func (h *ExpenseHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	expense, err := h.Repo.GetById(int64(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(expense); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
