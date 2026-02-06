package handler

import (
	"encoding/json"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
	"log"
	"strconv"
	"strings"
	"time"

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

	if min > 0 && max > 0 && min > max { //if min is greater than 0 and max is grater than 0 and min is grater than max then error
		http.Error(w, "min can't be grater than max", http.StatusBadRequest)
		return
	}

	layout := "2006-01-02"
	var from, to time.Time

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		var err error
		if from, err = time.Parse(layout, fromStr); err != nil {
			http.Error(w, "invalid 'from' data format", http.StatusBadRequest)
			return
		}
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		var err error
		if to, err = time.Parse(layout, toStr); err != nil {
			http.Error(w, "invalid 'to' data format", http.StatusBadRequest)
			return
		}
	}

	if !from.IsZero() && !to.IsZero() && from.After(to) {
		http.Error(w, "'from' date must be before 'to' date", http.StatusBadRequest)
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("q"))
	if search == "" {}

	offset := (page - 1) * limit

	filter := repository.ExpenseFilter{
		Min:    min,
		Max:    max,
		From:   from,
		To:     to,
		Search: search,
	}

	rows, err := h.Repo.Get(offset, limit, filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to get expense data", http.StatusInternalServerError)
		return
	}

	total, err := h.Repo.Count(filter)
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
