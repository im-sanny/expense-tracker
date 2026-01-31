package handler

import (
	"errors"
	er "expense-tracker/pkg/errors"
	"net/http"
	"strconv"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(int64(id))
	if err != nil {
		if errors.Is(err, er.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
