package handler

import (
	"context"
	"encoding/json"
	"errors"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
	"expense-tracker/pkg/apperrors"
	"expense-tracker/pkg/response"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ExpenseHandler struct {
	Repo repository.ExpenseRepoInterface
}

func NewHandler(repo repository.ExpenseRepoInterface) *ExpenseHandler {
	return &ExpenseHandler{Repo: repo}
}

// =======GET=======
func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, err := ParseExpenseQuery(r)
	if err != nil {
		response.WriteBadRequest(w, err)
		return
	}
	log.Printf("Parsed params: %+v", params)
	filter := repository.ExpenseFilter{
		Min:    params.Min,
		Max:    params.Max,
		From:   params.From,
		To:     params.To,
		Search: params.Search,
	}

	offset := CalculateOffset(params.Page, params.Limit)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// execute query concurrently
	expensesChan := make(chan []model.Expense, 1)
	countChan := make(chan int, 1)
	errChan := make(chan error, 2)

	go func() {
		expense, err := h.Repo.Get(offset, params.Limit, filter)
		if err != nil {
			errChan <- fmt.Errorf("%w: %v", apperrors.ErrFailedToGetExpenses, err)
			return
		}
		expensesChan <- expense
	}()

	go func() {
		total, err := h.Repo.Count(filter)
		if err != nil {
			errChan <- fmt.Errorf("%w: %v", apperrors.ErrFailedToCount, err)
			return
		}
		countChan <- total
	}()

	// wait for result
	var expenses []model.Expense
	var total int

	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			log.Printf("database error: %v", err)
			response.WriteInternalServerError(w, err)
			return
		case expenses = <-expensesChan:
		case total = <-countChan:
		case <-ctx.Done():
			response.WriteError(w, apperrors.ErrTimeOut, http.StatusGatewayTimeout)
			return
		}
	}

	// build response
	responseData := model.CountRes{
		Data:       expenses,
		Page:       params.Page,
		TotalPages: CalculateTotalPage(total, params.Limit),
		Total:      total,
	}

	if err := response.WriteSuccess(w, responseData); err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

// =======GETbyID=======
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

// =======POST=======
func (h *ExpenseHandler) Post(w http.ResponseWriter, r *http.Request) {
	var e model.Expense

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	create, err := h.Repo.Post(&e)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
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

// =======PUT=======
func (h *ExpenseHandler) Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var e model.Expense

	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := h.Repo.Put(int64(id), &e)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

// =======PATCH=======
func (h *ExpenseHandler) Patch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var e model.Expense
	if err = json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := h.Repo.Patch(int64(id), &e)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// "I use pointers so that only the fields actually sent in the PATCH request get updated. Fields not included are left unchanged."

// =======DELETE=======
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(int64(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
