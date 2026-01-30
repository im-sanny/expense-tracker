package handler

import (
	"encoding/json"

	"net/http"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.Repo.Get()
	if err != nil {
		http.Error(w, "incomplete expense data", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(rows); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
