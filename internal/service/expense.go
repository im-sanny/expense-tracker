package service

import (
	"context"
	"errors"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
	"expense-tracker/pkg/apperrors"
	"fmt"
	"time"
)

type ExpenseServiceInterface interface {
	Get(ctx context.Context, page, limit int, filter repository.ExpenseFilter) ([]model.Expense, error)
	GetById(ctx context.Context, id int64) (*model.Expense, error)
	Post(ctx context.Context, expense *model.Expense) (*model.Expense, error)
	Put(ctx context.Context, id int64, expense *model.Expense) (*model.Expense, error)
	Patch(ctx context.Context, id int64, expense *model.Expense) (*model.Expense, error)
	Delete(ctx context.Context, id int64) error
}

// service configuration
type ServiceConfig struct {
	MaxPage     int
	DefaultPage int
	TimeOut     time.Duration
}

// default configuration
var DefaultConfig = ServiceConfig{
	MaxPage:     100,
	DefaultPage: 1,
	TimeOut:     5 * time.Second,
}

type ExpenseService struct {
	repo   repository.ExpenseRepoInterface
	config ServiceConfig
}

func NewExpenseService(repo repository.ExpenseRepoInterface, cfg *ServiceConfig) *ExpenseService {
	if cfg == nil {
		cfg = &DefaultConfig
	}
	return &ExpenseService{repo: repo, config: *cfg}
}

func (s *ExpenseService) Get(ctx context.Context, page, limit int, filter repository.ExpenseFilter) (*model.CountRes, error) {
	if err := s.validatePagination(page, limit); err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, err)
	}

	offset := CalculateOffset(page, limit)

	ctx, cancel := context.WithTimeout(ctx, s.config.TimeOut)
	defer cancel()

	// execute query concurrently
	expensesChan := make(chan []model.Expense, 1)
	countChan := make(chan int, 1)
	errChan := make(chan error, 2)

	go func() {
		expense, err := s.repo.Get(ctx, offset, page, filter)
		if err != nil {
			errChan <- fmt.Errorf("%w: %v", apperrors.ErrFailedToGetExpenses, err)
			return
		}
		expensesChan <- expense
	}()

	go func() {
		total, err := s.repo.Count(filter)
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
			return nil, err
		case expenses = <-expensesChan:
		case total = <-countChan:
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %v", apperrors.ErrTimeOut, ctx.Err())
		}
	}

	// build response
	return &model.CountRes{
		Data:       expenses,
		Page:       page,
		TotalPages: CalculateTotalPage(total, limit),
		Total:      total,
	}, nil

}

func (s *ExpenseService) GetById(ctx context.Context, id int64) (*model.Expense, error) {
	return s.repo.GetById(ctx, id)
}

func (s *ExpenseService) Post(ctx context.Context, expense *model.Expense) (*model.Expense, error) {
	return s.repo.Post(ctx, expense)
}

func (s *ExpenseService) Put(ctx context.Context, id int64, expense *model.Expense) (*model.Expense, error) {
	return s.repo.Put(ctx, id, expense)
}

func (s *ExpenseService) Patch(ctx context.Context, id int64, expense *model.Expense) (*model.Expense, error) {
	return s.repo.Patch(ctx, id, expense)
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ExpenseService) validatePagination(page, limit int) error {
	if page < 1 {
		return errors.New("page must be at least 1")
	}

	if limit < 1 {
		return errors.New("limit must be at least 1")
	}

	if limit > s.config.MaxPage {
		return fmt.Errorf("limit cannot be exceed %d", s.config.MaxPage)
	}

	return nil
}

func (s *ExpenseService) validateExpense(expense *model.Expense) error {
	if expense == nil {
		return errors.New("expense cannot be nil")
	}

	if *expense.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if expense.Date.IsZero() {
		return errors.New("date cannot be zero")
	}

	if expense.Date.After(time.Now()) {
		return errors.New("date cannot be in the future")
	}

	if len(*expense.Note) > 500 { // Example validation
		return errors.New("note cannot exceed 500 characters")
	}

	return nil
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
