package handler

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/model"
	er "expense-tracker/pkg/errors"
	"log"
	"net/http"
)

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var e model.Expense

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	create, err := h.Repo.Post(&e)
	if err != nil {
		if errors.Is(err, er.ErrNotFound) {
			http.Error(w, "expense not created", http.StatusNoContent)
			return
		}
		log.Println(err)
		http.Error(w, "failed to insert data in expenses", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(create)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
