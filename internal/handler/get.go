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

	min := 0
	if minStr := r.URL.Query().Get("min"); minStr != "" {
		var err error
		if min, err = strconv.Atoi(minStr); err != nil {
			http.Error(w, "invalid min values", http.StatusBadRequest)
			return
		}
	}

	max := 0
	if maxStr := r.URL.Query().Get("max"); maxStr != "" {
		var err error
		if max, err = strconv.Atoi(maxStr); err != nil {
			http.Error(w, "invalid max value", http.StatusBadRequest)
			return
		}
	}

	if min > 0 && max > 0 && min > max {
		http.Error(w, "min can't be grater than max", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit

	rows, err := h.Repo.Get(offset, limit, min, max)
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
