package handler

import (
	"encoding/json"
	"expense-tracker/internal/model"
	"log"
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
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit // what is offset? why i'm doing this?

	rows, err := h.Repo.Get(offset, limit) // passing offset and limit to repo, and why not passing page?
	if err != nil {
		http.Error(w, "incomplete expense data", http.StatusInternalServerError)
		return
	}

	total, err := h.Repo.Count()
	if err != nil {
		log.Println(err)
		return
	}
	totalPages := (total + limit - 1) / limit

	type ExpenseCount struct {
		Data       []model.Expense `json:"data"`
		Page       int             `json:"page"`
		TotalPages int             `json:"total_pages"`
		Total      int             `json:"total"`
	}

	res := ExpenseCount{
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
