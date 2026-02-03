package handler

import (
	"encoding/json"
	"strconv"

	"net/http"
)

func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0{
		limit = 10
	}

	offset := (page - 1) * limit // what is offset? why i'm doing this?

	rows, err := h.Repo.Get(offset, limit) // passing offset and limit to repo, and why not passing page?
	if err != nil {
		http.Error(w, "incomplete expense data", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(rows); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
