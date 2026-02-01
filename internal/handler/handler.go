package handler

import (
	"expense-tracker/internal/repository"
)

type ExpenseHandler struct {
	Repo repository.ExpenseRepoInterface
}

func NewHandler(repo repository.ExpenseRepoInterface) *ExpenseHandler {
	return &ExpenseHandler{Repo: repo}
}
