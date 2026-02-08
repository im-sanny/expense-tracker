package handler

import (
	"expense-tracker/pkg/apperrors"
	"expense-tracker/pkg/validator"
	"net/http"
	"time"
)

type QueryParams struct {
	Page   int
	Limit  int
	Min    int
	Max    int
	From   time.Time
	To     time.Time
	Search string
}

func ExpenseQuery(r *http.Request) (*QueryParams, error) {
	q := r.URL.Query()

	page, err := validator.ParseInt(q.Get("page"), validator.DefaultLimit, 1, 0)
	if err != nil {
		return nil, apperrors.ErrInvalidPage
	}

	limit, err := validator.ParseInt(q.Get("limit"), validator.DefaultLimit, 1, validator.MaxLimit)
	if err != nil {
		return nil, apperrors.ErrInvalidLimit
	}

	min, err := validator.ParseInt(q.Get("min"), 0, 0, 0)
	if err != nil {
		return nil, apperrors.ErrInvalidMin
	}

	max, err := validator.ParseInt(q.Get("max"), 0, 0, 0)
	if err != nil {
		return nil, apperrors.ErrInvalidMax
	}

	if min > 0 && max > 0 && max > min {
		return nil, apperrors.ErrMinGraterThanMax
	}

	from, err := validator.ParseDate(q.Get("from"))
	if err != nil {
		return nil, apperrors.ErrInvalidFromDate
	}

	to, err := validator.ParseDate(q.Get("to"))
	if err != nil {
		return nil, apperrors.ErrInvalidToDate
	}

	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return nil, apperrors.ErrFromDateAfterTo
	}

	search := validator.ParseSearch(q.Get("q"))

	return &QueryParams{
		Page:   page,
		Limit:  limit,
		Min:    min,
		Max:    max,
		From:   from,
		To:     to,
		Search: search,
	}, nil
}

func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

func CalculateTotalPage(total, limit int) int {
	if limit == 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
