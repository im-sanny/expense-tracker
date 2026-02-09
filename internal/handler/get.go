package handler

import (
	"context"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
	"expense-tracker/pkg/apperrors"
	"expense-tracker/pkg/response"
	"fmt"
	"log"
	"time"

	"net/http"
)

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
