package handler

import (
	"encoding/json"
	"expense-tracker/internal/model"
	"strconv"

	"net/http"
)

func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit // Page is for humans. Offset is for databases, Offset = how many rows to SKIP.

	rows, err := h.Repo.Get(offset, limit) // sql only only cares about OFFSET and LIMIT that's why I'm not sending page or anything.
	if err != nil {
		http.Error(w, "failed to get expense data", http.StatusInternalServerError)
		return
	}

	total, err := h.Repo.Count()
	if err != nil {
		http.Error(w, "failed to get total count", http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	res := model.CountRes{
		Data:       rows,
		Page:       page,
		TotalPages: totalPages,
		Total:      total,
	}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
