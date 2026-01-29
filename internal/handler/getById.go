package handler

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"

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

	gId, err := h.Repo.GetById(int64(id), &e)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(gId)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
