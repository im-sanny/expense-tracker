package handler

import (
	"expense-tracker/internal/repository"
)

type Handler struct {
	Repo *repository.ExpenseRepo
}

func NewHandler(repo *repository.ExpenseRepo) *Handler {
	return &Handler{Repo: repo}
}
