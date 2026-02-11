package service

import (
	"context"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
	"time"
)

type ExpenseServiceInterface interface {
	Get(ctx context.Context, page, limit int, filter repository.ExpenseFilter) ([]model.Expense, error)
	GetById(ctx context.Context, id int64) (*model.Expense, error)
	Post(ctx context.Context, expense *model.Expense) (*model.Expense, error)
	Put(ctx context.Context, id int64, expense *model.Expense) (*model.Expense, error)
	Patch(ctx context.Context, id int64, expense *model.Expense) (*model.Expense, error)
	Delete(ctx context.Context, id int64, expense *model.Expense) error
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
