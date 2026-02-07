package service

import (
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
)

type ExpenseService struct {
	Repo repository.ExpenseRepo
}

func NewExpenseService(repo repository.ExpenseRepo) *ExpenseService {
	return &ExpenseService{Repo: repo}
}

func (s *ExpenseService) List(page, limit int, f repository.ExpenseFilter) (*model.CountRes, error) {
	offset := (page - 1) * limit

	rows, err := s.Repo.Get(offset, limit, f)
	if err != nil {
		return nil, err
	}

	total, err := s.Repo.Count(f)
	if err != nil {
		return nil, err
	}

	return &model.CountRes{
		Data:       rows,
		Page:       page,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}, nil

}
